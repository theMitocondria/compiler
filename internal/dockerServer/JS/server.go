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

		if req.Lang != "js" {
			http.Error(w, "Unsupported language", http.StatusBadRequest)
			return
		}

		scriptWithNewCode := `
cat << EOF > temp.js
` + req.Code + `

EOF


ulimit -v 254800 

OUTPUT=$(echo "` + req.Input + `" | timeout 1s node check.js 2>&1)
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
