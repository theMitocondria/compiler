package compile

import (
	"net/http"
	"sync"

	"github.com/theMitocondria/compiler/internal/http/handlers/compile"
	codeTypes "github.com/theMitocondria/compiler/internal/types/code"
)

// ccp : 8081 to 8085
// js : 86 to 90
// java
// py

var loadBalancer = codeTypes.LoadBalancer{
	Containers: []codeTypes.Container{
		// {ID: "cpp1", InUse: false, Lang: "cpp", Port: "8081"},
		// {ID: "cpp2", InUse: false, Lang: "cpp", Port: "8082"},
		// {ID: "cpp3", InUse: false, Lang: "cpp", Port: "8083"},
		// {ID: "cpp4", InUse: false, Lang: "cpp", Port: "8084"},
		// {ID: "cpp5", InUse: false, Lang: "cpp", Port: "8085"},
		// {ID: "JS1", InUse: false, Lang: "js", Port: "8086"},
		// {ID: "JS2", InUse: false, Lang: "js", Port: "8087"},
		// {ID: "JS3", InUse: false, Lang: "js", Port: "8088"},
		// {ID: "JS4", InUse: false, Lang: "js", Port: "8089"},
		// {ID: "JS5", InUse: false, Lang: "js", Port: "8090"},
		// {ID: "Java1", InUse: false, Lang: "java", Port: "8091"},
		// {ID: "Java2", InUse: false, Lang: "java", Port: "8092"},
		// {ID: "Java3", InUse: false, Lang: "java", Port: "8093"},
		// {ID: "Java4", InUse: false, Lang: "java", Port: "8094"},
		// {ID: "Java5", InUse: false, Lang: "java", Port: "8095"},
		{ID: "Py1", InUse: false, Lang: "py", Port: "8096"},
		// {ID: "Py2", InUse: false, Lang: "py", Port: "8097"},
		// {ID: "Py3", InUse: false, Lang: "py", Port: "8098"},
		// {ID: "Py4", InUse: false, Lang: "py", Port: "8099"},
		// {ID: "Py5", InUse: false, Lang: "py", Port: "8100"},
	},
	MaxLoad: 5,
}

var mu sync.Mutex

func InitRouter() *http.ServeMux {
	router := http.NewServeMux()
	router.HandleFunc("GET /api/v1/helloworld", compile.HelloWorld())
	router.HandleFunc("POST /api/v1/compile", compile.CompileCode(&loadBalancer, &mu))
	return router
}
