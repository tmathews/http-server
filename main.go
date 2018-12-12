package main

import (
	"flag"
	"fmt"
	"os"
	"path"
	"path/filepath"
	"strconv"

	"github.com/gin-contrib/gzip"
	"github.com/gin-gonic/gin"
)

func main() {
	var port int
	var dir string
	var defaultFilename string

	flag.IntVar(&port, "port", 8080, "Port to use")
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

	gin.SetMode(gin.ReleaseMode)
	r := gin.New()
	r.Use(gin.Logger())
	r.Use(gzip.Gzip(gzip.DefaultCompression))
	r.Use(func(c *gin.Context) {
		p := path.Join(dir, c.Request.URL.Path)
		ext := filepath.Ext(p)
		if len(ext) > 0 {
			if ext == ".wasm" {
				c.Header("Content-Type", "application/wasm")
			}
			c.File(p)
		} else if _, err := os.Stat(p); os.IsNotExist(err) && defaultFilename != "" {
			c.File(defaultFilename)
		} else {
			c.File(p)
		}
		c.Next()
	})
	r.Run(":" + strconv.Itoa(port))
}
