package dsocial

import (
    "github.com/pomack/jsonhelper.go/jsonhelper"
    "strings"
    "time"
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

func (p *Consumer) Accessible() bool {
    if p == nil {
        return false
    }
    if p.AllowLogin == false {
        return false
    }
    if p.DisableLoginAt > 0 && p.DisableLoginAt < time.Now().Unix() {
        return false
    }
    return true
}

func (p *Consumer) InitFromJSONObject(obj jsonhelper.JSONObject) {
    p.PersistableModel.InitFromJSONObject(obj)
    p.DomainName = obj.GetAsString("domain_name")
    p.HomePage = obj.GetAsString("home_page")
    p.AuthorizationPage = obj.GetAsString("authorization_page")
    p.ShortName = obj.GetAsString("short_name")
    p.Name = obj.GetAsString("name")
    p.Email = obj.GetAsString("email")
    p.PhoneNumber = obj.GetAsString("phone_number")
    p.IsTrusted = obj.GetAsBool("is_trusted")
    p.IsSuggested = obj.GetAsBool("is_suggested")
    p.AllowLogin = obj.GetAsBool("allow_login")
    p.DisableLoginAt = obj.GetAsInt64("disable_login_at")
}

func (p *Consumer) CleanFromUser(user *User, original *Consumer) {
    if original == nil {
        p.PersistableModel.CleanFromUser(user, nil)
    } else {
        p.PersistableModel.CleanFromUser(user, &original.PersistableModel)
    }
    if original == nil {
        p.IsTrusted = false
        p.IsSuggested = false
        p.AllowLogin = true
        p.DisableLoginAt = 0
    } else {
        p.ShortName = original.ShortName
        p.IsTrusted = original.IsTrusted
        p.IsSuggested = original.IsSuggested
        p.AllowLogin = original.AllowLogin
        p.DisableLoginAt = original.DisableLoginAt
    }
}

func (p *Consumer) Validate(createNew bool, errors map[string][]error) (isValid bool) {
    if errors == nil {
        errors = make(map[string][]error)
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

func (p *User) Accessible() bool {
    if p == nil {
        return false
    }
    if p.AllowLogin == false {
        return false
    }
    if p.DisableLoginAt > 0 && p.DisableLoginAt < time.Now().Unix() {
        return false
    }
    return true
}

func (p *User) InitFromJSONObject(obj jsonhelper.JSONObject) {
    p.PersistableModel.InitFromJSONObject(obj)
    p.Role = obj.GetAsInt32("role")
    p.Name = obj.GetAsString("name")
    p.Username = obj.GetAsString("username")
    p.Email = obj.GetAsString("email")
    p.PhoneNumber = obj.GetAsString("phone_number")
    p.Address = obj.GetAsString("address")
    p.ContactId = obj.GetAsString("contact_id")
    p.AllowLogin = obj.GetAsBool("allow_login")
    p.IsPayingUser = obj.GetAsBool("is_paying_user")
    p.Notes = obj.GetAsString("notes")
    p.DisableLoginAt = obj.GetAsInt64("disable_login_at")
}

func (p *User) CleanFromUser(user *User, original *User) {
    if original == nil {
        p.PersistableModel.CleanFromUser(user, nil)
    } else {
        p.PersistableModel.CleanFromUser(user, &original.PersistableModel)
    }
    if user == nil || user.Role != ROLE_ADMIN {
        if original == nil {
            p.Role = ROLE_STANDARD
        } else {
            p.Role = original.Role
        }
    }
    if original == nil {
        p.ContactId = ""
        p.AllowLogin = true
        p.IsPayingUser = false
        p.Notes = ""
        p.DisableLoginAt = 0
    } else {
        p.Username = original.Username
        p.ContactId = original.ContactId
        p.AllowLogin = original.AllowLogin
        p.IsPayingUser = original.IsPayingUser
        p.Notes = original.Notes
        p.DisableLoginAt = original.DisableLoginAt
    }
}

func (p *User) Validate(createNew bool, errors map[string][]error) (isValid bool) {
    if errors == nil {
        errors = make(map[string][]error)
    }
    p.PersistableModel.Validate(createNew, errors)
    p.Role = p.Role & (ROLE_STANDARD | ROLE_BUSINESS_SUPPORT | ROLE_TECHNICAL_SUPPORT | ROLE_BACKUP_OPERATOR | ROLE_SECURITY_OFFICER | ROLE_OWNER | ROLE_ADMIN)
    if p.Role == 0 {
        p.Role = ROLE_STANDARD
    }
    p.Name, _ = validateNonEmpty(p.Name, true, "name", errors)
    p.Username, _ = validateAlphaNumeric(p.Username, true, false, true, false, false, "username", errors)
    p.Email, _ = validateEmail(p.Email, false, "email", errors)
    p.PhoneNumber = strings.TrimSpace(p.PhoneNumber)
    p.Address = strings.TrimSpace(p.Address)
    p.ContactId, _ = validateId(p.ContactId, true, "contact_id", errors)
    isValid = len(errors) == 0
    return
}

func (p *ExternalUser) InitFromJSONObject(obj jsonhelper.JSONObject) {
    p.PersistableModel.InitFromJSONObject(obj)
    p.ConsumerId = obj.GetAsString("consumer_id")
    p.ExternalUserId = obj.GetAsString("external_user_id")
    p.Name = obj.GetAsString("name")
}

func (p *ExternalUser) CleanFromUser(user *User, original *ExternalUser) {
    if original == nil {
        p.PersistableModel.CleanFromUser(user, nil)
    } else {
        p.PersistableModel.CleanFromUser(user, &original.PersistableModel)
    }
    if original != nil {
        if p.ConsumerId == "" {
            p.ConsumerId = original.ConsumerId
        }
        if p.ExternalUserId == "" {
            p.ExternalUserId = original.ExternalUserId
        }
        if p.Name == "" {
            p.Name = original.Name
        }
    }
}

func (p *ExternalUser) Validate(createNew bool, errors map[string][]error) (isValid bool) {
    if errors == nil {
        errors = make(map[string][]error)
    }
    p.PersistableModel.Validate(createNew, errors)
    p.ConsumerId, _ = validateId(p.ConsumerId, false, "consumer_id", errors)
    p.ExternalUserId, _ = validateId(p.ExternalUserId, false, "external_user_id", errors)
    p.Name, _ = validateNonEmpty(p.Name, true, "name", errors)
    isValid = len(errors) == 0
    return
}
