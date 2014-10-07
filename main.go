package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
)

type OptionResult struct {
}

func (c *Config) optionsHandler(w http.ResponseWriter, h *http.Request) {
	var result OptionResult
	t, err := template.ParseFiles("templates/options.html")
	if err != nil {
		log.Println("Error: ", err)
		return
	}
	t.Execute(w, result)
}

func (c *Config) listHandler(w http.ResponseWriter, h *http.Request) {
	t, err := template.ParseFiles("templates/list.html")
	if err != nil {
		log.Println("Error: ", err)
		return
	}
	c.list.Ready = c.isReady()
	t.Execute(w, c.list)
}

func (c *Config) refreshHandler(w http.ResponseWriter, h *http.Request) {
	err := c.refreshList()
	if err != nil {
		log.Println(err)
	}

	http.Redirect(w, h, "/", http.StatusTemporaryRedirect)
}

type ScanResult struct {
	Success bool
}

func (c *Config) scanHandler(w http.ResponseWriter, h *http.Request) {
	var result ScanResult

	err := c.requestScan()
	result.Success = (err != nil)

	t, err := template.ParseFiles("templates/scan.html")
	if err != nil {
		log.Println("Error: ", err)
		return
	}

	t.Execute(w, result)
}

func (c *Config) scanWaitHandler(w http.ResponseWriter, h *http.Request) {
	// Wait for scan end
	c.waitUntilReady()
	fmt.Fprint(w, "ok")
}

func handleStyle(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "static/style.css")
}

func main() {

	config := DefaultConfig()
	config.refreshList()
	go config.handleScanRequests()

	http.HandleFunc("/", config.listHandler)
	http.HandleFunc("/options", config.optionsHandler)
	http.HandleFunc("/scan", config.scanHandler)
	http.HandleFunc("/scan/wait", config.scanWaitHandler)
	http.HandleFunc("/refresh", config.refreshHandler)

	http.HandleFunc("/static/style.css", handleStyle)
	http.Handle(config.targetPath, http.StripPrefix(config.targetPath, http.FileServer(http.Dir(config.dataDir))))

	log.Println("Listening on port 8080")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Println("Error serving http:", err)
	}
}
