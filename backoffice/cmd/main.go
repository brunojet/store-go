package main

import (
	"store-go/backoffice/internal"
	"store-go/shared/pkg"
)

func main() {
	pkg.HelloShared()
	internal.HelloBackoffice()
}
