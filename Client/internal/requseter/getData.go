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

func GetDataRequest(StorageID string, req *ReqStruct, offlineMode bool) error {
	var ans clientmodels.NewerData
	var ansOn clientmodels.NewerData
	var ansOf clientmodels.NewerData
	var err error
	var err2 error

	fmt.Println(offlineMode)

	ansOf, err = GetDataOfline(StorageID, req)
	if err != nil {
		fmt.Println(err)
	}

	TimeLoc, err := time.Parse(time.RFC850, ansOf.SaveTime)
	if err != nil {
		fmt.Println(err)
	}
	if !offlineMode {
		ansOn, err2 = CompareGetData(StorageID, TimeLoc, req)
	}
	if offlineMode || err2 != nil {
		ans = ansOf
	} else {
		ans = ansOn
	}
	err = inputfunc.ParceAnswer(ans.Data, StorageID, ans.DataType, ans.MetaInfo)
	if err != nil {
		fmt.Println(err)
		return err
	}
	return err
}

func GetDataOnline(StorageID string, addr string) (clientmodels.NewerData, error) {
	var DataGet clientmodels.NewerData
	var UserID string

	fmt.Println("Open token")
	file, err := os.Open("./Token.txt")
	if err != nil {
		fmt.Println(err)
		//logger.Log.Error("SAVE Error in opening file", zap.Error(err))
		return DataGet, err
	}
	defer file.Close()

	reader := bufio.NewReader(file)
	// Считываем строку текста
	UserID, err = reader.ReadString('\n')
	if err != nil {
		fmt.Println(err)
		return DataGet, err
	}
	//Выводим строку
	UserID = strings.TrimSuffix(UserID, "\n")
	fmt.Printf("%#v", UserID)

	var header metadata.MD
	md := metadata.New(map[string]string{"UserID": UserID})
	ctx := metadata.NewOutgoingContext(context.Background(), md)

	dial, err := grpc.NewClient(addr, grpc.WithTransportCredentials(insecure.NewCredentials())) //":8082"
	if err != nil {
		log.Fatal(err)
	}

	req := pb.NewKeeperClient(dial)
	ans, err := req.GiveData(ctx, &pb.GiveDataRequest{
		StorageID: StorageID,
	}, grpc.Header(&header))
	if err != nil {
		fmt.Println(err)
		return DataGet, err
	}

	DataGet.Data = ans.Data
	DataGet.DataType = ans.DataType
	DataGet.MetaInfo = ans.Metainfo
	DataGet.StorageID = StorageID
	DataGet.SaveTime = ans.Time
	return DataGet, nil
}

func GetDataOfline(StorageID string, req *ReqStruct) (clientmodels.NewerData, error) {
	var DataGet clientmodels.NewerData
	var err error
	DataGet.StorageID = StorageID

	DataGet, err = req.DBStor.Get(StorageID)
	if err != nil {
		return DataGet, err
	}

	DataGet.Data, err = req.BinStor.GetBinData(StorageID)
	if err != nil {
		return DataGet, err
	}
	return DataGet, nil
}

func CompareGetData(StorageID string, TimeLoc time.Time, req *ReqStruct) (clientmodels.NewerData, error) {
	var ansOn clientmodels.NewerData
	var err error

	ansOn, err = GetDataOnline(StorageID, req.ServAddr)
	if err != nil {
		fmt.Println(err)
		return ansOn, err
	}
	timeServ, err := time.Parse(time.RFC850, ansOn.SaveTime)
	if err != nil {
		fmt.Println(err)
		return ansOn, err
	}
	if timeServ.After(TimeLoc) {
		err = req.DBStor.MarkUnDone(StorageID)
		if err != nil {
			fmt.Println(err)
			return ansOn, err
		}
	}
	return ansOn, nil
}
