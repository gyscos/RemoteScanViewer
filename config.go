package main

type Config struct {
	list ItemList

	// dataDir is the directory where files are actually located
	dataDir string

	// targetPath is the virtual url prefix for pdf and thumbnails
	targetPath string
}

func DefaultConfig() Config {
	return Config{
		dataDir:    "/data/scans/",
		targetPath: "/scans/"}
}

func (c *Config) refreshList() error {
	return c.list.refresh(c.dataDir, c.targetPath)
}
