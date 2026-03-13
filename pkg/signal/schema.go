package signal

import (
	"encoding/json"
	"fmt"

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

	msg := "schema validation failed"
	if len(result.Errors()) > 0 {
		msg = fmt.Sprintf("%s: %s", msg, result.Errors()[0].String())
	}
	return fmt.Errorf(msg)
}
