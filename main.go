package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"mime"
	"net/http"
	"os"
	"path/filepath"
)

func main() {
	var port int
	var dir string
	var defaultFilename string

	flag.IntVar(&port, "port", 80, "Port to use")
	flag.StringVar(&dir, "dir", ".", "Relative or absolute directory location")
	flag.StringVar(&defaultFilename, "default", "", "File to serve on non-file request")
	flag.Parse()

	if !filepath.IsAbs(dir) {
		var err error
		dir, err = filepath.Abs(dir)
		if err != nil {
			fmt.Println(err.Error())
			os.Exit(1)
			return
		}
	}

	defaultFilename = filepath.Join(dir, defaultFilename)
	if stat, err := os.Stat(defaultFilename); os.IsNotExist(err) || stat.IsDir() {
		defaultFilename = ""
	}

	handler := Handler{
		Dir: dir,
		Default: defaultFilename,
	}
	err := http.ListenAndServe(fmt.Sprintf(":%d", port), &handler)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
}

type Handler struct {
	Dir string
	Default string
}

func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	filename := filepath.Join(h.Dir, r.URL.Path)

	// Do a check if we need to serve the default file instead for single page apps
	if h.Default != "" {
		if stat, err := os.Stat(filename); os.IsNotExist(err) || stat.IsDir() {
			filename = h.Default
		}
	}

	if ext := filepath.Ext(filename); len(ext) > 0 {
		w.Header().Add("Content-Type", mime.TypeByExtension(ext))
	}

	if stat, err := os.Stat(filename); os.IsNotExist(err) {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("File not found."))
		return
	} else if stat.IsDir() {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Directory check."))
		return
	} else if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	file, err := os.Open(filename)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}
	defer file.Close()

	w.WriteHeader(http.StatusOK)
	if _, err := io.Copy(w, file); err != nil {
		log.Println(err)
	}
}