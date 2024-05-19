package core

import "strings"

var CommentReplacer = strings.NewReplacer("//", "", " ", "", "\n", "")

func IsTimeField(snakeName string) bool {
	return strings.HasSuffix(snakeName, "_time")
}
