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
	parts := op.Decode()

	switch parts {
	case OpCodeParts{0x0, 0x0, 0x0, 0x0}:
		// NOOP
		break
	case OpCodeParts{0x0, 0x0, 0xE, 0x0}:
		// Clear Screen
		cpu.ClearScreen()
	default:
		return fmt.Errorf("unexpected opcode: 0x%X", op)
	}

	return nil
}
