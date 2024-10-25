/*
 * Schedrestd Rest API
 *
 * REST API to access Schedrestd
 *
 * API version: 1
 * Contact: williamlu9@gmail.com
 */

package auth

/*
#cgo CFLAGS: -std=c99
#cgo LDFLAGS: -lpam
#include <security/pam_appl.h>
#include <stdlib.h>

static struct pam_response *reply = NULL;

int
null_conv(int num_msg, const struct pam_message **msg,
           struct pam_response **resp, void *appdata_ptr)
{
    *resp = reply;
    return PAM_SUCCESS;
}

static struct pam_conv conv = {null_conv, NULL};

// return = true: authenticated
//          false: not authenticated
int
authenticate(char *user, char *password)
{
    pam_handle_t *pamh = NULL;
    int retval;

    if ((retval = pam_start("system-auth", user, &conv, &pamh)) != PAM_SUCCESS)
        return 0;
    reply = calloc(1, sizeof(struct pam_response));
    reply->resp = password;
    retval = pam_authenticate(pamh, 0);
    pam_end(pamh, PAM_SUCCESS);
    return (retval == PAM_SUCCESS ? 1 : 0);
}

*/
import "C"
import (
	"schedrestd/common"
	"schedrestd/common/jwt"
	"schedrestd/common/kvdb"
	"schedrestd/common/logger"
	"schedrestd/common/response"
	"schedrestd/config"
	"fmt"
	stdJwt "github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
	"time"
	"errors"
	"unsafe"
)

const (
	AUTHERROR = "Invalid username/password supplied"
)

// Handler ...
type Handler struct {
	jwt  *jwt.JWT
	db   *kvdb.KVStore
	conf *config.Config
}

// NewAuthHandler ...
func NewAuthHandler(jwt *jwt.JWT, db *kvdb.KVStore, conf *config.Config) *Handler {
	return &Handler{
		jwt:  jwt,
		db:   db,
		conf: conf,
	}
}

// @login
// @Description Logs user into the system
// @Tags    auth
// @Param   data  body  AuthReq  true  "Authenticate request"
// @Success 200 {object} TokenResp "Success"
// @Failure 400 {object} response.Response "Invalid username/password supplied"
// @Router /login [post] {}
func (h *Handler) LoginUser(c *gin.Context) {
	var req AuthReq
	var err error
	if err = c.BindJSON(&req); err != nil {
		response.ResErr(c, http.StatusBadRequest, err)
		return
	}

	username := C.CString(req.UserName)
	password := C.CString(req.Password)
	cAuth := C.authenticate(username, password)
	defer C.free(unsafe.Pointer(username))
	defer C.free(unsafe.Pointer(password))
	if int(cAuth) == 0 {
		logger.GetDefault().Errorf("Authentication failed for user: " + req.UserName)
		response.ResErr(c, http.StatusBadRequest, errors.New(AUTHERROR))
		return
	}

	timeOut := int64(req.Duration)
	if timeOut == 0 {
		timeOut, err = strconv.ParseInt(h.conf.Timeout, 10, 64)
                if err != nil {
			timeOut = 30 // Default value is 30 minutes
		}
	}
	// generate token
	claims := jwt.AClaims{
		req.UserName,
		stdJwt.StandardClaims{
			Issuer: "Schedrestd",
			IssuedAt: time.Now().Unix(),
		},
	}
	token, err := h.jwt.GenerateToken(claims)
	if err != nil {
		response.ResErr(c, http.StatusInternalServerError, err)
		return
	}

	// store token in db
	key := fmt.Sprintf("%v-%v", token, req.UserName)
	expireTime := time.Now().Add(time.Minute * time.Duration(timeOut)).Unix()
	if err = h.db.Put(common.BoltDBJWTTable,
		key,
		[]byte(strconv.FormatInt(expireTime, 10))); err != nil {
		response.ResErr(c, http.StatusInternalServerError, err)
		return
	}

	// send response
	resp := Token{
		Token:    token,
		UserName: req.UserName,
	}
	response.ResOKGinJson(c, &TokenResp{
		resp,
	})			
}

func (h *Handler) TokenRenew(c *gin.Context) {
    val,_ := c.Get(common.UserHeader)
    username := val.(string)
    claims := jwt.AClaims{
            username,
            stdJwt.StandardClaims{
                        Issuer: "Schedrestd",
                        IssuedAt: time.Now().Unix(),
                },
        }
        token, err := h.jwt.GenerateToken(claims)
        if err != nil {
                response.ResErr(c, http.StatusInternalServerError, err)
                return
        }

        // store token in db
        key := fmt.Sprintf("%v-%v", token, username)
        timeOut, err := strconv.ParseInt(h.conf.Timeout, 10, 64)
        if err != nil {
                timeOut = 30 // Default value is 30 min
        }
        expireTime := time.Now().Add(time.Minute * time.Duration(timeOut)).Unix()
        if err = h.db.Put(common.BoltDBJWTTable,
                key,
                []byte(strconv.FormatInt(expireTime, 10))); err != nil {
                response.ResErr(c, http.StatusInternalServerError, err)
                return
        }

        // send response
        resp := Token{
                Token:    token,
                UserName: username,
        }
        response.ResOKGinJson(c, &TokenResp{
                resp,
        })
}
