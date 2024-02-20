package management

import (
	"log"
)

func TenantsSettings(domain string, token string, params []byte) (string, error) {
	log.Print("Retrieving tenant settings...")
	return send("GET", apiURI(domain, "tenants", "settings"), token, nil, params)
}

func TenantsSettingsUpdate(domain string, token string, data []byte) (string, error) {
	log.Print("Updating tenant settings...")
	return send("PATCH", apiURI(domain, "tenants", "settings"), token, data, nil)
}
