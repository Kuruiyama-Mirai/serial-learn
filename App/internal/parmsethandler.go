package internal

import (
	"Demo/App/define"
	"Demo/App/tools"
	"errors"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

// 设置参数结构体
type Parm struct {
	ID         uint `json:"id" binding:"required"`
	Threshold  int  `json:"threshold"`
	AlarmSound int  `json:"alarmSound"`
	FaultSound int  `json:"faultSound"`
}

// 参数设置
func ParmSet(c *gin.Context) {
	//1.连接到串口

	var parm Parm
	// 解析并绑定JSON格式的请求数据到结构体
	if err := c.ShouldBindJSON(&parm); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":       400,
			"json error": err.Error(),
		})
		return
	}
	newFireData, err := Set350Parm(parm)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":  400,
			"error": err,
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "参数设置成功!",
		"info":    newFireData,
	})
}

// 根据传递的json参数来设置350
func Set350Parm(parm Parm) (define.LiveData[define.ElectFire], error) {

	client, err := tools.CreateModbusClient(int(parm.ID))
	if err != nil {
		return define.LiveData[define.ElectFire]{}, err
	}

	//先判断有没有参数，写入threshold
	if uint16(parm.Threshold) <= 45 {
		res := tools.IntToModbusBytes(parm.Threshold)
		_, err = client.WriteMultipleRegisters(uint16(50004), uint16(1), res)
		if err != nil {
			return define.LiveData[define.ElectFire]{}, err
		} else {
			log.Println("成功写入Threshold")
		}
	} else {
		return define.LiveData[define.ElectFire]{}, errors.New("threshold参数设置超过上限")
	}
	//写入AlarmSound
	if uint16(parm.AlarmSound) <= 2 {
		res := tools.IntToModbusBytes(parm.AlarmSound)
		_, err = client.WriteMultipleRegisters(uint16(50000), uint16(1), res)
		if err != nil {
			return define.LiveData[define.ElectFire]{}, err
		} else {
			log.Println("成功写入AlarmSound")
		}
	} else {
		return define.LiveData[define.ElectFire]{}, errors.New("alarmSound参数设置超过上限")
	}
	//写入faultSound
	if uint16(parm.FaultSound) <= 2 {
		res := tools.IntToModbusBytes(parm.FaultSound)
		_, err = client.WriteMultipleRegisters(uint16(50001), uint16(1), res)
		if err != nil {
			return define.LiveData[define.ElectFire]{}, err
		} else {
			log.Println("成功写入FaultSound")
		}
	} else {
		return define.LiveData[define.ElectFire]{}, errors.New("faultSound参数设置超过上限")
	}

	//再读一次数据做判断
	newAlarmSoundRes, err := client.ReadHoldingRegisters(uint16(50000), uint16(1))
	if err != nil {
		log.Println("读取报警声失败", err)
	}
	newFaultSoundRes, err := client.ReadHoldingRegisters(uint16(50001), uint16(1))
	if err != nil {
		log.Println("读取故障声失败", err)
	}
	newThresholdRes, err := client.ReadHoldingRegisters(uint16(50004), uint16(1))
	if err != nil {
		log.Println("读取温度门槛值失败", err)
	}

	newThreshold := tools.AnalyzeRegistersToUint16AndInt16[uint16](newThresholdRes)
	newAlarmSound := tools.AnalyzeRegistersToUint16AndInt16[uint16](newAlarmSoundRes)
	newFaultSound := tools.AnalyzeRegistersToUint16AndInt16[uint16](newFaultSoundRes)

	if (newThreshold == uint16(parm.Threshold)) && (newAlarmSound == uint16(parm.AlarmSound)) && (newFaultSound == uint16(parm.FaultSound)) {
		log.Println("全部写入成功")
	} else {
		log.Println("部分写入失败")
	}

	//给前端展示的是 写入后再次查询的数据
	newFireData := define.LiveData[define.ElectFire]{
		Id: int64(parm.ID),
		Data: define.ElectFire{
			Threshold:  uint16(newThreshold),
			AlarmSound: uint16(newAlarmSound),
			FaultSound: uint16(newFaultSound),
		},
	}

	return newFireData, nil
}
