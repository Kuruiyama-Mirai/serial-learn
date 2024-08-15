package internal

import (
	"Demo/App/define"
	"Demo/App/models"
	"Demo/App/tools"
	"log"
	"math"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

// 获取实时数据
func ShowRealTimeData(c *gin.Context) {
	// 在所有表中做一个查询id,匹配到了才进行查询，不然就报错，这样是多350操作
	// 每一个client只有读写接口了

	id, _ := strconv.Atoi(c.Query("id"))
	typeStr := c.Query("type")

	var newLiveData define.LiveData[define.DataInfo]
	var err error
	switch typeStr {
	case "350":
		newLiveData, err = GetPM350Data(id)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"code":  400,
				"error": err,
			})
			return
		}
	case "53A":
		newLiveData, err = Get53AData(id)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"code":  400,
				"error": err,
			})
			return
		}
	default:
		c.JSON(http.StatusBadRequest, gin.H{
			"code":  400,
			"error": "请输入要查询的设备类型",
		})
		return
	}
	newInsert := &models.SaveData{
		DeviceID:               int64(newLiveData.Id),
		Temperature:            float64(newLiveData.Data.Temperature),
		TemperatureAlarmStatus: int(newLiveData.Data.TemperatureAlarmStatus),
		DIStatus:               int(newLiveData.Data.DIStatus),
		SOEpointer:             int(newLiveData.Data.SOEpointer),
		VoltageImbalance:       float64(newLiveData.Data.VoltageImbalance),
		CurrentImbalance:       float64(newLiveData.Data.CurrentImbalance),
		CreatedAt:              time.Now(),
	}
	// 2.保存到数据库中
	res := models.DB.Create(&newInsert)
	if res.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":  400,
			"error": res.Error,
		})
		return
	} else {
		log.Printf("插入数据成功!")
	}
	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"info": newLiveData,
	})

}

// 根据表的Id来获取350的数据
func GetPM350Data(id int) (define.LiveData[define.DataInfo], error) {

	client, err := tools.CreateModbusClient(id)
	if err != nil {
		return define.LiveData[define.DataInfo]{}, err
	}

	//2.读取寄存器
	temperatureRes, err := client.ReadHoldingRegisters(uint16(86), uint16(2))
	if err != nil {
		log.Println("读取温度失败", err)
	}
	temperatureAlarmStatusRes, err := client.ReadHoldingRegisters(uint16(95), uint16(1))
	if err != nil {
		log.Println("读取温度报警状态失败", err)
	}
	dIStatusRes, err := client.ReadHoldingRegisters(uint16(96), uint16((1)))
	if err != nil {
		log.Println("读取DI状态失败", err)
	}
	SOEpointerRes, err := client.ReadHoldingRegisters(uint16(102), uint16(2))
	if err != nil {
		log.Println("读取SOE指针总数失败", err)
	}
	voltageImbalanceRes, err := client.ReadHoldingRegisters(uint16(1330), uint16(2))
	if err != nil {
		log.Println("读取电压不平衡度失败", err)
	}
	currentImbalanceRes, err := client.ReadHoldingRegisters(uint16(1332), uint16(2))
	if err != nil {
		log.Println("读取电流不平衡度失败", err)
	}
	// serialNoRes, err := client.ReadHoldingRegisters(uint16(9825), uint16(2))
	// if err != nil {
	// 	log.Println("读取设备序列号失败", err)
	// }
	temperature := tools.AnalyzeRegistersToFloat32(temperatureRes)
	temperatureAlarmStatus := tools.AnalyzeRegistersToUint16AndInt16[uint16](temperatureAlarmStatusRes)
	dIStatus := tools.AnalyzeRegistersToUint16AndInt16[uint16](dIStatusRes)
	SOEpointer := tools.AnalyzeRegistersToUint32AndInt32[uint32](SOEpointerRes)
	voltageImbalance := tools.AnalyzeRegistersToFloat32(voltageImbalanceRes)
	currentImbalance := tools.AnalyzeRegistersToFloat32(currentImbalanceRes)
	//serialNo := tools.AnalyzeRegistersToUint32AndInt32[uint32](serialNoRes)

	//定义返回的数据类型
	newLiveData := define.LiveData[define.DataInfo]{
		Id: int64(id),
		Data: define.DataInfo{
			Temperature:            temperature,
			TemperatureAlarmStatus: temperatureAlarmStatus,
			DIStatus:               dIStatus,
			SOEpointer:             SOEpointer,
			VoltageImbalance:       voltageImbalance,
			CurrentImbalance:       currentImbalance,
		},
	}

	return newLiveData, nil
}

// 获取53A的数据
func Get53AData(id int) (define.LiveData[define.DataInfo], error) {

	client, err := tools.CreateModbusClient(id)
	if err != nil {
		return define.LiveData[define.DataInfo]{}, err
	}

	//2.读取寄存器
	temperatureRes, err := client.ReadHoldingRegisters(uint16(86), uint16(2))
	if err != nil {
		log.Println("读取温度失败", err)
	}
	// temperatureAlarmStatusRes, err := client.ReadHoldingRegisters(uint16(95), uint16(1))
	// if err != nil {
	// 	log.Println("读取温度报警状态失败", err)
	// }
	dIStatusRes, err := client.ReadHoldingRegisters(uint16(96), uint16((1)))
	if err != nil {
		log.Println("读取DI状态失败", err)
	}
	SOEpointerRes, err := client.ReadHoldingRegisters(uint16(102), uint16(2))
	if err != nil {
		log.Println("读取SOE指针总数失败", err)
	}
	voltageImbalanceRes, err := client.ReadHoldingRegisters(uint16(1330), uint16(2))
	if err != nil {
		log.Println("读取电压不平衡度失败", err)
	}
	currentImbalanceRes, err := client.ReadHoldingRegisters(uint16(1332), uint16(2))
	if err != nil {
		log.Println("读取电流不平衡度失败", err)
	}
	// serialNoRes, err := client.ReadHoldingRegisters(uint16(9825), uint16(2))
	// if err != nil {
	// 	log.Println("读取设备序列号失败", err)
	// }
	temperature := tools.AnalyzeRegistersToFloat32(temperatureRes)
	// 53A
	//temperatureAlarmStatus := tools.AnalyzeRegistersToUint16AndInt16[uint16](temperatureAlarmStatusRes)
	dIStatus := tools.AnalyzeRegistersToUint16AndInt16[uint16](dIStatusRes)
	SOEpointer := tools.AnalyzeRegistersToUint32AndInt32[uint32](SOEpointerRes)
	voltageImbalance := tools.AnalyzeRegistersToFloat32(voltageImbalanceRes)
	currentImbalance := tools.AnalyzeRegistersToFloat32(currentImbalanceRes)
	//serialNo := tools.AnalyzeRegistersToUint32AndInt32[uint32](serialNoRes)

	//定义返回的数据类型
	newLiveData := define.LiveData[define.DataInfo]{
		Id: int64(id),
		Data: define.DataInfo{
			Temperature:            temperature,
			TemperatureAlarmStatus: math.MaxUint16,
			DIStatus:               dIStatus,
			SOEpointer:             SOEpointer,
			VoltageImbalance:       voltageImbalance,
			CurrentImbalance:       currentImbalance,
		},
	}

	return newLiveData, nil

}
