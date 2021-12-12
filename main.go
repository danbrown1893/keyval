package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"time"

	"keyval/cache"

	"github.com/gorilla/mux"
)

var localCache *cache.Cache

type PostBody struct {
	Key      string      `json:"key"`
	Value    interface{} `json:"value"`
	Deadline string      `json:"deadline"`
}

func main() {
	router := mux.NewRouter()
	InitCache()

	router.HandleFunc("/store", Set).Methods("POST")
	router.HandleFunc("/store", Get).Methods("GET")
	router.HandleFunc("/store", Delete).Methods("DELETE")
	http.Handle("/", router)
	server := &http.Server{
		Handler: router,
		Addr:    "127.0.0.1:3000",
	}

	log.Fatal(server.ListenAndServe())
}

func InitCache() {
	localCache = cache.NewCache()
}

func Set(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, fmtErr(err), http.StatusInternalServerError)
		return
	}
	var tmp PostBody

	if err := json.Unmarshal(body, &tmp); err != nil {
		http.Error(w, fmtErr(err), http.StatusInternalServerError)
		return
	}

	deadline, err := handleDeadline(tmp.Deadline)
	if err != nil {
		http.Error(w, fmtErr(err), http.StatusInternalServerError)
		return
	}

	if tmp.Key != "" && tmp.Value != "" {
		v, _, err := localCache.EnsureKey(tmp.Key, tmp.Value, deadline)
		if err != nil {
			http.Error(w, fmtErr(err), http.StatusInternalServerError)
			return
		}
		res, err := json.Marshal(v)
		if err != nil {
			http.Error(w, fmtErr(err), http.StatusInternalServerError)
			return
		}
		w.Write(res)
	}
}

func Get(w http.ResponseWriter, r *http.Request) {
	v, err := localCache.Get(r.URL.RawQuery)
	if err != nil {
		http.Error(w, fmtErr(err), http.StatusInternalServerError)
		return
	}
	enc, err := json.Marshal(v)
	if err != nil {
		http.Error(w, fmtErr(err), http.StatusInternalServerError)
		return
	}
	w.Write(enc)
}

func Delete(w http.ResponseWriter, r *http.Request) {
	localCache.Delete(r.URL.RawQuery)
	w.WriteHeader(http.StatusOK)
}

func fmtErr(err error) string {
	return fmt.Sprintf("Error: %s", err.Error())
}

func handleDeadline(ts string) (time.Time, error) {
	if ts == "" {
		return time.Time{}, nil
	}
	i, err := strconv.ParseInt(ts, 10, 64)
	if err != nil {
		return time.Time{}, err
	}

	return time.Unix(i, 0), nil
}
