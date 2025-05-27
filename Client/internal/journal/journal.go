package journal

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
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

func JournalGet() ([]clientmodels.JournalInfo, error) {
	var JourMass []clientmodels.JournalInfo
	var err error

	file, err := os.OpenFile(clientmodels.JournalFile, os.O_RDONLY, 0644)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	defer file.Close()

	reader := bufio.NewReader(file)
	dec := json.NewDecoder(reader) // (1)
	for {
		var Jour clientmodels.JournalInfo
		err := dec.Decode(&Jour) // (2)
		if err == io.EOF {
			break
		}
		if err != nil {
			fmt.Println(err)
			return nil, err
		}
		if !Jour.Sync {
			JourMass = append(JourMass, Jour)
		}
	}

	return JourMass, nil
}
