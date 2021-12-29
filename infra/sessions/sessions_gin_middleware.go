package sessions

import (
	"skysight/infra/fail"

	"github.com/gin-gonic/gin"
)

func SessionFilter() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		token, err := ctx.Cookie(KeySecToken)
		if err != nil {
			panic(fail.ErrUnauthenticated)
		}
		securityContextValue, found := TokenCache.Get(token)
		if !found {
			panic(fail.ErrUnauthenticated)
		}
		secCtx, ok := securityContextValue.(*Session)
		if !ok {
			panic(fail.ErrUnauthenticated)
		}
		InjectSessionIntoGinContext(ctx, secCtx)
		ctx.Next()
	}
}
