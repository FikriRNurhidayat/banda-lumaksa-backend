package version

import (
	"fmt"
	"time"
)

var (
	Version = "dev"
	Build   = fmt.Sprintf("%d", time.Now().UnixMilli())
)
