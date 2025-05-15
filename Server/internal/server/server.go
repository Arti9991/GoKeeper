package server

import (
	"fmt"
	"log"
	"net"

	"github.com/Arti9991/GoKeeper/server/internal/config"
	"github.com/Arti9991/GoKeeper/server/internal/logger"
	"github.com/Arti9991/GoKeeper/server/internal/server/proto"
	"go.uber.org/zap"
	"google.golang.org/grpc"
)

type Server struct {
	// структура с инфомрацией о сервере
	Config config.Config
	proto.UnimplementedKeeperServer
}

func InitServer() *Server {
	Serv := new(Server)

	Serv.Config = config.InitConf()

	logger.Initialize(Serv.Config.InFileLog)
	logger.Log.Info("Logger initialyzed!",
		zap.Bool("In file mode:", Serv.Config.InFileLog),
	)

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
	//interceptors := grpc.ChainUnaryInterceptor(atuhInterceptor, loggingInterceptor)
	s := grpc.NewServer()

	proto.RegisterKeeperServer(s, server)

	// получаем запрос gRPC
	if err := s.Serve(listen); err != nil {
		log.Fatal(err)
	}
	return nil
}
