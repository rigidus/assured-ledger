// Copyright 2020 Insolar Network Ltd.
// All rights reserved.
// This material is licensed under the Insolar License version 1.0,
// available at https://github.com/insolar/assured-ledger/blob/master/LICENSE.md.

package virtual

import (
	"context"
	"io"

	"github.com/ThreeDotsLabs/watermill"
	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/ThreeDotsLabs/watermill/pubsub/gochannel"

	"github.com/insolar/assured-ledger/ledger-core/v2/application/testwalletapi"
	"github.com/insolar/assured-ledger/ledger-core/v2/network/messagesender"
	"github.com/insolar/assured-ledger/ledger-core/v2/runner"

	"github.com/insolar/component-manager"

	"github.com/insolar/assured-ledger/ledger-core/v2/application/api"
	"github.com/insolar/assured-ledger/ledger-core/v2/certificate"
	"github.com/insolar/assured-ledger/ledger-core/v2/configuration"
	"github.com/insolar/assured-ledger/ledger-core/v2/contractrequester"
	"github.com/insolar/assured-ledger/ledger-core/v2/cryptography"
	"github.com/insolar/assured-ledger/ledger-core/v2/insolar"
	"github.com/insolar/assured-ledger/ledger-core/v2/insolar/bus"
	"github.com/insolar/assured-ledger/ledger-core/v2/insolar/jet"
	"github.com/insolar/assured-ledger/ledger-core/v2/insolar/jetcoordinator"
	"github.com/insolar/assured-ledger/ledger-core/v2/insolar/node"
	"github.com/insolar/assured-ledger/ledger-core/v2/insolar/pulse"
	"github.com/insolar/assured-ledger/ledger-core/v2/instrumentation/inslogger"
	"github.com/insolar/assured-ledger/ledger-core/v2/instrumentation/inslogger/logwatermill"
	"github.com/insolar/assured-ledger/ledger-core/v2/keystore"
	"github.com/insolar/assured-ledger/ledger-core/v2/logicrunner/artifacts"
	"github.com/insolar/assured-ledger/ledger-core/v2/logicrunner/pulsemanager"
	"github.com/insolar/assured-ledger/ledger-core/v2/metrics"
	"github.com/insolar/assured-ledger/ledger-core/v2/network/servicenetwork"
	"github.com/insolar/assured-ledger/ledger-core/v2/platformpolicy"
	"github.com/insolar/assured-ledger/ledger-core/v2/server/internal"
	"github.com/insolar/assured-ledger/ledger-core/v2/virtual"
)

type bootstrapComponents struct {
	CryptographyService        insolar.CryptographyService
	PlatformCryptographyScheme insolar.PlatformCryptographyScheme
	KeyStore                   insolar.KeyStore
	KeyProcessor               insolar.KeyProcessor
}

func initBootstrapComponents(ctx context.Context, cfg configuration.Configuration) bootstrapComponents {
	earlyComponents := component.NewManager(nil)

	keyStore, err := keystore.NewKeyStore(cfg.KeysPath)
	checkError(ctx, err, "failed to load KeyStore: ")

	platformCryptographyScheme := platformpolicy.NewPlatformCryptographyScheme()
	keyProcessor := platformpolicy.NewKeyProcessor()

	cryptographyService := cryptography.NewCryptographyService()
	earlyComponents.Register(platformCryptographyScheme, keyStore)
	earlyComponents.Inject(cryptographyService, keyProcessor)

	return bootstrapComponents{
		CryptographyService:        cryptographyService,
		PlatformCryptographyScheme: platformCryptographyScheme,
		KeyStore:                   keyStore,
		KeyProcessor:               keyProcessor,
	}
}

func initCertificateManager(
	ctx context.Context,
	cfg configuration.Configuration,
	cryptographyService insolar.CryptographyService,
	keyProcessor insolar.KeyProcessor,
) *certificate.CertificateManager {
	var certManager *certificate.CertificateManager
	var err error

	publicKey, err := cryptographyService.GetPublicKey()
	checkError(ctx, err, "failed to retrieve node public key")

	certManager, err = certificate.NewManagerReadCertificate(publicKey, keyProcessor, cfg.CertificatePath)
	checkError(ctx, err, "failed to start Certificate")

	return certManager
}

// initComponents creates and links all insolard components
func initComponents(
	ctx context.Context,
	cfg configuration.Configuration,
	cryptographyService insolar.CryptographyService,
	pcs insolar.PlatformCryptographyScheme,
	keyStore insolar.KeyStore,
	keyProcessor insolar.KeyProcessor,
	certManager insolar.CertificateManager,

) (*component.Manager, func()) {
	cm := component.NewManager(nil)

	// Watermill.
	var (
		wmLogger   *logwatermill.WatermillLogAdapter
		publisher  message.Publisher
		subscriber message.Subscriber
	)
	{
		wmLogger = logwatermill.NewWatermillLogAdapter(inslogger.FromContext(ctx))
		pubsub := gochannel.NewGoChannel(gochannel.Config{}, wmLogger)
		subscriber = pubsub
		publisher = pubsub
		// Wrapped watermill Publisher for introspection.
		publisher = internal.PublisherWrapper(ctx, cm, cfg.Introspection, publisher)
	}

	nw, err := servicenetwork.NewServiceNetwork(cfg, cm)
	checkError(ctx, err, "failed to start Network")

	metricsComp := metrics.NewMetrics(cfg.Metrics, metrics.GetInsolarRegistry("virtual"), "virtual")

	jc := jetcoordinator.NewJetCoordinator(cfg.Ledger.LightChainLimit, *certManager.GetCertificate().GetNodeRef())
	pulses := pulse.NewStorageMem()

	b := bus.NewBus(cfg.Bus, publisher, pulses, jc, pcs)
	artifactsClient := artifacts.NewClient(b)
	cachedPulses := artifacts.NewPulseAccessorLRU(pulses, artifactsClient, cfg.LogicRunner.PulseLRUSize)

	messageSender := messagesender.NewDefaultService(publisher, jc, pulses)

	runnerService := runner.NewService()
	err = runnerService.Init()
	checkError(ctx, err, "failed to initialize Runner Service")

	virtualDispatcher := virtual.NewDispatcher()
	virtualDispatcher.Runner = runnerService
	virtualDispatcher.MessageSender = messageSender

	contractRequester, err := contractrequester.New(
		b,
		pulses,
		jc,
		pcs,
	)
	checkError(ctx, err, "failed to start ContractRequester")

	availabilityChecker := api.NewNetworkChecker(cfg.AvailabilityChecker)

	API, err := api.NewRunner(
		&cfg.APIRunner,
		certManager,
		contractRequester,
		nw,
		nw,
		pulses,
		artifactsClient,
		jc,
		nw,
		availabilityChecker,
	)
	checkError(ctx, err, "failed to start ApiRunner")

	AdminAPIRunner, err := api.NewRunner(
		&cfg.AdminAPIRunner,
		certManager,
		contractRequester,
		nw,
		nw,
		pulses,
		artifactsClient,
		jc,
		nw,
		availabilityChecker,
	)
	checkError(ctx, err, "failed to start AdminAPIRunner")

	APIWrapper := api.NewWrapper(API, AdminAPIRunner)

	pm := pulsemanager.NewPulseManager()

	cm.Register(
		pcs,
		keyStore,
		cryptographyService,
		keyProcessor,
		certManager,
		virtualDispatcher,
		runnerService,
		APIWrapper,
		testwalletapi.NewTestWalletServer(cfg.TestWalletAPI, virtualDispatcher, pulses),
		availabilityChecker,
		nw,
		pm,
		cachedPulses,
	)

	components := []interface{}{
		b,
		publisher,
		contractRequester,
		artifactsClient,
		jc,
		pulses,

		jet.NewStore(),
		node.NewStorage(),
	}
	components = append(components, []interface{}{
		metricsComp,
		cryptographyService,
		keyProcessor,
	}...)

	cm.Inject(components...)

	err = cm.Init(ctx)
	checkError(ctx, err, "failed to init components")

	// this should be done after Init due to inject
	pm.AddDispatcher(virtualDispatcher.FlowDispatcher)

	return cm, startWatermill(
		ctx, wmLogger, subscriber, b,
		nw.SendMessageHandler,
		virtualDispatcher.FlowDispatcher.Process,
	)
}

func startWatermill(
	ctx context.Context,
	logger watermill.LoggerAdapter,
	sub message.Subscriber,
	b *bus.Bus,
	outHandler, inHandler message.NoPublishHandlerFunc,
) func() {
	inRouter, err := message.NewRouter(message.RouterConfig{}, logger)
	if err != nil {
		panic(err)
	}
	outRouter, err := message.NewRouter(message.RouterConfig{}, logger)
	if err != nil {
		panic(err)
	}

	outRouter.AddNoPublisherHandler(
		"OutgoingHandler",
		bus.TopicOutgoing,
		sub,
		outHandler,
	)

	inRouter.AddMiddleware(
		b.IncomingMessageRouter,
	)

	inRouter.AddNoPublisherHandler(
		"IncomingHandler",
		bus.TopicIncoming,
		sub,
		inHandler,
	)
	startRouter(ctx, inRouter)
	startRouter(ctx, outRouter)

	return stopWatermill(ctx, inRouter, outRouter)
}

func stopWatermill(ctx context.Context, routers ...io.Closer) func() {
	return func() {
		for _, r := range routers {
			err := r.Close()
			if err != nil {
				inslogger.FromContext(ctx).Error("Error while closing router", err)
			}
		}
	}
}

func startRouter(ctx context.Context, router *message.Router) {
	go func() {
		if err := router.Run(ctx); err != nil {
			inslogger.FromContext(ctx).Error("Error while running router", err)
		}
	}()
	<-router.Running()
}
