package main

import (
	"embed"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
)

//go:embed html/build
var htmlFS embed.FS

var ServerPid int

var server *Server

func main() {
	/*tail := exec.Command("cat", "/var/mc-test/fifo")
	out, oerr := tail.StdoutPipe()
	if oerr != nil {
		panic(oerr)
	}
	terr := tail.Start()
	println(tail.Process.Pid)
	defer tail.Process.Kill()
	defer out.Close()
	if terr != nil {
		panic(terr)
	}*/

	java := exec.Command("java", "-jar", "server.jar", "nogui")
	java.Dir = "/var/mc-test/"
	//java.Stdin = out
	jerr := java.Start()
	defer java.Process.Kill()
	if jerr != nil {
		panic(jerr)
	}
	println(java.Process.Pid)
	ServerPid = java.Process.Pid

	http.Handle("GET /", html())
	http.Handle("GET /logs/", logs())
	http.HandleFunc("POST /server/start", ServerStart)
	http.HandleFunc("POST /server/stop", ServerStop)
	http.HandleFunc("POST /server/command", command)

	println(fmt.Printf("Hosted at http://%s:%s\n", "172.23.80.175", "8080"))
	panic(http.ListenAndServe(":8080", nil))
}

func html() http.Handler {

	return http.FileServerFS(htmlFS)
}

func logs() http.Handler {
	return http.StripPrefix("/logs/", http.FileServerFS(os.DirFS("/var/mc-test/")))
}

func command(w http.ResponseWriter, r *http.Request) {
	var err error
	if server == nil {
		http.Error(w, "Server not started", http.StatusInternalServerError)
		return
	}
	if r.Header.Get("Content-Type") != "text/plain" {
		http.Error(w, "Content-Type not supported", http.StatusUnsupportedMediaType)
		return
	}

	b, err := io.ReadAll(r.Body)
	b = append(b, '\n')
	if err != nil {
		panic(err)
	}

	stdin, err := server.Command.StdinPipe()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	_, err = io.WriteString(stdin, string(b))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	_, _ = w.Write([]byte(fmt.Sprintf("Sent to server: \n%s", string(b))))
}

func ServerStart(w http.ResponseWriter, r *http.Request) {
	err := server.Start()
	if err != nil {
		println(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func ServerStop(w http.ResponseWriter, r *http.Request) {
	err := server.Stop()
	if err != nil {
		println(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
