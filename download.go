package geektime

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

// Downloader download the files
//
// monitor the download progress and display it with the console ui
// it also log the details
type Downloader struct {
	Config

	bus      *bus
	progress *courseProgress
	logger   *logger
	ui       *cui

	cellPhone string
	name      string
}

// NewDownloader creates the downloader
func NewDownloader(conf Config, cellPhone, name string) (*Downloader, error) {
	bus := &bus{}

	logger, err := newLogger(filepath.Join(conf.OutputDir, "download.log"))
	if err != nil {
		return nil, err
	}

	return &Downloader{
		Config:    conf,
		bus:       bus,
		progress:  newProgress(bus),
		cellPhone: cellPhone,
		name:      name,
		logger:    logger,
		ui:        newCUI(),
	}, nil
}

// Run run the downloader
func (d *Downloader) Run() error {
	err := d.readCookie()
	if err != nil {
		return err
	}

	err = d.findCourseID()
	if err != nil {
		return err
	}

	driver, err := newDriver(d.Config, d.bus)
	if err != nil {
		return err
	}

	d.progress.subscribeEvents()
	d.ui.subscribeEvents(d.bus)
	d.logger.subscribeEvents(d.bus)

	driver.Start()

	d.ui.run()
	return nil
}

func (d *Downloader) readCookie() error {
	err := d.readCookieFromFile()
	if err == nil {
		return nil
	}

	return d.login()
}

func (d *Downloader) readCookieFromFile() error {
	data, err := ioutil.ReadFile(d.cookieFilePath())
	if err != nil {
		return err
	}
	if len(data) == 0 {
		return errors.New("empty cookie from the file")
	}
	d.Cookie = strings.TrimSpace(string(data))
	return nil
}

func (d *Downloader) findCourseID() error {
	cs, err := allCoursesBought(defaultAllBoughtURL, d.Cookie)
	if err != nil {
		return nil
	}

	for _, c := range cs {
		for _, i := range c {
			if i.Title == d.name {
				d.CourseID = strconv.Itoa(i.ID)
				return nil
			}
		}
	}
	return fmt.Errorf(`no course of name:"%s"`, d.name)
}

func (d *Downloader) login() error {
	cookie, err := sendSMSCode(d.cellPhone)
	if err != nil {
		return fmt.Errorf("send the sms code error:%s", err)
	}

	smsCode := make([]byte, 16)
	fmt.Printf("input the sms code:")
	n, err := os.Stdin.Read(smsCode)
	if err != nil {
		return fmt.Errorf("read the sms code error:%s", err)
	}

	d.Cookie, err = login(cookie, d.cellPhone, string(smsCode[:n-1])) // trim \n
	writeFile(d.cookieFilePath(), []byte(d.Cookie))

	if err != nil {
		return fmt.Errorf("login error:%s", err)
	}

	return nil
}

func (d *Downloader) cookieFilePath() string {
	return filepath.Join(d.OutputDir, "cookie")
}
