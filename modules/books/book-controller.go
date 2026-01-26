package books

import "github.com/gin-gonic/gin"

type BookController interface {
	GetList(ctx *gin.Context)
	GetList2(ctx *gin.Context)
	GetByID(ctx *gin.Context)
	Create(ctx *gin.Context)
	Update(ctx *gin.Context)
	Delete(ctx *gin.Context)
	BulkDelete(ctx *gin.Context)
	Search(ctx *gin.Context)
	UploadImage(ctx *gin.Context)
}

type bookController struct {
	service BookService
}

func NewBookController(service BookService) BookController {
	return &bookController{service: service}
}

func (c *bookController) GetList(ctx *gin.Context)     { c.service.GetList(ctx) }
func (c *bookController) GetList2(ctx *gin.Context)    { c.service.GetList2(ctx) }
func (c *bookController) GetByID(ctx *gin.Context)     { c.service.GetByID(ctx) }
func (c *bookController) Create(ctx *gin.Context)      { c.service.Create(ctx) }
func (c *bookController) Update(ctx *gin.Context)      { c.service.Update(ctx) }
func (c *bookController) Delete(ctx *gin.Context)      { c.service.Delete(ctx) }
func (c *bookController) BulkDelete(ctx *gin.Context)  { c.service.BulkDelete(ctx) }
func (c *bookController) Search(ctx *gin.Context)      { c.service.Search(ctx) }
func (c *bookController) UploadImage(ctx *gin.Context) { c.service.UploadImage(ctx) }
