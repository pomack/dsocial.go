package dsocial

type UserPassword struct {
    PersistableModel `json:"model,omitempty,collapse"`
    UserId           string        `json:"user_id,omitempty"`
    HashType         HashAlgorithm `json:"hash_type,omitempty"`
    Salt             string        `json:"salt,omitempty"`
    HashValue        string        `json:"hash_value,omitempty"`
}

type UserKey struct {
    PersistableModel `json:"model,omitempty,collapse"`
    UserId           string        `json:"user_id,omitempty"`
    HashType         HashAlgorithm `json:"hash_type,omitempty"`
    PrivateKey       string        `json:"private_key,omitempty"`
}

type ConsumerPassword struct {
    PersistableModel `json:"model,omitempty,collapse"`
    ConsumerId       string        `json:"consumer_id,omitempty"`
    HashType         HashAlgorithm `json:"hash_type,omitempty"`
    PrivateKey       string        `json:"private_key,omitempty"`
}
