package common

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"strings"
	"time"
	"voxelprismatic/library-management-senior-project/web/user"
)

func CookieAuth(_ http.ResponseWriter, _ *http.Request) user.User {
	// TO-DO: Implement cookie authentication with JWT
	return user.User{
		Roles: user.UserRoleAdmin,
	}
}

func GenerateJWT(userData user.User, secret []byte) (string, error) {
	return generateJWTIssuedAt(userData, secret, time.Now().UTC().Unix())
}

func generateJWTIssuedAt(userData user.User, secret []byte, issuedAt int64) (string, error) {
	if len(secret) == 0 {
		return "", errors.New("secret is required")
	}

	header := map[string]string{
		"alg": "HS256",
		"typ": "JWT",
	}
	claims := map[string]string{
		"user_id": userData.ID.String(),
		"iat":     strconv.FormatInt(issuedAt, 10),
		"roles":   strconv.Itoa(int(userData.Roles)),
	}

	enc := base64.RawURLEncoding
	headerData, err := json.Marshal(header)
	if err != nil {
		return "", err
	}
	claimsData, err := json.Marshal(claims)
	if err != nil {
		return "", err
	}

	unsignedToken := enc.EncodeToString(headerData) + "." + enc.EncodeToString(claimsData)
	sig := hmac.New(sha256.New, secret)
	_, _ = sig.Write([]byte(unsignedToken))
	signature := enc.EncodeToString(sig.Sum(nil))
	return unsignedToken + "." + signature, nil
}

func validateJWTSignature(token string, secret []byte) bool {
	parts := strings.Split(token, ".")
	if len(parts) != 3 || len(secret) == 0 {
		return false
	}

	sig := hmac.New(sha256.New, secret)
	_, _ = sig.Write([]byte(parts[0] + "." + parts[1]))
	expected := base64.RawURLEncoding.EncodeToString(sig.Sum(nil))
	return hmac.Equal([]byte(parts[2]), []byte(expected))
}
