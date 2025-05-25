package clientcli

import (
	"fmt"
	"os"

	"github.com/Arti9991/GoKeeper/client/internal/requseter"
	"github.com/spf13/cobra"
)

// func StartCLI(message string) {
// 	rootCmd := &cobra.Command{
// 		Use:   "mycli",
// 		Short: "My CLI is a simple CLI application built with Cobra and Viper",
// 	}
// 	rootCmd.PersistentFlags().StringVarP(&message, "message", "m", "", "A custom message")
// 	viper.BindPFlag("message", rootCmd.PersistentFlags().Lookup("message"))
// 	viper.SetDefault("message", "Welcome to my CLI configured with Viper!")

// 	rootCmd.Run = func(cmd *cobra.Command, args []string) {
// 		message := viper.GetString("message")
// 		fmt.Println(message)
// 	}

// 	versionCmd := &cobra.Command{
// 		Use:   "version",
// 		Short: "Print the version number of my cli",
// 		Run: func(cmd *cobra.Command, args []string) {
// 			fmt.Println("mycli v0.1")
// 		},
// 	}

// 	sayHelloCmd := &cobra.Command{
// 		Use:   "sayhello",
// 		Short: "Say Hello",
// 		Run: func(cmd *cobra.Command, args []string) {
// 			fmt.Println("Hello!")
// 		},
// 	}

// 	rootCmd.AddCommand(versionCmd, sayHelloCmd)

// 	if err := rootCmd.Execute(); err != nil {
// 		fmt.Println(err)
// 		os.Exit(1)
// 	}
// }

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

var userLogin = &cobra.Command{
	Use:   "login",
	Short: "Login user",
	Long:  `Login user. Where 1st your login and 2nd your password`,
	Args:  cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		requseter.TestLoginRequest(args[0], args[1])
	},
}

var userRegister = &cobra.Command{
	Use:   "register",
	Short: "Register user",
	Long:  `Register user. Where 1st your login and 2nd your password`,
	Args:  cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		requseter.TestRegisterRequest(args[0], args[1])
	},
}

var saveData = &cobra.Command{
	Use:   "save",
	Short: "Save users data",
	Long:  `Save users data. Where 1st is data type (AUTH,CARD,TEXT,BINARY)`,
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		requseter.TestSaveDataRequest(args[0])
	},
}

var getData = &cobra.Command{
	Use:   "Get",
	Short: "Get users data",
	Long:  `Get users data. Where 1st is data ID`,
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		requseter.TestGetDataRequest(args[0])
	},
}

func StartCLI(req *requseter.ReqStruct) {

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.AddCommand(userLogin)
	rootCmd.AddCommand(userRegister)
	rootCmd.AddCommand(saveData)
	rootCmd.AddCommand(getData)
	// rootCmd.AddCommand(listCmd)
	// rootCmd.AddCommand(deleteCmd)
}
