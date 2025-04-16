package painter

import (
	"image"
	"image/color"

	"golang.org/x/exp/shiny/screen"
)

type Operation interface {
	Do(t screen.Texture) (ready bool)
}

type OperationList []Operation

type Figure struct {
	X, Y  float32
	Color color.Color
}

type State struct {
	BackgroundColor color.Color
	BGRects         [][5]float32
	Figures         []Figure
	OffsetX         float32
	OffsetY         float32
	LastOp          string
}

func (ol OperationList) Do(t screen.Texture) (ready bool) {
	for _, o := range ol {
		ready = o.Do(t) || ready
	}
	return
}

var UpdateOp = updateOp{}

type updateOp struct{}

func (op updateOp) Do(t screen.Texture) bool { return true }

type OperationFunc func(t screen.Texture)

func (f OperationFunc) Do(t screen.Texture) bool {
	f(t)
	return false
}

func WhiteFill(t screen.Texture) {
	t.Fill(t.Bounds(), color.White, screen.Src)
}

func GreenFill(t screen.Texture) {
	t.Fill(t.Bounds(), color.RGBA{G: 0xff, A: 0xff}, screen.Src)
}

type backgroundOp struct {
	state *State
	color color.Color
}

func NewBackgroundOp(state *State, c color.Color) Operation {
	return &backgroundOp{state: state, color: c}
}

func (op *backgroundOp) Do(t screen.Texture) bool {
	op.state.BackgroundColor = op.color
	op.state.LastOp = "background"

	t.Fill(t.Bounds(), op.color, screen.Src)

	for _, rect := range op.state.BGRects {
		drawRect(t, rect[0], rect[1], rect[2], rect[3], color.Black)
	}

	for _, fig := range op.state.Figures {
		drawFigure(t, fig.X+op.state.OffsetX, fig.Y+op.state.OffsetY, fig.Color)
	}

	return false
}

type bgRectOp struct {
	state          *State
	X1, Y1, X2, Y2 float32
}

func NewBGRectOp(state *State, x1, y1, x2, y2 float32) Operation {
	return &bgRectOp{state, x1, y1, x2, y2}
}

func (op *bgRectOp) Do(t screen.Texture) bool {
	op.state.BGRects = append(op.state.BGRects, [5]float32{op.X1, op.Y1, op.X2, op.Y2, 0})
	op.state.LastOp = "bgrect"

	drawRect(t, op.X1, op.Y1, op.X2, op.Y2, color.Black)
	return false
}

func drawRect(t screen.Texture, x1, y1, x2, y2 float32, c color.Color) {
	rect := t.Bounds()
	w := float32(rect.Dx())
	h := float32(rect.Dy())

	ix1 := int(x1 * w)
	iy1 := int(y1 * h)
	ix2 := int(x2 * w)
	iy2 := int(y2 * h)

	t.Fill(image.Rect(ix1, iy1, ix2, iy2), c, screen.Src)
}

type figureOp struct {
	state *State
	X, Y  float32
	Color color.Color
}

func NewFigureOp(state *State, x, y float32, c color.Color) Operation {
	return &figureOp{state, x, y, c}
}

func (op *figureOp) Do(t screen.Texture) bool {
	op.state.Figures = append(op.state.Figures, Figure{op.X, op.Y, op.Color})
	op.state.LastOp = "figure"

	drawFigure(t, op.X+op.state.OffsetX, op.Y+op.state.OffsetY, op.Color)
	return false
}

func drawFigure(t screen.Texture, x, y float32, c color.Color) {
	rect := t.Bounds()
	w := float32(rect.Dx())
	h := float32(rect.Dy())

	cx := int(x * w)
	cy := int(y * h)

	scaleFactor := float32(1.0)
	crossSize := int(scaleFactor * 100)
	armThickness := int(scaleFactor * 50)

	t.Fill(image.Rect(
		cx-crossSize/2, cy-armThickness/2,
		cx+crossSize/2, cy+armThickness/2,
	), c, screen.Src)

	t.Fill(image.Rect(
		cx-armThickness/2, cy-crossSize/2,
		cx+armThickness/2, cy+crossSize/2,
	), c, screen.Src)
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
	op.state.LastOp = "move"

	if op.state.BackgroundColor != nil {
		t.Fill(t.Bounds(), op.state.BackgroundColor, screen.Src)
	} else {
		t.Fill(t.Bounds(), color.Black, screen.Src)
	}

	for _, rect := range op.state.BGRects {
		drawRect(t, rect[0], rect[1], rect[2], rect[3], color.Black)
	}

	for _, fig := range op.state.Figures {
		drawFigure(t, fig.X+op.state.OffsetX, fig.Y+op.state.OffsetY, fig.Color)
	}

	return false
}

type resetOp struct {
	state *State
}

func NewResetOp(state *State) Operation {
	return &resetOp{state: state}
}

func (op *resetOp) Do(t screen.Texture) bool {
	op.state.BackgroundColor = color.Black
	op.state.BGRects = nil
	op.state.Figures = nil
	op.state.OffsetX = 0
	op.state.OffsetY = 0
	op.state.LastOp = "reset"

	t.Fill(t.Bounds(), color.Black, screen.Src)
	return false
}
