package secret

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/binary"
	"fmt"
	"time"

	"github.com/realglebivanov/hstd/hstdlib"
)

func GenerateClientUUIDs(rootSecret []byte) []string {
	var uuids []string

	for i := range hstdlib.XrayClientCount {
		uuids = append(uuids, GenerateClientUUID(i, rootSecret))
	}

	return uuids
}

func GenerateClientUUID(i int, rootSecret []byte) string {
	secret := deriveSubscriptionSecret(i, rootSecret)

	now := time.Now().UTC()
	epoch := now.Add(-3 * time.Hour)
	day := time.Date(epoch.Year(), epoch.Month(), epoch.Day(), 0, 0, 0, 0, time.UTC).Unix()

	h := sha256.New()
	binary.Write(h, binary.BigEndian, day)
	h.Write(secret)
	sum := h.Sum(nil)

	return fmt.Sprintf("%08x-%04x-%04x-%04x-%012x",
		sum[0:4], sum[4:6], sum[6:8], sum[8:10], sum[10:16])
}

func deriveSubscriptionSecret(index int, rootSecret []byte) []byte {
	secretMsg := fmt.Appendf(nil, "secret:%d", index)
	secretMAC := hmac.New(sha256.New, rootSecret)
	secretMAC.Write(secretMsg)

	return secretMAC.Sum(nil)
}
