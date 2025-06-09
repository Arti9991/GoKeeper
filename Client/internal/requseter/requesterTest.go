package requseter

import (
	"bufio"
	"bytes"
	"context"
	"crypto/sha256"
	"encoding/gob"
	"encoding/hex"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/Arti9991/GoKeeper/client/internal/clientmodels"
	pb "github.com/Arti9991/GoKeeper/client/internal/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

func SaveDataTest(DtInf clientmodels.NewerData, req *ReqStruct, offlineMode bool) (string, error) {

	hashData := sha256.Sum256(DtInf.Data)
	DtInf.StorageID = hex.EncodeToString(hashData[:])

	err := req.BinStor.SaveBinData(DtInf.StorageID, DtInf.Data)
	if err != nil {
		return " ", err
	}

	DtInf.SaveTime = time.Now().UTC().Format(time.RFC850)

	err = req.DBStor.SaveNew(DtInf.StorageID, DtInf)
	if err != nil {
		return " ", err
	}

	if !offlineMode {
		err = SendWithUpdate(DtInf.StorageID, DtInf, req, DtInf.Data)
		if err != nil {
			return " ", err
		}
	}

	fmt.Printf("\nNew data saved! StorageID for new data: %s\n", DtInf.StorageID)
	return DtInf.StorageID, nil
}

func TestGetDataRequest(StorageID string, req *ReqStruct) error {
	var UserID string

	fmt.Println("Open token")
	file, err := os.Open("./Token.txt")
	if err != nil {
		return err
	}
	reader := bufio.NewReader(file)
	// Считываем строку текста
	UserID, err = reader.ReadString('\n')
	if err != nil {
		return err
	}
	//Выводим строку
	UserID = strings.TrimSuffix(UserID, "\n")
	fmt.Printf("%#v", UserID)

	var header metadata.MD
	md := metadata.New(map[string]string{"UserID": UserID})
	ctx := metadata.NewOutgoingContext(context.Background(), md)

	dial, err := grpc.NewClient(":8082", grpc.WithTransportCredentials(req.Creds)) //":8082"
	if err != nil {
		return err
	}

	r := pb.NewKeeperClient(dial)
	ans, err := r.GiveData(ctx, &pb.GiveDataRequest{
		StorageID: StorageID,
	}, grpc.Header(&header))
	if err != nil {
		return err
	}

	switch ans.DataType {
	case "CARD":

		fmt.Println(ans.Data)
		//fmt.Println(buff2)
		fmt.Println(ans.Metainfo)

		out := pb.CardInfo{}
		dec := gob.NewDecoder(bytes.NewBuffer(ans.Data))
		dec.Decode(&out)
		fmt.Println(out.Number)
		fmt.Println(out.CVVcode)
		fmt.Println(out.ExpDate)
		fmt.Println(out.Holder)
	}
	return nil
}

func UpdateDataTest(DtInf clientmodels.NewerData, req *ReqStruct, offlineMode bool) error {

	DtInf.SaveTime = time.Now().UTC().Format(time.RFC850)

	err := UpdateDataOfline(DtInf, req)
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
