package requseter

import (
	"bufio"
	"bytes"
	"context"
	"encoding/gob"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/Arti9991/GoKeeper/client/internal/clientmodels"
	"github.com/Arti9991/GoKeeper/client/internal/inputfunc"
	"github.com/Arti9991/GoKeeper/client/internal/journal"
	pb "github.com/Arti9991/GoKeeper/client/internal/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"
)

func TestSaveDataRequest(Type string) error {
	//offlineMode := false

	var UserID string
	var Metainfo string

	data, err := inputfunc.ParceInput(Type)
	if err != nil {
		return err
	}
	fmt.Printf("Введите метаинформацию:")

	// открываем потоковое чтение из консоли
	reader := bufio.NewReader(os.Stdin)
	// читаем строку из консоли
	Metainfo, _ = reader.ReadString('\n')

	fmt.Println("Open token")
	file, err := os.Open(clientmodels.TokenFile)
	if err != nil {
		fmt.Println(err)
		//logger.Log.Error("SAVE Error in opening file", zap.Error(err))
		return err
	}
	reader = bufio.NewReader(file)
	// Считываем строку текста
	UserID, err = reader.ReadString('\n')
	if err != nil {
		fmt.Println(err)
		return err
	}
	//Выводим строку
	UserID = strings.TrimSuffix(UserID, "\n")
	fmt.Printf("%#v", UserID)

	// if !offlineMode {
	// 	SyncRequest(offlineMode)
	// }
	var header metadata.MD
	md := metadata.New(map[string]string{"UserID": UserID})
	ctx := metadata.NewOutgoingContext(context.Background(), md)

	dial, err := grpc.NewClient(":8082", grpc.WithTransportCredentials(insecure.NewCredentials())) //":8082"
	if err != nil {
		log.Fatal(err)
	}
	CurrTime := time.Now().Format(time.RFC850)
	req := pb.NewKeeperClient(dial)
	ans, err := req.SaveData(ctx, &pb.SaveDataRequest{
		Metainfo: Metainfo,
		DataType: Type,
		Time:     CurrTime,
		Data:     data,
	}, grpc.Header(&header))
	if err != nil {
		fmt.Println(err)
		return err
	}

	JrInf := clientmodels.JournalInfo{
		Opperation: "SAVE",
		StorageID:  ans.StorageID,
		DataType:   Type,
		MetaInfo:   Metainfo,
		SaveTime:   CurrTime,
	}

	err = journal.JournalSave(JrInf)

	fmt.Println(ans.StorageID)

	return nil
}

func TestGetDataRequest(StorageID string) error {
	var UserID string

	fmt.Println("Open token")
	file, err := os.Open("./Token.txt")
	if err != nil {
		fmt.Println(err)
		//logger.Log.Error("SAVE Error in opening file", zap.Error(err))
		return err
	}
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

	dial, err := grpc.NewClient(":8082", grpc.WithTransportCredentials(insecure.NewCredentials())) //":8082"
	if err != nil {
		log.Fatal(err)
	}

	req := pb.NewKeeperClient(dial)
	ans, err := req.GiveData(ctx, &pb.GiveDataRequest{
		StorageID: StorageID,
	}, grpc.Header(&header))
	if err != nil {
		fmt.Println(err)
		return err
	}

	switch ans.DataType {
	case "CARD":

		fmt.Println(ans.Data)
		//fmt.Println(buff2)
		fmt.Println(ans.Metainfo)

		out := clientmodels.CardInfo{}
		dec := gob.NewDecoder(bytes.NewBuffer(ans.Data))
		dec.Decode(&out)
		fmt.Println(out.Number)
		fmt.Println(out.CVVcode)
		fmt.Println(out.ExpDate)
		fmt.Println(out.Holder)
	}
	return nil
}
