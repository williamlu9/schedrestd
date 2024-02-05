/*
 * Schedrestd Rest API
 *
 * REST API to access Schedrestd
 *
 * API version: 1
 * Contact: williamlu9@gmail.com
 */

package file

import (
	"schedrestd/common"
	"schedrestd/common/logger"
	"schedrestd/common/response"
	"github.com/gin-gonic/gin"
	"net/http"
	"os/user"
	"os"
	"path"
	"errors"
	"strconv"
)

// Handler ...
type Handler struct {
}

// NewHostHandler ...
func NewFileHandler() *Handler {
	return &Handler{}
}

// @Summary   Upload a file to user home directory
// @Tags      file
// @Accept    multipart/form-data
// @Produce   json
// @Param     Authorization  header  string  true  "Token with Bearer started" default(Bearer <Add token here>)
// @Param     file  formData  file  true  "upload file"
// @Success   200 {object} FileResp  "Success"
// @Failure   400 {object} response.Response  "Bad request"
// @Failure   401 {object} response.Response  "Unauthorized user"
// @Failure   500 {object} response.Response  "Internal error"
// @Router    /file/upload [post]
// @Description Example upload request: 
// @Description curl -H "Authorization: Bearer $TOKEN" -H "Content-Type: multipart/form-data" -F "file=@/shared/testfile" "http://localhost:8088/sa/v1/file/upload"
func (h *Handler) UploadFile(c *gin.Context) {
	// src file
	file, err := c.FormFile("file")
	if err != nil {
		response.ResErr(c, http.StatusBadRequest, err)
		return			
	}
	logger.GetDefault().Infof("upload file: %v", file.Filename)
			
	// dst is user homedir
	val,_ := c.Get(common.UserHeader)
	username := val.(string)	
	user,err := user.Lookup(username)
	if err != nil {
		response.ResErr(c, http.StatusInternalServerError, err)
		return		
	}

	uid, _ := strconv.Atoi(user.Uid)
        gid, _ := strconv.Atoi(user.Gid)
	workdir := c.Query("dir")
	if len(workdir) == 0 {
		// destination is user's home
		if len(user.HomeDir) == 0 {
			response.ResErr(c, http.StatusInternalServerError, errors.New("User homedir not found"))
			return
		}
		workdir = user.HomeDir
	} else {
		if _, err := os.Stat(workdir); errors.Is(err, os.ErrNotExist) {
			err := os.MkdirAll(workdir, 0755)
			if err != nil {
				response.ResErr(c, http.StatusInternalServerError, err)
                                return
                        }
                        err = os.Chown(workdir, uid, gid)
                        if err != nil {
                                response.ResErr(c, http.StatusInternalServerError, err)
                                return
                        }
                }
	}

	dst := path.Join(workdir, file.Filename)
	logger.GetDefault().Infof("upload file to: %v", dst)
	
	// save file
	if err := c.SaveUploadedFile(file, dst); err != nil {
		response.ResErr(c, http.StatusInternalServerError, err)
		return		
	}
	
	// change file owner
	if err := os.Chown(dst, uid, gid); err != nil {
		response.ResErr(c, http.StatusInternalServerError, err)
		return			
	}

	resp := FileResponse {
		Path: dst,
	}
	response.ResOKGinJson(c, &FileResp{
		resp,
	})	
}

// @Summary   Download a file from user home directory
// @Tags      file
// @Produce   application/octet-stream
// @Param     Authorization  header  string  true  "Token with Bearer started" default(Bearer <Add token here>)
// @Param     file_name  path  string  true  "file_name"
// @Success   200 {object} file "Success"
// @Failure   400 {object} response.Response  "Bad request"
// @Failure   401 {object} response.Response  "Unauthorized user"
// @Failure   500 {object} response.Response  "Internal error"
// @Router    /file/download/{file_name} [get]
// @Description Example upload request: 
// @Description curl -X GET -H "Authorization: Bearer $TOKEN" "http://localhost:8088/sa/v1/file/download/testfile" > ./testfile
func (h *Handler) DownloadFile(c *gin.Context) {
	filename := c.Param("file_name")
	if len(filename) == 0 {
		response.ResErr(c, http.StatusBadRequest, errors.New("file_name is required"))
		return
	}	
	
	val,_ := c.Get(common.UserHeader)
	username := val.(string)	
	user,err := user.Lookup(username)
	if err != nil {
		response.ResErr(c, http.StatusInternalServerError, err)
		return		
	}

        workdir := c.Query("dir")
        if len(workdir) == 0 {
		if len(user.HomeDir) == 0 {
			response.ResErr(c, http.StatusInternalServerError, errors.New("User homedir not found"))
			return
		}
		workdir = user.HomeDir
	}
	srcPath := path.Join(workdir, filename)
	logger.GetDefault().Infof("download file: %v", srcPath)

    //c.File(srcPath)
    c.FileAttachment(srcPath, filename)
}
