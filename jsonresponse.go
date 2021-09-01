package main
import (
	// external
	"encoding/json"
	"net/http"
	"github.com/golang/gddo/httputil/header"
	"fmt"
	"io"
	"errors"
	"strconv"
	"strings"
)

// 	errorResponse(w, "Content Type is not application/json", http.StatusUnsupportedMediaType)
func errorResponse(w http.ResponseWriter, message string, httpStatusCode int) {
    w.Header().Set("Content-Type", "application/json")
    // w.Header().Set("Cache-Control", "['no-cache','no-store','must-revalidate']")
    // w.Header().Set("Pragma", "no-cache")

    w.WriteHeader(httpStatusCode)

    scaErrorResp  := new (ScanErrorResponse)
    sanErrorData := ScanErrorData{message, strconv.Itoa(httpStatusCode)}
    scaErrorResp.Error = sanErrorData

    // scaErrorResp := make(map[string]map[string]string)
    // sanErrorData := map[string]string{"message": message, "code": strconv.Itoa(httpStatusCode)}
    // scaErrorResp["error"] = sanErrorData
    jsonResp, _ := json.Marshal(scaErrorResp)
    w.Write(jsonResp)
}

func sendJsonResponse(w http.ResponseWriter, jsonResp []byte) {
    w.Header().Set("Content-Type", "application/json")
    // w.Header().Set("Cache-Control", "['no-cache','no-store','must-revalidate']")
    // w.Header().Set("Pragma", "no-cache")
    w.WriteHeader(200)
    w.Write(jsonResp)
}

type malformedRequest struct {
    status int
    msg    string
}

func (mr *malformedRequest) Error() string {
    return mr.msg
}

func validateContentType(w http.ResponseWriter, r *http.Request) error {

    if r.Header.Get("Content-Type") != "" {
        value, _ := header.ParseValueAndParams(r.Header, "Content-Type")
        if value != "application/json" {
            return &malformedRequest{status: http.StatusUnsupportedMediaType, msg: "Content-Type header is not application/json"}
        }
    } 
    return nil
}


func decodeJSONBody(w http.ResponseWriter, r *http.Request, dst interface{}) error {

    if r.Header.Get("Content-Type") != "" {
        value, _ := header.ParseValueAndParams(r.Header, "Content-Type")
        if value != "application/json" {
            return &malformedRequest{status: http.StatusUnsupportedMediaType, msg: "Content-Type header is not application/json"}
        }
    }

    r.Body = http.MaxBytesReader(w, r.Body, 1048576)

    dec := json.NewDecoder(r.Body)
    dec.DisallowUnknownFields()

    err := dec.Decode(&dst)
    if err != nil {
        var syntaxError *json.SyntaxError
        var unmarshalTypeError *json.UnmarshalTypeError

        switch {
        case errors.As(err, &syntaxError):
            msg := "Request body contains badly-formed JSON (at position "+ fmt.Sprintf("%d", syntaxError.Offset)+")"
            return &malformedRequest{status: http.StatusBadRequest, msg: msg}

        case errors.Is(err, io.ErrUnexpectedEOF):
            msg := "Request body contains badly-formed JSON"
            return &malformedRequest{status: http.StatusBadRequest, msg: msg}

        case errors.As(err, &unmarshalTypeError):
            msg := "Request body contains an invalid value for the "+unmarshalTypeError.Field+" field (at position "+ fmt.Sprintf("%d", syntaxError.Offset)+")"
            return &malformedRequest{status: http.StatusBadRequest, msg: msg}

        case strings.HasPrefix(err.Error(), "json: unknown field "):
            fieldName := strings.TrimPrefix(err.Error(), "json: unknown field ")
            msg := "Request body contains unknown field "+fieldName
            return &malformedRequest{status: http.StatusBadRequest, msg: msg}

        case errors.Is(err, io.EOF):
            msg := "Request body must not be empty"
            return &malformedRequest{status: http.StatusBadRequest, msg: msg}

        case err.Error() == "http: request body too large":
            msg := "Request body must not be larger than 1MB"
            return &malformedRequest{status: http.StatusBadRequest, msg: msg}
        default:
            return err
        }
    }

	err = dec.Decode(&struct{}{})
	if err != io.EOF {
        msg := "Request body must only contain a single JSON object"
        return &malformedRequest{status: http.StatusBadRequest, msg: msg}
    }

    return nil
}
