package geektime

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

	articleProgresses sync.Map // artice id -> *progress
	progress          progress

	articles sync.Map // article id -> articleVideo
}

func (cp *courseProgress) isEnd() bool {
	ok := true
	cp.articleProgresses.Range(func(_, v interface{}) bool {
		if !v.(*progress).isEnd() {
			ok = false
			return false
		}
		return true
	})

	return ok
}

func (cp *courseProgress) advance(d int32) {
	cp.progress.advance(d)
}

func (cp *courseProgress) abort(err error) {
	cp.progress.abort(err)
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

		ps := make([]*progress, len(as.articles))
		for i, a := range as.articles {
			p := &progress{name: a.Title, total: 1}
			cp.articleProgresses.Store(a.ID, p)
			ps[i] = p
		}
		b.post(eventUIProgressTotal, ps)
	})

	b.subscribe(eventArticleFinished, func(v interface{}) {
		articleID := v.(int)
		v, ok := cp.articleProgresses.Load(v)
		if !ok {
			panic(fmt.Sprintf("[BUG] article progress not exist:%d", articleID))
		}

		p := v.(*progress)
		p.advance(1)
		b.post(eventUIUpdateProgress, p)

		cp.advance(1)
	})

	b.subscribe(eventArticleVideo, func(v interface{}) {
		av := v.(articleVideo)
		if av.err == nil {
			cp.articles.LoadOrStore(av.articleID, av)
			return
		}

		v, ok := cp.articleProgresses.Load(av.articleID)
		if !ok {
			panic(fmt.Sprintf("[BUG] article progress not exist:%d", av.articleID))
		}

		p := v.(*progress)
		p.abort(wrapErr("fetch article video", av.err))
		b.post(eventUIUpdateProgress, p)

		cp.advance(1)
	})

	b.subscribe(eventPlayAuth, func(v interface{}) {
		pa := v.(playAuth)
		if pa.err == nil {
			return
		}

		v, ok := cp.articleProgresses.Load(pa.articleID)
		if !ok {
			panic(fmt.Sprintf("[BUG] article progress not exist:%d", pa.articleID))
		}

		p := v.(*progress)
		p.abort(wrapErr("play auth", pa.err))
		b.post(eventUIUpdateProgress, p)

		cp.advance(1)
	})

	b.subscribe(eventDownloadTS, func(v interface{}) {
		articleID := v.(downloadTS).articleID
		v, ok := cp.articleProgresses.Load(articleID)
		if !ok {
			panic(fmt.Sprintf("[BUG] article progress not exist:%d", articleID))
		}

		p := v.(*progress)
		p.advance(1)
		b.post(eventUIUpdateProgress, p)

		if p.isEnd() {
			cp.advance(1)
		}
	})

	b.subscribe(eventM3U8Parsed, func(v interface{}) {
		m := v.(m3u8)
		v, ok := cp.articleProgresses.Load(m.articleID)
		if !ok {
			panic(fmt.Sprintf("[BUG] article progress not exist:%s", m))
		}

		p := v.(*progress)
		if m.err == nil {
			p.name, p.total = m.name, int32(len(m.ts))
		} else {
			p.abort(fmt.Errorf("m3u8 parsed error:%s", m.err))
			cp.advance(1)
		}

		b.post(eventUIUpdateProgress, p)
	})

	b.subscribe(eventCreateArticleFoldFailed, func(v interface{}) {
		v, ok := cp.articleProgresses.Load(v)
		if !ok {
			panic(fmt.Sprintf("[BUG] article progress not exist:%s", v))
		}

		p := v.(*progress)
		p.abort(errors.New("create article fold failed"))
		b.post(eventUIUpdateProgress, p)

		cp.advance(1)
	})
}

func newProgress(bus *bus) *courseProgress {
	return &courseProgress{bus: bus}
}

func wrapErr(msg string, err error) error {
	return fmt.Errorf(msg+" error:%s", err)
}
