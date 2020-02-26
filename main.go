package main

import (
	"local.packages/db"
	"local.packages/server"
)

func main() {
	db.Init()
	server.Init()
	db.Close()
}
