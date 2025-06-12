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

// GiveDataList функция для выдачи списка с информацией о всех данных пользователя на сервере
func (s *Server) GiveDataList(ctx context.Context, in *pb.GiveDataListRequest) (*pb.GiveDataListResponce, error) {
	var res pb.GiveDataListResponce
	// получение UserID из контекста с интерцептора
	UserInfo := ctx.Value(interceptors.CtxKey).(servermodels.UserInfo)
	// если пользователь не авторизован, сообщаем ему об этом
	if !UserInfo.Register {
		return &res, status.Errorf(codes.Aborted, `Пользователь не авторизован`)
	}

	// получаем массив структур с ифнормацией о данных пользователя
	getData, err := s.InfoStor.GetDataList(UserInfo.UserID)
	if err != nil {
		logger.Log.Error("Error in get datainfo from DB", zap.Error(err))
		return &res, status.Error(codes.Aborted, `Ошибка в получении информации о данных`)
	}

	// записываем всю информацию в ответ
	for _, dataLine := range getData {
		var resLine pb.GiveDataListResponce_DataList
		resLine.StorageID = dataLine.StorageID
		resLine.DataType = dataLine.Type
		resLine.Metainfo = dataLine.MetaInfo
		resLine.Time = dataLine.SaveTime.Format(time.RFC850)
		res.DataList = append(res.DataList, &resLine)
	}

	return &res, nil
}
