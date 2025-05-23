package metrics

import (
	"github.com/givanov/rpi_export/pkg/config"
	"github.com/givanov/rpi_export/pkg/mbox"
	"github.com/prometheus/client_golang/prometheus"
	"go.uber.org/zap"
)

type RaspberryPiMboxCollector struct {
	vcRevisionMetric          *prometheus.Desc
	boardVersionMetric        *prometheus.Desc
	boardModelMetric          *prometheus.Desc
	powerStateMetric          *prometheus.Desc
	clockRateHzMetric         *prometheus.Desc
	clockRateMeasuredHzMetric *prometheus.Desc
	turboMetric               *prometheus.Desc
	temperatureCMetric        *prometheus.Desc
	maxTemperatureCMetric     *prometheus.Desc
	voltageMetric             *prometheus.Desc
	voltageMaxMetric          *prometheus.Desc
	voltageMinMetric          *prometheus.Desc
}

func NewRaspberryPiMboxCollector(cfg *config.Config) *RaspberryPiMboxCollector {
	constLabels := prometheus.Labels{
		"host": cfg.HostNameOverride,
	}
	return &RaspberryPiMboxCollector{
		vcRevisionMetric: prometheus.NewDesc("rpi_vc_revision",
			"Firmware revision of the VideoCore device.",
			nil, constLabels,
		),
		boardVersionMetric: prometheus.NewDesc("rpi_board_revision",
			"Board revision.",
			nil, constLabels,
		),
		boardModelMetric: prometheus.NewDesc("rpi_board_model",
			"Board model.",
			nil, constLabels,
		),
		powerStateMetric: prometheus.NewDesc("rpi_power_state",
			"Component power state (0: off, 1: on, 2: missing).",
			[]string{"id"}, constLabels,
		),
		clockRateHzMetric: prometheus.NewDesc("rpi_clock_rate_hz",
			"Clock rate in Hertz.",
			[]string{"id"}, constLabels,
		),
		clockRateMeasuredHzMetric: prometheus.NewDesc("rpi_clock_rate_measured_hz",
			"Measured clock rate in Hertz.",
			[]string{"id"}, constLabels,
		),
		turboMetric: prometheus.NewDesc("rpi_turbo",
			"Turbo state.",
			nil, constLabels,
		),
		temperatureCMetric: prometheus.NewDesc("rpi_temperature_c",
			"Temperature of the SoC in degrees Celsius.",
			[]string{"id"}, constLabels,
		),
		maxTemperatureCMetric: prometheus.NewDesc("rpi_max_temperature_c",
			"Maximum temperature of the SoC in degrees Celsius.",
			[]string{"id"}, constLabels,
		),
		voltageMetric: prometheus.NewDesc("rpi_voltage",
			"Current component voltage.",
			[]string{"id"}, constLabels,
		),
		voltageMaxMetric: prometheus.NewDesc("rpi_voltage_max",
			"Maximum supported component voltage.",
			[]string{"id"}, constLabels,
		),
		voltageMinMetric: prometheus.NewDesc("rpi_voltage_min",
			"Minimum supported component voltage.",
			[]string{"id"}, constLabels,
		),
	}
}

func (collector *RaspberryPiMboxCollector) Describe(ch chan<- *prometheus.Desc) {
	ch <- collector.vcRevisionMetric
	ch <- collector.boardVersionMetric
	ch <- collector.boardModelMetric
	ch <- collector.powerStateMetric
	ch <- collector.clockRateHzMetric
	ch <- collector.clockRateMeasuredHzMetric
	ch <- collector.turboMetric
	ch <- collector.temperatureCMetric
	ch <- collector.maxTemperatureCMetric
	ch <- collector.voltageMetric
	ch <- collector.voltageMaxMetric
	ch <- collector.voltageMinMetric
}

func (collector *RaspberryPiMboxCollector) Collect(ch chan<- prometheus.Metric) {
	m, err := mbox.Open()
	if err != nil {
		return
	}
	defer m.Close()

	/*
	 * Hardware.
	 */
	rev, err := m.GetFirmwareRevision()
	if err != nil {
		zap.L().Error("Error getting firmware revision", zap.Error(err))
		return
	}
	fwRevMetric := prometheus.MustNewConstMetric(collector.vcRevisionMetric, prometheus.GaugeValue, float64(rev))
	ch <- fwRevMetric

	rev, err = m.GetBoardRevision()
	if err != nil {
		zap.L().Error("Error getting board revision", zap.Error(err))
		return
	}
	boardRevMetric := prometheus.MustNewConstMetric(collector.boardVersionMetric, prometheus.GaugeValue, float64(rev))
	ch <- boardRevMetric

	// Extract board model from a revision number.
	// Revision 17 and later is Raspberry Pi 5 that does not have a board model.
	if (rev>>4)&0xff < 17 {
		model, err := m.GetBoardModel()
		if err != nil {
			zap.L().Error("Error getting board model", zap.Error(err))
			return
		}
		boardModelMetric := prometheus.MustNewConstMetric(collector.boardModelMetric, prometheus.GaugeValue, float64(model))
		ch <- boardModelMetric
	}

	/*
	 * Power.
	 */
	for id, label := range powerLabelsByID {
		powerState, err := m.GetPowerState(id)
		if err != nil {
			zap.L().Error("Error getting power state", zap.Error(err))
			return
		}
		boardPowerStateMetric := prometheus.MustNewConstMetric(collector.powerStateMetric, prometheus.GaugeValue, float64(powerState), label)
		ch <- boardPowerStateMetric
	}

	/*
	 * Clocks.
	 */
	for id, label := range clockLabelsByID {
		clockRate, err := m.GetClockRate(id)
		if err != nil {
			zap.L().Error("Error getting clock rate", zap.Error(err))
			return
		}
		clockRateMetric := prometheus.MustNewConstMetric(collector.clockRateHzMetric, prometheus.GaugeValue, float64(clockRate), label)
		ch <- clockRateMetric
	}

	for id, label := range clockLabelsByID {
		clockRate, err := m.GetClockRateMeasured(id)
		if err != nil {
			zap.L().Error("Error getting clock rate measured", zap.Error(err))
			return
		}
		clockRateMeasuredMetric := prometheus.MustNewConstMetric(collector.clockRateMeasuredHzMetric, prometheus.GaugeValue, float64(clockRate), label)
		ch <- clockRateMeasuredMetric
	}

	turbo, err := m.GetTurbo()
	if err != nil {
		zap.L().Error("Error getting turbo", zap.Error(err))
		return
	}
	turboVal := float64(0)
	if turbo {
		turboVal = float64(1)
	}

	turboMetric := prometheus.MustNewConstMetric(collector.turboMetric, prometheus.GaugeValue, turboVal)
	ch <- turboMetric

	/*
	 * Temperature sensors.
	 */

	// Current SoC temperature
	temp, err := m.GetTemperature()
	if err != nil {
		zap.L().Error("Error getting temperature", zap.Error(err))
		return
	}
	temperatureCMetric := prometheus.MustNewConstMetric(collector.temperatureCMetric, prometheus.GaugeValue, float64(temp), "soc")
	ch <- temperatureCMetric

	maxTemp, err := m.GetMaxTemperature()
	if err != nil {
		zap.L().Error("Error getting max temperature", zap.Error(err))
		return
	}

	// Max SoC temperature
	maxTemperatureCMetric := prometheus.MustNewConstMetric(collector.maxTemperatureCMetric, prometheus.GaugeValue, float64(maxTemp), "soc")
	ch <- maxTemperatureCMetric

	/*
	 * Voltages
	 */

	// Current voltages.
	for id, label := range voltageLabelsByID {
		volts, err := m.GetVoltage(id)
		if err != nil {
			zap.L().Error("Error getting voltage", zap.Error(err))
			return
		}
		voltageMetric := prometheus.MustNewConstMetric(collector.voltageMetric, prometheus.GaugeValue, float64(volts), label)
		ch <- voltageMetric
	}

	for id, label := range voltageLabelsByID {
		volts, err := m.GetMinVoltage(id)
		if err != nil {
			zap.L().Error("Error getting voltage", zap.Error(err))
			return
		}
		voltageMinMetric := prometheus.MustNewConstMetric(collector.voltageMinMetric, prometheus.GaugeValue, float64(volts), label)
		ch <- voltageMinMetric
	}

	for id, label := range voltageLabelsByID {
		volts, err := m.GetMaxVoltage(id)
		if err != nil {
			zap.L().Error("Error getting voltage", zap.Error(err))
			return
		}
		voltageMaxMetric := prometheus.MustNewConstMetric(collector.voltageMaxMetric, prometheus.GaugeValue, float64(volts), label)
		ch <- voltageMaxMetric
	}
}
