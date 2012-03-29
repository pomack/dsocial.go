package apiutil

import "errors"

const (
    SIGNATURE_METHOD_HMAC_SHA256   = "HmacSHA256"
    SIGNATURE_METHOD_HMAC_SHA1     = "HmacSHA1"
    MAX_VALID_TIMESTAMP_IN_SECONDS = 5 * 60

    DEFAULT_SIGNATURE_VERSION = "1"
    DEFAULT_SIGNATURE_METHOD  = SIGNATURE_METHOD_HMAC_SHA256
)

var ErrorInvalidURI = errors.New("Could not parse the URI")
var ErrorInvalidAccessKeyId = errors.New("Could not find the specified DSOCAccessKeyId")
var ErrorSignatureDoesNotMatch = errors.New("Signature does not match")
var ErrorTimestampTooOld = errors.New("Timestamp too old")
var ErrorRequestExpired = errors.New("Request Expired")
var ErrorExpiresOrTimestampRequired = errors.New("Expires or Timestamp Required")
var ErrorInvalidSignatureVersion = errors.New("Invalid SignatureVersion")
var ErrorInvalidSignatureMethod = errors.New("Invalid SignatureMethod")
