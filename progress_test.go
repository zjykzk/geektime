package geektimedl

import (
	"testing"

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

	cp.subscribeEvents(bus)

	bus.post(eventArticles, articles{articles: []article{article{VideoID: "1"}, article{VideoID: "2"}}})
	bus.post(eventM3U8Parsed, m3u8{videoID: "1", ts: []string{"1"}})
	bus.post(eventM3U8Parsed, m3u8{videoID: "2", ts: []string{"2"}})
	assert.False(t, cp.isEnd())

	go func() {
		bus.post(eventDownloadTS, downloadTS{videoID: "2"})
		bus.post(eventDownloadTS, downloadTS{videoID: "1"})
	}()

	cp.await()
}
