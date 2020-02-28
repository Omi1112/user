package main

import (
	"github.com/SeijiOmi/gin-tamplate/db"
	"github.com/SeijiOmi/gin-tamplate/server"
)

func main() {
	db.Init()
	server.Init()
	db.Close()
}
