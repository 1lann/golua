package lang

import "go/types"

type Type struct {
	types.Type
}

type Expression struct {
	Variable *Variable
	Type     Type
	Value    string
}

type Variable struct {
	Name  string
	Type  Type
	DeRef bool
}

type Call struct {
	Target string
	Input  []*Variable
	Output []*Variable
}

type PhiValue struct {
	Origin int
	Value  Expression
}

type Branch struct {
	Condition *Variable
}

type Block interface {
	EmitNew(dst Variable, params []Expression)
	EmitAssign(dst Variable, exp Expression)
	EmitPhi(dst Variable, phi []PhiValue)
	EmitCall(target string, input []Expression, output []Variable)
	EmitJump(id int)
	EmitReturn([]Expression)
	EmitIf(condition Expression, thenID int, elseID int)
}

type Function interface {
	StartBlock(id int) Block
}

// Lang represents a high level language that we're targeting.
type Lang interface {
	StartFunction(name string, inputs []Variable, outputs []Variable) Function
}
