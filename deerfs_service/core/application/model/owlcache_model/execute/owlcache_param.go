package execute

//response status type
type ResStatus int

//request command type
type CommandType string

//response 状态
const (
	SUCCESS         ResStatus = 200
	ERROR           ResStatus = 500
	NOT_FOUND       ResStatus = 404
	UNKNOWN_COMMAND ResStatus = 501
	NOT_PASS        ResStatus = 401
)

const (
	GET    CommandType = "get"
	EXIST  CommandType = "exist"
	SET    CommandType = "set"
	EXPIRE CommandType = "expire"
	DELETE CommandType = "delete"
	PASS   CommandType = "pass"
	PING   CommandType = "ping"
)

//status to string
func ResStatusToString(resstatus ResStatus) string {

	var s string
	switch resstatus {
	case SUCCESS:
		s = "SUCCESS"
	case ERROR:
		s = "ERROR"
	case NOT_FOUND:
		s = "NOT_FOUND"
	case UNKNOWN_COMMAND:
		s = "UNKNOWN_COMMAND"
	case NOT_PASS:
		s = "NOT_PASS"
	}

	return s
}
