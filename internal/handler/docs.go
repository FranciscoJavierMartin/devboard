// Package handler health
package handler

import (
	"fmt"
	"net/http"

	"github.com/MarceloPetrucio/go-scalar-api-reference"
)

// DocsHandler struct
type DocsHandler struct {
	html []byte
}

// NewDocsHandler builds the Scalar API reference page from the given OpenAPI spec
func NewDocsHandler(spec []byte) (*DocsHandler, error) {
	htmlContent, err := scalar.ApiReferenceHTML(&scalar.Options{
		SpecContent: string(spec),
		CustomOptions: scalar.CustomOptions{
			PageTitle: "Devboard API",
		},
		DarkMode: true,
	})
	if err != nil {
		return nil, fmt.Errorf("building docs html: %w", err)
	}

	return &DocsHandler{html: []byte(htmlContent)}, nil
}

// ServeHTTP serves the Scalar API reference page
func (handle *DocsHandler) ServeHTTP(writer http.ResponseWriter, _ *http.Request) {
	writer.Header().Set("Content-Type", "text/html; charset=utf-8")
	_, _ = writer.Write(handle.html)
}
