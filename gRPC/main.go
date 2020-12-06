package main

import (
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/reflection"
)

type SSLCerts struct {
	// certificate file path
	CertificateFP string `json:"certfilepath"`
	// key file path
	KeyFP string `json:"keyfilepath"`
}

type ConnAttr struct {
	Host      string   `json:"host"`
	Port      int      `json:"port"`
	EnableTLS bool     `json:"enableTLS"`
	Secrets   SSLCerts `json:"secrets"`
}

func newConnAttr() *ConnAttr {
	return &ConnAttr{
		Host:      "0.0.0.0",
		Port:      505001,
		EnableTLS: false,
		Secrets: SSLCerts{
			CertificateFP: "",
			KeyFP:         ""}}
}

const (
	ip   = "localhost"
	port = 505001
	tls  = false
)

func main() {
	cleanup := make(chan bool, 1)
	ch := make(chan os.Signal, 1)
	newServerwithServices(cleanup)
	signal.Notify(ch, os.Interrupt)
	<-ch
	cleanup <- true
	<-cleanup
}

func getServer(connAttr *ConnAttr) (*grpc.Server, *net.Listener, error) {
	lis, err := net.Listen("tcp", fmt.Sprintf("%s:%d", ip, port))
	if err != nil {
		return nil, nil, fmt.Errorf("Failed to open a TCP connection at port %d", port)
	}
	opts := []grpc.ServerOption{}

	if creds, err := credentials.NewServerTLSFromFile(connAttr.Secrets.CertificateFP, connAttr.Secrets.KeyFP); err != nil {
		return nil, nil, fmt.Errorf("Failed loading certificates: %s", err.Error())
	} else {
		opts = append(opts, grpc.Creds(creds))
	}
	s := grpc.NewServer(opts...)
	return s, &lis, nil
}

func newServerwithServices(cleanup chan bool) {
	connAttr := newConnAttr()
	s, lis, err := getServer(connAttr)
	if err != nil {
		log.Fatalf("%s", err.Error())
	}
	// Register the services here
	reflection.Register(s)
	go func() {
		fmt.Println("Starting the server")
		if err := s.Serve(*lis); err != nil {
			log.Fatalf("Failed to serve")
		}
	}()
	<-cleanup
	if err := (*lis).Close(); err != nil {
		log.Fatalf("Error on closing the listener : %v", err)
	}
	s.Stop()
	cleanup <- true
}
