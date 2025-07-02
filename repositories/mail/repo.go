package mail

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/vivasoft-ltd/go-ems/config"
	"github.com/vivasoft-ltd/go-ems/types"
)

type Repository struct {
	client *http.Client
	config *config.EmailConfig
}

func NewRepository(client *http.Client, config *config.EmailConfig) *Repository {
	return &Repository{
		client: client,
		config: config,
	}
}

func (repo *Repository) SendEmail(payload *types.EmailPayload) error {
	reqByte, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("failed to marshal email payload: %w", err)
	}

	req, err := http.NewRequest(http.MethodPost, repo.config.Url, bytes.NewBuffer(reqByte))
	if err != nil {
		return fmt.Errorf("failed to create email request: %w", err)
	}

	res, err := repo.client.Do(req)
	if err != nil {
		return fmt.Errorf("error sending email to %s: %w", payload.MailTo, err)
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return fmt.Errorf("email service returned status code %d for recipient %s", res.StatusCode, payload.MailTo)
	}

	return nil
}
