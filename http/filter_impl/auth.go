package filter_impl

import (
	"context"
	"errors"
	"fmt"
	"strconv"

	"github.com/glory-go/glory/grpc"
	ghttp "github.com/glory-go/glory/http"
	"github.com/glory-go/glory/log"
	"git.go-online.org.cn/GoOnline-IDE/grpc-files/auth-service"
	"github.com/dgrijalva/jwt-go"
)

/*
HttpTokenAuthMiddleware does token check and reject invalid request
Before using this middleware, pls make sure you have config like:

"auth-service":
    service_id: GoOnline-IDE-auth_service
    server_address: "dev.go-online.org.cn:31512"
    protocol: grpc

The configured client of auth-service in your glory.yaml

if success, this middleware will put userID that parsed from token into c.Ctx.
*/
func HttpTokenAuthMiddleware(c *ghttp.GRegisterController, f ghttp.HandleFunc) (err error) {
	// todo: new grpc once in http filter
	client := grpc.NewGrpcClient("auth-service")
	authClient := auth.NewAuthClient(client.GetConn())
	ctx := c.R.Context()
	key, err := authClient.GetJWTTokenKey(ctx, &auth.Empty{})
	if err != nil {
		log.CtxErrorf(ctx, "call auth.GetJWTTokenKey got error: %v", err)
		return err
	}

	token := c.R.Header.Get("Authorization")
	// 解析token
	data, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(key.JWTTokenKey), nil
	})
	if err != nil {
		log.CtxErrorf(ctx, "parse jwt meets error: %v", err)
		return err
	}
	// parse username from jwt
	if claims, ok := data.Claims.(jwt.MapClaims); ok && data.Valid {
		username := claims["sub"]
		if username == "" {
			log.CtxWarnf(ctx, "got invalid token contains empty username, token: %v", token)
			return errors.New("invalid token")
		}
		sub := username.(string)
		userID, err := strconv.ParseInt(sub, 10, 64)
		if err != nil {
			log.CtxErrorf(ctx, "%v is not valid int64 data", sub)
			return err
		}
		// if success, put userID parsed from token into c.Ctx
		c.Ctx = context.WithValue(c.Ctx, "userID", userID)
		return f(c)
	}
	log.CtxWarnf(ctx, "got invalid token, token: %v", token)
	return errors.New("invalid token")
}
