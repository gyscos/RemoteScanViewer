package main

import (
	"io/ioutil"
	"sort"
	"strings"
	"time"
)

type ItemList struct {
	Items []Item
}

type Item struct {
	Thumbnail string
	Filepath  string
	Title     string
	Date      string
	Id        int
	ModTime   time.Time
}

func (l *ItemList) Len() int {
	return len(l.Items)
}
func (l *ItemList) Swap(i, j int) {
	l.Items[i], l.Items[j] = l.Items[j], l.Items[i]
}
func (l *ItemList) Less(i, j int) bool {
	return l.Items[j].ModTime.Before(l.Items[i].ModTime)
}

func (l *ItemList) refresh(dataDir string, targetPath string) error {
	infos, err := ioutil.ReadDir(dataDir)
	if err != nil {
		return err
	}

	for _, info := range infos {
		if info.IsDir() {
			continue
		}
		if !strings.HasSuffix(info.Name(), ".pdf") {
			continue
		}

		var item Item
		item.Filepath = targetPath + info.Name()
		item.ModTime = info.ModTime()
		item.Date = info.ModTime().Format("2006-01-02")
		item.Thumbnail = targetPath + ".thumb/" + strings.TrimSuffix(info.Name(), ".pdf") + ".png"

		l.Items = append(l.Items, item)
	}

	sort.Sort(l)
	for i, _ := range l.Items {
		l.Items[i].Id = len(l.Items) - i
	}

	return nil
}
