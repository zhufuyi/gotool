package utils

import (
	"fmt"
	"strings"
)

// ListTypeNames 拼接类型名称
func ListTypeNames(names ...string) string {
	content := []string{fmt.Sprintf("%d types are supported:\n", len(names))}
	for _, name := range names {
		content = append(content, "    "+name+"\n")
	}

	return strings.Join(content, "")
}
