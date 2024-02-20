package management

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"testing"
	"time"
)

func TestBranding(t *testing.T) {
	var err error

	domain := os.Getenv("AUTH0_DOMAIN")
	clientId := os.Getenv("AUTH0_CLIENT_ID")
	clientSecret := os.Getenv("AUTH0_CLIENT_SECRET")

	token := GetToken(domain, clientId, clientSecret)

	initialBrandingResponse, err := Branding(domain, token, nil)

	if err != nil {
		t.Fatal(err)
	}

	t.Cleanup(func() {
		time.Sleep(time.Second)
		if _, err = BrandingUpdate(domain, token, []byte(initialBrandingResponse)); err != nil {
			t.Fatal(err)
		}
	})

	t.Run("Update", func(t *testing.T) {
		brandingData, _ := json.Marshal(map[string]interface{}{
			"logo_url":    "https://govfetch.com/img/logo/favicon.svg",
			"favicon_url": "https://govfetch.com/img/logo/favicon.svg",
			"colors": map[string]string{
				"primary":         "#006E75",
				"page_background": "#E4E4E7",
			},
		})

		_, err = BrandingUpdate(domain, token, brandingData)

		if err != nil {
			t.Errorf("Error: %v", err)
		}

		expectedError := fmt.Sprintf("exited with %d status code", http.StatusBadRequest)

		_, err = BrandingUpdate(domain, token, []byte("Invalid"))

		if err != nil && err.Error() != expectedError {
			t.Errorf("Error: %v != %s", err, expectedError)
		}
	})
}
