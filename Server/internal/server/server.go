package server

import (
	"context"
	"crypto/tls"
	"net"
	"os"
	"os/signal"
	"syscall"

	"github.com/Arti9991/GoKeeper/server/internal/config"
	"github.com/Arti9991/GoKeeper/server/internal/logger"
	"github.com/Arti9991/GoKeeper/server/internal/server/interceptors"
	"github.com/Arti9991/GoKeeper/server/internal/server/proto"
	"github.com/Arti9991/GoKeeper/server/internal/storage/binstor"
	"github.com/Arti9991/GoKeeper/server/internal/storage/pgstor"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

type Server struct {
	// структура с инфомрацией о сервере
	DBusers     *pgstor.DBUsersStor
	UserStor    pgstor.UserStor
	DBData      *pgstor.DBStor
	InfoStor    pgstor.InfoStorage
	BinStor     *binstor.BinStor
	BinStorFunc binstor.BinStrorFunc
	Config      config.Config
	Ctx         context.Context
	// встраиваем proto
	proto.UnimplementedKeeperServer
}

// InitServer функция инициализации сервера со всеми основными параметрами
func InitServer(ctx context.Context) (*Server, error) {
	var err error
	Serv := new(Server)

	Serv.Ctx = ctx
	Serv.Config = config.InitConf()

	logger.Initialize(Serv.Config.InFileLog)
	logger.Log.Info("Logger initialyzed!",
		zap.Bool("In file mode:", Serv.Config.InFileLog),
	)

	Serv.DBusers, err = pgstor.DBUsersInit(Serv.Config.DBAdr)
	if err != nil {
		logger.Log.Error("Error in creating users DB", zap.Error(err))
		return Serv, err
	}
	Serv.UserStor = Serv.DBusers

	Serv.DBData, err = pgstor.DBDataInit(Serv.Config.DBAdr)
	if err != nil {
		logger.Log.Error("Error in creating data DB", zap.Error(err))
		return Serv, err
	}
	Serv.InfoStor = Serv.DBData

	Serv.BinStor, err = binstor.NewBinStor(Serv.Config.StorageDir)
	if err != nil {
		logger.Log.Error("Error in creating binary storage", zap.Error(err))
		return Serv, err
	}
	Serv.BinStorFunc = Serv.BinStor

	return Serv, nil
}

// RunServer функция запуска сервера
func RunServer() error {

	// канал для сообщения о Shutdown
	shutCh := make(chan struct{})
	// контекст для ожидания системного сигнала на завершение работы
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()
	// инциализация сервера
	server, err := InitServer(ctx)
	if err != nil {
		logger.Log.Error("Error in InitServer")
		return err
	}
	// получение сертификатов для secure mode
	cert, err := tls.LoadX509KeyPair("server.crt", "server.key")
	if err != nil {
		logger.Log.Error("Error in load TLS files for secure mode", zap.Error(err))
		return err
	}

	// Настройки TLS
	config := &tls.Config{
		Certificates: []tls.Certificate{cert},
	}

	creds := credentials.NewTLS(config)

	// определяем адрес для сервера и слушаем его
	listen, err := net.Listen("tcp", server.Config.HostAddr)
	if err != nil {
		logger.Log.Error("Error in listening port!", zap.String("port:", server.Config.HostAddr), zap.Error(err))
		return err
	}
	// инициализация интерцепторов для логгирования и авторизации
	interceptors := grpc.ChainUnaryInterceptor(interceptors.AtuhInterceptor, interceptors.LoggingInterceptor)
	// создаём gRPC-сервер с интерцепторами и протоколами безопасности
	s := grpc.NewServer(grpc.Creds(creds), interceptors)
	// регистрация обработчика
	proto.RegisterKeeperServer(s, server)

	logger.Log.Info("Server initialyzed!", zap.String("Address:", server.Config.HostAddr))
	// запуск горутины на ожидание сигнала о выключении
	WaitShutdown(server, s, shutCh)

	// ожидание запросов в gRPC
	if err := s.Serve(listen); err != nil {
		logger.Log.Error("Error in gRPC Serve!", zap.Error(err))
		return err
	}
	// ожидание сообщения о Shutdown
	<-shutCh
	logger.Log.Info("Server stopped!")
	return nil
}

// WaitShutdown функция с горутиной, ожидающей сигнала об отключении сервера и
// выполняющей GracefulStop() для gRPC сервера
func WaitShutdown(server *Server, s *grpc.Server, shutCh chan struct{}) {
	go func() {
		var err error
		<-server.Ctx.Done()
		// получили сигнал os.Interrupt, запускаем процедуру graceful stop
		logger.Log.Info("Graceful shutdown...")

		s.GracefulStop()
		// отключаем соединение с базами данных
		err = server.DBData.DB.Close()
		if err != nil {
			logger.Log.Error("Error in closing datainfo Db", zap.Error(err))
		}
		server.DBusers.DB.Close()
		if err != nil {
			logger.Log.Error("Error in closing datainfo Db", zap.Error(err))
		}
		// сообщение о Shutdown
		close(shutCh)
	}()
}
