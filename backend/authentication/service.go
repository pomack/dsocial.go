package authentication

import (
    dm "github.com/pomack/dsocial.go/models/dsocial"
    "os"
)


type NextToken interface{}

type DataStore interface {
    SetUserPassword(userId, password string) (os.Error)
    SetConsumerPassword(consumerId, password string) (os.Error)
    SetUserKey(userId, hashType dm.HashAlgorithm, privateKey string) (os.Error)
    
    RetrieveUserPassword(userId string) (*dm.UserPassword, os.Error)
    RetrieveConsumerPassword(consumerId string) (*dm.ConsumerPassword, os.Error)
    RetrieveUserKeys(userId, next NextToken, maxResults int) ([]*dm.UserKey, NextToken, os.Error)
    
    ValidateUserPassword(userId, password string) (bool, os.Error)
    ValidateConsumerPassword(userId, password string) (bool, os.Error)
    ValidateUserKeyString(userId, data, signature string) (bool, os.Error)
    ValidateUserKeyBytes(userId string, data []byte, signature string) (bool, os.Error)
}

