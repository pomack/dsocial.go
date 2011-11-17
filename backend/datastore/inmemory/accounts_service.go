package inmemory

import (
    dm "github.com/pomack/dsocial.go/models/dsocial"
    ba "github.com/pomack/dsocial.go/backend/accounts"
    "os"
    "strconv"
    "time"
)

func (p *InMemoryDataStore) retrieveUserAccountCollection() (m *inMemoryCollection) {
    return p.retrieveCollection(_INMEMORY_USER_ACCOUNT_COLLECTION_NAME)
}

func (p *InMemoryDataStore) retrieveUserAccountIdForUsernameCollection() (m *inMemoryCollection) {
    return p.retrieveCollection(_INMEMORY_USER_ACCOUNT_ID_FOR_USERNAME_COLLECTION_NAME)
}

func (p *InMemoryDataStore) retrieveUserAccountIdForEmailCollection() (m *inMemoryCollection) {
    return p.retrieveCollection(_INMEMORY_USER_ACCOUNT_ID_FOR_EMAIL_COLLECTION_NAME)
}

func (p *InMemoryDataStore) retrieveConsumerAccountCollection() (m *inMemoryCollection) {
    return p.retrieveCollection(_INMEMORY_CONSUMER_ACCOUNT_COLLECTION_NAME)
}

func (p *InMemoryDataStore) retrieveExternalUserAccountCollection() (m *inMemoryCollection) {
    return p.retrieveCollection(_INMEMORY_EXTERNAL_USER_ACCOUNT_COLLECTION_NAME)
}

func (p *InMemoryDataStore) generateIdForAccount(collectionName string) (string) {
    nextId := collectionName + "/" + strconv.Itoa64(p.NextId)
    p.NextId++
    return nextId
}

func (p *InMemoryDataStore) CreateUserAccount(user *dm.User) (*dm.User, os.Error) {
    if user.Id == "" {
        user.Id = p.generateIdForAccount(_INMEMORY_USER_ACCOUNT_COLLECTION_NAME)
    }
    if _, ok := p.retrieve(_INMEMORY_USER_ACCOUNT_COLLECTION_NAME, user.Id); ok {
        return user, ERR_ACCOUNT_ALREADY_EXISTS_WITH_SPECIFIED_ID
    }
    if _, ok := p.retrieve(_INMEMORY_USER_ACCOUNT_ID_FOR_USERNAME_COLLECTION_NAME, user.Username); ok {
        return user, ERR_ACCOUNT_ALREADY_EXISTS_WITH_SPECIFIED_USERNAME
    }
    if _, ok := p.retrieve(_INMEMORY_USER_ACCOUNT_ID_FOR_EMAIL_COLLECTION_NAME, user.Email); ok {
        return user, ERR_ACCOUNT_ALREADY_EXISTS_WITH_SPECIFIED_EMAIL
    }
    p.store(user.Id, _INMEMORY_USER_ACCOUNT_COLLECTION_NAME, user.Id, user)
    p.store(user.Id, _INMEMORY_USER_ACCOUNT_ID_FOR_USERNAME_COLLECTION_NAME, user.Username, user.Id)
    p.store(user.Id, _INMEMORY_USER_ACCOUNT_ID_FOR_EMAIL_COLLECTION_NAME, user.Email, user.Id)
    return user, nil
}

func (p *InMemoryDataStore) UpdateUserAccount(user *dm.User) (*dm.User, os.Error) {
    return nil, nil
}

func (p *InMemoryDataStore) DeleteUserAccount(user *dm.User) (*dm.User, os.Error) {
    return nil, nil
}


func (p *InMemoryDataStore) CreateConsumerAccount(user *dm.Consumer) (*dm.User, os.Error) {
    return nil, nil
}

func (p *InMemoryDataStore) UpdateConsumerAccount(user *dm.Consumer) (*dm.User, os.Error) {
    return nil, nil
}

func (p *InMemoryDataStore) DeleteConsumerAccount(user *dm.Consumer) (*dm.User, os.Error) {
    return nil, nil
}


func (p *InMemoryDataStore) CreateExternalUserAccount(user *dm.ExternalUser) (*dm.ExternalUser, os.Error) {
    return nil, nil
}

func (p *InMemoryDataStore) UpdateExternalUserAccount(user *dm.ExternalUser) (*dm.ExternalUser, os.Error) {
    return nil, nil
}

func (p *InMemoryDataStore) DeleteExternalUserAccount(user *dm.ExternalUser) (*dm.ExternalUser, os.Error) {
    return nil, nil
}


func (p *InMemoryDataStore) RetrieveUserAccountById(id string) (*dm.User, os.Error) {
    return nil, nil
}

func (p *InMemoryDataStore) FindUserAccountByUsername(username string) (*dm.User, os.Error) {
    return nil, nil
}

func (p *InMemoryDataStore) FindUserAccountsByEmail(email string, next ba.NextToken, maxResults int) ([]*dm.User, ba.NextToken, os.Error) {
    return nil, nil, nil
}

func (p *InMemoryDataStore) FindUserAccountsByPhoneNumber(phoneNumber string, next ba.NextToken, maxResults int) ([]*dm.User, ba.NextToken, os.Error) {
    return nil, nil, nil
}


func (p *InMemoryDataStore) RetrieveConsumerAccountById(id string) (*dm.Consumer, os.Error) {
    return nil, nil
}

func (p *InMemoryDataStore) FindConsumerAccountByShortName(shortName string) (*dm.Consumer, os.Error) {
    return nil, nil
}

func (p *InMemoryDataStore) FindConsumerAccountsByDomainName(domainName string, next ba.NextToken, maxResults int) ([]*dm.Consumer, ba.NextToken, os.Error) {
    return nil, nil, nil
}

func (p *InMemoryDataStore) FindConsumerAccountsByName(name string, exact bool, next ba.NextToken, maxResults int) ([]*dm.Consumer, ba.NextToken, os.Error) {
    return nil, nil, nil
}


func (p *InMemoryDataStore) RetrieveExternalUserAccountById(id string) (*dm.ExternalUser, os.Error) {
    return nil, nil
}

func (p *InMemoryDataStore) RetrieveExternalUserAccountByConsumerAndExternalUserId(consumerId, externalUserId string) (*dm.ExternalUser, os.Error) {
    return nil, nil
}

func (p *InMemoryDataStore) FindExternalUserAccountsByConsumerId(consumerId string, next ba.NextToken, maxResults int) ([]*dm.ExternalUser, ba.NextToken, os.Error) {
    return nil, nil, nil
}

func (p *InMemoryDataStore) FindExternalUserAccountsByExternalUserId(externalUserId string, next ba.NextToken, maxResults int) ([]*dm.ExternalUser, ba.NextToken, os.Error) {
    return nil, nil, nil
}


func (p *InMemoryDataStore) AllowLoginByUserId(userId string) (bool, os.Error) {
    return false, nil
}

func (p *InMemoryDataStore) DisableLogin(userId string) (os.Error) {
    return nil
}

func (p *InMemoryDataStore) DisableLoginAt(userId string, at *time.Time) (os.Error) {
    return nil
}
