package v1alpha5

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strings"

	ekstypes "github.com/aws/aws-sdk-go-v2/service/eks/types"
)

// ResolveConflictsType determines how to resolve field value conflicts for an EKS add-on if a value was changed from default
type ResolveConflicts int

const (
	// None – EKS doesn't change the value. The update might fail.
	None ResolveConflicts = iota
	// Overwrite – EKS overwrites the changed value back to default.
	Overwrite
	// Preserve – EKS preserves the value.
	Preserve
)

var stringToResolveConflicts = map[string]ResolveConflicts{
	"none":      None,
	"overwrite": Overwrite,
	"preserve":  Preserve,
}

var resolveConflictsToString = map[ResolveConflicts]string{
	None:      "none",
	Overwrite: "overwrite",
	Preserve:  "preserve",
}

// MarshalJSON implements json.Marshaler
func (rc *ResolveConflicts) MarshalJSON() ([]byte, error) {
	buffer := bytes.NewBufferString(`"`)
	buffer.WriteString(resolveConflictsToString[*rc])
	buffer.WriteString(`"`)
	return buffer.Bytes(), nil
}

// UnmarshalJSON implements json.Unmarshaler
func (rc *ResolveConflicts) UnmarshalJSON(b []byte) error {
	var s string
	if err := json.Unmarshal(b, &s); err != nil {
		return err
	}
	resolveConflicts, ok := stringToResolveConflicts[strings.ToLower(s)]
	if !ok {
		return fmt.Errorf("%q is not a valid resolveConflict value", s)
	}
	*rc = resolveConflicts
	return nil
}

// ToEKSType converts internal ResolveConflicts type to AWS EKS ResolveConflicts type
func (rc *ResolveConflicts) ToEKSType() ekstypes.ResolveConflicts {
	switch *rc {
	case Overwrite:
		return ekstypes.ResolveConflictsOverwrite
	case Preserve:
		return ekstypes.ResolveConflictsPreserve
	default:
		return ekstypes.ResolveConflictsNone
	}
}
