package requseter

import (
	"bufio"
	"context"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/Arti9991/GoKeeper/client/internal/binstor"
	"github.com/Arti9991/GoKeeper/client/internal/clientmodels"
	"github.com/Arti9991/GoKeeper/client/internal/dbstor"
	"github.com/Arti9991/GoKeeper/client/internal/inputfunc"
	pb "github.com/Arti9991/GoKeeper/client/internal/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"
)

type ReqStruct struct {
	Ctx     context.Context
	Client  pb.KeeperClient
	BinStor *binstor.BinStor
	DBStor  *dbstor.DBStor
}

func NewRequester(addr string) *ReqStruct {
	var err error
	ReqStruct := new(ReqStruct)
	dial, err := grpc.NewClient(addr, grpc.WithTransportCredentials(insecure.NewCredentials())) //":8082"
	if err != nil {
		log.Fatal(err)
	}
	ReqStruct.Client = pb.NewKeeperClient(dial)

	ReqStruct.Ctx = context.Background()

	ReqStruct.BinStor = binstor.NewBinStor(clientmodels.StorageDir)

	ReqStruct.DBStor, err = dbstor.DbInit("Journal.db")
	if err != nil {
		fmt.Println(err)
	}

	return ReqStruct
}

func (r *ReqStruct) TestLogin() error {

	ans, err := r.Client.Loginuser(r.Ctx, &pb.LoginRequest{
		UserLogin:    "TestUserSECOND",
		UserPassword: "123456789",
	})
	if err != nil {
		return err
	}

	fmt.Println(ans.UserID)
	return nil
}

func TestLoginRequest(Login string, Password string) error {

	ctx := context.Background()

	dial, err := grpc.NewClient(":8082", grpc.WithTransportCredentials(insecure.NewCredentials())) //":8082"
	if err != nil {
		log.Fatal(err)
	}
	req := pb.NewKeeperClient(dial)
	ans, err := req.Loginuser(ctx, &pb.LoginRequest{
		UserLogin:    Login,
		UserPassword: Password,
	})
	if err != nil {
		return err
	}

	fmt.Println(ans.UserID)
	file, err := os.OpenFile("./Token.txt", os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		//logger.Log.Error("SAVE Error in opening file", zap.Error(err))
		return err
	}
	defer file.Close()

	n, err := file.Write([]byte(ans.UserID + "\n"))
	if err != nil || n == 0 {
		//logger.Log.Error("Error in saving to file", zap.Error(err))
		return err
	}
	return nil
}

func TestRegisterRequest(Login string, Password string) error {
	ctx := context.Background()

	dial, err := grpc.NewClient(":8082", grpc.WithTransportCredentials(insecure.NewCredentials())) //":8082"
	if err != nil {
		log.Fatal(err)
	}
	req := pb.NewKeeperClient(dial)
	ans, err := req.RegisterUser(ctx, &pb.RegisterRequest{
		UserLogin:    Login,
		UserPassword: Password,
	})
	if err != nil {
		return err
	}

	fmt.Println(ans.UserID)
	file, err := os.OpenFile("./Token.txt", os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		//logger.Log.Error("SAVE Error in opening file", zap.Error(err))
		return err
	}
	defer file.Close()

	n, err := file.Write([]byte(ans.UserID + "\n"))
	if err != nil || n == 0 {
		//logger.Log.Error("Error in saving to file", zap.Error(err))
		return err
	}
	return nil
}

func SyncRequest(req *ReqStruct, offlineMode bool) error {

	if offlineMode {
		return errors.New("cannot sync in offline mode")
	}

	var UserID string

	fmt.Println("Open token")
	file, err := os.Open(clientmodels.TokenFile)
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
	fmt.Printf("%#v\n", UserID)

	SnData, err := req.DBStor.GetForSync()
	if err != nil {
		fmt.Println(err)
		return err
	}

	for _, syncPart := range SnData {
		//fmt.Println(syncPart)

		syncPart.Data, err = req.BinStor.GetBinData(syncPart.StorageID)
		if err != nil {
			fmt.Println(err)
			return err
		}
		err = SendWithUpdate(syncPart.StorageID, syncPart, req, syncPart.Data)
		if err != nil {
			fmt.Println(err)
			return err
		}
	}
	return nil
}

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

	CurrTime := time.Now().Format(time.RFC850)

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
		return err
	}

	TimeLoc, err := time.Parse(time.RFC850, ansOf.SaveTime)
	if err != nil {
		fmt.Println(err)
	}
	if !offlineMode {
		ansOn, err2 = GetDataOnline(StorageID)
		fmt.Println(err2)
		timeServ, err := time.Parse(time.RFC850, ansOn.SaveTime)
		if err != nil {
			fmt.Println(err)
		}
		if timeServ.After(TimeLoc) {
			err = req.DBStor.MarkUnDone(StorageID)
			if err != nil {
				fmt.Println(err)
			}
		}
	}
	if offlineMode || err2 != nil {
		ans = ansOf
	} else {
		ans = ansOn
	}
	// Не обноваляется дата в локальном хранилище при получении более новой через sync
	///////////////////////////
	/////////////////////////
	///////////////////////
	///////////////////
	/////////////
	err = inputfunc.ParceAnswer(ans.Data, StorageID, ans.DataType, ans.MetaInfo)
	if err != nil {
		fmt.Println(err)
		return err
	}
	return nil
}

func GetDataOnline(StorageID string) (clientmodels.NewerData, error) {
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

func DeleteDataRequest(StorageID string, req *ReqStruct, offlineMode bool) error {

	var err error

	fmt.Println(offlineMode)

	if !offlineMode {
		err = DeleteDataOnline(StorageID)
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

func DeleteDataOnline(StorageID string) error {

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

	dial, err := grpc.NewClient(":8082", grpc.WithTransportCredentials(insecure.NewCredentials())) //":8082"
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
