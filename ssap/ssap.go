package ssap

import (
	"fmt"

	"github.com/1lann/golua/lang"
	"golang.org/x/tools/go/ssa"
)

func toExp(val ssa.Value) lang.Expression {
	cnst, ok := val.(*ssa.Const)
	if ok {
		return lang.Expression{
			Type:  lang.Type{cnst.Type()},
			Value: cnst.Value.ExactString(),
		}
	}

	return lang.Expression{
		Type: lang.Type{val.Type()},
		Variable: &lang.Variable{
			Name: val.Name(),
			Type: lang.Type{val.Type()},
		},
	}
}

func toVar(val ssa.Value) lang.Variable {
	return lang.Variable{
		Name: val.Name(),
		Type: lang.Type{val.Type()},
	}
}

func Process(pkg *ssa.Package, language lang.Lang) {
	for _, mem := range pkg.Members {
		f, ok := mem.(*ssa.Function)
		if !ok {
			continue
		}

		inputs := f.Signature.Params()
		outputs := f.Signature.Results()

		passedInput := make([]lang.Variable, inputs.Len())
		passedOutput := make([]lang.Variable, outputs.Len())

		for i := 0; i < inputs.Len(); i++ {
			input := inputs.At(i)
			passedInput[i] = lang.Variable{Name: input.Name(), Type: lang.Type{input.Type()}}
		}

		for i := 0; i < outputs.Len(); i++ {
			output := outputs.At(i)
			passedOutput[i] = lang.Variable{Name: output.Name(), Type: lang.Type{output.Type()}}
		}

		lf := language.StartFunction(f.Name(), passedInput, passedOutput)
		_ = lf

		for _, b := range f.Blocks {
			lb := lf.StartBlock(b.Index)

			extractions := make(map[string][]*ssa.Extract)

			for _, v := range b.Instrs {
				switch v := v.(type) {
				case *ssa.Extract:
					extractions[v.Tuple.Name()] = append(extractions[v.Tuple.Name()], v)
				}
			}

			for _, instr := range b.Instrs {
				switch instr := instr.(type) {
				case *ssa.DebugRef:
					// no-op
				case *ssa.UnOp:
					// lb.EmitCall(instr.Op.String(), []lang.Expression{
					// 	toExp(instr.X),
					// }, []lang.Variable{toVar(instr)})
				case *ssa.BinOp:
					lb.EmitCall(instr.Op.String(), []lang.Expression{
						toExp(instr.X),
						toExp(instr.Y),
					}, []lang.Variable{toVar(instr)})
				case *ssa.Call:
					switch f := instr.Call.Value.(type) {
					case *ssa.Builtin:
						inputs := make([]lang.Expression, len(instr.Call.Args))

						for i, arg := range instr.Call.Args {
							inputs[i] = toExp(arg)
						}

						outputs := make([]lang.Variable, len(extractions[f.Name()]))

						for _, ext := range extractions[f.Name()] {
							outputs[ext.Index] = toVar(ext)
						}

						lb.EmitCall(f.Name(), inputs, outputs)
					case *ssa.Function:
						inputs := make([]lang.Expression, len(instr.Call.Args))

						for i, arg := range instr.Call.Args {
							inputs[i] = toExp(arg)
						}

						var outputs []lang.Variable
						if _, found := extractions[f.Name()]; found {
							outputs = make([]lang.Variable, len(extractions[f.Name()]))

							for _, ext := range extractions[f.Name()] {
								outputs[ext.Index] = toVar(ext)
							}
						} else {
							sigOutputs := instr.Call.Signature().Results()
							outputs = make([]lang.Variable, sigOutputs.Len())

							for i := 0; i < sigOutputs.Len(); i++ {
								output := sigOutputs.At(i)
								outputs[i] = lang.Variable{Name: instr.Name(), Type: lang.Type{output.Type()}}
							}
						}

						lb.EmitCall(f.Name(), inputs, outputs)
					default:
						panic("unsupported type")
					}

					// ssaCallFunc, ok := instr.Call.Value.(*ssa.Function)
					// if !ok {
					// 	// continue
					// 	fmt.Println(reflect.TypeOf(instr.Call.Value))
					// 	panic("what")
					// }
				case *ssa.ChangeInterface:

				case *ssa.ChangeType:

				case *ssa.Convert:

				case *ssa.MakeInterface:

				case *ssa.Extract:

				case *ssa.Slice:

				case *ssa.Return:
					returns := make([]lang.Expression, len(instr.Results))
					for i, v := range instr.Results {
						returns[i] = toExp(v)
					}

					lb.EmitReturn(returns)

				case *ssa.RunDefers:

				case *ssa.Panic:

				case *ssa.Send:

				case *ssa.Store:

				case *ssa.If:
					// succ := 1
					// if fr.get(instr.Cond).(bool) {
					// 	succ = 0
					// }
					// fr.prevBlock, fr.block = fr.block, fr.block.Succs[succ]
					// return kJump

					lb.EmitIf(toExp(instr.Cond), instr.Block().Succs[0].Index, instr.Block().Succs[1].Index)

				case *ssa.Jump:
					lb.EmitJump(instr.Block().Succs[0].Index)

					// fr.prevBlock, fr.block = fr.block, fr.block.Succs[0]
					// return kJump

				case *ssa.Defer:
					// fn, args := prepareCall(fr, &instr.Call)
					// fr.defers = &deferred{
					// 	fn:    fn,
					// 	args:  args,
					// 	instr: instr,
					// 	tail:  fr.defers,
					// }

				case *ssa.Go:
					// fn, args := prepareCall(fr, &instr.Call)
					// atomic.AddInt32(&fr.i.goroutines, 1)
					// go func() {
					// 	call(fr.i, nil, instr.Pos(), fn, args)
					// 	atomic.AddInt32(&fr.i.goroutines, -1)
					// }()

				case *ssa.MakeChan:
					// fr.env[instr] = make(chan value, asInt(fr.get(instr.Size)))

				case *ssa.Alloc:
					// var addr *value
					// if instr.Heap {
					// 	// new
					// 	addr = new(value)
					// 	fr.env[instr] = addr
					// } else {
					// 	// local
					// 	addr = fr.env[instr].(*value)
					// }
					// *addr = zero(deref(instr.Type()))

				case *ssa.MakeSlice:
					// slice := make([]value, asInt(fr.get(instr.Cap)))
					// tElt := instr.Type().Underlying().(*types.Slice).Elem()
					// for i := range slice {
					// 	slice[i] = zero(tElt)
					// }
					// fr.env[instr] = slice[:asInt(fr.get(instr.Len))]

				case *ssa.MakeMap:
					// reserve := 0
					// if instr.Reserve != nil {
					// 	reserve = asInt(fr.get(instr.Reserve))
					// }
					// fr.env[instr] = makeMap(instr.Type().Underlying().(*types.Map).Key(), reserve)

				case *ssa.Range:
					// fr.env[instr] = rangeIter(fr.get(instr.X), instr.X.Type())

				case *ssa.Next:
					// fr.env[instr] = fr.get(instr.Iter).(iter).next()

				case *ssa.FieldAddr:
					// fr.env[instr] = &(*fr.get(instr.X).(*value)).(structure)[instr.Field]

				case *ssa.Field:
					// fr.env[instr] = fr.get(instr.X).(structure)[instr.Field]

				case *ssa.IndexAddr:
					// x := fr.get(instr.X)
					// idx := fr.get(instr.Index)
					// switch x := x.(type) {
					// case []value:
					// 	fr.env[instr] = &x[asInt(idx)]
					// case *value: // *array
					// 	fr.env[instr] = &(*x).(array)[asInt(idx)]
					// default:
					// 	panic(fmt.Sprintf("unexpected x type in IndexAddr: %T", x))
					// }

				case *ssa.Index:
					// fr.env[instr] = fr.get(instr.X).(array)[asInt(fr.get(instr.Index))]

				case *ssa.Lookup:
					// fr.env[instr] = lookup(instr, fr.get(instr.X), fr.get(instr.Index))

				case *ssa.MapUpdate:
					// m := fr.get(instr.Map)
					// key := fr.get(instr.Key)
					// v := fr.get(instr.Value)
					// switch m := m.(type) {
					// case map[value]value:
					// 	m[key] = v
					// case *hashmap:
					// 	m.insert(key.(hashable), v)
					// default:
					// 	panic(fmt.Sprintf("illegal map type: %T", m))
					// }

				case *ssa.TypeAssert:
					// fr.env[instr] = typeAssert(fr.i, instr, fr.get(instr.X).(iface))

				case *ssa.MakeClosure:
					// var bindings []value
					// for _, binding := range instr.Bindings {
					// 	bindings = append(bindings, fr.get(binding))
					// }
					// fr.env[instr] = &closure{instr.Fn.(*ssa.Function), bindings}

				case *ssa.Phi:
					// dst Variable, phi []PhiValue
					phis := make([]lang.PhiValue, len(instr.Edges))
					for i, pred := range instr.Block().Preds {
						phis[i] = lang.PhiValue{
							Origin: pred.Index,
							Value:  toExp(instr.Edges[i]),
						}

						// fmt.Println("phi:", reflect.TypeOf(pred))

						// phis[i] = lang.PhiValue{
						// 	Origin: bb.Index,
						// 	Value:  pred.Name(),
						// }
					}

					lb.EmitPhi(toVar(instr), phis)

					// for i, pred := range instr.Block().Preds {
					// 	if fr.prevBlock == pred {
					// 		fr.env[instr] = fr.get(instr.Edges[i])
					// 		break
					// 	}
					// }

				case *ssa.Select:
					// var cases []reflect.SelectCase
					// if !instr.Blocking {
					// 	cases = append(cases, reflect.SelectCase{
					// 		Dir: reflect.SelectDefault,
					// 	})
					// }
					// for _, state := range instr.States {
					// 	var dir reflect.SelectDir
					// 	if state.Dir == types.RecvOnly {
					// 		dir = reflect.SelectRecv
					// 	} else {
					// 		dir = reflect.SelectSend
					// 	}
					// 	var send reflect.Value
					// 	if state.Send != nil {
					// 		send = reflect.ValueOf(fr.get(state.Send))
					// 	}
					// 	cases = append(cases, reflect.SelectCase{
					// 		Dir:  dir,
					// 		Chan: reflect.ValueOf(fr.get(state.Chan)),
					// 		Send: send,
					// 	})
					// }
					// chosen, recv, recvOk := reflect.Select(cases)
					// if !instr.Blocking {
					// 	chosen-- // default case should have index -1.
					// }
					// r := tuple{chosen, recvOk}
					// for i, st := range instr.States {
					// 	if st.Dir == types.RecvOnly {
					// 		var v value
					// 		if i == chosen && recvOk {
					// 			// No need to copy since send makes an unaliased copy.
					// 			v = recv.Interface().(value)
					// 		} else {
					// 			v = zero(st.Chan.Type().Underlying().(*types.Chan).Elem())
					// 		}
					// 		r = append(r, v)
					// 	}
					// }
					// fr.env[instr] = r

				default:
					panic(fmt.Sprintf("unexpected instruction: %T", instr))
				}
			}
		}
	}
}
