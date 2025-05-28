package binstor

import (
	"fmt"
	"os"
	"sync"
)

// type BinStrorFunc interface {
// 	SaveBinData(userID string, storageID string, binData []byte) error
// 	UpdateBinData(userID string, storageID string, binData []byte) error
// 	GetBinData(userID string, storageID string) ([]byte, error)
// 	RemoveBinData(userID string, storageID string) error
// }

type BinStor struct {
	MainStor   map[string][]byte
	mu         sync.Mutex
	StorageDir string
}

func NewBinStor(StorageDir string) *BinStor {
	MainSt := make(map[string][]byte)

	err := os.Mkdir(StorageDir, 0644)
	if err != nil {
		//logger.Log.Error("Error in creating directory", zap.Error(err))
	}
	return &BinStor{MainStor: MainSt, StorageDir: StorageDir}
}

func (s *BinStor) SaveBinData(storageID string, binData []byte) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.MainStor[storageID] = binData

	// Также сохраняем данные на диск
	file, err := os.OpenFile(s.StorageDir+storageID, os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		//logger.Log.Error("SAVE Error in opening file", zap.Error(err))
		return err
	}
	defer file.Close()

	n, err := file.Write(binData)
	if err != nil || n == 0 {
		//logger.Log.Error("Error in saving to file", zap.Error(err))
		return err
	}

	return nil
}

func (s *BinStor) UpdateBinData(storageID string, binData []byte) error {

	fmt.Println("UpdateBinData")
	// Также сохраняем данные на диск
	file, err := os.OpenFile(s.StorageDir+storageID, os.O_WRONLY|os.O_TRUNC, 0644)
	if err != nil {
		//logger.Log.Error("UPDATE Error in opening file", zap.Error(err))
		return err
	}
	defer file.Close()

	n, err := file.Write(binData)
	if err != nil || n == 0 {
		//logger.Log.Error("Error in updating file", zap.Error(err))
		return err
	}

	return nil
}

func (s *BinStor) GetBinData(storageID string) ([]byte, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	var out []byte
	var err error
	out, ok := s.MainStor[storageID]
	if !ok {
		out, err = os.ReadFile(s.StorageDir + storageID)
		if err != nil {
			//logger.Log.Error("Error in reading from file", zap.Error(err))
			return nil, err
		}
	}

	return out, nil
}

func (s *BinStor) RemoveBinData(storageID string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	delete(s.MainStor, storageID)

	err := os.Remove(s.StorageDir + storageID)
	if err != nil {
		//logger.Log.Error("Error in deleting file", zap.Error(err))
		return err
	}
	return nil
}
