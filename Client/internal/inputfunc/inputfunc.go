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
		fmt.Printf("Введите номер карты: ")
		// открываем потоковое чтение из консоли
		reader := bufio.NewReader(os.Stdin)
		// читаем строку из консоли
		num, err := reader.ReadString('\n')
		fmt.Printf("Введите срок действия карты: ")
		date, err := reader.ReadString('\n')
		fmt.Printf("Введите CVV карты: ")
		CVV, err := reader.ReadString('\n')
		fmt.Printf("Введите владельца карты: ")
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
		fmt.Printf("Введите логин: ")
		// открываем потоковое чтение из консоли
		reader := bufio.NewReader(os.Stdin)
		// читаем строку из консоли
		log, err := reader.ReadString('\n')
		fmt.Printf("Введите пароль: ")
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
		fmt.Printf("\nВведите текст для сохранения: ")
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
	case "BINARY":
		fmt.Printf("\nВведите путь к файлу для загрузки: ")
		// открываем потоковое чтение из консоли
		reader := bufio.NewReader(os.Stdin)
		// читаем строку из консоли
		path, err := reader.ReadString('\n')
		path = strings.TrimSuffix(path, "\n")
		path = strings.TrimSuffix(path, "\r")
		fmt.Printf("%#v\n", path)

		data, err := os.ReadFile(path)
		if err != nil {
			fmt.Println(err)
			return nil, clientmodels.ErrorInput
		}
		if len(data) > 4194304 {
			return nil, clientmodels.ErrBigData
		}
		return data, nil
	default:
		return nil, clientmodels.ErrorInput
	}
}

func ParceAnswer(Data []byte, storageID string, Type string, Metainfo string) error {

	switch Type {
	case "CARD":
		//fmt.Println(Data)
		//fmt.Println(buff2)
		fmt.Printf("Information: %s\n", Metainfo)

		out := clientmodels.CardInfo{}
		dec := gob.NewDecoder(bytes.NewBuffer(Data))
		err := dec.Decode(&out)
		if err != nil {
			return err
		}
		fmt.Printf("Card numver: %s\n", out.Number)
		fmt.Printf("Card CVV: %s\n", out.CVVcode)
		fmt.Printf("Card exipre: %s\n", out.ExpDate)
		fmt.Printf("Card holder: %s\n", out.Holder)
		return nil

	case "AUTH":
		fmt.Printf("Information: %s\n", Metainfo)
		fmt.Println(Metainfo)

		out := clientmodels.LoginInfo{}
		dec := gob.NewDecoder(bytes.NewBuffer(Data))
		err := dec.Decode(&out)
		if err != nil {
			return err
		}
		fmt.Printf("Login: %s\n", out.Login)
		fmt.Printf("Password: %s\n", out.Password)
		return nil
	case "TEXT":
		fmt.Printf("Information: %s\n", Metainfo)
		fmt.Println(Metainfo)

		fmt.Printf("Saved text info: \n%s\n", string(Data))
		return nil
	case "BINARY":
		fmt.Printf("\nВведите путь к файлу для сохранения: ")
		// открываем потоковое чтение из консоли
		reader := bufio.NewReader(os.Stdin)
		// читаем строку из консоли
		path, err := reader.ReadString('\n')
		path = strings.TrimSuffix(path, "\n")
		path = strings.TrimSuffix(path, "\r")
		fmt.Printf("%#v\n", path)

		err = os.WriteFile(path, Data, 0644)
		if err != nil {
			fmt.Println(err)
			return clientmodels.ErrBigCantSave
		}

		fmt.Printf("Файл сохранен в: %s\n", path)
		return nil
	}
	return nil
}
