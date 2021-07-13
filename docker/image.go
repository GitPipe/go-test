package docker

import (
	"fmt"
	"log"
	"regexp"
	"strings"
)

var ecrRegistryRegex *regexp.Regexp

func init() {

	account := `[\d]{12}`
	region := `(us(-gov)?|ap|ca|cn|eu|sa)-(central|(north|south)?(east|west)?)-\d`
	r, err := regexp.Compile(fmt.Sprintf("%s.dkr.ecr.%s.amazonaws.com", account, region))
	if err != nil {
		log.Fatalf("cannot compile ecr registry regex")
	}
	ecrRegistryRegex = r
}

type Image struct {
	Registry   string
	Repository string
	Tag        string
}

// used by dvyukov/go-fuzz project
func Fuzz(data []byte) int {

	imageName := string(data)
	imageName2 := ToImage(imageName).String()
	if imageName != imageName2 {
		panic(fmt.Sprintf("image name %q does not equal to image name 2 %q", imageName, imageName2))
	}
	return 1
}

func ToImage(image string) Image {

	// get tag
	var img Image
	imageTagParts := strings.Split(image, ":")
	if len(imageTagParts) != 1 {
		img.Tag = strings.Join(imageTagParts[1:], ":")
	}

	// get repository if it exists
	imageWithoutVersion := imageTagParts[0]
	imageParts := strings.Split(imageWithoutVersion, "/")
	if len(imageParts) > 1 &&
		(strings.HasPrefix(imageParts[0], "localhost") || strings.Contains(imageParts[0], ".")) {
		// docker.io is default registry and is not included in repository names, if we included this
		// then it is hard to get image id (image id has to be searched by repository name)
		if imageParts[0] != "docker.io" {
			img.Registry = imageParts[0]
		}
		img.Repository = strings.Join(imageParts[1:], "/")
		return img
	}

	img.Repository = imageWithoutVersion
	return img
}

func (i Image) String() string {

	image := i.Repository
	if i.Tag != "" {
		image = fmt.Sprintf("%s:%s", image, i.Tag)
	}
	if i.Registry != "" {
		image = fmt.Sprintf("%s/%s", i.Registry, image)
	}
	return image
}
