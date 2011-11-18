package dsocial

type Authorization int
type AclType int

const (
    READ        Authorization = 1
    UPDATE      Authorization = 2
    DELETE      Authorization = 4
    SHARE_READ  Authorization = 8
    SHARE_WRITE Authorization = 16
    READ_ACL    Authorization = 32
    UPDATE_ACL  Authorization = 64
)

const (
    PUBLIC       AclType = 1
    NONANONYMOUS AclType = 2
    GROUP        AclType = 3
    INDIVIDUAL   AclType = 4
)

type Session struct {
    PersistableModel `json:"model,omitempty,collapse"`
    UserId           string   `json:"user_id,omitempty"`
    ConsumerId       string   `json:"consumer_id,omitempty"`
    ExternalUserId   string   `json:"external_user_id,omitempty"`
    ExpiresAt        int64    `json:"expires_at,omitempty"`
    IsConsumer       bool     `json:"is_consumer,omitempty"`
    IpAddress        string   `json:"ip_address,omitempty"`
    Scopes           []string `json:"scopes,omitempty"`
    ExtraData        string   `json:"extra_data,omitempty"`
    Name             string   `json:"name,omitempty"`
}

type AuthorizationToken struct {
    PersistableModel `json:"model,omitempty,collapse"`
    UserId           string   `json:"user_id,omitempty"`
    ConsumerId       string   `json:"consumer_id,omitempty"`
    ExternalUserId   string   `json:"external_user_id,omitempty"`
    ExpiresAt        int64    `json:"expires_at,omitempty"`
    IsConsumer       bool     `json:"is_consumer,omitempty"`
    IpAddress        string   `json:"ip_address,omitempty"`
    Scopes           []string `json:"scopes,omitempty"`
    ExtraData        string   `json:"extra_data,omitempty"`
    Name             string   `json:"name,omitempty"`
}

type AclKind struct {
    Rel      AclType `json:"rel,omitempty"`
    EntityId string  `json:"entity_id,omitempty"`
}

type AclEntry struct {
    PersistableModel   `json:"model,omitempty,collapse"`
    AuthorizationLevel int32      `json:"authorization_level,omitempty"`
    EntityIds          []*AclKind `json:"entity_ids,omitempty"`
    ForKeys            []string   `json:"for_keys,omitempty"`
}

type Acl struct {
    PersistableModel `json:"model,omitempty,collapse"`
    OwnerId          string      `json:"owner_id,omitempty"`
    Entries          []*AclEntry `json:"entries,omitempty"`
}

type AclPersistableModel struct {
    PersistableModel `json:"model,omitempty,collapse"`
    Acl              Acl `json:"acl,omitempty"`
}

func (p *Session) UID() (uid string) {
    if p != nil {
        if len(p.ExternalUserId) > 0 {
            uid = p.ExternalUserId
        } else if len(p.ConsumerId) > 0 {
            uid = p.ConsumerId
        } else {
            uid = p.UserId
        }
    }
    return
}

func (p *AuthorizationToken) UID() (uid string) {
    if p != nil {
        if len(p.ExternalUserId) > 0 {
            uid = p.ExternalUserId
        } else if len(p.ConsumerId) > 0 {
            uid = p.ConsumerId
        } else {
            uid = p.UserId
        }
    }
    return
}
