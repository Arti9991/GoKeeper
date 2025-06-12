package requseter

import (
	"bufio"
	"context"
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

// GetDataRequest метод для получения данных
func (req *ReqStruct) GetDataRequest(StorageID string, offlineMode bool) error {
	var ans clientmodels.NewerData
	var ansOn clientmodels.NewerData
	var ansOf clientmodels.NewerData
	var err error
	var err1 error
	var err2 error

	// получем локальные данные
	ansOf, err1 = GetDataOfline(StorageID, req)
	// парсим время локальных данных
	TimeLoc, err := time.Parse(time.RFC850, ansOf.SaveTime)
	if err == nil {
		TimeLoc = time.Now().UTC()
	}

	if !offlineMode {
		// если режим работы не офлайн, получаем данные с сервера
		ansOn, err2 = CompareGetData(StorageID, TimeLoc, req)
	}
	if offlineMode || err2 != nil {
		// если при запросе есть ошибка
		// или стоит флаг офлайн режима
		// выдаем локальные данные
		if err1 == nil {
			ans = ansOf
		} else {
			return err1
		}
	} else {
		// иначе выдаем онлайн данные
		ans = ansOn
	}
	// парсим полученные данные и выводим их на экран
	err = inputfunc.ParceAnswer(ans.Data, StorageID, ans.DataType, ans.MetaInfo)
	if err != nil {
		return err
	}
	return nil
}

// GetDataOnline функция получения данных с сервера
func GetDataOnline(StorageID string, addr string, req *ReqStruct) (clientmodels.NewerData, error) {
	var DataGet clientmodels.NewerData
	var UserID string
	// читаем токен из файла
	file, err := os.Open("./Token.txt")
	if err != nil {
		return DataGet, err
	}
	defer file.Close()

	reader := bufio.NewReader(file)
	// Считываем строку текста
	UserID, err = reader.ReadString('\n')
	if err != nil {
		return DataGet, err
	}
	//Выводим строку
	UserID = strings.TrimSuffix(UserID, "\n")

	// добавляем UserID в метаданные запроса
	var header metadata.MD
	md := metadata.New(map[string]string{"UserID": UserID})
	ctx := metadata.NewOutgoingContext(context.Background(), md)
	// инциализируем клиент запроса
	dial, err := grpc.NewClient(addr, grpc.WithTransportCredentials(req.Creds))
	if err != nil {
		log.Fatal(err)
	}
	// выполняем запрос на получение данных
	r := pb.NewKeeperClient(dial)
	ans, err := r.GiveData(ctx, &pb.GiveDataRequest{
		StorageID: StorageID,
	}, grpc.Header(&header))
	if err != nil {
		return DataGet, err
	}
	// выдаем полученные данные
	DataGet.Data = ans.Data
	DataGet.DataType = ans.DataType
	DataGet.MetaInfo = ans.Metainfo
	DataGet.StorageID = StorageID
	DataGet.SaveTime = ans.Time
	return DataGet, nil
}

// GetDataOfline функция получения данных, сохраненных локально
func GetDataOfline(StorageID string, req *ReqStruct) (clientmodels.NewerData, error) {
	var DataGet clientmodels.NewerData
	var err error
	DataGet.StorageID = StorageID

	// получем информацию о данных из базы
	DataGet, err = req.DBStor.Get(StorageID)
	if err != nil {
		return DataGet, err
	}
	// получаем сами данные из хранилища
	DataGet.Data, err = req.BinStor.GetBinData(StorageID)
	if err != nil {
		return DataGet, err
	}
	return DataGet, nil
}

// CompareGetData функция получения данных с сервера и сравнения их свежести с локальными данными
func CompareGetData(StorageID string, TimeLoc time.Time, req *ReqStruct) (clientmodels.NewerData, error) {
	var ansOn clientmodels.NewerData
	var err error

	// получаем данные с сервера
	ansOn, err = GetDataOnline(StorageID, req.ServAddr, req)
	if err != nil {
		return ansOn, err
	}
	// парсим время сохранения, полученнное с сервера
	timeServ, err := time.Parse(time.RFC850, ansOn.SaveTime)
	if err != nil {
		return ansOn, err
	}
	if timeServ.After(TimeLoc) {
		// если время сохранения данных на сервере
		// новее, ставим метку о необходимости
		// синхронизации данных
		err = req.DBStor.MarkUnDone(StorageID)
		if err != nil {
			return ansOn, err
		}
	}
	return ansOn, nil
}
