package paystack

import (
	"context"
	"fmt"
	"net/http"
	"time"
)

type BanksResponse struct {
	Status  bool   `json:"status"`
	Message string `json:"message"`
	Data    []Bank `json:"data"`
	Meta    Meta   `json:"meta"`
}

type Bank struct {
	Name        string    `json:"name"`
	Slug        string    `json:"slug"`
	Code        string    `json:"code"`
	LongCode    string    `json:"longcode"`
	Gateway     string    `json:"gateway"`
	PayWithBank bool      `json:"pay_with_bank"`
	Active      bool      `json:"active"`
	IsDeleted   bool      `json:"is_deleted"`
	Country     string    `json:"country"`
	Currency    string    `json:"currency"`
	Type        string    `json:"type"`
	ID          int       `json:"id"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
}

type Meta struct {
	Next     string `json:"next"`
	Previous string `json:"previous"`
	PerPage  int    `json:"perPage"`
}

func (c *client) GetBanks(ctx context.Context) (data BanksResponse, err error) {
	req, err := c.newRequest(ctx, http.MethodGet, "/bank")
	if err != nil {
		return data, fmt.Errorf("failed to create request: %w", err)
	}

	req.DecodeTo(&data)

	return data, c.do(ctx, req)
}
