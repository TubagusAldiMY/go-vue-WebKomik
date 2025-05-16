package middleware

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/TubagusAldiMY/go-vue-WebKomik/webkomik-backend/internal/config" // Sesuaikan path modul Anda
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

// Claims adalah struktur untuk custom claims di JWT Anda.
// Supabase menggunakan 'sub' untuk user ID.
// Anda bisa menambahkan 'role' atau claims lain jika sudah dikonfigurasi di Supabase JWT.
type Claims struct {
	UserID string `json:"sub"`            // Standar claim untuk Subject (User ID)
	Role   string `json:"role,omitempty"` // Contoh custom claim 'role'
	// Tambahkan claims lain yang mungkin ada di Supabase JWT Anda
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

		// Jika token valid, simpan informasi pengguna (misalnya UserID dan Role) ke dalam context Gin
		// sehingga bisa diakses oleh handler berikutnya.
		c.Set("userID", claims.UserID)
		if claims.Role != "" {
			c.Set("userRole", claims.Role)
		}

		c.Next() // Lanjutkan ke handler berikutnya
	}
}

// AdminRoleMiddleware memeriksa apakah pengguna yang terotentikasi memiliki peran 'admin'.
// Middleware ini harus dijalankan SETELAH AuthMiddleware.
func AdminRoleMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Ambil userRole dari context yang sudah di-set oleh AuthMiddleware
		userRoleVal, exists := c.Get("userRole")
		if !exists {
			// Ini bisa terjadi jika AuthMiddleware tidak dijalankan sebelumnya
			// atau jika 'userRole' tidak ada di token dan tidak di-set default.
			// AuthMiddleware kita saat ini tidak men-set default jika role tidak ada, jadi ini skenario valid.
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "Peran pengguna tidak ditemukan di konteks. Akses ditolak."})
			return
		}

		userRole, ok := userRoleVal.(string)
		if !ok {
			// Tipe userRole di context bukan string, ini seharusnya tidak terjadi.
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Kesalahan internal: tipe peran pengguna tidak valid."})
			return
		}

		if strings.ToLower(userRole) != "admin" {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "Akses ditolak: membutuhkan peran admin."})
			return
		}

		c.Next() // Pengguna adalah admin, lanjutkan ke handler berikutnya.
	}
}
