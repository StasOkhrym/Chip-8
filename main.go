package main

import (
	"chip-8-go/cpu"
	"chip-8-go/emulator"
	"fmt"
	"os"

	"github.com/veandco/go-sdl2/sdl"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Print("Please provide file")
		os.Exit(0)
	}

	fileName := os.Args[1]

	if sdlErr := sdl.Init(sdl.INIT_EVERYTHING); sdlErr != nil {
		panic(sdlErr)
	}
	defer sdl.Quit()

	var scaleModifier int32 = 10 //Adjust screen size

	window, windowErr := sdl.CreateWindow(
		"Chip 8 - "+fileName,
		sdl.WINDOWPOS_UNDEFINED,
		sdl.WINDOWPOS_UNDEFINED,
		cpu.SCREEN_WIDTH*scaleModifier,
		cpu.SCREEN_HEIGHT*scaleModifier,
		sdl.WINDOW_SHOWN,
	)
	if windowErr != nil {
		panic(windowErr)
	}
	defer window.Destroy()

	renderer, rendererErr := sdl.CreateRenderer(window, -1, 0)
	if rendererErr != nil {
		panic(rendererErr)
	}
	defer renderer.Destroy()

	c8, err := emulator.InitChip8(fileName, scaleModifier, renderer)
	if err != nil {
		panic(err)
	}

	runErr := c8.Run()
	if runErr != nil {
		panic(runErr)
	}
}
