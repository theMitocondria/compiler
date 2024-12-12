package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os/exec"
)

type CodeExecutionRequest struct {
	Code  string `json:"code"`
	Lang  string `json:"lang"`
	Input string `json:"input"`
}

type CodeExecutionResponse struct {
	Output string `json:"output"`
	Error  string `json:"error"`
}

func main() {
	router := http.NewServeMux()
	router.HandleFunc("/compile", runCode())

	server := http.Server{
		Addr:    "0.0.0.0:8080",
		Handler: router,
	}

	if err := server.ListenAndServe(); err != nil {
		log.Fatalf("failed to start server: %s", err.Error())
	}
}

func runCode() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
			return
		}

		var req CodeExecutionRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "Failed to parse request body", http.StatusBadRequest)
			return
		}

		if req.Code == "" || req.Lang == "" {
			http.Error(w, "Missing code or language in request", http.StatusBadRequest)
			return
		}

		if req.Lang != "py" {
			http.Error(w, "Unsupported language", http.StatusBadRequest)
			return
		}

		scriptWithNewCode := `
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
` +
			req.Code +
			`
EOF

# Syntax check the Python code
python3 -m py_compile temp.py 2>/tmp/syntax_error.log
EXIT_CODE=$?  # Capture the exit code of the syntax check

if [ $EXIT_CODE -ne 0 ]; then
    # If there's a syntax error, display it
    cat /tmp/syntax_error.log
else
    # Set a memory limit and run the Python script with a timeout
    ulimit -v 284800  # Limit memory to 284 MB
    OUTPUT=$(echo "` + req.Input + `" | timeout 1s python3 temp.py 2>&1)
    EXIT_CODE=$?

    # Handle different exit codes
    if [ $EXIT_CODE -eq 143 ]; then
        OUTPUT="TLE"
    elif [ $EXIT_CODE -eq 1 ]; then
        OUTPUT="Error: Runtime Error"
    elif [ $EXIT_CODE -eq 139 ]; then
        OUTPUT="Error: Memory Limit Exceeded"
    elif [ $EXIT_CODE -ne 0 ]; then
        OUTPUT="Error: Program terminated unexpectedly. $EXIT_CODE"
    else
        TLE="false"
    fi

    # Cleanup the Python script
    rm temp.py
    # Output the result
    echo "$OUTPUT"
fi`
		fmt.Println(scriptWithNewCode)
		cmd := exec.Command("sh", "-c", scriptWithNewCode)
		out, err := cmd.CombinedOutput()

		var response CodeExecutionResponse

		if err != nil {
			response.Error = err.Error()
		} else {
			response.Output = string(out)
		}

		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(response); err != nil {
			http.Error(w, "Error encoding response: "+err.Error(), http.StatusInternalServerError)
		}
	}
}
