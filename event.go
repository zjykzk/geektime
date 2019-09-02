package geektime

import (
	"fmt"
	"strings"
)

type event int

const (
	// eventCourse course fetched
	eventCourse event = iota
	// eventArticles articles is fetched
	eventArticles
	eventArticleFinished
	// eventArticleVideo get the video information of the article
	eventArticleVideo
	// eventPlayAuth get the play auth
	eventPlayAuth
	// eventPlayList get the play list which contains video url
	eventPlayList
	// eventDownloadTS the ts file is downloaded
	eventDownloadTS
	// eventM3U8Parsed the m3u8 file is parsed
	eventM3U8Parsed
	eventCreateArticleFoldFailed

	eventUIUpdateProgress
	eventUIProgressTotal

	eventCount
)

var (
	eventNames = []string{
		"course",
		"articles",
		"articleFinished",
		"articleVideo",
		"playAuth",
		"playList",
		"downloadTS",
		"m3u8Parsed",
		"createVideoFoldFailed",
		"uiUpdateProgress",
		"uiProgressTotal",
	}
)

func (e event) String() string {
	return eventNames[e]
}

type bus struct {
	eventHandlers [eventCount][]func(interface{})
}

func (b *bus) subscribe(e event, f func(interface{})) {
	b.eventHandlers[e] = append(b.eventHandlers[e], f)
}

func (b *bus) post(e event, v interface{}) {
	for _, f := range b.eventHandlers[e] {
		f(v)
	}
}

type courseRet struct {
	course course
	err    error
}

func (c courseRet) String() string {
	return fmt.Sprintf("course ret[course:%s, error:%s]", c.course.String(), errMsg(c.err))
}

type articles struct {
	articles []article
	err      error
}

func (a articles) String() string {
	as := make([]string, len(a.articles))
	for i := range a.articles {
		as[i] = a.articles[i].String()
	}
	return fmt.Sprintf("articles:[articles:{%s}, error:%s]", strings.Join(as, ","), errMsg(a.err))
}

type articleVideo struct {
	articleID int
	name      string
	id        string
	err       error
}

func (av articleVideo) String() string {
	return fmt.Sprintf(
		"articleVideo:[articleID:%d,name:%s,videoID:%s,error:%s]",
		av.articleID, av.name, av.id, errMsg(av.err),
	)
}

type playAuth struct {
	auth      videoPlayAuth
	articleID int
	videoID   string
	err       error
}

func (p playAuth) String() string {
	as := fmt.Sprintf("%+v", p.auth)
	return fmt.Sprintf(
		"playAuth:[auth:%s,artilceID:%d,videoID:%s,error:%s]",
		as, p.articleID, p.videoID, errMsg(p.err),
	)
}

type playListRet struct {
	list      playList
	articleID int
	err       error
}

func (p playListRet) String() string {
	ls := fmt.Sprintf("%+v", p.list)
	return fmt.Sprintf("playListRet:[list:%s,articleID:%d,error:%s]", ls, p.articleID, errMsg(p.err))
}

type downloadTS struct {
	url       string
	articleID int
	err       error
}

func (d downloadTS) String() string {
	return fmt.Sprintf("downloadTS:[articleID:%d,url:%s,error:%s", d.articleID, d.url, errMsg(d.err))
}

type m3u8 struct {
	name      string
	articleID int
	ts        []string
	m3u8URL   string
	outputDir string
	err       error
}

func (m m3u8) String() string {
	return fmt.Sprintf(
		"m3u8:[name:%s,articleID:%d,ts count:%d,m3u8URL:%s,outputDir:%s,error:%s]",
		m.name, m.articleID, len(m.ts), m.m3u8URL, m.outputDir, errMsg(m.err),
	)
}

func errMsg(err error) string {
	if err == nil {
		return ""
	}
	return err.Error()
}
