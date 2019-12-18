package server

import (
	"fmt"
	"log"
	"net/http"
	"testing"
)

func TestParseClaims(t *testing.T) {
	type args struct {
		tokenStr   string
		signingKey []byte
	}
	guestToken, err := GenerateToken("guest", false, 60, []byte("guest"))
	if err != nil {
		t.Errorf("GenerateToken() error = %v, wantErr nil", err)
		return
	}
	guestArgs := args{guestToken, []byte("guest")}
	adminToken, err := GenerateToken("admin", true, 60, []byte("admin"))
	if err != nil {
		t.Errorf("GenerateToken() error = %v, wantErr nil", err)
		return
	}
	adminArgs := args{adminToken, []byte("admin")}
	tests := []struct {
		name    string
		args    args
		want    *CustomClaims
		wantErr bool
	}{
		{"ExpiredToken",
			args{"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VybmFtZSI6ImF2ZXJhZ2Vqb2UiLCJpc19hZG1pbiI6ZmFsc2UsImlhdCI6MTU3NjYxNjE5MiwiZXhwIjoxNTc2NjE3MzkyfQ.x8XLjenNF4jF6tKempqE7PZUruM-Mgopf0WE4092Nlk",
				[]byte("secretKey")},
			nil,
			true,
		},
		{"ForgedToken",
			args{"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VybmFtZSI6ImF2ZXJhZ2Vqb2UiLCJpc19hZG1pbiI6ZmFsc2UsImlhdCI6MTU3NjYxNjE5MiwiZXhwIjoxNTc2NjE3MzkyfQ.x8XLjenNF4jF6tKempqE7PZUruM-Mgopf0WE403dqkk",
				[]byte("secretKey")},
			nil,
			true,
		},
		{"GuestToken",
			guestArgs,
			&CustomClaims{Username: "guest", IsAdmin: false},
			false,
		},
		{"AdminToken",
			adminArgs,
			&CustomClaims{Username: "admin", IsAdmin: true},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParseClaims(tt.args.tokenStr, tt.args.signingKey)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseClaims() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got == nil {
				if tt.want == nil {
					return
				}
				t.Errorf("ParseClaims() = %v, want %v", got, tt.want)
			}
			if got.Username != tt.want.Username ||
				got.IsAdmin != tt.want.IsAdmin {
				t.Errorf("ParseClaims() = %v, want %v", got, tt.want)
			}
		})
	}
}

func ExampleGenerateToken() {
	username := "mary.smith"
	isAdmin := false
	expiresInSeconds := 120
	secretKey := []byte("secretKey")
	token, err := GenerateToken(username, isAdmin, expiresInSeconds, secretKey)
	if err != nil {
		log.Fatal(err)
	}

	// For Bearer Authorization, add the token to your http request
	r, err := http.NewRequest(http.MethodGet, "/protectedAPI", nil)
	if err != nil {
		log.Fatal(err)
	}
	var bearer = "Bearer " + token
	r.Header.Add("Authorization", bearer)
	// Output:
}

func ExampleParseClaims() {
	username := "mary.smith"
	isAdmin := false
	expiresInSeconds := 120
	secretKey := []byte("secretKey")
	token, err := GenerateToken(username, isAdmin, expiresInSeconds, secretKey)
	if err != nil {
		log.Fatal(err)
	}
	claims, err := ParseClaims(token, secretKey)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("username: %s\n", claims.Username)
	fmt.Printf("isAdmin: %t\n", claims.IsAdmin)
	// Output:
	// username: mary.smith
	// isAdmin: false
}
