package paystack

import (
	"context"
	"fmt"
	"net/http"
)

type ResolveResponse struct {
	Status  bool   `json:"status"`
	Message string `json:"message"`
	Data    struct {
		AccountNumber string `json:"account_number"`
		AccountName   string `json:"account_name"`
		BankID        int    `json:"bank_id"`
	} `json:"data"`
}

func (c *client) Resolve(ctx context.Context, bankCode, accountNumber string) (data ResolveResponse, err error) {
	req, err := c.newRequest(ctx, http.MethodGet, "/bank/resolve")
	if err != nil {
		return data, fmt.Errorf("failed to create request: %w", err)
	}

	req.AddQueryParam("bank_code", bankCode)
	req.AddQueryParam("account_number", accountNumber)
	req.DecodeTo(&data)

	return data, c.do(ctx, req)
}
