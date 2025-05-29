package clientcli

import (
	"fmt"
	"os"

	"github.com/Arti9991/GoKeeper/client/internal/requseter"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	buildVersion string = "1.0.0"
	buildDate    string = "29.05.2025"
)

var ServAddr string

var rootCmd = &cobra.Command{
	Use:   "Keeper client",
	Short: "Client for GoKeeper service",
	Long: `This client implements functionality to save data 
		to GoKeeper server. Also this client provides possibility to save data in
		offline mode and syncronise it with server in normal mode`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Welcome to Keeper client! Use --help for usage.")
	},
}

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version number of Keeper",
	Long:  `Print the version number of Keeper`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("Client build version: %s\n", buildVersion)
		fmt.Printf("Client build date: %s\n", buildDate)
	},
}

var userLogin = &cobra.Command{
	Use:   "login",
	Short: "Login user",
	Long:  `Login user. Where 1st your login and 2nd your password`,
	Args:  cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		req := requseter.NewRequester(viper.GetString("addr"))
		requseter.LoginRequest(args[0], args[1], req)
	},
}

var userRegister = &cobra.Command{
	Use:   "register",
	Short: "Register user",
	Long:  `Register user. Where 1st your login and 2nd your password`,
	Args:  cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		req := requseter.NewRequester(viper.GetString("addr"))
		requseter.RegisterRequest(args[0], args[1], req)
	},
}

var userLogout = &cobra.Command{
	Use:   "logout",
	Short: "logout user",
	Long:  `logout user`,
	Run: func(cmd *cobra.Command, args []string) {
		req := requseter.NewRequester(viper.GetString("addr"))
		requseter.LogoutRequest(req)
	},
}

var saveData = &cobra.Command{
	Use:   "save",
	Short: "Save users data",
	Long:  `Save users data. Where 1st is data type (AUTH,CARD,TEXT,BINARY)`,
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		req := requseter.NewRequester(viper.GetString("addr"))
		requseter.SaveDataRequest(args[0], req, viper.GetBool("offlineMode"))
	},
}

var getData = &cobra.Command{
	Use:   "get",
	Short: "Get users data",
	Long:  `Get users data. Where 1st is data ID`,
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		req := requseter.NewRequester(viper.GetString("addr"))
		requseter.GetDataRequest(args[0], req, viper.GetBool("offlineMode"))
	},
}

var showTableLoc = &cobra.Command{
	Use:   "showLoc",
	Short: "show local users data",
	Long:  `show users data`,
	Run: func(cmd *cobra.Command, args []string) {
		req := requseter.NewRequester(viper.GetString("addr"))
		requseter.ShowDataLoc(req)
	},
}

var showTableOn = &cobra.Command{
	Use:   "showOn",
	Short: "show online users data",
	Long:  `show users data`,
	Run: func(cmd *cobra.Command, args []string) {
		req := requseter.NewRequester(viper.GetString("addr"))
		requseter.ShowDataOn(req, viper.GetBool("offlineMode"))
	},
}

var syncData = &cobra.Command{
	Use:   "sync",
	Short: "Sync users data",
	Long:  `Sync users data`,
	Run: func(cmd *cobra.Command, args []string) {
		req := requseter.NewRequester(viper.GetString("addr"))
		requseter.SyncRequest(req, viper.GetBool("offlineMode"))
	},
}

var deleteData = &cobra.Command{
	Use:   "delete",
	Short: "delete users data",
	Long:  `delete users data`,
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		req := requseter.NewRequester(viper.GetString("addr"))
		requseter.DeleteDataRequest(args[0], req, viper.GetBool("offlineMode"))
	},
}
var updateData = &cobra.Command{
	Use:   "update",
	Short: "update users data",
	Long:  `update users data`,
	Args:  cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		req := requseter.NewRequester(viper.GetString("addr"))
		requseter.UpdateDataRequest(args[0], args[1], req, viper.GetBool("offlineMode"))
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
	rootCmd.PersistentFlags().Bool("offlineMode", false, "Работа в офлайн режиме")
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
