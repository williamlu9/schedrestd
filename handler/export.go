package handler

import (
	"schedrestd/handler/auth"
	"schedrestd/handler/cmd"
	"schedrestd/handler/file"
	"go.uber.org/fx"
)

var Module = fx.Options(
	fx.Provide(NewHandler),
	fx.Provide(auth.NewAuthHandler),
	fx.Provide(cmd.NewCmdHandler),
	fx.Provide(file.NewFileHandler),
)
