package geektime

import (
	"os"
	"testing"
)

func TestOnlineAPI(t *testing.T) {
	cookie := "_ga=GA1.2.579002931.1565926025; GCID=5571b9d-d62bcc2-d94b2b8-627c523; modalclose=true; modalCloseDay=16; _gid=GA1.2.1309432447.1566809933; GCESS=BAUEAAAAAAcE4DLmeggBAwMEdJ9jXQkBAQYEVQ967AoEAAAAAAIEdJ9jXQsCBAAEBAAvDQABBKV3DwAMAQE-; Hm_lvt_022f847c4e3acd44d4a2481d9187f1e6=1566812553,1566870489,1566870540,1566870552; Hm_lpvt_022f847c4e3acd44d4a2481d9187f1e6=1566870552; SERVERID=3431a294a18c59fc8f5805662e2bd51e|1566870623|1566870488; _gat=1"

	t.Log(fetchCourse(defaultIntroURL, "98", cookie))
	t.Log(fetchArticles(defaultArticlesURL, "98", cookie))

	t.Log(fetchVideoOfArticle(defaultArticleURL, "98", cookie, 8336))
	auth, err := fetchVideoPlayAuth(
		defaultVideoPlayAuthURL, cookie, 8337, 1, "f4ee5bdae6c84b2c9c9f7f4f130c2b0e",
	)
	t.Logf("%+v", auth)
	if err != nil {
		t.Fatal(err)
	}

	t.Log(fetchPlayList(defaultPlayListURL, "f4ee5bdae6c84b2c9c9f7f4f130c2b0e", auth))
}

func TestFetchM3U8(t *testing.T) {
	m3u8Filename, err := download("https://static001.geekbang.org/files/vod/16e79badf28f43e8b07fb0f1adb54b1e/sd.m3u8", "testdata/fetchM3U8")
	t.Log(m3u8Filename, err)
	t.Log(parseM3u8(m3u8Filename))

	os.RemoveAll("testdata/fetchM3U8")
}
