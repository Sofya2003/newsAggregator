package api

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"net/http"
	"skillfactory/36/pkg/storage"
	"strconv"
)

type API struct {
	r  *mux.Router
	db storage.Interface
}

func (api *API) Router() *mux.Router {
	return api.r
}

// Регистрация методов API в маршрутизаторе запросов.
func (api *API) endpoints() {
	// получить n последних новостей
	api.r.HandleFunc("/news/{n}", api.posts).Methods(http.MethodGet, http.MethodOptions)
	// веб-приложение
	api.r.PathPrefix("/").Handler(http.StripPrefix("/", http.FileServer(http.Dir("./webapp"))))
}

func New(db storage.Interface) *API {
	api := API{db: db}
	api.r = mux.NewRouter()
	api.endpoints()
	return &api
}

func (api *API) posts(w http.ResponseWriter, r *http.Request) {
	s := mux.Vars(r)["n"]
	n, _ := strconv.Atoi(s)

	news, err := api.db.GetPosts(n)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(news)
}
