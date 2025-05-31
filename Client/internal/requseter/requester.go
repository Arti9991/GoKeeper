package requseter

import (
	"bufio"
	"context"
	"crypto/x509"
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/Arti9991/GoKeeper/client/internal/binstor"
	"github.com/Arti9991/GoKeeper/client/internal/clientmodels"
	"github.com/Arti9991/GoKeeper/client/internal/dbstor"
	pb "github.com/Arti9991/GoKeeper/client/internal/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

// структура со всеми данными для запросов
type ReqStruct struct {
	ServAddr string
	BinStor  *binstor.BinStor
	DBStor   *dbstor.DBStor
	Creds    credentials.TransportCredentials
}

// NewRequester функция для инциализации структуры для запросов
func NewRequester(addr string) (*ReqStruct, error) {
	var err error
	ReqStruct := new(ReqStruct)
	ReqStruct.ServAddr = addr

	// инициализируем бинарное хранилище
	ReqStruct.BinStor = binstor.NewBinStor(clientmodels.StorageDir)
	// инициализируем базу данных с информацией о данных
	ReqStruct.DBStor, err = dbstor.DbInit("Journal.db")
	if err != nil {
		return ReqStruct, err
	}

	// Загружаем сертификат, которому доверяем (тот, что сгенерирован на сервере)
	caCert, err := os.ReadFile("server.crt")
	if err != nil {
		return ReqStruct, err
	}

	// Создаём пул корневых сертификатов и добавляем туда server.crt
	certPool := x509.NewCertPool()
	if !certPool.AppendCertsFromPEM(caCert) {
		return ReqStruct, err
	}

	// Настраиваем TLS
	ReqStruct.Creds = credentials.NewClientTLSFromCert(certPool, "localhost") // CN должен совпадать с /CN= в server.crt

	return ReqStruct, nil
}

// LoginRequest метод авторизации пользователя
func (req *ReqStruct) LoginRequest(Login string, Password string) error {
	// создаем контекст
	ctx := context.Background()
	// открывем соединение
	dial, err := grpc.NewClient(req.ServAddr, grpc.WithTransportCredentials(req.Creds)) //req.ServAddr
	if err != nil {
		return err
	}
	// инициализируем клиент и вызываем метод авторизации
	r := pb.NewKeeperClient(dial)
	ans, err := r.Loginuser(ctx, &pb.LoginRequest{
		UserLogin:    Login,
		UserPassword: Password,
	})
	if err != nil {
		return err
	}
	// записываем полученый токен в файл
	file, err := os.OpenFile("./Token.txt", os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	n, err := file.Write([]byte(ans.UserID + "\n"))
	if err != nil || n == 0 {
		return err
	}
	fmt.Printf("\nАвторизация успешна! Пользователь: %s \t", Login)
	return nil
}

// RegisterRequest метод регистрации нового пользователя
func (req *ReqStruct) RegisterRequest(Login string, Password string) error {
	// проверяем корректность логина и пароля
	if len([]rune(Login)) < 4 {
		return clientmodels.ErrBadLogin
	}
	if len([]rune(Password)) < 6 {
		return clientmodels.ErrBadPassowrd
	}

	// создаем контекст
	ctx := context.Background()
	// открывем соединение
	dial, err := grpc.NewClient(req.ServAddr, grpc.WithTransportCredentials(req.Creds)) //req.ServAddr
	if err != nil {
		return err
	}
	// инициализируем клиент и вызываем метод регистрации
	r := pb.NewKeeperClient(dial)
	ans, err := r.RegisterUser(ctx, &pb.RegisterRequest{
		UserLogin:    Login,
		UserPassword: Password,
	})
	if err != nil {
		return err
	}
	// записываем полученый токен в файл
	file, err := os.OpenFile("./Token.txt", os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	n, err := file.Write([]byte(ans.UserID + "\n"))
	if err != nil || n == 0 {
		return err
	}
	fmt.Printf("\nРегистрация успешна! Пользователь: %s \t", Login)
	return nil
}

// LogoutRequest метод выхода пользователя из системы
// (удаляет токен сессии, все файлы из локального хранилища
// и информацию о них)
func (req *ReqStruct) LogoutRequest() error {
	var err error

	// спрашиваем подтверждение на выход из аккаунта
	fmt.Println("Вы точно хотите выйти из аккаунта? [y/n] (Все несинхронизированные данные будут далены!)")
	reader := bufio.NewReader(os.Stdin)
	// читаем путь из консоли
	agree, err := reader.ReadString('\n')
	agree = strings.ToLower(agree)
	if strings.Contains(agree, "n") {
		return clientmodels.ErrUserAbort
	}
	// обновляем таблицу
	err = req.DBStor.ReinitTable()
	if err != nil {
		return err
	}
	// чистим бинерное хранилище
	err = os.RemoveAll(clientmodels.StorageDir)
	if err != nil {
		return err
	}
	err = os.Mkdir(clientmodels.StorageDir, 0644)
	// удаляем файл токена
	err = os.Remove("./Token.txt")
	if err != nil {
		return err
	}
	fmt.Printf("\nВсе данные пользователя очищены!\n")
	return nil
}

// SyncRequest метод синхронизации локальных данных с сервером
// (по флагу синхронизации отправляет запросы
// на сохранение данных на сервере)
func (req *ReqStruct) SyncRequest(offlineMode bool) error {
	// проверяем флаг офлайн мода
	if offlineMode {
		return errors.New("cannot sync in offline mode")
	}
	// считываем токен из файла
	var UserID string
	file, err := os.Open(clientmodels.TokenFile)
	if err != nil {
		return err
	}
	defer file.Close()

	reader := bufio.NewReader(file)
	UserID, err = reader.ReadString('\n')
	if err != nil {
		return err
	}
	// удаляем суффикс
	UserID = strings.TrimSuffix(UserID, "\n")
	// получаем список данных на синхронизацию
	SnData, err := req.DBStor.GetForSync()
	if err != nil {
		return err
	}

	// проходим в цикле по списку и обновляем данные
	for _, syncPart := range SnData {

		syncPart.Data, err = req.BinStor.GetBinData(syncPart.StorageID)
		if err != nil {
			return err
		}
		err = SendWithUpdate(syncPart.StorageID, syncPart, req, syncPart.Data)
		if err != nil {
			return err
		}
	}
	fmt.Printf("\nДанные синхронизированы!\n")
	return nil
}
