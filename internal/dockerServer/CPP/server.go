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

// Update struct
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
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(CodeExecutionResponse{Error: "Error encoding response: " + err.Error()})
		}
	}
}

func giveCode(code, input string) string {
	scriptWithNewCode := `
#!/bin/bash

if [ ! -d /tmp ]; then
    mkdir -p /tmp
fi

# Create a unique temporary directory
TEMP_DIR="/tmp/tempdir"
mkdir -p $TEMP_DIR

# Create temporary files
TEMP_CPP="$TEMP_DIR/temp.cpp"
TEMP_EXE="$TEMP_DIR/temp.exe"
TEMP_INPUT="$TEMP_DIR/temp.txt"

cat << EOF > $TEMP_CPP
` + code + `
EOF

cat << EOF > $TEMP_INPUT
` + input + `
EOF

OUTPUT=$(g++ -DONLINE_JUDGE=true -O2 -Wall -Wextra -Werror -std=c++17 -o $TEMP_EXE $TEMP_CPP 2>&1)

if [ ! -f $TEMP_EXE ]; then
    rm -f $TEMP_CPP $TEMP_EXE $TEMP_INPUT
    rmdir $TEMP_DIR
    echo "1,$OUTPUT"
    exit 0
fi

ulimit -v 254800 
ulimit -t 5
ulimit -f 1000

OUTPUT=$(timeout 1s $TEMP_EXE < $TEMP_INPUT 2>&1)
EXIT_CODE=$?  # Capture the exit code of the last command

if [ $EXIT_CODE -eq 143 ]; then
    OUTPUT="1,Time limit exceeded."
elif [ $EXIT_CODE -eq 139 ]; then
    OUTPUT="1,Segmentation fault."
elif [ $EXIT_CODE -eq 136 ]; then
    OUTPUT="1,Floating point exception."
elif [ $EXIT_CODE -eq 134 ]; then
    OUTPUT="1,Aborted."
elif [ $EXIT_CODE -eq 127 ]; then
    OUTPUT="1,Command not found."
elif [ $EXIT_CODE -ge 128 ] && [ $EXIT_CODE -lt 256 ]; then
    SIGNAL=$((EXIT_CODE - 128))
    OUTPUT="1,Interrupted with signal $SIGNAL."
elif [ $EXIT_CODE -ne 0 ]; then
    OUTPUT="1,Exit code $EXIT_CODE."
else
    OUTPUT="0,$OUTPUT"
fi

rm -f $TEMP_CPP $TEMP_EXE $TEMP_INPUT
rmdir $TEMP_DIR
echo "$OUTPUT"

`
	return scriptWithNewCode
}
