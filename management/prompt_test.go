package management

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"testing"
	"time"
)

func TestPrompts(t *testing.T) {
	var err error

	domain := os.Getenv("AUTH0_DOMAIN")
	clientId := os.Getenv("AUTH0_CLIENT_ID")
	clientSecret := os.Getenv("AUTH0_CLIENT_SECRET")

	token := GetToken(domain, clientId, clientSecret)

	initialPromptsResponse, err := Prompts(domain, token)

	if err != nil {
		t.Fatal(err)
	}

	t.Cleanup(func() {
		time.Sleep(time.Second)
		if _, err = PromptsUpdate(domain, token, []byte(initialPromptsResponse)); err != nil {
			t.Fatal(err)
		}
	})

	t.Run("Update", func(t *testing.T) {
		promptsData, _ := json.Marshal(map[string]interface{}{
			"universal_login_experience": "new",
			"identifier_first":           true,
		})

		_, err = PromptsUpdate(domain, token, promptsData)

		if err != nil {
			t.Errorf("Error: %v", err)
		}

		expectedError := fmt.Sprintf("exited with %d status code", http.StatusBadRequest)

		_, err = PromptsUpdate(domain, token, []byte("Invalid"))

		if err != nil && err.Error() != expectedError {
			t.Errorf("Error: %v != %s", err, expectedError)
		}
	})
}

func TestPromptsCustomText(t *testing.T) {
	var err error

	domain := os.Getenv("AUTH0_DOMAIN")
	clientId := os.Getenv("AUTH0_CLIENT_ID")
	clientSecret := os.Getenv("AUTH0_CLIENT_SECRET")

	token := GetToken(domain, clientId, clientSecret)
	prompt := "consent"

	initialPromptsCustomTextResponse, err := PromptsCustomText(domain, token, prompt)

	if err != nil {
		t.Fatal(err)
	}

	t.Cleanup(func() {
		if _, err = PromptsCustomTextUpdate(domain, token, prompt, []byte(initialPromptsCustomTextResponse)); err != nil {
			t.Fatal(err)
		}
	})

	t.Run("CustomTextUpdate", func(t *testing.T) {
		consentData, _ := json.Marshal(map[string]interface{}{
			"consent": map[string]string{
				"pageTitle":              "Authorize | Gov Fetch",
				"title":                  "Hello",
				"messageMultipleTenants": "We are requesting access to your Gov Fetch account.",
				"messageSingleTenant":    "We are requesting access to your Gov Fetch account.",
			},
		})

		_, err = PromptsCustomTextUpdate(domain, token, prompt, consentData)

		if err != nil {
			t.Errorf("Error: %v", err)
		}

		expectedError := fmt.Sprintf("exited with %d status code", http.StatusBadRequest)

		_, err = PromptsCustomTextUpdate(domain, token, prompt, []byte("Invalid"))

		if err != nil && err.Error() != expectedError {
			t.Errorf("Error: %v != %s", err, expectedError)
		}
	})
}
