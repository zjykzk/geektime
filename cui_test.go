package geektime

import (
	"sync"
	"testing"
	"time"
)

func TestProgressCUI(t *testing.T) {
	p := &progress{name: "label", total: 100}
	ui := newProgressCUI(100, 10, p)
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

	ps := []*progress{{name: "label", total: 200}, {name: "LABEL", total: 100}}

	bus.post(eventUIProgressTotal, ps)
	pg := func(index int) {
		p := ps[index]
		for i := int32(0); i < p.total; i++ {
			time.Sleep(time.Millisecond * 10)
			p.advance(1)
			bus.post(eventUIUpdateProgress, p)
		}
		wg.Done()
	}

	wg.Add(2)
	go pg(0)
	go pg(1)
	cui.run()
}
