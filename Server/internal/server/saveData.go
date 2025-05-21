package server

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
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
func (s *Server) SaveData(ctx context.Context, in *pb.SaveDataRequest) (*pb.SaveDataResponse, error) {
	var res pb.SaveDataResponse
	var err error
	UserInfo := ctx.Value(interceptors.CtxKey).(servermodels.UserInfo)

	if !UserInfo.Register {
		return &res, status.Errorf(codes.Aborted, `Пользователь не авторизован`)
	}

	hashData := sha256.Sum256(in.Data)
	StorageID := hex.EncodeToString(hashData[:])

	fmt.Println(StorageID)
	fmt.Println(len(StorageID))

	var SaveDataInfo servermodels.SaveDataInfo
	SaveDataInfo.Data = in.Data
	SaveDataInfo.StorageID = StorageID
	SaveDataInfo.Type = in.DataType
	SaveDataInfo.MetaInfo = in.Metainfo
	SaveDataInfo.SaveTime, err = time.Parse(time.RFC850, in.Time)
	if err != nil {
		logger.Log.Error("Error in parse time from request setting own time", zap.Error(err))
		SaveDataInfo.SaveTime = time.Now()
	}

	getData, err := s.InfoStor.SaveNewData(UserInfo.UserID, SaveDataInfo)
	if err != nil {
		if err == servermodels.ErrNewerData {
			fmt.Println(getData)
			SaveDataInfo.StorageID = getData.StorageID
		} else {
			return &res, status.Error(codes.Aborted, `Ошибка в сохранении информации о данных`)
		}
	}

	s.BinStor.SaveBinData(UserInfo.UserID, StorageID, SaveDataInfo.Data)

	fmt.Println(UserInfo.UserID)
	fmt.Println("Input Data", in.Metainfo)

	res.StorageID = StorageID

	return &res, nil
}
