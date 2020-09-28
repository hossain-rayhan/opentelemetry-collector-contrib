// Copyright 2020, OpenTelemetry Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package awsecscontainermetrics

import (
	"io"
	"net/http"
	"strings"
	"testing"

	"github.com/pkg/errors"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
)

func TestClient(t *testing.T) {
	tr := &fakeRoundTripper{}
	baseURL := "http://localhost:8080"
	client := &clientImpl{
		baseURL:    baseURL,
		httpClient: http.Client{Transport: tr},
	}
	require.False(t, tr.closed)
	resp, err := client.Get("/stats")
	require.NoError(t, err)
	require.Equal(t, "hello", string(resp))
	require.True(t, tr.closed)
	require.Equal(t, baseURL+"/stats", tr.url)
	require.Equal(t, 1, len(tr.header))
	require.Equal(t, "application/json", tr.header["Content-Type"][0])
	require.Equal(t, "GET", tr.method)
}

func TestNewClientProvider(t *testing.T) {
	provider := NewClientProvider("http://localhost:8080", zap.NewNop())
	require.NotNil(t, provider)
	_, ok := provider.(*defaultClientProvider)
	require.True(t, ok)

	client := provider.BuildClient()
	require.Equal(t, "http://localhost:8080", string(client.(*clientImpl).baseURL))
}

func TestDefaultClient(t *testing.T) {
	endpoint := "http://localhost:8080"
	client := defaultClient(endpoint, zap.NewNop())
	require.NotNil(t, client.httpClient.Transport)
	require.Equal(t, endpoint, client.baseURL)
}

func TestBuildReq(t *testing.T) {
	p := &defaultClientProvider{
		endpoint: "http://localhost:8080",
		logger:   zap.NewNop(),
	}
	cl := p.BuildClient()
	req, err := cl.(*clientImpl).buildReq("/test")
	require.NoError(t, err)
	require.NotNil(t, req)
	require.Equal(t, "application/json", req.Header["Content-Type"][0])
}

func TestBuildBadReq(t *testing.T) {
	p := &defaultClientProvider{
		endpoint: "http://localhost:8080",
		logger:   zap.NewNop(),
	}
	cl := p.BuildClient()
	_, err := cl.(*clientImpl).buildReq(" ")
	require.Error(t, err)
}

func TestGetBad(t *testing.T) {
	p := &defaultClientProvider{
		endpoint: "http://localhost:8080",
		logger:   zap.NewNop(),
	}
	cl := p.BuildClient()
	resp, err := cl.(*clientImpl).Get(" ")
	require.Error(t, err)
	require.Nil(t, resp)
}

func TestFailedRT(t *testing.T) {
	tr := &fakeRoundTripper{failOnRT: true}
	baseURL := "http://localhost:8080"
	client := &clientImpl{
		baseURL:    baseURL,
		httpClient: http.Client{Transport: tr},
	}
	_, err := client.Get("/test")
	require.Error(t, err)
}

func TestErrOnRead(t *testing.T) {
	tr := &fakeRoundTripper{errOnRead: true}
	baseURL := "http://localhost:8080"
	client := &clientImpl{
		baseURL:    baseURL,
		httpClient: http.Client{Transport: tr},
		logger:     zap.NewNop(),
	}
	resp, err := client.Get("/foo")
	require.Error(t, err)
	require.Nil(t, resp)
}

func TestErrCode(t *testing.T) {
	tr := &fakeRoundTripper{errCode: true}
	baseURL := "http://localhost:9876"
	client := &clientImpl{
		baseURL:    baseURL,
		httpClient: http.Client{Transport: tr},
		logger:     zap.NewNop(),
	}
	resp, err := client.Get("/foo")
	require.Error(t, err)
	require.Nil(t, resp)
}

var _ http.RoundTripper = (*fakeRoundTripper)(nil)

type fakeRoundTripper struct {
	closed     bool
	header     http.Header
	method     string
	url        string
	failOnRT   bool
	errOnClose bool
	errOnRead  bool
	errCode    bool
}

func (f *fakeRoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	if f.failOnRT {
		return nil, errors.New("failOnRT == true")
	}
	f.header = req.Header
	f.method = req.Method
	f.url = req.URL.String()
	var reader io.Reader
	if f.errOnRead {
		reader = &failingReader{}
	} else {
		reader = strings.NewReader("hello")
	}
	statusCode := 200
	if f.errCode {
		statusCode = 503
	}
	return &http.Response{
		StatusCode: statusCode,
		Body: &fakeReadCloser{
			Reader: reader,
			onClose: func() error {
				f.closed = true
				if f.errOnClose {
					return errors.New("")
				}
				return nil
			},
		},
	}, nil
}

var _ io.Reader = (*failingReader)(nil)

type failingReader struct{}

func (f *failingReader) Read([]byte) (n int, err error) {
	return 0, errors.New("")
}

var _ io.ReadCloser = (*fakeReadCloser)(nil)

type fakeReadCloser struct {
	io.Reader
	onClose func() error
}

func (f *fakeReadCloser) Close() error {
	return f.onClose()
}
