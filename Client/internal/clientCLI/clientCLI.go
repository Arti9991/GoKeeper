package clientcli

import (
	"bufio"
	"fmt"
	"os"

	"github.com/Arti9991/GoKeeper/client/internal/requseter"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	buildVersion string = "1.0.0"
	buildDate    string = "31.05.2025"
)

var ServAddr string

var rootCmd = &cobra.Command{
	Use:   "Keeper client",
	Short: "Client for GoKeeper service",
	Long: `Этот клиент предоставляет возможность сохранять бинарные данные как в локальном хранилище,
	так и в защищенном хранлище на удаленном сервере.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Welcome to Keeper client! Use --help for usage.")
	},
}

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version number of Keeper",
	Long:  `Вывод версии и даты сборки клиента`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("Client build version: %s\n", buildVersion)
		fmt.Printf("Client build date: %s\n", buildDate)
	},
}

var userLogin = &cobra.Command{
	Use:   "login",
	Short: "Login user. Where 1st your login and 2nd your password",
	Long:  `Авторизация пользователя. Параметры: [1]Логин [2]Пароль`,
	Run: func(cmd *cobra.Command, args []string) {
		req, err := requseter.NewRequester(viper.GetString("addr"))
		if err != nil {
			fmt.Printf("\nError in starting requester: %s\n", err.Error())
			return
		}
		// просим ввести данные пользователя
		fmt.Printf("\nВведите данные для авторизации\n")
		fmt.Printf("Имя пользователя: ")
		// открываем потоковое чтение из консоли
		reader := bufio.NewReader(os.Stdin)
		// считывваем все данные из консоли
		Login, err := reader.ReadString('\n')

		fmt.Printf("Введите пароль: ")
		Password, err := reader.ReadString('\n')
		err = req.LoginRequest(Login, Password)
		if err != nil {
			fmt.Printf("\nError in LoginRequest: %s\n", err.Error())
			return
		}
		return
	},
}

var userRegister = &cobra.Command{
	Use:   "register",
	Short: "Register. Where 1st your login and 2nd your password",
	Long:  `Регистрация новго пользователя. Параметры: [1]Логин [2]Пароль`,
	Run: func(cmd *cobra.Command, args []string) {
		req, err := requseter.NewRequester(viper.GetString("addr"))
		if err != nil {
			fmt.Printf("\nError in starting requester: %s\n", err.Error())
			return
		}
		// просим ввести данные пользователя
		fmt.Printf("\nВведите данные для авторизации\n")
		fmt.Printf("Имя пользователя: ")
		// открываем потоковое чтение из консоли
		reader := bufio.NewReader(os.Stdin)
		// считывваем все данные из консоли
		Login, err := reader.ReadString('\n')

		fmt.Printf("Введите пароль: ")
		Password, err := reader.ReadString('\n')
		err = req.RegisterRequest(Login, Password)
		if err != nil {
			fmt.Printf("\nError in RegisterRequest: %s\n", err.Error())
			return
		}
		return
	},
}

var userLogout = &cobra.Command{
	Use:   "logout",
	Short: "Logout user",
	Long: `Выход пользователя из системы. Все локальные данные будут удалены! 
	Во избежание потери данных, нужно их предварительно синхронизировать с сервером!`,
	Run: func(cmd *cobra.Command, args []string) {
		req, err := requseter.NewRequester(viper.GetString("addr"))
		if err != nil {
			fmt.Printf("\nError in starting requester: %s\n", err.Error())
			return
		}
		err = req.LogoutRequest()
		if err != nil {
			fmt.Printf("\nError in LogoutRequest: %s\n", err.Error())
			return
		}
		return
	},
}

var saveData = &cobra.Command{
	Use:   "save",
	Short: "Save users data.  Where 1st is data type (AUTH,CARD,TEXT,BINARY)",
	Long: `	Сохранение пользовательских данных. 
	В единственном аргументе передается тип этих данных [1]Тип:
	AUTH - данные для авторизации(Логин и пароль),
	CARD - данные карты (номер, срок действия, CVV, держатель),
	TEXT - текстовая информация,
	BINARY - бинарные данные. Размером не более 4 Мб.
	Иные типы и параметры не поддерживаются!`,
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		req, err := requseter.NewRequester(viper.GetString("addr"))
		if err != nil {
			fmt.Printf("\nError in starting requester: %s\n", err.Error())
			return
		}
		err = req.SaveDataRequest(args[0], viper.GetBool("offlineMode"))
		if err != nil {
			fmt.Printf("\nError in SaveDataRequest: %s\n", err.Error())
			return
		}
		return
	},
}

var getData = &cobra.Command{
	Use:   "get",
	Short: "Get users data. Where 1st is data ID",
	Long: `Получение сохраненных данных. 
	В первом параметре [1] передается ID данных.`,
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		req, err := requseter.NewRequester(viper.GetString("addr"))
		if err != nil {
			fmt.Printf("\nError in starting requester: %s\n", err.Error())
			return
		}
		err = req.GetDataRequest(args[0], viper.GetBool("offlineMode"))
		if err != nil {
			fmt.Printf("\nError in GetDataRequest: %s\n", err.Error())
			return
		}
		return
	},
}

var showTableLoc = &cobra.Command{
	Use:   "showLoc",
	Short: "Show local users data.",
	Long:  `Отображает данные пользователя, сохраненные на этом устройстве.`,
	Run: func(cmd *cobra.Command, args []string) {
		req, err := requseter.NewRequester(viper.GetString("addr"))
		if err != nil {
			fmt.Printf("\nError in starting requester: %s\n", err.Error())
			return
		}
		err = req.ShowDataLoc()
		if err != nil {
			fmt.Printf("\nError in ShowDataLoc: %s\n", err.Error())
			return
		}
		return
	},
}

var showTableOn = &cobra.Command{
	Use:   "showOn",
	Short: "Show online users data.",
	Long:  `Отображает данные пользователя, сохраненные на сервере.`,
	Run: func(cmd *cobra.Command, args []string) {
		req, err := requseter.NewRequester(viper.GetString("addr"))
		if err != nil {
			fmt.Printf("\nError in starting requester: %s\n", err.Error())
			return
		}
		err = req.ShowDataOn(viper.GetBool("offlineMode"))
		if err != nil {
			fmt.Printf("\nError in ShowDataOn: %s\n", err.Error())
			return
		}
		return
	},
}

var syncData = &cobra.Command{
	Use:   "sync",
	Short: "Sync users data.",
	Long:  `Отправляет новые или обновленные локальные данные на сервер`,
	Run: func(cmd *cobra.Command, args []string) {
		req, err := requseter.NewRequester(viper.GetString("addr"))
		if err != nil {
			fmt.Printf("\nError in starting requester: %s\n", err.Error())
			return
		}
		err = req.SyncRequest(viper.GetBool("offlineMode"))
		if err != nil {
			fmt.Printf("\nError in SyncRequest: %s\n", err.Error())
			return
		}
		return
	},
}

var deleteData = &cobra.Command{
	Use:   "delete",
	Short: "Delete users data. Where 1st is data ID",
	Long:  `Удаление данных. Первый параметр [1], это ID данных.`,
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		req, err := requseter.NewRequester(viper.GetString("addr"))
		if err != nil {
			fmt.Printf("\nError in starting requester: %s\n", err.Error())
			return
		}
		err = req.DeleteDataRequest(args[0], viper.GetBool("offlineMode"))
		if err != nil {
			fmt.Printf("\nError in DeleteDataRequest: %s\n", err.Error())
			return
		}
		return
	},
}
var updateData = &cobra.Command{
	Use:   "update",
	Short: "Update users data. Where 1st is data ID",
	Long:  `Обновление данных. Первый параметр [1], это ID данных.`,
	Args:  cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		req, err := requseter.NewRequester(viper.GetString("addr"))
		if err != nil {
			fmt.Printf("\nError in starting requester: %s\n", err.Error())
			return
		}
		err = req.UpdateDataRequest(args[0], args[1], viper.GetBool("offlineMode"))
		if err != nil {
			fmt.Printf("\nError in UpdateDataRequest: %s\n", err.Error())
			return
		}
		return
	},
}

func StartCLI() {

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().StringVarP(&ServAddr, "addr", "a", ":8082", "Server address")
	_ = viper.BindPFlag("addr", rootCmd.PersistentFlags().Lookup("addr"))
	// Добавляем глобальный флаг к root-команде
	rootCmd.PersistentFlags().BoolP("offlineMode", "o", false, "Работа в офлайн режиме")
	_ = viper.BindPFlag("offlineMode", rootCmd.PersistentFlags().Lookup("offlineMode"))

	rootCmd.AddCommand(userLogin)
	rootCmd.AddCommand(userRegister)
	rootCmd.AddCommand(userLogout)
	rootCmd.AddCommand(saveData)
	rootCmd.AddCommand(getData)
	rootCmd.AddCommand(syncData)
	rootCmd.AddCommand(showTableLoc)
	rootCmd.AddCommand(showTableOn)
	rootCmd.AddCommand(versionCmd)
	rootCmd.AddCommand(deleteData)
	rootCmd.AddCommand(updateData)

	// rootCmd.AddCommand(listCmd)
	// rootCmd.AddCommand(deleteCmd)
}
