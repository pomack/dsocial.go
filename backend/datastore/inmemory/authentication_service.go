package inmemory

import (
    ba "github.com/pomack/dsocial.go/backend/authentication"
    dm "github.com/pomack/dsocial.go/models/dsocial"
    "os"
)

func (p *InMemoryDataStore) RetrieveUserPassword(userId string) (*dm.UserPassword, os.Error) {
    pwd, _ := p.retrieve(_INMEMORY_USER_PASSWORD_COLLECTION_NAME, userId)
    if pwd == nil {
        return nil, nil
    }
    return pwd.(*dm.UserPassword), nil
}

func (p *InMemoryDataStore) RetrieveConsumerKey(consumerKeyId string) (*dm.ConsumerKey, os.Error) {
    pwd, _ := p.retrieve(_INMEMORY_CONSUMER_KEYS_COLLECTION_NAME, consumerKeyId)
    if pwd == nil {
        return nil, nil
    }
    return pwd.(*dm.ConsumerKey), nil
}

func (p *InMemoryDataStore) RetrieveUserKey(userKeyId string) (*dm.UserKey, os.Error) {
    pwd, _ := p.retrieve(_INMEMORY_USER_KEYS_COLLECTION_NAME, userKeyId)
    if pwd == nil {
        return nil, nil
    }
    return pwd.(*dm.UserKey), nil
}

func (p *InMemoryDataStore) RetrieveConsumerKeys(consumerId string, next ba.NextToken, maxResults int) ([]*dm.ConsumerKey, ba.NextToken, os.Error) {
    m := p.retrieveStringMapCollection(consumerId, _INMEMORY_CONSUMER_KEYS_FOR_CONSUMER_ID_COLLECTION_NAME, consumerId)
    arr := make([]*dm.ConsumerKey, len(m))
    i := 0
    var err os.Error
    for k := range m {
        arr[i], err = p.RetrieveConsumerKey(k)
        i++
        if err != nil {
            break
        }
    }
    return arr, nil, err
}

func (p *InMemoryDataStore) RetrieveUserKeys(userId string, next ba.NextToken, maxResults int) ([]*dm.UserKey, ba.NextToken, os.Error) {
    m := p.retrieveStringMapCollection(userId, _INMEMORY_USER_KEYS_FOR_USER_ID_COLLECTION_NAME, userId)
    arr := make([]*dm.UserKey, len(m))
    i := 0
    var err os.Error
    for k := range m {
        arr[i], err = p.RetrieveUserKey(k)
        i++
        if err != nil {
            break
        }
    }
    return arr, nil, err
}

func (p *InMemoryDataStore) StoreUserPassword(password *dm.UserPassword) (*dm.UserPassword, os.Error) {
    p.store(password.UserId, _INMEMORY_USER_PASSWORD_COLLECTION_NAME, password.UserId, password)
    p.addToStringMapCollection(password.UserId, _INMEMORY_USER_KEYS_FOR_USER_ID_COLLECTION_NAME, password.UserId, password.Id, password.Id)
    return password, nil
}

func (p *InMemoryDataStore) StoreConsumerKey(key *dm.ConsumerKey) (*dm.ConsumerKey, os.Error) {
    p.store(key.ConsumerId, _INMEMORY_CONSUMER_KEYS_COLLECTION_NAME, key.Id, key)
    p.addToStringMapCollection(key.ConsumerId, _INMEMORY_CONSUMER_KEYS_FOR_CONSUMER_ID_COLLECTION_NAME, key.ConsumerId, key.Id, key.Id)
    return key, nil
}

func (p *InMemoryDataStore) StoreUserKey(key *dm.UserKey) (*dm.UserKey, os.Error) {
    p.store(key.UserId, _INMEMORY_USER_KEYS_COLLECTION_NAME, key.Id, key)
    return key, nil
}

func (p *InMemoryDataStore) DeleteUserPassword(userId string) (*dm.UserPassword, os.Error) {
    oldValue, _ := p.delete(userId, _INMEMORY_USER_PASSWORD_COLLECTION_NAME, userId)
    if oldValue != nil {
        pwd, _ := oldValue.(*dm.UserPassword)
        return pwd, nil
    }
    return nil, nil
}

func (p *InMemoryDataStore) DeleteConsumerKey(consumerKeyId string) (oldKey *dm.ConsumerKey, err os.Error) {
    oldValue, _ := p.delete("", _INMEMORY_CONSUMER_KEYS_COLLECTION_NAME, consumerKeyId)
    if oldValue != nil {
        if key, ok := oldValue.(*dm.ConsumerKey); ok {
            p.removeFromStringMapCollection(key.ConsumerId, _INMEMORY_CONSUMER_KEYS_FOR_CONSUMER_ID_COLLECTION_NAME, key.ConsumerId, consumerKeyId)
            oldKey = key
        }
    }
    return
}

func (p *InMemoryDataStore) DeleteUserKey(userKeyId string) (oldKey *dm.UserKey, err os.Error) {
    oldValue, _ := p.delete("", _INMEMORY_USER_KEYS_COLLECTION_NAME, userKeyId)
    if oldValue != nil {
        if key, ok := oldValue.(*dm.UserKey); ok {
            p.removeFromStringMapCollection(key.UserId, _INMEMORY_USER_KEYS_FOR_USER_ID_COLLECTION_NAME, key.UserId, userKeyId)
            oldKey = key
        }
    }
    return 
}

func (p *InMemoryDataStore) BackendAuthenticationDataStore() (ba.DataStore) {
    return p
}
