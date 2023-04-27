package main

import (
	"os"

	basic "com.aco.go.otel/1_basic"
	propagationSample "com.aco.go.otel/2_propagationSample"
	grpc "com.aco.go.otel/3_grpcExample"
)

func main() {

	if len(os.Args) < 2 {
		basic.Start()
	} else {
		switch os.Args[1] {
		case "--propa":
			propagationSample.Start()
		case "--grpc":
			grpc.Start()
		}
	}
}
