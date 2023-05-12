package main

import (
	"flag"
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

type Myhttpd struct {
}

func (router *Myhttpd) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	fulldir := staticDir + r.URL.Path
	info, err := os.Stat(fulldir)
	if err != nil {
		w.Write([]byte("<h1>404 not found</h1>"))
		return
	}
	if !info.IsDir() {
		http.FileServer(http.Dir(staticDir)).ServeHTTP(w, r)
		return
	}
	//if r.URL.Path[1:] == "" || info.IsDir()
	path := r.URL.Path
	if !strings.HasSuffix(path, "/") {
		path += "/"
	}
	dir, err := ioutil.ReadDir(fulldir)
	if err != nil {
		panic(err)
	}
	w.Write([]byte("<h1>Content of " + path + "</h1><hr>"))
	if path != "/" {
		back := strings.Split(path, "/")
		bs := strings.Join(back[:(len(back)-2)], "/")
		if !strings.HasPrefix(bs, "/") {
			bs = "/" + bs
		}
		w.Write([]byte("<a href=\"" + bs + "\">< back</a><br><br>"))
	}
	for _, v := range dir {
		if v.IsDir() {
			w.Write([]byte("<a href=\"" + path + v.Name() + "\">" + v.Name() + "/</a><br>"))
			continue
		}
		w.Write([]byte("<a href=\"" + path + v.Name() + "\">" + v.Name() + "</a><br>"))
	}
}

func main() {
	flag.Parse()

	mux := &Myhttpd{}
	serv := &http.Server{
		Addr:    addr,
		Handler: mux,
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
