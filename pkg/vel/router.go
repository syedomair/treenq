package vel

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"unsafe"

	"github.com/gorilla/schema"
)

type Handler[I, O any] func(ctx context.Context, i I) (O, *Error)

func NewHandler[I, O any](call Handler[I, O]) http.HandlerFunc {
	var iType I
	var oType O
	hasReqBody := unsafe.Sizeof(iType) != 0
	hasResBody := unsafe.Sizeof(oType) != 0

	decoder := schema.NewDecoder()

	return func(w http.ResponseWriter, r *http.Request) {
		// TODO: must be validated against memory leak, write a unit test to bench mark against 100k calls
		*r = *r.WithContext(RequestWithContext(r.Context(), r))
		*r = *r.WithContext(WriterWithContext(r.Context(), w))
		var i I

		if hasReqBody {
			if r.Method == "GET" {
				if err := decoder.Decode(&i, r.URL.Query()); err != nil {
					w.WriteHeader(http.StatusBadRequest)
					err = json.NewEncoder(w).Encode(Error{
						Code: "FAILED_DECODING_QUERY",
						Err:  err,
					})
					if err != nil {
						slog.Default().ErrorContext(r.Context(), "failed to write request marshal error", "err", err)
					}
					return
				}
			} else {
				if err := json.NewDecoder(r.Body).Decode(&i); err != nil {
					w.WriteHeader(http.StatusBadRequest)
					err = json.NewEncoder(w).Encode(Error{
						Code: "FAILED_DECODING_BODY",
						Err:  err,
					})
					if err != nil {
						slog.Default().ErrorContext(r.Context(), "failed to write request marshal error", "err", err)
					}
					return
				}
			}
		}

		res, callErr := call(r.Context(), i)
		if callErr != nil {
			*r = *r.WithContext(LogAttrsToContext(r.Context(), "err", callErr.Err, "code", callErr.Code, "errMsg", callErr.Message))
			status := http.StatusBadRequest
			if callErr.Code == "" {
				status = http.StatusInternalServerError
			}
			w.WriteHeader(status)
			err := json.NewEncoder(w).Encode(callErr)
			if err != nil {
				slog.Default().ErrorContext(r.Context(), "failed to write api call error", "err", err)
			}
			return
		}

		if hasResBody {
			if err := json.NewEncoder(w).Encode(res); err != nil {
				w.WriteHeader(http.StatusBadRequest)
				err = json.NewEncoder(w).Encode(Error{
					Code:    "FAILED_ENCODING",
					Message: err.Error(),
				})
				if err != nil {
					slog.Default().ErrorContext(r.Context(), "failed to write request marshal error", "err", err)
				}
			}
		}
	}
}

type Router struct {
	mux         *http.ServeMux
	middlewares []Middleware

	handlersMeta []HandlerMeta
}

func (r *Router) Mux() *http.ServeMux {
	return r.mux
}

func (r *Router) Use(m func(http.Handler) http.Handler) {
	r.middlewares = append(r.middlewares, m)
}

func (r *Router) Meta() []HandlerMeta {
	meta := make([]HandlerMeta, len(r.handlersMeta))
	copy(meta, r.handlersMeta)
	return meta
}

type HandlerMeta struct {
	Input       any
	Output      any
	OperationID string
	Method      string
}

type Error struct {
	Code    string            `json:"code"`
	Message string            `json:"message"`
	Meta    map[string]string `json:"meta"`
	Err     error             `json:"-"`
}

func (e *Error) Error() string {
	if e.Message == "" {
		return e.Code
	}
	return fmt.Sprintf("%s: %s", e.Code, e.Message)
}

func (e *Error) JsonString() string {
	// we assume it never gives an error,
	// either way we catch it in the other tests
	data, _ := json.Marshal(e)
	return string(data)
}

func NewRouter() *Router {
	mux := http.NewServeMux()
	mux.Handle("GET /healthz", NewHandler(func(ctx context.Context, _ struct{}) (struct{}, *Error) {
		return struct{}{}, nil
	}))

	return &Router{mux: mux}
}

type (
	Middleware func(http.Handler) http.Handler
)

func NoopMiddleware(h http.Handler) http.Handler {
	return h
}

func RegisterPost[I, O any](r *Router, operationID string, handler Handler[I, O], middlewares ...Middleware) {
	var i I
	var o O
	r.handlersMeta = append(r.handlersMeta, HandlerMeta{
		Input:       i,
		Output:      o,
		OperationID: operationID,
		Method:      "POST",
	})

	var h http.Handler = NewHandler(handler)

	RegisterHandler(r, "POST /"+operationID, h, middlewares...)
}

func RegisterGet[I, O any](r *Router, operationID string, handler Handler[I, O], middlewares ...Middleware) {
	var i I
	var o O
	r.handlersMeta = append(r.handlersMeta, HandlerMeta{
		Input:       i,
		Output:      o,
		OperationID: operationID,
		Method:      "GET",
	})

	var h http.Handler = NewHandler(handler)
	RegisterHandler(r, "GET /"+operationID, h, middlewares...)
}

func RegisterHandler(r *Router, pattern string, handler http.Handler, middlewares ...Middleware) {
	for i := range middlewares {
		handler = middlewares[i](handler)
	}
	for i := range r.middlewares {
		handler = r.middlewares[i](handler)
	}
	r.mux.Handle(pattern, handler)
}

func RegisterHandlerFunc(r *Router, pattern string, h http.HandlerFunc, middlewares ...Middleware) {
	var handler http.Handler = h
	for i := range middlewares {
		handler = middlewares[i](handler)
	}
	for i := range r.middlewares {
		handler = r.middlewares[i](handler)
	}
	r.mux.Handle(pattern, handler)
}
