package apiutil

// borrowed heavily from https://github.com/abneptis/GoAWS/blob/master/signer.go
// which is released under the Go language license
// see https://github.com/abneptis/GoAWS/blob/master/LICENSE

import (
    dm "github.com/pomack/dsocial.go/models/dsocial"
    "bytes"
    "crypto"
    "crypto/hmac"
    "encoding/base64"
    "hash"
    "http"
    "os"
    "strings"
    "strconv"
    "time"
    "url"
)

// A signer simply holds the access & secret access keys
// necessary for aws, and proivides helper functions
// to assist in generating an appropriate signature.
type Signer interface {
    AccessKey() string
    SignBytes(h crypto.Hash, buf []byte) (signature []byte, err os.Error)
    SignString(h crypto.Hash, s string) (signature string, err os.Error)
    SignEncoded(h crypto.Hash, s string, enc *base64.Encoding) (signature []byte, err os.Error)
    SignRequest(req *http.Request, expiresIn int64)
}

type signer struct {
    accessKey string
    secretAccessKey []byte
}

func NewSigner(accessKey string, secretAccessKey string) Signer {
    return &signer{
        accessKey: accessKey,
        secretAccessKey: bytes.NewBufferString(secretAccessKey).Bytes(),
    }
}

func (p *signer) AccessKey() string {
    return p.accessKey
}

// the core function of the Signer, generates the raw hmac of he bytes.
func (p *signer) SignBytes(h crypto.Hash, buf []byte) (signature []byte, err os.Error) {
    hasher := hmac.New(func() hash.Hash { return h.New() }, p.secretAccessKey)
	_, err = hasher.Write(buf)
	if err == nil {
		signature = hasher.Sum()
	}
	return
}

// Same as SignBytes, but with strings.
func (p *signer) SignString(h crypto.Hash, s string) (signature string, err os.Error) {
    buf, err := p.SignBytes(h, bytes.NewBufferString(s).Bytes())
	if err == nil {
		signature = string(buf)
	}
	return
}

// SignBytes, but will base64 encode based on the specified encoder.
func (p *signer) SignEncoded(h crypto.Hash, s string, enc *base64.Encoding) (signature []byte, err os.Error) {
	buf, err := p.SignBytes(h, bytes.NewBufferString(s).Bytes())
	if err == nil {
		signature = make([]byte, enc.EncodedLen(len(buf)))
		enc.Encode(signature, buf)
	}
	return
}

// Modifies the request for signing
// if expiresIn is set to 0, a Timestamp will be used, otherwise an expiration.
func (p *signer) SignRequest(req *http.Request, expiresIn int64) {
    qstring, err := url.ParseQuery(req.URL.RawQuery)
    if err != nil { return }
    qstring["SignatureVersion"] = []string{DEFAULT_SIGNATURE_VERSION}
    if _, ok := qstring["SignatureMethod"]; !ok || len(qstring["SignatureMethod"]) == 0 {
        qstring["SignatureMethod"] = []string{DEFAULT_SIGNATURE_METHOD}
    }
    if expiresIn > 0 {
        qstring["Expires"] = []string{strconv.Itoa64(time.Seconds()+expiresIn)}
    } else {
        qstring["Timestamp"] = []string{time.UTC().Format(dm.UTC_DATETIME_FORMAT)}
    }
    qstring["Signature"] = nil, false
    qstring["DSOCAccessKeyId"] = []string{p.accessKey}

  	var signature []byte
    req.URL.RawQuery = qstring.Encode()
    canonicalizedStringToSign, err := p.Canonicalize(req)
    if err != nil { return }
    //log.Printf("String-to-sign: '%s'", canonicalizedStringToSign)

  	switch qstring["SignatureMethod"][0] {
  	case SIGNATURE_METHOD_HMAC_SHA256:
  		signature, err = p.SignEncoded(crypto.SHA256, canonicalizedStringToSign, base64.StdEncoding)
  	case SIGNATURE_METHOD_HMAC_SHA1:
  		signature, err = p.SignEncoded(crypto.SHA1, canonicalizedStringToSign, base64.StdEncoding)
  	default:
  		err = os.NewError("Unknown SignatureMethod:" + req.Form.Get("SignatureMethod"))
  	}

  	if err == nil {
        req.URL.RawQuery += "&" + url.Values{"Signature": []string{string(signature)}}.Encode()
        req.RawURL = req.URL.String()
  	}
  	return
}

// Generates the canonical string-to-sign for dsocial services.
// You shouldn't need to use this directly.
func (p *signer) Canonicalize(req *http.Request) (out string, err os.Error) {
    fv, err := url.ParseQuery(req.URL.RawQuery)
    if err == nil {
        out = strings.Join([]string{req.Method, req.Host, req.URL.Path, SortedEscape(fv)}, "\n")
    }
    return
}


// Checks whether the request has a signature, validates it if it does
// and returns an error if signature is present but not valid
func (p *signer) CheckSignature(req *http.Request) (hasSignature, validSignature bool, err os.Error) {
    qstring, err := url.ParseQuery(req.URL.RawQuery)
    if err != nil {
        err = ErrorInvalidURI
        return
    }
    if qstring.Get("Signature") == "" || qstring.Get("DSOCAccessKeyId") == "" {
        return
    }
    hasSignature = true
    now := time.UTC().Seconds()
    if expiresStr := qstring.Get("Expires"); expiresStr != "" {
        expiresAt, _ := strconv.Atoi64(expiresStr)
        if expiresAt < now {
            err = ErrorRequestExpired
            return
        }
    } else if timestampStr := qstring.Get("Timestamp"); timestampStr != "" {
        timestamp, _ := time.Parse(dm.UTC_DATETIME_FORMAT, timestampStr)
        if timestamp == nil || timestamp.Seconds() - MAX_VALID_TIMESTAMP_IN_SECONDS > now || timestamp.Seconds() + MAX_VALID_TIMESTAMP_IN_SECONDS < now {
            err = ErrorTimestampTooOld
            return
        }
    } else {
        err = ErrorExpiresOrTimestampRequired
        return
    }
    if qstring.Get("SignatureVersion") != "" && qstring.Get("SignatureVersion") != DEFAULT_SIGNATURE_VERSION {
        err = ErrorInvalidSignatureVersion
        return
    }
    var h crypto.Hash
    signatureMethod := qstring.Get("SignatureMethod")
    if signatureMethod == "" {
        signatureMethod = DEFAULT_SIGNATURE_METHOD
    }
    switch signatureMethod {
        case SIGNATURE_METHOD_HMAC_SHA256:
            h = crypto.SHA256
      	case SIGNATURE_METHOD_HMAC_SHA1:
      	    h = crypto.SHA1
      	default:
      	    err = ErrorInvalidSignatureMethod
      	    return
    }
    originalSignature := qstring.Get("Signature")
    qstring["Signature"] = nil, false
    qstring["DSOCAccessKeyId"] = []string{p.accessKey}
    
  	var signature []byte
    req.URL.RawQuery = qstring.Encode()
    canonicalizedStringToSign, err := p.Canonicalize(req)
    if err != nil { return }
    //log.Printf("String-to-sign: '%s'", canonicalizedStringToSign)

	signature, err = p.SignEncoded(h, canonicalizedStringToSign, base64.StdEncoding)
	if err != nil || string(signature) != originalSignature {
	    err = ErrorSignatureDoesNotMatch
	} else {
	    validSignature = true
	}
  	return
}


