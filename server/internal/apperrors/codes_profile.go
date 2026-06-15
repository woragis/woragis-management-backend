package apperrors

const (
	CodeProfileGetV1ServiceNotFound = "PROFILE_GET_V1_SERVICE_NOT_FOUND"
	MsgProfileGetV1ServiceNotFound  = "Profile not found."

	CodeProfileGetV1ServiceLoadFailed = "PROFILE_GET_V1_SERVICE_LOAD_FAILED"
	MsgProfileGetV1ServiceLoadFailed  = "Failed to load profile."

	CodeProfilePatchV1ServiceDisplayNameEmpty = "PROFILE_PATCH_V1_SERVICE_DISPLAY_NAME_EMPTY"
	MsgProfilePatchV1ServiceDisplayNameEmpty  = "Display name is required."

	CodeProfilePatchV1ServiceUpdateFailed = "PROFILE_PATCH_V1_SERVICE_UPDATE_FAILED"
	MsgProfilePatchV1ServiceUpdateFailed  = "Failed to update profile."
)
