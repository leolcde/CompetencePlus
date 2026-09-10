package main

import (
	"net/http/httptest"
	"testing"
	"time"
)

func TestPassword(t *testing.T) {
	hash, _ := hashPassword("passwd")
	if hash == "passwd" {
		t.Fatal("password disclosed")
	}
	if !checkPassword(hash, "passwd") {
		t.Fatal("right password refused")
	}
	if checkPassword(hash, "autre") {
		t.Fatal("wong password accepted")
	}
}

func TestComputeAge(t *testing.T) {
	twenty_ago := time.Now().AddDate(-20, 0, 0)
	if computeAge(twenty_ago) != 20 {
		t.Fatalf("wait 20, get %d", computeAge(twenty_ago))
	}
}

func TestToken(t *testing.T) {
	t.Setenv("JWT_SECRET", "secret")

	token, _ := generateToken(7, "candidate")

	req := httptest.NewRequest("GET", "/", nil)
	req.Header.Set("Authorization", "Bearer "+token)

	id, role, err := parseToken(req)
	if err != nil || id != 7 || role != "candidate" {
		t.Fatalf("get id=%d role=%s err=%v", id, role, err)
	}
}

func TestTokenManquant(t *testing.T) {
	req := httptest.NewRequest("GET", "/", nil)
	if _, _, err := parseToken(req); err == nil {
		t.Fatal("token missing")
	}
}

func TestTokenFalsifie(t *testing.T) {
	t.Setenv("JWT_SECRET", "true-secret")
	token, _ := generateToken(1, "admin")

	t.Setenv("JWT_SECRET", "false-secret")
	req := httptest.NewRequest("GET", "/", nil)
	req.Header.Set("Authorization", "Bearer "+token)

	if _, _, err := parseToken(req); err == nil {
		t.Fatal("wong token accepted")
	}
}
