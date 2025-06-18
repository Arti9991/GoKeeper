package binstor

import (
	"os"
)

// структура с информацией о бинарном хранилище
type BinStor struct {
	StorageDir string
}

// NewBinStor инициализация бинарного хранилища
func NewBinStor(StorageDir string) *BinStor {

	err := os.Mkdir(StorageDir, 0644)
	if err != nil {
		//logger.Log.Error("Error in creating directory", zap.Error(err))
	}
	return &BinStor{StorageDir: StorageDir}
}

// SaveBinData функция сохранения данных в бинарное хранилище
func (s *BinStor) SaveBinData(storageID string, binData []byte) error {
	// открываем\создаем файл для записи
	file, err := os.OpenFile(s.StorageDir+storageID, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		return err
	}
	defer file.Close()
	// записываем данные
	n, err := file.Write(binData)
	if err != nil || n == 0 {
		return err
	}

	return nil
}

// SaveBinData функция обновления данных в бинарном хранилище
func (s *BinStor) UpdateBinData(storageID string, binData []byte) error {

	// открываем\создаем файл для записи
	file, err := os.OpenFile(s.StorageDir+storageID, os.O_WRONLY|os.O_TRUNC, 0644)
	if err != nil {
		return err
	}
	defer file.Close()
	// записываем данные
	n, err := file.Write(binData)
	if err != nil || n == 0 {
		return err
	}

	return nil
}

// GetBinData функция получения данных из бинарного хранилища
func (s *BinStor) GetBinData(storageID string) ([]byte, error) {

	var out []byte
	var err error
	// читаем файл из хранилища
	out, err = os.ReadFile(s.StorageDir + storageID)
	if err != nil {
		return nil, err
	}

	return out, nil
}

// RemoveBinData функция удаления данных из бинарного хранилища
func (s *BinStor) RemoveBinData(storageID string) error {
	err := os.Remove(s.StorageDir + storageID)
	if err != nil {
		return err
	}
	return nil
}
