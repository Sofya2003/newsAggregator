package main

import (
	"context"
	"log"
	"skillfactory/36/pkg/api"
	"skillfactory/36/pkg/storage"
	"skillfactory/36/pkg/storage/db"
	"time"
)

//const DBURL = "postgres://postgres:postgres@127.0.0.1:8081/posts"
const DBURL = "user=postgres password=Keks17sql dbname=postsDB sslmode=disable"

type server struct {
	db  storage.Interface
	api *api.API
}

type config struct {
	Period  int      `json:"request_period"`
	LinkArr []string `json:"rss"`
}

func main() {
	var serv server
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()
	db, err := db.New(ctx, DBURL)
	if err != nil {
		log.Fatal(err)
	}
	serv.db = db
}
