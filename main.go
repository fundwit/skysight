package main

import (
	"fmt"
	"net/http"
	"skysight/bizerror"
	"skysight/infra/tracing"
	"skysight/localize"

	"github.com/gin-gonic/gin"
	"github.com/opentracing/opentracing-go"
	"github.com/sirupsen/logrus"
)

func main() {
	logrus.Infoln("service start")

	tracer, closer, err := tracing.NewTracer()
	if err != nil {
		logrus.Fatalf("failed to build tracer %v\n", err)
	}
	opentracing.SetGlobalTracer(tracer)
	defer closer.Close()

	engine := gin.Default()

	engine.Use(localize.LocalizeMiddleware("./i18n"))
	engine.Use(tracing.TracingIngress())
	engine.Use(bizerror.ErrorHandling())

	engine.GET("/", func(c *gin.Context) {
		fmt.Println(localize.MustGetMessage("status-running"))
		c.String(http.StatusOK, localize.MustGetMessage("status-running"))
	})

	err = engine.Run(":80")
	if err != nil {
		panic(err)
	}
}
