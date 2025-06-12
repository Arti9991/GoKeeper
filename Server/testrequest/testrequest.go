package main

import (
	// ...

	"context"
	"fmt"
	"log"
	"time"

	pb "github.com/Arti9991/GoKeeper/server/internal/server/proto"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"
)

func main() {

	// устанавливаем соединение с сервером
	conn, err := grpc.NewClient(":8082", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()
	// получаем переменную интерфейсного типа UsersClient,
	// через которую будем отправлять сообщения
	c := pb.NewKeeperClient(conn)

	// функция, в которой будем отправлять сообщения
	BaseTestKeeper(c)
	// TestKeeperJSON(c)
}

// BaseTestKeeper функция для простейших POST и GET запросов
func BaseTestKeeper(c pb.KeeperClient) {
	var err error
	// набор тестовых данных
	// for _, user := range users {
	// md := metadata.New(map[string]string{"UserID": "8c537969b84ad4eb0a73e29b3f2a9030"})
	ctx := context.Background()

	respReg, err := c.Loginuser(ctx, &pb.LoginRequest{
		UserLogin:    "TestUserSECOND",
		UserPassword: "123456789",
	})
	if err != nil {
		log.Fatal(err)
	}

	// respReg, err := c.RegisterUser(ctx, &pb.RegisterRequest{
	// 	UserLogin:    "TestUserSECOND",
	// 	UserPassword: "123456789",
	// })
	// if err != nil {
	// 	log.Fatal(err)
	// }

	fmt.Println(respReg.UserID)
	var header metadata.MD
	md := metadata.New(map[string]string{"UserID": respReg.UserID})
	ctx2 := metadata.NewOutgoingContext(context.Background(), md)

	CurrTime := time.Now().Format(time.RFC850)
	savedData, err := c.SaveData(ctx2, &pb.SaveDataRequest{
		Data:     []byte("Hello there!"),
		DataType: "TEXT",
		Metainfo: "METAINFO",
		Time:     CurrTime,
	}, grpc.Header(&header))
	if err != nil {
		fmt.Println(err)
	}

	recievedData, err := c.GiveData(ctx2, &pb.GiveDataRequest{
		StorageID: savedData.StorageID,
	}, grpc.Header(&header))
	if err != nil {
		fmt.Println(err)
	}

	fmt.Printf("Recieved data: Type: %s, MetaInfo: %s, time: %s\n, Data: %s\n",
		recievedData.DataType,
		recievedData.Metainfo,
		recievedData.Time,
		string(recievedData.Data))

	_, err = c.UpdateData(ctx2, &pb.UpdateDataRequest{
		StorageID: savedData.StorageID,
		Data:      []byte("Hello there! UPDATE Not there... HERE!!!"),
		DataType:  "TEXT",
		Metainfo:  "second METAINFO updated",
		Time:      CurrTime,
	}, grpc.Header(&header))
	if err != nil {
		fmt.Println(err)
	}

	dataList, err := c.GiveDataList(ctx2, &pb.GiveDataListRequest{}, grpc.Header(&header))
	if err != nil {
		fmt.Println(err)
	}

	_, err = c.DeleteData(ctx2, &pb.DeleteDataRequest{
		StorageID: savedData.StorageID,
	}, grpc.Header(&header))
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println(dataList)
}
