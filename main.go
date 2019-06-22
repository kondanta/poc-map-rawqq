package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/gocolly/colly"
)

var baseURL = "https://hanascan.com/"

// Manga ...
type Manga struct {
	Title       string
	URL         string
	LastChapter int
	ChapterInfo *[]ChapterInfo
}

// ChapterInfo holds chapter's static information along with the image links
type ChapterInfo struct {
	URL           string
	LastUpdate    string
	ChapterTitle  string
	ChapterImages *[]ChapterImage
}

// ChapterImage holds image information of the specified  chapter
type ChapterImage struct {
	ImageNo  int
	ImageURL string
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

// ChapterInfoExtractor extracts the static information:
// Name, Last update and Chapter Name/No
func ChapterInfoExtractor(chapterURL string) *[]ChapterInfo {
	url := baseURL + chapterURL
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
	// extractChapterImages works with the address of the ChapterInfo variable
	// so we can call this function here and return one single unified variable.
	// extractChapterImages(&uniqInfo)
	return &uniqInfo
}

// extractChapterImages strip image links from the given chapter url
func extractChapterImages(chapInfo *[]ChapterInfo) {

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
			tmp.ImageURL = elem.Attr("data-original")
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

	// fmt.Println("%+v\n", (*chapInfo))
}

// TODO: Create robot.txt for crawling all existing mangas.
// TODO: Think about the database and its integration

// Searchmanga searches manga with the given name
func Searchmanga(mangaName string) *[]Manga {
	mangaName = querySanitizer(mangaName)
	query := "manga-list.html?m_status=&author=&group=&name=" + mangaName + "&genre=&ungenre="
	url := baseURL + query
	mangas := []Manga{}

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

	c.OnHTML(".row .top .media-body", func(e *colly.HTMLElement) {
		manga := Manga{}
		// Finding the latest chapter on search page
		e.ForEach("a", func(_ int, elem *colly.HTMLElement) {
			num := elem.Text
			// This bunch has also genres such as Action, Adventure. The last element is the
			// latest chapter.
			if res, err := strconv.Atoi(num); err == nil {
				manga.LastChapter = res
			}
		})

		// Create the view for search result
		e.ForEach("#tables", func(_ int, elem *colly.HTMLElement) {
			manga.URL = elem.ChildAttr("a", "href")
			manga.Title = elem.Text
			manga.ChapterInfo = nil /*ChapterInfoExtractor(manga.URL)*/
			mangas = append(mangas, manga)
		})

	})

	c.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting..", r.URL.String())
	})

	c.Visit(url)
	c.Wait()
	// fmt.Println("%+v\n", mangas)

	return &mangas
}

func querySanitizer(query string) string {
	return strings.Replace(query, " ", "+", -1)
}

func main() {
	// ChapterInfoExtractor is for mangapplizer. As for tracker, its not needed "yet"
	//chapInfo := ChapterInfoExtractor()
	//ExtractChapterImages(chapInfo)
	if len(os.Args) > 1 {
		Searchmanga(os.Args[1])
	} else {
		panic("Need a name for the manga that is going to be searched!")
	}
}
