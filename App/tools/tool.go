package tools

import (
	"encoding/binary"
	"errors"
	"log"
	"math"
	"os"
	"time"

	"github.com/goburrow/modbus"
	"gopkg.in/yaml.v3"
)

// 连接串口 调整架构，这里只返回一组handler,具体的modbus客户端交给接口来建立
func SerialConfig() ([]*modbus.RTUClientHandler, error) {
	// 构建配置文件的路径

	file, err := os.Open("config.yaml")
	if err != nil {
		log.Fatalf("error: %v", err)
	}
	defer file.Close()

	var devices Devices
	decoder := yaml.NewDecoder(file)
	if err := decoder.Decode(&devices); err != nil {
		log.Fatalf("error: %v", err)
	}

	// 创建ModbusRTU客户端
	handler := make([]*modbus.RTUClientHandler, len(devices.Devices))
	//创建modbus客户端
	//client := make([]modbus.Client, len(devices.Devices))

	for i, device := range devices.Devices {

		handler[i] = modbus.NewRTUClientHandler(device.Address)
		handler[i].SlaveId = byte(device.SlaveId) //这块是通讯地址，也即每个表都有对应的ID，地址读对了才能正确的通信
		handler[i].BaudRate = device.BaudRate
		handler[i].DataBits = device.DataBits
		handler[i].Parity = device.Parity //默认是偶校验
		handler[i].StopBits = device.StopBits
		handler[i].Timeout = time.Second * 10

		err := handler[i].Connect()
		if err != nil {
			log.Println("modbus连接失败", err)
		}
		log.Printf("modbus%d连接成功！！！", i)
		defer handler[i].Close()

		// client[i] = modbus.NewClient(handler[i])
	}

	return handler, nil
}

// 根据不同的通讯ID创建Modbus客户端
func CreateModbusClient(id int) (client modbus.Client, err error) {
	handler, err := SerialConfig()
	if err != nil {
		log.Println("串口连接失败", err)
		return nil, err
	}
	// 需要一个计数器来判断如果没有一个表匹配的情况
	count := 0

	for i := range handler {
		if handler[i].SlaveId == byte(id) {
			client = modbus.NewClient(handler[i])
			log.Printf("modbus:%d号客户端成功建立", id)
		} else {
			count++
		}
	}
	if count == len(handler) {
		return nil, errors.New("未找到匹配的设备")
	}
	return
}

// 寄存器值解析
// 转成int16和uint16
func AnalyzeRegistersToUint16AndInt16[T uint16 | int16](res []byte) T {
	return T(binary.BigEndian.Uint16(res))
}

// 转成int32和uint32
func AnalyzeRegistersToUint32AndInt32[T uint32 | int32](res []byte) T {
	return T(binary.BigEndian.Uint32(res))
}

// 转成float32
func AnalyzeRegistersToFloat32(res []byte) float32 {
	bits := binary.BigEndian.Uint32(res)
	value := math.Float32frombits(bits)
	return value
}

// 转成char
func AnalyzeRegistersToChar(res []byte) string {
	var charRes []rune
	for _, b := range res {
		charRes = append(charRes, rune(b))
	}
	return string(charRes)
}

// 写入寄存器的int转字节序列
func IntToModbusBytes[T int | int16 | int32 | uint16 | uint | uint32](value T) []byte {
	bytes := make([]byte, 2)
	binary.BigEndian.PutUint16(bytes, uint16(value))
	return bytes
}

type Device struct {
	Address  string `yaml:"address"`
	SlaveId  int    `yaml:"slaveId"`
	BaudRate int    `yaml:"baudRate"`
	DataBits int    `yaml:"dataBits"`
	Parity   string `yaml:"parity"`
	StopBits int    `yaml:"stopBits"`
}

// Devices holds a list of Device.
type Devices struct {
	Devices []Device `yaml:"Devices"`
}
