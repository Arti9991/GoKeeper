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

func ShowDataLoc(req *ReqStruct) error {
	err := req.DBStor.ShowTable()
	if err != nil {
		fmt.Println(err)
		return err
	}

	return nil
}

func ShowDataOn(req *ReqStruct, offlineMode bool) error {
	var err error
	var UserID string

	if offlineMode {
		return clientmodels.ErrNoOfflineList
	}

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
	UserID, err = reader.ReadString('\n')
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

	dial, err := grpc.NewClient(req.ServAddr, grpc.WithTransportCredentials(req.Creds)) //req.ServAddr
	if err != nil {
		log.Fatal(err)
	}

	r := pb.NewKeeperClient(dial)
	ans, err := r.GiveDataList(ctx, &pb.GiveDataListRequest{}, grpc.Header(&header))
	if err != nil {
		fmt.Println(err)
		return err
	}
	fmt.Printf("\nSaved data on server\n")
	fmt.Printf("%-5s %-64s %#-25v %-10s %-40s\n", "num", "StorageID", "MetaInfo", "Type", "SaveTime")
	for i, infoLine := range ans.DataList {
		fmt.Printf("%-5d %-64s %#-25v %-10s %-40s\n", i+1, infoLine.StorageID, infoLine.Metainfo, infoLine.DataType, infoLine.Time)
	}
	return nil
}
