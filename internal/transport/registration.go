package transport

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/OpsNexusHQ/opsnexus-common/models"
)

// RegistrationClient handles Agent registration with the OpsNexus backend.
type RegistrationClient struct {
	baseURL    string
	httpClient *http.Client
}

// NewRegistrationClient creates a registration client.
func NewRegistrationClient(baseURL string) *RegistrationClient {
	return &RegistrationClient{
		baseURL: strings.TrimRight(baseURL, "/"),
		httpClient: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

// Register registers an Agent with the OpsNexus backend.
func (c *RegistrationClient) Register(ctx context.Context, agent models.Agent) error {
	payload, err := json.Marshal(agent)
	if err != nil {
		return fmt.Errorf("marshal agent registration: %w", err)
	}

	url := c.baseURL + "/api/v1/agents/register"

	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodPost,
		url,
		bytes.NewReader(payload),
	)
	if err != nil {
		return fmt.Errorf("create registration request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("send registration request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		return fmt.Errorf("registration failed with HTTP status %d", resp.StatusCode)
	}

	return nil
}
