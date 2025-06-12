package inputfunc

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/Arti9991/GoKeeper/client/internal/clientmodels"
	"google.golang.org/protobuf/proto"

	//"github.com/Arti9991/GoKeeper/client/internal/proto"
	pb "github.com/Arti9991/GoKeeper/client/internal/proto"
)

// ParceInput функция для парсинга входных данных
func ParceInput(Type string) ([]byte, error) {
	var data []byte

	switch Type {
	case "CARD":
		// просим ввести данные с карты
		fmt.Printf("\nВведите данные карты\n")
		fmt.Printf("Введите номер карты: ")
		// открываем потоковое чтение из консоли
		reader := bufio.NewReader(os.Stdin)
		// считывваем все данные из консоли
		num, err := reader.ReadString('\n')
		if len([]rune(num)) < 8 {
			// номер карты слишком короткий
			return nil, clientmodels.ErrShortNum
		}
		fmt.Printf("Введите срок действия карты: ")
		date, err := reader.ReadString('\n')
		fmt.Printf("Введите CVV карты: ")
		CVV, err := reader.ReadString('\n')
		if len([]rune(CVV)) != 5 {
			// CVV не равен трем символам
			return nil, clientmodels.ErrBadCVV
		}
		fmt.Printf("Введите владельца карты: ")
		name, err := reader.ReadString('\n')
		if err != nil {
			return nil, clientmodels.ErrorInput
		}
		// убираем лишние суффиксы после ввода
		num = strings.TrimSuffix(num, "\n")
		date = strings.TrimSuffix(date, "\n")
		CVV = strings.TrimSuffix(CVV, "\n")
		name = strings.TrimSuffix(name, "\n")
		// сериализуем данные в proto формат
		struc := &pb.CardInfo{
			Number:  num,
			ExpDate: date,
			CVVcode: CVV,
			Holder:  name,
		}
		data, err = proto.Marshal(struc)
		if err != nil {
			return nil, err
		}

		return data, nil

	case "AUTH":
		fmt.Printf("\nВведите для авторизации\n")
		fmt.Printf("Введите логин: ")
		// открываем потоковое чтение из консоли
		reader := bufio.NewReader(os.Stdin)
		// считывваем все данные из консоли
		log, err := reader.ReadString('\n')
		fmt.Printf("Введите пароль: ")
		passw, err := reader.ReadString('\n')
		if err != nil {
			return nil, clientmodels.ErrorInput
		}
		// убираем лишние суффиксы после ввода
		log = strings.TrimSuffix(log, "\n")
		passw = strings.TrimSuffix(passw, "\n")
		// сериализуем данные в proto формат
		struc := &pb.AuthInfo{
			Login:    log,
			Password: passw,
		}
		data, err = proto.Marshal(struc)
		if err != nil {
			return nil, err
		}

		return data, nil

	case "TEXT":
		fmt.Printf("\nВведите текст для сохранения: ")
		// открываем потоковое чтение из консоли
		reader := bufio.NewReader(os.Stdin)
		// читаем текст из консоли
		txt, err := reader.ReadString('\n')
		if err != nil {
			return nil, clientmodels.ErrorInput
		}
		// приводим к типу byte
		data = []byte(txt)

		return data, nil
	case "BINARY":
		fmt.Printf("\nВведите путь к файлу для загрузки: ")
		// открываем потоковое чтение из консоли
		reader := bufio.NewReader(os.Stdin)
		// читаем путь к файлу из консоли
		path, err := reader.ReadString('\n')
		path = strings.TrimSuffix(path, "\n")
		path = strings.TrimSuffix(path, "\r")
		fmt.Printf("%#v\n", path)
		// читаем файл
		data, err := os.ReadFile(path)
		if err != nil {
			return nil, clientmodels.ErrorInput
		}
		// если объем файла слишком большой, сообщаем об этом
		if len(data) > 4194304 {
			return nil, clientmodels.ErrBigData
		}
		return data, nil
	default:
		return nil, clientmodels.ErrorInput
	}
}

// ParceAnswer функция для вывода всех данных на экран
func ParceAnswer(Data []byte, storageID string, Type string, Metainfo string) error {
	// отображаем метаданные
	fmt.Printf("\nInformation: %s\n", Metainfo)

	switch Type {
	case "CARD":
		// десериализуем полученные данные
		out := pb.CardInfo{}
		err := proto.Unmarshal(Data, &out)
		if err != nil {
			return err
		}
		// выводим данные карты на экран
		fmt.Printf("Card numver: %s\n", out.Number)
		fmt.Printf("Card CVV: %s\n", out.CVVcode)
		fmt.Printf("Card exipre: %s\n", out.ExpDate)
		fmt.Printf("Card holder: %s\n", out.Holder)
		return nil

	case "AUTH":
		// десериализуем полученные данные
		out := pb.AuthInfo{}
		err := proto.Unmarshal(Data, &out)
		if err != nil {
			return err
		}
		// выводим данные авторизации на экран
		fmt.Printf("Login: %s\n", out.Login)
		fmt.Printf("Password: %s\n", out.Password)
		return nil
	case "TEXT":
		// отображаем текст на экране
		fmt.Printf("Saved text info: \n%s\n", string(Data))
		return nil
	case "BINARY":
		// просим ввести путь, куда будет сохранен файл
		fmt.Printf("\nВведите путь к файлу для сохранения (включая его имя и расширение): ")
		// открываем потоковое чтение из консоли
		reader := bufio.NewReader(os.Stdin)
		// читаем путь из консоли
		path, err := reader.ReadString('\n')
		path = strings.TrimSuffix(path, "\n")
		path = strings.TrimSuffix(path, "\r")

		// сохрвняем файл по указанному пути
		err = os.WriteFile(path, Data, 0644)
		if err != nil {
			return clientmodels.ErrBigCantSave
		}
		// уведомляем что файл сохранен
		fmt.Printf("Файл сохранен в: %s\n", path)
		return nil
	}
	return nil
}
