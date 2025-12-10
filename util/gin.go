package util

import "github.com/gin-gonic/gin"

func QueryAll(ctx *gin.Context) map[string]string {
	query := ctx.Request.URL.Query()
	result := make(map[string]string)
	for key := range query {
		result[key] = query.Get(key)
	}
	return result
}
