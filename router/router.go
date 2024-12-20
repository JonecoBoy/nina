package router

import (
	"context"
	"crypto/tls"
	"encoding/json"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"net/url"
	"regexp"
)

type ServeMux struct {
	*http.ServeMux
}

type Handler func(w http.ResponseWriter, r *NinaRequest)

func (h Handler) ServeHTTP(writer http.ResponseWriter, request *NinaRequest) {
	// Call the original golang http handler
	h(writer, request)
}

type Middlewares []Middleware
type Middleware func(Handler) Handler

func NewRouter() *ServeMux {
	mux := &ServeMux{
		ServeMux: http.NewServeMux(),
	}
	return mux
}

type NinaRequest struct {
	*http.Request
	Header        http.Header
	ContentLength int64
	Form          *url.Values
	PostForm      *url.Values
	MultipartForm *multipart.Form
	RemoteAddr    string
	Method        string
	RequestURI    string
	ctx           context.Context
	tls           *tls.ConnectionState
	UserAgent     string
	Proto         string
	Host          string
	Pattern       map[string]string
	Params        *NinaParamsRequest
	Body          interface{}
}

type NinaParamsRequest struct {
	QueryString map[string]string
	UriParams   map[string]string
	Params      map[string]string
}

// variadic input
func (mux *ServeMux) GET(pattern string, handler Handler, middlewares ...Middleware) {
	finalHandler := applyMiddlewares(handler, middlewares...)
	mux.ServeMux.Handle("GET "+pattern, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}
		reqParams := getReqParams(r, pattern)
		params := &NinaParamsRequest{
			QueryString: reqParams["queryString"],
			UriParams:   reqParams["uriParams"],
			Params:      reqParams["params"],
		}

		ninaRequest := &NinaRequest{
			Request:       r,
			Header:        r.Header,
			Form:          &r.Form,
			Method:        r.Method,
			PostForm:      &r.PostForm,
			ctx:           r.Context(),
			ContentLength: r.ContentLength,
			tls:           r.TLS,
			Proto:         r.Proto,
			Host:          r.Host,
			Params:        params,
			UserAgent:     r.UserAgent(),
		}
		finalHandler(w, ninaRequest)
	}))
}

func (mux *ServeMux) POST(pattern string, handler Handler, middlewares ...Middleware) {
	finalHandler := applyMiddlewares(handler, middlewares...)
	mux.ServeMux.Handle("POST "+pattern, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		var parsedBody interface{}

		// Read and parse JSON body or form data
		if r.Header.Get("Content-Type") == "application/json" {
			bodyBytes, err := ioutil.ReadAll(r.Body)
			if err != nil {
				http.Error(w, "Unable to read body", http.StatusBadRequest)
				return
			}
			defer r.Body.Close()

			var jsonBody map[string]interface{}
			if err := json.Unmarshal(bodyBytes, &jsonBody); err != nil {
				http.Error(w, "Invalid JSON", http.StatusBadRequest)
				return
			}
			parsedBody = jsonBody
		} else {
			// Handle form data
			if err := r.ParseForm(); err != nil {
				http.Error(w, "Unable to parse form", http.StatusBadRequest)
				return
			}
			parsedBody = r.PostForm // Use parsed form data
		}

		reqParams := getReqParams(r, pattern)
		params := &NinaParamsRequest{
			QueryString: reqParams["queryString"],
			UriParams:   reqParams["uriParams"],
			Params:      reqParams["params"],
		}

		ninaRequest := &NinaRequest{
			Request:       r,
			Header:        r.Header,
			Form:          &r.Form,
			Method:        r.Method,
			PostForm:      &r.PostForm,
			ctx:           r.Context(),
			ContentLength: r.ContentLength,
			tls:           r.TLS,
			Proto:         r.Proto,
			Host:          r.Host,
			Params:        params,
			UserAgent:     r.UserAgent(),
			Body:          parsedBody, // Assign parsed body here
		}
		finalHandler(w, ninaRequest)
	}))
}

func getReqParams(r *http.Request, pattern string) map[string]map[string]string {
	qs := parseQueryString(r)
	uriParams := parseUriParams(r, pattern)

	params := make(map[string]map[string]string)
	params["params"] = make(map[string]string)
	params["queryString"] = qs
	params["uriParams"] = uriParams

	for key, value := range qs {
		params["params"][key] = value
	}
	for key, value := range uriParams {
		params["params"][key] = value
	}

	return params
}

func parseQueryString(r *http.Request) map[string]string {
	rawQS := r.URL.Query()
	qs := make(map[string]string)
	for key, values := range rawQS {
		if len(values) > 0 {
			qs[key] = values[0]
		}
	}
	return qs
}

func parseUriParams(r *http.Request, pattern string) map[string]string {
	re := regexp.MustCompile(`{([^}]*)}`)
	matches := re.FindAllStringSubmatch(pattern, -1)

	params := make(map[string]string)
	for _, match := range matches {
		params[match[1]] = r.PathValue(match[1])
	}

	return params
}

// variadic so can be any size of array
func applyMiddlewares(h Handler, middlewares ...Middleware) Handler {
	// in this normal order will be last middleware first!
	//for _, middleware := range middlewares {
	//	h = middleware(h)
	//}
	// in this order will load from the left to the right in the declaration
	for i := len(middlewares) - 1; i >= 0; i-- {
		h = middlewares[i](h)
	}
	return h
}
