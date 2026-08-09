package projects

import "testing"

func TestNormalizeDeploymentStatus(t *testing.T) {
	if got := normalizeDeploymentStatus(""); got != StatusQueued {
		t.Fatalf("expected queued status by default, got %q", got)
	}

	if got := normalizeDeploymentStatus(StatusBuilding); got != StatusBuilding {
		t.Fatalf("expected building status to be preserved, got %q", got)
	}
}
