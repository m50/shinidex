package passwords

import (
	"crypto/rand"
	"crypto/subtle"
	"encoding/base64"
	"errors"
	"fmt"
	"strings"

	"golang.org/x/crypto/argon2"
)

var (
	ErrInvalidHash         = errors.New("the encoded hash is not in the correct format")
	ErrIncompatibleVersion = errors.New("incompatible version of argon2")
	ErrNotAccepted         = errors.New("passwords do not match")
)

type params struct {
	memory      uint32
	iterations  uint32
	parallelism uint8
	keyLength   uint32
}

func ComparePasswords(encodedPasswd, plainTextPasswd string) error {
	p, salt, hash, err := decodeHash(encodedPasswd)
	if err != nil {
		return err
	}

	otherHash := hashPassword(plainTextPasswd, salt, p)

	if subtle.ConstantTimeCompare(hash, otherHash) == 1 {
		return nil
	}

	return ErrNotAccepted
}

func hashPassword(passwd string, salt []byte, p *params) []byte {
	return argon2.IDKey([]byte(passwd), salt, p.iterations, p.memory, p.parallelism, p.keyLength)
}

func HashPassword(passwd string) (string, error) {
	salt := make([]byte, 16)
	_, err := rand.Read(salt)
	if err != nil {
		return "", err
	}
	hash := hashPassword(passwd, salt, &params{
		memory:      64 * 1024,
		iterations:  3,
		parallelism: 2,
		keyLength:   32,
	})
	b64Salt := base64.RawStdEncoding.EncodeToString(salt)
	b64Hash := base64.RawStdEncoding.EncodeToString(hash)

	encodedHash := fmt.Sprintf("$argon2id$v=%d$m=%d,t=%d,p=%d$%s$%s", argon2.Version, 64*1024, 3, 2, b64Salt, b64Hash)
	return encodedHash, nil
}

func decodeHash(encodedHash string) (*params, []byte, []byte, error) {
	vals := strings.Split(encodedHash, "$")
	if len(vals) != 6 {
		return nil, nil, nil, ErrInvalidHash
	}

	var version int
	_, err := fmt.Sscanf(vals[2], "v=%d", &version)
	if err != nil {
		return nil, nil, nil, err
	}
	if version != argon2.Version {
		return nil, nil, nil, ErrIncompatibleVersion
	}

	p := &params{}
	_, err = fmt.Sscanf(vals[3], "m=%d,t=%d,p=%d", &p.memory, &p.iterations, &p.parallelism)
	if err != nil {
		return nil, nil, nil, err
	}

	salt, err := base64.RawStdEncoding.Strict().DecodeString(vals[4])
	if err != nil {
		return nil, nil, nil, err
	}

	hash, err := base64.RawStdEncoding.Strict().DecodeString(vals[5])
	if err != nil {
		return nil, nil, nil, err
	}

	p.keyLength = uint32(len(hash))
	return p, salt, hash, nil
}
