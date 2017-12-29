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
)

// messages
const (
	MessLoggedIn        = "login_success"
	MessLoggedOut       = "logged_out"
	MessConfirmMailSent = "confirm_mail_sent"
	MessRestoreMailSent = "restore_mail_sent"
)
