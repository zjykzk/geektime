package geektime

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func newFakeServer(t *testing.T) *httptest.Server {
	var host string
	server := httptest.NewServer(http.HandlerFunc(func(resp http.ResponseWriter, req *http.Request) {
		var content string
		switch path := req.URL.Path; {
		case strings.HasSuffix(path, "articles"):
			content = `{
	"error":[],
	"extra":[],
	"code":0,
	"data":{
		"list":[{
			"id":1,
			"article_title":"test article"
		}],
		"page":{
				"count":71,"more":false
			}
		}
}`
		case strings.HasSuffix(path, "article"):
			content = fmt.Sprintf(`{
	"error":[],
	"extra":[],
	"code":0,
	"data":{
	"video_preview": {
		"sd": {
			"url": "%s/sd.m3u8",
			"size": 18467992
		}
	},
	"article_title":"test article"
}}`, host)
		case strings.HasSuffix(path, "m3u8"):
			content = `#EXTM3U
#EXT-X-VERSION:3
#EXT-X-TARGETDURATION:10
#EXT-X-MEDIA-SEQUENCE:0
#EXTINF:10.000000,
bf43e8565761e1a459ef78cc5a1b3c64-1.ts
#EXTINF:10.000000,
fd09bcb049dc5d166ac20b031cb92a9f-2.ts`
		case strings.HasSuffix(path, "ts"):
			content = "ts ok"
		case strings.HasSuffix(path, "intro"):
			content = `{
	"error":[],
	"extra":[],
	"code":0,
	"data":{
	"column_title": "test column"
}}`
		case strings.HasSuffix(path, "playAuth"):
			content = "{}"
		}

		resp.Write([]byte(content))
	}))
	host = server.URL
	return server
}

func TestSimplify(t *testing.T) {
	assert.Equal(t, "02Python的发展历史与版本", simplify("02 | Python的发展历史与版本"))
	assert.Equal(t, "30答疑文章二用动态的观点看加锁", simplify("30 | 答疑文章（二）：用动态的观点看加锁"))
	invalidChars := `\:：*?？“”"，,`
	assert.Equal(t, "", simplify(invalidChars))
}

func testDownload(t *testing.T) {
	p, err := newDriver(Config{
		CourseID:  "98",
		OutputDir: "testdata",
		Cookie:    "_ga=GA1.2.579002931.1565926025; GCID=5571b9d-d62bcc2-d94b2b8-627c523; modalclose=true; modalCloseDay=16; _gid=GA1.2.1309432447.1566809933; GCESS=BAUEAAAAAAcE4DLmeggBAwMEdJ9jXQkBAQYEVQ967AoEAAAAAAIEdJ9jXQsCBAAEBAAvDQABBKV3DwAMAQE-; Hm_lvt_022f847c4e3acd44d4a2481d9187f1e6=1566812553,1566870489,1566870540,1566870552; Hm_lpvt_022f847c4e3acd44d4a2481d9187f1e6=1566870552; SERVERID=3431a294a18c59fc8f5805662e2bd51e|1566870623|1566870488; _gat=1",
	}, &bus{})

	if err != nil {
		t.Fatal(err)
	}

	p.Start()
}
