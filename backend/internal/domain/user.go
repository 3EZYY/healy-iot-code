package domain

import "time"

// User merepresentasikan pengguna yang bisa login ke dashboard.
type User struct {
	ID        string    `json:"id"`
	Username  string    `json:"username"`
	Password  string    `json:"-"` // Never serialize bcrypt hash to JSON
	CreatedAt time.Time `json:"created_at"`
}

// LoginRequest adalah payload untuk endpoint POST /api/auth/login.
type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// LoginResponse adalah respons setelah login berhasil.
type LoginResponse struct {
	Token     string `json:"token"`
	ExpiresAt int64  `json:"expires_at"` // Unix timestamp
}
