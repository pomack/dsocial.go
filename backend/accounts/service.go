package accounts

import (
    dm "github.com/pomack/dsocial.go/models/dsocial"
    "os"
    "time"
)


var (
    ERR_ACCOUNT_ALREADY_EXISTS_WITH_SPECIFIED_ID os.Error
    ERR_ACCOUNT_ALREADY_EXISTS_WITH_SPECIFIED_USERNAME os.Error
    ERR_ACCOUNT_ALREADY_EXISTS_WITH_SPECIFIED_SHORTNAME os.Error
    ERR_ACCOUNT_ALREADY_EXISTS_WITH_SPECIFIED_EMAIL os.Error
    ERR_ACCOUNT_ALREADY_EXISTS_WITH_SPECIFIED_DOMAIN_NAME os.Error
    ERR_ACCOUNT_MUST_SPECIFY_SHORTNAME os.Error
)

func init() {
    ERR_ACCOUNT_ALREADY_EXISTS_WITH_SPECIFIED_ID = os.NewError("Account already exists with specified id")
    ERR_ACCOUNT_ALREADY_EXISTS_WITH_SPECIFIED_USERNAME = os.NewError("Account already exists with specified username")
    ERR_ACCOUNT_ALREADY_EXISTS_WITH_SPECIFIED_SHORTNAME = os.NewError("Account already exists with specified short-name")
    ERR_ACCOUNT_ALREADY_EXISTS_WITH_SPECIFIED_EMAIL = os.NewError("Account already exists with specified email")
    ERR_ACCOUNT_ALREADY_EXISTS_WITH_SPECIFIED_DOMAIN_NAME = os.NewError("Account already exists with specified domain name")
    ERR_ACCOUNT_MUST_SPECIFY_SHORTNAME = os.NewError("Must specify short-name")
}


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

