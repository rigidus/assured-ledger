///
//    Copyright 2019 Insolar Technologies
//
//    Licensed under the Apache License, Version 2.0 (the "License");
//    you may not use this file except in compliance with the License.
//    You may obtain a copy of the License at
//
//        http://www.apache.org/licenses/LICENSE-2.0
//
//    Unless required by applicable law or agreed to in writing, software
//    distributed under the License is distributed on an "AS IS" BASIS,
//    WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
//    See the License for the specific language governing permissions and
//    limitations under the License.
///

package main

import (
	"context"
	"fmt"
	"math/rand"
	"time"

	"github.com/insolar/assured-ledger/ledger-core/v2/conveyor/smachine"
	"github.com/insolar/assured-ledger/ledger-core/v2/conveyor/smachine/ping-pong/example"
	"github.com/insolar/assured-ledger/ledger-core/v2/conveyor/sworker"
	"github.com/insolar/assured-ledger/ledger-core/v2/vanilla/synckit"
)

func main() {
	const scanCountLimit = 1e6

	signal := synckit.NewVersionedSignal()
	sm := smachine.NewSlotMachine(smachine.SlotMachineConfig{
		SlotPageSize:         1000,
		PollingPeriod:        10 * time.Millisecond,
		PollingTruncate:      1 * time.Microsecond,
		BoostNewSlotDuration: 10 * time.Millisecond,
		ScanCountLimit:       scanCountLimit,
	}, signal.NextBroadcast, signal.NextBroadcast, nil)

	sm.PutDependency("example.ServiceAdapterA", example.CreateServiceAdapterA())
	sm.PutDependency("example.catalogC", example.CreateCatalogC())
	sm.PutDependency("example.catalogD", example.CreateCatalogD())

	rand.Seed(time.Now().UnixNano())

	//i := 0
	for i := 0; i < 2; i++ {
		fmt.Printf(":: add new sm :: %d | %v>\n", i,
			sm.AddNew(context.Background(), &example.StateMachine3{}, smachine.CreateDefaultValues{}))
	}

	workerFactory := sworker.NewAttachableSimpleSlotWorker()

	startNano := time.Now().UnixNano()
	startBase := example.IterationCount

	iterBase := example.IterationCount
	iterStart := time.Now().UnixNano()

	neverSignal := synckit.NewNeverSignal()

	prev := 0
	for i := 0; ; i++ {
		var (
			repeatNow    bool
			nextPollTime time.Time
		)
		wakeupSignal := signal.Mark()
		workerFactory.AttachTo(sm, neverSignal, scanCountLimit, func(worker smachine.AttachedSlotWorker) {
			repeatNow, nextPollTime = sm.ScanOnce(0, worker)
		})

		//fmt.Printf("%03d %v ================================== slots=%v of %v\n", i, time.Now(), sm.OccupiedSlotCount(), sm.AllocatedSlotCount())
		//fmt.Printf("%03d %v =============== repeatNow=%v nextPollTime=%v\n", i, time.Now(), repeatNow, nextPollTime)

		if i >= prev+1e8/scanCountLimit {
			prev = i
			fmt.Printf("%03d ================================== iters=%d speedIter=%7.0f speedTotal=%7.0f\n", i, example.IterationCount,
				float64(example.IterationCount-iterBase)/(float64(time.Now().UnixNano()-iterStart)/float64(time.Second)),
				float64(example.IterationCount-startBase)/(float64(time.Now().UnixNano()-startNano)/float64(time.Second)),
			)
			fmt.Printf("%03d %v ================================== slots=%v of %v\n", i, time.Now(), sm.OccupiedSlotCount(), sm.AllocatedSlotCount())
			active, inactive := example.Limiter.GetCounts()
			fmt.Println("Limiter: ", active, inactive)

			iterBase = example.IterationCount
			iterStart = time.Now().UnixNano()
			if example.IterationCount < 1000000 {
				startNano = time.Now().UnixNano()
				startBase = example.IterationCount
			}
			//sm.Cleanup(worker)
			//fmt.Printf("%03d %v ================================== slots=%v of %v\n", i, time.Now(), sm.OccupiedSlotCount(), sm.AllocatedSlotCount())
		}

		if repeatNow {
			continue
		}
		if !nextPollTime.IsZero() {
			time.Sleep(time.Until(nextPollTime))
		} else {
			wakeupSignal.Wait()
			//time.Sleep(3 * time.Second)
		}
	}
}
