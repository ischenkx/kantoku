package uid

import (
	"github.com/google/uuid"
	"strings"
)

func Generate() string {
	return strings.ReplaceAll(uuid.New().String(), "-", "")
}
