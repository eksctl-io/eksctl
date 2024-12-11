package serverless

import (
	"encoding/json"
	"sort"

	"goformation/v4/cloudformation/types"
	"goformation/v4/cloudformation/utils"
)

// StateMachine_DefinitionUri is a helper struct that can hold either a String or S3Location value
type StateMachine_DefinitionUri struct {
	String **types.Value

	S3Location *StateMachine_S3Location
}

func (r StateMachine_DefinitionUri) value() interface{} {
	ret := []interface{}{}

	if r.String != nil {
		ret = append(ret, r.String)
	}

	if r.S3Location != nil {
		ret = append(ret, *r.S3Location)
	}

	sort.Sort(utils.ByJSONLength(ret)) // Heuristic to select best attribute
	if len(ret) > 0 {
		return ret[0]
	}

	return nil
}

func (r StateMachine_DefinitionUri) MarshalJSON() ([]byte, error) {
	return json.Marshal(r.value())
}

// Hook into the marshaller
func (r *StateMachine_DefinitionUri) UnmarshalJSON(b []byte) error {

	// Unmarshal into interface{} to check it's type
	var typecheck interface{}
	if err := json.Unmarshal(b, &typecheck); err != nil {
		return err
	}

	switch val := typecheck.(type) {

	case string:
		v, err := types.NewValueFromPrimitive(val)
		if err != nil {
			return err
		}
		r.String = &v

	case map[string]interface{}:
		val = val // This ensures val is used to stop an error

		json.Unmarshal(b, &r.S3Location)

	case []interface{}:

	}

	return nil
}
