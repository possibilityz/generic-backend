package main

import (
	"log"
	"net/http"

	"example.com/internal/configs"
	"example.com/internal/logger"
	"example.com/internal/services"
	"github.com/spf13/viper"
)

func main() {
	// Initialize Viper across the application
	configs.InitializeViper()

	// Initialize Logger across the application
	logger.InitializeZapCustomLogger()

	// Initialize Oauth2 Services
	services.InitializeOAuthGoogle()
	services.InitializeOAuthGithub()

	// Routes for the application
	http.HandleFunc("/", services.HandleMain)
	http.HandleFunc("/login-gl", services.HandleGoogleLogin)
	http.HandleFunc("/callback-gl", services.CallBackFromGoogle)
	http.HandleFunc("/login-gh", services.HandleGithubLogin)
	http.HandleFunc("/callback-gh", services.CallBackFromGithub)

	logger.Log.Info("Started running on http://localhost:" + viper.GetString("port"))
	log.Fatal(http.ListenAndServe(":"+viper.GetString("port"), nil))
}
