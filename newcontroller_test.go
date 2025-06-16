package main

import (
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

type controllerTestCase struct {
	name            string
	channelPreload  [][]byte
	requestBody     io.Reader
	expectedStatus  int
	expectedBody    string
	expectedChannel string
}

func TestNewController(t *testing.T) {
	tests := []controllerTestCase{
		{
			name:            "case 1",
			requestBody:     strings.NewReader("1\n+2\n3\n"),
			expectedStatus:  http.StatusAccepted,
			expectedBody:    "OK: ",
			expectedChannel: "1\n+2\n3\n",
		},
		{
			name:            "case 2",
			requestBody:     badReader{},
			expectedStatus:  http.StatusBadRequest,
			expectedBody:    "Bad Input",
			expectedChannel: "",
		},
		{
			name:            "case 3",
			channelPreload:  [][]byte{[]byte("already full")},
			requestBody:     strings.NewReader("1\n+2\n3\n"),
			expectedStatus:  http.StatusServiceUnavailable,
			expectedBody:    "Too Busy: ",
			expectedChannel: "",
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			ch := make(chan []byte, 1)

			for _, msg := range tc.channelPreload {
				ch <- msg
			}

			handler := NewController(ch)
			req := httptest.NewRequest(http.MethodPost, "/", tc.requestBody)
			w := httptest.NewRecorder()
			handler.ServeHTTP(w, req)

			resp := w.Result()
			defer resp.Body.Close()

			if resp.StatusCode != tc.expectedStatus {
				t.Errorf("expected status %d, got %d", tc.expectedStatus, resp.StatusCode)
			}

			body, _ := io.ReadAll(resp.Body)
			if !strings.HasPrefix(string(body), tc.expectedBody) {
				t.Errorf("expected body: `%s`, got: `%s`", tc.expectedBody, string(body))
			}
			if tc.expectedChannel != "" {
				select {
				case got := <-ch:
					if string(got) != tc.expectedChannel {
						t.Errorf("expected `%q`, got: `%q`", tc.expectedChannel, string(got))
					}
				default:
					t.Errorf("no data in channel")

				}
			}

		})
	}
}

type badReader struct{}

func (badReader) Read(p []byte) (int, error) {
	return 0, io.ErrUnexpectedEOF
}
func (badReader) Close() error {
	return nil
}
