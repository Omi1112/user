package main

import (
	"local.packages/db"
	"github.com/SeijiOmi/gin-tamplate/server"
)

func main() {
	db.Init()
	server.Init()
	db.Close()
}
