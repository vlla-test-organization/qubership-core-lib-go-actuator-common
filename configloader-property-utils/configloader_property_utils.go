package clpropertyutils

import (
	"fmt"
	"strconv"

	"github.com/vlla-test-organization/qubership-core-lib-go/v3/configloader"
)

const (
	defaultHttpBufferHeaderMaxSize = "10240"
	CSHttpBufferHeaderMaxSizeName  = "http.buffer.header.max.size"
)

func GetHttpBufferHeaderMaxSizeBytes() int {
	httpBufferHeaderMaxSizeBytes, err := strconv.Atoi(configloader.GetOrDefaultString(CSHttpBufferHeaderMaxSizeName, defaultHttpBufferHeaderMaxSize))
	if err != nil {
		panic(fmt.Errorf("cannot convert the value for %s property, err: %w", CSHttpBufferHeaderMaxSizeName, err))
	}
	return httpBufferHeaderMaxSizeBytes
}
