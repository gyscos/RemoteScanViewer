package main

import (
	"errors"
	"log"
	"time"
)

type Config struct {
	list ItemList

	// dataDir is the directory where files are actually located
	dataDir string

	// targetPath is the virtual url prefix for pdf and thumbnails
	targetPath string

	scanRequests chan struct{}
	scannerReady chan struct{}

	scanStart    time.Time
	scanDuration time.Duration

	dummyScan bool
}

func DefaultConfig() Config {
	config := Config{
		dataDir:      "/data/scans/",
		targetPath:   "/scans/",
		scanDuration: 10 * time.Second,
		scanRequests: make(chan struct{}),
		scannerReady: make(chan struct{}, 1)}

	config.scannerReady <- struct{}{}
	return config
}

func (c *Config) refreshList() error {
	return c.list.refresh(c.dataDir, c.targetPath)
}

func (c *Config) scan() {
	if c.dummyScan {
		time.Sleep(10 * time.Second)
	} else {
		scanImage(c.dataDir)
	}
	c.refreshList()
}

func (c *Config) isReady() bool {

	select {
	case <-c.scannerReady:
		c.scannerReady <- struct{}{}
		return true
	default:
		return false
	}
}

func (c *Config) waitUntilReady() {
	<-c.scannerReady
	c.scannerReady <- struct{}{}
}

func (c *Config) requestScan() error {
	select {
	case <-c.scannerReady:
		c.scanStart = time.Now()
		c.scanRequests <- struct{}{}
		return nil
	default:
		return errors.New("Scanner not ready!")
	}
}

func (c *Config) handleScanRequests() {
	log.Println("Handling requests")
	for {
		<-c.scanRequests
		c.scan()
		c.scanDuration = time.Since(c.scanStart)
		c.scannerReady <- struct{}{}
	}
	log.Println("Leaving handleScanRequests")
}
