package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	// Library Swagger

	"kasir-api/database"
	_ "kasir-api/docs"
	"kasir-api/routes"

	"github.com/spf13/viper"
)

type Config struct {
	Port   string `mapstructure:"PORT"`
	DBConn string `mapstructure:"DB_CONN"`
}

// @title           CodeWithUmam - Task Session 1
// @version         1.0
// @description     Task Untuk Session 1.
// @host            localhost:8080
// @BasePath        /
func main() {
	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	if _, err := os.Stat(".env"); err == nil {
		viper.SetConfigFile(".env")
		_ = viper.ReadInConfig()
	}

	config := Config{
		Port:   viper.GetString("PORT"),
		DBConn: viper.GetString("DB_CONN"),
	}

	db, err := database.InitDB(config.DBConn)
	if err != nil {
		log.Fatal("Failed to Initialize Database: ", err)
	}
	defer db.Close()

	routes.RegisterAllRoutes(db)

	fmt.Println("server running di localhost: " + config.Port)
	err = http.ListenAndServe(":"+config.Port, nil)
	if err != nil {
		fmt.Println("gagal running server")
	}
}
