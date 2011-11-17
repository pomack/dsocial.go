package dsocial

import (
    "os"
)

var (
    ERR_MUST_SPECIFY_ID os.Error
    ERR_INVALID_ID os.Error
    ERR_INVALID_EMAIL_ADDRESS os.Error
    ERR_REQUIRED_FIELD os.Error
    ERR_INVALID_FORMAT os.Error
)

func init() {
    ERR_MUST_SPECIFY_ID = os.NewError("Must Specify Id")
    ERR_INVALID_ID = os.NewError("Invalid Id")
    ERR_INVALID_EMAIL_ADDRESS = os.NewError("Invalid Email Address")
    ERR_REQUIRED_FIELD = os.NewError("Required Field")
    ERR_INVALID_FORMAT = os.NewError("Invalid Format")
}

