package geektime

import (
	"path/filepath"
)

type task interface {
	run()
}

type courseFetcher struct {
	bus                                   *bus
	introURL, courseID, cookie, outputDir string
}

func (f *courseFetcher) run() {
	c, err := fetchCourse(f.introURL, f.courseID, f.cookie)
	f.bus.post(eventCourse, courseRet{c, err})
}

type articlesFetcher struct {
	bus                           *bus
	articlesURL, courseID, cookie string
}

func (f *articlesFetcher) run() {
	as, err := fetchArticles(f.articlesURL, f.courseID, f.cookie)
	f.bus.post(eventArticles, articles{as, err})
}

type articleVideoFetcher struct {
	bus                          *bus
	articleURL, courseID, cookie string
	articleID                    int
	videoID                      string
}

func (f *articleVideoFetcher) run() {
	v, err := fetchVideoOfArticle(f.articleURL, f.courseID, f.cookie, f.articleID)
	f.bus.post(eventArticleVideo, articleVideo{f.articleID, v.Name, f.videoID, err})
}

type playAuthFetcher struct {
	bus                               *bus
	videoPlayAuthURL, videoID, cookie string
	articleID                         int
}

func (f *playAuthFetcher) run() {
	auth, err := fetchVideoPlayAuth(f.videoPlayAuthURL, f.cookie, f.articleID, 1, f.videoID)
	f.bus.post(eventPlayAuth, playAuth{auth, f.articleID, f.videoID, err})
}

type playListFetcher struct {
	bus                  *bus
	articleID            int
	playListURL, videoID string
	auth                 videoPlayAuth
}

func (f *playListFetcher) run() {
	l, err := fetchPlayList(f.playListURL, f.videoID, f.auth)
	f.bus.post(eventPlayList, playListRet{l, f.articleID, err})
}

type m3u8Fetcher struct {
	bus       *bus
	name      string
	articleID int
	urls      []string
	outputDir string
}

func (f *m3u8Fetcher) run() {
	var (
		url, path string
		err       error
	)

	for _, url = range f.urls {
		path, err = download(url, f.outputDir)
		if err == nil {
			break
		}
	}

	if err != nil {
		f.bus.post(eventM3U8Parsed, m3u8{articleID: f.articleID, err: err})
		return
	}

	writeFile(filepath.Join(f.outputDir, "m3u8_addr"), []byte(url))

	ts, err := parseM3u8(path)
	if err != nil {
		f.bus.post(eventM3U8Parsed, m3u8{articleID: f.articleID, err: err})
		return
	}

	f.bus.post(eventM3U8Parsed, m3u8{
		name:      f.name,
		articleID: f.articleID,
		ts:        ts,
		m3u8URL:   url,
		outputDir: f.outputDir,
	})
}

type tsDownloader struct {
	bus            *bus
	articleID      int
	url, outputDir string
}

func (d *tsDownloader) run() {
	_, err := download(d.url, d.outputDir)
	d.bus.post(eventDownloadTS, downloadTS{d.url, d.articleID, err})
}
