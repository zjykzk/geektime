package geektimedl

import (
	"sync"
	"testing"
	"time"
)

func TestProgressCUI(t *testing.T) {
	p := &progress{name: "label", total: 100}
	ui := newProgressCUI(100, p)
	t.Log(ui.content())
	p.advance(10)
	t.Log(ui.content())
	p.advance(10)
	t.Log(ui.content())
}

func TestCUI(t *testing.T) {
	cui, bus := newCUI(), &bus{}
	cui.subscribeEvents(bus)

	var wg sync.WaitGroup

	pg := func(name string, total int32) {
		p := &progress{name: name, total: total}
		bus.post(eventUINewProgress, p)
		for i := int32(0); i < total; i++ {
			time.Sleep(time.Millisecond * 10)
			p.advance(1)
			bus.post(eventUIUpdateProgress, nil)
		}
		wg.Done()
	}

	wg.Add(2)
	go pg("label", 200)
	go pg("LABEL", 100)
	go func() {
		wg.Wait()
		bus.post(eventUIProgressEnd, nil)
	}()
	cui.run()
}
