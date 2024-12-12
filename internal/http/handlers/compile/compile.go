package compile

import (
	"encoding/json"
	"net/http"

	codeTypes "github.com/theMitocondria/compiler/internal/types/code"
)

func HelloWorld() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello World."))
	}
}

func CompileCode() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req codeTypes.CompileCodeRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "Failed to parse request body", http.StatusBadRequest)
			return
		}

		var result codeTypes.CompiledCodeResponse

		response, err := 

	}
}
