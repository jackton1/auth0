[![Test](https://github.com/jackton1/auth0/actions/workflows/test.yml/badge.svg)](https://github.com/jackton1/auth0/actions/workflows/test.yml) 
![Coverage](https://img.shields.io/badge/Coverage-93.8%25-brightgreen)
[![Lint](https://github.com/jackton1/auth0/actions/workflows/lint.yml/badge.svg)](https://github.com/jackton1/auth0/actions/workflows/lint.yml)

# auth0
Auth0 API SDK 

## Installation

```bash
go get -u github.com/jackton1/auth0
```

## Usage

```golang
import (
    "github.com/jackton1/auth0/management"
    "os"
)

func main() {
    domain := os.Getenv("AUTH0_DOMAIN")
    clientId := os.Getenv("AUTH0_CLIENT_ID")
    clientSecret := os.Getenv("AUTH0_CLIENT_SECRET")

    token := management.GetToken(domain, clientID, clientSecret)

    users, err := management.Users(domain, token, nil)
}
```
