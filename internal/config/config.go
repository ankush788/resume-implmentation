package config

import (
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	DBUrl string
	RedisAddr     string
	RedisPassword string
	RedisDB       int
}

func Load() *Config {
	_ = godotenv.Load()

	dbUrl := os.Getenv("DATABASE_URL")
	if dbUrl == "" {
		// Log fatal error if DATABASE_URL is missing. Error will server crash.
		log.Fatal("DATABASE_URL missing")
	}


	redisAddr := os.Getenv("REDIS_ADDR") //This is the location of your Redis server.
	if redisAddr == "" {
		redisAddr = "127.0.0.1:6379"
	}

	redisPassword := os.Getenv("REDIS_PASSWORD") //This is the authentication password for Redis.

	redisDB := 0 //RedisDB = which logical database number inside Redis to use (default 0)
	 //Redis has multiple logical databases inside one server. 
	//Example:
	//DB 0 → cache
	//DB 1 → sessions
	//DB 2 → temporary data

	// if use give which db , he want to use  "REDIS_DB" variable in env (otherwise by default 0)
	if v := os.Getenv("REDIS_DB"); v != "" {
	if parsed, err := strconv.Atoi(v); err == nil {
        redisDB = parsed
    }
	}

return &Config{
    DBUrl:        dbUrl,
    RedisAddr:     redisAddr,
    RedisPassword: redisPassword,
    RedisDB:       redisDB,
}

}
