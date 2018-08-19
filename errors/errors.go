package errors

import (
    "encoding/json"
    "net/http"
)

type ApiError struct {
    Code int `json:"code"`
    Key string `json:"key"`
    Message string `json:"message"`
}

func NewApiError(code int, key string, message string) ApiError {
    return ApiError{code, key, message}
}

func Empty() ApiError {
    return ApiError{}
}

func NotFound() ApiError {
    return NewApiError(404, "not-found", "Endpoint not found.")
}

func InvalidJson() ApiError {
    return NewApiError(400, "invalid-json", "Received input is not valid JSON.")
}

func ValidationError(message string) ApiError {
    return NewApiError(400, "validation-error", message)
}

func (err ApiError) Error() string {
    return string(err.Json())
}

func (apiErr ApiError) Json() []byte {
    dict := make(map[string]ApiError)
    dict["error"] = apiErr
    json, err := json.Marshal(dict)
    if err != nil {
        panic(err)
    }
    return json
}

func (err ApiError) NotEmpty() bool {
    return err.Code != 0
}

func (err ApiError) WriteToResponse(response http.ResponseWriter) {
    response.WriteHeader(err.Code)
    response.Write(err.Json())
}
