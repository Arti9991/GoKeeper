package server

import (
	"context"
	"fmt"

	"github.com/Arti9991/GoKeeper/server/internal/logger"
	pb "github.com/Arti9991/GoKeeper/server/internal/server/proto"
	"github.com/Arti9991/GoKeeper/server/internal/server/servermodels"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// GetAddr получение исходного URL по укороченному
func (s *Server) Loginuser(ctx context.Context, in *pb.LoginRequest) (*pb.LoginResponce, error) {
	var res pb.LoginResponce

	// UserInfo := ctx.Value(interceptors.CtxKey).(servermodels.UserInfo)
	// UserExist := UserInfo.Register

	// if UserExist {
	// 	return &res, status.Errorf(codes.Aborted, `Пользователь уже авторизован`)
	// }

	fmt.Println(in.UserLogin)
	fmt.Println(in.UserPassword)

	codedPassw := servermodels.CodePassword(in.UserPassword)

	fmt.Println(codedPassw)

	//UserID := rand.Text()[0:16]

	//fmt.Println(UserID)

	UserID, basePassw, err := s.UserStor.GetUser(in.UserLogin)
	if err == servermodels.ErrorNoSuchUser {
		return &res, status.Error(codes.PermissionDenied, `Неверное имя пользователя или пароль`)
	} else if err != nil {
		logger.Log.Error("Error in get user from users DB", zap.Error(err))
		return &res, status.Error(codes.Unavailable, `Ошибка в получении пользователя`)
	}
	if basePassw != codedPassw {
		return &res, status.Error(codes.PermissionDenied, `Неверное имя пользователя или пароль`)
	}

	JWTstr, err := BuildJWTString(UserID)
	if err != nil {
		return &res, status.Error(codes.Unavailable, `Ошибка в создании токена`)
	}

	res.UserID = JWTstr

	return &res, nil
}
