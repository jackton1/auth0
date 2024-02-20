package management

import (
	"log"
)

func Users(domain string, token string, params []byte) (string, error) {
	log.Print("Retrieving users...")
	return send("GET", apiURI(domain, "users"), token, nil, params)
}

func UsersByEmail(domain string, token string, params []byte) (string, error) {
	log.Printf("Retrieving users with email...")
	return send("GET", apiURI(domain, "users-by-email"), token, nil, params)
}

func UserCreate(domain string, token string, data []byte) (string, error) {
	log.Printf("Creating user...")
	return send("POST", apiURI(domain, "users"), token, data, nil)
}

func UserDelete(domain string, token string, id string) (string, error) {
	log.Printf("Deleting user...")
	return send("DELETE", apiURI(domain, "users", id), token, nil, nil)
}
