package users

import "github.com/gin-gonic/gin"


func GetAllUsers2(c *gin.Context) { GetAllUsers2Service(c) }
func GetAllUsers(c *gin.Context)  { GetAllUsersService(c) }
func GetAllUsers3(c *gin.Context) { GetAllUsers3Service(c) }
func GetUserByID(c *gin.Context)  { GetUserByIDService(c) }
func CreateUser(c *gin.Context)   { CreateUserService(c) }
func UpdateUser(c *gin.Context)   { UpdateUserService(c) }
func DeleteUser(c *gin.Context)   { DeleteUserService(c) }
func SearchUsers(c *gin.Context)  { SearchUsersService(c) }
func GetUserStats(c *gin.Context) { GetUserStatsService(c) }
func Login(c *gin.Context)        { LoginService(c) }
func GetProfile(c *gin.Context)   { GetProfileService(c) }
