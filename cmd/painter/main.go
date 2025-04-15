package main

import (
	"fmt"
	"image/color"
	"net/http"

	"github.com/Homa4/architecture-lab-3/painter"
	"github.com/Homa4/architecture-lab-3/painter/lang"
	"github.com/Homa4/architecture-lab-3/ui"
)

func main() {
	state := &painter.State{
		BackgroundColor: color.Black,
		BGRects:         make([][5]float32, 0),
		Figures:         make([]painter.Figure, 0),
		OffsetX:         0,
		OffsetY:         0,
		LastOp:          "",
	}

	parser := lang.NewParser(state)
	opLoop := painter.Loop{}
	pv := ui.Visualizer{}
	pv.Title = "Simple painter"

	pv.OnScreenReady = opLoop.Start
	opLoop.Receiver = &pv

	go func() {
		http.Handle("/", lang.HttpHandler(&opLoop, parser))
		fmt.Println("Server running at http://localhost:17000")
		err := http.ListenAndServe("localhost:17000", nil)
		if err != nil {
			fmt.Printf("HTTP server error: %v\n", err)
		}
	}()

	fmt.Println("Starting painter application")
	pv.Main()

	opLoop.StopAndWait()
	fmt.Println("Application stopped")
}
