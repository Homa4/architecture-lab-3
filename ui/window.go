package ui

import (
	"image"
	"image/color"
	"log"

	"golang.org/x/exp/shiny/driver"
	"golang.org/x/exp/shiny/screen"
	"golang.org/x/image/draw"
	"golang.org/x/mobile/event/key"
	"golang.org/x/mobile/event/lifecycle"
	"golang.org/x/mobile/event/mouse"
	"golang.org/x/mobile/event/paint"
	"golang.org/x/mobile/event/size"
)

type Visualizer struct {
	Title         string
	Debug         bool
	OnScreenReady func(s screen.Screen)
	w             screen.Window
	tx            chan screen.Texture
	done          chan struct{}
	sz            image.Rectangle
	crossCenter   image.Point
}

const (
	crossSize    = 200
	armThickness = 50
)

func (pw *Visualizer) Main() {
	pw.tx = make(chan screen.Texture)
	pw.done = make(chan struct{})
	driver.Main(pw.run)
}

func (pw *Visualizer) Update(t screen.Texture) {
	pw.tx <- t
}

func (pw *Visualizer) run(s screen.Screen) {
	w, err := s.NewWindow(&screen.NewWindowOptions{
		Title:  pw.Title,
		Width:  800,
		Height: 800,
	})
	if err != nil {
		log.Fatal("Failed to initialize the app window:", err)
	}
	defer func() {
		w.Release()
		close(pw.done)
	}()

	if pw.OnScreenReady != nil {
		pw.OnScreenReady(s)
	}

	pw.w = w
	pw.crossCenter = image.Point{X: 400, Y: 400}

	events := make(chan any)
	go func() {
		for {
			e := w.NextEvent()
			if pw.Debug {
				log.Printf("new event: %v", e)
			}
			if detectTerminate(e) {
				close(events)
				break
			}
			events <- e
		}
	}()

	var t screen.Texture
	for {
		select {
		case e, ok := <-events:
			if !ok {
				return
			}
			pw.handleEvent(e, t)
		case t = <-pw.tx:
			w.Send(paint.Event{})
		}
	}
}

func detectTerminate(e any) bool {
	switch e := e.(type) {
	case lifecycle.Event:
		if e.To == lifecycle.StageDead {
			return true
		}
	case key.Event:
		if e.Code == key.CodeEscape {
			return true
		}
	}
	return false
}

func (pw *Visualizer) handleEvent(e any, t screen.Texture) {
	switch e := e.(type) {
	case size.Event:
		pw.sz = e.Bounds()
		pw.crossCenter = image.Point{
			X: pw.sz.Min.X + pw.sz.Dx()/2,
			Y: pw.sz.Min.Y + pw.sz.Dy()/2,
		}
	case error:
		log.Printf("ERROR: %s", e)
	case mouse.Event:
		if e.Direction == mouse.DirPress && e.Button == mouse.ButtonLeft {
			pw.crossCenter = image.Point{X: int(e.X), Y: int(e.Y)}
			pw.w.Send(paint.Event{})
		}
	case paint.Event:
		if t == nil {
			pw.drawDefaultUI()
		} else {
			pw.w.Scale(pw.sz, t, t.Bounds(), draw.Src, nil)
		}
		pw.w.Publish()
	}
}

func (pw *Visualizer) drawDefaultUI() {
	bgColor := color.RGBA{R: 0, G: 255, B: 0, A: 255}
	pw.w.Fill(pw.sz, bgColor, draw.Src)
	figureColor := color.RGBA{R: 0, G: 0, B: 255, A: 255}
	horz := image.Rect(
	 pw.crossCenter.X-crossSize/2,
	 pw.crossCenter.Y-armThickness/2,
	 pw.crossCenter.X+crossSize/2,
	 pw.crossCenter.Y+armThickness/2,
	)
	pw.w.Fill(horz, figureColor, draw.Src)
	vert := image.Rect(
	 pw.crossCenter.X-armThickness/2,
	 pw.crossCenter.Y-crossSize/2,
	 pw.crossCenter.X+armThickness/2,
	 pw.crossCenter.Y+crossSize/2,
	)
	pw.w.Fill(vert, figureColor, draw.Src)
   }
