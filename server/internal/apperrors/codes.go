package apperrors

const (
	CodeInternal = "INTERNAL_ERROR"
	MsgInternal  = "An unexpected error occurred."

	CodeReadyGetHandlerSqlGetterFailed = "READY_GET_HANDLER_SQL_GETTER_FAILED"
	MsgReadyGetHandlerSqlGetterFailed  = "Database connection is unavailable."

	CodeReadyGetHandlerDatabasePingFailed = "READY_GET_HANDLER_DATABASE_PING_FAILED"
	MsgReadyGetHandlerDatabasePingFailed  = "Database is not ready."

	CodeAdminAuthV1HandlerKeyMissing = "ADMIN_AUTH_V1_HANDLER_KEY_MISSING"
	MsgAdminAuthV1HandlerKeyMissing  = "Admin API key is required."

	CodeAdminAuthV1HandlerKeyInvalid = "ADMIN_AUTH_V1_HANDLER_KEY_INVALID"
	MsgAdminAuthV1HandlerKeyInvalid  = "Admin API key is invalid."
)
