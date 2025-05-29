package requseter

import (
	"bufio"
	"context"
	"crypto/sha256"
	"encoding/hex"
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

func SaveDataRequest(Type string, req *ReqStruct, offlineMode bool) error {
	var Metainfo string

	fmt.Println(offlineMode)

	data, err := inputfunc.ParceInput(Type)
	if err != nil {
		return err
	}
	fmt.Printf("Введите метаинформацию:")

	// открываем потоковое чтение из консоли
	reader := bufio.NewReader(os.Stdin)
	// читаем строку из консоли
	Metainfo, _ = reader.ReadString('\n')
	strings.TrimSuffix(Metainfo, "\n")

	hashData := sha256.Sum256(data)
	StorageID := hex.EncodeToString(hashData[:])

	err = req.BinStor.SaveBinData(StorageID, data)
	fmt.Println(err)

	CurrTime := time.Now().UTC().Format(time.RFC850)

	DtInf := clientmodels.NewerData{
		StorageID: StorageID,
		DataType:  Type,
		MetaInfo:  Metainfo,
		SaveTime:  CurrTime,
		Data:      data,
	}

	//err = journal.JournalSave(JrInf)

	// err = req.DBStor.ReinitTable()
	// fmt.Println(err)

	err = req.DBStor.SaveNew(StorageID, DtInf)
	fmt.Println(err)

	if !offlineMode {
		err = SendWithUpdate(StorageID, DtInf, req, data)
		if err != nil {
			fmt.Println(err)
			return err
		}
	}
	//fmt.Println(StorageID)

	return nil
}

func SendWithUpdate(StorageID string, DtInf clientmodels.NewerData, req *ReqStruct, data []byte) error {
	fmt.Println("SendWithUpdate")
	NewDt, err := SaveSend(DtInf, req, data)
	if err == nil {
		err2 := req.DBStor.MarkDone(StorageID)
		if err2 != nil {
			fmt.Println(err2)
		}
		return nil
	} else if err == clientmodels.ErrNewerData {
		//fmt.Println(NewDt)
		fmt.Println(clientmodels.ErrNewerData)
		DtInf2 := clientmodels.NewerData{
			StorageID: NewDt.StorageID,
			DataType:  NewDt.DataType,
			MetaInfo:  NewDt.MetaInfo,
			SaveTime:  NewDt.SaveTime,
		}
		err2 := req.DBStor.UpdateInfoNewer(NewDt.StorageID, DtInf2)
		if err2 != nil {
			fmt.Println(err2)
			return err2
		}
		fmt.Println("TO BIN DATA UPDATE")
		err2 = req.BinStor.UpdateBinData(NewDt.StorageID, NewDt.Data)
		if err2 != nil {
			fmt.Println(err2)
			return err2
		}
		return nil
	} else {
		fmt.Println(err)
		return err
	}
}
func SaveSend(JrInf clientmodels.NewerData, req *ReqStruct, data []byte) (clientmodels.NewerData, error) {
	var UserID string
	var NewerData clientmodels.NewerData

	fmt.Println("Open token")
	file, err := os.Open(clientmodels.TokenFile)
	if err != nil {
		fmt.Println(err)
		//logger.Log.Error("SAVE Error in opening file", zap.Error(err))
		return NewerData, err
	}
	reader := bufio.NewReader(file)
	// Считываем строку текста
	UserID, err = reader.ReadString('\n')
	if err != nil {
		fmt.Println(err)
		return NewerData, err
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
	ans, err := r.SaveData(ctx, &pb.SaveDataRequest{
		StorageID: JrInf.StorageID,
		Metainfo:  JrInf.MetaInfo,
		DataType:  JrInf.DataType,
		Time:      JrInf.SaveTime,
		Data:      data,
	}, grpc.Header(&header))
	if err != nil {
		fmt.Println(err)
		return NewerData, err
	}
	if ans.IsOutdated {
		fmt.Println("IS OUTDATED")
		NewerData.Data = ans.ReverseData.Data
		NewerData.DataType = ans.ReverseData.DataType
		NewerData.MetaInfo = ans.ReverseData.Metainfo
		NewerData.SaveTime = ans.ReverseData.Time
		NewerData.StorageID = ans.StorageID
		return NewerData, clientmodels.ErrNewerData
	}
	return NewerData, nil
}
