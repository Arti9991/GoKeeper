package main

import (
	// ...

	"context"
	"fmt"
	"log"

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
		UserLogin:    "TestUser",
		UserPassword: "123456",
	})
	if err != nil {
		log.Fatal(err)
	}

	// respReg, err := c.RegisterUser(ctx, &pb.RegisterRequest{
	// 	UserLogin:    "TestUser",
	// 	UserPassword: "123456",
	// })
	// if err != nil {
	// 	log.Fatal(err)
	// }

	fmt.Println(respReg.UserID)
	var header metadata.MD
	md := metadata.New(map[string]string{"UserID": respReg.UserID})
	ctx2 := metadata.NewOutgoingContext(context.Background(), md)
	// // добавляем пользователей
	_, err = c.SaveData(ctx2, &pb.SaveDataRequest{
		Id:       "1",
		Data:     []byte("string"),
		Metainfo: "This is meta info",
	}, grpc.Header(&header))
	if err != nil {
		fmt.Println(err)
	}
}
