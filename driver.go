package geektime

import (
	"errors"
	"path/filepath"
	"strings"
)

const (
	defaultWorkerCount = 512
)

// Config the configuration for creating pipeline
type Config struct {
	ArticlesURL      string
	ArticleURL       string
	IntroURL         string
	VideoPlayAuthURL string
	PlayListURL      string

	CourseID  string
	OutputDir string
	Cookie    string

	Overwrite bool
}

type driver struct {
	Config

	bus *bus

	apiExecutor      *executor
	downloadExecutor *executor
	playAuthExecutor *executor

	course course
}

func newDriver(conf Config, bus *bus) (*driver, error) {
	wrapErr := func(m string) error {
		return errors.New("create driver error:" + m)
	}

	if conf.CourseID == "" {
		return nil, wrapErr("empty course id")
	}

	if conf.Cookie == "" {
		return nil, wrapErr("empty cookie")
	}

	if conf.ArticlesURL == "" {
		conf.ArticlesURL = defaultArticlesURL
	}

	if conf.ArticleURL == "" {
		conf.ArticleURL = defaultArticleURL
	}

	if conf.IntroURL == "" {
		conf.IntroURL = defaultIntroURL
	}

	if conf.VideoPlayAuthURL == "" {
		conf.VideoPlayAuthURL = defaultVideoPlayAuthURL
	}

	if conf.PlayListURL == "" {
		conf.PlayListURL = defaultPlayListURL
	}

	return &driver{
		Config:           conf,
		bus:              bus,
		apiExecutor:      newExecutor(128, 256),
		playAuthExecutor: newExecutor(1, 256), // for the rate limit
		downloadExecutor: newExecutor(512, 4096),
	}, nil
}

// Start starts download
func (p *driver) Start() {
	p.subscribeEvents()

	p.apiExecutor.start()
	p.playAuthExecutor.start()
	p.downloadExecutor.start()

	p.apiExecutor.execute(&courseFetcher{p.bus, p.IntroURL, p.CourseID, p.Cookie, p.OutputDir})
	p.apiExecutor.execute(&articlesFetcher{p.bus, p.ArticlesURL, p.CourseID, p.Cookie})
}

// Abort shutdowns the executor
func (p *driver) Abort() {
	p.apiExecutor.shutdown()
	p.playAuthExecutor.shutdown()
	p.downloadExecutor.shutdown()
}

func (p *driver) subscribeEvents() {
	p.bus.subscribe(eventCourse, func(v interface{}) { p.onCourse(v.(courseRet)) })
	p.bus.subscribe(eventArticles, func(v interface{}) { p.onActicles(v.(articles)) })
	p.bus.subscribe(eventArticleVideo, func(v interface{}) { p.onActicleVideo(v.(articleVideo)) })
	p.bus.subscribe(eventPlayAuth, func(v interface{}) { p.onPlayAuth(v.(playAuth)) })
	p.bus.subscribe(eventPlayList, func(v interface{}) { p.onPlayList(v.(playListRet)) })
	p.bus.subscribe(eventM3U8Parsed, func(v interface{}) { p.onM3U8(v.(m3u8)) })
}

func (p *driver) onCourse(c courseRet) {
	if c.err != nil {
		return
	}
	p.course = c.course
	makeSureDirExist(filepath.Join(p.OutputDir, c.course.Title))
}

func (p *driver) onActicles(as articles) {
	if as.err != nil {
		return
	}

	for _, a := range as.articles {
		if isAudio(a) {
			p.procAudit(a)
		} else if a.VideoID != "" {
			p.apiExecutor.execute(&articleVideoFetcher{
				bus:        p.bus,
				articleURL: p.ArticleURL,
				courseID:   p.CourseID,
				cookie:     p.Cookie,
				articleID:  a.ID,
				videoID:    a.VideoID,
			})
		} else {
			p.bus.post(eventArticleFinished, a.ID)
		}
	}
}

func (p *driver) procAudit(a article) {
	outputDir := filepath.Join(p.OutputDir, p.course.Title, simplify(a.Title))
	err := makeSureDirExist(outputDir)
	if err != nil {
		p.bus.post(eventCreateArticleFoldFailed, a.ID)
		return
	}

	p.apiExecutor.execute(&m3u8Fetcher{
		p.bus,
		a.Title,
		a.ID,
		[]string{a.AuditM3U8URL},
		outputDir,
	})
}

func isAudio(a article) bool {
	return a.AuditM3U8URL != ""
}

func (p *driver) onActicleVideo(av articleVideo) {
	if av.err != nil {
		return
	}

	name := simplify(av.name)
	err := makeSureDirExist(filepath.Join(p.OutputDir, p.course.Title, name))
	if err != nil {
		p.bus.post(eventCreateArticleFoldFailed, av.articleID)
		return
	}

	p.apiExecutor.execute(&playAuthFetcher{p.bus, name, p.VideoPlayAuthURL, av.id, p.Cookie, av.articleID})
}

func (p *driver) onPlayAuth(pa playAuth) {
	if pa.err == nil {
		p.playAuthExecutor.execute(&playListFetcher{
			p.bus,
			pa.name,
			pa.articleID,
			p.PlayListURL,
			pa.videoID,
			pa.auth,
		})
	}
}

func (p *driver) onPlayList(l playListRet) {
	if l.err != nil {
		return
	}

	urls := make([]string, len(l.list.PlayInfoList.PlayInfo))
	for i, p := range l.list.PlayInfoList.PlayInfo {
		urls[i] = p.PlayURL
	}
	vp := filepath.Join(p.OutputDir, p.course.Title, l.name)

	p.apiExecutor.execute(&m3u8Fetcher{
		p.bus,
		l.name,
		l.articleID,
		urls,
		vp,
	})
}

func (p *driver) onM3U8(m m3u8) {
	if m.err != nil {
		return
	}

	path := m.m3u8URL[:strings.LastIndexByte(m.m3u8URL, '/')+1]
	for _, ts := range m.ts {
		p.downloadExecutor.execute(&tsDownloader{
			p.bus,
			m.articleID,
			path + ts,
			m.outputDir,
		})
	}
}

func simplify(t string) string {
	invalidChars := `| \:：*?？“”"，,()（）`

	var nr []rune
	for _, r := range t {

		if strings.IndexRune(invalidChars, r) >= 0 {
			continue
		}
		nr = append(nr, r)
	}
	return string(nr)
}
