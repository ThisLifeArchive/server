package episodes

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
)

const thisLifeURL = "https://www.thisamericanlife.org"
const episodesURL = thisLifeURL + "/archive"

type Episode struct {
	Name        string    `json:"name"`
	Description string    `json:"description,omitempty"`
	Date        time.Time `json:"date,omitempty"`
	PageURL     string    `json:"pageURL"`
	AudioURL    string    `json:"audioURL,omitempty"`
	Image       *Image    `json:"image,omitempty"`
}

type Image struct {
	Source string `json:"source"`
	Width  int    `json:"width"`
	Height int    `json:"height"`
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
		container := s.Find("header > .container")

		link := container.Find("h2 > a")
		name := link.Text()
		if len(name) == 0 {
			fmt.Printf("WARN: text not found in article element at position %d\n", i)
			return
		}

		pageURL, exists := link.Attr("href")
		if !exists {
			fmt.Printf("WARN: href not found on link at position %d with name `%s`\n", i, name)
			return
		}

		date := container.Find("span.date-display-single").Text()
		origMonth := strings.Split(date, " ")[0]
		month, hadSuffix := strings.CutSuffix(origMonth, ".")
		monthLayout := "January"
		if hadSuffix {
			monthLayout = "Jan."
			date = strings.Replace(date, origMonth, fmt.Sprintf("%s.", month[0:3]), 1)
		}
		d, err := time.Parse(fmt.Sprintf("%s 2, 2006", monthLayout), date)
		if err != nil {
			fmt.Printf("WARN: failed to parse date string `%s`: %v\n", date, err)
			return
		}

		desc := s.Find("div.content .field-item > p")

		episode := Episode{
			Name:        name,
			PageURL:     thisLifeURL + pageURL,
			Description: desc.Text(),
			Date:        d,
		}

		img := s.Find("header figure.episode-image > img")
		if img != nil {
			src, _ := img.Attr("src")
			width, _ := img.Attr("width")
			height, _ := img.Attr("height")
			w, err := strconv.Atoi(width)
			h, err := strconv.Atoi(height)
			if err == nil {
				episode.Image = &Image{
					Source: src,
					Width:  w,
					Height: h,
				}
			} else {
				fmt.Printf("failed to convert string to int: %v\n", err)
			}
		}

		eps = append(eps, &episode)
	})
	return eps, nil
}
