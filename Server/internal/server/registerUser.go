package server

import (
	"context"
	"crypto/rand"
	"fmt"
	"time"

	"github.com/Arti9991/GoKeeper/server/internal/server/interceptors"
	pb "github.com/Arti9991/GoKeeper/server/internal/server/proto"
	"github.com/Arti9991/GoKeeper/server/internal/server/servermodels"
	"github.com/golang-jwt/jwt/v4"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// GetAddr получение исходного URL по укороченному
func (s *Server) RegisterUser(ctx context.Context, in *pb.RegisterRequest) (*pb.RegisterResponce, error) {
	var res pb.RegisterResponce

	// UserInfo := ctx.Value(interceptors.CtxKey).(servermodels.UserInfo)
	// UserExist := UserInfo.Register

	// if UserExist {
	// 	return &res, status.Errorf(codes.Aborted, `Пользователь уже авторизован`)
	// }

	fmt.Println(in.UserLogin)
	fmt.Println(in.UserPassword)

	codedPassw := servermodels.CodePassword(in.UserPassword)

	fmt.Println(codedPassw)

	UserID := rand.Text()[0:16]

	fmt.Println(UserID)

	err := s.UserStor.SaveNewUser(UserID, in.UserLogin, codedPassw)
	if err != nil {
		if err == servermodels.ErrorUserAlready {
			return &res, status.Error(codes.Unavailable, `Пользователь уже зарегистрирован`)
		} else {
			return &res, status.Error(codes.Unavailable, `Ошибка в сохранении пользователя`)
		}
	}

	JWTstr, err := BuildJWTString(UserID)
	if err != nil {
		return &res, status.Error(codes.Unavailable, `Ошибка в создании токена`)
	}

	res.UserID = JWTstr

	return &res, nil
}

// BuildJWTString создаёт токен и возвращает его в виде строки.
func BuildJWTString(UserID string) (string, error) {
	// создаём новый токен с алгоритмом подписи HS256 и утверждениями — Claims
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, interceptors.Claims{
		RegisteredClaims: jwt.RegisteredClaims{
			// когда создан токен
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(interceptors.TOKENEXP)),
		},
		// собственное утверждение
		UserID: UserID,
	})

	// создаём строку токена
	tokenString, err := token.SignedString([]byte(interceptors.SECRETKEY))
	if err != nil {
		return "", err
	}

	// возвращаем строку токена
	return tokenString, nil
}
