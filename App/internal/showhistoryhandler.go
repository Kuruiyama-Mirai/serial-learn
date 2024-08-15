package internal

import (
	"Demo/App/define"
	"Demo/App/models"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// 展示历史数据
func ShowHistoryData(c *gin.Context) {
	startTimeStr := c.Query("start")
	endTimeStr := c.Query("end")
	page, err := strconv.Atoi(c.Query("page"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":   400,
			"分页输入错误": err,
		})
	}
	pageSize, err := strconv.Atoi(c.Query("pagesize"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":   400,
			"分页输入错误": err,
		})
	}

	res := make([]define.HistoryData, 0)

	newData, total, err := FindByDate(models.DB, startTimeStr, endTimeStr, page, pageSize)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":  400,
			"error": err,
		})
	}
	for i := range newData {
		history := define.HistoryData{
			DateTime: newData[i].CreatedAt.Format("2006-01-02 15:04:05"),
			Data: define.LiveData[define.DataInfo]{
				Id: int64(newData[i].DeviceID),
				Data: define.DataInfo{
					Temperature:            float32(newData[i].Temperature),
					TemperatureAlarmStatus: uint16(newData[i].TemperatureAlarmStatus),
					DIStatus:               uint16(newData[i].DIStatus),
					SOEpointer:             uint32(newData[i].SOEpointer),
					VoltageImbalance:       float32(newData[i].VoltageImbalance),
					CurrentImbalance:       float32(newData[i].CurrentImbalance),
				},
			},
		}
		res = append(res, history)
	}
	if res != nil {
		c.JSON(http.StatusOK, gin.H{
			"code":  200,
			"total": total,
			"info":  res,
		})
	} else {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":  400,
			"error": "数据查询失败,请检查日期是否正确",
		})
	}
}

// 分页查询函数
func FindByDate(db *gorm.DB, start, end string, page, pageSize int) (totladata []models.SaveData, num int64, err error) {
	var data []models.SaveData
	offset := (page - 1) * pageSize
	layout := "2006-01-02 15:04:05"
	startTime, err := time.Parse(layout, start)
	if err != nil {
		log.Println("Invalid start time")
		return nil, 0, err
	}

	endTime, err := time.Parse(layout, end)
	if err != nil {
		log.Println("Invalid end time")
		return nil, 0, err
	}
	var total int64
	db.Model(&data).Where("created_at BETWEEN ? AND ?", startTime.Format(layout), endTime.Format(layout)).Count(&total)
	res := db.Where("created_at BETWEEN ? AND ?", startTime.Format(layout), endTime.Format(layout)).Offset(offset).Limit(pageSize).Find(&data)

	if res.Error != nil {
		log.Println(res.Error)
	}
	return data, total, nil
}
