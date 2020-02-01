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
	"fmt"
	"github.com/pkg/errors"
	"go/ast"
	"go/parser"
	"go/token"
	"io/ioutil"
	"os"
	"reflect"
	"strings"
)

const (
	TemplateDirectory = "templates"

	filename = "src/github.com/insolar/assured-ledger/ledger-core/v2/conveyor/smachine/ping-pong/example/example_3.go"
	mainPkg   = "main"
	errorType = "error"
	MachineTypeGoPlugin
)

type RecvPair struct {
	Name 	string
	Type 	string
}

type Variant struct {
	Type    int
	Obj 	string
	Fun		string
	Str		string // string representation
}

const (
	MistakeType  = 1 + iota

	// type-ids for Variant struct
	SelectorType
	StringType

	// type-ids for Ret struct
	RetTypeCall
)

func (v Variant) Show() string {
	switch v.Type {
	case MistakeType:
		return fmt.Sprintf("MistakeType")
	case SelectorType:
		return fmt.Sprintf("(. %s %s)", v.Obj, v.Fun)
	case StringType:
		return fmt.Sprintf("%s", v.Str)
	default:
		return "Impossible Error"
	}
}

type Ret struct {
	Lvl  	string
	Str  	string
	Type 	int
	Sel  	Variant
	Args 	[]Variant
}

type FnState struct {
	Name 	string
	Recv    *RecvPair
	Pars    map[string]string // k:name, v:type
	Rets    []*Ret
}

// ParsedFile struct with prepared info we extract from source code
type ParsedFile struct {
	filename	string
	code        []byte
	fileSet     *token.FileSet
	node        *ast.File
	states      map[string]*FnState
}

func main() {
	dbg := true
	pathname:= fmt.Sprintf("%s/%s", os.Getenv("GOPATH"), filename)
	pf := ParseFile(pathname)
	uml := "@startuml"
	// Debug output
	if dbg {
		fmt.Printf("\n:: resource filename: %s", pf.filename)
	}
	for _, state := range pf.states {
		if dbg {
			fmt.Printf("\n\nfn: %s", state.Name) // Function name
			fmt.Printf("\nrecv: %s | %s", state.Recv.Name, state.Recv.Type) // Receiver
			for parName, parType := range state.Pars { // Parameters
				fmt.Printf("\npar name: %s | type: %s", parName, parType)
			}
		}
		for _, item := range state.Rets {
			if dbg {
				fmt.Printf("\n%s: ['%s']", item.Lvl, item.Str)
			}
			switch item.Type {
			case RetTypeCall:
				switch item.Sel.Fun {
				case "Stop":
					uml += fmt.Sprintf("\n%s --> [*]", state.Name)
				case "Jump":
					if 1 == len(state.Rets) && "Top" == item.Lvl { // One Top Level Jmp -> Init
						uml += fmt.Sprintf("\n[*] --> %s : %s", item.Args[0].Fun, state.Name)
					} else {
						uml += fmt.Sprintf("\n%s --> %s", state.Name, item.Args[0].Fun)
				    }
				case "RepeatOrJumpElse":
					uml += fmt.Sprintf("\n%s -[#RoyalBlue]-> %s : RepeatOr(Jump)Else", state.Name, item.Args[2].Fun)
					uml += fmt.Sprintf("\n%s -[#DarkGreen]-> %s : RepeatOrJump(Else)", state.Name, item.Args[3].Fun)
				default:
					fmt.Printf("\n(=> (. %s %s)", item.Sel.Obj, item.Sel.Fun)
					for _, arg := range item.Args {
						fmt.Printf("\n       %s", arg.Show())
					}
					fmt.Printf(")" )
				}
			default:
				fmt.Printf("\nError: Unknown RetType: %d", item.Type)
			}
			if dbg {
				fmt.Printf("\n(=> (. %s %s)", item.Sel.Obj, item.Sel.Fun)
				for _, arg := range item.Args {
					fmt.Printf("\n       %s", arg.Show())
				}
				fmt.Printf(")")
			}
		}
	}
	uml += "\n@enduml"
	fmt.Printf("\n\n\n\n\n~~~~~~~~~~~~~~~~~\n%s", uml)
}

// ParseFile parses a file as Go source code of a smart contract
// and returns it as `ParsedFile`
func ParseFile(fileName string) *ParsedFile {
	pf := &ParsedFile{
		filename:        fileName,
	}

	sourceCode, err := slurpFile(fileName)
	if err != nil {
		return nil
	}
	pf.code = sourceCode

	pf.fileSet = token.NewFileSet()
	node, err := parser.ParseFile(pf.fileSet, pf.filename, pf.code, parser.ParseComments)
	if err != nil {
		return nil
	}
	pf.node = node

	pf.states = make(map[string]*FnState)

	ast.Inspect(node, func(n ast.Node) bool {
		fn, ok := n.(*ast.FuncDecl)
		if ok {
			pf.parseMethod(fn)  // TODO: return value
		}
		return true
	})

	return pf
}

func (pf *ParsedFile) parseMethod(fn *ast.FuncDecl) {
	// I want to analise only method functions
	for _, fld := range fn.Recv.List {

		// I want analyse only method-functions
		if 1 != len(fld.Names) { // There is method function
			continue
		}

		// I want analyse only exported methods
		if !fn.Name.IsExported() {
			continue
		}

		// Receiver
		recv := &RecvPair{
			Name : fld.Names[0].Name,
			Type : fmt.Sprintf("%s", pf.code[fld.Type.Pos()-1:fld.Type.End()-1]),
		}

		// Parameters
		pars := make(map[string]string, 0)
		for _, par := range fn.Type.Params.List {
			pars[par.Names[0].Name] = fmt.Sprintf("%s", pf.code[par.Type.Pos()-1:par.Type.End()-1])
		}

		// I want to analyse only methods, who takes context
		if !isMethodTakesCtx(pars) {
			continue
		}

		// I want to analyze methods which have a `smashine.StateUpdate' result type
		res := fn.Type.Results.List[0].Type
		resSel, ok := res.(*ast.SelectorExpr)
		if !ok || "StateUpdate" != resSel.Sel.Name {
			continue
		}
		resXstr := fmt.Sprintf("%s", pf.code[resSel.X.Pos()-1:resSel.X.End()-1])
		if "smachine" != resXstr {
			continue
		}

		// Find all Return Statements in function content
		var rets = make([]*Ret, 0)
		for _, smth := range fn.Body.List { // ∀ fn.Body.List ← (or RetStmt (Inspect ...))
			retStmt, ok := smth.(*ast.ReturnStmt)
			if ok {
				// return from top-level statements of function
				rets = append(rets, collectRets(retStmt, pf.code, "Top")...)
			} else {
				ast.Inspect(smth, func(in ast.Node) bool {
					// Find Return Statements
					retStmt, ok := in.(*ast.ReturnStmt) // ←
					if ok {
						// return from deep-level function statememt
						rets = append(rets, collectRets(retStmt, pf.code, "Deep")...)
					} else {
						//fmt.Printf("\nin: %s", reflect.TypeOf(in))
					}
					return true
				})
			}
		}

		pf.states[fn.Name.Name] = &FnState{
			Name: fn.Name.Name,
			Recv: recv,
			Pars: pars,
			Rets: rets,
		}
	}
}

func collectRets(retStmt *ast.ReturnStmt, code []byte, level string) []*Ret {
	var acc []*Ret
	for _, ret := range retStmt.Results {
		item := &Ret{
			Lvl:	level,
			Str:	fmt.Sprintf("%s",code[ret.Pos()-1:ret.End()-1]),
		}
		for _, retNode := range retStmt.Results {
			switch retNode.(type) {
			case *ast.CallExpr:
				item.Type = RetTypeCall
				retCall := retNode.(*ast.CallExpr)
				switch retCall.Fun.(type) {
				case *ast.SelectorExpr:
					retSelector := retCall.Fun.(*ast.SelectorExpr)
					retX, ok := retSelector.X.(*ast.Ident)
					if ok {
						item.Sel.Obj = retX.Name
					} else { // may be not Ident !
						fmt.Printf("\nERR: UNKNOWN RETSELECTOR %s", reflect.TypeOf(retSelector.X))
					}
					item.Sel.Fun = retSelector.Sel.Name
					// Args
					accArgs := make([]Variant, 0)
					for _, retarg := range retCall.Args {
						switch retarg.(type) {
						case *ast.SelectorExpr:
							sel := retarg.(*ast.SelectorExpr)
							selName := fmt.Sprintf("%s", code[sel.X.Pos()-1:sel.X.End()-1])
							// arg := fmt.Sprintf("(. %s %s)", selName, sel.Sel.Name) // ::::::::::::::::::::::::::::::::::
							arg := Variant{
								Type: SelectorType,
								Obj: selName,
								Fun: sel.Sel.Name,
							}
							accArgs = append(accArgs, arg)
						case *ast.Ident:
							idn := retarg.(*ast.Ident)
							//arg := fmt.Sprintf("%s", idn.Name)
							arg := Variant{
								Type: StringType,
								Str:  idn.Name,
							}
							accArgs = append(accArgs, arg)
						default:
							fmt.Printf("\nERR: UNKNOWN RETARG('%s' of %s)", retarg, reflect.TypeOf(retarg))
						}
					}
					//item.Args = append(item.Args, accArgs)
					item.Args = accArgs
				default:
					fmt.Printf("\nERR: UNKNOWN RETSEL %s", fmt.Sprintf("%s", reflect.TypeOf(retCall.Fun)))
				}
			default:
				fmt.Printf("\nERR: UNKNOWN RETNODE %s", fmt.Sprintf("%s", reflect.TypeOf(retNode)))
			}
		}
		acc = append(acc, item)
	}
	return acc
}

func isMethodTakesCtx(pars map[string]string) bool {
	for _, parType := range pars {
		if strings.Contains(parType, "Context") {
			return true
		}
	}
	return false
}

func slurpFile(fileName string) ([]byte, error) {
	file, err := os.OpenFile(fileName, os.O_RDONLY, 0)
	if err != nil {
		return nil, errors.Wrap(err, "Can't open file '"+fileName+"'")
	}
	defer file.Close() //nolint: errcheck

	res, err := ioutil.ReadAll(file)
	if err != nil {
		return nil, errors.Wrap(err, "Can't read file '"+fileName+"'")
	}
	return res, nil
}
