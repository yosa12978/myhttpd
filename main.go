package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
)

var (
	addr      string
	staticDir string
)

func init() {
	flag.StringVar(&addr, "addr", "0.0.0.0:8080", "address")
	flag.StringVar(&staticDir, "dir", "./www", "static files directory")
}

func ServeStaticDir(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Path
	fulldir := staticDir + path

	info, err := os.Stat(fulldir)
	if err != nil {
		w.Write([]byte("<h1>404 not found</h1>"))
		return
	}
	if !info.IsDir() {
		http.FileServer(http.Dir(staticDir)).ServeHTTP(w, r)
		return
	}

	// Here the path is surely a directory
	if !strings.HasSuffix(path, "/") {
		path += "/"
	}
	dir, err := ioutil.ReadDir(fulldir)
	if err != nil {
		panic(err)
	}
	w.Write([]byte(fmt.Sprintf("<h1>Content of %s</h1><hr>", path)))
	if path != "/" {
		back := strings.Split(path, "/")
		bs := strings.Join(back[:(len(back)-2)], "/")
		if !strings.HasPrefix(bs, "/") {
			bs = "/" + bs
		}
		w.Write([]byte(fmt.Sprintf("<a href=\"%s\">< back</a><br><br>", bs)))
	}
	for _, v := range dir {
		if v.IsDir() {
			w.Write([]byte(fmt.Sprintf("<a href=\"%s%s\">%s/</a><br>", path, v.Name(), v.Name())))
			continue
		}
		w.Write([]byte(fmt.Sprintf("<a href=\"%s%s\">%s</a><br>", path, v.Name(), v.Name())))
	}
}

func main() {
	flag.Parse()

	serv := &http.Server{
		Addr:    addr,
		Handler: http.HandlerFunc(ServeStaticDir),
	}
	d, err := os.Stat(staticDir)
	if err != nil || !d.IsDir() {
		os.Mkdir(staticDir, os.ModePerm)
	}

	log.Println("server started")
	go serv.ListenAndServe()
	out := make(chan os.Signal, 1)
	signal.Notify(out, os.Interrupt, syscall.SIGTERM)
	<-out
	log.Println("server is shutted down")
}
