package main

import (
	backoffice "store-go/backoffice/pkg"
	shared "store-go/shared/pkg"
	store "store-go/store/pkg"
)

func main() {
	shared.HelloShared()
	backoffice.HelloBackoffice()
	store.HelloStore()
}
