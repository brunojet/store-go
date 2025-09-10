module github.com/brunojet/store-go/office

go 1.23.4

replace github.com/brunojet/store-go/shared => ../shared

require (
	github.com/brunojet/store-go/shared v0.0.0
	gorm.io/gorm v1.30.1
)

require (
	filippo.io/edwards25519 v1.1.0 // indirect
	github.com/go-sql-driver/mysql v1.8.1 // indirect
)

require (
	github.com/brunojet/infra-go v0.0.4 // indirect
	github.com/jinzhu/inflection v1.0.0 // indirect
	github.com/jinzhu/now v1.1.5 // indirect
	golang.org/x/text v0.28.0 // indirect
	gorm.io/driver/mysql v1.6.0
)
