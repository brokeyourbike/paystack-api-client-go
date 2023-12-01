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

//go:embed testdata/resolve-success.json
var resolveSuccess []byte

//go:embed testdata/resolve-fail-auth.json
var resolveFailAuth []byte

func TestResolveSuccess(t *testing.T) {
	mockHttpClient := paystack.NewMockHttpClient(t)
	client := paystack.NewClient("token", paystack.WithHTTPClient(mockHttpClient))

	resp := &http.Response{StatusCode: 200, Body: io.NopCloser(bytes.NewReader(resolveSuccess))}
	mockHttpClient.On("Do", mock.AnythingOfType("*http.Request")).Return(resp, nil).Once()

	got, err := client.ResolveBankAccount(context.TODO(), "123", "123456789")
	require.NoError(t, err)

	assert.True(t, got.Status)
}

func TestResolveFailedAuth(t *testing.T) {
	mockHttpClient := paystack.NewMockHttpClient(t)
	client := paystack.NewClient("token", paystack.WithHTTPClient(mockHttpClient))

	resp := &http.Response{StatusCode: 401, Body: io.NopCloser(bytes.NewReader(resolveFailAuth))}
	mockHttpClient.On("Do", mock.AnythingOfType("*http.Request")).Return(resp, nil).Once()

	got, err := client.ResolveBankAccount(context.TODO(), "123", "123456789")
	require.NoError(t, err)

	assert.False(t, got.Status)
}
