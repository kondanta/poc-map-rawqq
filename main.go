package main

import (
	"fmt"

	"github.com/gocolly/colly"
)

type ChapterInfo struct {
	URL           string
	LastUpdate    string
	ChapterTitle  string
	ChapterImages *[]ChapterImage
}

type ChapterImage struct {
	ImageNo  int
	ImageUrl string
}

// Uniqfy the Slices
func unique(intSlice []ChapterInfo) []ChapterInfo {
	keys := make(map[ChapterInfo]bool)
	list := []ChapterInfo{}
	for _, entry := range intSlice {
		if _, value := keys[entry]; !value {
			keys[entry] = true
			list = append(list, entry)
		}
	}
	return list
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
	})

	c.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting..", r.URL.String())
	})

	// Visiting
	c.Visit(url)
	c.Wait()
	uniqInfo := unique(info)
	return &uniqInfo
}

func ExtractChapterImages(chapInfo *[]ChapterInfo) {
	baseURL := "https://hanascan.com/"

	c := colly.NewCollector(
		colly.AllowedDomains("hanascan.com"),
		// colly.Async(true),
		// colly.MaxDepth(10),
	)

	c.Limit(&colly.LimitRule{
		DomainGlob: "*",
		// Parallelism: 10,
		Delay: 5,
	})
	index := -1
	c.OnHTML("#content", func(e *colly.HTMLElement) {
		images := []ChapterImage{}
		tmp := ChapterImage{}
		e.ForEach("img", func(count int, elem *colly.HTMLElement) {
			tmp.ImageNo = count + 1
			tmp.ImageUrl = elem.Attr("data-original")
			images = append(images, tmp)
		})
		fmt.Printf("Index: %d\n", index)
		(*chapInfo)[index].ChapterImages = &images
	})

	c.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting..", r.URL.String())
	})

	c.OnResponse(func(r *colly.Response) {
		index++
	})

	for _, data := range *chapInfo {
		c.Visit(baseURL + data.URL)
		c.Wait()
	}
	//c.Wait()
	// fmt.Printf("%+v\n", images)

	fmt.Println("%+v\n", (*chapInfo))
}

func main() {
	chapInfo := ChapterInfoExtractor()
	ExtractChapterImages(chapInfo)
}
