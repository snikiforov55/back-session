package main

import (
	"github.com/gorilla/mux"
	"github.com/snikiforov55/back-session/session"
	"github.com/snikiforov55/back-session/session/db"
	"log"
	"net/http"
	"os"
	"strconv"
)

func StartWebServer(port string, router *mux.Router) {
	log.Println("Starting HTTP service at " + port)
	http.Handle("/", router)
	err := http.ListenAndServe(":"+port, nil) // Goroutine will block here
	if err != nil {
		log.Println("An error occurred starting HTTP listener at port " + port)
		log.Println("Error: " + err.Error())
	}
}

type Config struct {
	SessionExpSec int    `json:"session_exp_sec"`
	ServicePort   string `json:"service_port"`
}

func getEnvOrString(key string, defaultValue string) string {
	if env, set := os.LookupEnv(key); set {
		return env
	} else {
		return defaultValue
	}
}

func getEnvOrInt(key string, defaultValue int) int {
	if env, set := os.LookupEnv(key); set {
		val, err := strconv.Atoi(env)
		if err != nil {
			return defaultValue
		} else {
			return val
		}
	} else {
		return defaultValue
	}
}
func NewConfig() *Config {
	return &Config{
		SessionExpSec: getEnvOrInt("SESSION_EXP_SEC", session.DefaultSessionExpirationSec),
		ServicePort:   getEnvOrString("SERVICE_PORT", "8090"),
	}
}
func NewRedisConfig() *db.RedisConfig {
	return &db.RedisConfig{RedisHost: getEnvOrString("REDIS_HOSTNAME", "localhost"),
		RedisPort:     getEnvOrString("REDIS_PORT", "6379"),
		RedisPassword: getEnvOrString("REDIS_PASSWORD", ""),
		RedisDb:       getEnvOrInt("REDIS_DB", 0),
	}
}
func main() {

	config := NewConfig()
	redisConfig := NewRedisConfig()

	redisClient := db.NewRedisClient(redisConfig, db.RandomString)

	server, err := session.NewServer(redisClient, config.SessionExpSec)
	if err != nil {
		log.Panicln("Failed to create a Session object. Error: " + err.Error())
		return
	}

	StartWebServer(config.ServicePort, server.Router)
}
