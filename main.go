package main

import (
	"flag"
	"os"
	"path"
	"path/filepath"
	"strconv"

	"github.com/gin-gonic/gin"
)

func main() {
	var port int
	var dir string
	var defaultFilename string

	flag.IntVar(&port, "port", 8080, "Port to use")
	flag.StringVar(&dir, "dir", "", "Relative or absolute directory location")
	flag.StringVar(&defaultFilename, "default", "", "File to serve on non-file")
	flag.Parse()

	if dir == "" {
		dir = "."
	}
	if !filepath.IsAbs(dir) {
		var err error
		dir, err = filepath.Abs(dir)
		if err != nil {
			panic(err)
		}
	}

	defaultFilename = filepath.Join(dir, defaultFilename)
	if stat, err := os.Stat(defaultFilename); os.IsNotExist(err) || stat.IsDir() {
		defaultFilename = ""
	}

	gin.SetMode(gin.ReleaseMode)
	r := gin.New()
	r.Use(gin.Logger())
	r.Use(func(c *gin.Context) {
		p := path.Join(dir, c.Request.URL.Path)
		if len(filepath.Ext(p)) > 0 {
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
