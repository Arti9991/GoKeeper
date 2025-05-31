package requseter

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/Arti9991/GoKeeper/client/internal/clientmodels"
	pb "github.com/Arti9991/GoKeeper/client/internal/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

// DeleteDataRequest метод удаления данных
func (req *ReqStruct) DeleteDataRequest(StorageID string, offlineMode bool) error {
	var err error
	// спрашиваем подтверждение на выход из аккаунта
	fmt.Println("Вы точно хотите удалить данные? [y/n] (В онлайн режиме данные на сервере будут также удалены!)")
	reader := bufio.NewReader(os.Stdin)
	// читаем путь из консоли
	agree, err := reader.ReadString('\n')
	agree = strings.ToLower(agree)
	if strings.Contains(agree, "n") {
		return clientmodels.ErrUserAbort
	}

	if !offlineMode {
		// если режим работы онлайн, то удаляем данные на сервере
		err = DeleteDataOnline(StorageID, req.ServAddr, req)
		if err != nil {
			fmt.Printf("\nНе удалось удалить данные на сервере: %s\n", err.Error())
		}
	}
	// удаляем информацию о данных из локальной базы
	err = req.DBStor.DeleteData(StorageID)
	if err != nil {
		return err
	}
	// удаляем сами данные из хранилища
	err = req.BinStor.RemoveBinData(StorageID)
	if err != nil {
		return err
	}
	return nil
}

// DeleteDataOnline функция удаления данных с сервера
func DeleteDataOnline(StorageID string, addr string, req *ReqStruct) error {

	// считываем токен из файла
	file, err := os.Open("./Token.txt")
	if err != nil {
		return err
	}
	defer file.Close()

	reader := bufio.NewReader(file)
	// Считываем строку текста
	UserID, err := reader.ReadString('\n')
	if err != nil {
		return err
	}
	// удаляем суффикс
	UserID = strings.TrimSuffix(UserID, "\n")

	// добавляем метаданные с UserID к запросу
	var header metadata.MD
	md := metadata.New(map[string]string{"UserID": UserID})
	ctx := metadata.NewOutgoingContext(context.Background(), md)
	// инциализируем клиент
	dial, err := grpc.NewClient(addr, grpc.WithTransportCredentials(req.Creds)) //req.ServAddr
	if err != nil {
		log.Fatal(err)
	}
	// вызываем метод удаления на сервере
	r := pb.NewKeeperClient(dial)
	_, err = r.DeleteData(ctx, &pb.DeleteDataRequest{
		StorageID: StorageID,
	}, grpc.Header(&header))
	if err != nil {
		return err
	}
	return nil
}
