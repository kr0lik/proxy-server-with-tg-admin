package helper

import (
	"encoding/base64"
	"encoding/binary"
	"errors"
	"hash/fnv"
	"time"
)

var ErrInviteTokenExpired = errors.New("token expired")

const InviteTokenTtl = time.Hour

func GenerateInviteToken(username string) string {
	if len(username) > 20 { //nolint: mnd
		username = username[:20]
	}

	h := fnv.New64a()
	_, _ = h.Write([]byte(username))
	hash := h.Sum64()

	now := time.Now().Unix()
	var ts uint32
	switch {
	case now < 0:
		ts = 0
	case now > int64(^uint32(0)):
		ts = ^uint32(0)
	default:
		ts = uint32(now)
	}

	buf := make([]byte, 12) //nolint: mnd
	for i := range 8 {
		buf[i] = byte((hash >> (56 - 8*i)) & 0xFF) //nolint: mnd
	}
	binary.BigEndian.PutUint32(buf[8:], ts)

	xor(buf)

	return base64.RawURLEncoding.EncodeToString(buf)
}

func CheckInviteToken(token string) error {
	buf, err := base64.RawURLEncoding.DecodeString(token)
	if err != nil || len(buf) != 12 {
		return errors.New("invalid token")
	}

	xor(buf)

	ts := binary.BigEndian.Uint32(buf[8:])
	if time.Since(time.Unix(int64(ts), 0)) > InviteTokenTtl {
		return ErrInviteTokenExpired
	}

	return nil
}

func xor(data []byte) {
	xorKey := []byte{0xA5, 0x5A, 0xC3, 0x3C, 0x7E, 0xE7, 0x11, 0x22}

	for i := range data {
		data[i] ^= xorKey[i%len(xorKey)]
	}
}
