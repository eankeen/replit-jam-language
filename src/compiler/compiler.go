package compiler

import . "parser"

type Compiler struct {
	Buff       []byte
	ScopeCount int
}

func (c *Compiler) append(buff []byte) {
	c.Buff = append(c.Buff, []byte(buff)...)
}

func (c *Compiler) colon() {
	c.append([]byte(":"))
}

func (c *Compiler) space() {
	c.append([]byte(" "))
}

func (c *Compiler) comma() {
	c.append([]byte(","))
}

func (c *Compiler) semicolon() {
	c.append([]byte(";"))
}

func (c *Compiler) newline() {
	c.append([]byte("\n"))
}

func (c *Compiler) openParen() {
	c.append([]byte("("))
}

func (c *Compiler) closeParen() {
	c.append([]byte(")"))
}

func (c *Compiler) openBrace() {
	c.append([]byte("{"))
}

func (c *Compiler) closeBrace() {
	c.append([]byte("}"))
}

func (c *Compiler) pushScope() {
	c.ScopeCount++
}

func (c *Compiler) popScope() {
	c.ScopeCount--
}

func (c *Compiler) indent() {
	for i := 0; i < c.ScopeCount; i++ {
		c.append([]byte("	"))
	}
}

func (c *Compiler) Statement(stmt Statement) {
	c.newline()
	switch stmt.(type) {
	case Declaration:
		c.declaration(stmt.(Declaration))
	case Return:
		c.rturn(stmt.(Return))
	case IfElseBlock:
		c.ifElse(stmt.(IfElseBlock))
	case Loop:
		c.loop(stmt.(Loop))
	case Assignment:
		c.assignment(stmt.(Assignment))
	case Switch:
		c.swtch(stmt.(Switch))
	default:
		c.indent()
		c.expression(stmt.(Expression))
		c.semicolon()
	}
}

func (c *Compiler) loop(loop Loop) {

	if loop.Type == InitCondLoop || loop.Type == InitCond {
		c.indent()
		c.openBrace()
		c.pushScope()
		c.Statement(loop.InitStatement)
	}

	c.indent()
	c.append([]byte("while"))
	c.openParen()

	if loop.Type == NoneLoop {
		c.append([]byte("1"))
	} else {
		c.expression(loop.Condition)
	}

	c.closeParen()

	if loop.Type == InitCondLoop {
		loop.Block.Statements = append(loop.Block.Statements, loop.LoopStatement)
	}

	c.block(loop.Block)

	if loop.Type == InitCondLoop || loop.Type == InitCond {
		c.popScope()
		c.newline()
		c.indent()
		c.closeBrace()
	}
}

func (c *Compiler) declaration(dec Declaration) {
	hasValues := len(dec.Values) > 0
	//print(dec.Values[0].(BasicLit).Value.Buff)
	switch len(dec.Types) {
	case 0:
	case 1:
		Type := dec.Types[0]

		switch Type.Type {
		case IdentifierType:
			for i, Var := range dec.Identifiers {
				c.indent()
				c.basicTypeNoArray(Type)
				c.append([]byte(" "))
				c.append(Var.Buff)
				if hasValues {
					c.append([]byte(" = "))
					c.expression(dec.Values[i])
				}
				c.semicolon()
				c.newline()
			}
		case FuncType:
			for i, Var := range dec.Identifiers {
				c.indent()
				c.append(dec.Values[i].(FunctionExpression).ReturnTypes[0].Identifier.Buff)
				c.space()
				c.append(Var.Buff)
				c.openParen()

				for _, Arg := range dec.Values[i].(FunctionExpression).Args {
					c.arg(Arg)
				}

				c.closeParen()
				c.block(dec.Values[i].(FunctionExpression).Block)
			}
		}
	default:
		for i, Var := range dec.Identifiers {
			Type := dec.Types[i]

			switch Type.Type {
			case IdentifierType:
				c.indent()
				c.basicTypeNoArray(Type)
				c.space()
				c.append(Var.Buff)

				if hasValues {
					c.append([]byte(" = "))
					c.expression(dec.Values[i])
				}

				c.semicolon()
				c.newline()

			case FuncType:
				c.indent()
				c.append(dec.Values[i].(FunctionExpression).ReturnTypes[0].Identifier.Buff)
				c.space()
				c.append(Var.Buff)
				c.openParen()

				for _, Arg := range dec.Values[i].(FunctionExpression).Args {
					c.arg(Arg)
				}

				c.closeParen()
				c.block(dec.Values[i].(FunctionExpression).Block)
			}
		}
	}
}

func (c *Compiler) rturn(rtrn Return) {
	c.indent()
	c.append([]byte("return"))
	c.space()

	if len(rtrn.Values) > 0 {
		c.expression(rtrn.Values[0])
	}

	c.semicolon()
}

func (c *Compiler) arg(arg ArgStruct) {
	c.basicTypeNoArray(arg.Type)
	c.space()
	c.identifier(arg.Identifier)
	c.comma()
	c.space()
}

func (c *Compiler) block(block Block) {
	c.openBrace()
	c.pushScope()
	for _, statement := range block.Statements {
		c.Statement(statement)
	}
	c.popScope()
	c.newline()
	c.indent()
	c.closeBrace()
}

func (c *Compiler) expression(expr Expression) {

	switch expr.(type) {
	case FunctionCall:
		c.functionCall(expr.(FunctionCall))
	default:
		c.append(expr.(BasicLit).Value.Buff)
	}
}

func (c *Compiler) functionCall(call FunctionCall) {
	c.append(call.Name.Buff)
	c.openParen()

	if len(call.Args) > 0 {
		c.expression(call.Args[0])

		for i := 1; i < len(call.Args); i++ {
			c.comma()
			c.space()
			c.expression(call.Args[i])
		}
	}
	c.closeParen()
}

func (c *Compiler) identifier(identifer Token) {
	c.append(identifer.Buff)
}

func (c *Compiler) basicTypeNoArray(Type TypeStruct) {
	c.append(Type.Identifier.Buff)

	for x := byte(0); x < Type.PointerIndex; x++ {
		c.append([]byte("*"))
	}
}

func (c *Compiler) ifElse(ifElse IfElseBlock) {

	if ifElse.HasInitStmt {
		c.indent()
		c.openBrace()
		c.pushScope()
		c.Statement(ifElse.InitStatement)
	}

	c.indent()
	for i, condition := range ifElse.Conditions {
		c.append([]byte("if("))
		c.expression(condition)
		c.append([]byte(")"))
		c.block(ifElse.Blocks[i])
		c.append([]byte(" else "))
	}

	c.block(ifElse.ElseBlock)

	if ifElse.HasInitStmt {
		c.popScope()
		c.newline()
		c.indent()
		c.closeBrace()
	}
}

func (c *Compiler) assignment(as Assignment) {

	for i, Var := range as.Variables {
		c.indent()
		c.expression(Var)

		if as.Op.SecondaryType != AddAdd && as.Op.SecondaryType != SubSub {
			c.space()
			c.operator(as.Op)
			c.space()
			if len(as.Values) > 1 {
				c.expression(as.Values[i])
			} else {
				c.expression(as.Values[0])
			}
		} else {
			c.operator(as.Op)
		}

		c.semicolon()
		c.newline()
	}
}

func (c *Compiler) swtch(swtch Switch) {
	if swtch.Type == InitCondSwitch {
		c.indent()
		c.openBrace()
		c.pushScope()
		c.Statement(swtch.InitStatement)
	}

	c.indent()
	c.append([]byte("switch"))
	c.openParen()

	if swtch.Type == NoneSwtch {
		c.append([]byte("1"))
	} else {
		c.expression(swtch.Condition)
	}

	c.closeParen()
	c.openBrace()

	for _, Case := range swtch.Cases {
		c.newline()
		c.indent()
		c.append([]byte("case"))
		c.space()
		c.expression(Case.Condition)
		c.colon()

		c.pushScope()
		for _, stmt := range Case.Statements {
			c.Statement(stmt)
		}
		c.popScope()
	}

	if swtch.HasDefaultCase {
		c.newline()
		c.indent()
		c.append([]byte("default"))
		c.colon()

		c.pushScope()
		for _, stmt := range swtch.DefaultCase.Statements {
			c.Statement(stmt)
		}
		c.popScope()
	}

	c.newline()
	c.indent()
	c.closeBrace()

	if swtch.Type == InitCondSwitch {
		c.popScope()
		c.newline()
		c.indent()
		c.closeBrace()
	}
}

func (c *Compiler) operator(op Token) {
	c.append(op.Buff)
}

/*
func (c *Compiler) structTypedef(st StructTypeStruct) {
	c.append("typedef struct {")

	for prop := range st.Props {
		switch prop.Type.Type {
		case IdentifierType:
			c.basicTypeNoArray(prop.Type)
			c.append(prop.Identifier.Buff)
			c.append(";\n")
		}
	}
}
*/