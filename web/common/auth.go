package common

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"errors"
	"net/http"
	"time"
	"voxelprismatic/library-management-senior-project/web/user"
)

// Token validity period in seconds (1 hour).
const jwtLifetimeSeconds int64 = 3600

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
	claims := map[string]any{
		"user_id": userData.ID.String(),
		"iat":     issuedAt,
		"exp":     issuedAt + jwtLifetimeSeconds,
		"roles":   userData.Roles,
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
	sig.Write([]byte(unsignedToken))
	signature := enc.EncodeToString(sig.Sum(nil))
	return unsignedToken + "." + signature, nil
}
