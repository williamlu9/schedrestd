package middleware

import (
	"schedrestd/common"
	"schedrestd/common/jwt"
	"schedrestd/common/kvdb"
	"schedrestd/common/logger"
	"schedrestd/config"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
	"strings"
	"time"
)

// JWTAuth implements jwt middleware for gin
func JWTAuth(jwt *jwt.JWT, db *kvdb.KVStore, conf *config.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		if strings.Contains(c.Request.RequestURI, "login") ||
				 strings.Contains(c.Request.RequestURI, "swagger"){
			c.Next()
			return
		}
		
		ginLogger := logger.GetLogger(c)

		// Check token in header
		token := c.Request.Header.Get(common.TokenHeader)
		if token == "" {
			c.JSONP(http.StatusForbidden, gin.H{
				"msg": common.NoTokenMsg,
			})
			ginLogger.Info("Request denied without token")

			c.Abort()
			return
		}

		ginLogger.Infof("Start to check token: %v", token)

		tokenSlice := strings.SplitN(token, " ", 2)
		if !(len(tokenSlice) == 2 && tokenSlice[0] == "Bearer") {
			c.JSON(http.StatusBadRequest, gin.H{
				"msg": "Token is invalid.",
			})
			c.Abort()
			return
		}

		// Refresh token variable
		token = tokenSlice[1]

		// Parse the token
		claims, err := jwt.ParseToken(token)
		if err != nil {
			if !strings.Contains(err.Error(), "expired") {
				c.JSON(http.StatusBadRequest, gin.H{
					"msg": "Request failed. " + err.Error(),
				})
				c.Abort()
				return
			}
		}
		
		// Get the user name from parsed claim
		userName := claims.Name
		// Set the user name into header
		c.Set(common.UserHeader, userName)

		// Refresh the token time
		var timeOut int64
		timeOut = 30
		if conf.Timeout != "" {
			timeOut1, err := strconv.ParseInt(conf.Timeout, 10, 64)
			if err != nil {
				logger.GetDefault().Warnf("The timeout value in schedrestd.yaml is invalid: %v", err.Error())
			} else {
				timeOut = timeOut1
			}
		}
		timeOut = timeOut * 60
		key := fmt.Sprintf("%v-%v", token, userName)
		recordTimeByte, err := db.Get(common.BoltDBJWTTable, key)
		if err != nil {
			logger.GetDefault().Warnf("Failed to get the current token record from bolt db: %v", err.Error())
			c.JSON(http.StatusBadRequest, gin.H{
				"msg": "Request failed. Token is invalid",
			})
			c.Abort()
			return
		}
		recordTime, err := strconv.ParseInt(string(recordTimeByte), 10, 64)
		if err != nil {
			logger.GetDefault().Warnf("Stored time is invalid: %v", err.Error())
			c.JSON(http.StatusBadRequest, gin.H{
				"msg": "Request failed. Token is invalid",
			})
			c.Abort()
			return
		}

		currentTime := time.Now().Unix()
		if currentTime-recordTime <= timeOut {
			// No expired
			newTime := fmt.Sprintf("%d", currentTime)
			db.Put(common.BoltDBJWTTable, key, []byte(newTime))
		} else {
			// Expired
			c.JSON(http.StatusForbidden, gin.H{
				"msg": common.ExpiredMsg,
			})
			c.Abort()
			return
		}

		c.Next()
	}
}
