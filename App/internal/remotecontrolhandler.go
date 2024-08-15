package internal

import (
	"Demo/App/define"
	"Demo/App/tools"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

// 设置参数结构体
type Contorl struct {
	ID        uint32 `json:"id" binding:"required"`
	InputType *int   `json:"inputType" binding:"required"`
}

// 远程遥控
func RemoteControl(c *gin.Context) {
	//1.连接到串口
	var cont Contorl
	var resType string
	// 解析并绑定JSON格式的请求数据到结构体
	if err := c.ShouldBindJSON(&cont); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":       400,
			"json error": err.Error(),
		})
		return
	}
	client, err := tools.CreateModbusClient(int(cont.ID))
	if err != nil {
		log.Println("串口连接失败", err)
		return
	}
	//读取DO1状态
	DOStatus, err := client.ReadHoldingRegisters(uint16(98), uint16(1))
	if err != nil {
		log.Println("读取设备序列号失败", err)
	}
	do := tools.AnalyzeRegistersToUint16AndInt16[uint16](DOStatus)
	//do1=0 就是bit1的遥合的状态 1就是遥分的状态
	do1 := (do >> 0) & 1

	if cont.InputType != nil {

		switch *cont.InputType {
		case define.RemoteClosingPreset:
			if do1 == 0 {
				resType = define.Type1
				_, err := client.WriteSingleCoil(uint16(60064), define.OPEN)
				if err != nil {
					c.JSON(http.StatusBadRequest, gin.H{
						"code":      400,
						"遥合预置失败,%v": err.Error(),
					})
					return
				}

			} else {
				c.JSON(http.StatusOK, gin.H{
					"code":    666,
					"message": "设备已经打开，无需再次打开",
				})
				return
			}
		case define.RemoteClosingExecution:
			resType = define.Type2
			_, err = client.WriteSingleCoil(uint16(60065), define.OPEN)
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{
					"code":      400,
					"遥合执行失败,%v": err.Error(),
				})
				return
			}

		case define.RemoteDvisionPreset:
			if do1 == 1 {
				resType = define.Type3
				_, err = client.WriteSingleCoil(uint16(60066), define.OPEN)
				if err != nil {
					c.JSON(http.StatusBadRequest, gin.H{
						"code":      400,
						"遥分预置失败,%v": err.Error(),
					})
					return
				}

			} else {
				c.JSON(http.StatusOK, gin.H{
					"code":    666,
					"message": "设备已经关闭，无需再次关闭",
				})
				return
			}
		case define.RemoteDvisionExecution:

			resType = define.Type4
			_, err = client.WriteSingleCoil(uint16(60067), define.OPEN)
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{
					"code":      400,
					"遥分执行失败,%v": err.Error(),
				})
				return
			}

		default:
			return
		}
	} else {
		c.JSON(http.StatusBadGateway, gin.H{
			"code":    500,
			"message": "请输入要遥控的类型",
		})
		return
	}
	newRemoteData := define.LiveData[define.RemoteData]{
		Id: int64(cont.ID),
		Data: define.RemoteData{
			InputType: uint16(*cont.InputType),
			Type:      resType,
			IsSuccess: true,
		},
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "远程遥控成功!",
		"info":    newRemoteData,
	})
}
