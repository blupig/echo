package main

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"
)

// rootHandler handles requests to /
func rootHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}

	fmt.Fprint(w, "<html><pre><b>routes:</b>\n"+
		"/         root (this route)\n"+
		"/cache    returns cacheable but delayed (500ms) response\n"+
		"/cpu      CPU-intensive operation\n"+
		"/exit     causes server process to exit immediately\n"+
		"/headers  returns request headers as JSON\n"+
		"/health   returns health info\n"+
		"/ip       returns client IP\n"+
		"</pre></html>")
}

func cacheHandler(w http.ResponseWriter, r *http.Request) {
	// Delayed response, but allows caching
	t := time.NewTimer(500 * time.Millisecond)
	<-t.C
	w.Header().Set("Cache-Control", "max-age=60") // 1 minute
	fmt.Fprint(w, "ok")
}

func cpuHandler(w http.ResponseWriter, r *http.Request) {
	// Handlers are already executed asynchronously
	fmt.Fprint(w, cpuLoad())
}

// cpuLoad performs CPU-intensive operations
func cpuLoad() string {
	for i := 0; i < 1000000; i++ {
		sha256.Sum256([]byte("abc"))
	}
	return fmt.Sprint("")
}

func exitHandler(w http.ResponseWriter, r *http.Request) {
	// Not a graceful shutdown
	log.Println("Exiting...")
	os.Exit(0)
}

func headersHandler(w http.ResponseWriter, r *http.Request) {
	// Join header values
	headers := make(map[string]string)
	for name, _ := range r.Header {
		headers[name] = r.Header.Get(name)
	}

	// Marshal with indent
	m, _ := json.MarshalIndent(headers, "", "    ")
	fmt.Fprint(w, string(m))
}

func healthHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "ok")
}

func ipHandler(w http.ResponseWriter, r *http.Request) {
	// Try to get remote IP from xff header first
	remote_addr := r.Header.Get("x-forwarded-for")

	// Use client IP if no xff header
	if remote_addr == "" {
		remote_addr = r.RemoteAddr
	}

	// Write JSON
	result := map[string]string{"remote_addr": remote_addr}
	m, _ := json.MarshalIndent(result, "", "    ")
	fmt.Fprint(w, string(m))
}

// main initializes application
func main() {
	// Configuration
	port := os.Getenv("PORT")
	if port == "" {
		port = "8000"
	}

	// Routes
	http.HandleFunc("/", rootHandler)
	http.HandleFunc("/cache", cacheHandler)
	http.HandleFunc("/cpu", cpuHandler)
	http.HandleFunc("/exit", exitHandler)
	http.HandleFunc("/headers", headersHandler)
	http.HandleFunc("/health", healthHandler)
	http.HandleFunc("/ip", ipHandler)

	// Start serving
	log.Println("Starting service...")
	http.ListenAndServe(":"+port, nil)
}
