package requseter

import (
	"crypto/rand"
	"crypto/x509"
	"fmt"
	"io/ioutil"
	"log"
	"testing"

	"github.com/Arti9991/GoKeeper/client/internal/binstor"
	"github.com/Arti9991/GoKeeper/client/internal/clientmodels"
	"github.com/Arti9991/GoKeeper/client/internal/dbstor"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/status"
)

func RequesterTest() *ReqStruct {
	var err error
	ReqStruct := new(ReqStruct)
	ReqStruct.ServAddr = ":8082"

	ReqStruct.BinStor = binstor.NewBinStor(clientmodels.StorageDir)

	ReqStruct.DBStor, err = dbstor.DbInit("Journal.db")
	if err != nil {
		fmt.Println(err)
	}

	// Загружаем сертификат, которому доверяем (тот, что сгенерирован на сервере)
	caCert, err := ioutil.ReadFile("./../../cmd/client/server.crt")
	if err != nil {

	}

	// Создаём пул корневых сертификатов и добавляем туда server.crt
	certPool := x509.NewCertPool()
	if !certPool.AppendCertsFromPEM(caCert) {
		log.Fatalf("Failed to add server.crt to cert pool")
	}

	// Настраиваем TLS
	ReqStruct.Creds = credentials.NewClientTLSFromCert(certPool, "localhost") // CN должен совпадать с /CN= в server.crt

	return ReqStruct
}

func TestRegisterLoginUser(t *testing.T) {

	req := RequesterTest()

	AllPassw := rand.Text()[:8]
	OneLogin := rand.Text()[:6]

	type want struct {
		serv_err1 error
		serv_err2 error
	}
	tests := []struct {
		Name          string
		UserLogin     string
		UserPassword  string
		UserLogin2    string
		UserPassword2 string
		anotherAuth   bool
		want          want
	}{
		{
			Name:         "Simple registration and login",
			UserLogin:    rand.Text()[:6],
			UserPassword: AllPassw,
			anotherAuth:  false,
			want: want{
				serv_err1: nil,
				serv_err2: nil,
			},
		},
		{
			Name:         "Simple registration with no password and login",
			UserLogin:    rand.Text()[:6],
			UserPassword: "",
			anotherAuth:  false,
			want: want{
				serv_err1: clientmodels.ErrBadPassowrd,
				serv_err2: status.Error(codes.PermissionDenied, `Неверное имя пользователя или пароль`),
			},
		},
		{
			Name:         "Simple registration short username and login",
			UserLogin:    rand.Text()[:3],
			UserPassword: AllPassw,
			anotherAuth:  false,
			want: want{
				serv_err1: clientmodels.ErrBadLogin,
				serv_err2: status.Error(codes.PermissionDenied, `Неверное имя пользователя или пароль`),
			},
		},
		{
			Name:          "Simple registration and login with another Login",
			UserLogin:     rand.Text()[:6],
			UserPassword:  AllPassw,
			UserLogin2:    "Some Another Login",
			UserPassword2: AllPassw,
			anotherAuth:   true,
			want: want{
				serv_err1: nil,
				serv_err2: status.Error(codes.PermissionDenied, `Неверное имя пользователя или пароль`),
			},
		},
		{
			Name:          "Simple registration and login with another password",
			UserLogin:     OneLogin,
			UserPassword:  AllPassw,
			UserLogin2:    OneLogin,
			UserPassword2: "Some Another Passw",
			anotherAuth:   true,
			want: want{
				serv_err1: nil,
				serv_err2: status.Error(codes.PermissionDenied, `Неверное имя пользователя или пароль`),
			},
		},
	}
	for _, test := range tests {
		err := RegisterRequest(test.UserLogin, test.UserPassword, req)
		require.Equal(t, test.want.serv_err1, err)
		if !test.anotherAuth {
			err = LoginRequest(test.UserLogin, test.UserPassword, req)
			require.Equal(t, test.want.serv_err2, err)
		} else {
			err = LoginRequest(test.UserLogin2, test.UserPassword2, req)
			require.Equal(t, test.want.serv_err2, err)
		}
	}
	err := LogoutRequest(req)
	require.NoError(t, err)
}

func TestBasicReq(t *testing.T) {
	req := RequesterTest()

	AllLogin := rand.Text()[:6]
	AllPassw := rand.Text()[:8]

	type want struct {
		serv_err1 error
	}
	tests := []struct {
		Name         string
		TestType     string
		OpType       string
		UserLogin    string
		UserPassword string
		offlineMode  bool
		Data         clientmodels.NewerData
		Data2        clientmodels.NewerData
		want         want
	}{
		{
			Name:         "--------Simple registration------------------",
			TestType:     "REG",
			UserLogin:    AllLogin,
			UserPassword: AllPassw,
			offlineMode:  false,
			want: want{
				serv_err1: nil,
			},
		},
		{
			Name:         "---------Simple save test data---------------",
			TestType:     "SAVE",
			OpType:       "TEXT",
			UserLogin:    AllLogin,
			UserPassword: AllPassw,
			offlineMode:  false,
			Data: clientmodels.NewerData{
				DataType: "TEXT",
				MetaInfo: "Metainfo for test",
				Data:     []byte("TEXT FOR SAVE"),
			},
			want: want{
				serv_err1: nil,
			},
		},
		{
			Name:         "-------------Update data test------------------",
			TestType:     "UPDATE",
			OpType:       "TEXT",
			UserLogin:    AllLogin,
			UserPassword: AllPassw,
			offlineMode:  false,
			Data: clientmodels.NewerData{
				DataType: "TEXT",
				MetaInfo: "Metainfo for test UPDATED",
				Data:     []byte("TEXT FOR SAVE AND UPDATE"),
			},
			Data2: clientmodels.NewerData{
				DataType: "TEXT",
				MetaInfo: "Metainfo for test UPDATED. META UPDATED!!!",
				Data:     []byte("TEXT FOR SAVE AND UPDATE. THIS TEXT IS UPDATED!!!"),
			},
			want: want{
				serv_err1: nil,
			},
		},
		{
			Name:         "------------Delete data test-----------------",
			TestType:     "DELETE",
			OpType:       "TEXT",
			UserLogin:    AllLogin,
			UserPassword: AllPassw,
			offlineMode:  false,
			Data: clientmodels.NewerData{
				DataType: "TEXT",
				MetaInfo: "Metainfo for test delete",
				Data:     []byte("TEXT FOR DELETE"),
			},
			want: want{
				serv_err1: nil,
			},
		},
		{
			Name:         "-----------Offline save and update data test-------------",
			TestType:     "OFFLINE",
			OpType:       "TEXT",
			UserLogin:    AllLogin,
			UserPassword: AllPassw,
			offlineMode:  true,
			Data: clientmodels.NewerData{
				DataType: "TEXT",
				MetaInfo: "Metainfo for test UPDATED",
				Data:     []byte("TEXT FOR SAVE AND UPDATE"),
			},
			Data2: clientmodels.NewerData{
				DataType: "TEXT",
				MetaInfo: "Metainfo for test UPDATED. META UPDATED!!!",
				Data:     []byte("TEXT FOR SAVE AND UPDATE. THIS TEXT IS UPDATED!!!"),
			},
			want: want{
				serv_err1: nil,
			},
		},
		{
			Name:         "-----------Show data test-------------",
			TestType:     "SHOW",
			OpType:       "TEXT",
			UserLogin:    AllLogin,
			UserPassword: AllPassw,
			Data: clientmodels.NewerData{
				DataType: "TEXT",
				MetaInfo: "Metainfo for test SHOW ONLINE",
				Data:     []byte("TEXT FOR SAVE AND SHOW ONLINE"),
			},
			Data2: clientmodels.NewerData{
				DataType: "TEXT",
				MetaInfo: "Metainfo for test SHOW OFFLINE",
				Data:     []byte("TEXT FOR SAVE AND SHOW OFFLINE"),
			},
			want: want{
				serv_err1: nil,
			},
		},
		// {
		// 	Name:         "Simple registration with no password and login",
		// 	UserLogin:    rand.Text()[:6],
		// 	UserPassword: "",
		// 	anotherAuth:  false,
		// 	want: want{
		// 		serv_err1: clientmodels.ErrBadPassowrd,
		// 		serv_err2: status.Error(codes.PermissionDenied, `Неверное имя пользователя или пароль`),
		// 	},
		// },
		// {
		// 	Name:         "Simple registration short username and login",
		// 	UserLogin:    rand.Text()[:3],
		// 	UserPassword: AllPassw,
		// 	anotherAuth:  false,
		// 	want: want{
		// 		serv_err1: clientmodels.ErrBadLogin,
		// 		serv_err2: status.Error(codes.PermissionDenied, `Неверное имя пользователя или пароль`),
		// 	},
		// },
		// {
		// 	Name:          "Simple registration and login with another Login",
		// 	UserLogin:     rand.Text()[:6],
		// 	UserPassword:  AllPassw,
		// 	UserLogin2:    "Some Another Login",
		// 	UserPassword2: AllPassw,
		// 	anotherAuth:   true,
		// 	want: want{
		// 		serv_err1: nil,
		// 		serv_err2: status.Error(codes.PermissionDenied, `Неверное имя пользователя или пароль`),
		// 	},
		// },
		// {
		// 	Name:          "Simple registration and login with another password",
		// 	UserLogin:     OneLogin,
		// 	UserPassword:  AllPassw,
		// 	UserLogin2:    OneLogin,
		// 	UserPassword2: "Some Another Passw",
		// 	anotherAuth:   true,
		// 	want: want{
		// 		serv_err1: nil,
		// 		serv_err2: status.Error(codes.PermissionDenied, `Неверное имя пользователя или пароль`),
		// 	},
		// },
	}
	for _, test := range tests {
		fmt.Printf("\n\n%s\n\n", test.Name)
		switch test.TestType {
		case "REG":
			err := RegisterRequest(test.UserLogin, test.UserPassword, req)
			require.Equal(t, test.want.serv_err1, err)
		case "SAVE":
			StorID, err := SaveDataTest(test.Data, req, test.offlineMode)
			require.Equal(t, test.want.serv_err1, err)

			err = GetDataRequest(StorID, req, test.offlineMode)
			require.Equal(t, test.want.serv_err1, err)
		case "UPDATE":
			StorID, err := SaveDataTest(test.Data, req, test.offlineMode)
			require.Equal(t, test.want.serv_err1, err)

			test.Data2.StorageID = StorID
			err = UpdateDataTest(test.Data2, req, test.offlineMode)
			require.Equal(t, test.want.serv_err1, err)

			err = GetDataRequest(StorID, req, test.offlineMode)
			require.Equal(t, test.want.serv_err1, err)

		case "DELETE":
			StorID, err := SaveDataTest(test.Data, req, test.offlineMode)
			require.Equal(t, test.want.serv_err1, err)

			err = DeleteDataRequest(StorID, req, test.offlineMode)
			require.Equal(t, test.want.serv_err1, err)
		case "OFFLINE":
			StorID, err := SaveDataTest(test.Data, req, test.offlineMode)
			require.Equal(t, test.want.serv_err1, err)

			test.Data2.StorageID = StorID
			err = UpdateDataTest(test.Data2, req, test.offlineMode)
			require.Equal(t, test.want.serv_err1, err)

			err = GetDataRequest(StorID, req, test.offlineMode)
			require.Equal(t, test.want.serv_err1, err)
		case "SHOW":
			_, err := SaveDataTest(test.Data, req, false)
			require.Equal(t, test.want.serv_err1, err)

			_, err = SaveDataTest(test.Data2, req, true)
			require.Equal(t, test.want.serv_err1, err)

			err = ShowDataLoc(req)
			require.Equal(t, test.want.serv_err1, err)
			err = ShowDataOn(req, false)
			require.Equal(t, test.want.serv_err1, err)
		}

	}

	err := LogoutRequest(req)
	require.NoError(t, err)
}
