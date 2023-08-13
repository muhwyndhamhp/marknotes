package typeext

import (
	"database/sql/driver"
	"encoding/json"
	"errors"

	"github.com/muhwyndhamhp/gotes-mx/utils/errs"
)

type JSONB map[string]interface{}

// Value Marshal
func (m JSONB) Value() (driver.Value, error) {
	return json.Marshal(m)
}

// Scan Unmarshal
func (m *JSONB) Scan(value interface{}) error {
	var source []byte
	_m := make(map[string]interface{})

	switch v := value.(type) {
	case []uint8:
		source = v
	case string:
		source = []byte(v)
	case nil:
		return nil
	default:
		return errors.New("incompatible type for StringInterfaceMap")
	}
	err := json.Unmarshal(source, &_m)
	if err != nil {
		return errs.Wrap(err)
	}
	*m = _m
	return nil
}

func ConvertStructToJSONB(value interface{}) (JSONB, error) {
	var result JSONB
	js, err := json.Marshal(value)
	if err != nil {
		return nil, errs.Wrap(err)
	}

	err = json.Unmarshal(js, &result)
	if err != nil {
		return nil, errs.Wrap(err)
	}

	return result, nil
}
