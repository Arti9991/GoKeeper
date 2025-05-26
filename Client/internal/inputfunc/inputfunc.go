package inputfunc

import (
	"bufio"
	"bytes"
	"encoding/gob"
	"fmt"
	"os"
	"strings"

	"github.com/Arti9991/GoKeeper/client/internal/clientmodels"
)

func ParceInput(Type string) ([]byte, error) {
	var data []byte

	switch Type {
	case "CARD":
		fmt.Printf("\nВведите данные карты\n")
		fmt.Printf("Введите номер карты:")
		// открываем потоковое чтение из консоли
		reader := bufio.NewReader(os.Stdin)
		// читаем строку из консоли
		num, err := reader.ReadString('\n')
		fmt.Printf("Введите срок действия карты:")
		date, err := reader.ReadString('\n')
		fmt.Printf("Введите CVV карты:")
		CVV, err := reader.ReadString('\n')
		fmt.Printf("Введите владельца карты:")
		name, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println(err)
			return nil, clientmodels.ErrorInput
		}

		num = strings.TrimSuffix(num, "\n")
		date = strings.TrimSuffix(date, "\n")
		CVV = strings.TrimSuffix(CVV, "\n")
		name = strings.TrimSuffix(name, "\n")

		struc := clientmodels.CardInfo{
			Number:  num,
			ExpDate: date,
			CVVcode: CVV,
			Holder:  name,
		}
		var buf bytes.Buffer
		enc := gob.NewEncoder(&buf)
		err = enc.Encode(struc)
		if err != nil {
			fmt.Println(err)
			return nil, err
		}

		data = buf.Bytes()

		return data, nil

	case "AUTH":
		fmt.Printf("\nВведите для авторизации\n")
		fmt.Printf("Введите логин")
		// открываем потоковое чтение из консоли
		reader := bufio.NewReader(os.Stdin)
		// читаем строку из консоли
		log, err := reader.ReadString('\n')
		fmt.Printf("Введите пароль")
		passw, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println(err)
			return nil, clientmodels.ErrorInput
		}

		log = strings.TrimSuffix(log, "\n")
		passw = strings.TrimSuffix(passw, "\n")

		struc := clientmodels.LoginInfo{
			Login:    log,
			Password: passw,
		}

		var buf bytes.Buffer
		enc := gob.NewEncoder(&buf)
		err = enc.Encode(struc)
		if err != nil {
			fmt.Println(err)
			return nil, err
		}

		data = buf.Bytes()
		fmt.Println(data)
		return data, nil

	case "TEXT":
		fmt.Printf("\nВведите текст для сохранения:")
		// открываем потоковое чтение из консоли
		reader := bufio.NewReader(os.Stdin)
		// читаем строку из консоли
		txt, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println(err)
			return nil, clientmodels.ErrorInput
		}

		data = []byte(txt)

		return data, nil
	default:
		return nil, clientmodels.ErrorInput
	}
}
