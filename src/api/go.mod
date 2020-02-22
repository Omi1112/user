module main

go 1.13

replace local.packages/db => ./db

replace local.packages/server => ./server

replace local.packages/controller => ./controller

replace local.packages/entity => ./entity

replace local.packages/service => ./service

require (
	github.com/gin-gonic/gin v1.4.0 // indirect
	github.com/jinzhu/gorm v1.9.10
	local.packages/controller v0.0.0-00010101000000-000000000000 // indirect
	local.packages/db v0.0.0-00010101000000-000000000000
	local.packages/entity v0.0.0-00010101000000-000000000000 // indirect
	local.packages/server v0.0.0-00010101000000-000000000000
	local.packages/service v0.0.0-00010101000000-000000000000 // indirect
)
