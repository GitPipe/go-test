package main

import (
	"fmt"
	"github.com/GitPipe/go-test/docker"
)

func main() {

	imageName := "test/test:latest"
	image := docker.ToImage(imageName)

	fmt.Printf("image: %q\n", image.String())
}
