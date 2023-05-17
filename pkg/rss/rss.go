package rss

import (
	"encoding/xml"
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"skillfactory/36/pkg/storage"
	"time"
)

type PostItem struct {
	Title   string
	Content string
	PubDate template.HTML
	Link    string
}
type XMLStruct struct {
	PostList []PostItem
}

func RSSStruct(link string) ([]storage.Post, error) {
	var posts XMLStruct
	if xmlBytes, err := getXML(link); err != nil {
		log.Printf("Failed to get XML: %v", err)
	} else {
		xml.Unmarshal(xmlBytes, &posts)
	}
	var lst []storage.Post
	for i := range posts.PostList {
		var elem storage.Post
		elem.Title = posts.PostList[i].Title
		elem.Content = posts.PostList[i].Content
		elem.Link = posts.PostList[i].Link
		t, err := time.Parse("Mon 2 Jan 2006 15:04:05 -0700", string(posts.PostList[i].PubDate))
		if err != nil {
			t, err = time.Parse("Mon 2 Jan 2006 15:04:05 GMT", string(posts.PostList[i].PubDate))
		}
		if err == nil {
			elem.PubTime = t.Unix()
		}
		lst = append(lst, elem)
	}
	return lst, nil
}

func getXML(url string) ([]byte, error) {

	res, err := http.Get(url)
	if err != nil {
		return []byte{}, fmt.Errorf("GET error: %v", err)
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return []byte{}, fmt.Errorf("Status error: %v", res.StatusCode)
	}

	data, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return []byte{}, fmt.Errorf("Read body: %v", err)
	}

	return data, nil
}
