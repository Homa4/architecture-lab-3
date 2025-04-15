package painter

import (
	"image"
	"image/color"

	"golang.org/x/exp/shiny/screen"
)

// Operation змінює вхідну текстуру.
type Operation interface {
	// Do виконує зміну операції, повертаючи true, якщо текстура вважається готовою для відображення.
	Do(t screen.Texture) (ready bool)
}

// OperationList групує список операції в одну.
type OperationList []Operation

type Figure struct {
	X, Y float32
}

// State зберігає стан малювання.
type State struct {
	BackgroundColorFill OperationFunc // Поточне зафарбування фону
	BGRect              *[4]float32   // Прямокутник (x1, y1, x2, y2)
	Figures             []Figure      // Список фігур
	OffsetX             float32       // Зміщення по X
	OffsetY             float32       // Зміщення по Y
}

func (ol OperationList) Do(t screen.Texture) (ready bool) {
	for _, o := range ol {
		ready = o.Do(t) || ready
	}
	return
}

// UpdateOp операція, яка не змінює текстуру, але сигналізує, що текстуру потрібно розглядати як готову.
var UpdateOp = updateOp{}

type updateOp struct{}

func (op updateOp) Do(t screen.Texture) bool { return true }

// OperationFunc використовується для перетворення функції оновлення текстури в Operation.
type OperationFunc func(t screen.Texture)

func (f OperationFunc) Do(t screen.Texture) bool {
	f(t)
	return false
}

// WhiteFill зафарбовує тестуру у білий колір. Може бути викоистана як Operation через OperationFunc(WhiteFill).
func WhiteFill(t screen.Texture) {
	t.Fill(t.Bounds(), color.White, screen.Src)
}

// GreenFill зафарбовує тестуру у зелений колір. Може бути викоистана як Operation через OperationFunc(GreenFill).
func GreenFill(t screen.Texture) {
	t.Fill(t.Bounds(), color.RGBA{G: 0xff, A: 0xff}, screen.Src)
}

type bgRectOp struct {
	state          *State
	X1, Y1, X2, Y2 float32
}

func NewBGRectOp(state *State, x1, y1, x2, y2 float32) Operation {
	return &bgRectOp{state, x1, y1, x2, y2}
}

func (op *bgRectOp) Do(t screen.Texture) bool {
	op.state.BGRect = &[4]float32{op.X1, op.Y1, op.X2, op.Y2}
	rect := t.Bounds()
	w := float32(rect.Dx())
	h := float32(rect.Dy())

	x1 := int(op.X1 * w)
	y1 := int(op.Y1 * h)
	x2 := int(op.X2 * w)
	y2 := int(op.Y2 * h)

	t.Fill(image.Rect(x1, y1, x2, y2), color.Black, screen.Src)
	return false
}

type figureOp struct {
	state *State
	X, Y  float32
}

func NewFigureOp(state *State, x, y float32) Operation {
	return &figureOp{state, x, y}
}

func (op *figureOp) Do(t screen.Texture) bool {
	op.state.Figures = append(op.state.Figures, Figure{op.X, op.Y})

	rect := t.Bounds()
	w := float32(rect.Dx())
	h := float32(rect.Dy())

	cx := int((op.X + op.state.OffsetX) * w)
	cy := int((op.Y + op.state.OffsetY) * h)

	size := 10
	t.Fill(image.Rect(cx-size, cy-size, cx+size, cy+size), color.RGBA{R: 255, A: 255}, screen.Src)
	return false
}

type moveOp struct {
	state *State
	X, Y  float32
}

func NewMoveOp(state *State, x, y float32) Operation {
	return &moveOp{state, x, y}
}

func (op *moveOp) Do(t screen.Texture) bool {
	op.state.OffsetX = op.X
	op.state.OffsetY = op.Y
	return false
}

type resetOp struct {
	state *State
}

func NewResetOp(state *State) Operation {
	return &resetOp{state: state}
}

func (op *resetOp) Do(t screen.Texture) bool {
	// Очистити стан
	op.state.BackgroundColorFill = OperationFunc(func(t screen.Texture) {
		t.Fill(t.Bounds(), color.Black, screen.Src)
	})
	op.state.BGRect = nil
	op.state.Figures = nil
	op.state.OffsetX = 0
	op.state.OffsetY = 0

	// Виконати чорне зафарбування
	op.state.BackgroundColorFill(t)
	return false
}
