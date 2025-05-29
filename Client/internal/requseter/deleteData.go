package requseter

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"os"
	"strings"

	pb "github.com/Arti9991/GoKeeper/client/internal/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"
)

func DeleteDataRequest(StorageID string, req *ReqStruct, offlineMode bool) error {

	var err error

	fmt.Println(offlineMode)

	if !offlineMode {
		err = DeleteDataOnline(StorageID, req.ServAddr)
		if err != nil {
			fmt.Println(err)
		}
	}

	err = req.DBStor.DeleteData(StorageID)
	if err != nil {
		fmt.Println(err)
		return err
	}
	err = req.BinStor.RemoveBinData(StorageID)
	if err != nil {
		fmt.Println(err)
		return err
	}
	return nil
}

func DeleteDataOnline(StorageID string, addr string) error {

	fmt.Println("Open token")
	file, err := os.Open("./Token.txt")
	if err != nil {
		fmt.Println(err)
		//logger.Log.Error("SAVE Error in opening file", zap.Error(err))
		return err
	}
	defer file.Close()

	reader := bufio.NewReader(file)
	// Считываем строку текста
	UserID, err := reader.ReadString('\n')
	if err != nil {
		fmt.Println(err)
		return err
	}
	//Выводим строку
	UserID = strings.TrimSuffix(UserID, "\n")
	fmt.Printf("%#v", UserID)

	var header metadata.MD
	md := metadata.New(map[string]string{"UserID": UserID})
	ctx := metadata.NewOutgoingContext(context.Background(), md)

	dial, err := grpc.NewClient(addr, grpc.WithTransportCredentials(insecure.NewCredentials())) //req.ServAddr
	if err != nil {
		log.Fatal(err)
	}

	r := pb.NewKeeperClient(dial)
	_, err = r.DeleteData(ctx, &pb.DeleteDataRequest{
		StorageID: StorageID,
	}, grpc.Header(&header))
	if err != nil {
		fmt.Println(err)
		return err
	}
	return nil
}
