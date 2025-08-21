package main

import (
	"store-go/shared/pkg"
	"store-go/store/internal"
)

func main() {
	pkg.HelloShared()
	internal.HelloStore()
}
