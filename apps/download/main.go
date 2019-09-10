package main

import (
	"errors"
	"flag"
	"fmt"

	"github.com/zjykzk/geektime"
)

var (
	course    string
	cellPhone string
	outputDir string
)

func init() {
	flag.StringVar(&course, "course", "", "course name")
	flag.StringVar(&cellPhone, "phone", "", "the cell phone")
	flag.StringVar(&outputDir, "output", ".", "output dir")
}

func main() {
	if err := parseParams(); err != nil {
		fmt.Println(err.Error())
		return
	}

	downloader, err := geektime.NewDownloader(
		geektime.Config{OutputDir: outputDir}, cellPhone, course,
	)

	if err != nil {
		fmt.Printf("create pipelie error:%s\n", err)
		return
	}

	err = downloader.Run()
	if err != nil {
		fmt.Printf("download error:%s\n", err)
	}
}

func parseParams() error {
	flag.Parse()
	if course == "" {
		return errors.New("empty course name")
	}

	if cellPhone == "" {
		return errors.New("empty cell phone")
	}

	return nil
}
