package server

import (
	"context"

	"github.com/Arti9991/GoKeeper/server/internal/logger"
	pb "github.com/Arti9991/GoKeeper/server/internal/server/proto"
	"github.com/Arti9991/GoKeeper/server/internal/server/servermodels"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// Loginuser функция для входа в ститему зарегистрированного пользователя
func (s *Server) Loginuser(ctx context.Context, in *pb.LoginRequest) (*pb.LoginResponce, error) {
	var res pb.LoginResponce

	// кодируем полученный пароль
	codedPassw := servermodels.CodePassword(in.UserPassword)

	// получаем информацию о пользователе из базы данных
	UserID, basePassw, err := s.UserStor.GetUser(in.UserLogin)
	if err == servermodels.ErrorNoSuchUser {
		logger.Log.Error("No such user", zap.Error(err))
		return &res, status.Error(codes.PermissionDenied, `Неверное имя пользователя или пароль`)
	} else if err != nil {
		logger.Log.Error("Error in get user from users DB", zap.Error(err))
		return &res, status.Error(codes.Unavailable, `Ошибка в получении пользователя`)
	}

	// проверяем совпадает ли хэш пароля с полученным
	if basePassw != codedPassw {
		logger.Log.Error("Bad password", zap.Error(err))
		return &res, status.Error(codes.PermissionDenied, `Неверное имя пользователя или пароль`)
	}
	// создаем JWT токен
	JWTstr, err := BuildJWTString(UserID)
	if err != nil {
		return &res, status.Error(codes.Unavailable, `Ошибка в создании токена`)
	}
	// отправляем JWT токен в ответ
	res.UserID = JWTstr

	return &res, nil
}
