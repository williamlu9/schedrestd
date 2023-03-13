package cmd

import (
	"schedrestd/common"
	"schedrestd/common/response"
	"errors"
	"github.com/gin-gonic/gin"
	"net/http"
	"os/exec"
	"os/user"
	"strconv"
	"syscall"
	"strings"
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
		cmdStr = cmdStr + strings.Join(cmdReq.Envs[:], " ") + " "
	}
	cmdStr = cmdStr + cmdReq.Command
	cmd := exec.Command("sh", "-c", cmdStr)

	// submit job as authenticated user
	val,_ := c.Get(common.UserHeader)
	username := val.(string)
	user, err := user.Lookup(username)
	if err != nil {
		response.ResErr(c, http.StatusInternalServerError, err)
		return
	}
	uid, _ := strconv.Atoi(user.Uid)
	gid, _ := strconv.Atoi(user.Gid)
	cmd.SysProcAttr = &syscall.SysProcAttr{}
	cmd.SysProcAttr.Credential =
		&syscall.Credential{Uid: uint32(uid), Gid: uint32(gid)}
	output, err := cmd.CombinedOutput()
	if err != nil {
		response.ResErr(c, http.StatusBadRequest, err)
		return
	}
	
	res.Output = string(output)
	response.ResOKGinJson(c, res) 
}
