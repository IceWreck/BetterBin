package handlers

import (
	"fmt"
	"net/http"

	"github.com/IceWreck/BetterBin/config"
)

// The logError() method is a generic helper for logging an error message.
func logError(app *config.Application, r *http.Request, err error) {
	app.Logger.Error().Stack().Err(err).Str("request_method", r.Method).Str("request_url", r.URL.String()).Msg("")
}

// The errorResponse() method is a generic helper for sending JSON-formatted error
// messages to the client with a given status code. Note that we're using an interface{}
// type for the message parameter, rather than just a string type, as this gives us
// more flexibility over the values that we can include in the response.
func errorResponse(app *config.Application, w http.ResponseWriter, r *http.Request, status int, message interface{}) {
	data := map[string]interface{}{"error": message}
	// Write the response using the writeJSON() helper. If this happens to return an
	// error then log it, and fall back to sending the client an empty response with a
	// 500 Internal Server Error status code.
	err := writeJSON(app, w, status, data, nil)
	if err != nil {
		logError(app, r, err)
		w.WriteHeader(500)
	}
}

// The serverErrorResponse() method will be used when our application encounters an
// unexpected problem at runtime. It logs the detailed error message, then uses the
// errorResponse() helper to send a 500 Internal Server Error status code and JSON
// response (containing a generic error message) to the client.
func serverErrorResponse(app *config.Application, w http.ResponseWriter, r *http.Request, err error) {
	logError(app, r, err)
	message := "the server encountered a problem and could not process your request"
	errorResponse(app, w, r, http.StatusInternalServerError, message)
}

// The notFoundResponse() method will be used to send a 404 Not Found status code and
// JSON response to the client.
func notFoundResponse(app *config.Application) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		message := "the requested resource could not be found"
		errorResponse(app, w, r, http.StatusNotFound, message)
	}
}

// The methodNotAllowedResponse() method will be used to send a 405 Method Not Allowed
// status code and JSON response to the client.
func methodNotAllowedResponse(app *config.Application) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		message := fmt.Sprintf("the %s method is not supported for this resource", r.Method)
		errorResponse(app, w, r, http.StatusMethodNotAllowed, message)
	}
}

// The badRequestResponse() method will be used to send a 400 Bad request status code
// and JSON response to the client.
func badRequestResponse(app *config.Application, w http.ResponseWriter, r *http.Request, err error) {
	errorResponse(app, w, r, http.StatusBadRequest, err.Error())
}

// Note that the errors parameter here has the type map[string]string, which is exactly
// the same as the errors map contained in our Validator type.
func failedValidationResponse(app *config.Application, w http.ResponseWriter, r *http.Request, errors map[string]string) {
	errorResponse(app, w, r, http.StatusUnprocessableEntity, errors)
}
