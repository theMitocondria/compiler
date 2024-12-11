package compile

import (
	"net/http"
)

func HelloWorld() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello World."))
	}
}

func CompileCode() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("compiled successfully."))
	}
}
