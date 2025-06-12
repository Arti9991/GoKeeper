package interceptors

import (
	"context"
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

// вспомогательные переменные
// (т.к. проект учебный, то секретный ключ задаем здесь)
const TOKENEXP = time.Hour * 24
const SECRETKEY = "supersecretkey"

type KeyContext string

var CtxKey = KeyContext("UserID")

// AtuhInterceptor перехватчик для получения информации об авторизации пользователя из метаданных
func AtuhInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {

	var UserID string
	var err error
	UserExist := true
	// получаем метаданные из входящего контекста
	md, ok := metadata.FromIncomingContext(ctx)
	if ok {
		// проверяем наличие нужной строки
		values := md.Get("UserID")
		if len(values) > 0 {
			UserIDJWT := values[0]
			// если строка другой длинны, значит неверный токен
			if len(UserIDJWT) != 143 {
				UserExist = false
			} else {
				// пробуем получить токен
				UserID, err = GetUserID(UserIDJWT)
				if err != nil {
					UserExist = false
				} else {
					// записываем полученный UserID в контекст
					newCtx := context.WithValue(ctx, CtxKey, servermodels.UserInfo{UserID: UserID, Register: UserExist})
					return handler(newCtx, req)
				}
			}
		} else if len(values) == 0 {
			UserExist = false
		}
	} else {
		UserExist = false
	}

	// Если какая-то из проверок не прошла и вызванный метод не авторизация/регистрация
	// то ставим ответ, что пользователь не авторизован
	if !UserExist && !IsLogRegMethod(info.FullMethod) {
		logger.Log.Info("Bad user token")
		return nil, status.Errorf(codes.PermissionDenied, `Данный пользователь не авторизован!`)
	}

	return handler(ctx, req)
}

// GetUserID функция получения UserID из токена
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

// IsLogRegMethod функция проверки вызываемого метода
func IsLogRegMethod(method string) bool {
	if strings.Contains(method, "Register") || strings.Contains(method, "Login") {
		return true
	} else {
		return false
	}
}
