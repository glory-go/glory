package extend_type

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"github.com/pkg/errors"
)

// JSONTagedGORMType is use to set user defined
type JSONTagedGORMType struct {
	val interface{}
}

// NewJSONTagedGORMType create a JSONTagedGORMType with given @value
func NewJSONTagedGORMType(value interface{}) *JSONTagedGORMType {
	return &JSONTagedGORMType{
		val: value,
	}
}

// Scan is internal exported function that can't be called by user 实现 sql.Scanner 接口，Scan 将 value 扫描至 Jsonb
func (j *JSONTagedGORMType) Scan(value interface{}) error {
	bytes, ok := value.([]byte)
	if !ok {
		return errors.New(fmt.Sprint("Failed to unmarshal JSONB value:", value))
	}

	var result interface{}
	err := json.Unmarshal(bytes, &result)
	(*j).val = result
	return err
}

// Value is internal exported function that can't be called by user  实现 driver.Valuer 接口，Value 返回 json value
func (j *JSONTagedGORMType) Value() (driver.Value, error) {
	return json.Marshal(j.val)
}

// GetInterface is called to get value from db
func (j *JSONTagedGORMType) GetInterface() interface{} {
	return j.val
}
