package dsocial

import (
    "hash"
    "crypto/md5"
    "crypto/sha1"
    "crypto/sha256"
    "crypto/sha512"
    "fmt"
    "rand"
    "strconv"
    "strings"
)

type HashAlgorithm int

const (
    MD5                    HashAlgorithm = iota
    SHA1                   HashAlgorithm = iota
    SHA224                 HashAlgorithm = iota
    SHA256                 HashAlgorithm = iota
    SHA384                 HashAlgorithm = iota
    SHA512                 HashAlgorithm = iota
    DEFAULT_HASH_ALGORITHM HashAlgorithm = SHA512
)

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
