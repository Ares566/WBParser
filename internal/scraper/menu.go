package scraper

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
)

type Categories struct {
	categoryURLs []string
}

func NewCategories() *Categories {
	return &Categories{}
}

type mainMenuModel []struct {
	ID      int    `json:"id"`
	Name    string `json:"name"`
	URL     string `json:"url"`
	Shard   string `json:"shard,omitempty"`
	Query   string `json:"query,omitempty"`
	Landing bool   `json:"landing,omitempty"`
	Childs  []struct {
		ID     int    `json:"id"`
		Parent int    `json:"parent"`
		Name   string `json:"name"`
		Seo    string `json:"seo,omitempty"`
		URL    string `json:"url"`
		Shard  string `json:"shard"`
		Query  string `json:"query"`
		Childs []struct {
			ID     int    `json:"id"`
			Parent int    `json:"parent"`
			Name   string `json:"name"`
			URL    string `json:"url"`
			Shard  string `json:"shard"`
			Query  string `json:"query"`
			Seo    string `json:"seo,omitempty"`
		} `json:"childs,omitempty"`
		IsDenyLink bool `json:"isDenyLink,omitempty"`
	} `json:"childs,omitempty"`
	Seo  string `json:"seo,omitempty"`
	Dest []int  `json:"dest,omitempty"`
}

func (c *Categories) Parse() (urls []string, err error) {
	mainMenuURL := "https://static.wbstatic.net/data/main-menu-ru-ru.json"

	var res *http.Response
	res, err = http.Get(mainMenuURL)
	if err != nil {
		return
	}

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return
	}

	var data mainMenuModel
	err = json.Unmarshal(body, &data)
	if err != nil {
		return
	}

	for _, mainCat := range data {
		for _, child := range mainCat.Childs {
			urls = append(urls, child.URL)
		}
	}

	return
}
