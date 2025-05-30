package server

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/Arti9991/GoKeeper/server/internal/coder"
	"github.com/Arti9991/GoKeeper/server/internal/config"
	"github.com/Arti9991/GoKeeper/server/internal/logger"
	"github.com/Arti9991/GoKeeper/server/internal/server/interceptors"
	"github.com/Arti9991/GoKeeper/server/internal/server/mocks"
	pb "github.com/Arti9991/GoKeeper/server/internal/server/proto"
	"github.com/Arti9991/GoKeeper/server/internal/server/servermodels"
	"github.com/Arti9991/GoKeeper/server/internal/storage/binstortest"
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
// (устарело!)mockgen --source=./internal/storage/binstor/binstor.go --destination=./internal/server/mocks/mocks_bin.go --package=mocks BinStrorFunc

// функция для инциализации тестовой структуры Server
func InitServerTest(UserStor pgstor.UserStor, InfoStor pgstor.InfoStorage) *Server {
	Serv := new(Server)

	Serv.Config = config.InitConfTest()

	logger.InitializeTest(Serv.Config.InFileLog)
	logger.Log.Info("Logger initialyzed!",
		zap.Bool("In file mode:", Serv.Config.InFileLog),
	)

	Serv.UserStor = UserStor
	Serv.InfoStor = InfoStor

	Serv.BinStorFunc = binstortest.NewBinStorTest()

	return Serv
}

// TestRegisterUser тестирование функции регистрации
func TestRegisterUser(t *testing.T) {
	// создаём контроллер
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// создаём объект-заглушку
	mockUsers := mocks.NewMockUserStor(ctrl)
	mockInfo := mocks.NewMockInfoStorage(ctrl)
	// создаем тестовую структуру
	serv := InitServerTest(mockUsers, mockInfo)

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
			Name:         "DB error registration",
			UserLogin:    "Test Login",
			UserPassword: "1234567890",
			err:          errors.New("какая-то шибка в сохранении пользователя"),
			want: want{
				UserID:   "XDOJ6FD32JUYVJJ4",
				serv_err: status.Error(codes.Unavailable, `Ошибка в сохранении пользователя`),
			},
		},
		{
			Name:         "User already exist registration",
			UserLogin:    "Test Login",
			UserPassword: "1234567890",
			err:          servermodels.ErrorUserAlready,
			want: want{
				UserID:   "XDOJ6FD32JUYVJJ4",
				serv_err: status.Error(codes.Unavailable, `Пользователь уже зарегистрирован`),
			},
		},
	}
	// создание заглушки контекста для запросов
	ctx := context.Background()

	for _, test := range tests {
		// задаем режим работы моков для сохранения пользователя
		mockUsers.EXPECT().
			SaveNewUser(gomock.Any(), gomock.Any(), gomock.Any()).
			Return(test.err).
			MaxTimes(1)
		// выполняем запрос
		t.Run(test.Name, func(t *testing.T) {
			ans, err := serv.RegisterUser(ctx, &pb.RegisterRequest{
				UserLogin:    test.UserLogin,
				UserPassword: test.UserPassword,
			})
			require.Equal(t, err, test.want.serv_err)
			if err == nil {
				// если ошибка нулевая, то должен быть ответ
				require.NotEmpty(t, ans)
			}
		})
	}
}

// TestLoginUser тестирование функции логина
func TestLoginUser(t *testing.T) {
	// создаём контроллер
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// создаём объект-заглушку
	mockUsers := mocks.NewMockUserStor(ctrl)
	mockInfo := mocks.NewMockInfoStorage(ctrl)

	// создаем тестовую структуру
	serv := InitServerTest(mockUsers, mockInfo)

	type want struct {
		serv_err error
	}
	tests := []struct {
		Name         string
		UserLogin    string
		UserPassword string
		GetPassword  string
		UserID       string
		err          error
		want         want
	}{
		{
			Name:         "Simple Login",
			UserLogin:    "Test Login",
			UserPassword: "1234567890",
			GetPassword:  "1234567890",
			UserID:       "XDOJ6FD32JUYVJJ4",
			err:          nil,
			want: want{
				serv_err: nil,
			},
		},
		{
			Name:         "Login with bad password",
			UserLogin:    "Test Login",
			UserPassword: "1234567890",
			GetPassword:  "0192837465",
			UserID:       "XDOJ6FD32JUYVJJ4",
			err:          nil,
			want: want{
				serv_err: status.Error(codes.PermissionDenied, `Неверное имя пользователя или пароль`),
			},
		},
		{
			Name:         "Login with bad login",
			UserLogin:    "Not Test Login",
			UserPassword: "1234567890",
			GetPassword:  "1234567890",
			UserID:       "XDOJ6FD32JUYVJJ4",
			err:          servermodels.ErrorNoSuchUser,
			want: want{
				serv_err: status.Error(codes.PermissionDenied, `Неверное имя пользователя или пароль`),
			},
		},
		{
			Name:         "Login with DB error",
			UserLogin:    "Test Login",
			UserPassword: "1234567890",
			GetPassword:  "1234567890",
			UserID:       "XDOJ6FD32JUYVJJ4",
			err:          errors.New("some db error"),
			want: want{
				serv_err: status.Error(codes.Unavailable, `Ошибка в получении пользователя`),
			},
		},
	}
	// создание заглушки контекста для запросов
	ctx := context.Background()

	for _, test := range tests {
		// кодированный пароль, который мы рассчитываем получить
		codedPass := servermodels.CodePassword(test.GetPassword)
		// задаем режим работы моков для получения пользователя
		mockUsers.EXPECT().
			GetUser(gomock.Any()).
			Return(test.UserID, codedPass, test.err).
			MaxTimes(1)
			// выполняем запрос
		t.Run(test.Name, func(t *testing.T) {
			ans, err := serv.Loginuser(ctx, &pb.LoginRequest{
				UserLogin:    test.UserLogin,
				UserPassword: test.UserPassword,
			})
			require.Equal(t, err, test.want.serv_err)
			if err == nil {
				// если ошибка нулевая, то должен быть ответ
				require.NotEmpty(t, ans)
			}

		})
	}
}

// TestSaveData тестирование функции сохранения данных
func TestSaveData(t *testing.T) {
	// создаём контроллер
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// создаём объект-заглушку
	mockUsers := mocks.NewMockUserStor(ctrl)
	mockInfo := mocks.NewMockInfoStorage(ctrl)

	// создаем тестовую структуру
	serv := InitServerTest(mockUsers, mockInfo)

	type want struct {
		serv_err error
	}
	CurrTime := time.Now().Format(time.RFC850)

	tests := []struct {
		Name        string
		DataToSend  *pb.SaveDataRequest
		ReverseData servermodels.SaveDataInfo
		UserID      string
		UserExist   bool
		err         error
		want        want
	}{
		{
			Name: "Simple data send",
			DataToSend: &pb.SaveDataRequest{
				StorageID: "StorageID",
				Data:      []byte("Hello there!"),
				DataType:  "TEXT",
				Metainfo:  "second METAINFO updated",
				Time:      CurrTime,
			},
			ReverseData: servermodels.SaveDataInfo{},
			UserID:      "XDOJ6FD32JUYVJJ4",
			UserExist:   true,
			err:         nil,
			want: want{
				serv_err: nil,
			},
		},
		{
			Name: "Data send with error in DB",
			DataToSend: &pb.SaveDataRequest{
				StorageID: "StorageID",
				Data:      []byte("Hello there!"),
				DataType:  "TEXT",
				Metainfo:  "second METAINFO updated",
				Time:      CurrTime,
			},
			ReverseData: servermodels.SaveDataInfo{},
			UserID:      "XDOJ6FD32JUYVJJ4",
			UserExist:   true,
			err:         errors.New("some db error"),
			want: want{
				serv_err: status.Error(codes.Aborted, `Ошибка в сохранении информации о данных`),
			},
		},
		{
			Name: "Data send when newer data in DB",
			DataToSend: &pb.SaveDataRequest{
				StorageID: "StorageID",
				Data:      []byte("Hello there!"),
				DataType:  "TEXT",
				Metainfo:  "second METAINFO updated",
				Time:      CurrTime,
			},
			ReverseData: servermodels.SaveDataInfo{
				UserID:    "XDOJ6FD32JUYVJJ4",
				StorageID: "StorageID",
				MetaInfo:  "NEWER METAINFO updated",
				SaveTime:  time.Now(),
				Type:      "TEXT",
				Data:      []byte("Hello there! UPDATED Not there! Here!!!"),
			},
			UserID:    "XDOJ6FD32JUYVJJ4",
			UserExist: true,
			err:       errors.New("some db error"),
			want: want{
				serv_err: status.Error(codes.Aborted, `Ошибка в сохранении информации о данных`),
			},
		},
		{
			Name: "Data send when user is not registered",
			DataToSend: &pb.SaveDataRequest{
				StorageID: "StorageID",
				Data:      []byte("Hello there!"),
				DataType:  "TEXT",
				Metainfo:  "second METAINFO updated",
				Time:      CurrTime,
			},
			ReverseData: servermodels.SaveDataInfo{},
			UserID:      "",
			UserExist:   false,
			err:         nil,
			want: want{
				serv_err: status.Errorf(codes.Aborted, `Пользователь не авторизован`),
			},
		},
	}
	// создание заглушки контекста для запросов
	ctx := context.Background()

	for _, test := range tests {
		// задаем режим работы моков для сохранения данных
		mockInfo.EXPECT().
			SaveNewData(gomock.Any(), gomock.Any()).
			Return(test.ReverseData, test.err).
			MaxTimes(1)

		t.Run(test.Name, func(t *testing.T) {
			// добавляем к контексту запроса информацию о UserID (заглушка для интерцептора)
			newCtx := context.WithValue(ctx, interceptors.CtxKey,
				servermodels.UserInfo{UserID: test.UserID, Register: test.UserExist})
			// выполняем запрос
			ans, err := serv.SaveData(newCtx, test.DataToSend)
			require.Equal(t, err, test.want.serv_err)
			if err == nil {
				// если ошибка нулевая, то должен быть ответ
				require.NotEmpty(t, ans)
				if ans.IsOutdated {
					// если в ответ были получены новые данные, сравниваем их с ожидаемыми
					require.Equal(t, ans.ReverseData.Data, test.ReverseData.Data)
					require.Equal(t, ans.ReverseData.DataType, test.ReverseData.Type)
					require.Equal(t, ans.ReverseData.Metainfo, test.ReverseData.MetaInfo)
				}
			}

		})
	}
}

// TestUpdateData тестирование функции обновления данных
func TestUpdateData(t *testing.T) {
	// создаём контроллер
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// создаём объект-заглушку
	mockUsers := mocks.NewMockUserStor(ctrl)
	mockInfo := mocks.NewMockInfoStorage(ctrl)

	// создаем тестовую структуру
	serv := InitServerTest(mockUsers, mockInfo)

	type want struct {
		serv_err error
	}
	CurrTime := time.Now().Format(time.RFC850)

	tests := []struct {
		Name       string
		DataToSend *pb.UpdateDataRequest
		UserID     string
		UserExist  bool
		err        error
		want       want
	}{
		{
			Name: "Simple data send",
			DataToSend: &pb.UpdateDataRequest{
				StorageID: "storageID",
				Data:      []byte("Hello there!"),
				DataType:  "TEXT",
				Metainfo:  "METAINFO",
				Time:      CurrTime,
			},
			UserID:    "XDOJ6FD32JUYVJJ4",
			UserExist: true,
			err:       nil,
			want: want{
				serv_err: nil,
			},
		},
		{
			Name: "Data send with error in DB",
			DataToSend: &pb.UpdateDataRequest{
				StorageID: "storageID",
				Data:      []byte("Hello there!"),
				DataType:  "TEXT",
				Metainfo:  "METAINFO",
				Time:      CurrTime,
			},
			UserID:    "XDOJ6FD32JUYVJJ4",
			UserExist: true,
			err:       errors.New("some db error"),
			want: want{
				serv_err: status.Error(codes.Aborted, `Ошибка в обновлении информации о данных`),
			},
		},
		{
			Name: "Data send when user is not registered",
			DataToSend: &pb.UpdateDataRequest{
				StorageID: "storageID",
				Data:      []byte("Hello there!"),
				DataType:  "TEXT",
				Metainfo:  "METAINFO",
				Time:      CurrTime,
			},
			UserID:    "",
			UserExist: false,
			err:       nil,
			want: want{
				serv_err: status.Errorf(codes.Aborted, `Пользователь не авторизован`),
			},
		},
	}
	// создание заглушки контекста для запросов
	ctx := context.Background()

	for _, test := range tests {
		// задаем режим работы моков для обновления данных
		mockInfo.EXPECT().
			UpdateData(gomock.Any(), gomock.Any()).
			Return(test.err).
			MaxTimes(1)

		t.Run(test.Name, func(t *testing.T) {
			// добавляем к контексту запроса информацию о UserID (заглушка для интерцептора)
			newCtx := context.WithValue(ctx, interceptors.CtxKey,
				servermodels.UserInfo{UserID: test.UserID, Register: test.UserExist})
			// выполняем запрос
			_, err := serv.UpdateData(newCtx, test.DataToSend)
			require.Equal(t, err, test.want.serv_err)

		})
	}
}

// TestGetData тестирование функции получения данных
func TestGetData(t *testing.T) {
	// создаём контроллер
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// создаём объект-заглушку
	mockUsers := mocks.NewMockUserStor(ctrl)
	mockInfo := mocks.NewMockInfoStorage(ctrl)

	// создаем тестовую структуру
	serv := InitServerTest(mockUsers, mockInfo)

	type want struct {
		serv_err error
	}

	Time := time.Now()
	CurrTime := Time.Format(time.RFC850)

	tests := []struct {
		Name        string
		ToSend      *pb.GiveDataRequest
		RecieveData servermodels.SaveDataInfo
		UserID      string
		StorageID   string
		UserExist   bool
		err         error
		want        want
	}{
		{
			Name: "Simple data request",
			ToSend: &pb.GiveDataRequest{
				StorageID: "StorageID",
			},
			RecieveData: servermodels.SaveDataInfo{
				Data:     []byte("Hello there!"),
				Type:     "TEXT",
				MetaInfo: "METAINFO",
				SaveTime: Time,
			},
			UserID:    "XDOJ6FD32JUYVJJ4",
			StorageID: "StorageID",
			UserExist: true,
			err:       nil,
			want: want{
				serv_err: nil,
			},
		},
		{
			Name: "Data request with error in DB",
			ToSend: &pb.GiveDataRequest{
				StorageID: "StorageID",
			},
			RecieveData: servermodels.SaveDataInfo{
				Data:     []byte("Hello there!"),
				Type:     "TEXT",
				MetaInfo: "METAINFO",
				SaveTime: Time,
			},
			UserID:    "XDOJ6FD32JUYVJJ4",
			StorageID: "StorageID",
			UserExist: true,
			err:       errors.New("some db error"),
			want: want{
				serv_err: status.Error(codes.Aborted, `Ошибка в получении информации о данных`),
			},
		},
		{
			Name: "Data request when user is not registered",
			ToSend: &pb.GiveDataRequest{
				StorageID: "StorageID",
			},
			RecieveData: servermodels.SaveDataInfo{
				Data:     []byte("Hello there!"),
				Type:     "TEXT",
				MetaInfo: "METAINFO",
				SaveTime: Time,
			},
			UserID:    "",
			StorageID: "StorageID",
			UserExist: false,
			err:       nil,
			want: want{
				serv_err: status.Errorf(codes.Aborted, `Пользователь не авторизован`),
			},
		},
	}
	// создание заглушки контекста для запросов
	ctx := context.Background()

	for _, test := range tests {
		// задаем режим работы моков для получения данных
		mockInfo.EXPECT().
			GetData(gomock.Any(), gomock.Any()).
			Return(test.RecieveData, test.err).
			MaxTimes(1)

		// кодируем данные перед предварительным сохранением
		encData, err := coder.Encrypt(test.RecieveData.Data)
		require.NoError(t, err)
		// предварительно сохраняем данные в тестовое бинарное хранилище
		err = serv.BinStorFunc.SaveBinData(test.UserID, test.StorageID, encData)
		require.NoError(t, err)

		t.Run(test.Name, func(t *testing.T) {
			// добавляем к контексту запроса информацию о UserID (заглушка для интерцептора)
			newCtx := context.WithValue(ctx, interceptors.CtxKey,
				servermodels.UserInfo{UserID: test.UserID, Register: test.UserExist})
			// выполняем запрос
			ans, err := serv.GiveData(newCtx, test.ToSend)
			require.Equal(t, err, test.want.serv_err)
			if err == nil {
				// если ошибка нет, то сравниваем полученные данные
				require.Equal(t, ans.Data, test.RecieveData.Data)
				require.Equal(t, ans.DataType, test.RecieveData.Type)
				require.Equal(t, ans.Metainfo, test.RecieveData.MetaInfo)
				require.Equal(t, ans.Time, CurrTime)
			}

		})
	}
}

// TestGetData тестирование функции получения данных
func TestGetDataList(t *testing.T) {
	// создаём контроллер
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// создаём объект-заглушку
	mockUsers := mocks.NewMockUserStor(ctrl)
	mockInfo := mocks.NewMockInfoStorage(ctrl)

	// создаем тестовую структуру
	serv := InitServerTest(mockUsers, mockInfo)

	type want struct {
		serv_err error
	}

	Time := time.Now()
	CurrTime := Time.Format(time.RFC850)

	tests := []struct {
		Name       string
		ToSend     *pb.GiveDataListRequest
		StagedData []servermodels.SaveDataInfo
		UserID     string
		StorageID  string
		UserExist  bool
		err        error
		want       want
	}{
		{
			Name:   "Simple data list with one line request",
			ToSend: &pb.GiveDataListRequest{},
			StagedData: []servermodels.SaveDataInfo{
				servermodels.SaveDataInfo{
					Data:      []byte("Hello there!"),
					Type:      "TEXT",
					MetaInfo:  "METAINFO",
					SaveTime:  Time,
					StorageID: "StorageID",
				},
			},
			UserID:    "XDOJ6FD32JUYVJJ4",
			StorageID: "StorageID",
			UserExist: true,
			err:       nil,
			want: want{
				serv_err: nil,
			},
		},
		{
			Name:   "Simple data list with several lines request",
			ToSend: &pb.GiveDataListRequest{},
			StagedData: []servermodels.SaveDataInfo{
				servermodels.SaveDataInfo{
					Data:      []byte("Hello there!"),
					Type:      "TEXT",
					MetaInfo:  "METAINFO",
					SaveTime:  Time,
					StorageID: "StorageID1",
				},
				servermodels.SaveDataInfo{
					Data:      []byte("Hello there! Number 2"),
					Type:      "TEXT",
					MetaInfo:  "METAINFO number 2",
					SaveTime:  Time,
					StorageID: "StorageID2",
				},
				servermodels.SaveDataInfo{
					Data:      []byte("Some number 3"),
					Type:      "TEXT",
					MetaInfo:  "Some METAINFO number 3",
					SaveTime:  Time,
					StorageID: "StorageID3",
				},
				servermodels.SaveDataInfo{
					Data:      []byte("Some card info 4"),
					Type:      "CARD",
					MetaInfo:  "Some METAINFO about card number 4",
					SaveTime:  Time,
					StorageID: "StorageID4",
				},
			},
			UserID:    "XDOJ6FD32JUYVJJ4",
			StorageID: "StorageID",
			UserExist: true,
			err:       nil,
			want: want{
				serv_err: nil,
			},
		},
		{
			Name:   "Simple data list with error in DB",
			ToSend: &pb.GiveDataListRequest{},
			StagedData: []servermodels.SaveDataInfo{
				servermodels.SaveDataInfo{
					Data:      []byte("Hello there!"),
					Type:      "TEXT",
					MetaInfo:  "METAINFO",
					SaveTime:  Time,
					StorageID: "StorageID",
				},
			},
			UserID:    "XDOJ6FD32JUYVJJ4",
			StorageID: "StorageID",
			UserExist: true,
			err:       errors.New("some db error"),
			want: want{
				serv_err: status.Error(codes.Aborted, `Ошибка в получении информации о данных`),
			},
		},
		{
			Name:   "Simple data list with unauthorized user",
			ToSend: &pb.GiveDataListRequest{},
			StagedData: []servermodels.SaveDataInfo{
				servermodels.SaveDataInfo{
					Data:      []byte("Hello there!"),
					Type:      "TEXT",
					MetaInfo:  "METAINFO",
					SaveTime:  Time,
					StorageID: "StorageID",
				},
			},
			UserID:    "XDOJ6FD32JUYVJJ4",
			StorageID: "StorageID",
			UserExist: false,
			err:       nil,
			want: want{
				serv_err: status.Errorf(codes.Aborted, `Пользователь не авторизован`),
			},
		},
	}

	// создание заглушки контекста для запросов
	ctx := context.Background()

	for _, test := range tests {
		// задаем режим работы моков для получения информации о данных
		mockInfo.EXPECT().
			GetDataList(gomock.Any()).
			Return(test.StagedData, test.err).
			MaxTimes(1)

		t.Run(test.Name, func(t *testing.T) {
			// добавляем к контексту запроса информацию о UserID (заглушка для интерцептора)
			newCtx := context.WithValue(ctx, interceptors.CtxKey,
				servermodels.UserInfo{UserID: test.UserID, Register: test.UserExist})
			// выполняем запроос
			ans, err := serv.GiveDataList(newCtx, test.ToSend)
			require.Equal(t, err, test.want.serv_err)
			if err == nil {
				// если ошибка нет, то сравниваем полученные данные
				for i, dataLine := range ans.DataList {
					require.Equal(t, dataLine.DataType, test.StagedData[i].Type)
					require.Equal(t, dataLine.Metainfo, test.StagedData[i].MetaInfo)
					require.Equal(t, dataLine.StorageID, test.StagedData[i].StorageID)
					require.Equal(t, dataLine.Time, CurrTime)
				}

			}

		})
	}
}

// TestDeleteData тестирование функции удаления данных
func TestDeleteData(t *testing.T) {
	// создаём контроллер
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// создаём объект-заглушку
	mockUsers := mocks.NewMockUserStor(ctrl)
	mockInfo := mocks.NewMockInfoStorage(ctrl)

	// создаем тестовую структуру
	serv := InitServerTest(mockUsers, mockInfo)

	type want struct {
		serv_err error
	}
	tests := []struct {
		Name       string
		DataToSend *pb.DeleteDataRequest
		UserID     string
		UserExist  bool
		err        error
		want       want
	}{
		{
			Name: "Simple data delete",
			DataToSend: &pb.DeleteDataRequest{
				StorageID: "storageID",
			},
			UserID:    "XDOJ6FD32JUYVJJ4",
			UserExist: true,
			err:       nil,
			want: want{
				serv_err: nil,
			},
		},
		{
			Name: "Data delete with db error",
			DataToSend: &pb.DeleteDataRequest{
				StorageID: "storageID",
			},
			UserID:    "XDOJ6FD32JUYVJJ4",
			UserExist: true,
			err:       errors.New("любая ошибка базы данных"),
			want: want{
				serv_err: status.Error(codes.Aborted, `Ошибка в удалении информации о данных`),
			},
		},
		{
			Name: "Data delete with unregistered user",
			DataToSend: &pb.DeleteDataRequest{
				StorageID: "storageID",
			},
			UserID:    "",
			UserExist: false,
			err:       nil,
			want: want{
				serv_err: status.Errorf(codes.Aborted, `Пользователь не авторизован`),
			},
		},
	}
	// создание заглушки контекста для запросов
	ctx := context.Background()

	for _, test := range tests {
		// задаем режим работы моков для удаления данных
		mockInfo.EXPECT().
			DeleteData(gomock.Any(), gomock.Any()).
			Return(test.err).
			MaxTimes(1)

		t.Run(test.Name, func(t *testing.T) {
			// добавляем к контексту запроса информацию о UserID (заглушка для интерцептора)
			newCtx := context.WithValue(ctx, interceptors.CtxKey,
				servermodels.UserInfo{UserID: test.UserID, Register: test.UserExist})
			// выполняем запрос
			_, err := serv.DeleteData(newCtx, test.DataToSend)
			require.Equal(t, err, test.want.serv_err)

		})
	}
}
