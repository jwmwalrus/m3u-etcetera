module github.com/jwmwalrus/m3u-etcetera

go 1.19

require (
	github.com/adrg/xdg v0.4.0
	github.com/dhowden/tag v0.0.0-20220618230019-adf36e896086
	github.com/go-gormigrate/gormigrate/v2 v2.0.2
	github.com/go-testfixtures/testfixtures/v3 v3.6.1
	github.com/godbus/dbus/v5 v5.1.0
	github.com/golang/protobuf v1.5.2
	github.com/gotk3/gotk3 v0.6.2-0.20211226093840-cf265f40b836
	github.com/jwmwalrus/bnp v1.12.0
	github.com/jwmwalrus/onerror v0.2.0
	github.com/jwmwalrus/seater v0.1.1
	github.com/nightlyone/lockfile v1.0.0
	github.com/pborman/getopt/v2 v2.1.0
	github.com/rodaine/table v1.0.1
	github.com/sirupsen/logrus v1.9.0
	github.com/stretchr/testify v1.7.1
	github.com/tinyzimmer/go-glib v0.0.25
	github.com/tinyzimmer/go-gst v0.2.33
	github.com/urfave/cli/v2 v2.11.2
	golang.org/x/exp v0.0.0-20220826205824-bd9bcdd0b820
	golang.org/x/text v0.3.7
	google.golang.org/grpc v1.49.0
	google.golang.org/protobuf v1.28.1
	gopkg.in/natefinch/lumberjack.v2 v2.0.0
	gorm.io/driver/sqlite v1.3.6
	gorm.io/gorm v1.23.8
)

require (
	github.com/cpuguy83/go-md2man/v2 v2.0.2 // indirect
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/jinzhu/inflection v1.0.0 // indirect
	github.com/jinzhu/now v1.1.5 // indirect
	github.com/mattn/go-pointer v0.0.1 // indirect
	github.com/mattn/go-sqlite3 v1.14.15 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	github.com/russross/blackfriday/v2 v2.1.0 // indirect
	github.com/xrash/smetrics v0.0.0-20201216005158-039620a65673 // indirect
	golang.org/x/net v0.0.0-20220826154423-83b083e8dc8b // indirect
	golang.org/x/sys v0.0.0-20220829200755-d48e67d00261 // indirect
	google.golang.org/genproto v0.0.0-20220822174746-9e6da59bd2fc // indirect
	gopkg.in/yaml.v2 v2.4.0 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)

// replace github.com/jwmwalrus/bnp => ../bnp
// replace github.com/jwmwalrus/onerror => ../onerror
// replace github.com/tinyzimmer/go-gst => github.com/jwmwalrus/go-gst v0.2.33-0.20220205191536-6a7b117cbaee
