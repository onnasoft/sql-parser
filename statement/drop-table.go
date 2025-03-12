package statement

import (
	"fmt"

	"github.com/asaskevich/govalidator"
	"github.com/onnasoft/ZenithSQL/protocol"
	"github.com/vmihailenco/msgpack/v5"
)

type DropTableStatement struct {
	TableName string `msgpack:"table_name" valid:"required,alphanumunderscore"`
}

func NewDropTableStatement(tableName string) (*DropTableStatement, error) {
	stmt := &DropTableStatement{TableName: tableName}

	if _, err := govalidator.ValidateStruct(stmt); err != nil {
		return nil, err
	}

	return stmt, nil
}

func (d DropTableStatement) Protocol() protocol.MessageType {
	return protocol.DropTable
}

func (d DropTableStatement) ToBytes() ([]byte, error) {
	return msgpack.Marshal(d)
}

func (d *DropTableStatement) FromBytes(data []byte) error {
	return msgpack.Unmarshal(data, d)
}

func (d DropTableStatement) String() string {
	return fmt.Sprintf("DropTableStatement{TableName: %s}", d.TableName)
}
