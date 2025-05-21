package server

import (
	"context"
	"errors"
	"fmt"
	"testing"

	"github.com/Arti9991/GoKeeper/server/internal/config"
	"github.com/Arti9991/GoKeeper/server/internal/logger"
	"github.com/Arti9991/GoKeeper/server/internal/server/mocks"
	pb "github.com/Arti9991/GoKeeper/server/internal/server/proto"
	"github.com/Arti9991/GoKeeper/server/internal/storage/pgstor"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// для базовых тестов производится генерация моков командой ниже
// mockgen --source=./internal/storage/pgstor/pgStor.go --destination=./internal/server/mocks/mocks_datainfo.go --package=mocks InfoStorage
// mockgen --source=./internal/storage/pgstor/pgStorUsers.go --destination=./internal/server/mocks/mocks_users.go --package=mocks UserStor

func InitServerTest(UserStor pgstor.UserStor, InfoStor pgstor.InfoStorage) *Server {
	Serv := new(Server)

	Serv.Config = config.InitConfTest()

	logger.Initialize(Serv.Config.InFileLog)
	logger.Log.Info("Logger initialyzed!",
		zap.Bool("In file mode:", Serv.Config.InFileLog),
	)

	Serv.UserStor = UserStor
	Serv.InfoStor = InfoStor

	return Serv
}

// func RunServerTest() error {
// 	server := InitServerTest()

// 	fmt.Println("Host addr", server.Config.HostAddr)
// 	// определяем адрес для сервера
// 	listen, err := net.Listen("tcp", server.Config.HostAddr)
// 	if err != nil {
// 		return err
// 	}
// 	// создаём gRPC-сервер без зарегистрированной службы
// 	interceptors := grpc.ChainUnaryInterceptor(interceptors.AtuhInterceptor, interceptors.LoggingInterceptor)
// 	s := grpc.NewServer(interceptors)

// 	proto.RegisterKeeperServer(s, server)

// 	// получаем запрос gRPC
// 	if err := s.Serve(listen); err != nil {
// 		log.Fatal(err)
// 	}
// 	return nil
// }

func TestRegisterUser(t *testing.T) {
	// создаём контроллер
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// создаём объект-заглушку
	mockUsers := mocks.NewMockUserStor(ctrl)
	mockInfo := mocks.NewMockInfoStorage(ctrl)

	serv := InitServerTest(mockUsers, mockInfo)
	// go func() {
	// 	err := RunServerTest()
	// 	require.NoError(t, err)
	// }()

	// message RegisterRequest {
	// 	string UserLogin = 1;
	// 	string UserPassword = 2;

	//   }

	//   message RegisterResponce {
	// 	string UserID = 1;
	//   }

	type want struct {
		UserID   string
		serv_err error
	}
	tests := []struct {
		Name         string
		UserLogin    string
		UserPassword string
		err          error
		want         want
	}{
		{
			Name:         "Simple registration",
			UserLogin:    "Test Login",
			UserPassword: "1234567890",
			err:          nil,
			want: want{
				UserID:   "XDOJ6FD32JUYVJJ4",
				serv_err: nil,
			},
		},
		{
			Name:         "Repeated registration",
			UserLogin:    "Test Login",
			UserPassword: "1234567890",
			err:          errors.New("Ошибка в сохранении пользователя"),
			want: want{
				UserID:   "XDOJ6FD32JUYVJJ4",
				serv_err: status.Error(codes.Unavailable, `Ошибка в сохранении пользователя`),
			},
		},
	}

	// устанавливаем соединение с сервером
	// conn, err := grpc.NewClient(":8081", grpc.WithTransportCredentials(insecure.NewCredentials()))
	// require.NoError(t, err)
	// defer conn.Close()
	ctx := context.Background()

	for _, test := range tests {
		fmt.Println("В тесте 1")
		// задаем режим работы моков (для POST главное отсутствие ошибки)
		mockUsers.EXPECT().
			SaveNewUser(gomock.Any(), gomock.Any(), gomock.Any()).
			Return(test.err).
			MaxTimes(1)

		t.Run(test.Name, func(t *testing.T) {
			fmt.Println("Запускаем")
			ans, err := serv.RegisterUser(ctx, &pb.RegisterRequest{
				UserLogin:    test.UserLogin,
				UserPassword: test.UserPassword,
			})
			fmt.Println(err)
			require.Equal(t, err, test.want.serv_err)
			if err == nil {
				require.NotEmpty(t, ans)
			}

			fmt.Println("Id user:", ans.UserID)
			fmt.Println("Закончили")
		})
	}
}
