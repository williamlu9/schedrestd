/*
 * Schedrestd Rest API
 *
 * REST API to access Schedrestd
 *
 * API version: 1
 * Contact: williamlu9@gmail.com
 */

package file

type FileResponse struct {
	Path string `json:"path"`
}

// FileResp ...
type FileResp struct {
	File FileResponse `json:"file"`
}
