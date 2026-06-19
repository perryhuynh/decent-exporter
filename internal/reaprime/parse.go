package reaprime

import (
	"encoding/json"
	"strconv"
	"time"
)

func ParseMachine(data []byte) (MachineSnapshot, error) {
	var raw map[string]any
	if err := json.Unmarshal(data, &raw); err != nil {
		return MachineSnapshot{}, err
	}
	state := asMap(raw["state"])
	return MachineSnapshot{
		Timestamp:              asTime(raw["timestamp"]),
		State:                  asString(state["state"]),
		Substate:               asString(state["substate"]),
		Flow:                   asFloat(raw["flow"]),
		Pressure:               asFloat(raw["pressure"]),
		TargetFlow:             asFloat(raw["targetFlow"]),
		TargetPressure:         asFloat(raw["targetPressure"]),
		MixTemperature:         asFloat(raw["mixTemperature"]),
		GroupTemperature:       asFloat(raw["groupTemperature"]),
		TargetMixTemperature:   asFloat(raw["targetMixTemperature"]),
		TargetGroupTemperature: asFloat(raw["targetGroupTemperature"]),
		ProfileFrame:           asFloat(raw["profileFrame"]),
		SteamTemperature:       asFloat(raw["steamTemperature"]),
	}, nil
}

func ParseScale(data []byte) (ScaleSnapshot, error) {
	var raw map[string]any
	if err := json.Unmarshal(data, &raw); err != nil {
		return ScaleSnapshot{}, err
	}
	scale := ScaleSnapshot{
		Timestamp:  asTime(raw["timestamp"]),
		Status:     asString(raw["status"]),
		Weight:     asFloat(raw["weight"]),
		WeightFlow: asFloat(raw["weightFlow"]),
	}
	if battery, ok := maybeFloat(raw["battery"]); ok {
		scale.Battery = &battery
	}
	if timer, ok := maybeFloat(raw["timerValue"]); ok {
		scale.TimerMS = &timer
	}
	return scale, nil
}

func ParseShotSettings(data []byte) (ShotSettings, error) {
	var raw map[string]any
	if err := json.Unmarshal(data, &raw); err != nil {
		return ShotSettings{}, err
	}
	return ShotSettings{
		SteamSetting:          asFloat(raw["steamSetting"]),
		TargetSteamTemp:       asFloat(raw["targetSteamTemp"]),
		TargetSteamDuration:   asFloat(raw["targetSteamDuration"]),
		TargetHotWaterTemp:    asFloat(raw["targetHotWaterTemp"]),
		TargetHotWaterVolume:  asFloat(raw["targetHotWaterVolume"]),
		TargetHotWaterSeconds: asFloat(raw["targetHotWaterDuration"]),
		TargetShotVolume:      asFloat(raw["targetShotVolume"]),
		GroupTemp:             asFloat(raw["groupTemp"]),
	}, nil
}

func ParseWater(data []byte) (WaterLevels, error) {
	var raw map[string]any
	if err := json.Unmarshal(data, &raw); err != nil {
		return WaterLevels{}, err
	}
	return WaterLevels{
		CurrentLevel: asFloat(raw["currentLevel"]),
		RefillLevel:  asFloat(raw["refillLevel"]),
	}, nil
}

func ParseDisplay(data []byte) (DisplayState, error) {
	var raw map[string]any
	if err := json.Unmarshal(data, &raw); err != nil {
		return DisplayState{}, err
	}
	supported := asMap(raw["platformSupported"])
	return DisplayState{
		WakeLockEnabled:            asBool(raw["wakeLockEnabled"]),
		WakeLockOverride:           asBool(raw["wakeLockOverride"]),
		Brightness:                 asFloat(raw["brightness"]),
		RequestedBrightness:        asFloat(raw["requestedBrightness"]),
		LowBatteryBrightnessActive: asBool(raw["lowBatteryBrightnessActive"]),
		BrightnessSupported:        asBool(supported["brightness"]),
		WakeLockSupported:          asBool(supported["wakeLock"]),
	}, nil
}

func ParseDevices(data []byte) (DevicesSnapshot, error) {
	var raw map[string]any
	if err := json.Unmarshal(data, &raw); err != nil {
		return DevicesSnapshot{}, err
	}
	out := DevicesSnapshot{
		Timestamp: asTime(raw["timestamp"]),
		Scanning:  asBool(raw["scanning"]),
	}
	status := asMap(raw["connectionStatus"])
	out.Phase = asString(status["phase"])
	if errObj := asMap(status["error"]); len(errObj) > 0 {
		out.ErrorKind = asString(errObj["kind"])
	}
	for _, item := range asSlice(raw["devices"]) {
		device := asMap(item)
		out.Devices = append(out.Devices, DeviceState{
			Type:      asString(device["type"]),
			State:     asString(device["state"]),
			Available: asBool(device["available"]),
		})
	}
	return out, nil
}

func asMap(v any) map[string]any {
	if m, ok := v.(map[string]any); ok {
		return m
	}
	return nil
}

func asSlice(v any) []any {
	if s, ok := v.([]any); ok {
		return s
	}
	return nil
}

func asString(v any) string {
	if s, ok := v.(string); ok {
		return s
	}
	return ""
}

func asBool(v any) bool {
	if b, ok := v.(bool); ok {
		return b
	}
	return false
}

func asFloat(v any) float64 {
	f, _ := maybeFloat(v)
	return f
}

func maybeFloat(v any) (float64, bool) {
	switch typed := v.(type) {
	case float64:
		return typed, true
	case int:
		return float64(typed), true
	case json.Number:
		f, err := typed.Float64()
		return f, err == nil
	case string:
		f, err := strconv.ParseFloat(typed, 64)
		return f, err == nil
	default:
		return 0, false
	}
}

func asTime(v any) time.Time {
	s := asString(v)
	if s == "" {
		return time.Time{}
	}
	for _, layout := range []string{time.RFC3339Nano, "2006-01-02T15:04:05.999999"} {
		t, err := time.Parse(layout, s)
		if err == nil {
			return t
		}
	}
	return time.Time{}
}
