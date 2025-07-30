package uuid

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"strings"
)

type UUID [16]byte

func NewV4() UUID {
	bytes := make([]byte, 16)

	_, _ = rand.Read(bytes)

	bytes[6] = (bytes[6] & 0x0f) | 0x40
	bytes[8] = (bytes[8] & 0x3f) | 0x80

	return UUID(bytes)
}

func (uuid UUID) String() string {
	return fmt.Sprintf(
		"%x-%x-%x-%x-%x",
		uuid[0:4],
		uuid[4:6],
		uuid[6:8],
		uuid[8:10],
		uuid[10:16],
	)
}

func Parse(str string) (UUID, error) {
	uuid := UUID{}

	if len(str) != 36 {
		return uuid, errors.New("Invalid UUID format: incorrect length")
	}

	if str[8] != '-' || str[13] != '-' || str[18] != '-' || str[23] != '-' {
		return uuid, errors.New("invalid UUID format: incorrect hyphen positions")
	}

	hexStr := strings.ReplaceAll(str, "-", "")
	if len(hexStr) != 32 {
		return uuid, errors.New("invalid UUID format: incorrect hex string length")
	}

	bytes, err := hex.DecodeString(hexStr)
	if err != nil {
		return uuid, fmt.Errorf("invalid UUID: failed to decode hex: %v", err)
	}

	if len(bytes) != 16 {
		return uuid, errors.New("invalid UUID: decoded length is not 16 bytes")
	}

	if (bytes[6] & 0xf0) != 0x40 {
		return uuid, errors.New("invalid UUID: not version 4")
	}
	if (bytes[8] & 0xc0) != 0x80 {
		return uuid, errors.New("invalid UUID: invalid variant")
	}

	uuid = UUID(bytes)

	return uuid, nil
}
