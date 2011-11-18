package authentication

import (
    dm "github.com/pomack/dsocial.go/models/dsocial"
    "os"
)


type NextToken interface{}

type DataStore interface {
    RetrieveUserPassword(userId string) (*dm.UserPassword, os.Error)
    RetrieveConsumerKey(consumerKeyId string) (*dm.ConsumerKey, os.Error)
    RetrieveUserKey(userKeyId string) (*dm.UserKey, os.Error)
    RetrieveConsumerKeys(consumerId string, next NextToken, maxResults int) ([]*dm.ConsumerKey, NextToken, os.Error)
    RetrieveUserKeys(userId string, next NextToken, maxResults int) ([]*dm.UserKey, NextToken, os.Error)
    
    StoreUserPassword(password *dm.UserPassword) (*dm.UserPassword, os.Error)
    StoreConsumerKey(key *dm.ConsumerKey) (*dm.ConsumerKey, os.Error)
    StoreUserKey(key *dm.UserKey) (*dm.UserKey, os.Error)
    
    DeleteUserPassword(userId string) (*dm.UserPassword, os.Error)
    DeleteConsumerKey(consumerKeyId string) (*dm.ConsumerKey, os.Error)
    DeleteUserKey(userKeyId string) (*dm.UserKey, os.Error)
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

func GenerateNewConsumerKey(ds DataStore, consumerId string) (*dm.ConsumerKey, os.Error) {
    pwd := &dm.ConsumerKey{ConsumerId: consumerId}
    pwd.GeneratePrivateKey()
    return ds.StoreConsumerKey(pwd)
}

func UpsertConsumerKey(ds DataStore, consumerId string) (*dm.ConsumerKey, os.Error) {
    pwd, err := ds.RetrieveConsumerKey(consumerId)
    if err != nil {
        return pwd, err
    }
    if pwd == nil {
        pwd = &dm.ConsumerKey{ConsumerId: consumerId}
    }
    pwd.GeneratePrivateKey()
    return ds.StoreConsumerKey(pwd)
}

func GenerateNewUserKey(ds DataStore, userId string) (*dm.UserKey, os.Error) {
    pwd := &dm.UserKey{UserId: userId}
    pwd.GeneratePrivateKey()
    return ds.StoreUserKey(pwd)
}

func UpsertUserKey(ds DataStore, userId string) (*dm.UserKey, os.Error) {
    pwd, err := ds.RetrieveUserKey(userId)
    if err != nil {
        return pwd, err
    }
    if pwd == nil {
        pwd = &dm.UserKey{UserId: userId}
    }
    pwd.GeneratePrivateKey()
    return ds.StoreUserKey(pwd)
}

func ValidateUserPassword(ds DataStore, userId, password string) (isValid bool, err os.Error) {
    pwd, err := ds.RetrieveUserPassword(userId)
    if pwd != nil {
        isValid = pwd.CheckPassword(password)
    }
    return
}

func ValidateConsumerKeyString(ds DataStore, consumerKeyId, data, signature string) (bool, os.Error) {
    pwd, err := ds.RetrieveConsumerKey(consumerKeyId)
    return pwd != nil && pwd.CheckHashedValue(data, signature), err
}

func ValidateConsumerKeyBytes(ds DataStore, consumerKeyId string, data []byte, signature string) (bool, os.Error) {
    pwd, err := ds.RetrieveConsumerKey(consumerKeyId)
    return pwd != nil && pwd.CheckHashedByteValue(data, signature), err
}

func ValidateUserKeyString(ds DataStore, userKeyId, data, signature string) (bool, os.Error) {
    pwd, err := ds.RetrieveUserKey(userKeyId)
    return pwd != nil && pwd.CheckHashedValue(data, signature), err
}

func ValidateUserKeyBytes(ds DataStore, userKeyId string, data []byte, signature string) (bool, os.Error) {
    pwd, err := ds.RetrieveUserKey(userKeyId)
    return pwd != nil && pwd.CheckHashedByteValue(data, signature), err
}
