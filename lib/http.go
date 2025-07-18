package lib

import (
	"bytes"
	"crypto/ed25519"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"petpet/logging"
)

type ResponseType uint8
const (
	PONG_RESPONSE_TYPE ResponseType = iota + 1
	ACKNOWLEDGE_RESPONSE_TYPE
	CHANNEL_MESSAGE_RESPONSE_TYPE
	CHANNEL_MESSAGE_WITH_SOURCE_RESPONSE_TYPE
	DEFERRED_CHANNEL_MESSAGE_WITH_SOURCE_RESPONSE_TYPE
	DEFERRED_UPDATE_MESSAGE_RESPONSE_TYPE
)

type InteractionType uint8
type InteractionTypeExtractor struct {
	Type InteractionType `json:"type"`
}
const (
	PING_INTERACTION_TYPE InteractionType = iota + 1
	APPLICATION_COMMAND_INTERACTION_TYPE
)

var (
	bodyPingResponse           = fmt.Appendf(nil, `{"type":%d}`, PONG_RESPONSE_TYPE)
	bodyUnknownCommandResponse = fmt.Appendf(nil, `{"type":%d,"data":{"content":"Oh uh.. It looks like you tried to use outdated/unknown slash command. Please report this bug to bot owner.","flags":%d}}`, CHANNEL_MESSAGE_WITH_SOURCE_RESPONSE_TYPE, MESSAGE_FLAG_EPHEMERAL)
)

type MessageFlags BitSet
const (
	_ = 1 << iota
	_
	_
	_
	_
	_
	MESSAGE_FLAG_EPHEMERAL
)

func (bot *Bot) DiscordRequestHandler(w http.ResponseWriter, r *http.Request) {
	verified := verifyDiscordRequest(r, ed25519.PublicKey(bot.PublicKey))
	if !verified {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		logging.Error("Unauthorized request to Discord endpoint", "remote_addr", r.RemoteAddr)
		return
	}

	buf := bot.jsonBufferPool.Get().(*[]byte)
	defer bot.jsonBufferPool.Put(buf)

	n, err := r.Body.Read(*buf)
	if err != nil && err != io.EOF {
		http.Error(w, "bad request - failed to read body payload", http.StatusBadRequest)
		logging.Error("Failed to read request body", "error", err)
		return
	}
	defer r.Body.Close()

	var extractor InteractionTypeExtractor
	if err := json.Unmarshal((*buf)[:n], &extractor); err != nil {
		http.Error(w, "bad request - invalid body json payload", http.StatusBadRequest)
		logging.Error("Failed to unmarshal request body", "error", err)
		return
	}

	switch extractor.Type {
	case PING_INTERACTION_TYPE:
		w.Header().Add("Content-Type", CONTENT_TYPE_JSON)
		w.Write(bodyPingResponse)
		return
	case APPLICATION_COMMAND_INTERACTION_TYPE:
		var interaction CommandInteraction
		if err := json.Unmarshal((*buf)[:n], &interaction); err != nil {
			http.Error(w, "bad request - failed to decode CommandInteraction", http.StatusBadRequest)
			logging.Error("Failed to unmarshal CommandInteraction", "error", err)
			return
		}
		bot.commandInteractionHandler(w, interaction)
		return
	}
}


func (bot *Bot) makeHttpRequestToDiscord(method string, url string, body any, files []DiscordFile, authRequired bool) error {
	var req *http.Request
	var err error

	if len(files) == 0 {
		// No files, send as JSON
		req, err = http.NewRequest(
			method,
			DISCORD_API_URL+url,
			marshalToReadCloser(body),
		)
		if err != nil {
			return err
		}
		req.Header.Set("Content-Type", "application/json")
	} else {
		// With files, use multipart/form-data
		var b bytes.Buffer
		w := multipart.NewWriter(&b)

		// Add payload_json part
		payload, err := json.Marshal(body)
		if err != nil {
			return err
		}
		if err := w.WriteField("payload_json", string(payload)); err != nil {
			return err
		}

		// Add files
		for i, file := range files {
			part, err := w.CreateFormFile(fmt.Sprintf("files[%d]", i), file.Filename)
			if err != nil {
				return err
			}
			if _, err := io.Copy(part, file.Reader); err != nil {
				return err
			}
		}

		if err := w.Close(); err != nil {
			return err
		}

		req, err = http.NewRequest(
			method,
			DISCORD_API_URL+url,
			io.NopCloser(&b),
		)
		if err != nil {
			return err
		}
		req.Header.Set("Content-Type", w.FormDataContentType())
	}

	if authRequired {
		req.Header.Set("Authorization", "Bot "+bot.Token)
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		logging.Error("HTTP request failed", "error", err)
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		errorBody := &bytes.Buffer{}
		if _, err := io.Copy(errorBody, resp.Body); err != nil {
			logging.Error("Failed to read error response body", "error", err)
			return fmt.Errorf("request failed with status %s", resp.Status)
		}
		logging.Error("Discord API request failed", "status", resp.Status, "body", errorBody.String())
	}
	return nil
}

// DiscordFile represents a file to upload to Discord.
type DiscordFile struct {
	Filename string
	Reader   io.Reader
}

func marshalToReadCloser(v any) *readCloserWrapper {
	data, err := json.Marshal(v)
	if err != nil {
		return &readCloserWrapper{Reader: nil}
	}
	return &readCloserWrapper{Reader: bytes.NewReader(data)}
}

type readCloserWrapper struct {
	Reader *bytes.Reader
}

func (rc *readCloserWrapper) Read(p []byte) (n int, err error) {
	if rc.Reader == nil {
		return 0, io.EOF
	}
	return rc.Reader.Read(p)
}

func (rc *readCloserWrapper) Close() error {
	return nil
}