module github.com/cavaliercoder/rpi_export

go 1.24

replace (
	github.com/cavaliercoder/rpi_export/pkg/export/prometheus => github.com/givanov/rpi_export/pkg/export/prometheus v0.0.0-20240914102107-ce572d9fd7be
	github.com/cavaliercoder/rpi_export/pkg/ioctl => github.com/givanov/rpi_export/pkg/ioctl v0.0.0-20240914102107-ce572d9fd7be
	github.com/cavaliercoder/rpi_export/pkg/mbox => github.com/givanov/rpi_export/pkg/mbox v0.0.0-20240914102107-ce572d9fd7be
)