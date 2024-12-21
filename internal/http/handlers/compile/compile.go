package compile

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"sync"
	"time"

	codeTypes "github.com/theMitocondria/compiler/internal/types/code"
	"github.com/theMitocondria/compiler/internal/utils/dataStructs"
)

// Port queues for each language
var cppPorts = dataStructs.InitializeQueue([]int{8081, 8082, 8083, 8084, 8085, 8086, 8087, 8088, 8089, 8090})
var jsPorts = dataStructs.InitializeQueue([]int{8091, 8092, 8093, 8094, 8095, 8096, 8097, 8098, 8099, 8100})
var javaPorts = dataStructs.InitializeQueue([]int{8101, 8102, 8103, 8104, 8105, 8106, 8107, 8108, 8109, 8110})
var pyPorts = dataStructs.InitializeQueue([]int{8111, 8112, 8113, 8114, 8115, 8116, 8117, 8118, 8119, 8120})

// Mutexes for port queues
var cppPortsMutex, jsPortsMutex, javaPortsMutex, pyPortsMutex sync.Mutex

// Job queue and worker pool
var jobQueue = make(chan CodeExecutionRequestWithResponse, 500) // Buffer size can be adjusted
var workers = 10                                                // Number of workers
var wg sync.WaitGroup

type CodeExecutionRequestWithResponse struct {
	Request  codeTypes.CodeExecutionRequest
	Response chan<- codeTypes.CodeExecutionResponse
}

func init() {
	// Start the workers
	for i := 0; i < workers; i++ {
		go worker()
	}
}

func pushPort(Lang string, emptyPort int) {
	switch Lang {
	case "cpp":
		cppPortsMutex.Lock()
		cppPorts.Push(emptyPort)
		cppPortsMutex.Unlock()
	case "js":
		jsPortsMutex.Lock()
		jsPorts.Push(emptyPort)
		jsPortsMutex.Unlock()
	case "java":
		javaPortsMutex.Lock()
		javaPorts.Push(emptyPort)
		javaPortsMutex.Unlock()
	case "python":
		pyPortsMutex.Lock()
		pyPorts.Push(emptyPort)
		pyPortsMutex.Unlock()
	}
}

func worker() {
	for reqWithResponse := range jobQueue {
		var response codeTypes.CodeExecutionResponse
		processRequest(reqWithResponse.Request, &response)
		reqWithResponse.Response <- response
		wg.Done()
	}
}

func processRequest(req codeTypes.CodeExecutionRequest, response *codeTypes.CodeExecutionResponse) {
	var emptyPort int
	var err error

	switch req.Lang {
	case "cpp":
		cppPortsMutex.Lock()
		emptyPort, err = cppPorts.Pop()
		cppPortsMutex.Unlock()
	case "js":
		jsPortsMutex.Lock()
		emptyPort, err = jsPorts.Pop()
		jsPortsMutex.Unlock()
	case "java":
		javaPortsMutex.Lock()
		emptyPort, err = javaPorts.Pop()
		javaPortsMutex.Unlock()
	case "python":
		pyPortsMutex.Lock()
		emptyPort, err = pyPorts.Pop()
		pyPortsMutex.Unlock()
	default:
		*response = codeTypes.CodeExecutionResponse{Error: "Unsupported language"}
		return
	}

	// Prepare the request body to forward
	reqBody, err := json.Marshal(req)

	if err != nil {
		*response = codeTypes.CodeExecutionResponse{Error: "Failed to marshal request body"}
		pushPort(req.Lang, emptyPort)
		return
	}

	var portName string
	if emptyPort-8080 <= 10 {
		portName = fmt.Sprintf("cpp%d", emptyPort-8080)
	} else if emptyPort-8080 <= 20 {
		portName = fmt.Sprintf("js%d", emptyPort-8090)
	} else if emptyPort-8080 <= 30 {
		portName = fmt.Sprintf("java%d", emptyPort-8100)
	} else {
		portName = fmt.Sprintf("py%d", emptyPort-8110)
	}
	// Make the POST request to the specified port
	postURL := fmt.Sprintf("http://%s:8080/compile", portName)
	fmt.Println(portName)

	httpResp, err := http.Post(postURL, "application/json", bytes.NewBuffer(reqBody))
	if err != nil {
		*response = codeTypes.CodeExecutionResponse{Error: err.Error()}
		pushPort(req.Lang, emptyPort)
		return
	}

	defer httpResp.Body.Close()

	// Read the response from the POST request
	respBody, err := io.ReadAll(httpResp.Body)

	if err != nil {
		*response = codeTypes.CodeExecutionResponse{Error: "Failed to read response body"}
		pushPort(req.Lang, emptyPort)
		return
	}

	var compilerResponse codeTypes.CodeExecutionResponse
	if err := json.Unmarshal(respBody, &compilerResponse); err != nil {
		*response = codeTypes.CodeExecutionResponse{Error: "Failed to unmarshal response body"}
		pushPort(req.Lang, emptyPort)

		return
	}

	*response = codeTypes.CodeExecutionResponse{
		Output: compilerResponse.Output,
		Error:  compilerResponse.Error,
	}

	// Push the port back to the queue after use
	pushPort(req.Lang, emptyPort)
}

func CompileCode() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req codeTypes.CodeExecutionRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]string{"output": "", "error": "Failed to parse the body"})
			return
		}

		responseChan := make(chan codeTypes.CodeExecutionResponse)

		wg.Add(1)
		jobQueue <- CodeExecutionRequestWithResponse{
			Request:  req,
			Response: responseChan,
		}

		select {
		case response := <-responseChan:
			if w.Header().Get("Content-Type") == "" {
				w.WriteHeader(http.StatusOK)
			}
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(response)
		case <-time.After(60 * time.Second): // Timeout after 30 seconds
			w.WriteHeader(http.StatusRequestTimeout)
			json.NewEncoder(w).Encode(map[string]string{"output": "", "error": "Compilation timed out please retry"})
		}
	}
}
