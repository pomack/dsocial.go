package authentication

import (
    dm "github.com/pomack/dsocial.go/models/dsocial"
    "os"
)

type NextToken interface{}

type DataStore interface {
    RetrieveUserPassword(userId string) (*dm.UserPassword, os.Error)
    RetrieveAccessKey(accessKeyId string) (*dm.AccessKey, os.Error)
    RetrieveConsumerKeys(consumerId string, next NextToken, maxResults int) ([]*dm.AccessKey, NextToken, os.Error)
    RetrieveUserKeys(userId string, next NextToken, maxResults int) ([]*dm.AccessKey, NextToken, os.Error)

    StoreUserPassword(password *dm.UserPassword) (*dm.UserPassword, os.Error)
    StoreAccessKey(key *dm.AccessKey) (*dm.AccessKey, os.Error)

    DeleteUserPassword(userId string) (*dm.UserPassword, os.Error)
    DeleteAccessKey(accessKeyId string) (*dm.AccessKey, os.Error)
}

func SetUserPassword(ds DataStore, userId, password string) (*dm.UserPassword, os.Error) {
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

func GenerateNewAccessKey(ds DataStore, userId, consumerId string) (*dm.AccessKey, os.Error) {
    pwd := &dm.AccessKey{UserId: userId, ConsumerId: consumerId}
    pwd.GeneratePrivateKey()
    pwd.GenerateId()
    return ds.StoreAccessKey(pwd)
}

func DeleteAccessKey(ds DataStore, accessKeyId string) (*dm.AccessKey, os.Error) {
    return ds.DeleteAccessKey(accessKeyId)
}

func ValidateUserPassword(ds DataStore, userId, password string) (isValid bool, err os.Error) {
    pwd, err := ds.RetrieveUserPassword(userId)
    if pwd != nil {
        isValid = pwd.CheckPassword(password)
    }
    return
}

func ValidateAccessKeyString(ds DataStore, hashAlgorithm dm.HashAlgorithm, accessKeyId, data, signature string) (bool, os.Error) {
    pwd, err := ds.RetrieveAccessKey(accessKeyId)
    return pwd != nil && pwd.CheckHashedValue(hashAlgorithm, data, signature), err
}

func ValidateAccessKeyBytes(ds DataStore, hashAlgorithm dm.HashAlgorithm, accessKeyId string, data []byte, signature string) (bool, os.Error) {
    pwd, err := ds.RetrieveAccessKey(accessKeyId)
    return pwd != nil && pwd.CheckHashedByteValue(hashAlgorithm, data, signature), err
}
