package main

import (
	"flag"
	"fmt"
	"github.com/gin-gonic/gin"
	"strconv"
)

func main() {
	var port int
	var dir string

	flag.IntVar(&port, "port", 8080, "Port to serve on.")
	flag.StringVar(&dir, "dir", "", "Directory to serve.")
	flag.Parse()

	r := gin.Default()
	r.Static("/", dir)
	fmt.Println("Starting server")
	r.Run(":"+strconv.Itoa(port))
}
