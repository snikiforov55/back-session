package main

import (
	"github.com/go-redis/redis"
	"github.com/gorilla/mux"
	"github.com/snikiforov55/back-session/session"
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
		log.Println("An error occured starting HTTP listener at port " + port)
		log.Println("Error: " + err.Error())
	}
}

type Config struct {
	RedisHost     string `json:"redis_host"`
	RedisPort     string `json:"redis_port"`
	RedisPassword string `json:"redis_password"`
	RedisDb       int    `json:"redis_db"`
	SessionExpSec int    `json:"session_exp_sec"`
	ServicePort   string `json:service_port`
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
	return &Config{RedisHost: getEnvOrString("REDIS_HOSTNAME", "localhost"),
		RedisPort:     getEnvOrString("REDIS_PORT", "6379"),
		RedisPassword: getEnvOrString("REDIS_PASSWORD", ""),
		RedisDb:       getEnvOrInt("REDIS_DB", 0),
		SessionExpSec: getEnvOrInt("SESSION_EXP_SEC", session.DefaultSessionExpirationSec),
		ServicePort:   getEnvOrString("SERVICE_PORT", "8090"),
	}
}

func main() {

	config := NewConfig()

	client := redis.NewClient(&redis.Options{
		Addr:     config.RedisHost + ":" + config.RedisPort,
		Password: config.RedisPassword, // no password set
		DB:       config.RedisDb,       // use default DB
	})

	server := session.NewServer(client, config.SessionExpSec)

	StartWebServer(config.ServicePort, server.Router)
}
