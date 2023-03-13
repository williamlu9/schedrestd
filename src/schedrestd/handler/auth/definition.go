/*
 * Schedrestd Rest API
 *
 * REST API to access Schedrestd
 *
 * API version: 1
 * Contact: williamlu9@gmail.com
 */

package auth

type Token struct {

	// User token used to be authenticated
	Token string `json:"token,omitempty"`

	// User name
	UserName string `json:"userName,omitempty"`
}

type TokenResp struct {
	Token Token `json:"token"`
}

type AuthReq struct {

	// User name
	UserName string `json:"username" binding:"required"`

	// Password
	Password string `json:"password" binding:"required"`
}
