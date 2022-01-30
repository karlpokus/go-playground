package main

import (
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
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

func runChild(fpath string, w http.ResponseWriter) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	codeDirPathHost, _ := filepath.Abs("code")
	codeDirPathVolume := fmt.Sprintf("%s:/code", codeDirPathHost)
	codePathContainer := fmt.Sprintf("/%s", fpath)
	args := []string{
		"run", "--rm", "-v", "/usr/local/go:/go", "-v", codeDirPathVolume,
		"ubuntu:20.04", "/go/bin/go", "run", codePathContainer,
	}
	cmd := exec.CommandContext(ctx, "docker", args...)
	cmd.Stdout = w
	cmd.Stderr = w // yolo
	err := cmd.Start()
	if err != nil {
		return err
	}
	err = cmd.Wait()
	if err != nil {
		return err
	}
	return nil
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
	err = runChild(fpath, w)
	if err != nil {
		// child stdout is connected to w
		// just log err
		log.Printf("Error running child: %s", err)
		return
	}
}

func test(w http.ResponseWriter, r *http.Request) {
	err := runChild("code/test.go", w)
	if err != nil {
		log.Printf("Error running child: %s", err)
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
	http.HandleFunc("/test", test)
	port := ":9000"
	log.Printf("listening on %s", port)
	log.Fatal(http.ListenAndServe(port, nil))
}
