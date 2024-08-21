package cpu

import (
	"fmt"
	"math/rand"
)

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

func (op OpCode) Execute(cpu *CPU) (bool, error) {
	opCode := op.Decode()

	// To keep track of opcodes which will jump
	// so no need to update ProgramCounter
	count := true

	switch {
	case opCode.n1 == 0x0:
		switch {
		case opCode.n3 == 0x0 && opCode.n4 == 0x0:
			// NOOP
			break
		case opCode.n3 == 0xE && opCode.n4 == 0x0:
			// Clear Screen
			cpu.ClearScreen()
		case opCode.n3 == 0xE && opCode.n4 == 0xE:
			// Return from subroutine
			cpu.ProgramCounter = cpu.Pop()
		default:
			fmt.Printf("Invalid opcode n1: 0x%x, n2: 0x%x, n3: 0x%x, n4: 0x%x\n", opCode.n1, opCode.n2, opCode.n3, opCode.n4)
		}
	case opCode.n1 == 0x1:
		// Jump to NNN
		address := uint16(op) & 0x0FFF
		cpu.ProgramCounter = address

		count = false

	case opCode.n1 == 0x2:
		// Call subroutine
		cpu.Push(cpu.ProgramCounter)
		address := uint16(op) & 0x0FFF
		cpu.ProgramCounter = address

		count = false

	case opCode.n1 == 0x3:
		// Skip next if VX == NN
		NN := uint8(op) & 0x0FF
		if cpu.VRegisters[opCode.n2] == NN {
			cpu.ProgramCounter += 2
		}
	case opCode.n1 == 0x4:
		// Skip next if VX != NN
		NN := uint8(op) & 0x0FF
		if cpu.VRegisters[opCode.n2] != NN {
			cpu.ProgramCounter += 2
		}
	case opCode.n1 == 0x5 && opCode.n4 == 0:
		// Skip next if VX == VY
		if cpu.VRegisters[opCode.n2] == cpu.VRegisters[opCode.n3] {
			cpu.ProgramCounter += 2
		}
	case opCode.n1 == 0x6:
		// Set VX = NN
		NN := opCode.n3<<4 | opCode.n4
		cpu.VRegisters[opCode.n2] = NN
	case opCode.n1 == 0x7:
		// Set VX += NN
		NN := opCode.n3<<4 | opCode.n4
		cpu.VRegisters[opCode.n2] += NN
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
			// VX += VY
			result, overflowed := OverflowAdd(cpu.VRegisters[opCode.n2], cpu.VRegisters[opCode.n3])
			cpu.VRegisters[opCode.n2] = result
			if overflowed {
				cpu.VRegisters[0xF] = 1
			} else {
				cpu.VRegisters[0xF] = 0
			}
		case 0x5:
			// VX -= VY
			result, overflowed := OverflowSub(cpu.VRegisters[opCode.n2], cpu.VRegisters[opCode.n3])
			cpu.VRegisters[opCode.n2] = result
			if overflowed {
				cpu.VRegisters[0xF] = 0
			} else {
				cpu.VRegisters[0xF] = 1
			}
		case 0x6:
			// VX >>= 1
			dropedBit := cpu.VRegisters[opCode.n2] & 1
			cpu.VRegisters[opCode.n2] >>= 1
			cpu.VRegisters[0xF] = dropedBit
		case 0x7:
			// VX = VY - VX
			result, overflowed := OverflowSub(cpu.VRegisters[opCode.n3], cpu.VRegisters[opCode.n2])
			cpu.VRegisters[opCode.n2] = result
			if overflowed {
				cpu.VRegisters[0xF] = 0
			} else {
				cpu.VRegisters[0xF] = 1
			}
		case 0xE:
			// VX <<= 1
			overflowedBit := (cpu.VRegisters[opCode.n2] >> 7) & 1
			cpu.VRegisters[opCode.n2] <<= 1
			cpu.VRegisters[0xF] = overflowedBit
		}
	case opCode.n1 == 0x9:
		// Skip if VX != VY
		if cpu.VRegisters[opCode.n2] != cpu.VRegisters[opCode.n3] {
			cpu.ProgramCounter += 2
		}
	case opCode.n1 == 0xA:
		// I = 0xNNN
		cpu.IndexRegister = uint16(op) & 0x0FFF
	case opCode.n1 == 0xB:
		// Jump to V0 + NNN
		cpu.ProgramCounter = uint16(cpu.VRegisters[0x0]) + uint16(op)&0x0FFF
	case opCode.n1 == 0xC:
		// VX = random & NN
		NN := opCode.n3<<3 | opCode.n4
		cpu.VRegisters[opCode.n2] = byte(rand.Uint32()) & NN
	case opCode.n1 == 0xD:
		// Draw Sprite
		x := cpu.VRegisters[opCode.n2]
		y := cpu.VRegisters[opCode.n3]
		height := opCode.n4

		var xLine uint16
		var yLine uint16
		var collision uint8 = 0

		for yLine = 0; yLine < uint16(height); yLine++ {
			pixels := cpu.Memory[cpu.IndexRegister+uint16(yLine)]

			for xLine = 0; xLine < 8; xLine++ {
				// Compute the pixel's position
				px := (x + uint8(xLine)) % SCREEN_WIDTH
				py := (y + uint8(yLine)) % SCREEN_HEIGHT

				// Fetch the current pixel value
				currentPixel := &cpu.Screen[py][px]

				// Check if the bit is set in the sprite data
				if (pixels & (0b1000_0000 >> xLine)) != 0 {
					// If pixel is already set, set VF to 1
					if *currentPixel {
						collision = 1
					}

					// Flip the pixel
					*currentPixel = !*currentPixel
				}
			}
		}

		// Set VF to 1 if there was a collision
		cpu.VRegisters[0xF] = collision
		cpu.shouldDraw = true

	case opCode.n1 == 0xE:
		switch {
		case opCode.n3 == 0x9 && opCode.n4 == 0xE:
			// Skip if key pressed
			if cpu.Keys[cpu.VRegisters[opCode.n2]] {
				cpu.ProgramCounter += 2
			}
		case opCode.n3 == 0xA && opCode.n4 == 0x1:
			// Skip if key released
			if !cpu.Keys[cpu.VRegisters[opCode.n2]] {
				cpu.ProgramCounter += 2
			}
		default:
			fmt.Printf("Invalid opcode n1: 0x%x, n2: 0x%x, n3: 0x%x, n4: 0x%x\n", opCode.n1, opCode.n2, opCode.n3, opCode.n4)
		}
	case opCode.n1 == 0xF:
		switch {
		case opCode.n3 == 0 && opCode.n4 == 0x7:
			// VX = DelayTimer
			cpu.VRegisters[opCode.n2] = cpu.DelayTimer
		case opCode.n3 == 0 && opCode.n4 == 0xA:
			// Wait for keypress
			pressed := false

			for i := 0; i < len(cpu.Keys); i++ {
				if cpu.Keys[i] {
					cpu.VRegisters[opCode.n2] = uint8(i)
					pressed = true
					break
				}
			}

			if !pressed {
				cpu.ProgramCounter -= 2
			}
		case opCode.n3 == 1 && opCode.n4 == 0x5:
			// DelayTimer = VX
			cpu.DelayTimer = cpu.VRegisters[opCode.n2]
		case opCode.n3 == 1 && opCode.n4 == 0x8:
			// SoundTimer = VX
			cpu.SoundTimer = cpu.VRegisters[opCode.n2]
		case opCode.n3 == 1 && opCode.n4 == 0xE:
			// I += VX
			result, overflow := OverflowAdd(cpu.IndexRegister, uint16(cpu.VRegisters[opCode.n2]))
			cpu.IndexRegister = result
			if overflow {
				cpu.IndexRegister = 0
			}
		case opCode.n3 == 2 && opCode.n4 == 9:
			// Set I to Font Address
			cpu.IndexRegister = uint16(cpu.VRegisters[opCode.n2]) * 5
		case opCode.n3 == 3 && opCode.n4 == 3:
			// Binary-Coded Decimal of VX stored in RAM
			VX := cpu.VRegisters[opCode.n2]

			cpu.Memory[cpu.IndexRegister] = VX / 100
			cpu.Memory[cpu.IndexRegister+1] = (VX / 10) % 10
			cpu.Memory[cpu.IndexRegister+2] = VX % 10

		case opCode.n3 == 0x5 && opCode.n4 == 0x5:
			// Store V0 - VX
			for i := 0; i <= int(opCode.n2); i++ {
				cpu.Memory[int(cpu.IndexRegister)+i] = cpu.VRegisters[i]

			}
			cpu.IndexRegister = uint16(opCode.n2) + 1
		case opCode.n3 == 0x6 && opCode.n4 == 0x5:
			// Load V0 - VX
			for i := 0; i <= int(opCode.n2); i++ {
				cpu.VRegisters[i] = cpu.Memory[int(cpu.IndexRegister)+i]

			}
			cpu.IndexRegister = uint16(opCode.n2) + 1
		default:
			fmt.Printf("Invalid opcode n1: 0x%x, n2: 0x%x, n3: 0x%x, n4: 0x%x\n", opCode.n1, opCode.n2, opCode.n3, opCode.n4)
		}

	default:
		fmt.Printf("Invalid opcode n1: 0x%x, n2: 0x%x, n3: 0x%x, n4: 0x%x\n", opCode.n1, opCode.n2, opCode.n3, opCode.n4)
	}
	return count, nil
}
