package execution

import "testing"

func TestValidateChangedFilesAllowsDeclaredScope(t *testing.T) {
	if err := ValidateChangedFiles([]string{"apps/web/src/App.tsx"}, []string{"apps/web/**"}, []string{"secrets/**"}); err != nil {
		t.Fatalf("expected allowed path: %v", err)
	}
}

func TestValidateChangedFilesRejectsForbiddenPath(t *testing.T) {
	if err := ValidateChangedFiles([]string{"secrets/token.txt"}, []string{"**"}, []string{"secrets/**"}); err != ErrScopeViolation {
		t.Fatalf("expected scope violation, got %v", err)
	}
}

func TestValidateChangedFilesRejectsTraversal(t *testing.T) {
	if err := ValidateChangedFiles([]string{"../outside.txt"}, []string{"**"}, nil); err != ErrScopeViolation {
		t.Fatalf("expected traversal violation, got %v", err)
	}
}
