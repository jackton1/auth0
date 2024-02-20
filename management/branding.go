package management

import (
	"log"
)

func Branding(domain string, token string, params []byte) (string, error) {
	log.Print("Retrieving branding...")
	return send("GET", apiURI(domain, "branding"), token, nil, params)
}

func BrandingUpdate(domain string, token string, data []byte) (string, error) {
	log.Print("Updating branding...")
	return send("PATCH", apiURI(domain, "branding"), token, data, nil)
}
