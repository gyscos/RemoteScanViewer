package main

import (
	"flag"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
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
	Success  bool
	Expected int
	Elapsed  int
}

func (c *Config) scanHandler(w http.ResponseWriter, h *http.Request) {
	var result ScanResult

	err := c.requestScan()
	result.Success = (err != nil)

	result.Expected = int(1000 * c.scanDuration.Seconds())
	result.Elapsed = int(1000 * time.Since(c.scanStart).Seconds())

	t, err := template.ParseFiles("templates/scan.html")
	if err != nil {
		log.Println("Error: ", err)
		return
	}

	t.Execute(w, result)
}

func (c *Config) removeHandler(w http.ResponseWriter, h *http.Request) {
	// Find page id
	target := strings.TrimPrefix(h.FormValue("target"), c.targetPath)

	pdfFile := c.dataDir + target
	thbFile := c.dataDir + ".thumb/" + strings.TrimSuffix(target, ".pdf") + ".png"
	os.Remove(pdfFile)
	os.Remove(thbFile)

	c.refreshList()

	fmt.Fprintf(w, "Ok")
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

	dummy := flag.Bool("dummy", false, "Wait instead of really scanning pages")
	port := flag.Int("port", 8080, "Port for HTTP interface")

	flag.Parse()

	config := DefaultConfig()
	config.dummyScan = *dummy
	config.refreshList()
	go config.handleScanRequests()

	http.HandleFunc("/", config.listHandler)
	http.HandleFunc("/options", config.optionsHandler)
	http.HandleFunc("/scan", config.scanHandler)
	http.HandleFunc("/scan/wait", config.scanWaitHandler)
	http.HandleFunc("/refresh", config.refreshHandler)
	http.HandleFunc("/remove", config.removeHandler)

	http.HandleFunc("/static/style.css", handleStyle)
	http.Handle(config.targetPath, http.StripPrefix(config.targetPath, http.FileServer(http.Dir(config.dataDir))))

	log.Println("Listening on port " + strconv.Itoa(*port))
	err := http.ListenAndServe(":"+strconv.Itoa(*port), nil)
	if err != nil {
		log.Println("Error serving http:", err)
	}
}
