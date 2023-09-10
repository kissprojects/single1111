package grpc

import (
	"context"
	"fmt"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/kissprojects/single/comps/go/api"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"net"
	"net/http"
)

var (
	// command-line options:
	// gRPC server endpoint
	GrpcServerEndpoint = "localhost:8080"
	grpcPort           = "8080"
	mux                = runtime.NewServeMux()
)

func New(port string) Adapter {
	SetPort(port)
	return Adapter{port: port}
}

func SetPort(port string) {
	grpcPort = port
	GrpcServerEndpoint = "localhost:" + port
}
func GetPort() string {
	return grpcPort
}

type appInterface interface {
	api.App
	Register(grpcServer *grpc.Server, mux *runtime.ServeMux) error
}

// Adapter struct to save apps that implement grpc adapters and set port that grpc server should run
type Adapter struct {
	port string
	apps []appInterface
}

func (g *Adapter) Add(app appInterface) {
	g.apps = append(g.apps, app)
}

// Run method that implements Adapter interface
func (g Adapter) Run() {
	var err error
	listen, err := net.Listen("tcp", fmt.Sprintf(":%s", g.port))
	if err != nil {
		log.Fatalf("failed to listen on %s: %v", g.port, err)
	}
	grpcServer := grpc.NewServer()

	// register all apps
	for _, app := range g.apps {
		app.Register(grpcServer, mux)
	}

	reflection.Register(grpcServer)
	log.Infof("GRPC server listening on %s\n", g.port)
	go generatePBFiles()
	InstallProtocPlugins()
	go runWebServer()
	if err := grpcServer.Serve(listen); err != nil {
		log.Fatalf("failed to serve gRPC over %s: %v", g.port, err)
	}
}

// runWebServer run webserver based on all protobuffers
func runWebServer() error {
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()
	if err := generateSwagger(); err != nil {
		log.WithError(err).Errorf("error while generate swagger")
	}
	addr := ":8081"
	// Start HTTP server (and proxy calls to gRPC server endpoint)
	log.Infof("Webserver listening on in prefix /api %s", addr)
	return http.ListenAndServe(addr, mux)
}