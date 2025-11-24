package errx

var (
	ErrInvalid = New(BadRequest, "The request body contains invalid or malformed JSON.")

	ErrLock      = New(BadRequest, "Resource is locked. Another process is running, please wait.")
	ErrColor     = New(BadRequest, "Hex color must be a valid string.")
	ErrUuid      = New(BadRequest, "The id must be a valid uuid.")
	ErrCategory  = New(BadRequest, "Category doesn't exists.")
	ErrTag       = New(BadRequest, "Tag doesn't exists.")
	ErrBitmask   = New(BadRequest, "Invalid bitmask value.")
	ErrRole      = New(BadRequest, "Role doesn't exists.")
	ErrPosition  = New(BadRequest, "Invalid position.")
	ErrTimezone  = New(BadRequest, "Timezone doesn't exists.")
	ErrTime      = New(BadRequest, "Invalid time format.")
	ErrNotEnough = New(BadRequest, "Not enough data to perform this action.")
	ErrJSONKey   = New(BadRequest, "JSON key can only contain characters and underscore.")

	// Authorization
	ErrToken = New(BadRequest, "Invalid or expired token.")
	ErrAuth  = New(BadRequest, "Missing or invalid Authorization header.")

	ErrUser        = New(BadRequest, "User doesn't exists.")
	ErrPassword    = New(BadRequest, "The password must be between 8 and 50 characters long and contain both uppercase and lowercase letters, as well as a number.")
	ErrEmail       = New(BadRequest, "Invalid email address.")
	ErrCredentials = New(BadRequest, "Invalid email or password.")
	ErrSession     = New(BadRequest, "Invalid or expired session.")
	ErrCodeLimit   = New(BadRequest, "Too many attempts. Start a new session and try again later.")
	ErrCode        = New(BadRequest, "Invalid or expired verification code.")
	ErrAuthLimit   = New(BadRequest, "Too many attempts, please try again later.")

	// Captch
	ErrCaptcha = New(BadRequest, "We couldn’t verify you’re human. Please try the security check again or reload the page.")
)
