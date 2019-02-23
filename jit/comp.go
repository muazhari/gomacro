/*
 * gomacro - A Go interpreter with Lisp-like macros
 *
 * Copyright (C) 2019 Massimiliano Ghilardi
 *
 *     This Source Code Form is subject to the terms of the Mozilla Public
 *     License, v. 2.0. If a copy of the MPL was not distributed with this
 *     file, You can obtain one at http://mozilla.org/MPL/2.0/.
 *
 *
 * comp.go
 *
 *  Created on Feb 10, 2019
 *      Author Massimiliano Ghilardi
 */

package jit

import (
	"github.com/cosmos72/gomacro/jit/asm"
)

type Comp struct {
	code        Code
	nextSoftReg SoftRegId
	arch        Arch
	asm.RegIdConfig
	asm *Asm
}

func New() *Comp {
	var c Comp
	return c.Init()
}

func NewArchId(archId ArchId) *Comp {
	var c Comp
	return c.InitArchId(archId)
}

func NewArch(arch Arch) *Comp {
	var c Comp
	return c.InitArch(arch)
}

func (c *Comp) Init() *Comp {
	return c.InitArchId(asm.ARCH_ID)
}

func (c *Comp) InitArchId(archId ArchId) *Comp {
	return c.InitArch(Archs[archId])
}

func (c *Comp) InitArch(arch Arch) *Comp {
	if arch == nil {
		errorf("unknown arch")
	}
	c.code = nil
	c.nextSoftReg = 0
	c.arch = arch
	c.RegIdConfig = arch.RegIdConfig()
	if c.asm != nil {
		c.asm.InitArch(arch)
	}
	return c
}

func (c *Comp) Arch() Arch {
	return c.arch
}

func (c *Comp) ArchId() ArchId {
	if c.arch == nil {
		return asm.NOARCH
	}
	return c.arch.Id()
}

func (c *Comp) Asm() *Asm {
	if c.asm == nil {
		// create asm.Asm on demand
		c.asm = asm.NewArch(c.arch)
	}
	return c.asm
}

// return compiled assembly code
func (c *Comp) Code() Code {
	return c.code
}

// discard compiled assembly code
func (c *Comp) ClearCode() {
	c.code = nil
}

func (c *Comp) AllocSoftReg(kind Kind) SoftReg {
	id := c.nextSoftReg
	c.nextSoftReg++
	return c.code.SoftReg(asm.ALLOC, id, kind)
}

func (c *Comp) FreeSoftReg(s SoftReg) {
	if s.Valid() {
		if s.id+1 == c.nextSoftReg {
			c.nextSoftReg--
		}
		c.code.SoftReg(asm.FREE, s.id, s.kind)
	}
}

func checkAssignable(e Expr) {
	switch e.(type) {
	case Reg, Mem, SoftReg:
		break
	default:
		errorf("cannot assign to %v", e)
	}
}

func (c *Comp) MakeVar(idx int, upn int, kind Kind) Mem {
	mem, err := MakeVar(idx, upn, kind, c.RegIdConfig)
	if err != nil {
		panic(err)
	}
	return mem
}

// compile list of statements
func (c *Comp) Compile(ts ...Stmt) {
	for _, t := range ts {
		c.Stmt(t)
	}
}
