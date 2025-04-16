package painter

import (
	"image"
	"image/color"
	"sync"
	"testing"
	"time"

	"golang.org/x/exp/shiny/screen"
	"golang.org/x/image/draw"
)

type mockReceiver struct {
	mu        sync.Mutex
	updates   []screen.Texture
	lastTex   screen.Texture
	updateCnt int
}

func (m *mockReceiver) Update(t screen.Texture) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.updates = append(m.updates, t)
	m.lastTex = t
	m.updateCnt++
}

type mockTexture struct{}

func (m *mockTexture) Release()                                                 {}
func (m *mockTexture) Size() image.Point                                        { return size }
func (m *mockTexture) Bounds() image.Rectangle                                  { return image.Rect(0, 0, size.X, size.Y) }
func (m *mockTexture) Upload(p image.Point, b screen.Buffer, r image.Rectangle) {}
func (m *mockTexture) Fill(r image.Rectangle, c color.Color, op draw.Op)        {}

type mockBuffer struct{}

func (m *mockBuffer) Release()                {}
func (m *mockBuffer) Size() image.Point       { return size }
func (m *mockBuffer) Bounds() image.Rectangle { return image.Rect(0, 0, size.X, size.Y) }
func (m *mockBuffer) RGBA() *image.RGBA       { return &image.RGBA{} }

type mockScreen struct{}

func (m *mockScreen) NewBuffer(size image.Point) (screen.Buffer, error) {
	return &mockBuffer{}, nil
}
func (m *mockScreen) NewTexture(size image.Point) (screen.Texture, error) {
	return &mockTexture{}, nil
}
func (m *mockScreen) NewWindow(*screen.NewWindowOptions) (screen.Window, error) {
	return nil, nil
}

type mockOperation struct {
	doFunc func(screen.Texture) bool
}

func (m *mockOperation) Do(t screen.Texture) bool {
	if m.doFunc != nil {
		return m.doFunc(t)
	}
	return false
}

func TestLoop_Start(t *testing.T) {
	l := &Loop{}
	scr := &mockScreen{}
	recv := &mockReceiver{}
	l.Receiver = recv

	l.Start(scr)

	if l.next == nil || l.prev == nil {
		t.Error("Loop should initialize next and prev textures")
	}

	l.StopAndWait()
}

func TestLoop_Post(t *testing.T) {
	l := &Loop{}
	scr := &mockScreen{}
	recv := &mockReceiver{}
	l.Receiver = recv

	l.Start(scr)
	defer l.StopAndWait()

	opCalled := false
	op := &mockOperation{
		doFunc: func(t screen.Texture) bool {
			opCalled = true
			return true
		},
	}

	l.Post(op)

	time.Sleep(100 * time.Millisecond)

	if !opCalled {
		t.Error("Posted operation was not executed")
	}

	if recv.updateCnt == 0 {
		t.Error("Receiver should have been updated after operation")
	}
}

func TestLoop_StopAndWait(t *testing.T) {
	l := &Loop{}
	scr := &mockScreen{}
	recv := &mockReceiver{}
	l.Receiver = recv

	l.Start(scr)

	stopped := make(chan struct{})

	slowOp := &mockOperation{
		doFunc: func(t screen.Texture) bool {
			time.Sleep(200 * time.Millisecond)
			return true
		},
	}
	l.Post(slowOp)

	go func() {
		time.Sleep(100 * time.Millisecond)
		l.StopAndWait()
		close(stopped)
	}()

	<-stopped

	if !l.stopReq {
		t.Error("Stop request flag should be set")
	}

	l.Post(&mockOperation{})

	if recv.updateCnt > 1 {
		t.Error("Only one operation should have been processed before stop")
	}
}

func TestLoop_OperationOrdering(t *testing.T) {
	l := &Loop{}
	scr := &mockScreen{}
	recv := &mockReceiver{}
	l.Receiver = recv

	l.Start(scr)
	defer l.StopAndWait()

	var executionOrder []int
	var mu sync.Mutex

	op1 := &mockOperation{
		doFunc: func(t screen.Texture) bool {
			mu.Lock()
			defer mu.Unlock()
			executionOrder = append(executionOrder, 1)
			return true
		},
	}

	op2 := &mockOperation{
		doFunc: func(t screen.Texture) bool {
			mu.Lock()
			defer mu.Unlock()
			executionOrder = append(executionOrder, 2)
			return true
		},
	}

	l.Post(op1)
	l.Post(op2)

	time.Sleep(100 * time.Millisecond)

	mu.Lock()
	defer mu.Unlock()
	if len(executionOrder) != 2 || executionOrder[0] != 1 || executionOrder[1] != 2 {
		t.Error("Operations should execute in the order they were posted")
	}
}

func TestLoop_NoUpdateWhenNotNeeded(t *testing.T) {
	l := &Loop{}
	scr := &mockScreen{}
	recv := &mockReceiver{}
	l.Receiver = recv

	l.Start(scr)
	defer l.StopAndWait()

	op := &mockOperation{
		doFunc: func(t screen.Texture) bool {
			return false
		},
	}

	l.Post(op)

	time.Sleep(100 * time.Millisecond)

	if recv.updateCnt > 0 {
		t.Error("Receiver should not be updated when operation returns false")
	}
}
