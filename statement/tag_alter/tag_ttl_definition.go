package tag_alter

import "fmt"

// TTLDefinition ttl_definition:
// TTL_DURATION = ttl_duration, TTL_COL = prop_name
// TTLDefinition represents the TTL definition in the ALTER TAG statement.
type TTLDefinition struct {
	ttlDuration int64
	ttlCol      string
}

// NewTTLDefinition creates a new TTLDefinition with the specified duration and column name.
func NewTTLDefinition(ttlDuration int64, ttlCol string) TTLDefinition {
	return TTLDefinition{
		ttlDuration: ttlDuration,
		ttlCol:      ttlCol,
	}
}

// GenerateTTlDefinitionStatement generates the TTL statement for the ALTER TAG statement.
func GenerateTTlDefinitionStatement(ttl TTLDefinition) (string, error) {
	if ttl.ttlDuration < 0 {
		return "", fmt.Errorf("TTL duration is required to be non zero")
	}
	if ttl.ttlCol == "" {
		return "", fmt.Errorf("TTL column name is required")
	}
	
	return fmt.Sprintf(`TTL_DURATION = %d, TTL_COL = "%s"`, ttl.ttlDuration, ttl.ttlCol), nil
}
