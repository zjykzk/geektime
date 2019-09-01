package geektimedl

import (
	"errors"
	"fmt"
	"sync"
	"sync/atomic"
)

type progress struct {
	name           string
	total, current int32
	err            error
}

func (p *progress) String() string {
	return fmt.Sprintf(
		"progress:[name:%s,total:%d,current:%d,err:%v]", p.name, p.total, p.current, p.err,
	)
}

func (p *progress) advance(d int32) {
	atomic.AddInt32(&p.current, d)
}

func (p *progress) abort(err error) {
	p.err = err
}

func (p *progress) isEnd() bool {
	return p.err != nil || atomic.LoadInt32(&p.current) >= atomic.LoadInt32(&p.total)
}

type courseProgress struct {
	bus *bus

	videoProgresses sync.Map // video id -> *progress
	progress        progress

	videos sync.Map // video id -> articleVideo

	finished chan struct{}
}

func (cp *courseProgress) videoName(videoID string) (string, bool) {
	v, ok := cp.videos.Load(videoID)
	if !ok {
		return "", false
	}
	return v.(articleVideo).name, true
}

func (cp *courseProgress) isEnd() bool {
	ok := true
	cp.videoProgresses.Range(func(_, v interface{}) bool {
		if !v.(*progress).isEnd() {
			ok = false
			return false
		}
		return true
	})

	return ok
}

func (cp *courseProgress) await() {
	<-cp.finished
}

func (cp *courseProgress) advance(d int32) {
	cp.progress.advance(d)
}

func (cp *courseProgress) abort(err error) {
	cp.progress.abort(err)
	cp.tryEnd()
}

func (cp *courseProgress) subscribeEvents() {
	b := cp.bus
	b.subscribe(eventCourse, func(v interface{}) {
		c := v.(courseRet)
		if c.err != nil {
			cp.abort(wrapErr("fetch course", c.err))
			return
		}

		cp.progress = progress{name: c.course.Title, total: int32(c.course.ArticleCount)}
	})

	b.subscribe(eventArticles, func(v interface{}) {
		as := v.(articles)
		if as.err != nil {
			cp.abort(wrapErr("fetch articles", as.err))
			return
		}

		for _, a := range as.articles {
			p := &progress{name: a.Title, total: 1}
			cp.videoProgresses.Store(a.VideoID, p)
			cp.bus.post(eventUINewProgress, p)
		}
	})

	b.subscribe(eventArticleVideo, func(v interface{}) {
		av := v.(articleVideo)
		if av.err == nil {
			cp.videos.LoadOrStore(av.id, av)
			return
		}

		v, ok := cp.videoProgresses.Load(av.id)
		if !ok {
			panic("[BUG] vide progress not exist:" + av.id)
		}

		p := v.(*progress)
		p.abort(wrapErr("fetch video", av.err))
		cp.bus.post(eventUIUpdateProgress, p.name)

		cp.advance(1)
		cp.tryEnd()
	})

	b.subscribe(eventPlayAuth, func(v interface{}) {
		pa := v.(playAuth)
		if pa.err == nil {
			return
		}

		v, ok := cp.videoProgresses.Load(pa.videoID)
		if !ok {
			panic("[BUG] vide progress not exist:" + pa.videoID)
		}

		p := v.(*progress)
		p.abort(wrapErr("play auth", pa.err))
		cp.bus.post(eventUIUpdateProgress, p.name)

		cp.advance(1)
		cp.tryEnd()
	})

	b.subscribe(eventDownloadTS, func(v interface{}) {
		videoID := v.(downloadTS).videoID
		v, ok := cp.videoProgresses.Load(videoID)
		if !ok {
			panic("[BUG] vide progress not exist:" + videoID)
		}

		p := v.(*progress)
		p.advance(1)
		cp.bus.post(eventUIUpdateProgress, p.name)

		if p.isEnd() {
			cp.advance(1)
			cp.tryEnd()
		}
	})

	b.subscribe(eventM3U8Parsed, func(v interface{}) {
		m := v.(m3u8)
		v, ok := cp.videoProgresses.Load(m.videoID)
		if !ok {
			panic("[BUG] vide progress not exist:" + m.videoID)
		}

		p := v.(*progress)
		p.name, p.total = m.name, int32(len(m.ts))
		cp.bus.post(eventUIUpdateProgress, p.name)
	})

	b.subscribe(eventCreateVideoFoldFailed, func(v interface{}) {
		v, ok := cp.videoProgresses.Load(v)
		if !ok {
			panic("[BUG] vide progress not exist:" + v.(string))
		}

		p := v.(*progress)
		p.abort(errors.New("create video fold failed"))
		cp.bus.post(eventUIUpdateProgress, p.name)
	})
}

func (cp *courseProgress) tryEnd() {
	if cp.progress.isEnd() {
		cp.bus.post(eventUIProgressEnd, nil)
		close(cp.finished)
	}
}

func newProgress(bus *bus) *courseProgress {
	return &courseProgress{bus: bus, finished: make(chan struct{})}
}

func wrapErr(msg string, err error) error {
	return fmt.Errorf(msg+" error:%s", err)
}
