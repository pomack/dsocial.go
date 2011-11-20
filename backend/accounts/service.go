package accounts

import (
    dm "github.com/pomack/dsocial.go/models/dsocial"
    "os"
    "time"
)

var (
    ERR_ACCOUNT_ALREADY_EXISTS_WITH_SPECIFIED_ID          os.Error
    ERR_ACCOUNT_ALREADY_EXISTS_WITH_SPECIFIED_USERNAME    os.Error
    ERR_ACCOUNT_ALREADY_EXISTS_WITH_SPECIFIED_SHORTNAME   os.Error
    ERR_ACCOUNT_ALREADY_EXISTS_WITH_SPECIFIED_EMAIL       os.Error
    ERR_ACCOUNT_ALREADY_EXISTS_WITH_SPECIFIED_DOMAIN_NAME os.Error
    ERR_ACCOUNT_MUST_SPECIFY_SHORTNAME                    os.Error
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

    CreateConsumerAccount(user *dm.Consumer) (*dm.Consumer, os.Error)
    UpdateConsumerAccount(user *dm.Consumer) (*dm.Consumer, os.Error)
    DeleteConsumerAccount(user *dm.Consumer) (*dm.Consumer, os.Error)

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
}

func AllowLoginByUserId(ds DataStore, userId string) (bool, os.Error) {
    user, err := ds.RetrieveUserAccountById(userId)
    if user != nil && user.AllowLogin && (user.DisableLoginAt <= 0 || user.DisableLoginAt < time.UTC().Seconds()) {
        return true, err
    }
    return false, err
}

func DisableLogin(ds DataStore, userId string) os.Error {
    user, err := ds.RetrieveUserAccountById(userId)
    if user == nil || err != nil {
        return err
    }
    user.AllowLogin = false
    _, err = ds.UpdateUserAccount(user)
    return err
}

func DisableLoginAt(ds DataStore, userId string, at *time.Time) os.Error {
    user, err := ds.RetrieveUserAccountById(userId)
    if user == nil || err != nil {
        return err
    }
    if user.AllowLogin {
        now := time.UTC().Seconds()
        if user.DisableLoginAt == 0 || user.DisableLoginAt < now || (at != nil && at.Seconds() < now) {
            user.AllowLogin = false
            user.DisableLoginAt = 0
        } else if at == nil {
            user.DisableLoginAt = 0
        } else {
            user.DisableLoginAt = at.Seconds()
        }
        _, err = ds.UpdateUserAccount(user)
    }
    return err
}
