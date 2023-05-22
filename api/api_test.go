package api

import (
	"google.golang.org/protobuf/proto"
)

type testCase struct {
	name        string
	fixturesDir fixturesDir
	req         proto.Message
	res         proto.Message
	wantErr     bool
}

type fixturesDir string

func (fg fixturesDir) FixturesDir() string {
	return string(fg)
}
