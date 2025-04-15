package main

import (
	"os"
	"strings"

	"github.com/Homa4/architecture-lab-3/lang"
	"github.com/Homa4/architecture-lab-3/painter"
	"golang.org/x/exp/shiny/screen"
	"golang.org/x/exp/shiny/driver"
)

// func main() {
// 	var (
// 		pv ui.Visualizer // Візуалізатор створює вікно та малює у ньому.

// 		// Потрібні для частини 2.
// 		opLoop painter.Loop // Цикл обробки команд.
// 		parser lang.Parser  // Парсер команд.
// 	)

// 	//pv.Debug = true
// 	pv.Title = "Simple painter"

// 	pv.OnScreenReady = opLoop.Start
// 	opLoop.Receiver = &pv

// 	go func() {
// 		http.Handle("/", lang.HttpHandler(&opLoop, &parser))
// 		_ = http.ListenAndServe("localhost:17000", nil)
// 	}()
// 	fmt.Printf("localhost:17000\n")
// 	pv.Main()
// 	opLoop.StopAndWait()
// }

func main() {
	driver.Main(func(s screen.Screen) {
		// Створення вікна
		w, err := s.NewWindow(&screen.NewWindowOptions{
			Width:  600,
			Height: 600,
			Title:  "Painter",
		})
		if err != nil {
			panic(err)
		}
		defer w.Release()

		// Створення текстури
		t, err := s.NewTexture(w.Bounds().Size())
		if err != nil {
			panic(err)
		}
		defer t.Release()

		// Підготуємо скрипт
		script := `
reset
bgrect 0.1 0.1 0.9 0.9
figure 0.5 0.5
move 0.1 0.1
update
`

		// Ініціалізація стану і парсера
		state := &painter.State{}
		parser := lang.Parser{State: state}

		ops, err := parser.Parse(strings.NewReader(script))
		if err != nil {
			panic(err)
		}

		// Виконання операцій
		for _, op := range ops {
			if op.Do(t) {
				w.Upload(w.Bounds().Min, t, t.Bounds())
				w.Publish()
			}
		}

		// Чекаємо закриття вікна
		for {
			e := w.NextEvent()
			switch e := e.(type) {
			case screen.CloseEvent:
				return
			case screen.MouseEvent:
				// За бажанням — додай реакцію на клік
			}
		}
	})
}
