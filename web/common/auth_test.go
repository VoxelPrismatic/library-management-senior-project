package common

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"strings"
	"testing"
	"time"

	"github.com/google/uuid"

	"voxelprismatic/library-management-senior-project/web/user"
)

func TestGenerateJWTIncludesRequiredClaims(t *testing.T) {
	u := user.User{
		ID:    uuid.MustParse("11111111-2222-3333-4444-555555555555"),
		Roles: user.UserRoleLibrarian,
	}
	secret := []byte("test-secret")
	issuedAt := int64(1712345678)

	token, err := generateJWTIssuedAt(u, secret, issuedAt)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	parts := strings.Split(token, ".")
	if len(parts) != 3 {
		t.Fatalf("expected 3 JWT parts, got %d", len(parts))
	}

	decodedClaims, err := base64.RawURLEncoding.DecodeString(parts[1])
	if err != nil {
		t.Fatalf("failed to decode claims: %v", err)
	}

	claimsAny := map[string]any{}
	if err := json.Unmarshal(decodedClaims, &claimsAny); err != nil {
		t.Fatalf("failed to unmarshal claims: %v", err)
	}

	if claimsAny["user_id"] != u.ID.String() {
		t.Fatalf("expected user_id %q, got %v", u.ID.String(), claimsAny["user_id"])
	}
	if claimsAny["iat"] != float64(1712345678) {
		t.Fatalf("expected iat %v, got %v", float64(1712345678), claimsAny["iat"])
	}
	if claimsAny["exp"] != float64(1712349278) {
		t.Fatalf("expected exp %v, got %v", float64(1712349278), claimsAny["exp"])
	}
	if claimsAny["roles"] != float64(2) {
		t.Fatalf("expected roles %v, got %v", float64(2), claimsAny["roles"])
	}

	sig := hmac.New(sha256.New, secret)
	sig.Write([]byte(parts[0] + "." + parts[1]))
	expectedSig := base64.RawURLEncoding.EncodeToString(sig.Sum(nil))
	if !hmac.Equal([]byte(parts[2]), []byte(expectedSig)) {
		t.Fatal("expected valid JWT signature")
	}
}

func TestGenerateJWTRequiresSecret(t *testing.T) {
	_, err := generateJWTIssuedAt(user.User{}, nil, 1712345678)
	if err == nil {
		t.Fatal("expected error when secret is empty")
	}
}

func TestGenerateJWTUsesCurrentTimeClaims(t *testing.T) {
	u := user.User{
		ID:    uuid.MustParse("aaaaaaaa-bbbb-cccc-dddd-eeeeeeeeeeee"),
		Roles: user.UserRoleAdmin,
	}
	secret := []byte("test-secret")

	start := time.Now().UTC().Unix()
	token, err := GenerateJWT(u, secret)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	end := time.Now().UTC().Unix()

	parts := strings.Split(token, ".")
	if len(parts) != 3 {
		t.Fatalf("expected 3 JWT parts, got %d", len(parts))
	}

	decodedClaims, err := base64.RawURLEncoding.DecodeString(parts[1])
	if err != nil {
		t.Fatalf("failed to decode claims: %v", err)
	}

	claimsAny := map[string]any{}
	if err := json.Unmarshal(decodedClaims, &claimsAny); err != nil {
		t.Fatalf("failed to unmarshal claims: %v", err)
	}

	iat, ok := claimsAny["iat"].(float64)
	if !ok {
		t.Fatalf("expected numeric iat, got %T", claimsAny["iat"])
	}
	if int64(iat) < start || int64(iat) > end {
		t.Fatalf("iat out of expected range: %v not in [%v,%v]", int64(iat), start, end)
	}

	exp, ok := claimsAny["exp"].(float64)
	if !ok {
		t.Fatalf("expected numeric exp, got %T", claimsAny["exp"])
	}
	if int64(exp-iat) != jwtLifetimeSeconds {
		t.Fatalf("unexpected exp lifetime delta: %v", int64(exp-iat))
	}
}
