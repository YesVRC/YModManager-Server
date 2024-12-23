package main

import (
	"embed"
	"github.com/a-h/templ"
	"log"
	"net/http"
	"os"
)

//go:embed public
var publicFS embed.FS

func main() {
	xfs := os.DirFS("X:\\YeSMP2")
	text, err := xfs.Open("user_jvm_args.txt")
	if err != nil {
		panic(err.Error())
	}
	stat, serr := text.Stat()
	if serr != nil {
		panic(serr.Error())
	}
	data := make([]byte, stat.Size())
	text.Read(data)
	log.Println(string(data))
	http.Handle("GET /public/", public())
	http.Handle("GET /", templ.Handler(base.Index()))
	http.Handle("GET /fs/", xOS())
	http.ListenAndServe(":8080", nil)
}

func public() http.Handler {
	return http.FileServerFS(publicFS)
}
func xOS() http.Handler {
	return http.StripPrefix("/fs/", http.FileServer(http.Dir("X:\\YeSMP2")))
}

type JavaVersion struct {
	Version int    `json:"version"`
	Title   string `json:"title"`
	Path    string `json:"path"`
}

type ServerType int

const (
	Vanilla ServerType = iota
	Forge
	Fabric
)

var ServerModloader = map[ServerType]string{
	Vanilla: "Vanilla",
	Forge:   "Forge",
	Fabric:  "Fabric",
}

type ServerInfo struct {
	Id         int        `json:"id"`
	McVersion  string     `json:"mc_version"`
	ServerType ServerType `json:"server_type"`
}
