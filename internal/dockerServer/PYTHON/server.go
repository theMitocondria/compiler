package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os/exec"
	"strings"
	"time"
)

type CodeExecutionRequest struct {
	Code  string `json:"code"`
	Input string `json:"input"`
}

type CodeExecutionResponse struct {
	Output struct {
		CodeOutput string `json:"code_output"`
		CodeError  string `json:"code_error"`
	} `json:"output"`
	Error string `json:"error"`
}

func main() {
	router := http.NewServeMux()
	router.HandleFunc("/compile", RunCode())

	server := http.Server{
		Addr:         ":8080",
		Handler:      router,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  15 * time.Second,
	}

	log.Println("Starting server on :8080")
	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("failed to start server: %s", err.Error())
	}
}

func RunCode() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req CodeExecutionRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			response := CodeExecutionResponse{
				Error: "Failed to parse request body",
			}
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(response)
			return
		}

		scriptWithNewCode := giveCode(req.Code, req.Input)
		cmd := exec.Command("sh", "-c", scriptWithNewCode)
		out, err := cmd.CombinedOutput()

		response := CodeExecutionResponse{}

		if err != nil {
			response.Error = err.Error()
		} else {
			// Parse output
			outputStr := string(out)
			parts := strings.SplitN(strings.TrimSpace(outputStr), ",", 2)

			if len(parts) == 2 {
				if parts[0] == "0" {
					response.Output.CodeOutput = parts[1]
				} else {
					response.Output.CodeError = parts[1]
				}
			}
		}

		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(response); err != nil {
			http.Error(w, "Error encoding response: "+err.Error(), http.StatusInternalServerError)
		}
	}
}

func giveCode(code, input string) string {
	scriptWithNewCode := `
#!/bin/bash

cat << 'EOF' > temp.py
import builtins
import os

builtins.eval = None
builtins.exec = None
builtins.open = None
os.setuid = None 
os.system = None
os.stetgid = None

# User code goes here
` + code + `
EOF

# Create input file
cat << EOF > temp_input.txt
` + input + `
EOF

# Syntax check the Python code
python3 -m py_compile temp.py 2>/tmp/syntax_error.log
EXIT_CODE=$?  # Capture the exit code of the syntax check

if [ $EXIT_CODE -ne 0 ]; then
    # If there's a syntax error, display it
    ERROR=$(cat /tmp/syntax_error.log)
	rm temp.py temp_input.txt
    echo "1,$ERROR"
else
    # Set a memory limit and run the Python script with a timeout
    ulimit -v 284800  # Limit memory to 284 MB
    OUTPUT=$(timeout 1s python3 temp.py < temp_input.txt 2>&1)
    EXIT_CODE=$?

    # Handle different exit codes
    if [ $EXIT_CODE -eq 143 ]; then
        OUTPUT="1,TLE"
    elif [ $EXIT_CODE -eq 1 ]; then
        OUTPUT="1,Error: Runtime Error"
    elif [ $EXIT_CODE -eq 139 ]; then
        OUTPUT="1,Error: Memory Limit Exceeded"
    elif [ $EXIT_CODE -ne 0 ]; then
        OUTPUT="1,Error: Program terminated unexpectedly. $EXIT_CODE"
    else
        OUTPUT="0,$OUTPUT"
    fi

    # Cleanup the Python script and input file
    rm temp.py temp_input.txt
    # Output the result
    echo "$OUTPUT"
fi
`
	return scriptWithNewCode
}
