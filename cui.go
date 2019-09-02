package geektime

import (
	"fmt"
	"os"
	"strings"
	"syscall"
	"unicode/utf8"

	"golang.org/x/sys/unix"
)

type progressCUI struct {
	label string

	progress      *progress
	progressWidth int
}

func newProgressCUI(width, nameWidth int, p *progress) *progressCUI {
	name := simplify(p.name)
	return &progressCUI{
		label:         name + strings.Repeat(" ", nameWidth-calcWidth(simplify(name))+1),
		progress:      p,
		progressWidth: width - nameWidth - 1,
	}
}

func calcWidth(s string) (w int) {
	for _, c := range s {
		if c == '“' || c == '”' {
			return 2
		}

		switch l := utf8.RuneLen(c); l {
		case 1, 2:
			w++
		default:
			w += 2
		}
	}
	return
}

func (ui *progressCUI) content() string {
	p := ui.progress
	if p.err != nil {
		return fmt.Sprintf("%s\x1b[31mFailed\x1b[0m:%s", ui.label, p.err)
	}

	width := ui.progressWidth - 7

	percent := float64(p.current) / float64(p.total)
	finishedCount := int(float64(width) * percent)

	r := strings.Repeat
	return fmt.Sprintf(
		"%s%3d%%[%s>%s]", ui.label, int(100*percent), r("=", finishedCount), r(" ", width-finishedCount),
	)
}

type cui struct {
	puis []*progressCUI

	width       int
	lineCount   int
	currentLine int

	progressChan chan *progress
	exitChan     chan struct{}
}

func newCUI() *cui {
	width := 200
	ws, err := unix.IoctlGetWinsize(syscall.Stdout, unix.TIOCGWINSZ)
	if err == nil {
		width = int(ws.Col)
	}

	return &cui{
		width:        width,
		progressChan: make(chan *progress, 1024),
		exitChan:     make(chan struct{}),
	}
}

func (ui *cui) run() {
	hideCursor()
	defer showCursor()

	go func() {
		for {
		OUT:
			select {
			case <-ui.exitChan:
				return
			case p, ok := <-ui.progressChan:
				if !ok {
					return
				}

				ui.drawLine(ui.findLineNo(p))
				for _, p := range ui.puis {
					if !p.progress.isEnd() {
						break OUT
					}
				}
				close(ui.exitChan)
			}

		}
	}()

	ui.wait()
	ui.drawLeft()
	ui.moveTo(len(ui.puis))
}

func (ui *cui) findLineNo(p *progress) int {
	for i, p0 := range ui.puis {
		if p0.progress == p {
			return i
		}
	}
	panic("cannot find progress:" + p.String())
}

func (ui *cui) moveTo(lineNo int) {
	switch {
	case lineNo > ui.currentLine:
		moveDown(lineNo - ui.currentLine)
	case lineNo < ui.currentLine:
		moveUp(ui.currentLine - lineNo)
	}
}

func (ui *cui) drawLine(lineNo int) {
	ui.moveTo(lineNo)
	ereaseCurrentLine()
	fmt.Println(ui.puis[lineNo].content())
	ui.currentLine = lineNo + 1
}

func (ui *cui) wait() {
	select {
	case <-ui.exitChan:
		close(ui.progressChan)
	case <-make(chan os.Signal, 1):
	}
}

func (ui *cui) drawLeft() {
	for p := range ui.progressChan {
		ui.drawLine(ui.findLineNo(p))
	}
}

func (ui *cui) subscribeEvents(bus *bus) {
	bus.subscribe(eventUIUpdateProgress, func(v interface{}) {
		ui.progressChan <- v.(*progress)
	})

	bus.subscribe(eventUIProgressTotal, func(v interface{}) {
		ps := v.([]*progress)
		ui.puis = make([]*progressCUI, len(ps))
		width := calcMaxWidth(ps)
		for i, p := range ps {
			ui.puis[i] = newProgressCUI(ui.width, width, p)
			ui.drawLine(i)
		}
	})
}

func calcMaxWidth(ps []*progress) (maxWidth int) {
	for _, p := range ps {
		w := calcWidth(p.name)
		if maxWidth < w {
			maxWidth = w
		}
	}
	return
}
