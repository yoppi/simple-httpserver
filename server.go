package main

import (
	"flag"
	"fmt"
	"net/http"
	"time"
	"log"
)

func logDate(now time.Time) string {
	return fmt.Sprintf("%02d/%3s/%04d:%02d:%02d:%02d", now.Day(), now.Month(), now.Year(), now.Hour(), now.Minute(), now.Second())
}

func accessLog(r *http.Request, code int, size interface{}) {
	fmt.Printf("%s - - [%s] \"%s %s\" %d %v\n", r.RemoteAddr, logDate(time.Now()), r.Method, r.URL.Path, code, size)
}

func simpleHTTPHandler() func(http.ResponseWriter, *http.Request) {
	dir := http.Dir(".")
	return func(w http.ResponseWriter, r *http.Request) {
		path := r.URL.Path[1:]
		f, err := dir.Open(path)
		if err != nil {
			accessLog(r, http.StatusNotFound, "-")
			http.NotFound(w, r)
		} else {
			defer f.Close()
			stat, _ := f.Stat()
			accessLog(r, http.StatusOK, stat.Size())
			if stat.IsDir() {
				http.ServeFile(w, r, path)
			} else {
				http.ServeContent(w, r, path, stat.ModTime(), f)
			}
		}
	}
}

func parseArgs() string {
	var port string
	flag.StringVar(&port, "port", "8000", "Serving HTTP port number")
	flag.Parse()
	return port
}

func main() {
	http.HandleFunc("/", simpleHTTPHandler())
	port := parseArgs()
	fmt.Printf("Serving HTTP on localhost port %v\n", port)
	err := http.ListenAndServe(":" + port, nil)
	if err != nil {
		panic(err)
	}
}
