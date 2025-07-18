package lib

import (
	"crypto/ed25519"
	"encoding/base64"
	"encoding/hex"
	"errors"
	"strings"
	"sync"
)

type Bot struct {
	ApplicationID    Snowflake
	PublicKey        ed25519.PublicKey
	Token            string

	jsonBufferPool   *sync.Pool
	commands         *SharedMap[string, Command]
}

func CreateBot(botToken string, publicKey string) Bot {
	decodedKey, err := hex.DecodeString(publicKey)
	if err != nil {
		panic("failed to decode public key: " + err.Error())
	}

	id, err := extractUserIDFromToken(botToken)
	if err != nil {
		panic("failed to extract bot user ID from bot token: " + err.Error())
	}

	return Bot{
		ApplicationID: id,
		Token:         botToken,
		PublicKey:     decodedKey,
		jsonBufferPool: &sync.Pool{
			New: func() any {
				buf := make([]byte, 8192)
				return &buf
			},
		},
		commands: NewSharedMap[string, Command](),
	}
}

func extractUserIDFromToken(token string) (Snowflake, error) {
	strs := strings.Split(token, ".")
	if len(strs) == 0 {
		return 0, errors.New("token is not in a valid format")
	}

	hexID := strings.Replace(strs[0], "Bot ", "", 1)

	byteID, err := base64.RawStdEncoding.DecodeString(hexID)
	if err != nil {
		return 0, err
	}

	return StringToSnowflake(string(byteID))
}