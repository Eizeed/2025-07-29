package dotenv

import (
	"log"
	"os"
	"strings"
)

func DotEnv() {
	bytes, err := os.ReadFile("./.env")
	if err != nil {
		log.Fatalln("Error parsing .env:", err)
	}

	blob := string(bytes)
	lines := strings.Lines(blob)
	for line := range lines {
		if strings.Contains(line, "=") {
			parts := strings.SplitN(line, "=", 2)
			key := strings.TrimSpace(parts[0])
			val := strings.TrimSpace(parts[1])

			os.Setenv(key, val)
		}
	}
}
