package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"time"

	"github.com/gocolly/colly"
)

// main() contains code adapted from example found in Colly's docs:
// http://go-colly.org/docs/examples/basic/
type GGEpisodes struct {
	Title    string
	Director string
	Rating   string
}

func createJson(entry []GGEpisodes) {
	jsonFile, _ := json.MarshalIndent(entry, "", " ")
	_ = ioutil.WriteFile("output.json", jsonFile, 0644)
}

func main() {
	c := colly.NewCollector()
	c.SetRequestTimeout(120 * time.Second)

	episode := make([]GGEpisodes, 0)
	//create callback for links
	c.OnHTML("table.wikitable.plainrowheaders.wikiepisodetable", func(e *colly.HTMLElement) {
		e.ForEach("tr.vevent", func(_ int, e *colly.HTMLElement) {
			newEntry := GGEpisodes{}
			newEntry.Title = e.ChildText("a")
			newEntry.Director = e.ChildText("c")
			newEntry.Rating = e.ChildText("g")
			episode = append(episode, newEntry)
		})
	})

	c.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting", r.URL)
	})

	c.OnResponse(func(r *colly.Response) {
		fmt.Println("Response Received", r.Request.URL)
	})

	c.OnError(func(r *colly.Response, e error) {
		fmt.Println("Error:", e)
	})

	// start scraping
	c.Visit("https://en.wikipedia.org/wiki/List_of_The_Golden_Girls_episodes")

	createJson(episode)
}
