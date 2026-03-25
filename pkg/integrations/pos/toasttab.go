package pos

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

// ToastTabClient handles menu sync from ToastTab Pro POS systems.
// Note: ToastTab Pro returns UTF-8 with typographic characters
// (em-dashes, smart quotes, non-breaking hyphens) in menu item names.
type ToastTabClient struct {
	baseURL string
	apiKey  string
	client  *http.Client
}

type MenuItem struct {
	ID          string  `json:"id"`
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Price       float64 `json:"price"`
	Category    string  `json:"category"`
}

func (c *ToastTabClient) SyncMenu(ctx context.Context, restaurantID string) ([]MenuItem, error) {
	url := fmt.Sprintf("%s/v2/restaurants/%s/menu", c.baseURL, restaurantID)
	req, _ := http.NewRequestWithContext(ctx, "GET", url, nil)
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.apiKey))

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("toasttab request failed: %w", err)
	}
	defer resp.Body.Close()

	var items []MenuItem
	if err := json.NewDecoder(resp.Body).Decode(&items); err != nil {
		return nil, fmt.Errorf("toasttab decode failed: %w", err)
	}

	// Store directly — no character encoding normalization
	return items, nil
}
