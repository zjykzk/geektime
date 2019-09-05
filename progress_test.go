package geektime

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestProgress(t *testing.T) {
	var size int32 = 10000
	p := progress{total: size}
	p.advance(10)
	assert.False(t, p.isEnd())

	p.advance(size - 10)
	assert.True(t, p.isEnd())
}

func TestCourseProgress(t *testing.T) {
	bus := &bus{}
	cp := newProgress(bus)

	cp.subscribeEvents()

	bus.post(eventCourse, courseRet{course{"TestCourseProgress", 2, 1}, nil})
	bus.post(eventArticles, articles{articles: []article{article{Title: "a1", ID: 1}, article{Title: "a2", ID: 2}}})
	bus.post(eventM3U8Parsed, m3u8{articleID: 1, ts: []string{"1"}})
	bus.post(eventM3U8Parsed, m3u8{articleID: 2, ts: []string{"2"}})
	assert.False(t, cp.isEnd())

	go func() {
		bus.post(eventDownloadTS, downloadTS{articleID: 2})
		bus.post(eventDownloadTS, downloadTS{articleID: 1})
	}()

	start := time.Now()
	for !cp.isEnd() {
		if time.Since(start) > time.Second {
			t.Fatal("timeout")
		}
	}
}
