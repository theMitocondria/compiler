package codeTypes

import "sync"

type CodeExecutionRequest struct {
	Code  string `json:"code"`
	Lang  string `json:"lang"`
	Input string `json:"input"`
}

type CodeExecutionResponse struct {
	Output struct {
		CodeOutput string `json:"code_output"`
		CodeError  string `json:"code_error"`
	} `json:"output"`
	Error string `json:"error"`
}

type Container struct {
	ID    string
	InUse bool
	Lang  string
	Port  string
}

type LoadBalancer struct {
	Containers []Container
	MaxLoad    int
}

type RequestQueue struct {
	Queue []CodeExecutionRequest
	Mu    sync.Mutex
}
