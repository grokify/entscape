package schema

import (
	"encoding/json"

	"github.com/invopop/jsonschema"
)

// GenerateJSONSchema generates a JSON Schema for the Schema type.
func GenerateJSONSchema() ([]byte, error) {
	r := &jsonschema.Reflector{
		DoNotReference: true,
	}

	s := r.Reflect(&Schema{})
	s.ID = "https://github.com/grokify/entscape/schema/entscape.schema.json"
	s.Title = "Entscape Schema"
	s.Description = "JSON schema for entscape visualization data"

	return json.MarshalIndent(s, "", "  ")
}
