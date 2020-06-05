package definition

// Definition represents a JSON Schema definition
type Definition struct {
	Ref                  string                 `json:"$ref,omitempty"`
	Items                *Definition            `json:"items,omitempty"`
	Required             []string               `json:"required,omitempty"`
	Properties           map[string]*Definition `json:"properties,omitempty"`
	PreferredOrder       []string               `json:"preferredOrder,omitempty"`
	AdditionalProperties interface{}            `json:"additionalProperties,omitempty"`
	Type                 string                 `json:"type,omitempty"`
	AnyOf                []*Definition          `json:"anyOf,omitempty"`
	Description          string                 `json:"description,omitempty"`
	HTMLDescription      string                 `json:"x-intellij-html-description,omitempty"`
	Default              interface{}            `json:"default,omitempty"`
	Examples             []string               `json:"examples,omitempty"`
	Enum                 []string               `json:"enum,omitempty"`
}
