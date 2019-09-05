package geektime

import (
	"os"
	"testing"
)

func TestOnlineAPI(t *testing.T) {
	cookie := "SERVERID=3431a294a18c59fc8f5805662e2bd51e|1567660915|1567660915;;GCID=69e4831-a62deff-8ca13c7-4e06a82; expires=Sat, 05-Oct-2019 05:22:00 GMT; Max-Age=2591999; path=/; domain=.geekbang.org;GCESS=BAcEhV6o_QIEvZtwXQQEAC8NAAgBAwEEpXcPAAkBAQsCBAAGBNpU._gDBL2bcF0MAQEKBAAAAAAFBAAAAAA-; expires=Sun, 15-Sep-2019 05:23:09 GMT; Max-Age=864000; path=/; domain=.geekbang.org;"

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

func TestLogin(t *testing.T) {
	//cookie, err := sendSMSCode("13524581737")
	//t.Log(cookie, err)
	t.Log(login("SERVERID=3431a294a18c59fc8f5805662e2bd51e|1567660915|1567660915;;GCID=69e4831-a62deff-8ca13c7-4e06a82; expires=Sat, 05-Oct-2019 05:22:00 GMT; Max-Age=2591999; path=/; domain=.geekbang.org;", "13524581737", "587715"))
}

func TestAllCoursesBought(t *testing.T) {
	t.Log(allCoursesBought(defaultAllBoughtURL, "GCID=94f2a1a-fb688cd-2f10228-c69936a; _ga=GA1.2.223491743.1565588313; _gid=GA1.2.1720172199.1567569814; GCESS=BAYEMN2sUwkBAQsCBAABBKV3DwAMAQEHBKFFcNQKBAAAAAAFBAAAAAADBKStcF0IAQMCBKStcF0EBAAvDQA-; SERVERID=1fa1f330efedec1559b3abbcb6e30f50|1567671630|1567671630"))
}
