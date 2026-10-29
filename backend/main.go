package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	netHttp "net/http"
	"os"
	"strings"

	"backend/database"
	"backend/handlers"
	"backend/repositories"
	"backend/services"

	"github.com/spf13/viper"
)

type Config struct {
	Port        string `mapstructure:"PORT"`
	DatabaseUrl string `mapstructure:"DATABASE_URL"`
}

func corsMiddleware(next netHttp.Handler) netHttp.Handler {
	return netHttp.HandlerFunc(func(w netHttp.ResponseWriter, r *netHttp.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if r.Method == "OPTIONS" {
			w.WriteHeader(netHttp.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func main() {
	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	mux := http.NewServeMux()

	if _, err := os.Stat(".env"); err == nil {
		viper.SetConfigFile(".env")
		_ = viper.ReadInConfig()
	}

	config := Config{
		Port:        viper.GetString("PORT"),
		DatabaseUrl: viper.GetString("DATABASE_URL"),
	}

	db, errDB := database.Connect(config.DatabaseUrl)
	if errDB != nil {
		log.Fatal("Failed to initialize database:", errDB)
	}

	defer db.Close()

	productRepo := repositories.NewProductRepository(db)
	productService := services.NewProductService(productRepo)
	productHandler := handlers.NewProductHandler(productService)

	mux.HandleFunc("/health", func(w netHttp.ResponseWriter, r *netHttp.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{
			"status":  "OK",
			"message": "API Running",
		})
	})

	mux.HandleFunc("/api/products", productHandler.HandleProducts)
	mux.HandleFunc("/api/products/{id}", productHandler.HandleProductById)
	mux.HandleFunc("/api/products/{id}/stock", productHandler.AdjustStock)
	mux.HandleFunc("/api/products/{id}/movements", productHandler.GetMovements)

	handler := corsMiddleware(mux)

	fmt.Println("Server running di localhost:" + config.Port)

	err := netHttp.ListenAndServe(":"+config.Port, handler)
	if err != nil {
		fmt.Println("Failed to run server")
	}
}
