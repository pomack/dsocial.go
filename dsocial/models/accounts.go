package dsocial

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
