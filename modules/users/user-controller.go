package users

import "github.com/gin-gonic/gin"

type UserController interface {
	Create(ctx *gin.Context)
	Update(ctx *gin.Context)
	Delete(ctx *gin.Context)
	GetList(ctx *gin.Context)
	GetList2(ctx *gin.Context)
	GetByID(ctx *gin.Context)
	Search(ctx *gin.Context)
	Login(ctx *gin.Context)
	GetProfile(ctx *gin.Context)
	GetStats(ctx *gin.Context)
}

type userController struct {
	service UserService
}

func NewUserController(service UserService) UserController {
	return &userController{service: service}
}

func (c *userController) Create(ctx *gin.Context)     { c.service.Create(ctx) }
func (c *userController) Update(ctx *gin.Context)     { c.service.Update(ctx) }
func (c *userController) Delete(ctx *gin.Context)     { c.service.Delete(ctx) }
func (c *userController) GetList(ctx *gin.Context)    { c.service.GetList(ctx) }
func (c *userController) GetList2(ctx *gin.Context)   { c.service.GetList2(ctx) }
func (c *userController) GetByID(ctx *gin.Context)    { c.service.GetByID(ctx) }
func (c *userController) Search(ctx *gin.Context)     { c.service.Search(ctx) }
func (c *userController) Login(ctx *gin.Context)      { c.service.Login(ctx) }
func (c *userController) GetProfile(ctx *gin.Context) { c.service.GetProfile(ctx) }
func (c *userController) GetStats(ctx *gin.Context)   { c.service.GetStats(ctx) }
