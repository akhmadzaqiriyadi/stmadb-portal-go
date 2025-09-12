package main

import (
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	
	// Import dari proyek Anda
	"github.com/akhmadzaqiriyadi/stmadb-portal-go/internal/database"
	"github.com/akhmadzaqiriyadi/stmadb-portal-go/internal/router"
	_ "github.com/akhmadzaqiriyadi/stmadb-portal-go/docs" // Import kosong untuk Swagger docs
)

// init() berjalan sekali sebelum fungsi main()
func init() {
	// Setup Viper untuk membaca konfigurasi dari file .env
	viper.SetConfigFile(".env")
	if err := viper.ReadInConfig(); err != nil {
		logrus.Fatalf("Error reading config file: %v", err)
	}

	// Setup Logrus untuk logging terstruktur
	logrus.SetFormatter(&logrus.JSONFormatter{})
	logrus.SetLevel(logrus.InfoLevel)
}

// @title           STMADB Portal Backend API
// @version         1.0
// @description     This is the API documentation for STMADB Portal Backend.
// @host            localhost:3000
// @BasePath        /api/v1
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
func main() {
	// Ambil port dari file .env, gunakan "3000" sebagai default
	port := viper.GetString("PORT")
	if port == "" {
		port = "3000"
	}

	// Inisialisasi koneksi database menggunakan Prisma Client
	logrus.Info("Connecting to database...")
	dbClient := database.NewClient()
	defer func() {
		if err := dbClient.Disconnect(); err != nil {
			logrus.Errorf("Failed to disconnect from database: %v", err)
		}
	}()
	logrus.Info("🗄️ Database connected successfully")

	// Setup router yang berisi semua endpoint API
	r := router.SetupRouter(dbClient)

	// Mulai server
	logrus.Infof("🚀 Server starting on port %s", port)
	logrus.Infof("📚 API Documentation available at http://localhost:%s/swagger/index.html", port)
	if err := r.Run(":" + port); err != nil {
		logrus.Fatalf("Failed to run server: %v", err)
	}
}