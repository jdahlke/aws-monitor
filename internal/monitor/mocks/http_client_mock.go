package mocks

import "net/http"

type MockClient struct {
	DoFunc        func(req *http.Request) (*http.Response, error)
	DoFuncInvoked bool
}

func (m *MockClient) Do(req *http.Request) (*http.Response, error) {
	m.DoFuncInvoked = true

	return m.DoFunc(req)
}
