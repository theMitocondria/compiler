package compile

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"sync"

	codeTypes "github.com/theMitocondria/compiler/internal/types/code"
	"github.com/theMitocondria/compiler/internal/utils/dataStructs"
)

var cppPorts = dataStructs.InitializeQueue([]int{8081, 8082, 8083, 8084, 8085})
var jsPorts = dataStructs.InitializeQueue([]int{8086, 8087, 8088, 8089, 8090})
var javaPorts = dataStructs.InitializeQueue([]int{8091, 8092, 8093, 8094, 8095})
var pyPorts = dataStructs.InitializeQueue([]int{8096, 8097, 8098, 8099, 8100})
var inExecution = 0

var mutex sync.Mutex

func pushPort(Lang string, emptyPort int) {
	switch Lang {
	case "cpp":
		cppPorts.Push(emptyPort)
	case "js":
		jsPorts.Push(emptyPort)
		fmt.Println(jsPorts)
	case "java":
		javaPorts.Push(emptyPort)
	case "python":
		pyPorts.Push(emptyPort)
	}
}

func CompileCode() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req codeTypes.CodeExecutionRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			json.NewEncoder(w).Encode(map[string]string{"output": "", "error": "Failed to parse the body"})
			return
		}

		var response codeTypes.CodeExecutionResponse
		var wg sync.WaitGroup

		wg.Add(1)
		go func() {
			defer wg.Done()

			for inExecution == 5 {
			}

			mutex.Lock()
			inExecution++
			mutex.Unlock()

			var emptyPort int
			var err error

			switch req.Lang {
			case "cpp":
				emptyPort, err = cppPorts.Pop()
			case "js":
				emptyPort, err = jsPorts.Pop()
			case "java":
				emptyPort, err = javaPorts.Pop()
			case "py":
				emptyPort, err = pyPorts.Pop()
			default:
				response.Error = "Unsupported language"
				json.NewEncoder(w).Encode(response)
				mutex.Lock()
				inExecution--
				pushPort(req.Lang, emptyPort)
				mutex.Unlock()
				return
			}

			if err != nil {
				response.Error = "No available ports"
				json.NewEncoder(w).Encode(response)
				mutex.Lock()
				pushPort(req.Lang, emptyPort)
				inExecution--
				mutex.Unlock()
				return
			}

			// Prepare the request body to forward
			reqBody, err := json.Marshal(req)
			if err != nil {
				response.Error = "Failed to marshal request body"
				json.NewEncoder(w).Encode(response)
				mutex.Lock()
				pushPort(req.Lang, emptyPort)
				inExecution--
				mutex.Unlock()
				return
			}

			// Make the POST request to the specified port
			postURL := fmt.Sprintf("http://localhost:%d/compile", emptyPort)
			httpResp, err := http.Post(postURL, "application/json", bytes.NewBuffer(reqBody))
			fmt.Println(postURL)
			// fmt.Println(string(reqBody))
			// print(err)
			if err != nil {
				// response.Error = "Failed to make POST request"
				response.Error = err.Error()
				json.NewEncoder(w).Encode(response)
				mutex.Lock()
				pushPort(req.Lang, emptyPort)
				inExecution--
				mutex.Unlock()
				return
			}
			defer httpResp.Body.Close()

			// Read the response from the POST request
			respBody, err := io.ReadAll(httpResp.Body)
			if err != nil {
				response.Error = "Failed to read response body"
				json.NewEncoder(w).Encode(response)
				mutex.Lock()
				pushPort(req.Lang, emptyPort)
				inExecution--
				mutex.Unlock()
				return
			}

			// response.Output = respBody.output)
			print(respBody)

			// Push the port back to the queue after use

			mutex.Lock()
			inExecution--
			pushPort(req.Lang, emptyPort)
			mutex.Unlock()

			json.NewEncoder(w).Encode(response)
		}()

		wg.Wait()
	}
}

// 			fmt.Println(fmt.Sprintf("http://localhost:%s/compile", container.Port))
// 			resp, err := http.Post(fmt.Sprintf("http://localhost:%s/compile", container.Port), "application/json", bytes.NewBuffer(payload))
// 			if err != nil {
// 				errorChan <- fmt.Errorf("Failed to send request to container")
// 				return
// 			}

// 			defer resp.Body.Close()
// 			respBody, err := io.ReadAll(resp.Body)
// 			if err != nil {
// 				errorChan <- fmt.Errorf("Failed to read response from container")
// 				return
// 			}

// 			var result codeTypes.CodeExecutionResponse
// 			if err := json.Unmarshal(respBody, &result); err != nil {
// 				errorChan <- fmt.Errorf("Failed to unmarshal response")
// 				return
// 			}

// 			resultChan <- result
// 		}()

// 		select {
// 		case result := <-resultChan:
// 			w.Header().Set("Content-Type", "application/json")
// 			json.NewEncoder(w).Encode(result)
// 		case err := <-errorChan:
// 			json.NewEncoder(w).Encode(map[string]string{"output": "", "error": err.Error()})
// 		}
// 	}
// }
