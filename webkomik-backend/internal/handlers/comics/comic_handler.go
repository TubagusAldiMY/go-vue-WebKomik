package comics

import (
	"log"
	"net/http"
	"strconv"

	"github.com/TubagusAldiMY/go-vue-WebKomik/webkomik-backend/internal/database"
	"github.com/TubagusAldiMY/go-vue-WebKomik/webkomik-backend/internal/middleware" // Import middleware for role checks
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

// CreateComicHandler menangani pembuatan komik baru.
// Dapat diakses oleh admin dan creator.
func CreateComicHandler(c *gin.Context) {
	var input CreateComicInput // Struct untuk binding dan validasi

	// Bind JSON body ke struct input dan validasi
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Input tidak valid", "details": err.Error()})
		return
	}

	// Ambil userID dari context yang di-set oleh AuthMiddleware
	userIDVal, exists := c.Get("userID")
	if !exists {
		// Seharusnya tidak terjadi jika AuthMiddleware bekerja dengan benar
		c.JSON(http.StatusInternalServerError, gin.H{"error": "UserID tidak ditemukan di context"})
		return
	}
	userID, ok := userIDVal.(string)
	if !ok || userID == "" {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Format UserID tidak valid"})
		return
	}

	// Buat objek models.Comic dari input
	comicData := models.Comic{
		Title:         input.Title,
		Description:   input.Description,
		AuthorName:    input.AuthorName,
		GenreID:       input.GenreID,
		CoverImageURL: input.CoverImageURL,
		// UploadedByAdminID akan diisi oleh fungsi database dari parameter userID
	}

	createdComic, err := database.CreateComic(c.Request.Context(), comicData, userID)
	if err != nil {
		c.Error(err)
		log.Printf("Error saat membuat komik: %v\nInput: %+v\nUserID: %s\n", err, input, userID)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menyimpan komik baru"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"data": createdComic})
}

// UpdateComicHandler menangani pembaruan komik yang sudah ada.
// Dapat diakses oleh admin dan creator, dengan tambahan pemeriksaan creator hanya bisa update komik mereka sendiri.
func UpdateComicHandler(c *gin.Context) {
	// 1. Ambil ID komik dari URL
	comicIDStr := c.Param("id")
	comicID, err := strconv.ParseInt(comicIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID komik tidak valid"})
		return
	}

	// 2. Ambil data komik yang ada untuk pemeriksaan hak akses
	existingComic, err := database.GetComicByID(c.Request.Context(), comicID)
	if err != nil {
		c.Error(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengambil detail komik"})
		return
	}
	if existingComic == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Komik tidak ditemukan"})
		return
	}

	// 3. Ambil userID dan userRole dari context
	userIDVal, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "UserID tidak ditemukan di context"})
		return
	}
	userID, ok := userIDVal.(string)
	if !ok || userID == "" {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Format UserID tidak valid"})
		return
	}

	// 4. Pemeriksaan hak akses: admin dapat update semua komik, creator hanya miliknya
	if !middleware.UserHasRole(c, middleware.RoleAdmin) {
		// Jika bukan admin, periksa apakah user adalah pemilik komik
		if existingComic.UploadedByAdminID == nil || *existingComic.UploadedByAdminID != userID {
			c.JSON(http.StatusForbidden, gin.H{"error": "Anda tidak memiliki hak untuk memperbarui komik ini"})
			return
		}
	}

	// 5. Bind dan validasi input
	var input UpdateComicInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Input tidak valid", "details": err.Error()})
		return
	}

	// 6. Siapkan updates yang akan dikirim ke database
	updates := make(map[string]interface{})

	// Hanya masukkan field yang ada di request
	if input.Title != nil {
		updates["title"] = *input.Title
	}
	if input.Description != nil {
		updates["description"] = input.Description
	}
	if input.AuthorName != nil {
		updates["author_name"] = input.AuthorName
	}
	if input.GenreID != nil {
		updates["genre_id"] = input.GenreID
	}
	if input.CoverImageURL != nil {
		updates["cover_image_url"] = input.CoverImageURL
	}

	// 7. Lakukan update di database jika ada field yang diupdate
	if len(updates) == 0 {
		// Tidak ada perubahan yang diminta
		c.JSON(http.StatusOK, gin.H{"data": existingComic, "message": "Tidak ada perubahan yang dilakukan"})
		return
	}

	// 8. Lakukan update di database
	updatedComic, err := database.UpdateComic(c.Request.Context(), comicID, updates, userID)
	if err != nil {
		c.Error(err)
		log.Printf("Error saat memperbarui komik ID %d: %v\nInput: %+v\nUserID: %s\n", comicID, err, input, userID)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal memperbarui komik"})
		return
	}

	// 9. Kembalikan data komik yang telah diperbarui
	c.JSON(http.StatusOK, gin.H{"data": updatedComic})
}
