package geektime

import "testing"

func TestCheckM3U8Download(t *testing.T) {
	var notCompletes, completes []string
	for _, p := range listM3U8Paths("/Users/zenk/tools/geektime/src/github.com/geektime/test-course/移动端自动化测试实战") {
		ok, err := doesDownloadAllTS(p)
		t.Logf("check %s, exist:%t, error:%v", p, ok, err)
		if !ok {
			notCompletes = append(notCompletes, p)
		} else {
			completes = append(completes, p)
		}
	}
	t.Log(notCompletes)
	t.Logf("complete count:%d", len(completes))
}
