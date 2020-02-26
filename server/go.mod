module server

go 1.13

replace local.packages/db => ../db

replace local.packages/server => ../server

replace local.packages/controller => ../controller

replace local.packages/entity => ../entity

replace local.packages/service => ../service

require (
	github.com/gin-gonic/gin v1.5.0
	github.com/stretchr/testify v1.5.1
	local.packages/controller v0.0.0-00010101000000-000000000000
	local.packages/db v0.0.0-00010101000000-000000000000
	local.packages/server v0.0.0-00010101000000-000000000000
)
