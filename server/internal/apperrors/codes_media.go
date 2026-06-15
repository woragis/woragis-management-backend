package apperrors

const (
	CodeMediaListV1ServiceLoadFailed = "MEDIA_LIST_V1_SERVICE_LOAD_FAILED"
	MsgMediaListV1ServiceLoadFailed  = "Failed to load media assets."

	CodeMediaGetV1HandlerPathIDInvalid = "MEDIA_GET_V1_HANDLER_PATH_ID_INVALID"
	MsgMediaGetV1HandlerPathIDInvalid  = "Media id is invalid."

	CodeMediaGetV1ServiceNotFound = "MEDIA_GET_V1_SERVICE_NOT_FOUND"
	MsgMediaGetV1ServiceNotFound  = "Media asset not found."

	CodeMediaPostV1ServiceFileMissing = "MEDIA_POST_V1_SERVICE_FILE_MISSING"
	MsgMediaPostV1ServiceFileMissing  = "File is required."

	CodeMediaPostV1ServiceFileTooLarge = "MEDIA_POST_V1_SERVICE_FILE_TOO_LARGE"
	MsgMediaPostV1ServiceFileTooLarge  = "File exceeds maximum size."

	CodeMediaPostV1ServiceCreateFailed = "MEDIA_POST_V1_SERVICE_CREATE_FAILED"
	MsgMediaPostV1ServiceCreateFailed  = "Failed to store media asset."

	CodeMediaDeleteV1ServiceNotFound = "MEDIA_DELETE_V1_SERVICE_NOT_FOUND"
	MsgMediaDeleteV1ServiceNotFound  = "Media asset not found."
)
