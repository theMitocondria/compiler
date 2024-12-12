package compile

import (
	"bytes"
	"encoding/json"
	"io"
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
			json.NewEncoder(w).Encode(map[string]string{"output": "", "error": "Failed to parse the body"})
		}

		payload, err := json.Marshal(req)

		if err != nil {
			json.NewEncoder(w).Encode(map[string]string{"output": "", "error": "Failed to marshal request payload"})
		}

		resp, err := http.Post("http://localhost:8081/compile", "application/json", bytes.NewBuffer(payload))
		if err != nil {
			json.NewEncoder(w).Encode(map[string]string{"output": "", "error": "Failed to send request to container"})
			return
		}

		defer resp.Body.Close()
		respBody, err := io.ReadAll(resp.Body)
		if err != nil {
			json.NewEncoder(w).Encode(map[string]string{"output": "", "error": "Failed to read response from container"})
			return
		}
		var result codeTypes.CompiledCodeResponse
		if err := json.Unmarshal(respBody, &result); err != nil {
			json.NewEncoder(w).Encode(map[string]string{"output": "", "error": "Failed to unmarshal response"})
			return
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(result)
	}
}
