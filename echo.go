// MIT License
//
// Copyright (c) 2018 Yunzhu Li
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in all
// copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
// SOFTWARE.

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

var testMode = false
var srcCommit string
var exitToken string

// main initializes application and starts serving requests
func main() {
	// Configuration
	exitToken = os.Getenv("EXIT_TOKEN")
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
	log.Println("Starting service on " + port)
	log.Println("Source commit: " + srcCommit)

	if !testMode {
		http.ListenAndServe(":"+port, nil)
	}
}

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
		"/exit     causes server process to exit immediately, requires a token in header X-Exit-Token\n"+
		"/headers  returns request headers as JSON\n"+
		"/health   returns health info\n"+
		"/ip       returns client IP (use X-Forwarded-For header if exists, then remote IP)\n"+
		"\n"+
		"Source commit: "+srcCommit+"\n"+
		"</pre></html>")
}

func cacheHandler(w http.ResponseWriter, r *http.Request) {
	// Delayed response, but allows caching
	timer := time.NewTimer(500 * time.Millisecond)
	<-timer.C
	w.Header().Set("Cache-Control", "public, max-age=10") // 10 seconds

	// Response with current time
	t := time.Now()
	fmt.Fprint(w, t.Format("060102_150405"))
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
	// Validate exit token
	token := r.Header.Get("X-Exit-Token")
	if token == "" || token != exitToken {
		w.WriteHeader(http.StatusUnauthorized)
		fmt.Fprint(w, "Invalid exit token\n")
		return
	}

	// Not a graceful shutdown
	log.Println("Exiting due to /exit")
	os.Exit(0)
}

func headersHandler(w http.ResponseWriter, r *http.Request) {
	// Join header values
	headers := make(map[string]string)
	for name := range r.Header {
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
	remoteAddr := r.Header.Get("x-forwarded-for")

	// Use client IP if no xff header
	if remoteAddr == "" {
		remoteAddr = r.RemoteAddr
	}

	// Write JSON
	result := map[string]string{"remote_addr": remoteAddr}
	m, _ := json.MarshalIndent(result, "", "    ")
	fmt.Fprint(w, string(m))
}
