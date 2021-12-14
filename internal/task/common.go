package task

import (
	"fmt"
	"net/url"
	"strconv"

	"google.golang.org/grpc"
)

func getGrpcOpts() (opts []grpc.DialOption) {
	opts = append(opts, grpc.WithInsecure())
	return
}

func parseIDs(ids []string) (parsed []int64, err error) {
	for _, v := range ids {
		var aux int64
		if aux, err = strconv.ParseInt(v, 10, 64); err != nil {
			return
		}
		if aux < 1 {
			err = fmt.Errorf("Found invalid ID: %v", aux)
			return
		}
		parsed = append(parsed, aux)
	}
	return
}

func parseLocations(locations []string) (parsed []string, err error) {
	for _, v := range locations {
		var u *url.URL
		if u, err = url.Parse(v); err != nil {
			return
		}
		if u.Scheme == "" {
			u.Scheme = "file"
		}
		parsed = append(parsed, u.String())
	}
	return
}
