package routes

// errors
const (
	ErrDB               = "database_error"
	ErrUnexpected       = "unexpected_error"
	ErrWrongCredentials = "wrong_credentials_error"
	ErrNotConfirmed     = "user_not_confirmed"
	ErrUnequalPasswords = "passwords_unequal"
	ErrUserExists       = "user_exists"
	ErrMail             = "mail_error"
	ErrWrongGrade       = "wrong_grade_error"
	ErrIllegalInput     = "illegal_input_error"
	ErrFS               = "filesystem_error"
	ErrIllegalDate      = "illegal_date_error"
	ErrNotReady         = "server_not_ready_error"
	ErrNoAssets         = "no_assets_error"
)

// messages
const (
	MessLoggedIn        = "login_success"
	MessLoggedOut       = "logged_out"
	MessConfirmMailSent = "confirm_mail_sent"
	MessRestoreMailSent = "restore_mail_sent"
	MessChangedPassword = "password_changed"
)
