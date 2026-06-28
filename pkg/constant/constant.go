package constant

// Context keys for storing values in context
const (
	USER_CTX    = "user"
	USER_ID_CTX = "user_id"
)

// User roles
const (
	ROLE_USER  = "user"
	ROLE_ADMIN = "admin"
)

// Response messages
const (
	MSG_SUCCESS              = "SUCCESS"
	MSG_CREATED              = "CREATED"
	MSG_UPDATED              = "UPDATED"
	MSG_DELETED              = "DELETED"
	MSG_INTERNAL_SERVER_ERROR = "INTERNAL_SERVER_ERROR"
	MSG_BAD_REQUEST          = "BAD_REQUEST"
	MSG_UNAUTHORIZED         = "UNAUTHORIZED"
	MSG_FORBIDDEN            = "FORBIDDEN"
	MSG_NOT_FOUND            = "NOT_FOUND"
)
