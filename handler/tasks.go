package handler

import (
	"github.com/DenrianWeiss/catballCoin/model"
	"github.com/DenrianWeiss/catballCoin/service"
	"github.com/gin-gonic/gin"
	"net/http"
)

func GetTaskCount(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, gin.H{
		"task_count": service.PeekTaskId(),
	})
}

func AddTask(ctx *gin.Context) {
	key, _ := ctx.Params.Get("key")
	if key != service.GlobalConfig.RpcKey {
		ctx.JSON(http.StatusForbidden, gin.H{
			"error": "key",
		})
		return
	}
	taskInfo := &model.CoinTask{}
	err := ctx.ShouldBindJSON(taskInfo)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "request",
		})
		return
	}
	err = service.AddTaskToDatabase(taskInfo)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": "database",
		})
		return
	}
	taskInfo.TaskID = service.GetTaskId()
	go service.NewCoinHarvestTask(taskInfo)
	ctx.JSON(http.StatusOK, gin.H{
		"error": "",
		"id": taskInfo.TaskID,
	})
}
