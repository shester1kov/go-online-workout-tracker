package utils

import "github.com/gosimple/slug"

func GenerateSlug(name string) string {
	return slug.Make(name)
}
