package batched

import (
	"github.com/ischenkx/kantoku/pkg/lib/impl/data/dependency/inmem"
	"testing"
)

func newInMemDeps() *inmem.Manager {
	return inmem.New()
}

func BatchedIntegrationTest(t *testing.T) {

}
