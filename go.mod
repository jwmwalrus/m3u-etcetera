module github.com/jwmwalrus/m3u-etcetera

go 1.18

require (
	github.com/adrg/xdg v0.4.0
	github.com/dhowden/tag v0.0.0-20220618230019-adf36e896086
	github.com/go-gormigrate/gormigrate/v2 v2.0.2
	github.com/go-testfixtures/testfixtures/v3 v3.6.1
	github.com/godbus/dbus/v5 v5.1.0
	github.com/gotk3/gotk3 v0.6.2-0.20211226093840-cf265f40b836
	github.com/jwmwalrus/bnp v1.10.0
	github.com/jwmwalrus/onerror v0.1.1
	github.com/jwmwalrus/seater v0.1.1
	github.com/nightlyone/lockfile v1.0.0
	github.com/pborman/getopt/v2 v2.1.0
	github.com/rodaine/table v1.0.1
	github.com/sirupsen/logrus v1.8.1
	github.com/stretchr/testify v1.7.1
	github.com/tinyzimmer/go-glib v0.0.25
	github.com/tinyzimmer/go-gst v0.2.33
	github.com/urfave/cli/v2 v2.10.2
	golang.org/x/exp v0.0.0-20220613132600-b0d781184e0d
	google.golang.org/grpc v1.47.0
	google.golang.org/protobuf v1.28.0
	gopkg.in/natefinch/lumberjack.v2 v2.0.0
	gorm.io/driver/sqlite v1.3.4
	gorm.io/gorm v1.23.6
)

require (
	github.com/cpuguy83/go-md2man/v2 v2.0.2 // indirect
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/golang/protobuf v1.5.2 // indirect
	github.com/jinzhu/inflection v1.0.0 // indirect
	github.com/jinzhu/now v1.1.5 // indirect
	github.com/mattn/go-pointer v0.0.1 // indirect
	github.com/mattn/go-sqlite3 v1.14.13 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	github.com/russross/blackfriday/v2 v2.1.0 // indirect
	github.com/xrash/smetrics v0.0.0-20201216005158-039620a65673 // indirect
	golang.org/x/net v0.0.0-20220621193019-9d032be2e588 // indirect
	golang.org/x/sys v0.0.0-20220615213510-4f61da869c0c // indirect
	golang.org/x/text v0.3.7 // indirect
	google.golang.org/genproto v0.0.0-20220621134657-43db42f103f7 // indirect
	gopkg.in/yaml.v2 v2.4.0 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)

// replace github.com/jwmwalrus/bnp => ../bnp
// replace github.com/tinyzimmer/go-gst => github.com/jwmwalrus/go-gst v0.2.33-0.20220205191536-6a7b117cbaee
