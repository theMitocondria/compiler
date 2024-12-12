package compile

import (
	"net/http"

	"github.com/theMitocondria/compiler/internal/http/handlers/compile"
)

func InitRouter() *http.ServeMux {
	router := http.NewServeMux()
	router.HandleFunc("GET /api/v1/helloworld", compile.HelloWorld())
	router.HandleFunc("POST /api/v1/compile", compile.CompileCode())
	return router
}
