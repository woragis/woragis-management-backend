package apperrors

const (
	CodeContentVideoCreateInvalid = "CONTENT_LEETCODE_VIDEO_CREATE_V1_INVALID"
	MsgContentVideoCreateInvalid  = "Invalid leetcode video payload"

	CodeContentVideoNotFound = "CONTENT_LEETCODE_VIDEO_GET_V1_NOT_FOUND"
	MsgContentVideoNotFound  = "Leetcode video not found"

	CodeContentThumbnailCreateInvalid = "CONTENT_LEETCODE_THUMBNAIL_CREATE_V1_INVALID"
	MsgContentThumbnailCreateInvalid  = "Invalid thumbnail payload"

	CodeContentThumbnailNotFound = "CONTENT_LEETCODE_THUMBNAIL_GET_V1_NOT_FOUND"
	MsgContentThumbnailNotFound  = "Thumbnail not found"

	CodeContentThumbnailGenerateInvalid = "CONTENT_LEETCODE_THUMBNAIL_GENERATE_V1_INVALID"
	MsgContentThumbnailGenerateInvalid  = "Thumbnail cannot be generated in its current state"

	CodeContentTemplateCreateInvalid = "CONTENT_LEETCODE_TEMPLATE_CREATE_V1_INVALID"
	MsgContentTemplateCreateInvalid  = "Invalid prompt template payload"

	CodeContentTemplateNotFound = "CONTENT_LEETCODE_TEMPLATE_GET_V1_NOT_FOUND"
	MsgContentTemplateNotFound  = "Prompt template not found"

	CodeCreativesClientUnavailable = "CREATIVES_CLIENT_UNAVAILABLE"
	MsgCreativesClientUnavailable  = "Creatives API is not configured"

	CodeWhatsappWorkerUnavailable = "WHATSAPP_WORKER_UNAVAILABLE"
	MsgWhatsappWorkerUnavailable  = "WhatsApp worker is not configured"

	CodeContentSettingsNotFound = "CONTENT_LEETCODE_SETTINGS_GET_V1_NOT_FOUND"
	MsgContentSettingsNotFound  = "LeetCode channel settings not found"

	CodeContentWhatsappTplNotFound = "CONTENT_WHATSAPP_TEMPLATE_GET_V1_NOT_FOUND"
	MsgContentWhatsappTplNotFound  = "WhatsApp template not found"

	CodeContentDispatchEmpty = "CONTENT_LEETCODE_DISPATCH_V1_EMPTY"
	MsgContentDispatchEmpty  = "No message to dispatch for this request"

	CodeWorkerAuthInvalid = "WORKER_AUTH_V1_INVALID"
	MsgWorkerAuthInvalid  = "Invalid worker API key"
)
