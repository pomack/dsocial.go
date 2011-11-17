package dsocial

import (
    "os"
    "strings"
    "unicode"
    "url"
)

func validateId(id string, allowEmpty bool, idName string, errors map[string][]os.Error) (newId string, err os.Error) {
    if idName == "" {
        idName = "id"
    }
    newId = strings.TrimSpace(id)
    if newId == "" {
        if !allowEmpty {
            err = ERR_MUST_SPECIFY_ID
        }
    }
    if err == nil {
        for _, rune := range newId {
            switch {
            case unicode.IsLetter(rune), unicode.IsNumber(rune):
            default:
                switch rune {
                case '/', '-', '_', '@', '.':
                default:
                    err = ERR_INVALID_ID
                }
            }
            if err != nil {
                break
            }
        }
    }
    if err != nil && errors != nil {
        errors[idName] = []os.Error{err}
    }
    return
}

func validateEmail(email string, allowEmpty bool, idName string, errors map[string][]os.Error) (newEmail string, err os.Error) {
    if idName == "" {
        idName = "email"
    }
    newEmail = strings.TrimSpace(email)
    if newEmail == "" {
        if !allowEmpty {
            err = ERR_MUST_SPECIFY_ID
        }
    }
    if err == nil {
        hasCharacterBeforeAt, hasAt, hasCharacterAfterAt, hasDotAfterAt, hasCharacterAfterDot := false, false, false, false, false
        for _, rune := range newEmail {
            switch {
            case unicode.IsLetter(rune), unicode.IsNumber(rune):
                if hasDotAfterAt {
                    hasCharacterAfterDot = true
                } else if hasAt {
                    hasCharacterAfterAt = true
                } else {
                    hasCharacterBeforeAt = true
                }
            default:
                switch rune {
                case '@':
                    if hasAt {
                        err = ERR_INVALID_EMAIL_ADDRESS
                        break
                    } else {
                        if !hasCharacterBeforeAt {
                            err = ERR_INVALID_EMAIL_ADDRESS
                        }
                        hasAt = true
                    }
                case '.':
                    if hasCharacterAfterAt {
                        hasDotAfterAt = true
                    }
                case '-', '_', '+':
                    if hasDotAfterAt {
                        hasCharacterAfterDot = true
                    } else if hasAt {
                        hasCharacterAfterAt = true
                    } else {
                        hasCharacterBeforeAt = true
                    }
                default:
                    err = ERR_INVALID_EMAIL_ADDRESS
                }
            }
            if err != nil {
                break
            }
        }
        if !hasCharacterAfterDot {
            err = ERR_INVALID_FORMAT
        }
    }
    if err != nil && errors != nil {
        errors[idName] = []os.Error{err}
    }
    return
}

func validateNonEmpty(s string, trimSpace bool, idName string, errors map[string][]os.Error) (newS string, err os.Error) {
    if trimSpace {
        newS = strings.TrimSpace(s)
    } else {
        newS = s
    }
    if newS == "" {
        err = ERR_MUST_SPECIFY_ID
        if errors != nil {
            errors[idName] = []os.Error{err}
        }
    }
    return
}

func validateAlphaNumeric(s string, trimSpace, allowEmpty, mustStartAlpha, allowSpace, allowNewline bool, idName string, errors map[string][]os.Error) (newS string, err os.Error) {
    if trimSpace {
        newS = strings.TrimSpace(s)
    } else {
        newS = s
    }
    if newS == "" {
        if !allowEmpty {
            err = ERR_REQUIRED_FIELD
        }
    } else {
        for i, rune := range newS {
            if i == 0 && mustStartAlpha {
                if !unicode.IsLetter(rune) {
                    err = ERR_INVALID_FORMAT
                    break
                }
                continue
            }
            switch {
            case unicode.IsLetter(rune), unicode.IsNumber(rune):
            default:
                switch rune {
                case ' ':
                    if !allowSpace {
                        err = ERR_INVALID_FORMAT
                    }
                case '\r', '\n':
                    if !allowNewline {
                        err = ERR_INVALID_FORMAT
                    }
                default:
                    err = ERR_INVALID_FORMAT
                }
            }
            if err != nil {
                break
            }
        }
    }
    if err != nil && errors != nil {
        errors[idName] = []os.Error{err}
    }
    return
}

func validateDomainName(name string, allowEmpty bool, idName string, errors map[string][]os.Error) (newName string, err os.Error) {
    newName = strings.TrimSpace(name)
    if newName == "" {
        if !allowEmpty {
            err = ERR_REQUIRED_FIELD
        }
    } else {
        hasDot, hasCharacterAfterDot := false, false
        for i, rune := range newName {
            if i == 0 {
                if !unicode.IsLetter(rune) {
                    err = ERR_INVALID_FORMAT
                    break
                }
                continue
            }
            switch {
            case unicode.IsLetter(rune), unicode.IsNumber(rune):
                hasCharacterAfterDot = hasDot
            default:
                switch rune {
                case '.':
                    hasDot = true
                case '-':
                default:
                    err = ERR_INVALID_FORMAT
                }
            }
            if err != nil {
                break
            }
        }
        if !hasDot || !hasCharacterAfterDot {
            err = ERR_INVALID_FORMAT
        }
    }
    if err == nil {
        parts := strings.Split(newName, ".")
        l := len(parts)
        if l < 2 || len(join(".", parts[l-2], parts[l-1])) > 64 {
            // registered domain names can only be up to 64 characters long including global namespace (example.com)
            err = ERR_INVALID_FORMAT
        }
    }
    if err != nil && errors != nil {
        errors[idName] = []os.Error{err}
    }
    return
}

func validateUrl(uri string, allowEmpty bool, idName string, errors map[string][]os.Error) (newUri string, err os.Error) {
    newUri = strings.TrimSpace(uri)
    if newUri == "" {
        if !allowEmpty {
            err = ERR_REQUIRED_FIELD
        }
    } else {
        if _, err = url.ParseWithReference(newUri); err != nil {
            err = ERR_INVALID_FORMAT
        }
    }
    if err != nil && errors != nil {
        errors[idName] = []os.Error{err}
    }
    return
}
