module github.com/jwmwalrus/m3u-etcetera

go 1.22

toolchain go1.22.2

require (
	github.com/dhowden/tag v0.0.0-20240417053706-3d75831295e8
	github.com/diamondburned/gotk4/pkg v0.2.2
	github.com/go-gormigrate/gormigrate/v2 v2.1.2
	github.com/go-gst/go-glib v1.0.1
	github.com/go-gst/go-gst v1.0.0
	github.com/go-testfixtures/testfixtures/v3 v3.6.1
	github.com/godbus/dbus/v5 v5.1.0
	github.com/jwmwalrus/bnp v1.22.1
	github.com/jwmwalrus/gear-pieces v0.10.2
	github.com/jwmwalrus/quorum v0.11.3
	github.com/jwmwalrus/rtcycler v0.7.0
	github.com/nightlyone/lockfile v1.0.0
	github.com/rodaine/table v1.2.0
	github.com/stretchr/testify v1.9.0
	github.com/urfave/cli/v2 v2.27.1
	golang.org/x/text v0.14.0
	google.golang.org/grpc v1.63.2
	google.golang.org/protobuf v1.33.0
	gorm.io/driver/sqlite v1.5.5
	gorm.io/gorm v1.25.9
)

require (
	github.com/KarpelesLab/weak v0.1.1 // indirect
	github.com/adrg/xdg v0.4.0 // indirect
	github.com/cpuguy83/go-md2man/v2 v2.0.4 // indirect
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/golang/protobuf v1.5.4 // indirect
	github.com/jinzhu/inflection v1.0.0 // indirect
	github.com/jinzhu/now v1.1.5 // indirect
	github.com/mattn/go-pointer v0.0.1 // indirect
	github.com/mattn/go-sqlite3 v1.14.22 // indirect
	github.com/niemeyer/pretty v0.0.0-20200227124842-a10e7caefd8e // indirect
	github.com/pborman/getopt/v2 v2.1.0 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	github.com/russross/blackfriday/v2 v2.1.0 // indirect
	github.com/xrash/smetrics v0.0.0-20240312152122-5f08fbb34913 // indirect
	go4.org/unsafe/assume-no-moving-gc v0.0.0-20231121144256-b99613f794b6 // indirect
	golang.org/x/net v0.24.0 // indirect
	golang.org/x/sync v0.7.0 // indirect
	golang.org/x/sys v0.19.0 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20240415180920-8c6c420018be // indirect
	gopkg.in/check.v1 v1.0.0-20200227125254-8fa46927fb4f // indirect
	gopkg.in/natefinch/lumberjack.v2 v2.2.1 // indirect
	gopkg.in/yaml.v2 v2.4.0 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)

// replace github.com/jwmwalrus/rtcycler => ../rtcycler

// replace github.com/jwmwalrus/bnp => ../bnp

// replace github.com/jwmwalrus/gear-pieces => ../gear-pieces

// replace github.com/tinyzimmer/go-gst => ../../repos/go-gst
