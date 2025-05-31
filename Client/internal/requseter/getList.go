package requseter

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/Arti9991/GoKeeper/client/internal/clientmodels"
	pb "github.com/Arti9991/GoKeeper/client/internal/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

// ShowDataLoc метод отображения информации о локальных данных
func (req *ReqStruct) ShowDataLoc() error {
	err := req.DBStor.ShowTable()
	if err != nil {
		return err
	}

	return nil
}

// ShowDataOn метод отображения информации о данных на сервере
func (req *ReqStruct) ShowDataOn(offlineMode bool) error {
	var err error
	var UserID string

	if offlineMode {
		// в офлайн режиме метод недоступен
		return clientmodels.ErrNoOfflineList
	}

	// считываем токен
	file, err := os.Open("./Token.txt")
	if err != nil {
		return err
	}
	defer file.Close()

	reader := bufio.NewReader(file)
	// Считываем строку текста
	UserID, err = reader.ReadString('\n')
	if err != nil {
		return err
	}
	// удаляем лишний суффикс
	UserID = strings.TrimSuffix(UserID, "\n")
	// добавляем UserID к метаданным запроса
	var header metadata.MD
	md := metadata.New(map[string]string{"UserID": UserID})
	ctx := metadata.NewOutgoingContext(context.Background(), md)
	// инициализируем клиент
	dial, err := grpc.NewClient(req.ServAddr, grpc.WithTransportCredentials(req.Creds)) //req.ServAddr
	if err != nil {
		return err
	}
	// вызываем метод для получения списка всех сохраненных данных
	r := pb.NewKeeperClient(dial)
	ans, err := r.GiveDataList(ctx, &pb.GiveDataListRequest{}, grpc.Header(&header))
	if err != nil {
		return err
	}
	// выводим список с информацией о данных на экран
	fmt.Printf("\nSaved data on server\n")
	fmt.Printf("%-5s %-64s %#-25v %-10s %-40s\n", "num", "StorageID", "MetaInfo", "Type", "SaveTime")
	for i, infoLine := range ans.DataList {
		fmt.Printf("%-5d %-64s %#-25v %-10s %-40s\n", i+1, infoLine.StorageID, infoLine.Metainfo, infoLine.DataType, infoLine.Time)
	}
	return nil
}
