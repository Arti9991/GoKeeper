package binstortest

import (
	"sync"
)

type BinStorTest struct {
	MainStor map[string][]byte
	mu       sync.Mutex
}

func NewBinStorTest() *BinStorTest {
	MainSt := make(map[string][]byte)

	return &BinStorTest{MainStor: MainSt}
}
func (s *BinStorTest) SaveBinData(userID string, storageID string, binData []byte) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.MainStor[storageID] = binData
	return nil
}

func (s *BinStorTest) UpdateBinData(userID string, storageID string, binData []byte) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.MainStor[storageID] = binData

	return nil
}

func (s *BinStorTest) GetBinData(userID string, storageID string) ([]byte, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	var out []byte
	out = s.MainStor[storageID]
	return out, nil
}

func (s *BinStorTest) RemoveBinData(userID string, storageID string) error {
	return nil
}
