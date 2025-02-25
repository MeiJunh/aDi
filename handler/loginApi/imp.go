package loginApi

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
)

// LoginHandlerImp handler模版 -- 需要进行cookie校验
type LoginHandlerImp struct {
}

// Hello Example handler functions
func (l *LoginHandlerImp) Hello(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, gin.H{"message": "Hello from GetHello"})
}

// GreetTest post 测试
func (l *LoginHandlerImp) GreetTest(ctx *gin.Context) {
	name := ctx.PostForm("name")
	ctx.JSON(http.StatusOK, gin.H{"message": fmt.Sprintf("Hello, %s", name)})
}
