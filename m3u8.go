package geektime

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

func parseM3u8(m3u8Path string) ([]string, error) {
	f, err := os.Open(m3u8Path)
	if err != nil {
		return nil, err
	}

	defer f.Close()
	data, err := ioutil.ReadAll(f)
	if err != nil {
		return nil, err
	}

	var tsFilenames []string
	for _, l := range strings.Split(string(data), "\n") {
		l = strings.TrimSpace(l)
		if l == "" || strings.HasPrefix(l, "#") {
			continue
		}
		tsFilenames = append(tsFilenames, l)
	}

	return tsFilenames, nil
}

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

func doesDownloadAllTS(m3u8Path string) (bool, error) {
	ts, err := parseM3u8(m3u8Path)
	if err != nil {
		return false, err
	}

	dir := filepath.Dir(m3u8Path)
	for _, t := range ts {
		if ok, err := fileExists(filepath.Join(dir, t)); err == nil {
			if !ok {
				return false, nil
			}
		} else {
			return false, err
		}
	}

	return true, nil
}

// M3U8Converter converts the videos with the m3u8 format to the mp4
type M3U8Converter struct {
	outputDir, inputDir string
}

// NewM3U8Converter creates the converter to convert the m3u8-video to the mp4
func NewM3U8Converter(inputDir, outputDir string) (*M3U8Converter, error) {
	if inputDir == "" {
		return nil, errors.New("empty input dir")
	}

	if outputDir == "" {
		return nil, errors.New("empty output dir")
	}

	err := makeSureDirExist(outputDir)
	if err != nil {
		return nil, err
	}

	return &M3U8Converter{outputDir, inputDir}, nil
}

// ToMP4 converts to the mp4 format
func (c *M3U8Converter) ToMP4() {
	c.convert(m3u8ToMP4, func(m3u8Path string) string {
		return filepath.Join(c.outputDir, filepath.Base(filepath.Dir(m3u8Path))) + ".mp4"
	})
}

// ToMP3 converts to the mp3 format
func (c *M3U8Converter) ToMP3() {
	c.convert(m3u8ToMP3, func(m3u8Path string) string {
		return filepath.Join(c.outputDir, filepath.Base(filepath.Dir(m3u8Path))) + ".mp3"
	})
}

func (c *M3U8Converter) convert(
	convert func(in, out string) (string, error),
	filename func(m3u8Path string) string,
) {
	for _, m := range listM3U8Paths(c.inputDir) {
		outputPath := filename(m)
		fmt.Printf("convert %s, to %s\n", m, outputPath)
		ret, err := convert(m, outputPath)

		msg := "success!"
		if err != nil {
			msg = err.Error()
		}
		fmt.Printf("result:%s, message:%s\n", ret, msg)
	}
}
