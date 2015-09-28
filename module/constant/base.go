package constant

const (
	AUTH_TYPE_PUBLIC = 0
	AUTH_TYPE_LOGIN  = 1
	AUTH_TYPE_AUTH   = 2

	PERMISSION_STATUS_OK            = 0
	PERMISSION_STATUS_NEED_LOGIN    = 1
	PERMISSION_STATUS_NO_PERMISSION = 2
)

const (
	CTX_KEY_USER = "ctx-user"

	// session key
	SS_KEY_USER = "ss-user"
)
