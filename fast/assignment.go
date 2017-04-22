/*
 * gomacro - A Go interpreter with Lisp-like macros
 *
 * Copyright (C) 2017 Massimiliano Ghilardi
 *
 *     This program is free software you can redistribute it and/or modify
 *     it under the terms of the GNU General Public License as published by
 *     the Free Software Foundation, either version 3 of the License, or
 *     (at your option) any later version.
 *
 *     This program is distributed in the hope that it will be useful,
 *     but WITHOUT ANY WARRANTY; without even the implied warranty of
 *     MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 *     GNU General Public License for more details.
 *
 *     You should have received a copy of the GNU General Public License
 *     along with this program.  If not, see <http//www.gnu.org/licenses/>.
 *
 * declaration.go
 *
 *  Created on Apr 01, 2017
 *      Author Massimiliano Ghilardi
 */

package fast

import (
	"go/ast"
	"go/token"
)

// Assign compiles an *ast.AssignStmt into an assignment to one or more place
func (c *Comp) Assign(node *ast.AssignStmt) {
	lhs, rhs := node.Lhs, node.Rhs

	if node.Tok == token.DEFINE {
		c.DeclVarsShort(lhs, rhs)
		return
	}

	ln, rn := len(lhs), len(rhs)
	if ln > 1 && rn == 1 {
		c.Errorf("unimplemented: assignment of multiple places with a single multi-valued expression: %v", node)
	} else if ln != rn {
		c.Errorf("invalid assignment, cannot assign %d values to %d places: %v", node)
	} else {
		for i := range lhs {
			c.Assign1(lhs[i], node.Tok, rhs[i])
		}
	}
}

// Assign1 compiles a single assignment to a place
func (c *Comp) Assign1(lhs ast.Expr, op token.Token, rhs ast.Expr) {
	node := &ast.AssignStmt{Lhs: []ast.Expr{lhs}, Tok: op, Rhs: []ast.Expr{rhs}} // only for nice error messages

	place := c.Place(lhs)
	if place.Fun != nil {
		c.Errorf("unimplemented: assignment of composite place (only assignment of variables is implemented): %v", node)
	}
	va := &place.Var
	init := c.Expr(rhs)

	panicking := true
	defer func() {
		if !panicking {
			return
		}
		rec := recover()
		c.Errorf("error compiling assignment: %v\n    %v", node, rec)
	}()
	c.SetVar(va, op, init)
	panicking = false
}

// SetVar compiles an assignment to a variable
func (c *Comp) SetVar(va *Var, op token.Token, init *Expr) {
	if init.Type == nil {
		c.Errorf("invalid operator %v on <%v>", op, init.Type)
	}
	c.varSetOp(va, op, init)
}

// Place compiles the left-hand-side of an assignment
func (c *Comp) Place(lhs ast.Expr) *Place {
	return c.PlaceOrAddress(lhs, false)
}

// PlaceOrAddress compiles the left-hand-side of an assignment or the location of an address-of
func (c *Comp) PlaceOrAddress(lhs ast.Expr, addressof bool) *Place {
	switch lhs := lhs.(type) {
	case *ast.Ident:
		name := lhs.Name
		if name == "_" {
			if addressof {
				c.Errorf("cannot take the address of _")
				return nil
			}
			return &Place{}
		}
		sym := c.Resolve(name)
		class := sym.Desc.Class()
		if class != VarBind && class != IntBind {
			if addressof {
				c.Errorf("cannot take the address of %s: %v", class, name)
			} else {
				c.Errorf("cannot assign to %s: %v", class, name)
			}
			return nil
		}
		return &Place{Var: *sym.AsVar()}
	default:
		if addressof {
			c.Errorf("unimplemented: address of non-identifier: %v", lhs)
		} else {
			c.Errorf("unimplemented: assignment of non-identifier: %v", lhs)
		}
		return nil
	}
}

// LookupVar compiles the left-hand-side of an assignment, in case it's an identifier (i.e. a variable name)
func (c *Comp) LookupVar(name string) *Var {
	if name == "_" {
		return &Var{}
	}
	sym := c.Resolve(name)
	class := sym.Desc.Class()
	if class != VarBind && class != IntBind {
		c.Errorf("cannot assign to %s: %v", class, name)
	}
	return sym.AsVar()
}