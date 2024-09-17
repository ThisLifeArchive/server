package episodes

import (
	"fmt"
	"net/http"
	"time"

	"github.com/PuerkitoBio/goquery"
)

const episodesURL = "https://www.thisamericanlife.org/archive"

type Episode struct {
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Date        time.Time `json:"date"`
	PageURL     string    `json:"pageURL"`
	AudioURL    string    `json:"audioURL,omitempty"`
}

func List() ([]*Episode, error) {
	resp, err := http.Get(episodesURL)
	if err != nil {
		return nil, fmt.Errorf("http GET request failed: %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected http status %d: %s", resp.StatusCode, resp.Status)
	}
	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read html: %w", err)
	}
	var eps []*Episode
	doc.Find("article.node").Each(func(i int, s *goquery.Selection) {
		link := s.Find("header > .container > h2 > a")
		name := link.Text()
		if len(name) == 0 {
			fmt.Printf("WARN: text not found in article element at position %d\n", i)
			return
		}
		pageURL, exists := link.Attr("href")
		if !exists {
			fmt.Printf("WARN: href not found on link at position %d with name `%s`", i, name)
			return
		}
		episode := Episode{
			Name:    name,
			PageURL: pageURL,
		}
		eps = append(eps, &episode)
	})
	return eps, nil
}
