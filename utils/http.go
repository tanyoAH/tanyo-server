package utils

import (
	"encoding/json"
	"fmt"
	"github.com/tanyoAH/tanyo-server/config"
	"net/http"
	"reflect"
	"time"
)

// Response contains the attributes found in an API response
type APIResponse struct {
	Message string      `json:"message"`
	Success bool        `json:"success"`
	Data    interface{} `json:"data"`
	Debug   string      `json:"debug,omitempty"`
}

type CheckableRequest interface {
	Parameters() error
}

func JSONSuccess(w http.ResponseWriter, data interface{}, message string) {
	if message == "" {
		message = "ok"
	}
	resp := APIResponse{
		Message: message,
		Success: true,
		Data:    data,
	}
	jsonWriter(w, resp, http.StatusOK)
}

func JSONError(w http.ResponseWriter, data interface{}, message string, debug string, statusCode int) {
	if message == "" {
		message = "error"
	}
	resp := APIResponse{
		Message: message,
		Success: false,
		Data:    data,
		Debug:   debug,
	}
	jsonWriter(w, resp, statusCode)
}

func JSONInternalError(w http.ResponseWriter, message string, debug string) {
	if message == "" {
		message = "error"
	}
	resp := APIResponse{
		Message: message,
		Success: false,
		Data:    nil,
		Debug:   debug,
	}
	jsonWriter(w, resp, http.StatusInternalServerError)
}

func JSONBadRequestError(w http.ResponseWriter, message string, debug string) {
	if message == "" {
		message = "bad_request"
	}
	resp := APIResponse{
		Message: message,
		Success: false,
		Data:    nil,
		Debug:   debug,
	}
	jsonWriter(w, resp, http.StatusBadRequest)
}

func JSONNotFoundError(w http.ResponseWriter, message string, debug string) {
	if message == "" {
		message = "not_found"
	}
	resp := APIResponse{
		Message: message,
		Success: false,
		Data:    nil,
		Debug:   debug,
	}
	jsonWriter(w, resp, http.StatusNotFound)
}

func JSONForbiddenError(w http.ResponseWriter, message string, debug string) {
	if message == "" {
		message = "forbidden"
	}
	resp := APIResponse{
		Message: message,
		Success: false,
		Data:    nil,
		Debug:   debug,
	}
	jsonWriter(w, resp, http.StatusForbidden)
}

func JSONDetailed(w http.ResponseWriter, resp APIResponse, statusCode int) {
	if statusCode == 0 {
		statusCode = http.StatusOK
	}
	jsonWriter(w, resp, statusCode)
}

func jsonWriter(w http.ResponseWriter, d interface{}, c int) {
	//dj, err := json.MarshalIndent(d, "", "  ")
	dj, err := json.Marshal(d)
	if err != nil {
		http.Error(w, "Error creating JSON response", http.StatusInternalServerError)
		fmt.Println(err)
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(c)
	fmt.Fprintf(w, "%s", dj)
}

func JSONDecodeAndCatchForAPI(w http.ResponseWriter, r *http.Request, outStruct interface{}) error {
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&outStruct)
	if err != nil {
		JSONBadRequestError(w, "Invalid JSON", "")
		return err
	}
	if !isCheckableRequest(outStruct) {
		return nil
	}
	method := reflect.ValueOf(outStruct).MethodByName("Parameters").Interface().(func() error)
	err = method()
	if err != nil {
		JSONBadRequestError(w, "", err.Error())
		return err
	}
	return nil
}

func isCheckableRequest(checkAgainst interface{}) bool {
	reader := reflect.TypeOf((*CheckableRequest)(nil)).Elem()
	return reflect.TypeOf(checkAgainst).Implements(reader)
}

func GetVisitorIPv4(r *http.Request) string {
	if config.Conf.IsDebugMode() {
		return "54.218.27.6"
	} else {
		return r.Header.Get("X-Forwarded-For")
	}
}

func SetJavascriptAccessibleCookie(w http.ResponseWriter, cookieName string, cookieValue interface{}) {
	http.SetCookie(w, &http.Cookie{
		Name:    cookieName,
		Value:   fmt.Sprint(cookieValue),
		Path:    "/",
		Expires: time.Now().AddDate(10, 0, 0),
	})
}

func SetHTTPOnlyCookie(w http.ResponseWriter, cookieName string, cookieValue interface{}) {
	http.SetCookie(w, &http.Cookie{
		Name:     cookieName,
		Value:    fmt.Sprint(cookieValue),
		Path:     "/",
		HttpOnly: true,
		Expires:  time.Now().AddDate(10, 0, 0),
	})
}

func DeleteCookie(w http.ResponseWriter, cookieName string) {
	http.SetCookie(w, &http.Cookie{
		Name:   cookieName,
		MaxAge: -1,
		Path:   "/",
	})
}

func TemporaryRedirect(w http.ResponseWriter, r *http.Request, url string) {
	http.Redirect(w, r, url, http.StatusTemporaryRedirect)
}

func InternalServerErrorView(w http.ResponseWriter, r *http.Request) {
	htmlWriter(w, r, "<div>Internal server error</div>", http.StatusInternalServerError)
}

func UnauthorizedErrorView(w http.ResponseWriter, r *http.Request) {
	htmlWriter(w, r, "<div>Unauthorized request error</div>", http.StatusUnauthorized)
}

func htmlWriter(w http.ResponseWriter, r *http.Request, html string, status int) {
	w.Header().Set("Content-Type", "text/html")
	w.WriteHeader(status)
	fmt.Fprintf(w, "%s", html)
}
