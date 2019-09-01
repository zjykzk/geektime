package main

import (
	"errors"
	"flag"
	"fmt"

	geektimedl "github.com/geektime-dl"
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

	downloader, err := geektimedl.NewDownloader(geektimedl.Config{
		CourseID: courseID,
		Cookie:   cookie,
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
