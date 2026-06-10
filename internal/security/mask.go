package security

import "strings"

func MaskPasswordArg(line string) string {
	parts := strings.Fields(line)
	for i, part := range parts {
		if strings.HasPrefix(part, "-p") && len(part) > 2 {
			parts[i] = "-p******"
		}
	}
	return strings.Join(parts, " ")
}
