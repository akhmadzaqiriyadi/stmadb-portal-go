package main

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"

	// Path modul dari go.mod
	"github.com/akhmadzaqiriyadi/stmadb-portal-go/internal/database" 
)

// init() akan berjalan sebelum fungsi main()
func init() {
	// Setup Viper untuk membaca file .env
	viper.SetConfigFile(".env")
	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("Error reading config file, %s", err)
	}

	// Setup Logrus
	logrus.SetFormatter(&logrus.JSONFormatter{})
	logrus.SetLevel(logrus.InfoLevel)
}

func main() {
	// Ambil port dari konfigurasi, dengan default 8080
	port := viper.GetString("PORT")
	if port == "" {
		port = "8080"
	}

	// Inisialisasi koneksi database
	logrus.Info("Connecting to database...")
	dbClient := database.NewClient()
	defer func() {
		if err := dbClient.Disconnect(); err != nil {
			logrus.Fatalf("Failed to disconnect from database: %v", err)
		}
	}()
	logrus.Info("🗄️ Database connected successfully")


	// Inisialisasi Gin router
	router := gin.Default()

	// Definisikan route health check
	router.GET("/api/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status":  "OK",
			"message": "STMADB Portal Go Backend is running!",
		})
	})

	// Jalankan server
	logrus.Infof("🚀 Server starting on port %s", port)
	if err := router.Run(":" + port); err != nil {
		logrus.Fatalf("Failed to run server: %v", err)
	}
}