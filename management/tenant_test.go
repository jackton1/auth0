package management

import (
	"encoding/json"
	"fmt"
	"github.com/tidwall/gjson"
	"net/http"
	"os"
	"testing"
	"time"
)

func TestTenantsSettings(t *testing.T) {
	var err error

	domain := os.Getenv("AUTH0_DOMAIN")
	clientId := os.Getenv("AUTH0_CLIENT_ID")
	clientSecret := os.Getenv("AUTH0_CLIENT_SECRET")

	token := GetToken(domain, clientId, clientSecret)

	tenantSettingParams, _ := json.Marshal(map[string]string{
		"fields": "flags",
	})

	initialTenantsSettingsResponse, err := TenantsSettings(domain, token, tenantSettingParams)

	initialTenantData, _ := json.Marshal(map[string]interface{}{
		"flags": map[string]bool{
			"enable_custom_domain_in_emails": gjson.Get(initialTenantsSettingsResponse, "enable_custom_domain_in_emails").Bool(),
		},
	})

	if err != nil {
		t.Fatal(err)
	}

	t.Cleanup(func() {
		time.Sleep(time.Second)
		if _, err = TenantsSettingsUpdate(domain, token, initialTenantData); err != nil {
			t.Fatal(err)
		}
	})

	t.Run("Update", func(t *testing.T) {
		tenantData, _ := json.Marshal(map[string]interface{}{
			"flags": map[string]bool{
				"enable_custom_domain_in_emails": false,
			},
		})

		_, err = TenantsSettingsUpdate(domain, token, tenantData)

		if err != nil {
			t.Errorf("Error: %v", err)
		}

		expectedError := fmt.Sprintf("exited with %d status code", http.StatusBadRequest)

		_, err = TenantsSettingsUpdate(domain, token, []byte("Invalid"))

		if err != nil && err.Error() != expectedError {
			t.Errorf("Error: %v != %s", err, expectedError)
		}
	})
}
