package main

import (
	"errors"
	"flag"
	"fmt"

	"github.com/geektime"
)

var (
	courseID  string
	cookie    string
	outputDir string
)

func init() {
	flag.StringVar(&courseID, "courseID", "", "course id")
	flag.StringVar(&cookie, "cookie", "", "cookie after login")
	flag.StringVar(&outputDir, "output", ".", "output dir")
}

func main() {
	if err := parseParams(); err != nil {
		fmt.Println(err.Error())
		return
	}

	downloader, err := geektime.NewDownloader(geektime.Config{
		CourseID:  courseID,
		Cookie:    cookie,
		OutputDir: outputDir,
	})

	if err != nil {
		fmt.Printf("create pipelie error:%s\n", err)
		return
	}

	downloader.Run()
}

func parseParams() error {
	flag.Parse()
	if courseID == "" {
		return errors.New("empty course id")
	}

	if cookie == "" {
		return errors.New("empty cookie")
	}

	return nil
}
