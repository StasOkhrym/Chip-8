package emulator

import (
	"chip-8-go/cpu"
	"fmt"
	"os"

	sdl "github.com/veandco/go-sdl2/sdl"
)

type Chip8 struct {
	beeper *Beeper
	cpu    *cpu.CPU
	fqHz   uint32

	renderer      *sdl.Renderer
	scaleModifier int32
}

func InitChip8(fileName string, scaleModifier int32, renderer *sdl.Renderer) (*Chip8, error) {
	beeper, err := NewBeeper()
	if err != nil {
		return nil, err
	}

	cpu := cpu.NewCPU()

	c8 := &Chip8{
		beeper:        beeper,
		cpu:           cpu,
		fqHz:          1000 / 60, // 60Hz
		scaleModifier: scaleModifier,
		renderer:      renderer,
	}
	loadErr := c8.LoadProgram(fileName)
	if loadErr != nil {
		return nil, loadErr
	}

	return c8, nil

}

func (c *Chip8) Run() error {
	for {
		draw, beep, err := c.cpu.Tick()
		if err != nil {
			return err
		}

		if beep {
			c.Beep()
		}

		if draw {
			c.Draw()
		}

		c.pollKeyPad()
		sdl.Delay(c.fqHz)
	}
}

func (c *Chip8) pollKeyPad() {
	for event := sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
		isPressed := isKeyPressed(event)
		switch et := event.(type) {
		case *sdl.QuitEvent:
			os.Exit(0)
		case *sdl.KeyboardEvent:
			switch et.Keysym.Sym {
			case sdl.K_1:
				c.cpu.SetKey(0x1, isPressed)
			case sdl.K_2:
				c.cpu.SetKey(0x2, isPressed)
			case sdl.K_3:
				c.cpu.SetKey(0x3, isPressed)
			case sdl.K_4:
				c.cpu.SetKey(0xC, isPressed)
			case sdl.K_q:
				c.cpu.SetKey(0x4, isPressed)
			case sdl.K_w:
				c.cpu.SetKey(0x5, isPressed)
			case sdl.K_e:
				c.cpu.SetKey(0x6, isPressed)
			case sdl.K_r:
				c.cpu.SetKey(0xD, isPressed)
			case sdl.K_a:
				c.cpu.SetKey(0x7, isPressed)
			case sdl.K_s:
				c.cpu.SetKey(0x8, isPressed)
			case sdl.K_d:
				c.cpu.SetKey(0x9, isPressed)
			case sdl.K_f:
				c.cpu.SetKey(0xE, isPressed)
			case sdl.K_z:
				c.cpu.SetKey(0xA, isPressed)
			case sdl.K_x:
				c.cpu.SetKey(0x0, isPressed)
			case sdl.K_c:
				c.cpu.SetKey(0xB, isPressed)
			case sdl.K_v:
				c.cpu.SetKey(0xF, isPressed)
			}
		}
	}
}

func (c *Chip8) Draw() {
	c.renderer.SetDrawColor(0, 0, 0, 255)
	c.renderer.Clear()

	for j := 0; j < len(c.cpu.Screen); j++ {
		for i := 0; i < len(c.cpu.Screen[j]); i++ {
			if c.cpu.Screen[j][i] {
				c.renderer.SetDrawColor(0, 255, 0, 255)
			} else {
				c.renderer.SetDrawColor(0, 0, 0, 255)
			}
			c.renderer.FillRect(
				&sdl.Rect{
					Y: int32(j) * c.scaleModifier,
					X: int32(i) * c.scaleModifier,
					W: c.scaleModifier,
					H: c.scaleModifier,
				},
			)
		}
	}

	c.renderer.Present()
}

func (c *Chip8) Beep() {
	c.beeper.Play()
}

func isKeyPressed(event sdl.Event) bool {
	var isPressed bool

	if event.GetType() == sdl.KEYDOWN {
		isPressed = true
	} else {
		isPressed = false
	}

	return isPressed
}

func (c *Chip8) LoadProgram(fileName string) error {
	file, fileErr := os.OpenFile(fileName, os.O_RDONLY, 0777)
	if fileErr != nil {
		return fileErr
	}
	defer file.Close()

	stat, statErr := file.Stat()
	if statErr != nil {
		return statErr
	}
	if int64(len(c.cpu.Memory)-cpu.START_ADDR) < stat.Size() {
		return fmt.Errorf("ROM file size is too big")
	}

	buffer := make([]byte, stat.Size())
	if _, readErr := file.Read(buffer); readErr != nil {
		return readErr
	}

	copy(c.cpu.Memory[cpu.START_ADDR:], buffer)
	return nil
}
