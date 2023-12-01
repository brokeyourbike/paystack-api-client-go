package paystack_test

import (
	"bytes"
	"context"
	_ "embed"
	"io"
	"net/http"
	"testing"

	"github.com/brokeyourbike/paystack-api-client-go"
	"github.com/sirupsen/logrus"
	logrustest "github.com/sirupsen/logrus/hooks/test"
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

	logger, hook := logrustest.NewNullLogger()
	logger.SetLevel(logrus.DebugLevel)

	client := paystack.NewClient("token", paystack.WithHTTPClient(mockHttpClient), paystack.WithBaseURL("https://a.com"), paystack.WithLogger(logger))

	resp := &http.Response{StatusCode: 200, Body: io.NopCloser(bytes.NewReader(resolveSuccess))}
	mockHttpClient.On("Do", mock.MatchedBy(func(req *http.Request) bool { return req.Host == "a.com" })).Return(resp, nil).Once()

	got, err := client.ResolveBankAccount(context.TODO(), "123", "123456789")
	require.NoError(t, err)

	assert.True(t, got.Status)
	assert.Equal(t, 2, len(hook.Entries))
}

func TestResolveFailedAuth(t *testing.T) {
	mockHttpClient := paystack.NewMockHttpClient(t)
	client := paystack.NewClient("token", paystack.WithHTTPClient(mockHttpClient))

	resp := &http.Response{StatusCode: 401, Body: io.NopCloser(bytes.NewReader(resolveFailAuth))}
	mockHttpClient.On("Do", mock.MatchedBy(func(req *http.Request) bool { return req.Host == "api.paystack.co" })).Return(resp, nil).Once()

	got, err := client.ResolveBankAccount(context.TODO(), "123", "123456789")
	require.NoError(t, err)

	assert.False(t, got.Status)
}
