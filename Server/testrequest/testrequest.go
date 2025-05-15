package main

import (
	// ...

	"context"
	"fmt"
	"log"

	pb "github.com/Arti9991/GoKeeper/server/internal/server/proto"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {

	// устанавливаем соединение с сервером
	conn, err := grpc.NewClient(":8080", grpc.WithTransportCredentials(insecure.NewCredentials()))
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
	// набор тестовых данных
	// for _, user := range users {
	// md := metadata.New(map[string]string{"UserID": "8c537969b84ad4eb0a73e29b3f2a9030"})
	ctx := context.Background()

	//var header metadata.MD
	// добавляем пользователей
	_, err := c.SaveData(ctx, &pb.SaveDataRequset{
		Id:       "1",
		Data:     []byte("string"),
		Metainfo: "This is meta info",
	})
	if err != nil {
		fmt.Println(err)
	}
}
