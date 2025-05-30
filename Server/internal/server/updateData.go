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

// UpdateData функция для принудительного(!!!) обновления данных
func (s *Server) UpdateData(ctx context.Context, in *pb.UpdateDataRequest) (*pb.UpdateDataResponse, error) {
	// инициализация ответа
	var res pb.UpdateDataResponse
	var err error
	// получение UserID из контекста с интерцептора
	UserInfo := ctx.Value(interceptors.CtxKey).(servermodels.UserInfo)
	// если пользователь не авторизован, сообщаем ему об этом
	if !UserInfo.Register {
		return &res, status.Errorf(codes.Aborted, `Пользователь не авторизован`)
	}

	// кодируем полученные данные
	encData, err := coder.Encrypt(in.Data)
	if err != nil {
		logger.Log.Error("Error in encoding data.", zap.Error(err))
		return &res, status.Error(codes.Aborted, `Ошибка в кодировании данных на сервере`)
	}
	// заполняем структуру для сохранени данных
	var SaveDataInfo servermodels.SaveDataInfo
	SaveDataInfo.Data = encData
	SaveDataInfo.StorageID = in.StorageID
	SaveDataInfo.Type = in.DataType
	SaveDataInfo.MetaInfo = in.Metainfo
	SaveDataInfo.SaveTime, err = time.Parse(time.RFC850, in.Time)
	if err != nil {
		// если при парсинге времени ошибка, то ставим текущее
		logger.Log.Error("Error in parse time from request setting own time", zap.Error(err))
		SaveDataInfo.SaveTime = time.Now()
	}

	// обновляем информацию о данных
	err = s.InfoStor.UpdateData(UserInfo.UserID, SaveDataInfo)
	if err != nil {
		return &res, status.Error(codes.Aborted, `Ошибка в обновлении информации о данных`)
	}
	// обновляем сами данные
	s.BinStorFunc.UpdateBinData(UserInfo.UserID, in.StorageID, SaveDataInfo.Data)
	if err != nil {
		logger.Log.Error("Error in updating binary data", zap.Error(err))
		return &res, status.Error(codes.Aborted, `Ошибка в обновлении бинарных данных`)
	}

	return &res, nil
}
