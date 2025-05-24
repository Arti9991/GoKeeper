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

// GetAddr получение исходного URL по укороченному
func (s *Server) DeleteData(ctx context.Context, in *pb.DeleteDataRequest) (*pb.DeleteDataResponce, error) {
	var res pb.DeleteDataResponce
	//var err error
	UserInfo := ctx.Value(interceptors.CtxKey).(servermodels.UserInfo)

	if !UserInfo.Register {
		return &res, status.Errorf(codes.Aborted, `Пользователь не авторизован`)
	}

	err := s.InfoStor.DeleteData(UserInfo.UserID, in.StorageID)
	if err != nil {
		logger.Log.Error("Error in delete datainfo from DB", zap.Error(err))
		return &res, status.Error(codes.Aborted, `Ошибка в удалении информации о данных`)
	}

	err = s.BinStorFunc.RemoveBinData(UserInfo.UserID, in.StorageID)
	if err != nil {
		logger.Log.Error("Error in remove data from binary storage", zap.Error(err))
		return &res, status.Error(codes.Aborted, `Ошибка в удалении бинарных данных`)
	}

	return &res, nil
}
