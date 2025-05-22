package server

import (
	"context"
	"time"

	"github.com/Arti9991/GoKeeper/server/internal/logger"
	"github.com/Arti9991/GoKeeper/server/internal/server/interceptors"
	pb "github.com/Arti9991/GoKeeper/server/internal/server/proto"
	"github.com/Arti9991/GoKeeper/server/internal/server/servermodels"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// GetAddr получение исходного URL по укороченному
func (s *Server) GiveData(ctx context.Context, in *pb.GiveDataRequest) (*pb.GiveDataResponce, error) {
	var res pb.GiveDataResponce
	//var err error
	UserInfo := ctx.Value(interceptors.CtxKey).(servermodels.UserInfo)

	if !UserInfo.Register {
		return &res, status.Errorf(codes.Aborted, `Пользователь не авторизован`)
	}

	getData, err := s.InfoStor.GetData(UserInfo.UserID, in.StorageID)
	if err != nil {
		logger.Log.Error("Error in get datainfo from DB", zap.Error(err))
		return &res, status.Error(codes.Aborted, `Ошибка в получении информации о данных`)
	}

	binData, err := s.BinStorFunc.GetBinData(UserInfo.UserID, in.StorageID)
	if err != nil {
		logger.Log.Error("Error in get data from bin storage", zap.Error(err))
		return &res, status.Error(codes.Aborted, `Ошибка в получении данных из хранилища`)
	}

	res.Data = binData
	res.DataType = getData.Type
	res.Metainfo = getData.MetaInfo
	res.Time = getData.SaveTime.Format(time.RFC850)

	return &res, nil
}
