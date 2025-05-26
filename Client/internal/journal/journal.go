package journal

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/Arti9991/GoKeeper/client/internal/clientmodels"
)

func JournalSave(JrInf clientmodels.JournalInfo) error {
	var SaveInf []byte
	file, err := os.OpenFile(clientmodels.JournalFile, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		//logger.Log.Error("SAVE Error in opening file", zap.Error(err))
		fmt.Println(err)
		return err
	}
	defer file.Close()

	SaveInf, err = json.Marshal(JrInf)
	if err != nil {
		fmt.Println(err)
		return err
	}
	SaveInf = append(SaveInf, byte('\n'))

	n, err := file.Write(SaveInf)
	if err != nil || n == 0 {
		fmt.Println(err)
		return err
	}
	return nil
}
