package requseter

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/Arti9991/GoKeeper/client/internal/binstor"
	"github.com/Arti9991/GoKeeper/client/internal/clientmodels"
	"github.com/Arti9991/GoKeeper/client/internal/dbstor"
	pb "github.com/Arti9991/GoKeeper/client/internal/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type ReqStruct struct {
	ServAddr string
	BinStor  *binstor.BinStor
	DBStor   *dbstor.DBStor
}

func NewRequester(addr string) *ReqStruct {
	var err error
	ReqStruct := new(ReqStruct)
	ReqStruct.ServAddr = addr

	ReqStruct.BinStor = binstor.NewBinStor(clientmodels.StorageDir)

	ReqStruct.DBStor, err = dbstor.DbInit("Journal.db")
	if err != nil {
		fmt.Println(err)
	}

	return ReqStruct
}

func TestLogin(req *ReqStruct) error {
	ctx := context.Background()

	dial, err := grpc.NewClient(req.ServAddr, grpc.WithTransportCredentials(insecure.NewCredentials())) //req.ServAddr
	if err != nil {
		log.Fatal(err)
	}

	r := pb.NewKeeperClient(dial)
	ans, err := r.Loginuser(ctx, &pb.LoginRequest{
		UserLogin:    "TestUserSECOND",
		UserPassword: "123456789",
	})
	if err != nil {
		return err
	}

	fmt.Println(ans.UserID)
	return nil
}

func LoginRequest(Login string, Password string, req *ReqStruct) error {

	ctx := context.Background()

	dial, err := grpc.NewClient(req.ServAddr, grpc.WithTransportCredentials(insecure.NewCredentials())) //req.ServAddr
	if err != nil {
		log.Fatal(err)
	}
	r := pb.NewKeeperClient(dial)
	ans, err := r.Loginuser(ctx, &pb.LoginRequest{
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

func RegisterRequest(Login string, Password string, req *ReqStruct) error {
	ctx := context.Background()

	dial, err := grpc.NewClient(req.ServAddr, grpc.WithTransportCredentials(insecure.NewCredentials())) //req.ServAddr
	if err != nil {
		log.Fatal(err)
	}
	r := pb.NewKeeperClient(dial)
	ans, err := r.RegisterUser(ctx, &pb.RegisterRequest{
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

func LogoutRequest(req *ReqStruct) error {
	var err error

	err = req.DBStor.ReinitTable()
	if err != nil {
		return err
	}

	err = os.Remove(clientmodels.StorageDir)
	if err != nil {
		return err
	}
	err = os.Mkdir(clientmodels.StorageDir, 0644)

	err = os.Remove("./Token.txt")
	if err != nil {
		return err
	}

	return err
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
