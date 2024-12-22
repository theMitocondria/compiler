package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	compile "github.com/theMitocondria/compiler/internal/http/handlers/compile"
)

func main() {

	// Setup router
	http.HandleFunc("POST /api/v1/compile", compile.CompileCode())
	http.HandleFunc("POST /api/v1/submit", compile.SubmitCode())
	http.HandleFunc("GET /api/v1/systemTesting" , compile.SystemTesting())

	// Setup server
	server := &http.Server{
		Addr:    ":3000",
		Handler: nil,
	}

	fmt.Println("server started")

	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		err := server.ListenAndServe()
		if err != nil && err != http.ErrServerClosed {
			log.Fatalf("failed to start server: %v", err)
		}
	}()

	<-done
	fmt.Print("shutting down the server gracefully")

	ctx, cancel := context.WithTimeout(context.Background(), 4*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		fmt.Println("faled to shutdown server")
	}

	fmt.Print("shutting down the server gracefull finally")
}

/*
compile ki request ->
	to purana wala compile ka code chlal

agr testing ki aagai to ->
	compile wale prre rdicret kr de bas
	input ka loop ispe chalo or check kro

naisubmit ki requset->
  user se mila code and language
  yha prr db ko call mari or submit k test acses mangva liye with input and outputs
  yha se testing ko call mari with code lang and testcases array

system testing mai
	hum dege usko testcases  or jin bacho n us question ko submit kr diya
	inke upper loop chala kr har users kacode, lang, testcases array ko testing k pass behjio

*/
