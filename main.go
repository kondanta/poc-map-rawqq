package main

import (
	"fmt"

	"github.com/gocolly/colly"
)

type ChapterInfo struct {
	URL          string
	LastUpdate   string
	ChapterTitle string
}

func ChapterInfoExtractor() *[]ChapterInfo {
	url := "https://hanascan.com/manga-chiyu-mahou-no-machigatta-tsukaikata-senjou-wo-kakeru-kaifuku-youin-raw.html"
	info := []ChapterInfo{}

	c := colly.NewCollector(
		colly.AllowedDomains("hanascan.com"),
		colly.Async(true),
		colly.MaxDepth(10),
	)

	c.Limit(&colly.LimitRule{
		DomainGlob:  "*",
		Parallelism: 10,
		Delay:       5,
	})

	c.OnHTML("#list-chapters p", func(e *colly.HTMLElement) {
		tmp := ChapterInfo{}
		// Each P element in list-chapters has 3 spans. Index 0: contains <a>, index 1: contains latest update time. index 2
		// is not useful
		e.ForEach("span", func(count int, elem *colly.HTMLElement) {
			switch count {
			case 0:
				tmp.ChapterTitle = elem.ChildAttr("a.chapter", "title")
				tmp.URL = elem.ChildAttr("a.chapter", "href")
			case 1:
				tmp.LastUpdate = elem.ChildText("i time")
			}
		})
		info = append(info, tmp)
		fmt.Printf("%+v\n", tmp)
	})

	c.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting..", r.URL.String())
	})

	// Visiting
	c.Visit(url)
	c.Wait()
	return &info
}

func main() {
	ChapterInfoExtractor()
}
