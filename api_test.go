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
	t.Log(allCoursesBought(defaultAllBoughtURL, "_ga=GA1.2.579002931.1565926025; MEIQIA_TRACK_ID=1Q32IKd5hEYC9yGzdJG9A3faXxR; LF_ID=1575772763277-9041709-359974; MEIQIA_VISIT_ID=1VhN6rYFfd6ZlVVCyBRmAz95zWs; gksskpitn=563745b7-2304-419a-b961-ccf3e3c3ba67; _gid=GA1.2.2005389193.1590845167; Hm_lvt_59c4ff31a9ee6263811b23eb921a5083=1590845167; GCID=48aa174-77487b0-d41444f-ac3d4f2; GRID=48aa174-77487b0-d41444f-ac3d4f2; GCESS=BQkBAQcE46.YhAYELFjIJQsCBQAFBAAAAAAEBAAvDQADBGe8014KBAAAAAAMAQEIAQMCBGe8014BCKV3DwAAAAAA; Hm_lvt_022f847c4e3acd44d4a2481d9187f1e6=1590845167,1590845195,1590845708,1590934633; Hm_lpvt_59c4ff31a9ee6263811b23eb921a5083=1590934637; gk_process_ev={%22count%22:4%2C%22utime%22:1590934613121%2C%22referrer%22:%22https://time.geekbang.org/%22%2C%22target%22:%22%22}; Hm_lpvt_022f847c4e3acd44d4a2481d9187f1e6=1590934637; SERVERID=3431a294a18c59fc8f5805662e2bd51e|1590934923|1590934606"))
}

func TestHasDownload(t *testing.T) {

	fn, existed, err := hasDownload(
		"https://media001.geekbang.org/36c4b682d1cf4e469eef440d5e69546f/b978afa8938a4386b901270c329ee5f2-9306829bcc29ce295e81d30664d4a7ae-ld-00157.ts",
		"/Users/zenk/tools/geektime/src/github.com/geektime/test-course/移动端自动化测试实战/58测试报告定制",
	)
	t.Log(fn, existed, err)
}
