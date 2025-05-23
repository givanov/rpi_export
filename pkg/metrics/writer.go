/*
Package metrics provides utilities for authoring a Prometheus metric exporter for the metrics
found in the Mailbox Property Interface of a Raspberry Pi.
*/
package metrics

import (
	"github.com/givanov/rpi_export/pkg/mbox"
)

var voltageLabelsByID = map[mbox.VoltageID]string{
	mbox.VoltageIDCore:   "core",
	mbox.VoltageIDSDRAMC: "sdram_c",
	mbox.VoltageIDSDRAMI: "sdram_i",
	mbox.VoltageIDSDRAMP: "sdram_p",
}

var powerLabelsByID = map[mbox.PowerDeviceID]string{
	mbox.PowerDeviceIDSDCard: "sd_card",
	mbox.PowerDeviceIDUART0:  "uart0",
	mbox.PowerDeviceIDUART1:  "uart1",
	mbox.PowerDeviceIDUSBHCD: "usb_hcd",
	mbox.PowerDeviceIDI2C0:   "i2c0",
	mbox.PowerDeviceIDI2C1:   "i2c1",
	mbox.PowerDeviceIDI2C2:   "i2c2",
	mbox.PowerDeviceIDSPI:    "spi",
	mbox.PowerDeviceIDCCP2TX: "ccp2tx",
}

var clockLabelsByID = map[mbox.ClockID]string{
	mbox.ClockIDEMMC:     "emmc",
	mbox.ClockIDUART:     "uart",
	mbox.ClockIDARM:      "arm",
	mbox.ClockIDCore:     "core",
	mbox.ClockIDV3D:      "v3d",
	mbox.ClockIDH264:     "h264",
	mbox.ClockIDISP:      "isp",
	mbox.ClockIDSDRAM:    "sdram",
	mbox.ClockIDPixel:    "pixel",
	mbox.ClockIDPWM:      "pwm",
	mbox.ClockIDHEVC:     "hevc",
	mbox.ClockIDEMMC2:    "emmc2",
	mbox.ClockIDM2MC:     "m2mc",
	mbox.ClockIDPixelBVB: "pixel_bvb",
}
