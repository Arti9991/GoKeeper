package interceptors

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/Arti9991/GoKeeper/server/internal/logger"
	"github.com/Arti9991/GoKeeper/server/internal/server/servermodels"
	"github.com/golang-jwt/jwt/v4"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

// Claims — структура утверждений, которая включает стандартные утверждения и
// одно пользовательское UserID
type Claims struct {
	jwt.RegisteredClaims
	UserID string
}

const TOKENEXP = time.Hour * 24
const SECRETKEY = "supersecretkey"

type KeyContext string

var CtxKey = KeyContext("UserID")

// перехватчик для получения информации об авторизации пользователя из метаданных
func AtuhInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {

	var UserID string
	var err error
	UserExist := true

	md, ok := metadata.FromIncomingContext(ctx)
	if ok {
		values := md.Get("UserID")
		if len(values) > 0 {
			UserIDJWT := values[0]

			fmt.Println(len(UserIDJWT))
			if len(UserIDJWT) != 143 {
				UserExist = false
			} else {
				UserID, err = GetUserID(UserIDJWT)
				if err != nil {
					UserExist = false
				} else {
					newCtx := context.WithValue(ctx, CtxKey, servermodels.UserInfo{UserID: UserID, Register: UserExist})
					return handler(newCtx, req)
				}
				//fmt.Println("User ID un interceptor is:", UserID)
			}
		} else if len(values) == 0 {
			UserExist = false
		}
	} else {
		UserExist = false
	}
	fmt.Println(info.FullMethod)
	fmt.Println(UserExist)
	if !UserExist && !IsLogRegMethod(info.FullMethod) {
		logger.Log.Info("Bad user token")
		return nil, status.Errorf(codes.PermissionDenied, `Данный пользователь не авторизован!`)
	}

	//newCtx := context.WithValue(ctx, CtxKey, servermodels.UserInfo{UserID: UserID, Register: UserExist})
	return handler(ctx, req)
}

func GetUserID(tokenString string) (string, error) {
	claims := &Claims{}
	token, err := jwt.ParseWithClaims(tokenString, claims,
		func(t *jwt.Token) (interface{}, error) {
			return []byte(SECRETKEY), nil
		})
	if err != nil {
		return "", err
	}

	if !token.Valid {
		return "", servermodels.ErrorBadToken
	}

	return claims.UserID, nil
}

func IsLogRegMethod(method string) bool {
	if strings.Contains(method, "Register") || strings.Contains(method, "Login") {
		return true
	} else {
		return false
	}
}
