package main

import (
	"context"

	"github.com/redpanda-data/benthos/v4/public/service"

	// Import full suite of FOSS connect plugins
	_ "github.com/redpanda-data/connect/public/bundle/free/v4"

	// Or, in order to import both FOSS and enterprise plugins, replace the
	// above with:
	// _ "github.com/redpanda-data/connect/public/bundle/enterprise/v4"

	// Add your plugin packages here
	_ "github.com/insidegreen/redpanda-connect-plugin-example/input/numa-reduce"
)

func main() {
	service.RunCLI(context.Background())
}
