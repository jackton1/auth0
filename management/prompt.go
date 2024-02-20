package management

import (
	"log"
)

func Prompts(domain string, token string) (string, error) {
	log.Print("Retrieving prompts...")
	return send("GET", apiURI(domain, "prompts"), token, nil, nil)
}

func PromptsUpdate(domain string, token string, data []byte) (string, error) {
	log.Print("Updating universal login...")
	return send("PATCH", apiURI(domain, "prompts"), token, data, nil)
}

func PromptsCustomText(domain string, token string, prompt string) (string, error) {
	log.Printf("Retrieving '%s' prompt...", prompt)
	return send("GET", apiURI(domain, "prompts", prompt, "custom-text", "en"), token, nil, nil)
}

func PromptsCustomTextUpdate(domain string, token string, prompt string, data []byte) (string, error) {
	log.Printf("Updating '%s' prompt...", prompt)
	return send("PUT", apiURI(domain, "prompts", prompt, "custom-text", "en"), token, data, nil)
}
