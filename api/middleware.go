package api

import (
	"errors"
	"moneytor/token"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

const authPayloadKey = "auth_payload"

func authMiddleware(tokenMaker token.Maker) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		authHeader := ctx.GetHeader("Authorization")
		if len(authHeader) == 0 {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, errResponse(errors.New("authorization header is missing")))
			return
		}

		fields := strings.Fields(authHeader)
		if len(fields) != 2 || !strings.EqualFold(fields[0], "bearer") {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, errResponse(errors.New("invalid authorization header format")))
			return
		}

		payload, err := tokenMaker.VerifyToken(fields[1])
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, errResponse(err))
			return
		}

		ctx.Set(authPayloadKey, payload)
		ctx.Next()
	}
}
