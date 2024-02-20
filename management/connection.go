package management

import (
	"log"
)

func Connections(domain string, token string, params []byte) (string, error) {
	log.Print("Retrieving connections...")
	return send("GET", apiURI(domain, "connections"), token, nil, params)
}

func ConnectionCreate(domain string, token string, data []byte, ignoreStatusCode ...int) (string, error) {
	log.Print("Creating connection...")
	return send("POST", apiURI(domain, "connections"), token, data, nil, ignoreStatusCode...)
}

func ConnectionUpdate(domain string, token string, id string, data []byte) (string, error) {
	log.Print("Updating connection...")
	return send("PATCH", apiURI(domain, "connections", id), token, data, nil)
}

func ConnectionDelete(domain string, token string, id string) (string, error) {
	log.Print("Deleting connection...")
	return send("DELETE", apiURI(domain, "connections", id), token, nil, nil)
}
