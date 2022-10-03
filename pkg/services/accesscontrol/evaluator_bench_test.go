package accesscontrol

import (
	"testing"
)

func benchEvaluator(b *testing.B, resourceCount, permissionsPerResource int) {
	permissions, _ := generatePermissions(b, resourceCount, permissionsPerResource)
	permMap := GroupScopesByAction(permissions)
	b.ResetTimer()

	for n := 0; n < b.N; n++ {
		for i := range permissions {
			EvalPermission(permissions[i].Action, permissions[i].Scope).Evaluate(permMap)
		}
	}
}

func BenchmarkMapHasAccess_100_100(b *testing.B)  { benchEvaluator(b, 100, 100) }
func BenchmarkMapHasAccess_1000_100(b *testing.B) { benchEvaluator(b, 1000, 100) }
