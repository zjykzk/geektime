package geektimedl

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

func newProgressCUI(width int, p *progress) *progressCUI {
	labelWidth := width / 5
	return &progressCUI{
		label: p.name + strings.Repeat(" ", labelWidth-calcWidth(p.name)),

		progress:      p,
		progressWidth: width - labelWidth,
	}
}

func calcWidth(s string) (w int) {
	for _, c := range s {
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

	updateChan      chan string
	newProgressChan chan *progress
	exitChan        chan struct{}
}

func newCUI() *cui {
	width := 200
	ws, err := unix.IoctlGetWinsize(syscall.Stdout, unix.TIOCGWINSZ)
	if err == nil {
		width = int(ws.Col)

	}

	return &cui{
		width:           width,
		updateChan:      make(chan string, 1024),
		newProgressChan: make(chan *progress, 1024),
		exitChan:        make(chan struct{}),
	}
}

func (ui *cui) run() {
	hideCursor()

	go func() {
	OUT:
		for {
			var line int
			select {
			case <-ui.exitChan:
				close(ui.newProgressChan)
				close(ui.updateChan)
				break OUT
			case p := <-ui.newProgressChan:
				ui.puis = append(ui.puis, newProgressCUI(ui.width, p))
				line = len(ui.puis) - 1
			case name := <-ui.updateChan:
				line = ui.findLineNo(name)
			}

			ui.drawLine(line)
		}
	}()

	ui.wait()
	ui.drawLeft()
	ui.moveTo(len(ui.puis))
	showCursor()
}

func (ui *cui) findLineNo(name string) int {
	for i, p := range ui.puis {
		if p.progress.name == name {
			return i
		}
	}

	return -1
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
		break
	case <-make(chan os.Signal, 1):
		break
	}
}

func (ui *cui) drawLeft() {
	for p := range ui.newProgressChan {
		ui.puis = append(ui.puis, newProgressCUI(ui.width, p))
		ui.drawLine(len(ui.puis))
	}
	for n := range ui.updateChan {
		ui.drawLine(ui.findLineNo(n))
	}
}

func (ui *cui) subscribeEvents(bus *bus) {
	bus.subscribe(eventUINewProgress, func(v interface{}) {
		ui.newProgressChan <- v.(*progress)
	})
	bus.subscribe(eventUIUpdateProgress, func(v interface{}) {
		ui.updateChan <- v.(string)
	})
	bus.subscribe(eventUIProgressEnd, func(interface{}) {
		close(ui.exitChan)
	})
}
