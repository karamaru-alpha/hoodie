package core

import "strings"

var CommentReplacer = strings.NewReplacer("//", "", " ", "", "\n", "")
