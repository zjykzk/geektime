package geektime

import (
	"os"
	"testing"
)

func TestOnlineAPI(t *testing.T) {
	cookie := "GCID=0326b6b-9a5dacc-ed688ec-162b3dc; expires=Fri, 04-Oct-2019 11:09:29 GMT; Max-Age=2592000; path=/; domain=.geekbang.org; HttpOnly GCESS=BAsCBAAIAQMMAQEGBEAunYkJAQEHBMdmi90CBGmbb10FBAAAAAABBKV3DwADBGmbb10EBAAvDQAKBAAAAAA-; expires=Sat, 14-Sep-2019 11:09:29 GMT; Max-Age=864000; path=/; domain=.geekbang.org; HttpOnly SERVERID=1fa1f330efedec1559b3abbcb6e30f50|1567595369|1567595313;Path=/"

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
