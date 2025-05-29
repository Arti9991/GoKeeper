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
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"
)

func UpdateDataRequest(StorageID string, Type string, req *ReqStruct, offlineMode bool) error {

	var err error

	fmt.Println(offlineMode)

	data, err := inputfunc.ParceInput(Type)
	if err != nil {
		return err
	}

	fmt.Printf("\nВведите метаинформацию:")
	// открываем потоковое чтение из консоли
	reader := bufio.NewReader(os.Stdin)
	// читаем строку из консоли
	Metainfo, _ := reader.ReadString('\n')
	strings.TrimSuffix(Metainfo, "\n")

	CurrTime := time.Now().Format(time.RFC850)

	DtInf := clientmodels.NewerData{
		StorageID: StorageID,
		DataType:  Type,
		MetaInfo:  Metainfo,
		SaveTime:  CurrTime,
		Data:      data,
	}
	if !offlineMode {
		err = UpdateDataOnline(DtInf, req)
		if err != nil {
			fmt.Println(err)
		}
	}

	err = UpdateDataOfline(DtInf, req)
	if err != nil {
		fmt.Println(err)
		return err
	}
	return nil
}

func UpdateDataOnline(DtInf clientmodels.NewerData, req *ReqStruct) error {
	var UserID string

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

	dial, err := grpc.NewClient(req.ServAddr, grpc.WithTransportCredentials(insecure.NewCredentials())) //req.ServAddr
	if err != nil {
		log.Fatal(err)
	}

	r := pb.NewKeeperClient(dial)
	_, err = r.UpdateData(ctx, &pb.UpdateDataRequest{
		StorageID: DtInf.StorageID,
		Metainfo:  DtInf.MetaInfo,
		DataType:  DtInf.DataType,
		Time:      DtInf.SaveTime,
		Data:      DtInf.Data,
	}, grpc.Header(&header))
	if err != nil {
		fmt.Println(err)
		return err
	}

	return nil
}

func UpdateDataOfline(DtInf clientmodels.NewerData, req *ReqStruct) error {
	var err error

	err = req.DBStor.UpdateInfoNewer(DtInf.StorageID, DtInf)
	if err != nil {
		return err
	}

	err = req.BinStor.UpdateBinData(DtInf.StorageID, DtInf.Data)
	if err != nil {
		return err
	}
	return nil
}
