package binstor

import (
	"os"
	"sync"

	"github.com/Arti9991/GoKeeper/server/internal/logger"
	"go.uber.org/zap"
)

type BinStrorFunc interface {
	SaveBinData(userID string, storageID string, binData []byte) error
	UpdateBinData(userID string, storageID string, binData []byte) error
	GetBinData(userID string, storageID string) ([]byte, error)
	RemoveBinData(userID string, storageID string) error
}

type BinStor struct {
	MainStor   map[string][]byte
	mu         sync.Mutex
	StorageDir string
}

func NewBinStor(StorageDir string) *BinStor {
	MainSt := make(map[string][]byte)

	err := os.Mkdir(StorageDir, 0644)
	if err != nil {
		logger.Log.Error("Error in creating directory", zap.Error(err))
	}
	return &BinStor{MainStor: MainSt, StorageDir: StorageDir}
}

func (s *BinStor) SaveBinData(userID string, storageID string, binData []byte) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.MainStor[storageID] = binData

	// Также сохраняем данные на диск
	//err := os.WriteFile(s.StorageDir+userID+"_"+storageID, binData, 0644)
	file, err := os.OpenFile(s.StorageDir+userID+"_"+storageID, os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		logger.Log.Error("SAVE Error in opening file", zap.Error(err))
		return err
	}
	n, err := file.Write(binData)
	if err != nil || n == 0 {
		logger.Log.Error("Error in saving to file", zap.Error(err))
		return err
	}

	return nil
}

func (s *BinStor) UpdateBinData(userID string, storageID string, binData []byte) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.MainStor[storageID] = binData

	// Также сохраняем данные на диск
	file, err := os.OpenFile(s.StorageDir+userID+"_"+storageID, os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		logger.Log.Error("UPDATE Error in opening file", zap.Error(err))
		return err
	}
	n, err := file.Write(binData)
	if err != nil || n == 0 {
		logger.Log.Error("Error in updating file", zap.Error(err))
		return err
	}

	return nil
}

func (s *BinStor) GetBinData(userID string, storageID string) ([]byte, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	var out []byte
	out, ok := s.MainStor[storageID]
	if !ok {
		file, err := os.OpenFile(s.StorageDir+userID+"_"+storageID, os.O_RDONLY, 0644)
		if err != nil {
			logger.Log.Error("GET Error in opening file", zap.Error(err))
			return nil, err
		}
		n, err := file.Read(out)
		if err != nil || n == 0 {
			logger.Log.Error("Error in reading from file", zap.Error(err))
			return nil, err
		}
		return out, nil
	}

	return out, nil
}

func (s *BinStor) RemoveBinData(userID string, storageID string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	delete(s.MainStor, storageID)

	err := os.Remove(s.StorageDir + userID + "_" + storageID)
	if err != nil {
		logger.Log.Error("Error in deleting file", zap.Error(err))
		return err
	}
	return nil
}
