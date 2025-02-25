package env

import (
	"fmt"
	"os"
)

func GetEnv(key string) string {
	value := os.Getenv(key)
	if value == "" {
		panic(fmt.Sprintf("env %s is null", key))
	}
	return value

}
