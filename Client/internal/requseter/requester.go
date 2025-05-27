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

	"github.com/Arti9991/GoKeeper/client/internal/binstor"
	"github.com/Arti9991/GoKeeper/client/internal/clientmodels"
	"github.com/Arti9991/GoKeeper/client/internal/dbstor"
	"github.com/Arti9991/GoKeeper/client/internal/inputfunc"
	"github.com/Arti9991/GoKeeper/client/internal/journal"
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
	n, err := file.Write([]byte(ans.UserID + "\n"))
	if err != nil || n == 0 {
		//logger.Log.Error("Error in saving to file", zap.Error(err))
		return err
	}
	return nil
}

func SyncRequest() error {
	//offlineMode := false

	var UserID string

	fmt.Println("Open token")
	file, err := os.Open(clientmodels.TokenFile)
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

	Jrn, err := journal.JournalGet()
	if err != nil {
		fmt.Println(err)
		return err
	}
	fmt.Println(Jrn)
	// if !offlineMode {

	// }

	// for _, JrStr := range Jrn {
	// 	var header metadata.MD
	// 	md := metadata.New(map[string]string{"UserID": UserID})
	// 	ctx := metadata.NewOutgoingContext(context.Background(), md)

	// 	dial, err := grpc.NewClient(":8082", grpc.WithTransportCredentials(insecure.NewCredentials())) //":8082"
	// 	if err != nil {
	// 		log.Fatal(err)
	// 	}
	// 	CurrTime := time.Now().Format(time.RFC850)
	// 	req := pb.NewKeeperClient(dial)
	// 	ans, err := req.SaveData(ctx, &pb.SaveDataRequest{
	// 		Metainfo: JrStr.MetaInfo,
	// 		DataType: JrStr.MetaInfo,
	// 		Time:     JrStr.SaveTime,
	// 		Data:
	// 	}, grpc.Header(&header))
	// 	if err != nil {
	// 		fmt.Println(err)
	// 		return err
	// 	}
	// }
	// JrInf := clientmodels.JournalInfo{
	// 	Opperation: "SAVE",
	// 	StorageID:  ans.StorageID,
	// 	DataType:   Type,
	// 	MetaInfo:   Metainfo,
	// 	SaveTime:   CurrTime,
	// }

	// err = journal.JournalSave(JrInf)

	// fmt.Println(ans.StorageID)

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

	JrInf := clientmodels.JournalInfo{
		Opperation: "SAVE",
		StorageID:  StorageID,
		DataType:   Type,
		MetaInfo:   Metainfo,
		SaveTime:   CurrTime,
		Sync:       false,
	}

	//err = journal.JournalSave(JrInf)

	//err = db.ReinitTable()
	// fmt.Println(err)

	err = req.DBStor.SaveNew(StorageID, JrInf)
	fmt.Println(err)

	if !offlineMode {
		err = SendWithUpdate(StorageID, JrInf, req, data)
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
	var err error

	fmt.Println(offlineMode)

	if !offlineMode {
		ans, err = GetDataOnline(StorageID)
	}
	if offlineMode || err != nil {
		ans, err = GetDataOfline(StorageID, req)
		if err != nil {
			fmt.Println(err)
			return err
		}
	}
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
