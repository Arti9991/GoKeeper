package binstor

import (
	"os"
	"sync"

	"github.com/Arti9991/GoKeeper/server/internal/logger"
	"go.uber.org/zap"
)

var (
	StorageDir = "D:/Course/Practicum/GoKeeper/Server/storage/"
)

type BinStor struct {
	MainStor map[string][]byte
	mu       sync.Mutex
}

func NewBinStor() *BinStor {
	MainSt := make(map[string][]byte)

	err := os.Mkdir(StorageDir, 0644)
	if err != nil {
		logger.Log.Error("Error in creating directory", zap.Error(err))
	}
	return &BinStor{MainStor: MainSt}
}

func (s *BinStor) SaveBinData(userID string, storageID string, binData []byte) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.MainStor[storageID] = binData

	// Также сохраняем данные на диск
	err := os.WriteFile(StorageDir+userID+"_"+storageID, binData, 0644)
	if err != nil {
		logger.Log.Error("Error in saving file", zap.Error(err))
	}
	return nil
}

func (s *BinStor) GetBinData(storageID string) ([]byte, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	var out []byte
	out, ok := s.MainStor[storageID]
	if !ok {
		out, err := os.ReadFile(StorageDir + storageID)
		if err != nil {
			logger.Log.Error("Error in reading file", zap.Error(err))
		}
		return out, nil
	}

	return out, nil
}
