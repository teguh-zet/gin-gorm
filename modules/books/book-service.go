package books

import (
	"strconv"

	"gin-gonic/helpers"

	"github.com/cloudinary/cloudinary-go/v2"
	"github.com/cloudinary/cloudinary-go/v2/api/uploader"
	"github.com/gin-gonic/gin"
)

// get all book reguler
func GetAllBooksService(ctx *gin.Context) {
	var book []Book

	if err := helpers.DB.Find(&book).Error; err != nil {
		helpers.InternalServerError(ctx, "Failed to fect book", err.Error())
		return
	}
	helpers.SuccessResponse(ctx, "Book retrieved succesfully", book)
}

// GetAllBooks2 godoc
// @Summary      Lihat Semua Buku (Pagination & Filter)
// @Description  Menampilkan daftar buku dengan fitur pagination, limit, sorting, dan filter ketersediaan.
// @Tags         books
// @Accept       json
// @Produce      json
// @Param        page      query int    false "Halaman ke berapa (Default: 1)"
// @Param        limit     query int    false "Jumlah data per halaman (Default: 10, Max: 100)"
// @Param        sort_by   query string false "Kolom sorting (id, title, author). Default: id"
// @Param        order     query string false "Arah urutan (ASC/DESC). Default: DESC"
// @Param        available query bool   false "Filter stok tersedia? (true/false)"
// @Success      200       {object} map[string]interface{} "Data Buku dengan Pagination"
// @Failure      500       {object} map[string]interface{} "Internal Server Error"
// @Router       /books/all [get]
func GetAllBooks2Service(c *gin.Context) {
	// 1. Ambil Parameter
	page := c.DefaultQuery("page", "1")
	limit := c.DefaultQuery("limit", "10")
	sortBy := c.DefaultQuery("sort_by", "id")
	order := c.DefaultQuery("order", "DESC")

	// [FITUR BARU] Parameter available
	available := c.DefaultQuery("available", "false")

	// 2. Validasi Angka (Pagination)
	pageNum, err := strconv.Atoi(page)
	if err != nil || pageNum < 1 {
		pageNum = 1
	}

	limitNum, err := strconv.Atoi(limit)
	if err != nil || limitNum < 1 {
		limitNum = 10
	}
	if limitNum > 100 {
		limitNum = 100
	}

	offset := (pageNum - 1) * limitNum

	// 3. Validasi Sorting
	allowed := map[string]bool{
		"title":  true,
		"author": true,
		"id":     true,
	}
	if !allowed[sortBy] {
		sortBy = "id"
	}
	if order != "ASC" && order != "DESC" {
		order = "DESC"
	}

	// ==========================================================
	// 4. KONSTRUKSI QUERY (Penting!)
	// ==========================================================

	var books []Book
	var total int64

	// Kita inisialisasi query GORM tanpa mengeksekusinya dulu
	query := helpers.DB.Model(&Book{})

	// [LOGIC FILTER] Jika available=true, tambahkan WHERE stock > 0
	if available == "true" {
		query = query.Where("stock > 0")
	}

	// 5. Hitung Total Data (Sesuai Filter)
	// Penting: Count harus dilakukan SETELAH filter where, tapi SEBELUM limit/offset
	query.Count(&total)

	// 6. Eksekusi Query dengan Pagination & Sorting
	if err := query.Order(sortBy + " " + order).
		Offset(offset).
		Limit(limitNum).
		Find(&books).Error; err != nil {
		helpers.InternalServerError(c, "Failed to fetch book", err.Error())
		return
	}

	// 7. Hitung Total Page
	totalPages := (total + int64(limitNum) - 1) / int64(limitNum)

	// 8. Kirim Response
	helpers.SuccessResponse(c, "Books retrieved successfully", gin.H{
		"data": books,
		"pagination": gin.H{
			"total":        total,
			"page":         pageNum,
			"limit":        limitNum,
			"total_pages":  totalPages,
			"has_next":     pageNum < int(totalPages),
			"has_previous": pageNum > 1,
		},
		"sorting": gin.H{
			"sort_by":   sortBy,
			"order":     order,
			"available": available, // Info filter dikembalikan juga
		},
	})
}

// GetBookByID godoc
// @Summary      Lihat Detail Buku
// @Description  Menampilkan detail satu buku berdasarkan ID.
// @Tags         books
// @Accept       json
// @Produce      json
// @Param        id   path int true "Book ID"
// @Success      200  {object} Book
// @Failure      400  {object} map[string]interface{} "ID Salah"
// @Failure      404  {object} map[string]interface{} "Buku tidak ditemukan"
// @Failure      500  {object} map[string]interface{} "Internal Server Error"
// @Router       /books/{id} [get]
func GetBookByIDService(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		helpers.BadRequestError(c, "invalid book ID", "ID must be a number")
		return
	}
	var user Book
	if err := helpers.DB.First(&user, id).Error; err != nil {
		if err.Error() == "record not found" {
			helpers.NotFoundError(c, "user not found ")
			return
		}
		helpers.InternalServerError(c, "Failed fecth user ", err.Error())
		return

	}
	helpers.SuccessResponse(c, "User retrieved successfully", user)
}

// CreateBook godoc
// @Summary      Create a new book
// @Description  Create a new book with title, author, and stock
// @Tags         admin
// @Accept       json
// @Produce      json
// @Param        Authorization header string true "Bearer Token"
// @Param        request body CreateBookRequest true "Book Data"
// @Success      201  {object} Book
// @Failure      400  {object} map[string]interface{}
// @Failure      401  {object} map[string]interface{}
// @Security     BearerAuth
// @Router       /books [post]
func CreateBookService(ctx *gin.Context) {
	var req CreateBookRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		helpers.ValidationError(ctx, err.Error())
		return
	}
	book := Book{
		Title:  req.Title,
		Author: req.Author,
		Stock:  req.Stock,
	}
	if err := helpers.DB.Create(&book).Error; err != nil {
		helpers.InternalServerError(ctx, "Failed to create book", err.Error())
		return
	}
	helpers.CreatedResponse(ctx, "book created succesfully", book)
}

// UpdateBook godoc
// @Summary      Update Data Buku (Admin)
// @Description  Mengubah data buku (Judul, Penulis, Stok). Khusus Admin.
// @Tags         admin
// @Accept       json
// @Produce      json
// @Param        Authorization header string true "Bearer Token"
// @Param        id      path int true "Book ID"
// @Param        request body UpdateBookRequest true "Data Update Buku"
// @Success      200     {object} Book
// @Failure      400     {object} map[string]interface{} "Validasi Error / ID Salah"
// @Failure      401     {object} map[string]interface{} "Unauthorized"
// @Failure      403     {object} map[string]interface{} "Forbidden"
// @Failure      404     {object} map[string]interface{} "Buku tidak ditemukan"
// @Failure      500     {object} map[string]interface{} "Internal Server Error"
// @Security     BearerAuth
// @Router       /admin/books/{id} [put]
func UpdateBookService(ctx *gin.Context) {
	idParam := ctx.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		helpers.BadRequestError(ctx, "invalid Book ID", "ID must be a number ")
		return
	}
	var book Book
	if err := helpers.DB.First(&book, id).Error; err != nil {
		if err.Error() == "record not found" {
			helpers.NotFoundError(ctx, "Book not found")
			return
		}
		helpers.InternalServerError(ctx, "Failed to fetch book", err.Error())
		return
	}
	var req UpdateBookRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		helpers.ValidationError(ctx, err.Error())
		return
	}
	// agar tidak terjadi data kosong karena hanya update field tertentu
	updates := make(map[string]interface{})
	if req.Title != "" {
		updates["title"] = req.Title
	}
	if req.Author != "" {
		updates["author"] = req.Author
	}
	if req.Stock != 0 {
		updates["stock"] = req.Stock
	}
	//setidaknya ada satu field yang diupdate
	if len(updates) == 0 {
		helpers.BadRequestError(ctx, "No field to update", "at least one field must be provided")
	}
	if err := helpers.DB.Model(&book).Updates(updates).Error; err != nil {
		helpers.InternalServerError(ctx, "failed to update book", err.Error())
		return
	}

	helpers.DB.First(&book, id)
	helpers.SuccessResponse(ctx, "Book updated successfully", book)

}

// DeleteBook godoc
// @Summary      Hapus Buku (Admin)
// @Description  Menghapus buku dari database (Soft Delete). Khusus Admin.
// @Tags         admin
// @Accept       json
// @Produce      json
// @Param        Authorization header string true "Bearer Token"
// @Param        id  path int true "Book ID"
// @Success      200 {object} Book "Data buku yang dihapus"
// @Failure      400 {object} map[string]interface{} "ID Salah"
// @Failure      401 {object} map[string]interface{} "Unauthorized"
// @Failure      403 {object} map[string]interface{} "Forbidden"
// @Failure      404 {object} map[string]interface{} "Buku tidak ditemukan"
// @Failure      500 {object} map[string]interface{} "Internal Server Error"
// @Security     BearerAuth
// @Router       /admin/books/{id} [delete]
func DeleteBookService(ctx *gin.Context) {
	idParam := ctx.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		helpers.BadRequestError(ctx, "Invalid book id", "id must be a number")
		return
	}
	var book Book
	if err := helpers.DB.First(&book, id).Error; err != nil {
		if err.Error() == "record not found" {
			helpers.NotFoundError(ctx, "Book not found")
			return
		}
		helpers.InternalServerError(ctx, "failed to fetch book", err.Error())
	}
	if err := helpers.DB.Delete(&book).Error; err != nil {
		helpers.InternalServerError(ctx, " failed to delete book", err.Error())
		return
	}
	helpers.SuccessResponse(ctx, "book deleted successfully", book)

}

func BulkDeleteBooksService(c *gin.Context) {
	var req struct {
		IDs []uint `json:"ids" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		helpers.ValidationError(c, err.Error())
		return
	}
	if err := helpers.DB.Delete(&[]Book{}, req.IDs).Error; err != nil {
		helpers.InternalServerError(c, "failed to delete book", err.Error())
		return
	}
	helpers.SuccessResponse(c, "books deleted succesfully", gin.H{"delete_count": len(req.IDs)})

}

// SearchBooks godoc
// @Summary      Cari Buku
// @Description  Mencari buku berdasarkan keyword yang cocok dengan Judul ATAU Penulis.
// @Tags         books
// @Accept       json
// @Produce      json
// @Param        title query string true "Keyword pencarian (Judul atau Penulis)"
// @Success      200   {array}  Book
// @Failure      400   {object} map[string]interface{} "Parameter title wajib diisi"
// @Failure      500   {object} map[string]interface{} "Internal Server Error"
// @Router       /books/search [get]
func SearchBooksService(c *gin.Context) {
	query := c.Query("title")
	if query == "" {
		helpers.BadRequestError(c, "search query required ", "parameter harus diisi dengan title")
		return
	}
	var books []Book
	if err := helpers.DB.Where("title LIKE ? OR author LIKE ?", "%"+query+"%", "%"+query+"%").Find(&books).Error; err != nil {
		helpers.InternalServerError(c, "failed to search book ", err.Error())
		return
	}
	helpers.SuccessResponse(c, "search result", books)

}

// UploadBookImage godoc
// @Summary      Upload Gambar Sampul Buku
// @Description  Mengunggah file gambar (jpg/png) untuk sampul buku.
// @Tags         admin
// @Accept       multipart/form-data
// @Produce      json
// @Param        Authorization header string true "Bearer Token"
// @Param        id      path   int  true "Book ID"
// @Param        image   formData file true "File Gambar"
// @Success      200     {object} map[string]interface{} "URL Gambar"
// @Failure      400     {object} map[string]interface{} "Gagal Upload"
// @Failure      500     {object} map[string]interface{} "Internal Server Error"
// @Security     BearerAuth
// @Router       /admin/books/{id}/image [patch]
func UploadBookImageService(c *gin.Context) {
	id := c.Param("id")

	// 1. Ambil file dari form-data
	fileHeader, err := c.FormFile("image")
	if err != nil {
		helpers.BadRequestError(c, "No file uploaded", "Key 'image' is required")
		return
	}

	// 2. Buka File agar bisa dibaca stream-nya
	file, err := fileHeader.Open()
	if err != nil {
		helpers.InternalServerError(c, "Failed to open file", err.Error())
		return
	}
	defer file.Close() // Wajib ditutup setelah selesai

	cldName := helpers.GetConfig("CLOUDINARY_CLOUD_NAME")
	cldKey := helpers.GetConfig("CLOUDINARY_API_KEY")
	cldSecret := helpers.GetConfig("CLOUDINARY_API_SECRET")
	cld, _ := cloudinary.NewFromParams(cldName, cldKey, cldSecret)

	// 4. Upload ke Cloudinary
	resp, err := cld.Upload.Upload(c, file, uploader.UploadParams{
		Folder:   "library-app", // Nanti gambar masuk ke folder ini di Cloudinary
		PublicID: "book-" + id,  // (Opsional) Nama file di cloud
	})

	if err != nil {
		helpers.InternalServerError(c, "Failed to upload to cloud", err.Error())
		return
	}

	// 5. Update Database (Simpan URL HTTPS dari Cloudinary)
	// resp.SecureURL adalah alamat gambar yang bisa diakses internet
	if err := helpers.DB.Model(&Book{}).Where("id = ?", id).
		Update("image_url", resp.SecureURL).Error; err != nil {
		helpers.InternalServerError(c, "Failed to update database", err.Error())
		return
	}

	helpers.SuccessResponse(c, "Image uploaded successfully", gin.H{"url": resp.SecureURL})
}
