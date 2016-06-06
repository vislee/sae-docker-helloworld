package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
)

type testHandler struct {
}

func (t *testHandler) errCode(w http.ResponseWriter, error string, code int) {
	w.Header().Set("Server", "lee-1.0 @SAE:wenqiang3")
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.Header().Set("X-Content-Type-Options", "nosniff")
	w.WriteHeader(code)
	fmt.Fprintln(w, error)
}

func (t *testHandler) envs(w http.ResponseWriter, req *http.Request) {
	envs := os.Environ()
	bd := "env:"
	for _, v := range envs {
		bd = fmt.Sprintf("%s\n%s", bd, v)
	}
	w.Header().Set("Server", "lee-1.0 @SAE:wenqiang3")
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.Header().Set("Content-Length", fmt.Sprintf("%d", len(bd)))
	io.WriteString(w, bd)
}

func (t *testHandler) hostName(w http.ResponseWriter, req *http.Request) {
	hn, err := os.Hostname()
	if err != nil {
		t.errCode(w, "500 internal error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Server", "lee-1.0 @SAE:wenqiang3")
	w.Header().Set("Content-Length", fmt.Sprintf("%d", len(hn)))
	io.WriteString(w, hn)
}

func (t *testHandler) pingHandler(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Server", "lee-1.0 @SAE:wenqiang3")
	w.Header().Set("Content-Length", "4")
	io.WriteString(w, "PANG")
}

func (t testHandler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	switch req.URL.Path {
	case "/ping":
		t.pingHandler(w, req)
	case "/hostname":
		t.hostName(w, req)
	case "/ttenvs":
		t.envs(w, req)
	default:
		t.errCode(w, "404 page not found", http.StatusNotFound)
	}
}

func main() {
	var t testHandler
	http.ListenAndServe(":5050", t)
}
