package main

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"skillfactory/36/pkg/api"
	"skillfactory/36/pkg/rss"
	"skillfactory/36/pkg/storage"
	"skillfactory/36/pkg/storage/db"
	"time"
)

const DBURL = "postgres://postgres:postgres@127.0.0.1:8081/posts"

type server struct {
	db  storage.Interface
	api *api.API
}

type config struct {
	Period  int      `json:"request_period"`
	LinkArr []string `json:"rss"`
}

func main() {
	chanErr, chanPost := make(chan error), make(chan []storage.Post)
	conFile, err := ioutil.ReadFile("./config.json")
	if err != nil {
		log.Fatal(err)
	}
	var config config
	err = json.Unmarshal(conFile, &config)
	if err != nil {
		log.Fatal(err)
	}
	rssLinks := rssJSON("config.json", chanErr)
	for i := range rssLinks.LinkArr {
		go postParse(rssLinks.LinkArr[i], config.Period, chanErr, chanPost)
	}
	var serv server
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()
	newDB, err := db.New(ctx, DBURL)
	if err != nil {
		log.Fatal(err)
	}
	serv.db = newDB
	go func() { // обработка новостей
		for allPosts := range chanPost {
			for idx := range allPosts {
				newDB.AddPost(allPosts[idx])
			}
		}
	}()
	go func() { // обработка ошибок
		for err := range chanErr {
			log.Println("Ошибка:", err)
		}
	}()
	err = http.ListenAndServe(":80", serv.api.Router())
	if err != nil {
		log.Fatal(err)
	}
}

func postParse(link string, dur int, chanErr chan<- error, chanPost chan<- []storage.Post) {
	for {
		postList, err := rss.RSSStruct(link)
		if err != nil {
			chanErr <- err
			continue
		}
		chanPost <- postList
		time.Sleep(time.Duration(dur) * time.Minute)
	}
}

func rssJSON(file string, chanErr chan<- error) config {
	jsFile, err := os.Open(file)
	if err != nil {
		chanErr <- err
	}
	defer jsFile.Close()
	byteValue, _ := ioutil.ReadAll(jsFile)
	var linkList config
	json.Unmarshal(byteValue, &linkList)
	return linkList
}
