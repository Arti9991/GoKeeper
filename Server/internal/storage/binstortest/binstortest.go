package binstortest

import (
	"sync"
)

// структура для тестового хранилища бинарных данных
// (без сохранения на диск)
type BinStorTest struct {
	MainStor map[string][]byte
	mu       sync.Mutex
}

// NewBinStorTest инциализация тестового хранилища бинарных данных
func NewBinStorTest() *BinStorTest {
	MainSt := make(map[string][]byte)

	return &BinStorTest{MainStor: MainSt}
}

// SaveBinData функция сохранения бинарных данных
func (s *BinStorTest) SaveBinData(userID string, storageID string, binData []byte) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.MainStor[storageID] = binData
	return nil
}

// UpdateBinData функция обновления бинарных данных
func (s *BinStorTest) UpdateBinData(userID string, storageID string, binData []byte) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.MainStor[storageID] = binData

	return nil
}

// GetBinData функция получения бинарных данных
func (s *BinStorTest) GetBinData(userID string, storageID string) ([]byte, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	var out []byte
	out = s.MainStor[storageID]
	return out, nil
}

// RemoveBinData функция удаления бинарных данных
func (s *BinStorTest) RemoveBinData(userID string, storageID string) error {
	return nil
}
