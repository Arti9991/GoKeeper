package requseter

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/Arti9991/GoKeeper/client/internal/clientmodels"
	"github.com/Arti9991/GoKeeper/client/internal/inputfunc"
	pb "github.com/Arti9991/GoKeeper/client/internal/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

// UpdateDataRequest функция принудительного обновления данных
func (req *ReqStruct) UpdateDataRequest(StorageID string, Type string, offlineMode bool) error {
	var err error

	// парсим входные данные для обновления
	data, err := inputfunc.ParceInput(Type)
	if err != nil {
		return err
	}
	// просим ввести новую метаинформацию
	fmt.Printf("\nВведите метаинформацию:")
	// открываем потоковое чтение из консоли
	reader := bufio.NewReader(os.Stdin)
	// читаем строку из консоли
	Metainfo, _ := reader.ReadString('\n')
	strings.TrimSuffix(Metainfo, "\n")

	// сохраняем время обновления информации
	CurrTime := time.Now().UTC().Format(time.RFC850)
	// создаем структуру с новыми данными
	DtInf := clientmodels.NewerData{
		StorageID: StorageID,
		DataType:  Type,
		MetaInfo:  Metainfo,
		SaveTime:  CurrTime,
		Data:      data,
	}
	// обновляем локальные данные
	err = UpdateDataOfline(DtInf, req)
	if err != nil {
		return err
	}

	if !offlineMode {
		// если режим работы online явно не указан
		// обновляем данные на сервере
		err2 := UpdateDataOnline(DtInf, req)
		if err2 != nil {
			return err2
		}
		// если все успешно, ставим метку о свежести
		// локальных данных
		err2 = req.DBStor.MarkDone(DtInf.StorageID)
		if err2 != nil {
			return err2
		}
	}
	fmt.Printf("\nДанные обновлены!\n")
	return nil
}

// UpdateDataOnline функция обновления данных на сервере
func UpdateDataOnline(DtInf clientmodels.NewerData, req *ReqStruct) error {
	var UserID string

	// считваем токен
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
	// инициазируем клиент для запроса
	dial, err := grpc.NewClient(req.ServAddr, grpc.WithTransportCredentials(req.Creds)) //req.ServAddr
	if err != nil {
		log.Fatal(err)
	}
	// выполняем запрос на обновление данных
	r := pb.NewKeeperClient(dial)
	_, err = r.UpdateData(ctx, &pb.UpdateDataRequest{
		StorageID: DtInf.StorageID,
		Metainfo:  DtInf.MetaInfo,
		DataType:  DtInf.DataType,
		Time:      DtInf.SaveTime,
		Data:      DtInf.Data,
	}, grpc.Header(&header))
	if err != nil {
		return err
	}

	return nil
}

// UpdateDataOfline функция обновления данных на клиенте
func UpdateDataOfline(DtInf clientmodels.NewerData, req *ReqStruct) error {
	var err error
	// ообновляем инфомрацию о данных
	err = req.DBStor.UpdateInfo(DtInf.StorageID, DtInf)
	if err != nil {
		if err == clientmodels.ErrNoSuchRows {
			// если информации нет, сохраняем ее
			err2 := req.DBStor.SaveNew(DtInf.StorageID, DtInf)
			if err2 != nil {
				return err2
			}
		} else {
			return err
		}
	}
	// обновляем сами данные в хранилище
	err = req.BinStor.UpdateBinData(DtInf.StorageID, DtInf.Data)
	if err != nil {
		err2 := req.BinStor.SaveBinData(DtInf.StorageID, DtInf.Data)
		return err2
	}
	return nil
}
