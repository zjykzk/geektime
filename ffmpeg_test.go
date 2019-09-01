package geektimedl

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func listM3U8Paths(dir string) []string {
	var m3u8Paths []string
	filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() {
			return nil
		}

		if strings.HasSuffix(path, ".m3u8") {
			m3u8Paths = append(m3u8Paths, path)
		}
		return nil
	})
	return m3u8Paths
}

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
