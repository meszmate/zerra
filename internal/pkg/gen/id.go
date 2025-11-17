package gen

import "crypto/rand"

const charset = "ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

func RID(length int) string {
	b := make([]byte, length)
	_, err := rand.Read(b)
	if err != nil {
		panic(err)
	}

	for i := range length {
		b[i] = charset[int(b[i])%len(charset)]
	}

	return string(b)
}
