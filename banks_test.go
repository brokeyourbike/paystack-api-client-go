package paystack_test

import (
	"bytes"
	"context"
	_ "embed"
	"io"
	"net/http"
	"testing"

	"github.com/brokeyourbike/paystack-api-client-go"
	"github.com/stretchr/testify/assert"
	mock "github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

//go:embed testdata/banks-success.json
var banksSuccess []byte

func TestBanksSuccess(t *testing.T) {
	mockHttpClient := paystack.NewMockHttpClient(t)
	client := paystack.NewClient("token", paystack.WithHTTPClient(mockHttpClient))

	resp := &http.Response{StatusCode: 200, Body: io.NopCloser(bytes.NewReader(banksSuccess))}
	mockHttpClient.On("Do", mock.AnythingOfType("*http.Request")).Return(resp, nil).Once()

	got, err := client.GetBanks(context.TODO())
	require.NoError(t, err)

	assert.True(t, got.Status)
	assert.Len(t, got.Data, 5)
}
