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

func findMainClass(code string) (string, error) {
	cmd := exec.Command("python3", "/usr/local/bin/server.py", code)
	out, err := cmd.CombinedOutput()
	if err != nil {
		return "", err
	}
	className := strings.TrimSpace(string(out))
	if className == "No main method found" {
		return "", nil
	}
	return className, nil
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

		response := CodeExecutionResponse{}
		className, err := findMainClass(req.Code)
		if err != nil {
			response.Error = "Failed to find main class: " + err.Error()
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(response)
			return
		}
		if className == "" {
			response.Error = "No main method found in the provided code"
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(response)
			return
		}

		scriptWithNewCode := giveCode(req.Code, req.Input, className)
		cmd := exec.Command("sh", "-c", scriptWithNewCode)
		out, err := cmd.CombinedOutput()

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

func giveCode(code, input string, className string) string {
	scriptWithNewCode := `
#!/bin/bash

if [ ! -d /tmp ]; then
    mkdir -p /tmp
fi

# Create unique temporary directory
TEMP_DIR="/tmp/tempdir"
mkdir -p $TEMP_DIR

# Create Java file
cat << EOF > "$TEMP_DIR/Main.java"
` + code + `
EOF

# Extract class name using Python script
CLASS_NAME=` + className + `

if [ "$CLASS_NAME" == "No main method found" ]; then
    echo "1,Error: No class with a main method found."
    rm -rf $TEMP_DIR
    exit 0
fi

rm -f $TEMP_DIR/Main.java

# Create input file
cat << EOF > "$TEMP_DIR/input.txt"
` + input + `
EOF

cat << EOF > "$TEMP_DIR/$CLASS_NAME.java"
` + code + `
EOF
# Compile Java file
javac "$TEMP_DIR/$CLASS_NAME.java" 2> "$TEMP_DIR/compile_errors.txt"
EXIT_CODE=$?

if [ $EXIT_CODE -ne 0 ]; then
    ERROR=$(cat "$TEMP_DIR/compile_errors.txt")
    rm -rf $TEMP_DIR
    echo "1,$ERROR"
    exit 0
fi


# Execute Java program
OUTPUT=$(timeout 1s java -XX:+UseSerialGC -XX:TieredStopAtLevel=1 -XX:NewRatio=5 -Xms8M -Xmx128M -Xss64M -DONLINE_JUDGE=true -classpath $TEMP_DIR $CLASS_NAME < "$TEMP_DIR/input.txt" 2>&1)
EXIT_CODE=$?

if [ $EXIT_CODE -eq 124 ] || [ $EXIT_CODE -eq 143 ]; then
    OUTPUT="1,Time limit exceeded"
elif [ $EXIT_CODE -eq 137 ]; then
    OUTPUT="1,Memory limit exceeded"
elif [ $EXIT_CODE -ne 0 ]; then
    OUTPUT="1,Runtime Error: $OUTPUT"
else
    OUTPUT="0,$OUTPUT"
fi

# Cleanup
rm -rf $TEMP_DIR
echo "$OUTPUT"

`
	return scriptWithNewCode
}
