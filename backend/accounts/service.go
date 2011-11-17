package accounts

import (
    dm "github.com/pomack/dsocial.go/models/dsocial"
    "os"
    "time"
)


type NextToken interface{}

type DataStore interface {
    CreateUserAccount(user *dm.User) (*dm.User, os.Error)
    UpdateUserAccount(user *dm.User) (*dm.User, os.Error)
    DeleteUserAccount(user *dm.User) (*dm.User, os.Error)
    
    CreateConsumerAccount(user *dm.Consumer) (*dm.User, os.Error)
    UpdateConsumerAccount(user *dm.Consumer) (*dm.User, os.Error)
    DeleteConsumerAccount(user *dm.Consumer) (*dm.User, os.Error)
    
    CreateExternalUserAccount(user *dm.ExternalUser) (*dm.ExternalUser, os.Error)
    UpdateExternalUserAccount(user *dm.ExternalUser) (*dm.ExternalUser, os.Error)
    DeleteExternalUserAccount(user *dm.ExternalUser) (*dm.ExternalUser, os.Error)
    
    RetrieveUserAccountById(id string) (*dm.User, os.Error)
    FindUserAccountByUsername(username string) (*dm.User, os.Error)
    FindUserAccountsByEmail(email string, next NextToken, maxResults int) ([]*dm.User, NextToken, os.Error)
    FindUserAccountsByPhoneNumber(phoneNumber string, next NextToken, maxResults int) ([]*dm.User, NextToken, os.Error)
    
    RetrieveConsumerAccountById(id string) (*dm.Consumer, os.Error)
    FindConsumerAccountByShortName(shortName string) (*dm.Consumer, os.Error)
    FindConsumerAccountsByDomainName(domainName string, next NextToken, maxResults int) ([]*dm.Consumer, NextToken, os.Error)
    FindConsumerAccountsByName(name string, exact bool, next NextToken, maxResults int) ([]*dm.Consumer, NextToken, os.Error)
    
    RetrieveExternalUserAccountById(id string) (*dm.ExternalUser, os.Error)
    RetrieveExternalUserAccountByConsumerAndExternalUserId(consumerId, externalUserId string) (*dm.ExternalUser, os.Error)
    FindExternalUserAccountsByConsumerId(consumerId string, next NextToken, maxResults int) ([]*dm.ExternalUser, NextToken, os.Error)
    FindExternalUserAccountsByExternalUserId(externalUserId string, next NextToken, maxResults int) ([]*dm.ExternalUser, NextToken, os.Error)
    
    AllowLoginByUserId(userId string) (bool, os.Error)
    DisableLogin(userId string) (os.Error)
    DisableLoginAt(userId string, at *time.Time) (os.Error)
}
