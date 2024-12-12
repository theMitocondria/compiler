package main

import (
	"encoding/json"
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

		var req CodeExecutionRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			response := CodeExecutionResponse{
				Output: "",
				Error:  "Failed to parse request body",
			}
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(response)
			return
		}

		scriptWithNewCode := giveCode(req.Code, req.Input)

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
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(CodeExecutionResponse{Output: "", Error: "Error encoding response: " + err.Error()})
		}
	}
}

func giveCode(code, input string) string {
	scriptWithNewCode := `
cat << EOF > temp.js
` + code + `

EOF


ulimit -v 254800 

OUTPUT=$(echo "` + input + `" | timeout 1s node check.js 2>&1)
EXIT_CODE=$?  # Capture the exit code of the last command

if [ $EXIT_CODE -eq 143  ]; then
    OUTPUT="TLE"
elif [ $EXIT_CODE -eq 139 ]; then
    OUTPUT="Error: Segmentation fault (invalid memory access)."
elif [ $EXIT_CODE -eq 136 ]; then
    OUTPUT="Error: Floating-point exception."
elif [ $EXIT_CODE -eq 134 ]; then
    OUTPUT="Error: Program aborted."
elif [ $EXIT_CODE -eq 127 ]; then
    OUTPUT="Error: Executable not found."
elif [ $EXIT_CODE -ge 128 ] && [ $EXIT_CODE -lt 256 ]; then
    SIGNAL=$((EXIT_CODE - 128))
    OUTPUT="Error: Program terminated by signal $SIGNAL."
fi

if [ -f temp.js ];then 
    rm temp.js
fi


echo "$OUTPUT , $EXIT_CODE"

`
	return scriptWithNewCode
}
