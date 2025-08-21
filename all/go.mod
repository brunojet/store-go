module store-go/all

go 1.23.4

require (
	store-go/backoffice v0.0.0
	store-go/shared v0.0.0
	store-go/store v0.0.0
)

replace store-go/shared => ../shared

replace store-go/backoffice => ../backoffice

replace store-go/store => ../store
