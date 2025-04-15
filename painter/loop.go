package painter

import (
	"image"
	"sync"

	"golang.org/x/exp/shiny/screen"
)

type Receiver interface {
	Update(t screen.Texture)
}

type Loop struct {
	Receiver Receiver
	next     screen.Texture
	prev     screen.Texture
	mq       messageQueue
	stop     chan struct{}
	stopReq  bool
	wg       sync.WaitGroup
}

var size = image.Pt(400, 400)

func (l *Loop) Start(s screen.Screen) {
	l.next, _ = s.NewTexture(size)
	l.prev, _ = s.NewTexture(size)
	l.stop = make(chan struct{})
	l.mq = messageQueue{
		queue: make(chan Operation, 100),
	}

	l.wg.Add(1)
	go func() {
		defer l.wg.Done()
		l.eventLoop()
	}()
}

func (l *Loop) eventLoop() {
	for {
		select {
		case <-l.stop:
			return
		case op := <-l.mq.queue:
			if update := op.Do(l.next); update {
				l.Receiver.Update(l.next)
				l.next, l.prev = l.prev, l.next
			}
		}
	}
}

func (l *Loop) Post(op Operation) {
	if l.stopReq {
		return
	}
	l.mq.push(op)
}

func (l *Loop) StopAndWait() {
	if l.stopReq {
		return
	}
	l.stopReq = true
	close(l.stop)
	l.wg.Wait()
}

type messageQueue struct {
	queue chan Operation
}

func (mq *messageQueue) push(op Operation) {
	mq.queue <- op
}

func (mq *messageQueue) pull() Operation {
	return <-mq.queue
}

func (mq *messageQueue) empty() bool {
	return len(mq.queue) == 0
}
