package main

import (
	"fmt"
	"log"
	"time"

	"github.com/Arti9991/GoKeeper/server/internal/server"
)

var (
	buildVersion string = "N/A"
	buildDate    string = "N/A"
	buildCommit  string = "N/A"
)

func main() {

	fmt.Printf("Server build version: %s\n", buildVersion)
	fmt.Printf("Server build date: %s\n", buildDate)
	fmt.Printf("Server build commit: %s\n", buildCommit)

	time.Sleep(1 * time.Second)

	err := server.RunServer()
	if err != nil {
		log.Fatal(err)
	}
}

/*
docker run -p 8082:8082 keeper:v00

protoc --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative proto/server.proto

openssl genrsa -out server.key 2048

openssl req -new -x509 -sha256   -key server.key   -out server.crt   -days 3650   -config openssl-san.cnf
*/

/* TODO:
1. Авторизация пользователя на клиенте не через аргументы, а через консольный ввод - DONE
2. Генерация ЛИЧНОГО ключа для каждого пользователя к GSM для бинарного хранилища
3. Добавить метод download на скачивание данных с сервера в локальное бинарное хранилище
4. Подтверждение на обновление локальных данных при синхронизации
5. Добавить пользовательский интерфейс.
*/
