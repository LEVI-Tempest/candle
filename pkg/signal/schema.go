package signal

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/xeipuuv/gojsonschema"
)

// ValidateReportSchema validates Report against JSON schema.
// ValidateReportSchema 使用 JSON schema 校验 Report。
func ValidateReportSchema(report Report, schemaPath string) error {
	data, err := json.Marshal(report)
	if err != nil {
		return err
	}

	schemaLoader := gojsonschema.NewReferenceLoader("file://" + schemaPath)
	docLoader := gojsonschema.NewBytesLoader(data)
	result, err := gojsonschema.Validate(schemaLoader, docLoader)
	if err != nil {
		return err
	}
	if result.Valid() {
		return nil
	}

	errs := result.Errors()
	msgs := make([]string, 0, len(errs))
	for _, e := range errs {
		msgs = append(msgs, e.String())
	}
	return fmt.Errorf("schema validation failed: %s", strings.Join(msgs, "; "))
}
