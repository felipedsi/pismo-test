package handler

import (
	"log"
	"net/http"

	"github.com/go-chi/render"
)

type ErrorResponse struct {
	Err            error `json:"-"` // low-level runtime error
	HTTPStatusCode int   `json:"-"` // http response status code

	StatusText string `json:"status"`          // user-level status message
	ErrorText  string `json:"error,omitempty"` // application-level error message
}

func (e *ErrorResponse) Render(w http.ResponseWriter, r *http.Request) error {
	render.Status(r, e.HTTPStatusCode)
	return nil
}

func errorInvalidRequest(err error, errorText string) render.Renderer {
	log.Printf("Invalid request error: %s, %s", err, errorText)

	return &ErrorResponse{
		Err:            err,
		HTTPStatusCode: 400,
		StatusText:     "Invalid request",
		ErrorText:      errorText,
	}
}

func errorNotFound(err error, errorText string) render.Renderer {
	log.Printf("Not found error: %s, %s", err, errorText)

	return &ErrorResponse{
		Err:            err,
		HTTPStatusCode: 404,
		StatusText:     "Not found",
		ErrorText:      errorText,
	}
}
