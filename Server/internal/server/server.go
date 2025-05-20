package server

import (
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
)

type Server struct {
	// структура с инфомрацией о сервере
	DBusers *pgstor.DBUsersStor
	DBData  *pgstor.DBStor
	BinStor *binstor.BinStor
	Config  config.Config
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

	Serv.DBData, err = pgstor.DBDataInit(Serv.Config.DBAdr)
	if err != nil {
		logger.Log.Error("Error in creating data DB", zap.Error(err))
		return Serv
	}

	Serv.BinStor = binstor.NewBinStor()

	return Serv
}

func RunServer() error {
	server := InitServer()

	fmt.Println("Host addr", server.Config.HostAddr)
	// определяем адрес для сервера
	listen, err := net.Listen("tcp", server.Config.HostAddr)
	if err != nil {
		return err
	}
	// создаём gRPC-сервер без зарегистрированной службы
	interceptors := grpc.ChainUnaryInterceptor(interceptors.AtuhInterceptor, interceptors.LoggingInterceptor)
	s := grpc.NewServer(interceptors)

	proto.RegisterKeeperServer(s, server)

	// получаем запрос gRPC
	if err := s.Serve(listen); err != nil {
		log.Fatal(err)
	}
	return nil
}
