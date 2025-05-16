package comics

import (
	"net/http"

	"github.com/TubagusAldiMY/go-vue-WebKomik/webkomik-backend/internal/database"
	"github.com/TubagusAldiMY/go-vue-WebKomik/webkomik-backend/internal/models" // Import models
	"github.com/gin-gonic/gin"
)

// GetAllComicsHandler menangani permintaan untuk mendapatkan semua komik.
func GetAllComicsHandler(c *gin.Context) {
	comicsList, err := database.GetAllComics(c.Request.Context()) // Menggunakan context dari request Gin
	if err != nil {
		c.Error(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengambil data komik"})
		return
	}

	if comicsList == nil {
		comicsList = []models.Comic{} // Gunakan models.Comic
	}

	c.JSON(http.StatusOK, gin.H{"data": comicsList})
}
