package cpu

import "fmt"

type OpCode uint16

// OpCode decoding binary example on high nibble (n2)
//
//        n1   n2   n3   n4
//       1010 1011 1100 1101  (op)
// AND   0000 1111 0000 0000  (0x0F00)
// ------------------------------------
//       0000 1011 0000 0000  (result)
// >> 8  0000 0000 0000 1011  (n1)

type OpCodeParts struct {
	n1, n2, n3, n4 uint8
}

func (op OpCode) Decode() OpCodeParts {
	return OpCodeParts{
		n1: uint8(op >> 12),
		n2: uint8((op >> 8) & 0x0F),
		n3: uint8((op >> 4) & 0x0F),
		n4: uint8(op & 0x0F),
	}
}

func (op OpCode) Execute(cpu *CPU) error {
	opCode := op.Decode()

	switch {
	case opCode.n1 == 0x0:
		switch {
		case opCode.n2 == 0x0 && opCode.n3 == 0x0 && opCode.n4 == 0x0:
			// NOOP
			break
		case opCode.n2 == 0x0 && opCode.n3 == 0xE && opCode.n4 == 0x0:
			// Clear Screen
			cpu.ClearScreen()
		case opCode.n2 == 0x0 && opCode.n3 == 0xE && opCode.n4 == 0xE:
			// Return from subroutine
			cpu.ProgramCounter = cpu.Pop()
		default:
			fmt.Printf("Invalid opcode %X\n", op)
		}
	case opCode.n1 == 0x1:
		// Jump to NNN
		cpu.ProgramCounter = uint16(op) & 0x0FFF
	case opCode.n1 == 0x2:
		// Call subroutine
		cpu.Push(uint16(op) & 0x0FFF)
	case opCode.n1 == 0x3:
		// Skip next if VX == NN
		if cpu.VRegisters[opCode.n2] == opCode.n3+opCode.n4 {
			cpu.ProgramCounter += 2
		}
	case opCode.n1 == 0x4:
		// Skip next if VX != NN
		if cpu.VRegisters[opCode.n2] != opCode.n3+opCode.n4 {
			cpu.ProgramCounter += 2
		}
	case opCode.n1 == 0x5:
		// Skip next if VX == VY
		if cpu.VRegisters[opCode.n2] != cpu.VRegisters[opCode.n3] {
			cpu.ProgramCounter += 2
		}
	case opCode.n1 == 0x6:
		// Set VX = NN
		cpu.VRegisters[opCode.n2] = opCode.n3 + opCode.n4
	case opCode.n1 == 0x7:
		// Set VX += NN
		cpu.VRegisters[opCode.n2] += opCode.n3 + opCode.n4
	case opCode.n1 == 0x8:
		switch opCode.n4 {
		case 0x0:
			// Set VX = VY
			cpu.VRegisters[opCode.n2] = cpu.VRegisters[opCode.n3]
		case 0x1:
			// Set VX |= VY
			cpu.VRegisters[opCode.n2] |= cpu.VRegisters[opCode.n3]
		case 0x2:
			// Set VX &= VY
			cpu.VRegisters[opCode.n2] &= cpu.VRegisters[opCode.n3]
		case 0x3:
			// Set VX ^= VY
			cpu.VRegisters[opCode.n2] ^= cpu.VRegisters[opCode.n3]
		case 0x4:
		case 0x5:
		case 0x6:
		case 0x7:
		case 0xE:
		}
	case opCode.n1 == 0x9:
	case opCode.n1 == 0xA:
	case opCode.n1 == 0xB:
	case opCode.n1 == 0xC:
	case opCode.n1 == 0xD:
	case opCode.n1 == 0xE:
	case opCode.n1 == 0xF:

	default:
		return fmt.Errorf("unexpected opcode: 0x%X", op)
	}

	return nil
}
