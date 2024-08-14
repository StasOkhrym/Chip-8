package cpu

const RAM_SIZE = 4096
const NUM_REGS = 16
const STACK_SIZE = 16
const NUM_KEYS = 16

const SCREEN_WIDTH = 64
const SCREEN_HEIGHT = 32

const START_ADDR = 0x200

const FONTSET_SIZE = 80

var FONTSET = [FONTSET_SIZE]uint8{
	0xF0, 0x90, 0x90, 0x90, 0xF0, // 0
	0x20, 0x60, 0x20, 0x20, 0x70, // 1
	0xF0, 0x10, 0xF0, 0x80, 0xF0, // 2
	0xF0, 0x10, 0xF0, 0x10, 0xF0, // 3
	0x90, 0x90, 0xF0, 0x10, 0x10, // 4
	0xF0, 0x80, 0xF0, 0x10, 0xF0, // 5
	0xF0, 0x80, 0xF0, 0x90, 0xF0, // 6
	0xF0, 0x10, 0x20, 0x40, 0x40, // 7
	0xF0, 0x90, 0xF0, 0x90, 0xF0, // 8
	0xF0, 0x90, 0xF0, 0x10, 0xF0, // 9
	0xF0, 0x90, 0xF0, 0x90, 0x90, // A
	0xE0, 0x90, 0xE0, 0x90, 0xE0, // B
	0xF0, 0x80, 0x80, 0x80, 0xF0, // C
	0xE0, 0x90, 0x90, 0x90, 0xE0, // D
	0xF0, 0x80, 0xF0, 0x80, 0xF0, // E
	0xF0, 0x80, 0xF0, 0x80, 0x80, // F
}

type CPU struct {
	// The 4096 bytes of memory.
	//
	// Memory Map:
	// +---------------+= 0xFFF (4095) End of Chip-8 RAM
	// |               |
	// |               |
	// |               |
	// |               |
	// |               |
	// | 0x200 to 0xFFF|
	// |     Chip-8    |
	// | Program / Data|
	// |     Space     |
	// |               |
	// |               |
	// |               |
	// +- - - - - - - -+= 0x600 (1536) Start of ETI 660 Chip-8 programs
	// |               |
	// |               |
	// |               |
	// +---------------+= 0x200 (512) Start of most Chip-8 programs
	// | 0x000 to 0x1FF|
	// | Reserved for  |
	// |  interpreter  |
	// +---------------+= 0x000 (0) Start of Chip-8 RAM
	Ram [RAM_SIZE]uint8

	ProgramCounter uint16

	StackPointer uint16
	Stack        [STACK_SIZE]uint16

	screen [SCREEN_HEIGHT][SCREEN_WIDTH]bool

	// CHIP-8 has 16 8-bit data registers named from V0 to VF. The VF
	// register doubles as a carry flag.
	Registers [NUM_REGS]uint8
	// The address register, which is named I, is 16 bits wide and is used
	// with several opcodes that involve memory operations.
	IndexRegister uint16

	Keys [NUM_KEYS]bool

	DelayTimer uint8
	SoundTimer uint8
}

func NewCPU() *CPU {
	cpu := &CPU{
		ProgramCounter: START_ADDR,
		Ram:            [RAM_SIZE]uint8{},
		StackPointer:   0,
		Stack:          [STACK_SIZE]uint16{},
		screen:         [SCREEN_HEIGHT][SCREEN_WIDTH]bool{},
		Registers:      [NUM_REGS]uint8{},
		IndexRegister:  0,
		Keys:           [NUM_KEYS]bool{false},
		DelayTimer:     0,
		SoundTimer:     0,
	}
	copy(cpu.Ram[:FONTSET_SIZE], FONTSET[:])

	return cpu
}

func (c *CPU) Tick() (bool, bool, error) {
	op := c.GetOpCode()

	err := c.execute(op)

	c.TickTimers()

	return false, c.shouldBeep(), err
}

func (c *CPU) execute(op uint16) error {
	err := OpCode(op).Execute(c)
	return err
}

func (c *CPU) TickTimers() {
	if c.DelayTimer > 0 {
		c.DelayTimer -= 1
	}
	if c.SoundTimer > 0 {
		c.SoundTimer -= 1
	}
}

func (c *CPU) shouldBeep() bool {
	return c.SoundTimer == 1
}

func (c *CPU) GetOpCode() uint16 {
	high_byte := c.Ram[c.ProgramCounter]
	low_byte := c.Ram[c.ProgramCounter+1]
	c.ProgramCounter += 2
	return (uint16(high_byte) << 8) | uint16(low_byte)
}

func (c *CPU) Push(val uint16) {
	c.Stack[c.StackPointer] = val
	c.StackPointer += 1
}

func (c *CPU) Pop() uint16 {
	c.StackPointer -= 1
	return c.Stack[c.StackPointer]
}

func (c *CPU) SetKey(num uint8, isPressed bool) {
	c.Keys[num] = isPressed
}

func (c *CPU) Screen() [SCREEN_HEIGHT][SCREEN_WIDTH]bool {
	return c.screen
}

func (cpu *CPU) ClearScreen() {
	cpu.screen = [SCREEN_HEIGHT][SCREEN_WIDTH]bool{}
}
