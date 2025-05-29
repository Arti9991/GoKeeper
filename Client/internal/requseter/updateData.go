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

	CurrTime := time.Now().UTC().Format(time.RFC850)

	DtInf := clientmodels.NewerData{
		StorageID: StorageID,
		DataType:  Type,
		MetaInfo:  Metainfo,
		SaveTime:  CurrTime,
		Data:      data,
	}

	err = UpdateDataOfline(DtInf, req)
	if err != nil {
		fmt.Println(err)
		return err
	}

	if !offlineMode {
		err2 := UpdateDataOnline(DtInf, req)
		if err2 != nil {
			fmt.Println(err2)
			return err2
		}
		err2 = req.DBStor.MarkDone(DtInf.StorageID)
		if err2 != nil {
			fmt.Println(err2)
			return err2
		}
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

	dial, err := grpc.NewClient(req.ServAddr, grpc.WithTransportCredentials(req.Creds)) //req.ServAddr
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

	err = req.DBStor.UpdateInfo(DtInf.StorageID, DtInf)
	if err != nil {
		if err == clientmodels.ErrNoSuchRows {
			err2 := req.DBStor.SaveNew(DtInf.StorageID, DtInf)
			if err2 != nil {
				return err2
			}
		} else {
			return err
		}
	}

	err = req.BinStor.UpdateBinData(DtInf.StorageID, DtInf.Data)
	if err != nil {
		fmt.Println(err)
		err2 := req.BinStor.SaveBinData(DtInf.StorageID, DtInf.Data)
		return err2
	}
	return nil
}
