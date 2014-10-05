package main

import "log"

type Config struct {
	list ItemList

	// dataDir is the directory where files are actually located
	dataDir string

	// targetPath is the virtual url prefix for pdf and thumbnails
	targetPath string

	scanRequests chan struct{}
}

func DefaultConfig() Config {
	return Config{
		dataDir:      "/data/scans/",
		targetPath:   "/scans/",
		scanRequests: make(chan struct{}, 5)}
}

func (c *Config) refreshList() error {
	return c.list.refresh(c.dataDir, c.targetPath)
}

func (c *Config) scan() {
	log.Println("Scanning...")
	scanImage(c.dataDir)
	c.refreshList()
}

func (c *Config) requestScan() {
	c.scanRequests <- struct{}{}
}

func (c *Config) handleScanRequests() {
	log.Println("Handling requests")
	for {
		<-c.scanRequests
		c.scan()
	}
	log.Println("Leaving handleScanRequests")
}
