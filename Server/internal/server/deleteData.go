package server

import (
	"context"

	"github.com/Arti9991/GoKeeper/server/internal/logger"
	"github.com/Arti9991/GoKeeper/server/internal/server/interceptors"
	pb "github.com/Arti9991/GoKeeper/server/internal/server/proto"
	"github.com/Arti9991/GoKeeper/server/internal/server/servermodels"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// DeleteData функция для принудительного полного удаления данных
func (s *Server) DeleteData(ctx context.Context, in *pb.DeleteDataRequest) (*pb.DeleteDataResponce, error) {
	var res pb.DeleteDataResponce
	// получение UserID из контекста с интерцептора
	UserInfo := ctx.Value(interceptors.CtxKey).(servermodels.UserInfo)
	// если пользователь не авторизован, сообщаем ему об этом
	if !UserInfo.Register {
		return &res, status.Errorf(codes.Aborted, `Пользователь не авторизован`)
	}

	// удаляем информацию о данных из базы
	err := s.InfoStor.DeleteData(UserInfo.UserID, in.StorageID)
	if err != nil {
		logger.Log.Error("Error in delete datainfo from DB", zap.Error(err))
		return &res, status.Error(codes.Aborted, `Ошибка в удалении информации о данных`)
	}

	// удаляем сами данные из хранилища
	err = s.BinStorFunc.RemoveBinData(UserInfo.UserID, in.StorageID)
	if err != nil {
		logger.Log.Error("Error in remove data from binary storage", zap.Error(err))
		return &res, status.Error(codes.Aborted, `Ошибка в удалении бинарных данных`)
	}

	return &res, nil
}
