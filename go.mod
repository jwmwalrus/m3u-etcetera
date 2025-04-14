module github.com/jwmwalrus/m3u-etcetera

go 1.23.1

toolchain go1.24.2

require (
	github.com/dhowden/tag v0.0.0-20240417053706-3d75831295e8
	github.com/diamondburned/gotk4/pkg v0.3.1
	github.com/go-gormigrate/gormigrate/v2 v2.1.4
	github.com/go-gst/go-glib v1.4.0
	github.com/go-gst/go-gst v1.4.0
	github.com/go-testfixtures/testfixtures/v3 v3.6.1
	github.com/godbus/dbus/v5 v5.1.0
	github.com/jwmwalrus/bnp v1.23.1
	github.com/jwmwalrus/gear-pieces v0.10.4
	github.com/jwmwalrus/quorum v0.11.5
	github.com/jwmwalrus/rtcycler v0.7.2
	github.com/nightlyone/lockfile v1.0.0
	github.com/rodaine/table v1.3.0
	github.com/stretchr/testify v1.10.0
	github.com/urfave/cli/v3 v3.1.1
	golang.org/x/text v0.24.0
	google.golang.org/grpc v1.71.1
	google.golang.org/protobuf v1.36.6
	gorm.io/driver/sqlite v1.5.7
	gorm.io/gorm v1.25.12
)

require (
	github.com/KarpelesLab/weak v0.1.1 // indirect
	github.com/adrg/xdg v0.5.3 // indirect
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/jinzhu/inflection v1.0.0 // indirect
	github.com/jinzhu/now v1.1.5 // indirect
	github.com/mattn/go-pointer v0.0.1 // indirect
	github.com/mattn/go-sqlite3 v1.14.27 // indirect
	github.com/pborman/getopt/v2 v2.1.0 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	go4.org/unsafe/assume-no-moving-gc v0.0.0-20231121144256-b99613f794b6 // indirect
	golang.org/x/exp v0.0.0-20250305212735-054e65f0b394 // indirect
	golang.org/x/net v0.38.0 // indirect
	golang.org/x/sync v0.13.0 // indirect
	golang.org/x/sys v0.32.0 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20250404141209-ee84b53bf3d0 // indirect
	gopkg.in/natefinch/lumberjack.v2 v2.2.1 // indirect
	gopkg.in/yaml.v2 v2.4.0 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)

// replace github.com/jwmwalrus/rtcycler => ../rtcycler

// replace github.com/jwmwalrus/bnp => ../bnp

// replace github.com/jwmwalrus/gear-pieces => ../gear-pieces

// replace github.com/tinyzimmer/go-gst => ../../repos/go-gst
