module github.com/brunojet/store-go/shared

go 1.23.4

require (
	github.com/jinzhu/inflection v1.0.0 // indirect
	github.com/jinzhu/now v1.1.5 // indirect
	golang.org/x/sys v0.28.0 // indirect
	golang.org/x/text v0.21.0 // indirect
)

require (
	github.com/brunojet/store-go/infra v0.0.0
	gorm.io/gorm v1.30.1
)

replace github.com/brunojet/store-go/infra => ../infra
