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
	ContentEncoding      string                 `json:"contentEncoding,omitempty"`
	OneOf                []*Definition          `json:"oneOf,omitempty"`
	Description          string                 `json:"description,omitempty"`
	HTMLDescription      string                 `json:"x-intellij-html-description,omitempty"`
	KubernetesGvk        []*GroupVersionKind    `json:"x-kubernetes-group-version-kind,omitempty"`
	Default              interface{}            `json:"default,omitempty"`
	Examples             []string               `json:"examples,omitempty"`
	Enum                 []string               `json:"enum,omitempty"`
}

type GroupVersionKind struct {
	Group   string `json:"group"`
	Version string `json:"version"`
	Kind    string `json:"kind"`
}
