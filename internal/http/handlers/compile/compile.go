// package compile

// import (
// 	"bytes"
// 	"encoding/json"
// 	"fmt"
// 	"io"
// 	"net/http"

// 	codeTypes "github.com/theMitocondria/compiler/internal/types/code"
// )

// func HelloWorld() http.HandlerFunc {
// 	return func(w http.ResponseWriter, r *http.Request) {
// 		w.Write([]byte("Hello World."))
// 	}
// }

// func CompileCode() http.HandlerFunc {
// 	return func(w http.ResponseWriter, r *http.Request) {
// 		var req codeTypes.CodeExecutionRequest
// 		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
// 			json.NewEncoder(w).Encode(map[string]string{"output": "", "error": "Failed to parse the body"})
// 		}

// 		payload, err := json.Marshal(req)

// 		if err != nil {
// 			json.NewEncoder(w).Encode(map[string]string{"output": "", "error": "Failed to marshal request payload"})
// 		}

// 		resp, err := http.Post("http://localhost:8089/compile", "application/json", bytes.NewBuffer(payload))
// 		if err != nil {
// 			json.NewEncoder(w).Encode(map[string]string{"output": "", "error": "Failed to send request to container"})
// 			return
// 		}

// 		defer resp.Body.Close()
// 		respBody, err := io.ReadAll(resp.Body)
// 		if err != nil {
// 			json.NewEncoder(w).Encode(map[string]string{"output": "", "error": "Failed to read response from container"})
// 			return
// 		}
// 		var result codeTypes.CodeExecutionResponse
// 		if err := json.Unmarshal(respBody, &result); err != nil {
// 			json.NewEncoder(w).Encode(map[string]string{"output": "", "error": "Failed to unmarshal response"})
// 			return
// 		}
// 		w.Header().Set("Content-Type", "application/json")
// 		json.NewEncoder(w).Encode(result)
// 	}
// }

package compile

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"sync"

	codeTypes "github.com/theMitocondria/compiler/internal/types/code"
)

func HelloWorld() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello World."))
	}
}

func CompileCode(loadBalancer *codeTypes.LoadBalancer, mu *sync.Mutex) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req codeTypes.CodeExecutionRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			json.NewEncoder(w).Encode(map[string]string{"output": "", "error": "Failed to parse the body"})
			return
		}

		mu.Lock()
		defer mu.Unlock()

		// Find an available container
		var container *codeTypes.Container
		for i := range loadBalancer.Containers {
			if loadBalancer.Containers[i].Lang == req.Lang && !loadBalancer.Containers[i].InUse {
				container = &loadBalancer.Containers[i]
				break
			}
		}

		if container == nil {
			json.NewEncoder(w).Encode(map[string]string{"output": "", "error": "No available container"})
			return
		}

		// Mark the container as in use
		container.InUse = true

		// Execute the code
		go func() {
			defer func() {
				mu.Lock()
				container.InUse = false
				mu.Unlock()
			}()

			payload, err := json.Marshal(req)
			if err != nil {
				mu.Lock()
				json.NewEncoder(w).Encode(map[string]string{"output": "", "error": "Failed to marshal request payload"})
				mu.Unlock()
				return
			}

			fmt.Println(fmt.Sprintf("http://localhost:%s/compile", container.Port))
			resp, err := http.Post(fmt.Sprintf("http://localhost:%s/compile", container.Port), "application/json", bytes.NewBuffer(payload))
			if err != nil {
				mu.Lock()
				json.NewEncoder(w).Encode(map[string]string{"output": "", "error": "Failed to send request to container"})
				mu.Unlock()
				return
			}

			defer resp.Body.Close()
			respBody, err := io.ReadAll(resp.Body)
			if err != nil {
				mu.Lock()
				json.NewEncoder(w).Encode(map[string]string{"output": "", "error": "Failed to read response from container"})
				mu.Unlock()
				return
			}

			var result codeTypes.CodeExecutionResponse
			if err := json.Unmarshal(respBody, &result); err != nil {
				mu.Lock()
				json.NewEncoder(w).Encode(map[string]string{"output": "", "error": "Failed to unmarshal response"})
				mu.Unlock()
				return
			}

			mu.Lock()
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(result)
			mu.Unlock()
		}()
	}
}
