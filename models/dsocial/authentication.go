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

type ConsumerKey struct {
    PersistableModel `json:"model,omitempty,collapse"`
    ConsumerId       string        `json:"consumer_id,omitempty"`
    HashType         HashAlgorithm `json:"hash_type,omitempty"`
    PrivateKey       string        `json:"private_key,omitempty"`
}

func (p HashAlgorithm) Hasher() (hasher hash.Hash) {
    switch p {
    case MD5:
        hasher = md5.New()
    case SHA1:
        hasher = sha1.New()
    case SHA224:
        hasher = sha256.New224()
    case SHA256:
        hasher = sha256.New()
    case SHA384:
        hasher = sha512.New384()
    case SHA512:
        hasher = sha512.New()
    default:
        hasher = sha512.New()
    }
    return
}

func (p *UserPassword) SetPassword(password string) {
    p.Salt = generateSalt(120)
    p.HashType = DEFAULT_HASH_ALGORITHM
    p.HashValue = generateHashedValued(p.HashType, p.Salt, password)
}

func (p *UserPassword) CheckPassword(password string) (bool) {
    return checkHashedValue(p.HashType, p.Salt, password, p.HashValue)
}

func (p *UserKey) GeneratePrivateKey() {
    p.PrivateKey = generateSalt(512)
    p.HashType = DEFAULT_HASH_ALGORITHM
}

func (p *UserKey) CheckHashedValue(testData, testHashedValue string) (bool) {
    return checkHashedValue(p.HashType, p.PrivateKey, testData, testHashedValue)
}

func (p *UserKey) CheckHashedByteValue(testData []byte, testHashedValue string) (bool) {
    return checkHashedByteValue(p.HashType, p.PrivateKey, testData, testHashedValue)
}

func (p *ConsumerKey) GeneratePrivateKey() {
    p.PrivateKey = generateSalt(512)
    p.HashType = DEFAULT_HASH_ALGORITHM
}

func (p *ConsumerKey) CheckHashedValue(testData, testHashedValue string) (bool) {
    return checkHashedValue(p.HashType, p.PrivateKey, testData, testHashedValue)
}

func (p *ConsumerKey) CheckHashedByteValue(testData []byte, testHashedValue string) (bool) {
    return checkHashedByteValue(p.HashType, p.PrivateKey, testData, testHashedValue)
}

func generateSalt(length int) (string) {
    l := length / 8 + 1
    arr := make([]string, l)
    for i := 0; i < l; i++ {
        s := strconv.Uitob(uint(rand.Uint32()), 16)
        for len(s) < 8 {
            s = "0" + s
        }
        arr[i] = s
    }
    salt := strings.Join(arr, "")
    l = len(salt)
    if l > length {
        salt = salt[l-length:]
    }
    return salt
}

func generateHashedValued(hashType HashAlgorithm, salt, testData string) (string) {
    hasher := hashType.Hasher()
    hasher.Write([]byte(salt))
    hasher.Write([]byte(testData))
    return fmt.Sprintf("%x", hasher.Sum())
}

func checkHashedValue(hashType HashAlgorithm, salt, testData, testHashedValue string) (bool) {
    return testHashedValue == generateHashedValued(hashType, salt, testData)
}

func generateHashedByteValued(hashType HashAlgorithm, salt string, testData []byte) (string) {
    hasher := hashType.Hasher()
    hasher.Write([]byte(salt))
    hasher.Write(testData)
    return fmt.Sprintf("%x", hasher.Sum())
}

func checkHashedByteValue(hashType HashAlgorithm, salt string, testData []byte, testHashedValue string) (bool) {
    return testHashedValue == generateHashedByteValued(hashType, salt, testData)
}
