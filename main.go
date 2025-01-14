//
// Copyright © 2022 Schedrestd
//
package main

import (
	"schedrestd/common/cron"
	"schedrestd/common/jwt"
	"schedrestd/common/kvdb"
	"schedrestd/common/logger"
	"schedrestd/config"
	"schedrestd/handler"
	"schedrestd/router"
	"schedrestd/service"
	"context"
	"go.uber.org/fx"
)

// @title Schedrestd Rest API
// @version v1.0
// @description REST API to access Schedrestd

// @contact.name Schedrestd Support
// @contact.url http://teraproc.com/
// @contact.email williamlu9@gmail.com

// @BasePath /sa/v1

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @tokenUrl  https://example.com/sa/v1/login
func main() {
	var errC router.ErrorChan
	defaultLogger := logger.GetDefault()

	var serviceObj service.Service
	app := fx.New(
		config.Module,
		kvdb.Module,
		cron.Module,
		service.Module,
		jwt.Module,
		handler.Module,
		router.Module,
		fx.Populate(&errC),
		fx.Populate(&serviceObj),
	)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	if err := serviceObj.Initial(ctx); err != nil {
		defaultLogger.Fatalf("Failed to initialize the services: %v", err.Error())
		return
	}

	err := app.Start(context.Background())
	if err != nil {
		defaultLogger.Fatalf("Start server error: " + err.Error())
		return
	}

	select {
	case signal := <- app.Done():
		defaultLogger.Infof("Server stopped with signal %v", signal)
	case err := <-errC:
		defaultLogger.Fatalf("server stopped with error %v", err.Error())
	}
}
