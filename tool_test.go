package geektimedl

import "testing"

func TestCheckM3U8Download(t *testing.T) {
	var notCompletes, completes []string
	for _, p := range listM3U8Paths("testdata") {
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
