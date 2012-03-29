package dsocial

import "errors"

var (
    ERR_MUST_SPECIFY_ID       error
    ERR_INVALID_ID            error
    ERR_INVALID_EMAIL_ADDRESS error
    ERR_REQUIRED_FIELD        error
    ERR_INVALID_FORMAT        error
)

func init() {
    ERR_MUST_SPECIFY_ID = errors.New("Must Specify Id")
    ERR_INVALID_ID = errors.New("Invalid Id")
    ERR_INVALID_EMAIL_ADDRESS = errors.New("Invalid Email Address")
    ERR_REQUIRED_FIELD = errors.New("Required Field")
    ERR_INVALID_FORMAT = errors.New("Invalid Format")
}
