package main

import (
	"errors"
	"flag"
	"fmt"

	"github.com/geektime"
)

var (
	courseID  string
	cellPhone string
	outputDir string
)

func init() {
	flag.StringVar(&courseID, "courseID", "", "course id")
	flag.StringVar(&cellPhone, "phone", "", "the cell phone")
	flag.StringVar(&outputDir, "output", ".", "output dir")
}

func main() {
	if err := parseParams(); err != nil {
		fmt.Println(err.Error())
		return
	}

	downloader, err := geektime.NewDownloader(geektime.Config{
		CourseID:  courseID,
		OutputDir: outputDir,
	}, cellPhone)

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
	if courseID == "" {
		return errors.New("empty course id")
	}

	if cellPhone == "" {
		return errors.New("empty cell phone")
	}

	return nil
}
