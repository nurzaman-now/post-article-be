package controllers

import (
	"database/sql"
	"fmt"
	"post-article-be/database"
	"post-article-be/handlers"
	"post-article-be/models"
	"post-article-be/requests/post"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

type PostController struct{}

// InitVoucherController membuat controller baru dengan dependency injection
func InitPostController() *PostController {
	return &PostController{}
}

func (c *PostController) Create(ctx *gin.Context) {
	// Menginisialisasi struct untuk menyimpan data dari request
	var req post.CreatePostRequest

	// Validasi Input JSON
	// ShouldBindJSON akan mencocokkan request body dengan CreatePostRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		handlers.ResponseValidation(ctx, err)
		return
	}
	// Inisialisasi struct models.Post dari data request
	post := models.Post{
		Title:       req.Title,
		Content:     req.Content,
		Category:    req.Category,
		CreatedDate: time.Now(),
		Status:      req.Status,
	}
	// Query SQL untuk menyimpan data ke database
	_, err := database.DB.Exec(
		"INSERT INTO posts (title, content, category, created_date, status) VALUES (?, ?, ?, ?, ?)",
		post.Title, post.Content, post.Category, post.CreatedDate, post.Status,
	)
	if err != nil {
		handlers.ResponseError(ctx, nil, "Gagal membuat article")
		return
	}

	// Kembalikan berhasil
	handlers.ResponseSuccess(ctx, nil, gin.H{}, "Berhasil dibuat article", 201)
}

func (c *PostController) Get(ctx *gin.Context) {
	limitStr := ctx.Param("limit")
	offsetStr := ctx.Param("offset")

	limit, err := strconv.Atoi(limitStr)
	if err != nil {
		handlers.ResponseError(ctx, nil, "Limit harus berupa angka", 422)
		return
	}

	offset, err := strconv.Atoi(offsetStr)
	if err != nil {
		handlers.ResponseError(ctx, nil, "Offset harus berupa angka", 422)
		return
	}

	status := ctx.Query("status")

	var rows *sql.Rows
	var queryErr error
	if status != "" {
		rows, queryErr = database.DB.Query("SELECT id, title, content, category, created_date, updated_date, status FROM posts WHERE status = ? LIMIT ? OFFSET ?", status, limit, offset)
	} else {
		rows, queryErr = database.DB.Query("SELECT id, title, content, category, created_date, updated_date, status FROM posts LIMIT ? OFFSET ?", limit, offset)
	}
	if queryErr != nil {
		handlers.ResponseError(ctx, nil, "Gagal mendapatkan article")
		return
	}
	defer rows.Close()

	posts := []models.Post{}
	// Iterasi setiap baris hasil query
	for rows.Next() {
		var post models.Post
		// Scan nilai dari setiap kolom ke struct GetPostResponse
		err := rows.Scan(&post.ID, &post.Title, &post.Content, &post.Category, &post.CreatedDate, &post.UpdatedDate, &post.Status)
		if err != nil {
			handlers.ResponseError(ctx, nil, "Gagal mendapatkan article")
			return
		}
		posts = append(posts, post)
	}
	if err = rows.Err(); err != nil {
		handlers.ResponseError(ctx, nil, "Gagal mendapatkan article")
		return
	}

	// get total data per status
	rows, err = database.DB.Query("SELECT COUNT(*) FROM posts where status = ?", status)
	if err != nil {
		handlers.ResponseError(ctx, nil, "Gagal mendapatkan article")
		return
	}
	defer rows.Close()

	total := 0
	for rows.Next() {
		err := rows.Scan(&total)
		if err != nil {
			handlers.ResponseError(ctx, nil, "Gagal mendapatkan article")
			return
		}
	}
	if err = rows.Err(); err != nil {
		handlers.ResponseError(ctx, nil, "Gagal mendapatkan article")
		return
	}

	// Kembalikan berhasil
	handlers.ResponseSuccess(ctx, gin.H{"total": total}, posts, "Berhasil mendapatkan article")
}

// RouteGet bertindak sebagai dispatcher manual untuk HTTP GET request.
// Ini diperlukan karena Gin tidak mendukung wildcard routing dengan parameter nama yang berbeda
// pada tingkat segmen path yang sama (misal: /:limit/:offset vs /:id) dan akan memicu panic saat start-up.
// Dengan ini, kita menggunakan catch-all wildcard "/*action" dan memetakan parameter secara manual.
func (c *PostController) RouteGet(ctx *gin.Context) {
	action := ctx.Param("action")
	action = strings.Trim(action, "/")
	if action == "" {
		handlers.ResponseError(ctx, nil, "Endpoint tidak valid", 404)
		return
	}
	parts := strings.Split(action, "/")
	if len(parts) == 2 {
		ctx.Params = append(ctx.Params, gin.Param{Key: "limit", Value: parts[0]})
		ctx.Params = append(ctx.Params, gin.Param{Key: "offset", Value: parts[1]})
		c.Get(ctx)
	} else if len(parts) == 1 {
		ctx.Params = append(ctx.Params, gin.Param{Key: "id", Value: parts[0]})
		c.Show(ctx)
	} else {
		handlers.ResponseError(ctx, nil, "Endpoint tidak ditemukan", 404)
	}
}

func (c *PostController) Show(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		handlers.ResponseError(ctx, nil, "ID harus berupa angka", 422)
		return
	}

	// Query SQL untuk mendapatkan data dari tabel posts berdasarkan id
	rows, err := database.DB.Query("SELECT id, title, content, category, created_date, updated_date, status FROM posts WHERE id = ?", id)
	if err != nil {
		handlers.ResponseError(ctx, nil, "Gagal mendapatkan article")
		return
	}
	// Menutup koneksi database setelah selesai
	defer rows.Close()

	post := models.Post{}
	found := false
	// Iterasi setiap baris hasil query
	for rows.Next() {
		found = true
		// Scan nilai dari setiap kolom ke struct Post
		err := rows.Scan(&post.ID, &post.Title, &post.Content, &post.Category, &post.CreatedDate, &post.UpdatedDate, &post.Status)
		if err != nil {
			handlers.ResponseError(ctx, nil, "Gagal mendapatkan article")
			return
		}
	}
	if err = rows.Err(); err != nil {
		handlers.ResponseError(ctx, nil, "Gagal mendapatkan article")
		return
	}

	if !found {
		handlers.ResponseError(ctx, nil, "Article tidak ditemukan", 404)
		return
	}

	// Kembalikan berhasil
	handlers.ResponseSuccess(ctx, gin.H{"total": 1}, post, "Berhasil mendapatkan article")
}

func (c *PostController) Update(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		handlers.ResponseError(ctx, nil, "ID harus berupa angka", 422)
		return
	}

	// Menginisialisasi struct untuk menyimpan data dari request
	var req post.UpdatePostRequest

	// Validasi Input JSON
	// ShouldBindJSON akan mencocokkan request body dengan UpdatePostRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		handlers.ResponseValidation(ctx, err)
		return
	}

	// Kondisi jika ada field yang tidak dikirim untuk diupdate
	var queryParts []string
	var args []any

	if req.Title != "" {
		queryParts = append(queryParts, "title = ?")
		args = append(args, req.Title)
	}
	if req.Content != "" {
		queryParts = append(queryParts, "content = ?")
		args = append(args, req.Content)
	}
	if req.Category != "" {
		queryParts = append(queryParts, "category = ?")
		args = append(args, req.Category)
	}
	if req.Status != "" {
		queryParts = append(queryParts, "status = ?")
		args = append(args, req.Status)
	}

	// Jika tidak ada field yang dikirim untuk diupdate, langsung sukses tanpa eksekusi query
	if len(queryParts) == 0 {
		handlers.ResponseSuccess(ctx, nil, gin.H{}, "Tidak ada data yang diperbarui")
		return
	}

	queryParts = append(queryParts, "updated_date = ?")
	args = append(args, time.Now())

	query := fmt.Sprintf("UPDATE posts SET %s WHERE id = ?", strings.Join(queryParts, ", "))
	args = append(args, id)

	_, err = database.DB.Exec(query, args...)
	if err != nil {
		handlers.ResponseError(ctx, nil, "Gagal memperbarui article")
		return
	}

	// Kembalikan berhasil
	handlers.ResponseSuccess(ctx, nil, gin.H{}, "Berhasil memperbarui article")
}

func (c *PostController) Delete(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		handlers.ResponseError(ctx, nil, "ID harus berupa angka", 422)
		return
	}

	// Query SQL untuk menghapus data dari tabel posts berdasarkan id
	_, err = database.DB.Exec("DELETE FROM posts WHERE id = ?", id)
	if err != nil {
		handlers.ResponseError(ctx, nil, "Gagal menghapus article")
		return
	}

	// Kembalikan berhasil
	handlers.ResponseSuccess(ctx, nil, gin.H{}, "Berhasil menghapus article")
}
