package comics

import (
	"log"
	"net/http"
	"strconv"

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

// GetComicDetailHandler menangani permintaan untuk mendapatkan detail satu komik.
func GetComicDetailHandler(c *gin.Context) {
	comicIDStr := c.Param("id") // Mengambil "id" dari URL path, contoh: /comics/123
	comicID, err := strconv.ParseInt(comicIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID komik tidak valid"})
		return
	}

	// 1. Ambil detail komik dasar
	comic, err := database.GetComicByID(c.Request.Context(), comicID)
	if err != nil {
		c.Error(err) // Log error server
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengambil detail komik"})
		return
	}
	if comic == nil { // Komik tidak ditemukan (GetComicByID mengembalikan nil, nil)
		c.JSON(http.StatusNotFound, gin.H{"error": "Komik tidak ditemukan"})
		return
	}

	// 2. Ambil chapters untuk komik ini
	chapters, err := database.GetChaptersByComicID(c.Request.Context(), comicID)
	if err != nil {
		c.Error(err)
		// Mungkin tidak fatal jika chapter gagal diambil, tergantung kebutuhan
		// c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengambil chapter komik"})
		// return
		log.Printf("Peringatan: Gagal mengambil chapters untuk comic ID %d: %v\n", comicID, err)
		// Tetap lanjutkan dengan chapters kosong jika ada error, atau Anda bisa memilih untuk gagal total.
		comic.Chapters = []models.Chapter{}
	} else {
		// 3. Untuk setiap chapter, ambil halamannya (pages)
		for i, chapter := range chapters {
			pages, err := database.GetPagesByChapterID(c.Request.Context(), chapter.ID)
			if err != nil {
				c.Error(err)
				log.Printf("Peringatan: Gagal mengambil pages untuk chapter ID %d: %v\n", chapter.ID, err)
				// chapters[i].Pages tetap nil atau []models.Page{} (default)
				chapters[i].Pages = []models.Page{} // Pastikan array kosong jika gagal
			} else {
				chapters[i].Pages = pages
			}
		}
		comic.Chapters = chapters
	}

	c.JSON(http.StatusOK, gin.H{"data": comic})
}
