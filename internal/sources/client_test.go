package sources

import "net/http"

type MockDoType func(req *http.Request) (*http.Response, error)

type MockClient struct {
	DoFunc MockDoType
}

func (m *MockClient) Do(req *http.Request) (*http.Response, error) {
	return m.DoFunc(req)
}
