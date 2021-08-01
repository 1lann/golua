package lua

import (
	"bytes"
	"strconv"
	"strings"

	"github.com/1lann/golua/lang"
)

type Lua struct {
	Functions map[string]*Function
}

func New() *Lua {
	return &Lua{
		Functions: make(map[string]*Function),
	}
}

type Function struct {
	Name      string
	Variables map[string]struct{}
	Inputs    []lang.Variable
	Outputs   []lang.Variable
	Blocks    []*Block
}

var _ lang.Function = (*Function)(nil)

type Block struct {
	ID         int
	Function   *Function
	Statements []string
}

var _ lang.Block = (*Block)(nil)

func (l *Lua) StartFunction(name string, inputs []lang.Variable, outputs []lang.Variable) lang.Function {
	f := &Function{
		Name:      name,
		Variables: make(map[string]struct{}),
		Inputs:    inputs,
		Outputs:   outputs,
	}
	l.Functions[name] = f
	return f
}

func (f *Function) StartBlock(id int) lang.Block {
	b := &Block{
		ID:         id,
		Function:   f,
		Statements: make([]string, 0),
	}
	f.Blocks = append(f.Blocks, b)
	if id != len(f.Blocks)-1 {
		panic("block ID mismatch")
	}

	return b
}

// EmitNew(dst *Variable, params []Expression)
// EmitAssign(dst *Variable, exp Expression)
// EmitPhi(dst *Variable, phi []PhiValue)
// EmitCall(target string, input []Expression, output []*Variable)
// EmitJump(id int)
// EmitReturn([]Expression)
// EmitIf(condition Expression, thenID int, elseID int)

func (b *Block) ensure(v lang.Variable) {
	for _, input := range b.Function.Inputs {
		if input.Name == v.Name {
			return
		}
	}

	b.Function.Variables[v.Name] = struct{}{}
}

func (b *Block) EmitNew(dst lang.Variable, params []lang.Expression) {
	b.ensure(dst)
	// TODO: implement lol
}

func (b *Block) EmitAssign(dst lang.Variable, exp lang.Expression) {
	b.ensure(dst)
	b.add(dst.Name + " = " + expToStr(exp))
}

func (b *Block) EmitPhi(dst lang.Variable, phi []lang.PhiValue) {
	b.ensure(dst)

	dict := new(bytes.Buffer)
	dict.WriteString("({")

	for i, v := range phi {
		dict.WriteString("[" + strconv.Itoa(v.Origin) + "] = " + expToStr(v.Value))
		if i != len(phi)-1 {
			dict.WriteString(", ")
		}
	}

	dict.WriteString("})[pid]")

	b.add(dst.Name + " = " + dict.String())
}

var binCalls = map[string]string{
	"+":  "+",
	"-":  "-",
	"*":  "*",
	"/":  "/",
	"%":  "%",
	"&":  "&",
	"|":  "|",
	"^":  "^",
	"<<": "<<",
	">>": ">>",
	"&&": "and",
	"||": "or",
	"==": "==",
	"!=": "~=",
	"<":  "<",
	">":  ">",
	"<=": "<=",
	">=": ">=",
}

func (b *Block) EmitCall(target string, input []lang.Expression, output []lang.Variable) {
	inputs := make([]string, len(input))
	outputs := make([]string, len(output))

	for i, exp := range input {
		inputs[i] = expToStr(exp)
	}
	for i, v := range output {
		b.ensure(v)
		outputs[i] = v.Name
	}

	if mapping, found := binCalls[target]; found {
		b.add(outputs[0] + " = " + inputs[0] + " " + mapping + " " + inputs[1])
		return
	}

	if target == "println" {
		target = "print"
	}

	if len(outputs) == 0 {
		b.add(target + "(" + strings.Join(inputs, ",") + ")")
	} else {
		b.add(strings.Join(outputs, ",") + " = " + target + "(" + strings.Join(inputs, ",") + ")")
	}
}

func (b *Block) EmitJump(id int) {
	b.emitJump(id, "")
}

func (b *Block) emitJump(id int, prefix string) {
	b.add(prefix + "cid, pid = " + strconv.Itoa(id) + ", " + strconv.Itoa(b.ID))
}

func (b *Block) EmitReturn(output []lang.Expression) {
	outputs := make([]string, len(output))
	for i, exp := range output {
		outputs[i] = expToStr(exp)
	}

	b.add("return " + strings.Join(outputs, ","))
}

func (b *Block) EmitIf(condition lang.Expression, thenID int, elseID int) {
	b.add("if " + expToStr(condition) + " then")
	b.emitJump(thenID, "\t")
	b.add("else")
	b.emitJump(elseID, "\t")
	b.add("end")
}

func expToStr(exp lang.Expression) string {
	if exp.Variable != nil {
		return exp.Variable.Name
	}

	return exp.Value
}

func (b *Block) add(statement string) {
	b.Statements = append(b.Statements, statement)
}

func (l *Lua) Assemble() []byte {
	buf := new(bytes.Buffer)

	for _, f := range l.Functions {
		inputNames := make([]string, len(f.Inputs))
		for i, v := range f.Inputs {
			inputNames[i] = v.Name
		}

		buf.WriteString("function " + f.Name + "(" + strings.Join(inputNames, ",") + ")\n")

		locals := make([]string, 0, len(f.Variables)+1)
		locals = append(locals, "pid")
		// initialize local variables
		for v := range f.Variables {
			locals = append(locals, v)
		}

		buf.WriteString("\tlocal " + strings.Join(locals, ", ") + "\n")
		buf.WriteString("\tlocal cid = 0\n")

		buf.WriteString("\twhile true do\n")
		// define blocks
		for i, b := range f.Blocks {
			if i == 0 {
				buf.WriteString("\t\tif cid == 0 then\n")
			} else {
				buf.WriteString("\t\telseif cid == " + strconv.Itoa(i) + " then\n")
			}

			for _, s := range b.Statements {
				buf.WriteString("\t\t\t" + s + "\n")
			}
		}
		buf.WriteString("\t\tend\n")
		buf.WriteString("\tend\n")

		buf.WriteString("end\n\n")
	}

	buf.WriteString("init()\n")
	buf.WriteString("main()\n")

	return buf.Bytes()
}

var _ lang.Lang = (*Lua)(nil)
