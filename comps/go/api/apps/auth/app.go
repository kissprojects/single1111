package auth

import (
	"github.com/kissprojects/single/comps/go/api/adapters/rest"
	"google.golang.org/grpc"
	"net/http"
)

var App AppModule

type AppModule struct{}

func (a AppModule) GetRouters() *[]rest.RouteGroup {
	return nil
}

func (a AppModule) GetRouterGroup() *[]rest.RouteGroup {
	return getRoutes()
}

func (a AppModule) GetMiddlewares() []func(http.Handler) http.Handler {
	return getMiddlewares()
}

func (a AppModule) Register(_ *grpc.Server) {

}

func (a AppModule) BeforeStart() {}

func (a AppModule) AfterStart() {
	AddPublicRouter(http.MethodPost, "/users/login")
	AddPublicRouter(http.MethodPost, "/users/register")
}
