package server

import (
	"context"
	"fmt"
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
func (s *Server) UpdateData(ctx context.Context, in *pb.UpdateDataRequest) (*pb.UpdateDataResponse, error) {
	var res pb.UpdateDataResponse
	var err error
	UserInfo := ctx.Value(interceptors.CtxKey).(servermodels.UserInfo)

	if !UserInfo.Register {
		return &res, status.Errorf(codes.Aborted, `Пользователь не авторизован`)
	}

	var SaveDataInfo servermodels.SaveDataInfo
	SaveDataInfo.Data = in.Data
	SaveDataInfo.StorageID = in.StorageID
	SaveDataInfo.Type = in.DataType
	SaveDataInfo.MetaInfo = in.Metainfo
	SaveDataInfo.SaveTime, err = time.Parse(time.RFC850, in.Time)
	if err != nil {
		logger.Log.Error("Error in parse time from request setting own time", zap.Error(err))
		SaveDataInfo.SaveTime = time.Now()
	}

	err = s.InfoStor.UpdateData(UserInfo.UserID, SaveDataInfo)
	if err != nil {
		return &res, status.Error(codes.Aborted, `Ошибка в обновлении информации о данных`)
	}

	s.BinStorFunc.UpdateBinData(UserInfo.UserID, in.StorageID, SaveDataInfo.Data)
	if err != nil {
		logger.Log.Error("Error in updating binary data", zap.Error(err))
		return &res, status.Error(codes.Aborted, `Ошибка в обновлении бинарных данных`)
	}

	fmt.Println(UserInfo.UserID)
	fmt.Println("Input Data", in.Metainfo)

	return &res, nil
}
