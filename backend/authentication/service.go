package authentication

import dm "github.com/pomack/dsocial.go/models/dsocial"

type NextToken interface{}

type DataStore interface {
    RetrieveUserPassword(userId string) (*dm.UserPassword, error)
    RetrieveAccessKey(accessKeyId string) (*dm.AccessKey, error)
    RetrieveConsumerKeys(consumerId string, next NextToken, maxResults int) ([]*dm.AccessKey, NextToken, error)
    RetrieveUserKeys(userId string, next NextToken, maxResults int) ([]*dm.AccessKey, NextToken, error)

    StoreUserPassword(password *dm.UserPassword) (*dm.UserPassword, error)
    StoreAccessKey(key *dm.AccessKey) (*dm.AccessKey, error)

    DeleteUserPassword(userId string) (*dm.UserPassword, error)
    DeleteAccessKey(accessKeyId string) (*dm.AccessKey, error)
}

func SetUserPassword(ds DataStore, userId, password string) (*dm.UserPassword, error) {
    pwd, err := ds.RetrieveUserPassword(userId)
    if err != nil {
        return pwd, err
    }
    if pwd == nil {
        pwd = &dm.UserPassword{UserId: userId}
    }
    pwd.SetPassword(password)
    return ds.StoreUserPassword(pwd)
}

func GenerateNewAccessKey(ds DataStore, userId, consumerId string) (*dm.AccessKey, error) {
    pwd := &dm.AccessKey{UserId: userId, ConsumerId: consumerId}
    pwd.GeneratePrivateKey()
    pwd.GenerateId()
    return ds.StoreAccessKey(pwd)
}

func DeleteAccessKey(ds DataStore, accessKeyId string) (*dm.AccessKey, error) {
    return ds.DeleteAccessKey(accessKeyId)
}

func ValidateUserPassword(ds DataStore, userId, password string) (isValid bool, err error) {
    pwd, err := ds.RetrieveUserPassword(userId)
    if pwd != nil {
        isValid = pwd.CheckPassword(password)
    }
    return
}

func ValidateAccessKeyString(ds DataStore, hashAlgorithm dm.HashAlgorithm, accessKeyId, data, signature string) (bool, error) {
    pwd, err := ds.RetrieveAccessKey(accessKeyId)
    return pwd != nil && pwd.CheckHashedValue(hashAlgorithm, data, signature), err
}

func ValidateAccessKeyBytes(ds DataStore, hashAlgorithm dm.HashAlgorithm, accessKeyId string, data []byte, signature string) (bool, error) {
    pwd, err := ds.RetrieveAccessKey(accessKeyId)
    return pwd != nil && pwd.CheckHashedByteValue(hashAlgorithm, data, signature), err
}
