package codeTypes

type CompileCodeRequest struct {
	Code  string
	Lang  string
	Input string
}

type CompiledCodeResponse struct {
	Output string
	Error  string
}
