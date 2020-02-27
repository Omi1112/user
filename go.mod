module github.com/SeijiOmi/gin-tamplate

go 1.13

replace local.packages/db => ./db

// replace local.packages/server => ./server

replace local.packages/controller => ./controller

replace local.packages/entity => ./entity

replace local.packages/service => ./service

require (
	cloud.google.com/go v0.37.4 // indirect
	github.com/jinzhu/gorm v1.9.12
	github.com/modern-go/concurrent v0.0.0-20180306012644-bacd9c7ef1dd // indirect
	github.com/modern-go/reflect2 v1.0.1 // indirect
	github.com/stretchr/testify v1.5.1
	golang.org/x/net v0.0.0-20190503192946-f4e77d36d62c // indirect
	gopkg.in/check.v1 v1.0.0-20180628173108-788fd7840127 // indirect
	gopkg.in/go-playground/validator.v8 v8.18.2 // indirect
	local.packages/db v0.0.0-00010101000000-000000000000
	local.packages/entity v0.0.0-00010101000000-000000000000 // indirect
	// local.packages/server v0.0.0-00010101000000-000000000000
	local.packages/service v0.0.0-00010101000000-000000000000 // indirect
)
