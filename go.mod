module github.com/jwmwalrus/m3u-etcetera

go 1.21

require (
	github.com/dhowden/tag v0.0.0-20230630033851-978a0926ee25
	github.com/go-gormigrate/gormigrate/v2 v2.1.0
	github.com/go-testfixtures/testfixtures/v3 v3.6.1
	github.com/godbus/dbus/v5 v5.1.0
	github.com/gotk3/gotk3 v0.6.2
	github.com/jwmwalrus/bnp v1.16.2
	github.com/jwmwalrus/quorum v0.11.2
	github.com/jwmwalrus/rtcycler v0.6.2
	github.com/nightlyone/lockfile v1.0.0
	github.com/rodaine/table v1.1.0
	github.com/stretchr/testify v1.8.4
	github.com/tinyzimmer/go-glib v0.0.25
	github.com/tinyzimmer/go-gst v0.2.33
	github.com/urfave/cli/v2 v2.25.7
	golang.org/x/text v0.12.0
	google.golang.org/grpc v1.57.0
	google.golang.org/protobuf v1.31.0
	gorm.io/driver/sqlite v1.5.3
	gorm.io/gorm v1.25.4
)

require (
	github.com/adrg/xdg v0.4.0 // indirect
	github.com/cpuguy83/go-md2man/v2 v2.0.2 // indirect
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/golang/protobuf v1.5.3 // indirect
	github.com/jinzhu/inflection v1.0.0 // indirect
	github.com/jinzhu/now v1.1.5 // indirect
	github.com/mattn/go-pointer v0.0.1 // indirect
	github.com/mattn/go-sqlite3 v1.14.17 // indirect
	github.com/niemeyer/pretty v0.0.0-20200227124842-a10e7caefd8e // indirect
	github.com/pborman/getopt/v2 v2.1.0 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	github.com/russross/blackfriday/v2 v2.1.0 // indirect
	github.com/xrash/smetrics v0.0.0-20201216005158-039620a65673 // indirect
	golang.org/x/net v0.14.0 // indirect
	golang.org/x/sys v0.11.0 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20230822172742-b8732ec3820d // indirect
	gopkg.in/check.v1 v1.0.0-20200227125254-8fa46927fb4f // indirect
	gopkg.in/natefinch/lumberjack.v2 v2.2.1 // indirect
	gopkg.in/yaml.v2 v2.4.0 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)

// replace github.com/jwmwalrus/rtcycler => ../rtcycler
// replace github.com/jwmwalrus/bnp => ../bnp
// replace github.com/tinyzimmer/go-gst => ../../repos/go-gst
