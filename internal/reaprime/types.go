package reaprime

import "time"

type MachineSnapshot struct {
	Timestamp              time.Time
	State                  string
	Substate               string
	Flow                   float64
	Pressure               float64
	TargetFlow             float64
	TargetPressure         float64
	MixTemperature         float64
	GroupTemperature       float64
	TargetMixTemperature   float64
	TargetGroupTemperature float64
	ProfileFrame           float64
	SteamTemperature       float64
}

type ScaleSnapshot struct {
	Timestamp  time.Time
	Status     string
	Weight     float64
	WeightFlow float64
	Battery    *float64
	TimerMS    *float64
}

type ShotSettings struct {
	SteamSetting          float64
	TargetSteamTemp       float64
	TargetSteamDuration   float64
	TargetHotWaterTemp    float64
	TargetHotWaterVolume  float64
	TargetHotWaterSeconds float64
	TargetShotVolume      float64
	GroupTemp             float64
}

type WaterLevels struct {
	CurrentLevel float64
	RefillLevel  float64
}

type DisplayState struct {
	WakeLockEnabled            bool
	WakeLockOverride           bool
	Brightness                 float64
	RequestedBrightness        float64
	LowBatteryBrightnessActive bool
	BrightnessSupported        bool
	WakeLockSupported          bool
}

type DeviceState struct {
	Type      string
	State     string
	Available bool
}

type DevicesSnapshot struct {
	Timestamp time.Time
	Scanning  bool
	Phase     string
	Devices   []DeviceState
	ErrorKind string
}

type Snapshot struct {
	Now time.Time

	Streams map[string]StreamState

	Machine     *MachineSnapshot
	Scale       *ScaleSnapshot
	Shot        *ShotSettings
	Water       *WaterLevels
	Display     *DisplayState
	Devices     *DevicesSnapshot
	Transitions map[string]uint64
}

type StreamState struct {
	Connected       bool
	LastMessageTime time.Time
	Messages        uint64
	Reconnects      uint64
	Errors          uint64
}
