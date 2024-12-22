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

var inputValues = []string{"khatri", "dhurv", "gautam", "navya", "mehta", "shubham", "sahil", "saurabh", "siddharth", "siddharth"}
var outputValues = []string{"irtahk", "vruhd", "matuag", "ayvan", "athem", "mahbuhs", "lihas", "hbaruas", "htrahddis", "htrahddis"}

func SubmitCode() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req codeTypes.CodeExecutionRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]string{"output": "", "error": "Failed to parse the body"})
			return
		}

		var compilerResponse codeTypes.CodeExecutionResponse

		for i, input := range inputValues {
			req.Input = input
			responseChan := make(chan codeTypes.CodeExecutionResponse)

			wg.Add(1)
			jobQueue <- CodeExecutionRequestWithResponse{
				Request:  req,
				Response: responseChan,
			}

			select {
			case response := <-responseChan:

				if response.Error != "" {
					compilerResponse.Error = response.Error
					w.WriteHeader(http.StatusOK)
					json.NewEncoder(w).Encode(compilerResponse)
					return
				} else if response.Output.CodeError != "" {
					compilerResponse.Output.CodeError = response.Output.CodeError
					w.WriteHeader(http.StatusOK)
					json.NewEncoder(w).Encode(compilerResponse)
					return
				} else if response.Output.CodeOutput != outputValues[i] {
					compilerResponse.Output.CodeError = fmt.Sprintf("Output mismatch for input %s: expected %s, got %s", input, outputValues[i], response.Output.CodeOutput)
					w.WriteHeader(http.StatusOK)
					json.NewEncoder(w).Encode(compilerResponse)
					return
				}
			case <-time.After(60 * time.Second): // Timeout after 60 seconds
				w.WriteHeader(http.StatusRequestTimeout)
				json.NewEncoder(w).Encode(map[string]string{"output": "", "error": "Compilation timed out please retry"})
				return
			}
		}

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"output": map[string]string{
				"code_output": "All test cases passed",
				"code_error":  "",
			},
			"error": "",
		})

	}
}

type UserCode struct {
	UserName string
	Code     string
	Lang     string
}

func SystemTesting() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		usersCode := []UserCode{
			{
				UserName: "User1",
				Code:     "def reverse_string(s):\n    return s[::-1]\ninput_string = input()\nreversed_string = reverse_string(input_string)\nprint(reversed_string)",
				Lang:     "python",
			},
			{
				UserName: "User2",
				Code: `
function reverseString(s) {
    return s.split("").reverse().join("");
}

var inputString ;
inputString = prompt();
const reversedString = reverseString(inputString);
console.log("Reversed string in JavaScript: " + reversedString);
                `,
				Lang: "js",
			},
			{
				UserName: "User3",
				Code: `
#include <iostream>
#include <string>
#include <algorithm>

std::string reverseString(const std::string &s) {
    std::string reversed = s;
    std::reverse(reversed.begin(), reversed.end());
    return reversed;
}

int main() {
    std::string inputString;
    std::cin >> inputString;
    std::string reversedString = reverseString(inputString);
    std::cout << reversedString;
    return 0;
}
                `,
				Lang: "cpp",
			},
			{
				UserName: "User4",
				Code: `
def reverse_string(s):
    return s[::-1]
input_string = input()
reversed_string = reverse_string(input_string)
print(f"Reversed string in Python: {reversed_string}")
                `,
				Lang: "python",
			},
		}

		failedUsers := []string{}

		for _, userCode := range usersCode {
			req := codeTypes.CodeExecutionRequest{
				Code: userCode.Code,
				Lang: userCode.Lang,
			}
			reqBody, err := json.Marshal(req)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				json.NewEncoder(w).Encode(map[string]string{"error": "Failed to marshal request body"})
				return
			}
			fmt.Println(req.Code)
			resp, err := http.Post("http://localhost:3000/api/v1/submit", "application/json", bytes.NewBuffer(reqBody))
			if err != nil {
				// failedUsers = append(failedUsers, userCode.UserName)
				continue
			}

			defer resp.Body.Close()
			respBody, err := io.ReadAll(resp.Body)
			if err != nil {
				// failedUsers = append(failedUsers, userCode.UserName)
				fmt.Println(err)
				continue
			}

			// var result map[string]string
			var result codeTypes.CodeExecutionResponse
			if err := json.Unmarshal(respBody, &result); err != nil {
				// failedUsers = append(failedUsers, userCode.UserName)
				fmt.Println(err)
				continue
			}

			// fmt.Println((result))

			if result.Output.CodeError != "" {
				failedUsers = append(failedUsers, userCode.UserName)
			}
		}

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(failedUsers)
	}
}
