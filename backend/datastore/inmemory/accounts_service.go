package inmemory

import (
    dm "github.com/pomack/dsocial.go/models/dsocial"
    ba "github.com/pomack/dsocial.go/backend/accounts"
    "os"
    "strconv"
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

func (p *InMemoryDataStore) generateIdForAccount(collectionName string) string {
    nextId := collectionName + "/" + strconv.Itoa64(p.NextId)
    p.NextId++
    return nextId
}

func (p *InMemoryDataStore) CreateUserAccount(user *dm.User) (*dm.User, os.Error) {
    if user == nil {
        return nil, nil
    }
    if user.Id == "" {
        user.Id = p.generateIdForAccount(_INMEMORY_USER_ACCOUNT_COLLECTION_NAME)
    }
    if _, ok := p.retrieve(_INMEMORY_USER_ACCOUNT_COLLECTION_NAME, user.Id); ok {
        return user, ba.ERR_ACCOUNT_ALREADY_EXISTS_WITH_SPECIFIED_ID
    }
    if _, ok := p.retrieve(_INMEMORY_USER_ACCOUNT_ID_FOR_USERNAME_COLLECTION_NAME, user.Username); ok {
        return user, ba.ERR_ACCOUNT_ALREADY_EXISTS_WITH_SPECIFIED_USERNAME
    }
    if _, ok := p.retrieve(_INMEMORY_USER_ACCOUNT_ID_FOR_EMAIL_COLLECTION_NAME, user.Email); ok {
        return user, ba.ERR_ACCOUNT_ALREADY_EXISTS_WITH_SPECIFIED_EMAIL
    }
    p.store(user.Id, _INMEMORY_USER_ACCOUNT_COLLECTION_NAME, user.Id, user)
    p.store(user.Id, _INMEMORY_USER_ACCOUNT_ID_FOR_USERNAME_COLLECTION_NAME, user.Username, user.Id)
    p.store(user.Id, _INMEMORY_USER_ACCOUNT_ID_FOR_EMAIL_COLLECTION_NAME, user.Email, user.Id)
    return user, nil
}

func (p *InMemoryDataStore) UpdateUserAccount(user *dm.User) (*dm.User, os.Error) {
    if user == nil {
        return nil, nil
    }
    if user.Id == "" {
        return user, dm.ERR_MUST_SPECIFY_ID
    }
    oldUserI, _ := p.retrieve(_INMEMORY_USER_ACCOUNT_COLLECTION_NAME, user.Id)
    var oldUser *dm.User = nil
    if oldUserI != nil {
        oldUser = oldUserI.(*dm.User)
    }
    if oldUser == nil || oldUser.Username != user.Username {
        if _, ok := p.retrieve(_INMEMORY_USER_ACCOUNT_ID_FOR_USERNAME_COLLECTION_NAME, user.Username); ok {
            return user, ba.ERR_ACCOUNT_ALREADY_EXISTS_WITH_SPECIFIED_USERNAME
        }
    }
    if oldUser == nil || oldUser.Email != user.Email {
        if _, ok := p.retrieve(_INMEMORY_USER_ACCOUNT_ID_FOR_EMAIL_COLLECTION_NAME, user.Email); ok {
            return user, ba.ERR_ACCOUNT_ALREADY_EXISTS_WITH_SPECIFIED_EMAIL
        }
    }
    p.store(user.Id, _INMEMORY_USER_ACCOUNT_COLLECTION_NAME, user.Id, user)
    p.store(user.Id, _INMEMORY_USER_ACCOUNT_ID_FOR_USERNAME_COLLECTION_NAME, user.Username, user.Id)
    p.store(user.Id, _INMEMORY_USER_ACCOUNT_ID_FOR_EMAIL_COLLECTION_NAME, user.Email, user.Id)
    if oldUser != nil {
        if oldUser.Username != user.Username {
            p.delete(_INMEMORY_USER_ACCOUNT_ID_FOR_USERNAME_COLLECTION_NAME, oldUser.Username)
        }
        if oldUser.Email != user.Email {
            p.delete(_INMEMORY_USER_ACCOUNT_ID_FOR_EMAIL_COLLECTION_NAME, oldUser.Email)
        }
    }
    return nil, nil
}

func (p *InMemoryDataStore) DeleteUserAccount(user *dm.User) (*dm.User, os.Error) {
    if user == nil {
        return nil, nil
    }
    oldUserI, _ := p.retrieve(_INMEMORY_USER_ACCOUNT_COLLECTION_NAME, user.Id)
    var oldUser *dm.User = nil
    if oldUserI != nil {
        oldUser = oldUserI.(*dm.User)
    }
    if oldUser != nil {
        p.delete(_INMEMORY_USER_ACCOUNT_COLLECTION_NAME, oldUser.Id)
        p.delete(_INMEMORY_USER_ACCOUNT_ID_FOR_USERNAME_COLLECTION_NAME, oldUser.Username)
        p.delete(_INMEMORY_USER_ACCOUNT_ID_FOR_EMAIL_COLLECTION_NAME, oldUser.Email)
    }
    return oldUser, nil
}

func (p *InMemoryDataStore) CreateConsumerAccount(user *dm.Consumer) (*dm.Consumer, os.Error) {
    if user == nil {
        return nil, nil
    }
    if user.Id == "" {
        user.Id = p.generateIdForAccount(_INMEMORY_CONSUMER_ACCOUNT_COLLECTION_NAME)
    }
    if _, ok := p.retrieve(_INMEMORY_CONSUMER_ACCOUNT_COLLECTION_NAME, user.Id); ok {
        return user, ba.ERR_ACCOUNT_ALREADY_EXISTS_WITH_SPECIFIED_ID
    }
    if _, ok := p.retrieve(_INMEMORY_CONSUMER_ACCOUNT_ID_FOR_SHORTNAME_COLLECTION_NAME, user.ShortName); ok {
        return user, ba.ERR_ACCOUNT_ALREADY_EXISTS_WITH_SPECIFIED_SHORTNAME
    }
    p.store(user.Id, _INMEMORY_CONSUMER_ACCOUNT_COLLECTION_NAME, user.Id, user)
    p.store(user.Id, _INMEMORY_CONSUMER_ACCOUNT_ID_FOR_SHORTNAME_COLLECTION_NAME, user.ShortName, user.Id)
    p.addToStringMapCollection(user.Id, _INMEMORY_CONSUMER_ACCOUNT_IDS_FOR_NAME_COLLECTION_NAME, user.Name, user.Id, user.Id)
    p.addToStringMapCollection(user.Id, _INMEMORY_CONSUMER_ACCOUNT_IDS_FOR_DOMAIN_NAME_COLLECTION_NAME, user.DomainName, user.Id, user.Id)
    return user, nil
}

func (p *InMemoryDataStore) UpdateConsumerAccount(user *dm.Consumer) (*dm.Consumer, os.Error) {
    if user == nil {
        return nil, nil
    }
    if user.Id == "" {
        return user, dm.ERR_MUST_SPECIFY_ID
    }
    oldUserI, _ := p.retrieve(_INMEMORY_CONSUMER_ACCOUNT_COLLECTION_NAME, user.Id)
    var oldUser *dm.Consumer = nil
    if oldUserI != nil {
        oldUser = oldUserI.(*dm.Consumer)
    }
    if oldUser == nil || oldUser.ShortName != user.ShortName {
        if _, ok := p.retrieve(_INMEMORY_CONSUMER_ACCOUNT_ID_FOR_SHORTNAME_COLLECTION_NAME, user.ShortName); ok {
            return user, ba.ERR_ACCOUNT_ALREADY_EXISTS_WITH_SPECIFIED_SHORTNAME
        }
    }
    if oldUser == nil || oldUser.DomainName != user.DomainName {
        if _, ok := p.retrieve(_INMEMORY_CONSUMER_ACCOUNT_IDS_FOR_DOMAIN_NAME_COLLECTION_NAME, user.DomainName); ok {
            return user, ba.ERR_ACCOUNT_ALREADY_EXISTS_WITH_SPECIFIED_DOMAIN_NAME
        }
    }
    p.store(user.Id, _INMEMORY_CONSUMER_ACCOUNT_COLLECTION_NAME, user.Id, user)
    p.store(user.Id, _INMEMORY_CONSUMER_ACCOUNT_ID_FOR_SHORTNAME_COLLECTION_NAME, user.ShortName, user.Id)
    p.addToStringMapCollection(user.Id, _INMEMORY_CONSUMER_ACCOUNT_IDS_FOR_NAME_COLLECTION_NAME, user.Name, user.Id, user.Id)
    p.addToStringMapCollection(user.Id, _INMEMORY_CONSUMER_ACCOUNT_IDS_FOR_DOMAIN_NAME_COLLECTION_NAME, user.DomainName, user.Id, user.Id)
    if oldUser != nil {
        if oldUser.ShortName != user.ShortName {
            p.delete(_INMEMORY_CONSUMER_ACCOUNT_ID_FOR_SHORTNAME_COLLECTION_NAME, oldUser.ShortName)
        }
        if oldUser.Name != user.Name {
            p.removeFromStringMapCollection(user.Id, _INMEMORY_CONSUMER_ACCOUNT_IDS_FOR_NAME_COLLECTION_NAME, oldUser.Name, user.Id)
        }
        if oldUser.DomainName != user.DomainName {
            p.removeFromStringMapCollection(user.Id, _INMEMORY_CONSUMER_ACCOUNT_IDS_FOR_DOMAIN_NAME_COLLECTION_NAME, oldUser.DomainName, user.Id)
        }
    }
    return user, nil
}

func (p *InMemoryDataStore) DeleteConsumerAccount(user *dm.Consumer) (*dm.Consumer, os.Error) {
    if user == nil {
        return nil, nil
    }
    if user.Id == "" {
        return user, dm.ERR_MUST_SPECIFY_ID
    }
    oldUserI, _ := p.retrieve(_INMEMORY_CONSUMER_ACCOUNT_COLLECTION_NAME, user.Id)
    var oldUser *dm.Consumer = nil
    if oldUserI != nil {
        oldUser = oldUserI.(*dm.Consumer)
    }
    if oldUser != nil {
        p.delete(_INMEMORY_CONSUMER_ACCOUNT_COLLECTION_NAME, user.Id)
        p.delete(_INMEMORY_CONSUMER_ACCOUNT_ID_FOR_SHORTNAME_COLLECTION_NAME, oldUser.ShortName)
        p.removeFromStringMapCollection(user.Id, _INMEMORY_CONSUMER_ACCOUNT_IDS_FOR_NAME_COLLECTION_NAME, oldUser.Name, user.Id)
        p.removeFromStringMapCollection(user.Id, _INMEMORY_CONSUMER_ACCOUNT_IDS_FOR_DOMAIN_NAME_COLLECTION_NAME, oldUser.DomainName, user.Id)
    }
    return user, nil
}

func (p *InMemoryDataStore) CreateExternalUserAccount(user *dm.ExternalUser) (*dm.ExternalUser, os.Error) {
    if user == nil {
        return nil, nil
    }
    if user.Id == "" {
        user.Id = p.generateIdForAccount(_INMEMORY_EXTERNAL_USER_ACCOUNT_COLLECTION_NAME)
    }
    if _, ok := p.retrieve(_INMEMORY_EXTERNAL_USER_ACCOUNT_COLLECTION_NAME, user.Id); ok {
        return user, ba.ERR_ACCOUNT_ALREADY_EXISTS_WITH_SPECIFIED_ID
    }
    p.store(user.Id, _INMEMORY_EXTERNAL_USER_ACCOUNT_COLLECTION_NAME, user.Id, user)
    p.addToStringMapCollection(user.Id, _INMEMORY_EXTERNAL_ACCOUNT_IDS_FOR_CONSUMER_ID_COLLECTION_NAME, user.ConsumerId, user.Id, user.Id)
    p.addToStringMapCollection(user.Id, _INMEMORY_EXTERNAL_ACCOUNT_IDS_FOR_EXTERNAL_USER_ID_COLLECTION_NAME, user.ExternalUserId, user.Id, user.Id)
    return user, nil
}

func (p *InMemoryDataStore) UpdateExternalUserAccount(user *dm.ExternalUser) (*dm.ExternalUser, os.Error) {
    if user == nil {
        return nil, nil
    }
    if user.Id == "" {
        return user, dm.ERR_MUST_SPECIFY_ID
    }
    oldUserI, _ := p.retrieve(_INMEMORY_EXTERNAL_USER_ACCOUNT_COLLECTION_NAME, user.Id)
    var oldUser *dm.ExternalUser = nil
    if oldUserI != nil {
        oldUser = oldUserI.(*dm.ExternalUser)
    }
    p.store(user.Id, _INMEMORY_EXTERNAL_USER_ACCOUNT_COLLECTION_NAME, user.Id, user)
    p.addToStringMapCollection(user.Id, _INMEMORY_EXTERNAL_ACCOUNT_IDS_FOR_CONSUMER_ID_COLLECTION_NAME, user.ConsumerId, user.Id, user.Id)
    p.addToStringMapCollection(user.Id, _INMEMORY_EXTERNAL_ACCOUNT_IDS_FOR_EXTERNAL_USER_ID_COLLECTION_NAME, user.ExternalUserId, user.Id, user.Id)
    if oldUser != nil {
        if oldUser.ConsumerId != user.ConsumerId {
            p.removeFromStringMapCollection(user.Id, _INMEMORY_EXTERNAL_ACCOUNT_IDS_FOR_CONSUMER_ID_COLLECTION_NAME, oldUser.ConsumerId, user.Id)
        }
        if oldUser.ExternalUserId != user.ExternalUserId {
            p.removeFromStringMapCollection(user.Id, _INMEMORY_EXTERNAL_ACCOUNT_IDS_FOR_EXTERNAL_USER_ID_COLLECTION_NAME, oldUser.ExternalUserId, user.Id)
        }
    }
    return user, nil
}

func (p *InMemoryDataStore) DeleteExternalUserAccount(user *dm.ExternalUser) (*dm.ExternalUser, os.Error) {
    if user == nil {
        return nil, nil
    }
    if user.Id == "" {
        return user, dm.ERR_MUST_SPECIFY_ID
    }
    oldUserI, _ := p.retrieve(_INMEMORY_EXTERNAL_USER_ACCOUNT_COLLECTION_NAME, user.Id)
    var oldUser *dm.ExternalUser = nil
    if oldUserI != nil {
        oldUser = oldUserI.(*dm.ExternalUser)
    }
    if oldUser != nil {
        p.delete(_INMEMORY_EXTERNAL_USER_ACCOUNT_COLLECTION_NAME, user.Id)
        p.removeFromStringMapCollection(user.Id, _INMEMORY_EXTERNAL_ACCOUNT_IDS_FOR_CONSUMER_ID_COLLECTION_NAME, oldUser.ConsumerId, user.Id)
        p.removeFromStringMapCollection(user.Id, _INMEMORY_EXTERNAL_ACCOUNT_IDS_FOR_EXTERNAL_USER_ID_COLLECTION_NAME, oldUser.ExternalUserId, user.Id)
    }
    return user, nil
}

func (p *InMemoryDataStore) RetrieveUserAccountById(id string) (*dm.User, os.Error) {
    user, _ := p.retrieve(_INMEMORY_USER_ACCOUNT_COLLECTION_NAME, id)
    if user != nil {
        return user.(*dm.User), nil
    }
    return nil, nil
}

func (p *InMemoryDataStore) FindUserAccountByUsername(username string) (*dm.User, os.Error) {
    uid, _ := p.retrieveString(_INMEMORY_USER_ACCOUNT_ID_FOR_USERNAME_COLLECTION_NAME, username)
    if uid == "" {
        return nil, nil
    }
    return p.RetrieveUserAccountById(uid)
}

func (p *InMemoryDataStore) FindUserAccountsByEmail(email string, next ba.NextToken, maxResults int) ([]*dm.User, ba.NextToken, os.Error) {
    uid, _ := p.retrieveString(_INMEMORY_USER_ACCOUNT_ID_FOR_EMAIL_COLLECTION_NAME, email)
    if uid == "" {
        return nil, nil, nil
    }
    user, err := p.RetrieveUserAccountById(uid)
    if user == nil {
        return make([]*dm.User, 0), nil, err
    }
    return []*dm.User{user}, nil, err
}

func (p *InMemoryDataStore) RetrieveConsumerAccountById(id string) (*dm.Consumer, os.Error) {
    user, _ := p.retrieve(_INMEMORY_CONSUMER_ACCOUNT_COLLECTION_NAME, id)
    if user != nil {
        return user.(*dm.Consumer), nil
    }
    return nil, nil
}

func (p *InMemoryDataStore) FindConsumerAccountByShortName(shortName string) (*dm.Consumer, os.Error) {
    uid, _ := p.retrieveString(_INMEMORY_CONSUMER_ACCOUNT_ID_FOR_SHORTNAME_COLLECTION_NAME, shortName)
    if uid == "" {
        return nil, nil
    }
    return p.RetrieveConsumerAccountById(uid)
}

func (p *InMemoryDataStore) FindConsumerAccountsByDomainName(domainName string, next ba.NextToken, maxResults int) ([]*dm.Consumer, ba.NextToken, os.Error) {
    m := p.retrieveStringMapCollection("", _INMEMORY_CONSUMER_ACCOUNT_IDS_FOR_DOMAIN_NAME_COLLECTION_NAME, domainName)
    arr := make([]*dm.Consumer, len(m))
    i := 0
    for k := range m {
        arr[i], _ = p.RetrieveConsumerAccountById(k)
        i++
    }
    return arr, nil, nil
}

func (p *InMemoryDataStore) FindConsumerAccountsByName(name string, exact bool, next ba.NextToken, maxResults int) ([]*dm.Consumer, ba.NextToken, os.Error) {
    // TODO handle non-exact name matches
    m := p.retrieveStringMapCollection("", _INMEMORY_CONSUMER_ACCOUNT_IDS_FOR_NAME_COLLECTION_NAME, name)
    arr := make([]*dm.Consumer, len(m))
    i := 0
    for k := range m {
        arr[i], _ = p.RetrieveConsumerAccountById(k)
        i++
    }
    return arr, nil, nil
}

func (p *InMemoryDataStore) RetrieveExternalUserAccountById(id string) (*dm.ExternalUser, os.Error) {
    user, _ := p.retrieve(_INMEMORY_EXTERNAL_USER_ACCOUNT_COLLECTION_NAME, id)
    if user != nil {
        return user.(*dm.ExternalUser), nil
    }
    return nil, nil
}

func (p *InMemoryDataStore) RetrieveExternalUserAccountByConsumerAndExternalUserId(consumerId, externalUserId string) (*dm.ExternalUser, os.Error) {
    return nil, nil
}

func (p *InMemoryDataStore) FindExternalUserAccountsByConsumerId(consumerId string, next ba.NextToken, maxResults int) ([]*dm.ExternalUser, ba.NextToken, os.Error) {
    m := p.retrieveStringMapCollection("", _INMEMORY_EXTERNAL_ACCOUNT_IDS_FOR_CONSUMER_ID_COLLECTION_NAME, consumerId)
    arr := make([]*dm.ExternalUser, len(m))
    i := 0
    for k := range m {
        arr[i], _ = p.RetrieveExternalUserAccountById(k)
        i++
    }
    return arr, nil, nil
}

func (p *InMemoryDataStore) FindExternalUserAccountsByExternalUserId(externalUserId string, next ba.NextToken, maxResults int) ([]*dm.ExternalUser, ba.NextToken, os.Error) {
    m := p.retrieveStringMapCollection("", _INMEMORY_EXTERNAL_ACCOUNT_IDS_FOR_EXTERNAL_USER_ID_COLLECTION_NAME, externalUserId)
    arr := make([]*dm.ExternalUser, len(m))
    i := 0
    for k := range m {
        arr[i], _ = p.RetrieveExternalUserAccountById(k)
        i++
    }
    return arr, nil, nil
}

func (p *InMemoryDataStore) BackendAccountsDataStore() ba.DataStore {
    return p
}
