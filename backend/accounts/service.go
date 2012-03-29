package accounts

import (
    "errors"
    dm "github.com/pomack/dsocial.go/models/dsocial"
    "time"
)

var (
    ERR_ACCOUNT_ALREADY_EXISTS_WITH_SPECIFIED_ID          error
    ERR_ACCOUNT_ALREADY_EXISTS_WITH_SPECIFIED_USERNAME    error
    ERR_ACCOUNT_ALREADY_EXISTS_WITH_SPECIFIED_SHORTNAME   error
    ERR_ACCOUNT_ALREADY_EXISTS_WITH_SPECIFIED_EMAIL       error
    ERR_ACCOUNT_ALREADY_EXISTS_WITH_SPECIFIED_DOMAIN_NAME error
    ERR_ACCOUNT_MUST_SPECIFY_SHORTNAME                    error
)

func init() {
    ERR_ACCOUNT_ALREADY_EXISTS_WITH_SPECIFIED_ID = errors.New("Account already exists with specified id")
    ERR_ACCOUNT_ALREADY_EXISTS_WITH_SPECIFIED_USERNAME = errors.New("Account already exists with specified username")
    ERR_ACCOUNT_ALREADY_EXISTS_WITH_SPECIFIED_SHORTNAME = errors.New("Account already exists with specified short-name")
    ERR_ACCOUNT_ALREADY_EXISTS_WITH_SPECIFIED_EMAIL = errors.New("Account already exists with specified email")
    ERR_ACCOUNT_ALREADY_EXISTS_WITH_SPECIFIED_DOMAIN_NAME = errors.New("Account already exists with specified domain name")
    ERR_ACCOUNT_MUST_SPECIFY_SHORTNAME = errors.New("Must specify short-name")
}

type NextToken interface{}

type DataStore interface {
    CreateUserAccount(user *dm.User) (*dm.User, error)
    UpdateUserAccount(user *dm.User) (*dm.User, error)
    DeleteUserAccount(user *dm.User) (*dm.User, error)

    CreateConsumerAccount(user *dm.Consumer) (*dm.Consumer, error)
    UpdateConsumerAccount(user *dm.Consumer) (*dm.Consumer, error)
    DeleteConsumerAccount(user *dm.Consumer) (*dm.Consumer, error)

    CreateExternalUserAccount(user *dm.ExternalUser) (*dm.ExternalUser, error)
    UpdateExternalUserAccount(user *dm.ExternalUser) (*dm.ExternalUser, error)
    DeleteExternalUserAccount(user *dm.ExternalUser) (*dm.ExternalUser, error)

    RetrieveUserAccountById(id string) (*dm.User, error)
    FindUserAccountByUsername(username string) (*dm.User, error)
    FindUserAccountsByEmail(email string, next NextToken, maxResults int) ([]*dm.User, NextToken, error)

    RetrieveConsumerAccountById(id string) (*dm.Consumer, error)
    FindConsumerAccountByShortName(shortName string) (*dm.Consumer, error)
    FindConsumerAccountsByDomainName(domainName string, next NextToken, maxResults int) ([]*dm.Consumer, NextToken, error)
    FindConsumerAccountsByName(name string, exact bool, next NextToken, maxResults int) ([]*dm.Consumer, NextToken, error)

    RetrieveExternalUserAccountById(id string) (*dm.ExternalUser, error)
    RetrieveExternalUserAccountByConsumerAndExternalUserId(consumerId, externalUserId string) (*dm.ExternalUser, error)
    FindExternalUserAccountsByConsumerId(consumerId string, next NextToken, maxResults int) ([]*dm.ExternalUser, NextToken, error)
    FindExternalUserAccountsByExternalUserId(externalUserId string, next NextToken, maxResults int) ([]*dm.ExternalUser, NextToken, error)
}

func AllowLoginByUserId(ds DataStore, userId string) (bool, error) {
    user, err := ds.RetrieveUserAccountById(userId)
    if user != nil && user.AllowLogin && (user.DisableLoginAt <= 0 || user.DisableLoginAt < time.Now().Unix()) {
        return true, err
    }
    return false, err
}

func DisableLogin(ds DataStore, userId string) error {
    user, err := ds.RetrieveUserAccountById(userId)
    if user == nil || err != nil {
        return err
    }
    user.AllowLogin = false
    _, err = ds.UpdateUserAccount(user)
    return err
}

func DisableLoginAt(ds DataStore, userId string, at time.Time) error {
    user, err := ds.RetrieveUserAccountById(userId)
    if user == nil || err != nil {
        return err
    }
    if user.AllowLogin {
        now := time.Now().Unix()
        if user.DisableLoginAt == 0 || user.DisableLoginAt < now || at.Unix() < now {
            user.AllowLogin = false
            user.DisableLoginAt = 0
        } else if at.IsZero() {
            user.DisableLoginAt = 0
        } else {
            user.DisableLoginAt = at.Unix()
        }
        _, err = ds.UpdateUserAccount(user)
    }
    return err
}
