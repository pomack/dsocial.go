package apiutil

import (
    "os"
)

const (
    SIGNATURE_METHOD_HMAC_SHA256 = "HmacSHA256"
    SIGNATURE_METHOD_HMAC_SHA1 = "HmacSHA1"
    MAX_VALID_TIMESTAMP_IN_SECONDS = 5 * 60
    
    DEFAULT_SIGNATURE_VERSION = "1"
    DEFAULT_SIGNATURE_METHOD = SIGNATURE_METHOD_HMAC_SHA256
)

var ErrorInvalidURI os.Error = os.NewError("Could not parse the URI")
var ErrorSignatureDoesNotMatch os.Error = os.NewError("Signature does not match")
var ErrorTimestampTooOld os.Error = os.NewError("Timestamp too old")
var ErrorRequestExpired os.Error = os.NewError("Request Expired")
var ErrorExpiresOrTimestampRequired os.Error = os.NewError("Expires or Timestamp Required")
var ErrorInvalidSignatureVersion os.Error = os.NewError("Invalid SignatureVersion")
var ErrorInvalidSignatureMethod os.Error = os.NewError("Invalid SignatureMethod")

