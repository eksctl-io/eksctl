package serverless

import (
	"encoding/json"
	"sort"

	"github.com/awslabs/goformation/v4/cloudformation/utils"
)

// StateMachine_Properties is a helper struct that can hold either a CloudWatchEventEvent, EventBridgeRuleEvent, ScheduleEvent, or ApiEvent value
type StateMachine_Properties struct {
	CloudWatchEventEvent *StateMachine_CloudWatchEventEvent
	EventBridgeRuleEvent *StateMachine_EventBridgeRuleEvent
	ScheduleEvent        *StateMachine_ScheduleEvent
	ApiEvent             *StateMachine_ApiEvent
}

func (r StateMachine_Properties) value() interface{} {
	ret := []interface{}{}

	if r.CloudWatchEventEvent != nil {
		ret = append(ret, *r.CloudWatchEventEvent)
	}

	if r.EventBridgeRuleEvent != nil {
		ret = append(ret, *r.EventBridgeRuleEvent)
	}

	if r.ScheduleEvent != nil {
		ret = append(ret, *r.ScheduleEvent)
	}

	if r.ApiEvent != nil {
		ret = append(ret, *r.ApiEvent)
	}

	sort.Sort(utils.ByJSONLength(ret)) // Heuristic to select best attribute
	if len(ret) > 0 {
		return ret[0]
	}

	return nil
}

func (r StateMachine_Properties) MarshalJSON() ([]byte, error) {
	return json.Marshal(r.value())
}

// Hook into the marshaller
func (r *StateMachine_Properties) UnmarshalJSON(b []byte) error {

	// Unmarshal into interface{} to check it's type
	var typecheck interface{}
	if err := json.Unmarshal(b, &typecheck); err != nil {
		return err
	}

	switch val := typecheck.(type) {

	case map[string]interface{}:
		val = val // This ensures val is used to stop an error

		json.Unmarshal(b, &r.CloudWatchEventEvent)

		json.Unmarshal(b, &r.EventBridgeRuleEvent)

		json.Unmarshal(b, &r.ScheduleEvent)

		json.Unmarshal(b, &r.ApiEvent)

	case []interface{}:

	}

	return nil
}
