package db

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"errors"
	"net/http"
	"strings"
	"time"
)

const JWT_LIFETIME = int64(WEEK/time.Second) * 4 // One month login time
const JWT_SECRET = "The Super Duper Secret Hash that should be stored elsewhere:tm:"

var JWT_ENC = base64.RawURLEncoding

func CookieAuth(w http.ResponseWriter, r *http.Request) UserPartial {
	cookie, err := r.Cookie("tok")
	if err != nil {
		w.Header().Set("X-Auth-Stage", "cookie.get")
		w.Header().Set("X-Auth-Reason", err.Error())
		return UserPartial{}
	}

	data, stage, err := ValidateJWT(cookie.String())
	if err != nil {
		w.Header().Set("X-Auth-Stage", stage)
		w.Header().Set("X-Auth-Reason", err.Error())
		return UserPartial{}
	}

	if data == nil {
		w.Header().Set("X-Auth-Stage", "auth.return")
		w.Header().Set("X-Auth-Reason", "nil")
		return UserPartial{}
	}

	return *data
}

func (u *User) NewJWT() string {
	header := ToJsonB64(map[string]string{
		"alg": "HS256",
		"typ": "JWT",
	})

	issuedAt := time.Now().Unix()
	partial := u.Partial()
	partial.SetTimestamp(issuedAt)
	claims := ToJsonB64(partial)

	unsignedToken := header + "." + claims
	sig := hmac.New(sha256.New, []byte(JWT_SECRET))
	sig.Write([]byte(unsignedToken))
	sum := JWT_ENC.EncodeToString(sig.Sum(nil))
	return unsignedToken + "." + sum
}

func ToJsonB64(dataMap any) string {
	dataStr, err := json.Marshal(dataMap)
	if err != nil {
		// Should be unreachable
		panic(err)
	}

	return JWT_ENC.EncodeToString(dataStr)
}

func ValidateJWT(jwt string) (*UserPartial, string, error) {
	parts := strings.Split(jwt, ".")
	if len(parts) != 3 {
		return nil, "jwt.split", errors.New("malformed jwt")
	}

	sig := hmac.New(sha256.New, []byte(JWT_SECRET))
	_, _ = sig.Write([]byte(parts[0] + "." + parts[1]))
	expected := JWT_ENC.EncodeToString(sig.Sum(nil))
	if !hmac.Equal([]byte(parts[2]), []byte(expected)) {
		return nil, "jwt.sig", errors.New("signature mismatch")
	}

	data, err := JWT_ENC.DecodeString(parts[1])
	if err != nil {
		return nil, "jwt.b64decode", err
	}

	ret := UserPartial{}
	if err := json.Unmarshal(data, &ret); err != nil {
		return nil, "jwt.unmarshal", err
	}

	expiresAt := time.Unix(ret.ExpiresAt, 0)
	if time.Now().After(expiresAt) {
		return nil, "user.security", errors.New("expired")
	}

	return &ret, "", nil
}
