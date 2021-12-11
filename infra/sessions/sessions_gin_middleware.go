package sessions

import (
	"skysight/bizerror"

	"github.com/gin-gonic/gin"
)

func SessionFilter() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		token, err := ctx.Cookie(KeySecToken)
		if err != nil {
			panic(bizerror.ErrUnauthenticated)
		}
		securityContextValue, found := TokenCache.Get(token)
		if !found {
			panic(bizerror.ErrUnauthenticated)
		}
		secCtx, ok := securityContextValue.(*Session)
		if !ok {
			panic(bizerror.ErrUnauthenticated)
		}
		InjectSessionIntoGinContext(ctx, secCtx)
		ctx.Next()
	}
}
