package main

import (
	"math/rand"
	"net/http"

	"github.com/DataDog/datadog-go/statsd"
	log "github.com/sirupsen/logrus"
)

func main() {

	client, err := statsd.New("127.0.0.1:8125",
		statsd.WithNamespace("echoServer."),            // prefix every metric with the app name
		statsd.WithTags([]string{"region:us-east-1a"}), // send the EC2 availability zone as a tag with every metric
	)
	if err != nil {
		log.Fatal(err)
	}

	// Create our server
	logger := log.New()
	server := Server{
		logger: logger,
		client: client,
	}

	// Start the server
	server.ListenAndServe()
}

// Server represents our server.
type Server struct {
	logger *log.Logger
	client *statsd.Client
}

// ListenAndServe starts the server
func (s *Server) ListenAndServe() {
	s.logger.Info("echo server is starting on port 8080...")
	http.HandleFunc("/", s.echo)
	http.ListenAndServe(":8080", nil)
}

// Echo echos back the request as a response
func (s *Server) echo(writer http.ResponseWriter, request *http.Request) {
	writer.Header().Set("Access-Control-Allow-Origin", "*")
	writer.Header().Set("Access-Control-Allow-Headers", "Content-Range, Content-Disposition, Content-Type, ETag")

	// 30% chance of failure
	if rand.Intn(100) < 30 {
		s.logger.Error("Unlucky Request")
		s.client.Count("Request_failed.counter", 1, []string{"env:production", "partition:1", "partition:2"}, 1)
		writer.WriteHeader(500)
		writer.Write([]byte("a chaos monkey broke your server"))
		return
	}

	// Happy path
	s.client.Count("Request_success.counter", 1, []string{"env:production", "partition:1", "partition:2"}, 1)
	writer.WriteHeader(200)
	request.Write(writer)
}
