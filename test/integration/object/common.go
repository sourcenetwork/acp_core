package object

import "github.com/sourcenetwork/acp_core/test/util"

var timestamp = util.MustDateTimeToProto("2024-01-01 00:00:00")

var metadata map[string]string = map[string]string{
	"test": "1234",
}
