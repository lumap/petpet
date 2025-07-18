package utils

import (
	"net/http"
	"petpet/lib"
	"strconv"
	"strings"
)

func IsLinkAnImageURL(url string) (bool, error) {
	resp, err := http.Head(url)
    if err != nil {
        return false, err
    }
    defer resp.Body.Close()

    contentType := resp.Header.Get("Content-Type")

    return strings.HasPrefix(contentType, "image/"), nil
}

func MakeAvatarURL(userId lib.Snowflake, hash string) string {
	if hash == "" {
		return lib.DISCORD_CDN_URL + "/embed/avatars/" + strconv.FormatUint(uint64(userId>>22)%6, 10) + ".png"
	}

	return lib.DISCORD_CDN_URL + "/avatars/" + strconv.FormatUint(uint64(userId), 10) + "/" + hash + ".png"
}