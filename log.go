package geektimedl

import (
	"fmt"
	"io"
	"os"
)

type logger struct {
	writer io.WriteCloser
}

func newLogger(name string) (*logger, error) {
	f, err := os.Create(name)
	if err != nil {
		return nil, err
	}
	return &logger{f}, nil
}

func (l *logger) subscribeEvents(bus *bus) {
	for i := eventCourse; i < eventCount; i++ {
		bus.subscribe(i, l.print)
	}
}

func (l *logger) print(v interface{}) {
	l.writer.Write([]byte(fmt.Sprintf("[DEBUG] %s\n", v)))
}

func (l *logger) shutdown() {
	l.writer.Close()
}
