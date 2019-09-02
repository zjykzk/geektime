package geektime

import "path/filepath"

// Downloader download the files
//
// monitor the download progress and display it with the console ui
// it also log the details
type Downloader struct {
	bus      *bus
	driver   *driver
	progress *courseProgress
	logger   *logger
	ui       *cui
}

// NewDownloader creates the downloader
func NewDownloader(conf Config) (*Downloader, error) {
	bus := &bus{}

	driver, err := newDriver(conf, bus)
	if err != nil {
		return nil, err
	}
	logger, err := newLogger(filepath.Join(conf.OutputDir, "download.log"))
	if err != nil {
		return nil, err
	}

	return &Downloader{
		bus:      bus,
		driver:   driver,
		progress: newProgress(bus),
		logger:   logger,
		ui:       newCUI(),
	}, nil
}

// Run run the downloader
func (d *Downloader) Run() {
	d.progress.subscribeEvents()
	d.ui.subscribeEvents(d.bus)
	d.logger.subscribeEvents(d.bus)

	d.driver.Start()

	d.ui.run()
}
