package main

import (
	"fmt"

	"github.com/gocolly/colly"
)

func main() {
	url := "https://hanascan.com/manga-chiyu-mahou-no-machigatta-tsukaikata-senjou-wo-kakeru-kaifuku-youin-raw.html"

	c := colly.NewCollector(
		colly.AllowedDomains("hanascan.com"),
	)

	c.OnHTML(".container .row .manga-info", func(e *colly.HTMLElement) {
		fmt.Println(e.ChildText("h1"))
		fmt.Println(e.ChildText("i"))
		//fmt.Println(e)
		e.ForEach("li", func(_ int, elem *colly.HTMLElement) {
			fmt.Println(elem.Text)
		})
		//fmt.Println(e.Text)
	})

	c.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting..", r.URL.String())
	})

	// Visiting
	c.Visit(url)
}
