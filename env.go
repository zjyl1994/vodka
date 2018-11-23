package vodka

import (
	"os"
)

func mustGetEnv(key string)(value string){
	value = os.Getenv(key)
	if len(value) == 0 {
		Logger.Fatalf("Environment variables \"%s\" is empty!\n",key)
	}
	return value
}