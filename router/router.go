package router

import (
	"schedrestd/common/jwt"
	"schedrestd/common/kvdb"
	"schedrestd/common/logger"
	"schedrestd/config"
	"schedrestd/handler"
	"schedrestd/handler/auth"
	"schedrestd/handler/cmd"
	"schedrestd/handler/file"
	"schedrestd/router/middleware"
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"go.uber.org/fx"
	"net"
	"net/http"
	"os"
	"strconv"
	"strings"
	docs "schedrestd/docs"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// AllHandler ...
// Add all handlers here for REST API
type AllHandler struct {
	fx.In
	IndexHandlers     *handler.Handler
	AuthHandlers      *auth.Handler
	CmdHandlers       *cmd.Handler
	FileHandlers      *file.Handler
}

// ErrorChan ...
type ErrorChan chan error

// NewRouter defines the routers of REST
func NewRouter(conf *config.Config, jwt *jwt.JWT, db *kvdb.KVStore, handler AllHandler) *gin.Engine {
	log := logger.GetDefault()

	gin.SetMode(gin.ReleaseMode)
	g := gin.New()
	g.Use(middleware.LoggerToFile(log))
	g.Use(middleware.JWTAuth(jwt, db, conf))
	g.Use(gin.Recovery())

	// Set the log writer for gin
	gin.DefaultWriter = &logger.GinLoggerWriter{
		Logger: log,
		IsErr:  false,
	}
	gin.DefaultErrorWriter = &logger.GinLoggerWriter{
		Logger: log,
		IsErr:  true,
	}

	registerAPI(g, conf, log, handler)

	return g
}

// NewStarter starts the http server with the port
func NewStarter(lifecycle fx.Lifecycle, engine *gin.Engine, config *config.Config) ErrorChan {
	errorChan := make(chan error, 1)
	saLogger := logger.GetDefault()

	lifecycle.Append(
		fx.Hook{
			OnStart: func(ctx context.Context) error {
				ssl := false
				if config.Ssl == "1" {
					if _, err := os.Stat(config.Cert); os.IsNotExist(err) {
						saLogger.Panicf("Cert file %s doe not exist.", config.Cert)
						panic("Cert file not exist")
					}
					if _, err := os.Stat(config.Key); os.IsNotExist(err) {
						saLogger.Panicf("Key file %s doe not exist.", config.Key)
						panic("Key file not exist")
					}
					ssl = true
				}

				var port string
				if ssl {
					port = config.HttpsPort
					if port == "" {
						port = "8043"
					}
				} else {
					port = config.HttpPort
					if port == "" {
						port = "8088"
					}
				}
				_, err := strconv.ParseInt(port, 10, 64)
				if err != nil {
					saLogger.Panicf("Port %v is invalid.", port)
					panic("Port is invalid")
				}

				port = fmt.Sprintf(":%v", port)
				listener, err := net.Listen("tcp", port)
				if err != nil {
					close(errorChan)
					return err
				}

				// Start the http server
				go func() {
					defer close(errorChan)

					if ssl {
						err = http.ServeTLS(listener, engine, config.Cert, config.Key)
					} else {
						err = http.Serve(listener, engine)
					}
					if err != nil {
						errorChan <- err
					}
				}()

				return nil
			},
		},
	)
	return errorChan
}

func registerAPI(g *gin.Engine, conf *config.Config, log logger.AipLogger, handler AllHandler) {
	if !strings.HasPrefix(conf.WebUrlPath, "/") {
		log.Panicf("The web url path is wrong: %v", conf.WebUrlPath)
		panic("Wrong web_url_path config.")
	}

	root := g.Group(conf.WebUrlPath).Group("/sa/v1")
	
	// Index
	root.GET("/", handler.IndexHandlers.Index)
	
	//doc
	docs.SwaggerInfo.BasePath = "/sa/v1"
	root.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))	

	// Auth login
	root.POST("/login", handler.AuthHandlers.LoginUser)

        // Token renew
        tokenR := root.Group("/token")
        tokenR.GET("/renew", handler.AuthHandlers.TokenRenew)

	// Cmd
	cmdR := root.Group("/cmd")
	cmdR.POST("/run", handler.CmdHandlers.RunCommand)

	// File
	fileR := root.Group("/file")
	fileR.GET("/download/:file_name", handler.FileHandlers.DownloadFile)
	fileR.POST("/upload", handler.FileHandlers.UploadFile)	
}

var Module = fx.Options(
	fx.Provide(NewRouter),
	fx.Provide(NewStarter),
)
