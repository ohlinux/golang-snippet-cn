package main

import (
	. "db"
	. "fetcher"
)

// test
func main() {
	module := &PreFetch{"param", 1, 1, "scm", "inf/odp/orp/test-app", "1.0.1.0"}

	FetchPackage(module)
}
