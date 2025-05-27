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
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"
)

func SendWithUpdate(StorageID string, JrInf clientmodels.JournalInfo, req *ReqStruct, data []byte) error {

	NewDt, err := SaveSend(JrInf, req, data)
	if err == nil {
		err2 := req.DBStor.MarkDone(StorageID)
		if err2 != nil {
			fmt.Println(err2)
		}
		return nil
	} else if err == clientmodels.ErrNewerData {
		fmt.Println(NewDt)
		JrInf2 := clientmodels.JournalInfo{
			Opperation: "SAVE",
			StorageID:  NewDt.StorageID,
			DataType:   NewDt.DataType,
			MetaInfo:   NewDt.MetaInfo,
			SaveTime:   NewDt.SaveTime,
			Sync:       true,
		}
		err2 := req.DBStor.UpdateInfoNewer(NewDt.StorageID, JrInf2)
		if err != nil {
			fmt.Println(err2)
			return err2
		}
		err2 = req.BinStor.SaveBinData(NewDt.StorageID, NewDt.Data)
		if err != nil {
			fmt.Println(err2)
			return err2
		}
		return nil
	} else {
		fmt.Println(err)
		return err
	}
}
func SaveSend(JrInf clientmodels.JournalInfo, req *ReqStruct, data []byte) (clientmodels.NewerData, error) {
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

	dial, err := grpc.NewClient(":8082", grpc.WithTransportCredentials(insecure.NewCredentials())) //":8082"
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
		NewerData.Data = ans.ReverseData.Data
		NewerData.DataType = ans.ReverseData.DataType
		NewerData.MetaInfo = ans.ReverseData.Metainfo
		NewerData.SaveTime = ans.ReverseData.Time
		NewerData.StorageID = ans.StorageID
		return NewerData, clientmodels.ErrNewerData
	}
	return NewerData, nil
}
