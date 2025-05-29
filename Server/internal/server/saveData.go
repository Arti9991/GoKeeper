package server

import (
	"context"
	"fmt"
	"strings"
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

// GetAddr получение исходного URL по укороченному
func (s *Server) SaveData(ctx context.Context, in *pb.SaveDataRequest) (*pb.SaveDataResponse, error) {
	var res pb.SaveDataResponse
	res.ReverseData = new(pb.SaveDataResponse_ReverseData)

	var err error
	UserInfo := ctx.Value(interceptors.CtxKey).(servermodels.UserInfo)

	if !UserInfo.Register {
		return &res, status.Errorf(codes.Aborted, `Пользователь не авторизован`)
	}

	StorageID := in.StorageID

	fmt.Println(in.Data)
	encData, err := coder.Encrypt(in.Data)
	if err != nil {
		logger.Log.Error("Error in encoding data.", zap.Error(err))
		return &res, status.Error(codes.Aborted, `Ошибка в кодировании данных на сервере`)
	}

	fmt.Println(encData)
	// заполняем структуру для сохранени данных
	var SaveDataInfo servermodels.SaveDataInfo
	SaveDataInfo.Data = encData
	SaveDataInfo.StorageID = in.StorageID
	SaveDataInfo.Type = in.DataType
	SaveDataInfo.MetaInfo = in.Metainfo
	SaveDataInfo.SaveTime, err = time.Parse(time.RFC850, in.Time)
	if err != nil {
		logger.Log.Error("Error in parse time from request setting own time", zap.Error(err))
		SaveDataInfo.SaveTime = time.Now()
	}

	// сохраняем новые данные
	getData, err := s.InfoStor.SaveNewData(UserInfo.UserID, SaveDataInfo)
	if err != nil {
		// если возвращена ошибка, что на сервере данные свежее (по времени)
		if err == servermodels.ErrNewerData {
			fmt.Println(getData)
			// выставляем ответный флаг что на сервере данные свежее
			res.IsOutdated = true
			// и получаем обновленные данные из бинарного харнилища
			getUpdateData, err2 := s.BinStorFunc.GetBinData(UserInfo.UserID, getData.StorageID)
			// если файл в бинарном хранилище отсутствует (по каким-либо причинам)
			if err2 != nil && strings.Contains(err2.Error(), "no such file") {
				// то возвращаем полученные данные
				getUpdateData = in.Data
				// и сохраняем их на в бинарное хранилище
				err3 := s.BinStorFunc.SaveBinData(UserInfo.UserID, getData.StorageID, in.Data)
				if err3 != nil {
					return &res, status.Error(codes.Aborted, `Ошибка в получении обновленных бинарных данных`)
				}
			} else if err2 != nil {
				logger.Log.Error("Error in getting newer binary data", zap.Error(err))
				return &res, status.Error(codes.Aborted, `Ошибка в получении обновленных бинарных данных`)
			}
			// записываем ответную структуру
			res.ReverseData.Data = getUpdateData
			res.ReverseData.DataType = getData.Type
			res.ReverseData.Metainfo = getData.MetaInfo
			res.ReverseData.Time = getData.SaveTime.Format(time.RFC850)
			res.StorageID = getData.StorageID
			return &res, nil
		} else {
			return &res, status.Error(codes.Aborted, `Ошибка в сохранении информации о данных`)
		}
	}

	// если же изначально данных на сервере не было
	// то ставим флаг что пришедшие данные не устарели
	res.IsOutdated = false
	// и сохраняем данные в бинарное хранилище
	err2 := s.BinStorFunc.SaveBinData(UserInfo.UserID, StorageID, SaveDataInfo.Data)
	if err2 != nil {
		logger.Log.Error("Error in saving binary data", zap.Error(err2))
		return &res, status.Error(codes.Aborted, `Ошибка в сохранении бинарных данных`)
	}

	fmt.Println(UserInfo.UserID)
	fmt.Println("Input Data", in.Metainfo)

	res.StorageID = StorageID

	return &res, nil
}
