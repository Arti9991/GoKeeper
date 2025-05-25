package requseter

import (
	"context"
	"fmt"
	"log"
	"os"

	pb "github.com/Arti9991/GoKeeper/client/internal/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type ReqStruct struct {
	Ctx    context.Context
	Client pb.KeeperClient
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
	n, err := file.Write([]byte(ans.UserID))
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
	n, err := file.Write([]byte(ans.UserID))
	if err != nil || n == 0 {
		//logger.Log.Error("Error in saving to file", zap.Error(err))
		return err
	}
	return nil
}

// ctx := context.Background()

// 	respReg, err := c.Loginuser(ctx, &pb.LoginRequest{
// 		UserLogin:    "TestUserSECOND",
// 		UserPassword: "123456789",
// 	})
// 	if err != nil {
// 		log.Fatal(err)
// 	}

// 	// respReg, err := c.RegisterUser(ctx, &pb.RegisterRequest{
// 	// 	UserLogin:    "TestUserSECOND",
// 	// 	UserPassword: "123456789",
// 	// })
// 	// if err != nil {
// 	// 	log.Fatal(err)
// 	// }

// 	fmt.Println(respReg.UserID)
// 	var header metadata.MD
// 	md := metadata.New(map[string]string{"UserID": respReg.UserID})
// 	ctx2 := metadata.NewOutgoingContext(context.Background(), md)

// 	CurrTime := time.Now().Format(time.RFC850)
// 	savedData, err := c.SaveData(ctx2, &pb.SaveDataRequest{
// 		Data:     []byte("Hello there!"),
// 		DataType: "TEXT",
// 		Metainfo: "METAINFO",
// 		Time:     CurrTime,
// 	}, grpc.Header(&header))
// 	if err != nil {
// 		fmt.Println(err)
// 	}

// 	recievedData, err := c.GiveData(ctx2, &pb.GiveDataRequest{
// 		StorageID: savedData.StorageID,
// 	}, grpc.Header(&header))
// 	if err != nil {
// 		fmt.Println(err)
// 	}

// 	fmt.Printf("Recieved data: Type: %s, MetaInfo: %s, time: %s\n, Data: %s\n",
// 		recievedData.DataType,
// 		recievedData.Metainfo,
// 		recievedData.Time,
// 		string(recievedData.Data))

// 	_, err = c.UpdateData(ctx2, &pb.UpdateDataRequest{
// 		StorageID: savedData.StorageID,
// 		Data:      []byte("Hello there! UPDATE Not there... HERE!!!"),
// 		DataType:  "TEXT",
// 		Metainfo:  "second METAINFO updated",
// 		Time:      CurrTime,
// 	}, grpc.Header(&header))
// 	if err != nil {
// 		fmt.Println(err)
// 	}

// 	dataList, err := c.GiveDataList(ctx2, &pb.GiveDataListRequest{}, grpc.Header(&header))
// 	if err != nil {
// 		fmt.Println(err)
// 	}

// 	_, err = c.DeleteData(ctx2, &pb.DeleteDataRequest{
// 		StorageID: savedData.StorageID,
// 	}, grpc.Header(&header))
// 	if err != nil {
// 		fmt.Println(err)
// 	}

// 	fmt.Println(dataList)
