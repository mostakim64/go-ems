package conn

import (
	"net/http"
	"time"

	"github.com/vivasoft-ltd/go-ems/config"
)

var emailClient *http.Client

func ConnectEmail() {
	config := config.Email()
	timeout := config.Timeout * time.Second
	emailClient = newHTTPClient(timeout, 50)
}

func EmailClient() *http.Client {
	return emailClient
}
