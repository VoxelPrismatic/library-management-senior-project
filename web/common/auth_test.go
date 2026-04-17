package common

import (
	"encoding/base64"
	"encoding/json"
	"strings"
	"testing"

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

	claims := map[string]string{}
	if err := json.Unmarshal(decodedClaims, &claims); err != nil {
		t.Fatalf("failed to unmarshal claims: %v", err)
	}

	if claims["user_id"] != u.ID.String() {
		t.Fatalf("expected user_id %q, got %q", u.ID.String(), claims["user_id"])
	}
	if claims["iat"] != "1712345678" {
		t.Fatalf("expected iat %q, got %q", "1712345678", claims["iat"])
	}
	if claims["roles"] != "2" {
		t.Fatalf("expected roles %q, got %q", "2", claims["roles"])
	}
	if !validateJWTSignature(token, secret) {
		t.Fatal("expected valid JWT signature")
	}
}

func TestGenerateJWTRequiresSecret(t *testing.T) {
	_, err := generateJWTIssuedAt(user.User{}, nil, 1712345678)
	if err == nil {
		t.Fatal("expected error when secret is empty")
	}
}
