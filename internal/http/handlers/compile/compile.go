package compile

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"sync"
	"time"

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
		t := time.Now()
		slog.Info("Current second", "second", t.Nanosecond())

		mu.Lock()

		// Find an available container
		var container *codeTypes.Container
		for i := range loadBalancer.Containers {
			if loadBalancer.Containers[i].Lang == req.Lang && !loadBalancer.Containers[i].InUse {
				container = &loadBalancer.Containers[i]
				break
			}
		}

		if container == nil {
			mu.Unlock()
			json.NewEncoder(w).Encode(map[string]string{"output": "", "error": "No available container"})

			return
		}

		// Mark the container as in use
		container.InUse = true
		mu.Unlock()

		// Channel to receive the result
		resultChan := make(chan codeTypes.CodeExecutionResponse)
		errorChan := make(chan error)

		go func() {
			defer func() {
				mu.Lock()
				container.InUse = false
				mu.Unlock()
			}()

			payload, err := json.Marshal(req)
			if err != nil {
				errorChan <- fmt.Errorf("Failed to marshal request payload")
				return
			}

			fmt.Println(fmt.Sprintf("http://localhost:%s/compile", container.Port))
			resp, err := http.Post(fmt.Sprintf("http://localhost:%s/compile", container.Port), "application/json", bytes.NewBuffer(payload))
			if err != nil {
				errorChan <- fmt.Errorf("Failed to send request to container")
				return
			}

			defer resp.Body.Close()
			respBody, err := io.ReadAll(resp.Body)
			if err != nil {
				errorChan <- fmt.Errorf("Failed to read response from container")
				return
			}

			var result codeTypes.CodeExecutionResponse
			if err := json.Unmarshal(respBody, &result); err != nil {
				errorChan <- fmt.Errorf("Failed to unmarshal response")
				return
			}

			resultChan <- result
		}()

		select {
		case result := <-resultChan:
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(result)
		case err := <-errorChan:
			json.NewEncoder(w).Encode(map[string]string{"output": "", "error": err.Error()})
		}
	}
}