// Copyright 2014-2018 Brett Vickers. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package disasm implements a 6502 instruction set
// disassembler.
package disasm

import (
	"fmt"
	"strings"

	"github.com/cjr29/go6502/cpu"
)

// Theme is a struct of color escape codes used to colorize the output
// of the disassembler.
type Theme struct {
	Addr       string
	Code       string
	Inst       string
	Operand    string
	RegName    string
	RegValue   string
	ReqEqual   string
	RegFlagOn  string
	RegFlagOff string
	Annotation string
	Reset      string
}

// Disassembler formatting for addressing modes
var modeFormat = []string{
	"#$%s",    // IMM
	"%s",      // IMP
	"$%s",     // REL
	"$%s",     // ZPG
	"$%s,X",   // ZPX
	"$%s,Y",   // ZPY
	"$%s",     // ABS
	"$%s,X",   // ABX
	"$%s,Y",   // ABY
	"($%s)",   // IND
	"($%s,X)", // IDX
	"($%s),Y", // IDY
	"%s",      // ACC
}

var hex = "0123456789ABCDEF"

type Flags uint8

const (
	ShowAddress Flags = 1 << iota
	ShowCode
	ShowInstruction
	ShowRegisters
	ShowCycles
	ShowAnnotations

	ShowBasic = ShowAddress | ShowCode | ShowInstruction | ShowAnnotations
	ShowFull  = ShowAddress | ShowCode | ShowInstruction | ShowRegisters | ShowCycles
)

// Disassemble the machine code at memory address addr. Return a string
// representing the disassembled instruction and the address of the next
// instruction.
func Disassemble(c *cpu.CPU, addr uint16, flags Flags, anno string, theme *Theme) (line string, next uint16) {
	opcode := c.Mem.LoadByte(addr)
	inst := c.InstSet.Lookup(opcode)
	next = addr + uint16(inst.Length)
	line = ""

	if (flags & ShowAddress) != 0 {
		line += fmt.Sprintf("%04X- ", addr)
	}

	if (flags & ShowCode) != 0 {
		var csbuf [3]byte
		c.Mem.LoadBytes(addr, csbuf[:next-addr])
		line += fmt.Sprintf("%-8s  ", codeString(csbuf[:next-addr]))
	}

	if (flags & ShowInstruction) != 0 {
		var buf [2]byte
		operand := buf[:inst.Length-1]
		c.Mem.LoadBytes(addr+1, operand)
		if inst.Mode == cpu.REL {
			// Convert relative offset to absolute address.
			operand = buf[:]
			braddr := int(addr) + int(inst.Length) + byteToInt(operand[0])
			operand[0] = byte(braddr)
			operand[1] = byte(braddr >> 8)
		}

		// Return string composed of CPU instruction and operand.
		line += fmt.Sprintf("%s   "+modeFormat[inst.Mode], inst.Name, hexString(operand))

		// Pad to next column using uncolorized version of the operand.
		dummy := fmt.Sprintf(modeFormat[inst.Mode], hexString(operand))
		line += strings.Repeat(" ", 9-len(dummy))
	}

	if (flags & ShowRegisters) != 0 {
		r := c.Reg
		line += fmt.Sprintf("A=%02X X=%02X Y=%02X PS=[%s] SP=%02X PC=%04X ",
			r.A, r.X, r.Y, getStatusBits(&r), r.SP, r.PC)
	}

	if (flags & ShowCycles) != 0 {
		line += fmt.Sprintf("C=%d", c.Cycles)
	}

	if (flags&ShowAnnotations) != 0 && anno != "" {
		line += " ; " + anno
	}

	return line, next
}

// GetRegisterString returns a string describing the contents of the 6502
// registers.
func GetRegisterString(r *cpu.Registers) string {
	return fmt.Sprintf("A=%02X X=%02X Y=%02X PS=[%s] SP=%02X PC=%04X",
		r.A, r.X, r.Y, getStatusBits(r), r.SP, r.PC)
}

// GetCompactRegisterString returns a compact string describing the contents
// of the 6502 registers. It excludes the program counter and stack pointer.
func GetCompactRegisterString(r *cpu.Registers) string {
	return fmt.Sprintf("A=%02X X=%02X Y=%02X PS=[%s]", r.A, r.X, r.Y, getStatusBits(r))
}

func codeString(b []byte) string {
	switch len(b) {
	case 1:
		return fmt.Sprintf("%02X", b[0])
	case 2:
		return fmt.Sprintf("%02X %02X", b[0], b[1])
	case 3:
		return fmt.Sprintf("%02X %02X %02X", b[0], b[1], b[2])
	default:
		return ""
	}
}

// Return a hexadecimal string representation of the byte slice.
func hexString(b []byte) string {
	hexlen := len(b) * 2
	hexbuf := make([]byte, hexlen)
	j := hexlen - 1
	for _, n := range b {
		hexbuf[j] = hex[n&0xf]
		hexbuf[j-1] = hex[n>>4]
		j -= 2
	}
	return string(hexbuf)
}

func getStatusBits(r *cpu.Registers) string {
	v := func(bit bool, ch byte) byte {
		if bit {
			return ch
		}
		return '-'
	}
	b := []byte{
		v(r.Sign, 'N'),
		v(r.Zero, 'Z'),
		v(r.Carry, 'C'),
		v(r.InterruptDisable, 'I'),
		v(r.Decimal, 'D'),
		v(r.Overflow, 'V'),
	}
	return string(b)
}

func byteToInt(b byte) int {
	if b >= 0x80 {
		return int(b) - 256
	}
	return int(b)
}
