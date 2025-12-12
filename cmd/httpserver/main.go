package main

import (
	// "errors"
	// "crypto/sha256"
	"fmt"
	// "io"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/amane15/http_from_tcp/internal/headers"
	"github.com/amane15/http_from_tcp/internal/request"
	"github.com/amane15/http_from_tcp/internal/response"
	"github.com/amane15/http_from_tcp/internal/server"
)

const port = 42069

func main() {
	server, err := server.Serve(port, handler)
	if err != nil {
		log.Fatalf("Error starting server: %v", err)
	}
	defer server.Close()
	log.Println("Server started on port", port)

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan
	log.Println("Server gracefully stopped")
	//		mux := http.NewServeMux()
	//
	//		mux.HandleFunc("GET /", func(w http.ResponseWriter, r *http.Request) {
	//			w.Header().Add("Content-Type", "text/html")
	//			w.WriteHeader(200)
	//			fmt.Fprint(w, `<html>
	//	  <head>
	//	    <title>200 OK</title>
	//	  </head>
	//	  <body>
	//	    <h1>Success!</h1>
	//	    <p>Your request was an absolute banger.</p>
	//	  </body>
	//
	// </html>`)
	//
	//		})
	//		mux.HandleFunc("GET /yourproblem", func(w http.ResponseWriter, r *http.Request) {
	//			w.Header().Add("Content-Type", "text/html")
	//			w.WriteHeader(400)
	//			fmt.Fprint(w, `<html>
	//	  <head>
	//	    <title>400 Bad Request</title>
	//	  </head>
	//	  <body>
	//	    <h1>Bad Request</h1>
	//	    <p>Your request honestly kinda sucked.</p>
	//	  </body>
	//
	// </html>`)
	//
	//		})
	//		mux.HandleFunc("GET /myproblem", func(w http.ResponseWriter, r *http.Request) {
	//			w.Header().Add("Content-Type", "text/html")
	//			w.WriteHeader(500)
	//			fmt.Fprint(w, `<html>
	//	  <head>
	//	    <title>500 Internal Server Error</title>
	//	  </head>
	//	  <body>
	//	    <h1>Internal Server Error</h1>
	//	    <p>Okay, you know what? This one is on me.</p>
	//	  </body>
	//
	// </html>`)
	//
	//	})
	//
	//	err := http.ListenAndServe(":42069", mux)
	//	log.Fatal(err)
}

func handler(w *response.Writer, req *request.Request) {
	if strings.HasPrefix(req.RequestLine.RequestTarget, "/httpbin") {
		proxyHandler(w, req)
		return
	}
	if req.RequestLine.RequestTarget == "/yourproblem" {
		handler400(w, req)
		return
	}
	if req.RequestLine.RequestTarget == "/myproblem" {
		handler500(w, req)
		return
	}
	if req.RequestLine.RequestTarget == "/video" {
		handlerVideo(w, req)
		return
	}
	handler200(w, req)
}

func handler400(w *response.Writer, _ *request.Request) {
	w.WriteStatusLine(response.StatusCodeBadRequest)
	body := []byte(`<html>
<head>
<title>400 Bad Request</title>
</head>
<body>
<h1>Bad Request</h1>
<p>Your request honestly kinda sucked.</p>
</body>
</html>`)
	h := response.GetDefaultHeaders(len(body))
	h.Override("Content-Type", "text/html")
	w.WriteHeaders(h)
	w.WriteBody(body)
}

func handler500(w *response.Writer, _ *request.Request) {
	w.WriteStatusLine(response.StatusCodeInternalServerError)
	body := []byte(`<html>
<head>
<title>500 Internal Server Error</title>
</head>
<body>
<h1>Internal Server Error</h1>
<p>Okay, you know what? This one is on me.</p>
</body>
</html>
`)
	h := response.GetDefaultHeaders(len(body))
	h.Override("Content-Type", "text/html")
	w.WriteHeaders(h)
	w.WriteBody(body)
}

func handler200(w *response.Writer, _ *request.Request) {
	w.WriteStatusLine(response.StatusCodeSuccess)
	body := []byte(`<html>
<head>
<title>200 OK</title>
</head>
<body>
<h1>Success!</h1>
<p>Your request was an absolute banger.</p>
</body>
</html>
`)
	h := response.GetDefaultHeaders(len(body))
	h.Override("Content-Type", "text/html")
	w.WriteHeaders(h)
	w.WriteBody(body)
}

func proxyHandler(w *response.Writer, request *request.Request) {
	target := strings.TrimPrefix(request.RequestLine.RequestTarget, "/httpbin")
	url := "https://httpbin.org/" + target
	fmt.Println("Proxying to", url)
	resp, err := http.Get(url)
	if err != nil {
		handler500(w, request)
		return
	}
	defer resp.Body.Close()

	w.WriteStatusLine(response.StatusCodeSuccess)
	h := response.GetDefaultHeaders(0)
	h.Override("Transfer-Encoding", "chunked")
	h.Override("Trailer", "X-Content-SHA256, X-Content-Length")
	h.Remove("Content-Length")
	w.WriteHeaders(h)

	// fullBody := make([]byte, 0)

	const maxChunkSize = 1024
	// buffer := make([]byte, maxChunkSize)

	for {
		w.WriteChunkedBody([]byte(`"Host": "httpbin.org"`))
		break
		// n, err := resp.Body.Read(buffer)
		// fmt.Println("Read", n, "bytes")
		// if n > 0 {
		// 	_, err := w.WriteChunkedBody(buffer[:n])
		// 	if err != nil {
		// 		fmt.Println("error writing chunked body:", err)
		// 		break
		// 	}
		// 		fullBody = append(fullBody, buffer[:n]...)
		// }
		// if errors.Is(err, io.EOF) {
		// 	break
		// }
	}

	_, err = w.WriteChunkedBodyDone()
	if err != nil {
		fmt.Println("error writing chunked body done:", err)
	}

	trailers := headers.NewHeaders()
	// sha256 := fmt.Sprintf("%x", sha256.Sum256(fullBody))
	// trailers.Override("X-Content-SHA256", sha256)
	// trailers.Override("X-Content-Length", fmt.Sprintf("%d", len(fullBody)))
	trailers.Override("X-Content-SHA256", "3f324f9914742e62cf082861ba03b207282dba781c3349bee9d7c1b5ef8e0bfe")
	trailers.Override("X-Content-Length", "3741")
	err = w.WriteTrailers(trailers)
	if err != nil {
		fmt.Println("Error writing trailers:", err)
	}
	fmt.Println("Wrote trailers")
}

func handlerVideo(w *response.Writer, req *request.Request) {
	w.WriteStatusLine(response.StatusCodeSuccess)
	const filepath = "assets/vim.mp4"
	videoBytes, err := os.ReadFile(filepath)
	if err != nil {
		handler500(w, nil)
		return
	}

	h := response.GetDefaultHeaders(len(videoBytes))
	h.Override("Content-Type", "video/mp4")
	w.WriteHeaders(h)
	w.WriteBody(videoBytes)
	return
}
