package accesscontrol

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func generatePermissions(b *testing.B, resourceCount, permissionPerResource int) (map[string][]string, map[string]bool) {
	res := map[string][]string{}
	ids := make(map[string]bool, resourceCount)

	for p := 0; p < permissionPerResource; p++ {
		action := fmt.Sprintf("resources:action%v", p)
		for r := 0; r < resourceCount; r++ {
			scope := fmt.Sprintf("resources:id:%v", r)
			res[action] = append(res[action], scope)
			ids[scope] = true
		}
	}

	return res, ids
}

func benchGetTrieMetadata(b *testing.B, resourceCount, permissionPerResource int) {
	permissions, ids := generatePermissions(b, resourceCount, permissionPerResource)
	trie := TrieFromMap(permissions)
	metas := make(map[string]Metadata)
	b.ResetTimer()

	for n := 0; n < b.N; n++ {
		for id := range ids {
			metas[id] = trie.Metadata(id)
		}
		assert.Len(b, metas, resourceCount)
		for _, resourceMetadata := range metas {
			assert.Len(b, resourceMetadata, permissionPerResource)
		}
	}
}

// Lots of permissions
func BenchmarkTrieMetadata_10_1000(b *testing.B)   { benchGetTrieMetadata(b, 10, 1000) }
func BenchmarkTrieMetadata_10_10000(b *testing.B)  { benchGetTrieMetadata(b, 10, 10000) }
func BenchmarkTrieMetadata_10_100000(b *testing.B) { benchGetTrieMetadata(b, 10, 100000) }
func BenchmarkTrieMetadata_10_1000000(b *testing.B) {
	if testing.Short() {
		b.Skip("Skipping benchmark in short mode")
	}
	benchGetTrieMetadata(b, 10, 1000000)
}

// Lots of resources
func BenchmarkTrieMetadata_1000_10(b *testing.B)   { benchGetTrieMetadata(b, 1000, 10) }
func BenchmarkTrieMetadata_10000_10(b *testing.B)  { benchGetTrieMetadata(b, 10000, 10) }
func BenchmarkTrieMetadata_100000_10(b *testing.B) { benchGetTrieMetadata(b, 100000, 10) }
func BenchmarkTrieMetadata_1000000_10(b *testing.B) {
	if testing.Short() {
		b.Skip("Skipping benchmark in short mode")
	}
	benchGetTrieMetadata(b, 1000000, 10)
}
