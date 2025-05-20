package server

import (
	"context"
	"fmt"
	"time"

	"github.com/Arti9991/GoKeeper/server/internal/server/interceptors"
	pb "github.com/Arti9991/GoKeeper/server/internal/server/proto"
	"github.com/Arti9991/GoKeeper/server/internal/server/servermodels"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// GetAddr получение исходного URL по укороченному
func (s *Server) SaveData(ctx context.Context, in *pb.SaveDataRequest) (*pb.SaveDataResponse, error) {
	var res pb.SaveDataResponse

	UserInfo := ctx.Value(interceptors.CtxKey).(servermodels.UserInfo)

	if !UserInfo.Register {
		return &res, status.Errorf(codes.Aborted, `Пользователь не авторизован`)
	}

	CurrTime := time.Now()

	var SaveDataInfo servermodels.SaveDataInfo

	SaveDataInfo.MetaInfo = "meta info here as we see NUMBER 2"
	SaveDataInfo.StorageID = "test"
	SaveDataInfo.Type = "TEXT"
	SaveDataInfo.SaveTime = CurrTime

	getData, err := s.DBData.SaveNewData(UserInfo.UserID, SaveDataInfo)
	if err != nil {
		if err == servermodels.ErrNewerData {
			fmt.Println(getData)
		} else {
			return &res, status.Error(codes.Aborted, `Ошибка в сохранении информации о данных`)
		}
	}

	fmt.Println(UserInfo.UserID)
	fmt.Println("Input ID", in.Id)
	fmt.Println("Input Data", in.Metainfo)

	return &res, nil
}
