package server

import (
	"context"
	"time"

	"github.com/Arti9991/GoKeeper/server/internal/coder"
	"github.com/Arti9991/GoKeeper/server/internal/logger"
	"github.com/Arti9991/GoKeeper/server/internal/server/interceptors"
	pb "github.com/Arti9991/GoKeeper/server/internal/server/proto"
	"github.com/Arti9991/GoKeeper/server/internal/server/servermodels"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// GiveData функция для получения данных, хранящихся на сервере
func (s *Server) GiveData(ctx context.Context, in *pb.GiveDataRequest) (*pb.GiveDataResponce, error) {
	var res pb.GiveDataResponce
	// получение UserID из контекста с интерцептора
	UserInfo := ctx.Value(interceptors.CtxKey).(servermodels.UserInfo)
	// если пользователь не авторизован, сообщаем ему об этом
	if !UserInfo.Register {
		return &res, status.Errorf(codes.Aborted, `Пользователь не авторизован`)
	}

	// получаем информацию о данных из базы
	getData, err := s.InfoStor.GetData(UserInfo.UserID, in.StorageID)
	if err != nil {
		logger.Log.Error("Error in get datainfo from DB", zap.Error(err))
		return &res, status.Error(codes.Aborted, `Ошибка в получении информации о данных`)
	}
	// получаем сами данные в бинарном формате
	binData, err := s.BinStorFunc.GetBinData(UserInfo.UserID, in.StorageID)
	if err != nil {
		logger.Log.Error("Error in get data from bin storage", zap.Error(err))
		return &res, status.Error(codes.Aborted, `Ошибка в получении данных из хранилища`)
	}

	// декодируем данные
	decData, err := coder.Derypt(binData)
	if err != nil {
		logger.Log.Error("Error in decoding data.", zap.Error(err))
		return &res, status.Error(codes.Aborted, `Ошибка в декодировании данных на сервере`)
	}

	// записываем ответ с полученными данными
	res.Data = decData
	res.DataType = getData.Type
	res.Metainfo = getData.MetaInfo
	res.Time = getData.SaveTime.Format(time.RFC850)

	return &res, nil
}
