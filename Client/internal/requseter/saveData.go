package requseter

import (
	"bufio"
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/Arti9991/GoKeeper/client/internal/clientmodels"
	"github.com/Arti9991/GoKeeper/client/internal/inputfunc"
	pb "github.com/Arti9991/GoKeeper/client/internal/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

// SaveDataRequest метод для сохранения данных
func (req *ReqStruct) SaveDataRequest(Type string, offlineMode bool) error {
	var Metainfo string

	// парсим данные на входе
	data, err := inputfunc.ParceInput(Type)
	if err != nil {
		return err
	}
	// просим ввести метаинформацию для данных
	fmt.Printf("Введите метаинформацию:")

	reader := bufio.NewReader(os.Stdin)
	Metainfo, _ = reader.ReadString('\n')
	strings.TrimSuffix(Metainfo, "\n")

	// создаем StorageID для данных
	hashData := sha256.Sum256(data)
	StorageID := hex.EncodeToString(hashData[:])
	// сохраняем данные в локальное бинарное хранилище
	err = req.BinStor.SaveBinData(StorageID, data)
	if err != nil {
		return err
	}
	// записываем время сохранения
	CurrTime := time.Now().UTC().Format(time.RFC850)
	// создаем структуру с информацией о данных
	DtInf := clientmodels.NewerData{
		StorageID: StorageID,
		DataType:  Type,
		MetaInfo:  Metainfo,
		SaveTime:  CurrTime,
		Data:      data,
	}
	// сохраняем ифнормацию в локальную таблицу
	err = req.DBStor.SaveNew(StorageID, DtInf)
	if err != nil {
		return err
	}
	// если офлайн режим отключен, отправляем данные на сервер
	if !offlineMode {
		err = SendWithUpdate(StorageID, DtInf, req, data)
		if err != nil {
			fmt.Printf("\nNew data saved locally! StorageID for new data: %s\n", StorageID)
			return err
		}
	}
	// возращаем сообщение, что данные успешно сохранены с таким StorageID
	fmt.Printf("\nNew data saved! StorageID for new data: %s\n", StorageID)
	return nil
}

// SendWithUpdate функция отправки данных, с их обновлением, если на сервере более свежие данные
func SendWithUpdate(StorageID string, DtInf clientmodels.NewerData, req *ReqStruct, data []byte) error {
	// отправляем данные на сервер
	NewDt, err := SaveSend(DtInf, req, data)
	if err == nil {
		// если ошибок нет, то ставим, что данные синхронизированы
		err2 := req.DBStor.MarkDone(StorageID)
		if err2 != nil {
			return err2
		}
		return nil
	} else if err == clientmodels.ErrNewerData {
		// если ошибка, что данные на сервере новее
		// то сохраняем обновленные данные на клиент
		DtInf2 := clientmodels.NewerData{
			StorageID: NewDt.StorageID,
			DataType:  NewDt.DataType,
			MetaInfo:  NewDt.MetaInfo,
			SaveTime:  NewDt.SaveTime,
		}
		// обновляем информацию о данных
		err2 := req.DBStor.UpdateInfoNewer(NewDt.StorageID, DtInf2)
		if err2 != nil {
			return err2
		}
		// обновляем сами данные, полученные в ответе от сервера
		err2 = req.BinStor.UpdateBinData(NewDt.StorageID, NewDt.Data)
		if err2 != nil {
			return err2
		}
		return nil
	} else {
		return err
	}
}

// SaveSend функция отправки данных на сервер
func SaveSend(JrInf clientmodels.NewerData, req *ReqStruct, data []byte) (clientmodels.NewerData, error) {
	var UserID string
	var NewerData clientmodels.NewerData

	// считываем токен с UserID
	file, err := os.Open(clientmodels.TokenFile)
	if err != nil {
		return NewerData, err
	}
	defer file.Close()
	reader := bufio.NewReader(file)
	// Считываем строку текста
	UserID, err = reader.ReadString('\n')
	if err != nil {
		return NewerData, err
	}
	// убираем лишний суффикс
	UserID = strings.TrimSuffix(UserID, "\n")

	// добавляем UserID к метаданным запроса
	var header metadata.MD
	md := metadata.New(map[string]string{"UserID": UserID})
	ctx := metadata.NewOutgoingContext(context.Background(), md)
	// создаем клиент для отправки
	dial, err := grpc.NewClient(req.ServAddr, grpc.WithTransportCredentials(req.Creds))
	if err != nil {
		return NewerData, err
	}
	// вызываем метод отправки данных
	r := pb.NewKeeperClient(dial)
	ans, err := r.SaveData(ctx, &pb.SaveDataRequest{
		StorageID: JrInf.StorageID,
		Metainfo:  JrInf.MetaInfo,
		DataType:  JrInf.DataType,
		Time:      JrInf.SaveTime,
		Data:      data,
	}, grpc.Header(&header))
	if err != nil {
		return NewerData, err
	}
	if ans.IsOutdated {
		// если стоит флаг, что данные на сервере новее
		// пишем их выходные параметры с соответствующей
		// ошибкой
		NewerData.Data = ans.ReverseData.Data
		NewerData.DataType = ans.ReverseData.DataType
		NewerData.MetaInfo = ans.ReverseData.Metainfo
		NewerData.SaveTime = ans.ReverseData.Time
		NewerData.StorageID = ans.StorageID
		return NewerData, clientmodels.ErrNewerData
	}
	return NewerData, nil
}
