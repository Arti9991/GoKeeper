package server

import (
	"context"
	"fmt"

	pb "github.com/Arti9991/GoKeeper/server/internal/server/proto"
)

// GetAddr получение исходного URL по укороченному
func (s *Server) SaveData(ctx context.Context, in *pb.SaveDataRequset) (*pb.SaveDataResponse, error) {
	var res pb.SaveDataResponse

	fmt.Println("Input ID", in.Id)
	fmt.Println("Input Data", in.Metainfo)

	return &res, nil
}
