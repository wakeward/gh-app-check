package ghclient

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/google/go-github/v89/github"
)

// PermissionsMap flattens go-github's InstallationPermissions struct into the
// api_key -> access map used by the rules engine and toxic-combination eval.
func PermissionsMap(perms *github.InstallationPermissions) (map[string]string, error) {
	if perms == nil {
		return map[string]string{}, nil
	}
	raw, err := json.Marshal(perms)
	if err != nil {
		return nil, fmt.Errorf("marshal installation permissions: %w", err)
	}
	var decoded map[string]*string
	if err := json.Unmarshal(raw, &decoded); err != nil {
		return nil, fmt.Errorf("decode installation permissions: %w", err)
	}
	out := make(map[string]string, len(decoded))
	for key, value := range decoded {
		if value == nil {
			continue
		}
		access := strings.ToLower(strings.TrimSpace(*value))
		if access == "" {
			continue
		}
		out[key] = access
	}
	return out, nil
}
