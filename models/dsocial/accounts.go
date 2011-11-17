package dsocial

import (
    "os"
    "strings"
)

type UserRole int32

const (
    ROLE_ANONYMOUS = 1 << iota
    ROLE_STANDARD
    ROLE_BUSINESS_SUPPORT
    ROLE_TECHNICAL_SUPPORT
    ROLE_BACKUP_OPERATOR
    ROLE_SECURITY_OFFICER
    ROLE_OWNER = 1 << 29
    ROLE_ADMIN = 1 << 30
)

type Consumer struct {
    PersistableModel  `json:"model,omitempty,collapse"`
    DomainName        string `json:"domain_name,omitempty"`
    HomePage          string `json:"home_page,omitempty"`
    AuthorizationPage string `json:"authorization_page,omitempty"`
    ShortName         string `json:"short_name,omitempty"`
    Name              string `json:"name,omitempty"`
    Email             string `json:"email,omitempty"`
    PhoneNumber       string `json:"phone_number,omitempty"`
    IsTrusted         bool   `json:"is_trusted,omitempty"`
    IsSuggested       bool   `json:"is_suggested,omitempty"`
    AllowLogin        bool   `json:"allow_login,omitempty"`
    DisableLoginAt    int64  `json:"disable_login_at,omitempty"`
}

type User struct {
    PersistableModel `json:"model,omitempty,collapse"`
    Role             int32  `json:"role,omitempty"`
    Name             string `json:"name,omitempty"`
    Username         string `json:"username,omitempty"`
    Email            string `json:"email,omitempty"`
    PhoneNumber      string `json:"phone_number,omitempty"`
    Address          string `json:"address,omitempty"`
    ContactId        string `json:"contact_id,omitempty"`
    AllowLogin       bool   `json:"allow_login,omitempty"`
    IsPayingUser     bool   `json:"is_paying_user,omitempty"`
    Notes            string `json:"notes,omitempty"`
    DisableLoginAt   int64  `json:"disable_login_at,omitempty"`
}

type ExternalUser struct {
    PersistableModel `json:"model,omitempty,collapse"`
    ConsumerId       string `json:"consumer_id,omitempty"`
    ExternalUserId   string `json:"external_user_id,omitempty"`
    Name             string `json:"name,omitempty"`
}

func (p *Consumer) Validate(createNew bool, errors map[string][]os.Error) (isValid bool) {
    if errors == nil {
        errors = make(map[string][]os.Error)
    }
    p.PersistableModel.Validate(createNew, errors)
    p.DomainName, _ = validateDomainName(p.DomainName, false, "domain_name", errors)
    p.HomePage, _ = validateUrl(p.HomePage, false, "home_page", errors)
    p.AuthorizationPage, _ = validateUrl(p.AuthorizationPage, false, "authorization_page", errors)
    p.ShortName, _ = validateAlphaNumeric(p.ShortName, true, false, true, false, false, "short_name", errors)
    p.Name, _ = validateNonEmpty(p.Name, true, "name", errors)
    p.Email, _ = validateEmail(p.Email, false, "email", errors)
    p.PhoneNumber = strings.TrimSpace(p.PhoneNumber)
    isValid = len(errors) == 0
    return
}

func (p *User) Validate(createNew bool, errors map[string][]os.Error) (isValid bool) {
    if errors == nil {
        errors = make(map[string][]os.Error)
    }
    p.PersistableModel.Validate(createNew, errors)
    p.Role = p.Role & ROLE_STANDARD & ROLE_BUSINESS_SUPPORT & ROLE_TECHNICAL_SUPPORT & ROLE_BACKUP_OPERATOR & ROLE_SECURITY_OFFICER & ROLE_OWNER & ROLE_ADMIN
    if p.Role == 0 {
        p.Role = ROLE_STANDARD
    }
    p.Username, _ = validateAlphaNumeric(p.Username, true, false, true, false, false, "short_name", errors)
    p.Email, _ = validateEmail(p.Email, false, "email", errors)
    p.PhoneNumber = strings.TrimSpace(p.PhoneNumber)
    p.Address = strings.TrimSpace(p.Address)
    p.ContactId, _ = validateId(p.ContactId, true, "contact_id", errors)
    isValid = len(errors) == 0
    return
}

func (p *ExternalUser) Validate(createNew bool, errors map[string][]os.Error) (isValid bool) {
    if errors == nil {
        errors = make(map[string][]os.Error)
    }
    p.PersistableModel.Validate(createNew, errors)
    p.ConsumerId, _ = validateId(p.ConsumerId, false, "consumer_id", errors)
    p.ExternalUserId, _ = validateId(p.ExternalUserId, false, "external_user_id", errors)
    p.Name, _ = validateNonEmpty(p.Name, true, "name", errors)
    isValid = len(errors) == 0
    return
}

