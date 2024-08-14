package cpu_test

import (
	CPU "chip-8-go/cpu"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewChip(t *testing.T) {
	assert := assert.New(t)

	cpu := CPU.NewCPU()

	assert.Equal(uint16(CPU.START_ADDR), cpu.ProgramCounter, "ProgramCounter should be set to START_ADDR")
	assert.Equal(uint8(0x00), cpu.Ram[CPU.START_ADDR], "Initial RAM value at START_ADDR should be 0")
	assert.Equal(uint16(0), cpu.StackPointer, "StackPointer should be initialized to 0")
	assert.Equal([CPU.SCREEN_HEIGHT][CPU.SCREEN_WIDTH]bool{}, cpu.Screen(), "Screen should be cleared")

	for i := 0; i < CPU.FONTSET_SIZE; i++ {
		assert.Equal(CPU.FONTSET[i], cpu.Ram[i], "FONTSET should be correctly loaded into RAM")
	}
}

func TestTickTimers(t *testing.T) {
	assert := assert.New(t)

	cpu := CPU.NewCPU()

	cpu.DelayTimer = 10
	cpu.SoundTimer = 5

	cpu.TickTimers()

	assert.Equal(uint8(9), cpu.DelayTimer, "DelayTimer should decrement by 1")
	assert.Equal(uint8(4), cpu.SoundTimer, "SoundTimer should decrement by 1")

	for i := 0; i < 9; i++ {
		cpu.TickTimers()
	}

	assert.Equal(uint8(0), cpu.DelayTimer, "DelayTimer should reach 0")
	assert.Equal(uint8(0), cpu.SoundTimer, "SoundTimer should reach 0")
}

func TestGetOpCode(t *testing.T) {
	assert := assert.New(t)

	cpu := CPU.NewCPU()

	cpu.Ram[cpu.ProgramCounter] = 0xAB
	cpu.Ram[cpu.ProgramCounter+1] = 0xCD

	opcode := cpu.GetOpCode()

	assert.Equal(uint16(0xABCD), opcode, "getOpCode() should return the correct opcode")
	assert.Equal(uint16(CPU.START_ADDR+2), cpu.ProgramCounter, "ProgramCounter should advance by 2 after getOpCode()")
}

func TestSetKey(t *testing.T) {
	assert := assert.New(t)

	cpu := CPU.NewCPU()

	cpu.SetKey(0x1, true)
	assert.True(cpu.Keys[0x1], "Key 0x1 should be set to pressed")

	cpu.SetKey(0x1, false)
	assert.False(cpu.Keys[0x1], "Key 0x1 should be set to not pressed")
}
