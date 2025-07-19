package lib

import (
	"bytes"
	"crypto/ed25519"
	"encoding/hex"
	"fmt"
	"io"
	"net/http"
)

func StringToSnowflake(s string) (Snowflake, error) {
	var id Snowflake
	
	if _, err := fmt.Sscanf(s, "%d", &id); err != nil {
		return 0, fmt.Errorf("failed to parse string to uint64: %w", err)
	}

	return id, nil
}

func verifyDiscordRequest(r *http.Request, key ed25519.PublicKey) bool {
	var msg bytes.Buffer

	signature := r.Header.Get("X-Signature-Ed25519")
	if signature == "" {
		return false
	}

	sig, err := hex.DecodeString(signature)
	if err != nil {
		return false
	}

	if len(sig) != ed25519.SignatureSize || sig[63]&224 != 0 {
		return false
	}

	timestamp := r.Header.Get("X-Signature-Timestamp")
	if timestamp == "" {
		return false
	}

	msg.WriteString(timestamp)

	defer CloseBody(r.Body)
	var body bytes.Buffer


	defer func() {
		r.Body = io.NopCloser(&body)
	}()

	if _, err = io.Copy(&msg, io.TeeReader(r.Body, &body)); err != nil {
		return false
	}

	return ed25519.Verify(key, msg.Bytes(), sig)
}
