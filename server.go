package main

import (
	"context"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/exec"
	"time"
)

func file(b []byte) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Write(b)
	}
}

func files() ([]byte, []byte, error) {
	html, err := os.ReadFile("index.html")
	if err != nil {
		return nil, nil, err
	}
	js, err := os.ReadFile("index.js")
	if err != nil {
		return nil, nil, err
	}
	return html, js, err
}

func code(w http.ResponseWriter, r *http.Request) {
	b, err := ioutil.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	fpath := "code/test.go"
	err = os.WriteFile(fpath, b, 0644)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	// execute on host for now
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	cmd := exec.CommandContext(ctx, "go", "run", fpath)
	cmd.Stdout = w
	err = cmd.Start()
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	err = cmd.Wait()
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
}

func main() {
	html, js, err := files()
	if err != nil {
		log.Fatal(err)
	}
	http.HandleFunc("/", file(html))
	http.HandleFunc("/index.js", file(js))
	http.HandleFunc("/code", code)
	port := ":9000"
	log.Printf("listening on %s", port)
	log.Fatal(http.ListenAndServe(port, nil))
}
