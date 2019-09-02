package geektime

import (
	"path/filepath"
	"testing"
)

func TestToMP4(t *testing.T) {
	m3u8Paths := listM3U8Paths("testdata")

	outputPaths := make([]string, len(m3u8Paths))
	for i, p := range m3u8Paths {
		outputPaths[i] = filepath.Dir(p)
	}

	for i, mp := range m3u8Paths {
		op := outputPaths[i] + ".mp4"
		t.Logf("convert %s to MP4, output:%s", mp, op)
		ret, err := m3u8ToMP4(mp, op)
		t.Logf("convert ret:%s, error:%v", ret, err)
	}
}
