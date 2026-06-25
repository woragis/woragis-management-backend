package apperrors

const (
	CodeProjectListV1ServiceLoadFailed = "PROJECT_LIST_V1_SERVICE_LOAD_FAILED"
	MsgProjectListV1ServiceLoadFailed  = "Failed to load projects."

	CodeProjectGetV1HandlerPathIDInvalid = "PROJECT_GET_V1_HANDLER_PATH_ID_INVALID"
	MsgProjectGetV1HandlerPathIDInvalid  = "Project id is invalid."

	CodeProjectGetV1ServiceNotFound = "PROJECT_GET_V1_SERVICE_NOT_FOUND"
	MsgProjectGetV1ServiceNotFound  = "Project not found."

	CodeProjectGetV1ServiceLoadFailed = "PROJECT_GET_V1_SERVICE_LOAD_FAILED"
	MsgProjectGetV1ServiceLoadFailed  = "Failed to load project."

	CodeProjectPostV1ServiceNameEmpty = "PROJECT_POST_V1_SERVICE_NAME_EMPTY"
	MsgProjectPostV1ServiceNameEmpty  = "Project name is required."

	CodeProjectPostV1ServiceSlugInvalid = "PROJECT_POST_V1_SERVICE_SLUG_INVALID"
	MsgProjectPostV1ServiceSlugInvalid  = "Project slug is invalid."

	CodeProjectPostV1ServiceCreateFailed = "PROJECT_POST_V1_SERVICE_CREATE_FAILED"
	MsgProjectPostV1ServiceCreateFailed  = "Failed to create project."

	CodeProjectPatchV1ServiceNotFound = "PROJECT_PATCH_V1_SERVICE_NOT_FOUND"
	MsgProjectPatchV1ServiceNotFound  = "Project not found."

	CodeProjectPatchV1ServiceUpdateFailed = "PROJECT_PATCH_V1_SERVICE_UPDATE_FAILED"
	MsgProjectPatchV1ServiceUpdateFailed  = "Failed to update project."

	CodeProjectDeleteV1ServiceNotFound = "PROJECT_DELETE_V1_SERVICE_NOT_FOUND"
	MsgProjectDeleteV1ServiceNotFound  = "Project not found."

	CodeProjectDeleteV1ServiceDeleteFailed = "PROJECT_DELETE_V1_SERVICE_DELETE_FAILED"
	MsgProjectDeleteV1ServiceDeleteFailed  = "Failed to delete project."

	CodeProjectLinkPostV1ServiceProjectNotFound = "PROJECT_LINK_POST_V1_SERVICE_PROJECT_NOT_FOUND"
	MsgProjectLinkPostV1ServiceProjectNotFound  = "Project not found."

	CodeProjectLinkPostV1ServiceURLInvalid = "PROJECT_LINK_POST_V1_SERVICE_URL_INVALID"
	MsgProjectLinkPostV1ServiceURLInvalid  = "Link url is required."

	CodeProjectLinkPostV1ServiceCreateFailed = "PROJECT_LINK_POST_V1_SERVICE_CREATE_FAILED"
	MsgProjectLinkPostV1ServiceCreateFailed  = "Failed to create project link."

	CodeProjectLinkDeleteV1ServiceNotFound = "PROJECT_LINK_DELETE_V1_SERVICE_NOT_FOUND"
	MsgProjectLinkDeleteV1ServiceNotFound  = "Project link not found."

	CodeProjectDomainPostV1ServiceDomainInvalid = "PROJECT_DOMAIN_POST_V1_SERVICE_DOMAIN_INVALID"
	MsgProjectDomainPostV1ServiceDomainInvalid  = "Domain is required."

	CodeProjectDomainDeleteV1ServiceNotFound = "PROJECT_DOMAIN_DELETE_V1_SERVICE_NOT_FOUND"
	MsgProjectDomainDeleteV1ServiceNotFound  = "Project domain not found."

	CodeProjectSecretPostV1ServiceNameEmpty = "PROJECT_SECRET_POST_V1_SERVICE_NAME_EMPTY"
	MsgProjectSecretPostV1ServiceNameEmpty  = "Secret name is required."

	CodeProjectSecretPostV1ServiceValueEmpty = "PROJECT_SECRET_POST_V1_SERVICE_VALUE_EMPTY"
	MsgProjectSecretPostV1ServiceValueEmpty  = "Secret value is required."

	CodeProjectSecretPostV1ServiceEncryptFailed = "PROJECT_SECRET_POST_V1_SERVICE_ENCRYPT_FAILED"
	MsgProjectSecretPostV1ServiceEncryptFailed  = "Failed to encrypt secret."

	CodeProjectSecretGetV1ServiceNotFound = "PROJECT_SECRET_GET_V1_SERVICE_NOT_FOUND"
	MsgProjectSecretGetV1ServiceNotFound  = "Project secret not found."

	CodeProjectSecretGetV1ServiceDecryptFailed = "PROJECT_SECRET_GET_V1_SERVICE_DECRYPT_FAILED"
	MsgProjectSecretGetV1ServiceDecryptFailed  = "Failed to decrypt secret."

	CodeProjectGalleryPostV1ServiceMediaInvalid = "PROJECT_GALLERY_POST_V1_SERVICE_MEDIA_INVALID"
	MsgProjectGalleryPostV1ServiceMediaInvalid  = "Media asset id is invalid."

	CodeProjectGalleryDeleteV1ServiceNotFound = "PROJECT_GALLERY_DELETE_V1_SERVICE_NOT_FOUND"
	MsgProjectGalleryDeleteV1ServiceNotFound  = "Gallery item not found."

	CodeProjectEnvPostV1ServiceKeyEmpty = "PROJECT_ENV_POST_V1_SERVICE_KEY_EMPTY"
	MsgProjectEnvPostV1ServiceKeyEmpty  = "Environment variable key is required."

	CodeProjectEnvPostV1ServiceCreateFailed = "PROJECT_ENV_POST_V1_SERVICE_CREATE_FAILED"
	MsgProjectEnvPostV1ServiceCreateFailed  = "Failed to create environment variable."

	CodeProjectEnvDeleteV1ServiceNotFound = "PROJECT_ENV_DELETE_V1_SERVICE_NOT_FOUND"
	MsgProjectEnvDeleteV1ServiceNotFound  = "Environment variable not found."

	CodeProjectSecretUnlockRequired = "PROJECT_SECRET_UNLOCK_REQUIRED"
	MsgProjectSecretUnlockRequired  = "Secret unlock password is required to change access from secret."

	CodeProjectSecretUnlockInvalid = "PROJECT_SECRET_UNLOCK_INVALID"
	MsgProjectSecretUnlockInvalid  = "Invalid secret unlock password."

	CodeProjectSecretUnlockUnavailable = "PROJECT_SECRET_UNLOCK_UNAVAILABLE"
	MsgProjectSecretUnlockUnavailable  = "Secret unlock is not configured on the server."

	CodeProjectSecretFeaturedBlocked = "PROJECT_SECRET_FEATURED_BLOCKED"
	MsgProjectSecretFeaturedBlocked  = "Secret projects cannot be featured."

	CodeProjectGalleryPostV1ServiceMediaTypeInvalid = "PROJECT_GALLERY_POST_V1_SERVICE_MEDIA_TYPE_INVALID"
	MsgProjectGalleryPostV1ServiceMediaTypeInvalid  = "Gallery only supports images, GIFs, and videos."
)
