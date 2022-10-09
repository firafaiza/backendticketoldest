package controller

import (
	"errors"
	"fmt"
	"path/filepath"

	"net/http"
	"strconv"

	"gorm.io/gorm"
	"ticket.narindo.com/delivery/api"

	"ticket.narindo.com/model"
	"ticket.narindo.com/usecase"

	"github.com/gin-gonic/gin"
)

type TicketController struct {
	router   *gin.Engine
	ucTicket usecase.TicketUseCase
	api.BaseApi
}

func (t *TicketController) createTicket(c *gin.Context) {
	var newTicket model.Ticket
	csId := c.PostForm("csId")
	ticketSubject := c.PostForm("ticketSubject")
	departmentId := c.PostForm("departmentId")
	ticketMessage := c.PostForm("ticketMessage")
	priorityId := c.PostForm("priorityId")
	picId := c.PostForm("picId")

	convDepartmentId, err := strconv.Atoi(departmentId)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   err,
			"message": "Failed converting department id",
		})
	}
	convPriorityId, err := strconv.Atoi(priorityId)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   err,
			"message": "Failed converting priority id",
		})
	}
	convPicId, err := strconv.Atoi(picId)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   err,
			"message": "Failed converting pic id",
		})
	}

	// ini buat gagal
	// kalo file uploadnya perlu pake form
	// sedangkan kalo userRole, dia harus dari json
	// var userRole model.UserRole
	// err = t.ParseRequestBody(c, &userRole)
	// if err != nil {
	// 	c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
	// 		"status":  "FAILED",
	// 		"message": err,
	// 	})
	// 	return
	// }

	newTicket.CsId = csId
	newTicket.TicketSubject = ticketSubject
	newTicket.DepartmentId = convDepartmentId
	newTicket.TicketMessage = ticketMessage
	newTicket.PriorityId = convPriorityId
	newTicket.PicId = convPicId

	ticket, err := t.ucTicket.CreateTicket(&newTicket)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"status":  "FAILED",
			"message": err,
		})
		return
	}

	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   err,
			"message": "Failed to upload",
		})
		return
	}

	extension := filepath.Ext(file.Filename)

	if err := c.SaveUploadedFile(file, "attachment/"+ticket.TicketId+"."+extension); err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"message": "Unable to save the file",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  "SUCCESS",
		"message": newTicket,
	})
}

func (t *TicketController) ListAllTicket(c *gin.Context) {
	result, err := t.ucTicket.ListAllTicket()
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"status":  "FAILED",
			"message": "Error when retrieving list ticket",
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"status":  "SUCCESS",
		"message": result,
	})
}

func (t *TicketController) ListTicketByUserId(c *gin.Context) {
	param := c.Request.URL.Query()
	userId := param["uid"][0]
	roleId := param["rid"][0]
	intRoleId, _ := strconv.Atoi(roleId)
	var userRole model.UserRole
	userRole.UserId = userId
	userRole.RoleId = intRoleId

	print, err := t.ucTicket.ListByUser(&userRole)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"status":  "FAILED",
			"message": err,
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"status":  "SUCCESS",
		"message": print,
	})
}

func (t *TicketController) ListTicketByDepartmentId(c *gin.Context) {
	departmentId := c.Param("id")
	integerDeptId, _ := strconv.Atoi(departmentId)
	print, err := t.ucTicket.ListByDepartment(integerDeptId)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"status":  "FAILED",
			"message": err,
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"status":  "SUCCESS",
		"message": print,
	})
}

func (t *TicketController) getTicketById(c *gin.Context) {
	ticketId := c.Param("id")
	print, err := t.ucTicket.GetById(ticketId)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
				"status":  "FAILED",
				"message": "Error ticket not found",
			})
			return
		}
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"status":  "FAILED",
			"message": err,
		})
		return

	}
	c.JSON(http.StatusOK, gin.H{
		"status":  "SUCCESS",
		"message": print,
	})
}

func (t *TicketController) getTicketSummary(c *gin.Context) {
	statusId := c.Param("status")
	convStatusId, err := strconv.Atoi(statusId)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"status":  "FAILED",
			"message": "Error converting status id",
		})
		return
	}

	var userRole model.UserRole
	err = t.ParseRequestBody(c, &userRole)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"status":  "FAILED",
			"message": err,
		})
		return
	}

	print, err := t.ucTicket.GetSummaryTicket(convStatusId, &userRole)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
				"status":  "FAILED",
				"message": "Error ticket not found",
			})
			return
		}
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"status":  "FAILED",
			"message": err,
		})
		return

	}
	c.JSON(http.StatusOK, gin.H{
		"status":  "SUCCESS",
		"message": print,
	})
}

func (t *TicketController) listTicketSortedBy(c *gin.Context) {
	var userRole model.UserRole
	err := t.ParseRequestBody(c, &userRole)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"status":  "FAILED",
			"message": err,
		})
		return
	}

	orderBy := c.Param("orderBy")
	field := c.Param("field")
	print, err := t.ucTicket.ListTicketSortedBy(&userRole, orderBy, field)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
				"status":  "FAILED",
				"message": "Error ticket not found",
			})
			return
		}
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"status":  "FAILED",
			"message": err,
		})
		return

	}
	c.JSON(http.StatusOK, gin.H{
		"status":  "SUCCESS",
		"message": print,
	})
}

// func (t *TicketController) getTicketSummaryByDate(c *gin.Context) { // BELUM BISA
// 	statusId := c.Param("statusid")
// 	convStatusId, err := strconv.Atoi(statusId)
// 	if err != nil {
// 		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
// 			"status":  "FAILED",
// 			"message": "Error converting status id",
// 		})
// 		return
// 	}

// 	param1 := c.Param("startdate") + " 23:59:59 +0700 WIB"
// 	param2 := c.Param("enddate") + " 23:59:59 +0700 WIB"
// 	// starDate, err := time.Parse("2006-01-02 15:04:05 -0700 MST", param1)
// 	// if err != nil {
// 	// 	c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
// 	// 		"status":  "FAILED",
// 	// 		"message": "Error converting date",
// 	// 	})
// 	// 	return
// 	// }

// 	// endDate, err := time.Parse("2006-01-02 15:04:05 -0700 MST", param2)
// 	// if err != nil {
// 	// 	c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
// 	// 		"status":  "FAILED",
// 	// 		"message": "Error converting date",
// 	// 	})
// 	// 	return
// 	// }

// 	var userRole model.UserRole
// 	err = t.ParseRequestBody(c, &userRole)
// 	if err != nil {
// 		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
// 			"status":  "FAILED",
// 			"message": err,
// 		})
// 		return
// 	}

// 	print, err := t.ucTicket.GetSummaryTicketByDate(convStatusId, &userRole, param1, param2)
// 	if err != nil {
// 		if errors.Is(err, gorm.ErrRecordNotFound) {
// 			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
// 				"status":  "FAILED",
// 				"message": "Error ticket not found",
// 			})
// 			return
// 		}
// 		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
// 			"status":  "FAILED",
// 			"message": err,
// 		})
// 		return

// 	}
// 	c.JSON(http.StatusOK, gin.H{
// 		"status":  "SUCCESS",
// 		"message": print,
// 	})
// }

func (t *TicketController) updatePIC(c *gin.Context) {
	var updateInput struct {
		TicketId string `json:"ticketId"`
		PicId    int    `json:"picId"`
	}

	if err := c.BindJSON(&updateInput); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "BAD REQUEST",
			"message": err.Error(),
		})
	} else {
		fmt.Println(updateInput)
		err := t.ucTicket.UpdatePIC(updateInput.TicketId, updateInput.PicId)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
				"status":  "FAILED",
				"message": "Error when assign new PIC",
			})
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"status":  "SUCCESS",
			"message": "complain assigned to a new pic",
		})
	}
}

func (t *TicketController) updateStatus(c *gin.Context) {
	var updateInput struct {
		TicketId string `json:"ticketId"`
		StatusId int    `json:"statusId"`
	}

	if err := c.BindJSON(&updateInput); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "BAD REQUEST",
			"message": err.Error(),
		})
	} else {
		fmt.Println(updateInput)
		err := t.ucTicket.UpdateStatus(updateInput.TicketId, updateInput.StatusId)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
				"status":  "FAILED",
				"message": "Error when update status ticket",
			})
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"status":  "SUCCESS",
			"message": "updated status success",
		})
	}
}

func NewTicketController(router *gin.Engine, ucTicket usecase.TicketUseCase) *TicketController {
	controller := TicketController{
		router:   router,
		ucTicket: ucTicket,
	}

	// router.PUT("/updatePIC/:userId", controller.updatePIC)

	protectedGroup := router.Group("api/ticket")

	protectedGroup.POST("", controller.createTicket)
	protectedGroup.GET("/list", controller.ListTicketByUserId)
	protectedGroup.GET("/department/:id", controller.ListTicketByDepartmentId)

	protectedGroup.GET("/:id", controller.getTicketById)
	protectedGroup.POST("/sort", controller.listTicketSortedBy)
	protectedGroup.GET("/summary/:statusid", controller.getTicketSummary)
	// protectedGroup.GET("/summarydate/:statusid/:startdate/:enddate", controller.getTicketSummaryByDate)

	protectedGroup.PUT("/update-pic", controller.updatePIC)
	protectedGroup.PUT("/update-status", controller.updateStatus)
	protectedGroup.POST("/listp/:orderBy/:field", controller.listTicketSortedBy)

	return &controller
}
