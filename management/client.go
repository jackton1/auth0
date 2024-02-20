package management

import (
	"log"
)

func Client(domain string, token string, id string, params []byte) (string, error) {
	log.Print("Retrieving client...")
	return send("GET", apiURI(domain, "clients", id), token, nil, params)
}

func ClientUpdate(domain string, token string, id string, data []byte) (string, error) {
	log.Print("Updating client settings...")
	return send("PATCH", apiURI(domain, "clients", id), token, data, nil)
}
