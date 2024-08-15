package define

// 返回的实时数据类型
type LiveData[T ElectFire | DataInfo | RemoteData] struct {
	Id   int64 `json:"id"`
	Data T     `json:"data"`
}

// 设备的实时数据信息
type DataInfo struct {
	Temperature            float32 `json:"temperature"`
	TemperatureAlarmStatus uint16  `json:"temperatureAlarmStatus"`
	DIStatus               uint16  `json:"DIStatus"`
	SOEpointer             uint32  `json:"SOEpointer"`
	VoltageImbalance       float32 `json:"voltageImbalance"`
	CurrentImbalance       float32 `json:"currentImbalance"`
}

// 历史数据信息
type HistoryData struct {
	DateTime string             `json:"datetime"`
	Data     LiveData[DataInfo] `json:"data"`
}

// 设备的电气火灾数据
type ElectFire struct {
	Threshold  uint16 `json:"threshold"`
	AlarmSound uint16 `json:"alarmSound"`
	FaultSound uint16 `json:"faultSound"`
}

// 远程遥控数据
type RemoteData struct {
	InputType uint16 `json:"inputType"`
	Type      string `json:"type"`
	IsSuccess bool   `json:"isSuccess"`
}

// 远程遥控参数
var OPEN uint16 = 0xFF00
var CLOSE uint16 = 0x0000

const (
	RemoteClosingPreset    int = 1 //遥合预置
	RemoteClosingExecution int = 2 //遥合执行
	RemoteDvisionPreset    int = 3 //遥分预置
	RemoteDvisionExecution int = 4 // 遥分执行
)
const (
	Type1 = "遥合预置"
	Type2 = "遥合执行"
	Type3 = "遥分预置"
	Type4 = "遥分执行"
)
