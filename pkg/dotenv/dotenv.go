package dotenv

import (
	"os"
	"strings"
)

func DotEnv() error {
	bytes, err := os.ReadFile("./.env")
	if err != nil {
		return err
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

	return nil
}
