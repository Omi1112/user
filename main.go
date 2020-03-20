package main

import (
	"github.com/SeijiOmi/user/db"
	"github.com/SeijiOmi/user/server"
)

func main() {
	db.Init()
	server.Init()
	db.Close()
}
