package exporter

import (
	"strconv"

	"github.com/perryhuynh/decent-exporter/internal/decent"
	"github.com/prometheus/client_golang/prometheus"
)

type BuildInfo struct {
	Version string
	Commit  string
	Date    string
}

type Collector struct {
	store *decent.Store
	build BuildInfo

	buildInfo           *prometheus.Desc
	streamConnected     *prometheus.Desc
	streamLastMessage   *prometheus.Desc
	streamMessages      *prometheus.Desc
	streamReconnects    *prometheus.Desc
	streamErrors        *prometheus.Desc
	machineState        *prometheus.Desc
	machineTemperature  *prometheus.Desc
	machineProfileFrame *prometheus.Desc
	shotsTotal          *prometheus.Desc
	shotDuration        *prometheus.Desc
	shotPeakPressure    *prometheus.Desc
	shotPeakFlow        *prometheus.Desc
	shotAverageFlow     *prometheus.Desc
	shotSetting         *prometheus.Desc
	waterLevel          *prometheus.Desc
	displayBool         *prometheus.Desc
	displayBrightness   *prometheus.Desc
	deviceCount         *prometheus.Desc
	devicesScanning     *prometheus.Desc
	machineTransitions  *prometheus.Desc
}

func NewCollector(store *decent.Store, build BuildInfo) *Collector {
	labels := prometheus.Labels{}
	return &Collector{
		store: store,
		build: build,

		buildInfo:           prometheus.NewDesc("decent_exporter_build_info", "Build information for decent-exporter.", []string{"version", "commit", "date"}, labels),
		streamConnected:     prometheus.NewDesc("decent_stream_connected", "Whether a Decent WebSocket stream is currently connected.", []string{"stream"}, labels),
		streamLastMessage:   prometheus.NewDesc("decent_stream_last_message_timestamp_seconds", "Unix timestamp of the last message received from a Decent stream.", []string{"stream"}, labels),
		streamMessages:      prometheus.NewDesc("decent_stream_messages_total", "Total messages received from a Decent stream since exporter start.", []string{"stream"}, labels),
		streamReconnects:    prometheus.NewDesc("decent_stream_reconnects_total", "Total successful Decent stream connections since exporter start.", []string{"stream"}, labels),
		streamErrors:        prometheus.NewDesc("decent_stream_errors_total", "Total Decent stream errors since exporter start.", []string{"stream"}, labels),
		machineState:        prometheus.NewDesc("decent_machine_state", "Current machine state as a one-hot gauge.", []string{"state", "substate"}, labels),
		machineTemperature:  prometheus.NewDesc("decent_machine_temperature_celsius", "Machine temperatures in Celsius.", []string{"sensor"}, labels),
		machineProfileFrame: prometheus.NewDesc("decent_machine_profile_frame", "Current machine profile frame.", nil, labels),
		shotsTotal:          prometheus.NewDesc("decent_shots_total", "Completed espresso shots observed since exporter start.", nil, labels),
		shotDuration:        prometheus.NewDesc("decent_shot_duration_seconds", "Duration of the last completed espresso shot in seconds.", nil, labels),
		shotPeakPressure:    prometheus.NewDesc("decent_shot_peak_pressure_bar", "Peak pressure during the last completed espresso shot in bar.", nil, labels),
		shotPeakFlow:        prometheus.NewDesc("decent_shot_peak_flow_ml_per_second", "Peak machine flow during the last completed espresso shot in ml/s.", nil, labels),
		shotAverageFlow:     prometheus.NewDesc("decent_shot_average_flow_ml_per_second", "Mean machine flow during the last completed espresso shot in ml/s.", nil, labels),
		shotSetting:         prometheus.NewDesc("decent_shot_setting", "Current shot setting values.", []string{"setting"}, labels),
		waterLevel:          prometheus.NewDesc("decent_water_level_millimeters", "Water tank level values in millimeters.", []string{"kind"}, labels),
		displayBool:         prometheus.NewDesc("decent_display_state", "Display boolean state as one-hot gauges.", []string{"state"}, labels),
		displayBrightness:   prometheus.NewDesc("decent_display_brightness_percent", "Display brightness percentage.", []string{"kind"}, labels),
		deviceCount:         prometheus.NewDesc("decent_devices", "Number of devices by bounded type, state, and availability.", []string{"type", "state", "available"}, labels),
		devicesScanning:     prometheus.NewDesc("decent_devices_scanning", "Whether Decent is scanning for devices.", []string{"phase", "error_kind"}, labels),
		machineTransitions:  prometheus.NewDesc("decent_machine_state_transitions_total", "Machine state transitions observed since exporter start.", []string{"state"}, labels),
	}
}

func (c *Collector) Describe(ch chan<- *prometheus.Desc) {
	for _, desc := range []*prometheus.Desc{
		c.buildInfo, c.streamConnected, c.streamLastMessage, c.streamMessages, c.streamReconnects, c.streamErrors,
		c.machineState, c.machineTemperature, c.machineProfileFrame,
		c.shotsTotal, c.shotDuration, c.shotPeakPressure, c.shotPeakFlow, c.shotAverageFlow,
		c.shotSetting, c.waterLevel, c.displayBool, c.displayBrightness, c.deviceCount, c.devicesScanning, c.machineTransitions,
	} {
		ch <- desc
	}
}

func (c *Collector) Collect(ch chan<- prometheus.Metric) {
	snap := c.store.Snapshot()
	ch <- prometheus.MustNewConstMetric(c.buildInfo, prometheus.GaugeValue, 1, c.build.Version, c.build.Commit, c.build.Date)

	for stream, st := range snap.Streams {
		ch <- prometheus.MustNewConstMetric(c.streamConnected, prometheus.GaugeValue, boolValue(st.Connected), stream)
		if !st.LastMessageTime.IsZero() {
			ch <- prometheus.MustNewConstMetric(c.streamLastMessage, prometheus.GaugeValue, float64(st.LastMessageTime.Unix()), stream)
		}
		ch <- prometheus.MustNewConstMetric(c.streamMessages, prometheus.CounterValue, float64(st.Messages), stream)
		ch <- prometheus.MustNewConstMetric(c.streamReconnects, prometheus.CounterValue, float64(st.Reconnects), stream)
		ch <- prometheus.MustNewConstMetric(c.streamErrors, prometheus.CounterValue, float64(st.Errors), stream)
	}

	if snap.Machine != nil {
		m := snap.Machine
		ch <- prometheus.MustNewConstMetric(c.machineState, prometheus.GaugeValue, 1, cleanLabel(m.State), cleanLabel(m.Substate))
		ch <- prometheus.MustNewConstMetric(c.machineTemperature, prometheus.GaugeValue, m.MixTemperature, "mix")
		ch <- prometheus.MustNewConstMetric(c.machineTemperature, prometheus.GaugeValue, m.GroupTemperature, "group")
		ch <- prometheus.MustNewConstMetric(c.machineTemperature, prometheus.GaugeValue, m.TargetMixTemperature, "target_mix")
		ch <- prometheus.MustNewConstMetric(c.machineTemperature, prometheus.GaugeValue, m.TargetGroupTemperature, "target_group")
		ch <- prometheus.MustNewConstMetric(c.machineTemperature, prometheus.GaugeValue, m.SteamTemperature, "steam")
		ch <- prometheus.MustNewConstMetric(c.machineProfileFrame, prometheus.GaugeValue, m.ProfileFrame)
	}

	ch <- prometheus.MustNewConstMetric(c.shotsTotal, prometheus.CounterValue, float64(snap.ShotsTotal))
	if snap.LastShot != nil {
		shot := snap.LastShot
		ch <- prometheus.MustNewConstMetric(c.shotDuration, prometheus.GaugeValue, shot.Duration.Seconds())
		ch <- prometheus.MustNewConstMetric(c.shotPeakPressure, prometheus.GaugeValue, shot.PeakPressure)
		ch <- prometheus.MustNewConstMetric(c.shotPeakFlow, prometheus.GaugeValue, shot.PeakFlow)
		ch <- prometheus.MustNewConstMetric(c.shotAverageFlow, prometheus.GaugeValue, shot.AverageFlow)
	}

	if snap.Shot != nil {
		emitShot(ch, c.shotSetting, *snap.Shot)
	}
	if snap.Water != nil {
		ch <- prometheus.MustNewConstMetric(c.waterLevel, prometheus.GaugeValue, snap.Water.CurrentLevel, "current")
		ch <- prometheus.MustNewConstMetric(c.waterLevel, prometheus.GaugeValue, snap.Water.RefillLevel, "refill")
	}
	if snap.Display != nil {
		d := snap.Display
		ch <- prometheus.MustNewConstMetric(c.displayBool, prometheus.GaugeValue, boolValue(d.WakeLockEnabled), "wake_lock_enabled")
		ch <- prometheus.MustNewConstMetric(c.displayBool, prometheus.GaugeValue, boolValue(d.WakeLockOverride), "wake_lock_override")
		ch <- prometheus.MustNewConstMetric(c.displayBool, prometheus.GaugeValue, boolValue(d.LowBatteryBrightnessActive), "low_battery_brightness_active")
		ch <- prometheus.MustNewConstMetric(c.displayBool, prometheus.GaugeValue, boolValue(d.BrightnessSupported), "brightness_supported")
		ch <- prometheus.MustNewConstMetric(c.displayBool, prometheus.GaugeValue, boolValue(d.WakeLockSupported), "wake_lock_supported")
		ch <- prometheus.MustNewConstMetric(c.displayBrightness, prometheus.GaugeValue, d.Brightness, "actual")
		ch <- prometheus.MustNewConstMetric(c.displayBrightness, prometheus.GaugeValue, d.RequestedBrightness, "requested")
	}
	if snap.Devices != nil {
		counts := map[[3]string]float64{}
		for _, device := range snap.Devices.Devices {
			key := [3]string{cleanLabel(device.Type), cleanLabel(device.State), strconv.FormatBool(device.Available)}
			counts[key]++
		}
		for key, value := range counts {
			ch <- prometheus.MustNewConstMetric(c.deviceCount, prometheus.GaugeValue, value, key[0], key[1], key[2])
		}
		ch <- prometheus.MustNewConstMetric(c.devicesScanning, prometheus.GaugeValue, boolValue(snap.Devices.Scanning), cleanLabel(snap.Devices.Phase), cleanLabel(snap.Devices.ErrorKind))
	}
	for state, count := range snap.Transitions {
		ch <- prometheus.MustNewConstMetric(c.machineTransitions, prometheus.CounterValue, float64(count), cleanLabel(state))
	}
}

func emitShot(ch chan<- prometheus.Metric, desc *prometheus.Desc, shot decent.ShotSettings) {
	values := map[string]float64{
		"steam_setting":             shot.SteamSetting,
		"target_steam_temp":         shot.TargetSteamTemp,
		"target_steam_duration":     shot.TargetSteamDuration,
		"target_hot_water_temp":     shot.TargetHotWaterTemp,
		"target_hot_water_volume":   shot.TargetHotWaterVolume,
		"target_hot_water_duration": shot.TargetHotWaterSeconds,
		"target_shot_volume":        shot.TargetShotVolume,
		"group_temp":                shot.GroupTemp,
	}
	for name, value := range values {
		ch <- prometheus.MustNewConstMetric(desc, prometheus.GaugeValue, value, name)
	}
}

func boolValue(v bool) float64 {
	if v {
		return 1
	}
	return 0
}

func cleanLabel(v string) string {
	if v == "" {
		return "unknown"
	}
	return v
}
