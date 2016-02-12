/*
Copyright 2015 The Kubernetes Authors All rights reserved.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package parsers

import (
	"github.com/docker/distribution/reference"
)

const (
	defaultImageTag = "latest"
)

// parseImageName parses a docker image string into two parts: repo and tag.
// If tag is empty, return the defaultImageTag.
func ParseImageName(image string) (string, string) {
	ref, _ := reference.Parse(image)
	switch ref := ref.(type) {
	default:
		return ref.String(), defaultImageTag
	case reference.NamedTagged:
		return ref.Name(), ref.Tag()
	case reference.Named:
		return ref.Name(), defaultImageTag
	}
}
