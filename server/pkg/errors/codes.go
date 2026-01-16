package errors

// Error code format: SERVICE_CATEGORY_NUMBER
// Each code should be used in exactly one place for easy tracking

const (
	// AUTH - Authentication errors (1000-1099)
	AUTH_JWT_INVALID_SIGNATURE    = "MGMT_1001"
	AUTH_JWT_EXPIRED              = "MGMT_1002"
	AUTH_JWT_MISSING_CLAIMS       = "MGMT_1003"
	AUTH_JWT_MALFORMED            = "MGMT_1004"
	AUTH_TOKEN_MISSING            = "MGMT_1020"
	AUTH_TOKEN_INVALID_FORMAT     = "MGMT_1021"
	AUTH_UNAUTHORIZED             = "MGMT_1022"
	AUTH_SERVICE_UNAVAILABLE      = "MGMT_1030"
	AUTH_SERVICE_ERROR            = "MGMT_1031"

	// CSRF - CSRF token errors (2000-2099)
	CSRF_TOKEN_EXPIRED            = "MGMT_2001"
	CSRF_TOKEN_INVALID            = "MGMT_2002"
	CSRF_TOKEN_MISSING            = "MGMT_2003"
	CSRF_TOKEN_GENERATION_FAILED  = "MGMT_2004"
	CSRF_TOKEN_MISMATCH           = "MGMT_2005"

	// DB - Database errors (3000-3099)
	DB_CONNECTION_FAILED          = "MGMT_3001"
	DB_QUERY_FAILED               = "MGMT_3002"
	DB_TRANSACTION_FAILED         = "MGMT_3003"
	DB_RECORD_NOT_FOUND           = "MGMT_3004"
	DB_DUPLICATE_ENTRY            = "MGMT_3005"
	DB_CONSTRAINT_VIOLATION       = "MGMT_3006"

	// VALIDATION - Input validation errors (4000-4099)
	VALIDATION_INVALID_INPUT      = "MGMT_4001"
	VALIDATION_MISSING_FIELD      = "MGMT_4002"
	VALIDATION_FIELD_TOO_LONG     = "MGMT_4003"
	VALIDATION_FIELD_TOO_SHORT    = "MGMT_4004"
	VALIDATION_INVALID_UUID       = "MGMT_4005"
	VALIDATION_INVALID_DATE       = "MGMT_4006"
	VALIDATION_INVALID_URL        = "MGMT_4007"
	VALIDATION_INVALID_EMAIL      = "MGMT_4008"
	VALIDATION_INVALID_SLUG       = "MGMT_4009"

	// REDIS - Redis/Cache errors (5000-5099)
	REDIS_CONNECTION_FAILED       = "MGMT_5001"
	REDIS_GET_FAILED              = "MGMT_5002"
	REDIS_SET_FAILED              = "MGMT_5003"
	REDIS_DELETE_FAILED           = "MGMT_5004"

	// PROJECT - Project specific errors (6000-6099)
	PROJECT_NOT_FOUND             = "MGMT_6001"
	PROJECT_SLUG_EXISTS           = "MGMT_6002"
	PROJECT_UPDATE_FAILED         = "MGMT_6003"
	PROJECT_DELETE_FAILED         = "MGMT_6004"

	// EXPERIENCE - Work experience errors (6100-6199)
	EXPERIENCE_NOT_FOUND          = "MGMT_6101"
	EXPERIENCE_EXISTS             = "MGMT_6102"
	EXPERIENCE_INVALID_DATE_RANGE = "MGMT_6103"
	EXPERIENCE_UPDATE_FAILED      = "MGMT_6104"

	// LANGUAGE - Language errors (6200-6299)
	LANGUAGE_NOT_FOUND            = "MGMT_6201"
	LANGUAGE_EXISTS               = "MGMT_6202"
	LANGUAGE_INVALID_LEVEL        = "MGMT_6203"

	// CLIENT - Client management errors (6300-6399)
	CLIENT_NOT_FOUND              = "MGMT_6301"
	CLIENT_EXISTS                 = "MGMT_6302"
	CLIENT_UPDATE_FAILED          = "MGMT_6303"

	// FINANCE - Financial record errors (6400-6499)
	FINANCE_NOT_FOUND             = "MGMT_6401"
	FINANCE_INVALID_AMOUNT        = "MGMT_6402"
	FINANCE_UPDATE_FAILED         = "MGMT_6403"

	// API_KEY - API key management errors (6500-6599)
	API_KEY_NOT_FOUND             = "MGMT_6501"
	API_KEY_INVALID               = "MGMT_6502"
	API_KEY_EXPIRED               = "MGMT_6503"
	API_KEY_GENERATION_FAILED     = "MGMT_6504"
	API_KEY_REVOKED               = "MGMT_6505"

	// CHAT - Chat/AI integration errors (6600-6699)
	CHAT_INVALID_MESSAGE          = "MGMT_6601"
	CHAT_CONTEXT_TOO_LARGE        = "MGMT_6602"
	CHAT_AI_SERVICE_ERROR         = "MGMT_6603"

	// TESTIMONIAL - Testimonial errors (6700-6799)
	TESTIMONIAL_NOT_FOUND         = "MGMT_6701"
	TESTIMONIAL_EXISTS            = "MGMT_6702"

	// AI_SERVICE - AI service errors (8000-8099)
	AI_SERVICE_UNAVAILABLE        = "MGMT_8001"
	AI_SERVICE_REQUEST_FAILED     = "MGMT_8002"
	AI_SERVICE_INVALID_RESPONSE   = "MGMT_8003"
	AI_SERVICE_TIMEOUT            = "MGMT_8004"

	// SERVER - Server/System errors (9000-9099)
	SERVER_INTERNAL_ERROR         = "MGMT_9001"
	SERVER_SERVICE_UNAVAILABLE    = "MGMT_9002"
	SERVER_TIMEOUT                = "MGMT_9003"
	SERVER_CONTEXT_CANCELLED      = "MGMT_9004"
	SERVER_RATE_LIMIT_EXCEEDED    = "MGMT_9005"
)

// Error messages - human-readable descriptions
var errorMessages = map[string]string{
	// Authentication
	AUTH_JWT_INVALID_SIGNATURE: "JWT token signature is invalid",
	AUTH_JWT_EXPIRED:           "JWT token has expired",
	AUTH_JWT_MISSING_CLAIMS:    "JWT token is missing required claims",
	AUTH_JWT_MALFORMED:         "JWT token is malformed",
	AUTH_TOKEN_MISSING:         "Authentication token is missing",
	AUTH_TOKEN_INVALID_FORMAT:  "Authentication token has invalid format",
	AUTH_UNAUTHORIZED:          "Unauthorized access",
	AUTH_SERVICE_UNAVAILABLE:   "Authentication service is unavailable",
	AUTH_SERVICE_ERROR:         "Authentication service error",

	// CSRF
	CSRF_TOKEN_EXPIRED:           "CSRF token has expired",
	CSRF_TOKEN_INVALID:           "CSRF token validation failed",
	CSRF_TOKEN_MISSING:           "CSRF token is missing from request",
	CSRF_TOKEN_GENERATION_FAILED: "Failed to generate CSRF token",
	CSRF_TOKEN_MISMATCH:          "CSRF token does not match stored value",

	// Database
	DB_CONNECTION_FAILED:    "Failed to connect to database",
	DB_QUERY_FAILED:         "Database query execution failed",
	DB_TRANSACTION_FAILED:   "Database transaction failed",
	DB_RECORD_NOT_FOUND:     "Requested record not found",
	DB_DUPLICATE_ENTRY:      "Record already exists",
	DB_CONSTRAINT_VIOLATION: "Database constraint violation",

	// Validation
	VALIDATION_INVALID_INPUT:  "Input validation failed",
	VALIDATION_MISSING_FIELD:  "Required field is missing",
	VALIDATION_FIELD_TOO_LONG: "Field value exceeds maximum length",
	VALIDATION_FIELD_TOO_SHORT: "Field value is below minimum length",
	VALIDATION_INVALID_UUID:   "Invalid UUID format",
	VALIDATION_INVALID_DATE:   "Invalid date format",
	VALIDATION_INVALID_URL:    "Invalid URL format",
	VALIDATION_INVALID_EMAIL:  "Invalid email format",
	VALIDATION_INVALID_SLUG:   "Invalid slug format",

	// Redis
	REDIS_CONNECTION_FAILED: "Failed to connect to Redis",
	REDIS_GET_FAILED:        "Failed to retrieve data from cache",
	REDIS_SET_FAILED:        "Failed to store data in cache",
	REDIS_DELETE_FAILED:     "Failed to delete data from cache",

	// Projects
	PROJECT_NOT_FOUND:     "Project not found",
	PROJECT_SLUG_EXISTS:   "Project with this slug already exists",
	PROJECT_UPDATE_FAILED: "Failed to update project",
	PROJECT_DELETE_FAILED: "Failed to delete project",

	// Experience
	EXPERIENCE_NOT_FOUND:          "Work experience not found",
	EXPERIENCE_EXISTS:             "Work experience already exists",
	EXPERIENCE_INVALID_DATE_RANGE: "Invalid date range for work experience",
	EXPERIENCE_UPDATE_FAILED:      "Failed to update work experience",

	// Language
	LANGUAGE_NOT_FOUND:     "Language not found",
	LANGUAGE_EXISTS:        "Language already exists",
	LANGUAGE_INVALID_LEVEL: "Invalid language proficiency level",

	// Client
	CLIENT_NOT_FOUND:     "Client not found",
	CLIENT_EXISTS:        "Client already exists",
	CLIENT_UPDATE_FAILED: "Failed to update client",

	// Finance
	FINANCE_NOT_FOUND:      "Financial record not found",
	FINANCE_INVALID_AMOUNT: "Invalid financial amount",
	FINANCE_UPDATE_FAILED:  "Failed to update financial record",

	// API Key
	API_KEY_NOT_FOUND:         "API key not found",
	API_KEY_INVALID:           "API key is invalid",
	API_KEY_EXPIRED:           "API key has expired",
	API_KEY_GENERATION_FAILED: "Failed to generate API key",
	API_KEY_REVOKED:           "API key has been revoked",

	// Chat
	CHAT_INVALID_MESSAGE:   "Invalid chat message format",
	CHAT_CONTEXT_TOO_LARGE: "Chat context exceeds maximum size",
	CHAT_AI_SERVICE_ERROR:  "AI service error in chat",

	// Testimonial
	TESTIMONIAL_NOT_FOUND: "Testimonial not found",
	TESTIMONIAL_EXISTS:    "Testimonial already exists",

	// AI Service
	AI_SERVICE_UNAVAILABLE:     "AI service is unavailable",
	AI_SERVICE_REQUEST_FAILED:  "AI service request failed",
	AI_SERVICE_INVALID_RESPONSE: "AI service returned invalid response",
	AI_SERVICE_TIMEOUT:         "AI service request timeout",

	// Server
	SERVER_INTERNAL_ERROR:      "Internal server error occurred",
	SERVER_SERVICE_UNAVAILABLE: "Service is temporarily unavailable",
	SERVER_TIMEOUT:             "Request timeout",
	SERVER_CONTEXT_CANCELLED:   "Request was cancelled",
	SERVER_RATE_LIMIT_EXCEEDED: "Rate limit exceeded",
}

// GetMessage returns the human-readable message for an error code
func GetMessage(code string) string {
	if msg, ok := errorMessages[code]; ok {
		return msg
	}
	return "Unknown error occurred"
}
