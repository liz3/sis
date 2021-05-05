package main

import (
	siss "sis/server"
	"github.com/kataras/iris/v12"
)

func main() {
	srv := siss.BuildServer()
	siss.BuildRoutes(srv)
	srv.Run(iris.Addr(":8080"), iris.WithoutServerError(iris.ErrServerClosed))
}
