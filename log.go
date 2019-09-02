package geektime

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
)

type logger struct {
	writer io.WriteCloser
}

func newLogger(name string) (*logger, error) {
	err := makeSureDirExist(filepath.Dir(name))
	if err != nil {
		return nil, err
	}

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
