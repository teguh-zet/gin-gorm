package book_controller

import (
	"gin-gonic/database"
	"gin-gonic/helpers"
	"gin-gonic/models"
	"strconv"

	"github.com/gin-gonic/gin"
)
// get all book reguler
func GetAllBooks(ctx *gin.Context) {
	var book []models.Book

	if err := database.DB.Find(&book).Error;err!=nil{
		helpers.InternalServerError(ctx,"Failed to fect book",err.Error())
		return
	}
	helpers.SuccessResponse(ctx, "Book retrieved succesfully",book)
}

//get all book with pagination and sorting
func GetAllBooks2(c *gin.Context){
	page:= c.DefaultQuery("page","1")
	limit := c.DefaultQuery("limit","10")

	sortBy:= c.DefaultQuery("sort_by", "id")
	order := c.DefaultQuery("order", "DESC")

	pageNum, err := strconv.Atoi(page)
	if err != nil || pageNum <1{
		pageNum = 1
	}

	limitNum, err := strconv.Atoi(limit)
	if err!= nil || limitNum < 1{
		limitNum =10
	}
	if limitNum > 100{
		limitNum =100
	}

	offset := (pageNum -1) * limitNum
	// validasi sort
	allowed := map[string]bool{
		"title": true,
		"author": true,
		"id": true,

	}
	if !allowed[sortBy]{
		sortBy = "id"
	}
	if order != "ASC" && order != "DESC"{
		order = "DESC"
	}
	var books []models.Book
	var total int64

	if err := database.DB.Order(sortBy +" "+order).Offset(offset).Limit(limitNum).Find(&books).Error;
	err!= nil{
		helpers.InternalServerError(c, " failed to fetch book", err.Error())
	}
	database.DB.Model(&models.Book{}).Count(&total)
	totalPages :=(total + int64(limitNum)-1 ) / int64(limitNum)
	helpers.SuccessResponse(c,"Books retrivied successfully", gin.H{
		"data":         books,
        "pagination": gin.H{
            "total":        total,
            "page":         pageNum,
            "limit":        limitNum,
            "total_pages":  totalPages,
            "has_next":     pageNum < int(totalPages),
            "has_previous": pageNum > 1,
        },
        "sorting": gin.H{
            "sort_by": sortBy,
            "order":   order,
        },
	})
}


func GetBookByID(c *gin.Context){
	idParam := c.Param("id")
	id, err := strconv.ParseUint(idParam,10,32)
	if err!=nil {
		helpers.BadRequestError(c,"invalid book ID","ID must be a number")
		return
	}
	var user models.Book
	if err:= database.DB.First(&user,id).Error; err!=nil{
		if err.Error() == "record not found"{
			helpers.NotFoundError(c, "user not found ")
			return
		}
		helpers.InternalServerError(c,"Failed fecth user ",err.Error())
		return
	
	}
	helpers.SuccessResponse(c,"User retrieved successfully", user )
}

func CreateBook(ctx *gin.Context){
	var req models.CreateBookRequest
	if err := ctx.ShouldBindBodyWithJSON(&req);err!=nil{
		helpers.ValidationError(ctx,err.Error())
		return
	}
	book:= models.Book{
		Title: req.Title,
		Author: req.Author,
	}
	if err := database.DB.Create(&book).Error; err!=nil{
		helpers.InternalServerError(ctx,"Failed to create book",err.Error())
		return
	}
	helpers.CreatedResponse(ctx,"book created succesfully",book)
}

func UpdateBook(ctx *gin.Context){
	idParam := ctx.Param("id")
	id, err := strconv.ParseUint(idParam,10,31)
	if err!=nil {
		helpers.BadRequestError(ctx,"invalid Book ID","ID must be a number ")
		return
	}
	var book models.Book
	if err :=database.DB.First(&book,id).Error; err!=nil{
		if err.Error() == "record not found"{
			helpers.NotFoundError(ctx,"User not found")
			return
		}
		helpers.InternalServerError(ctx,"Failed to fetch book", err.Error())
		return
	}
	var req models.CreateBookRequest
	if err := ctx.ShouldBindBodyWithJSON(&req);err!= nil{
		helpers.ValidationError(ctx, err.Error())
		return
	}
	// agar tidak terjadi data kosong karena hanya update field tertentu
	updates := make(map[string]interface{})
	if req.Title != ""{
		updates["title"] = req.Title		
	}
	if req.Author != ""{
		updates["author"] = req.Author		
	}

	if len(updates) ==0{
		helpers.BadRequestError(ctx, "No field to update", "at least one field must be provided")
	}
	if err := database.DB.Model(&book).Updates(updates).Error; err!=nil{
		helpers.InternalServerError(ctx,"failed to update book",err.Error())
		return
	}

	database.DB.First(&book,id)
	helpers.SuccessResponse(ctx, "Book updated successfully",book)

}

func DeleteBook(ctx *gin.Context){
	idParam := ctx.Param("id")
	id,err:= strconv.ParseUint(idParam,10,32)
	if err!= nil{
		helpers.BadRequestError(ctx,"Invalid book id","id must be a number")
		return
	}
	var book models.Book
	if err := database.DB.First(&book,id).Error; err!= nil{
		if err.Error() == "record not found"{
			helpers.NotFoundError(ctx, "Book not found")
			return
		}
		helpers.InternalServerError(ctx, "failed to fetch book",err.Error())
	}
	if err := database.DB.Delete(&book).Error; err!=nil{
		helpers.InternalServerError(ctx, " failed to delete book", err.Error())
		return
	}
	helpers.SuccessResponse(ctx, "book deleted successfully",book)


}

func BulkDeleteBooks(c *gin.Context){
	var req struct{
		IDs []uint `json:"ids" binding:"required"` 
	}
	if err := c.ShouldBindBodyWithJSON(&req); err !=nil{
		helpers.ValidationError(c,err.Error())
		return
	}
	if err := database.DB.Delete(&[]models.Book{}, req.IDs).Error;err != nil{
		helpers.InternalServerError(c, "failed to delete book", err.Error())
		return
	}
	helpers.SuccessResponse(c,"books deleted succesfully",gin.H{"delete_count" : len(req.IDs)})

}

func SearchBooks(c *gin.Context){
	query := c.Query("title")
	if query ==""{
		helpers.BadRequestError(c, "search query required ", "parameter harus diisi dengan title")
		return
	}
	var books []models.Book
	if err := database.DB.Where("title LIKE ? OR author LIKE ?","%"+query+"%", "%"+query+"%").Find(&books).Error;
	err != nil{
		helpers.InternalServerError(c, "failed to search book ",err.Error())
		return
	}
	helpers.SuccessResponse(c, "search result", books)

}