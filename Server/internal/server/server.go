package server

import (
	"crypto/tls"
	"fmt"
	"log"
	"net"

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
	proto.UnimplementedKeeperServer
}

func InitServer() *Server {
	var err error
	Serv := new(Server)

	Serv.Config = config.InitConf()

	logger.Initialize(Serv.Config.InFileLog)
	logger.Log.Info("Logger initialyzed!",
		zap.Bool("In file mode:", Serv.Config.InFileLog),
	)

	Serv.DBusers, err = pgstor.DBUsersInit(Serv.Config.DBAdr)
	if err != nil {
		logger.Log.Error("Error in creating users DB", zap.Error(err))
		return Serv
	}
	Serv.UserStor = Serv.DBusers

	Serv.DBData, err = pgstor.DBDataInit(Serv.Config.DBAdr)
	if err != nil {
		logger.Log.Error("Error in creating data DB", zap.Error(err))
		return Serv
	}
	Serv.InfoStor = Serv.DBData

	Serv.BinStor = binstor.NewBinStor(Serv.Config.StorageDir)
	Serv.BinStorFunc = Serv.BinStor

	return Serv
}

func RunServer() error {
	server := InitServer()

	cert, err := tls.LoadX509KeyPair("server.crt", "server.key")
	if err != nil {
		log.Fatalf("Failed to load server certificate: %v", err)
	}

	// Настройки TLS
	config := &tls.Config{
		Certificates: []tls.Certificate{cert},
	}

	creds := credentials.NewTLS(config)

	fmt.Println("Host addr", server.Config.HostAddr)
	// определяем адрес для сервера
	listen, err := net.Listen("tcp", server.Config.HostAddr)
	if err != nil {
		return err
	}
	// создаём gRPC-сервер без зарегистрированной службы
	interceptors := grpc.ChainUnaryInterceptor(interceptors.AtuhInterceptor, interceptors.LoggingInterceptor)
	s := grpc.NewServer(grpc.Creds(creds), interceptors)

	proto.RegisterKeeperServer(s, server)

	// получаем запрос gRPC
	if err := s.Serve(listen); err != nil {
		log.Fatal(err)
	}
	return nil
}
