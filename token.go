package crypto_exchange

import (
	"crypto/sha512"
	"encoding/hex"

	"github.com/gofrs/uuid"
	"github.com/golang-jwt/jwt/v4"
)

func generateAuthorizationToken(accessKey, secretKey string, query *string) (string, error) {
	claimMap := jwt.MapClaims{}
	claimMap["access_key"] = accessKey

	nonce, err := uuid.NewV4()
	if err != nil {
		return "", err
	}
	claimMap["nonce"] = nonce
	if query != nil {
		hash := sha512.New()
		hash.Write([]byte(*query))
		hashedQuery := hex.EncodeToString(hash.Sum(nil))

		claimMap["query_hash"] = hashedQuery
		claimMap["query_hash_alg"] = "SHA512"
	}

	claim := jwt.NewWithClaims(jwt.SigningMethodHS256, claimMap)
	token, err := claim.SignedString([]byte(secretKey))
	if err != nil {
		return "", err
	}

	return "Bearer " + token, nil
}
