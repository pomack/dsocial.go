package dsocial

import (
    "crypto/md5"
    "crypto/sha1"
    "crypto/sha256"
    "crypto/sha512"
    "fmt"
    "hash"
    "rand"
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

var validSaltCharacters = []string{
    "0", "1", "2", "3", "4", "5", "6", "7", "8", "9",
    "a", "b", "c", "d", "e", "f", "g", "h", "i", "j", "k", "l", "m", "n", "o", "p", "q", "r", "s", "t", "u", "v", "w", "x", "y", "z",
    "A", "B", "C", "D", "E", "F", "G", "H", "I", "J", "K", "L", "M", "N", "O", "P", "Q", "R", "S", "T", "U", "V", "W", "X", "Y", "Z",
    "!", "@", "(", ")", "-", "_", "|", ";", ":", ",", ".", "$", "*", "[", "]", "{", "}",
}

type UserPassword struct {
    PersistableModel `json:"model,omitempty,collapse"`
    UserId           string        `json:"user_id,omitempty"`
    HashType         HashAlgorithm `json:"hash_type,omitempty"`
    Salt             string        `json:"salt,omitempty"`
    HashValue        string        `json:"hash_value,omitempty"`
}

type AccessKey struct {
    PersistableModel `json:"model,omitempty,collapse"`
    UserId           string `json:"user_id,omitempty"`
    ConsumerId       string `json:"consumer_id,omitempty"`
    PrivateKey       string `json:"private_key,omitempty"`
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

func (p HashAlgorithm) IsValid() (valid bool) {
    switch p {
    case MD5, SHA1, SHA224, SHA256, SHA384, SHA512:
        valid = true
    }
    return
}

func NewUserPassword(userId, password string) *UserPassword {
    p := &UserPassword{UserId: userId}
    p.SetPassword(password)
    return p
}

func (p *UserPassword) SetPassword(password string) {
    p.Salt = generateSalt(120)
    p.HashType = DEFAULT_HASH_ALGORITHM
    p.HashValue = generateHashedValued(p.HashType, p.Salt, password)
}

func (p *UserPassword) CheckPassword(password string) bool {
    return checkHashedValue(p.HashType, p.Salt, password, p.HashValue)
}

func NewAccessKey(userId, consumerId string) *AccessKey {
    key := &AccessKey{UserId: userId, ConsumerId: consumerId}
    key.GeneratePrivateKey()
    key.GenerateId()
    return key
}

func (p *AccessKey) GeneratePrivateKey() {
    p.PrivateKey = generateSalt(512)
}

func (p *AccessKey) GenerateId() {
    p.Id = generateSalt(40)
}

func (p *AccessKey) CheckHashedValue(hashType HashAlgorithm, testData, testHashedValue string) bool {
    return checkHashedValue(hashType, p.PrivateKey, testData, testHashedValue)
}

func (p *AccessKey) CheckHashedByteValue(hashType HashAlgorithm, testData []byte, testHashedValue string) bool {
    return checkHashedByteValue(hashType, p.PrivateKey, testData, testHashedValue)
}

func generateSalt(l int) string {
    chars := len(validSaltCharacters)
    arr := make([]string, l)
    for i := 0; i < l; i++ {
        arr[i] = validSaltCharacters[rand.Intn(chars)]
    }
    return strings.Join(arr, "")
}

func generateHashedValued(hashType HashAlgorithm, salt, testData string) string {
    hasher := hashType.Hasher()
    hasher.Write([]byte(salt))
    hasher.Write([]byte(testData))
    return fmt.Sprintf("%x", hasher.Sum())
}

func checkHashedValue(hashType HashAlgorithm, salt, testData, testHashedValue string) bool {
    return testHashedValue == generateHashedValued(hashType, salt, testData)
}

func generateHashedByteValued(hashType HashAlgorithm, salt string, testData []byte) string {
    hasher := hashType.Hasher()
    hasher.Write([]byte(salt))
    hasher.Write(testData)
    return fmt.Sprintf("%x", hasher.Sum())
}

func checkHashedByteValue(hashType HashAlgorithm, salt string, testData []byte, testHashedValue string) bool {
    return testHashedValue == generateHashedByteValued(hashType, salt, testData)
}
