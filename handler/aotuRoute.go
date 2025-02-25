package handler

import (
	"aDi/log"
	"github.com/gin-gonic/gin"
	"net/http"
	"reflect"
	"strings"
)

// CheckDynamicRouter 需要校验的route
func CheckDynamicRouter(controller interface{}) gin.HandlerFunc {
	controllerValue := reflect.ValueOf(controller)

	return func(ctx *gin.Context) {
		// Extract the path part from the request URL
		// Assuming the path is something like "/api/hello" -> "hello"
		pathParts := strings.Split(ctx.Request.URL.Path, "/")
		if len(pathParts) < 2 {
			log.Errorf("err path:%s", ctx.Request.URL.Path)
			ctx.JSON(http.StatusNotFound, gin.H{"error": "Invalid path"})
			return
		}

		methodName := strings.Title(pathParts[1])
		method := controllerValue.MethodByName(methodName)

		if method.IsValid() {
			// 进行token校验
			// Call the method with context as the parameter
			rspList := method.Call([]reflect.Value{reflect.ValueOf(ctx)})
			if len(rspList) > 0 {
				ctx.JSON(http.StatusOK, rspList[0].Interface())
			}
		} else {
			ctx.JSON(http.StatusNotFound, gin.H{"error": "Method not found"})
			log.Errorf("invalid path:%s", ctx.Request.URL.Path)
		}
		return
	}
}

// NoCheckDynamicRouter 需要校验的route
func NoCheckDynamicRouter(controller interface{}) gin.HandlerFunc {
	controllerValue := reflect.ValueOf(controller)

	return func(ctx *gin.Context) {
		// Extract the path part from the request URL
		// Assuming the path is something like "/api/hello" -> "hello"
		pathParts := strings.Split(strings.Trim(ctx.Request.URL.Path, "/"), "/")
		if len(pathParts) < 2 {
			log.Errorf("err path:%s", ctx.Request.URL.Path)
			ctx.JSON(http.StatusNotFound, gin.H{"error": "Invalid path"})
			return
		}

		methodName := pathParts[1]
		// methodName := strings.Title(pathParts[1])
		method := controllerValue.MethodByName(methodName)

		if method.IsValid() {
			// Call the method with context as the parameter
			rspList := method.Call([]reflect.Value{reflect.ValueOf(ctx)})
			if len(rspList) > 0 {
				ctx.JSON(http.StatusOK, rspList[0].Interface())
			}
		} else {
			log.Errorf("invalid path:%s", ctx.Request.URL.Path)
			ctx.JSON(http.StatusNotFound, gin.H{"error": "Method not found"})
		}
		return
	}
}
