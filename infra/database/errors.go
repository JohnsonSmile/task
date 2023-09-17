package database

import "errors"

var (
	UserAddressAlreadyBindWithTwitter = errors.New("user address already bind with a twitter account")
	TwitterAlreadyBindWithUserAddress = errors.New("twitter account already bind with a user address")
)
