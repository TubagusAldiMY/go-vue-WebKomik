package middleware

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/TubagusAldiMY/go-vue-WebKomik/webkomik-backend/internal/config" // Sesuaikan path modul Anda
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

// Role constants for the application
const (
	RoleAdmin   = "admin"
	RoleCreator = "creator"
	RoleUser    = "user"
)

type AppMetadata struct {
	Role string `json:"role,omitempty"`
	// Tambahkan field lain dari app_metadata jika perlu
}

// Claims struct
type Claims struct {
	UserID      string      `json:"sub"`
	AppMeta     AppMetadata `json:"app_metadata,omitempty"` // Untuk membaca dari app_metadata
	jwt.RegisteredClaims
}

// AuthMiddleware membuat Gin middleware untuk otentikasi JWT.
// Ia akan menggunakan JWT secret dari konfigurasi.
func AuthMiddleware(cfg *config.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Authorization header dibutuhkan"})
			return
		}

		// Token biasanya dikirim sebagai "Bearer <token>"
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Format Authorization header salah. Harusnya: Bearer <token>"})
			return
		}
		tokenString := parts[1]

		// Parse dan validasi token
		claims := &Claims{}
		token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
			// Pastikan algoritma signing adalah yang diharapkan (misalnya HMAC)
			// Supabase biasanya menggunakan HS256 untuk JWT yang ditandatangani dengan secret.
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("metode signing tidak diharapkan: %v", token.Header["alg"])
			}
			return []byte(cfg.SupabaseJWTSecret), nil
		})

		if err != nil {
			if err == jwt.ErrSignatureInvalid {
				c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Signature token tidak valid"})
				return
			}
			// Tangani error lain seperti token kedaluwarsa, token belum aktif, dll.
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Token tidak valid atau kedaluwarsa: " + err.Error()})
			return
		}

		if !token.Valid {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Token tidak valid"})
			return
		}

		// Set user ID from the token
		c.Set("userID", claims.UserID)

		// Jika token valid, simpan informasi pengguna (UserID dan Role) ke dalam context Gin
		if claims.AppMeta.Role != "" {
			c.Set("userRole", claims.AppMeta.Role)
		} else {
			// Set default role to 'user' if no role is specified
			c.Set("userRole", RoleUser)
		}

		c.Next() // Lanjutkan ke handler berikutnya
	}
}

// RoleMiddleware creates a middleware that checks if the authenticated user has one of the allowed roles
// This middleware must be run AFTER AuthMiddleware
func RoleMiddleware(allowedRoles ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get userRole from context set by AuthMiddleware
		userRoleVal, exists := c.Get("userRole")
		if !exists {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "Peran pengguna tidak ditemukan di konteks. Akses ditolak."})
			return
		}

		userRole, ok := userRoleVal.(string)
		if !ok {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Kesalahan internal: tipe peran pengguna tidak valid."})
			return
		}

		// Convert userRole to lowercase for case-insensitive comparison
		userRole = strings.ToLower(userRole)

		// Check if the user's role is in the list of allowed roles
		for _, role := range allowedRoles {
			if userRole == strings.ToLower(role) {
				c.Next() // User has an allowed role, continue to the next handler
				return
			}
		}

		// If we get here, the user doesn't have any of the allowed roles
		c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "Akses ditolak: tidak memiliki peran yang diizinkan."})
	}
}

// AdminRoleMiddleware is a convenience wrapper for RoleMiddleware that only allows admins
// This is kept for backward compatibility
func AdminRoleMiddleware() gin.HandlerFunc {
	return RoleMiddleware(RoleAdmin)
}

// AdminOrCreatorRoleMiddleware allows both admin and creator roles
// This is useful for endpoints that should be accessible by both admins and content creators
func AdminOrCreatorRoleMiddleware() gin.HandlerFunc {
	return RoleMiddleware(RoleAdmin, RoleCreator)
}

// UserHasRole is a helper function to check if a user has a specific role
// This can be used in handlers for more granular control
func UserHasRole(c *gin.Context, role string) bool {
	userRoleVal, exists := c.Get("userRole")
	if !exists {
		return false
	}

	userRole, ok := userRoleVal.(string)
	if !ok {
		return false
	}

	return strings.ToLower(userRole) == strings.ToLower(role)
}
