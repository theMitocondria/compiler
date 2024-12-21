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

# Create JavaScript file
cat << EOF > temp.js
` + code + `
EOF


ulimit -v 254800 

OUTPUT=$(echo "` + input + `" | timeout 1s node temp.js 2>&1)
EXIT_CODE=$? 

# Handle different exit codes
if [ $EXIT_CODE -eq 143 ]; then
    OUTPUT="1,TLE"
elif [ $EXIT_CODE -eq 139 ]; then
    OUTPUT="1,Error: Segmentation fault (invalid memory access)."
elif [ $EXIT_CODE -eq 136 ]; then
    OUTPUT="1,Error: Floating-point exception."
elif [ $EXIT_CODE -eq 134 ]; then
    OUTPUT="1,Error: Program aborted."
elif [ $EXIT_CODE -eq 127 ]; then
    OUTPUT="1,Error: Executable not found."
elif [ $EXIT_CODE -ge 128 ] && [ $EXIT_CODE -lt 256 ]; then
    SIGNAL=$((EXIT_CODE - 128))
    OUTPUT="1,Error: Program terminated by signal $SIGNAL."
elif [ $EXIT_CODE -ne 0 ]; then
    OUTPUT="1,Error: Exit code $EXIT_CODE."
else
    OUTPUT="0,$OUTPUT"
fi

# Cleanup temporary files
rm -f temp.js temp_input.txt

echo "$OUTPUT"
`
	return scriptWithNewCode
}
