package cmd

import (
	"schedrestd/common"
	"schedrestd/common/response"
	"errors"
	"github.com/gin-gonic/gin"
	"net/http"
	"os/exec"
	"os/user"
	"bytes"
)

const (
	SIGKILL int = 9
	SIGSTOP int = 19
	SIGCONT int = 18
)

// Handler ...
type Handler struct {
}

// NewCmdHandler ...
func NewCmdHandler() *Handler {
	return &Handler{}
}

// @Summary   Run a command
// @Tags      cmd
// @Accept    json
// @Produce   json
// @Param     Authorization  header  string  true  "Token with Bearer started" default(Bearer <Add token here>)
// @Param     data  body  CmdRun  true  "Run Command"
// @Success   200 {object} RunCmdResp "Success"
// @Failure   400 {object} response.Response  "Bad request"
// @Failure   401 {object} response.Response  "Unauthorized user"
// @Failure   403 {object} response.Response  "Permission denied"
// @Failure   500 {object} response.Response  "Internal server error"
// @Router    /cmd/run [post]
func (h *Handler) RunCommand(c *gin.Context) {
	var cmdReq CmdRun
	var res RunCommandResponse
	var stdout, stderr bytes.Buffer
	if err := c.BindJSON(&cmdReq); err != nil {
		response.ResErr(c, http.StatusBadRequest, err)
		return
	}

	if len(cmdReq.Command) == 0 {
		response.ResErr(c, http.StatusBadRequest, errors.New("Command is required"))
		return
	}

	cmdStr := ""
	if len(cmdReq.Cwd) > 0 {
		cmdStr = cmdStr + "cd " + cmdReq.Cwd + ";"
	}
	if len(cmdReq.Envs) > 0 {
		for _, env := range cmdReq.Envs {
			cmdStr = cmdStr + "export " + env + ";"
		}
	}
	cmdStr = cmdStr + cmdReq.Command

	// submit job as authenticated user
	val,_ := c.Get(common.UserHeader)
	username := val.(string)
	_, err := user.Lookup(username)
	if err != nil {
		response.ResErr(c, http.StatusInternalServerError, err)
		return
	}
	cmd := exec.Command("su", "-s", "/bin/bash", "-", username, "-c", cmdStr)
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	err = cmd.Run()
	
	res.Output = string(stdout.Bytes())
	res.Error  = string(stderr.Bytes())
	res.ExitCode = 0
	if err != nil {
		if exitError, ok := err.(*exec.ExitError); ok {
			res.ExitCode = exitError.ExitCode()
		}
	}
	response.ResOKGinJson(c, res) 
}
