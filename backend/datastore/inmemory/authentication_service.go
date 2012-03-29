package inmemory

import (
    ba "github.com/pomack/dsocial.go/backend/authentication"
    dm "github.com/pomack/dsocial.go/models/dsocial"
)

func (p *InMemoryDataStore) RetrieveUserPassword(userId string) (*dm.UserPassword, error) {
    pwd, _ := p.retrieve(_INMEMORY_USER_PASSWORD_COLLECTION_NAME, userId)
    if pwd == nil {
        return nil, nil
    }
    return pwd.(*dm.UserPassword), nil
}

func (p *InMemoryDataStore) RetrieveAccessKey(accessKeyId string) (*dm.AccessKey, error) {
    pwd, _ := p.retrieve(_INMEMORY_ACCESS_KEYS_COLLECTION_NAME, accessKeyId)
    if pwd == nil {
        return nil, nil
    }
    return pwd.(*dm.AccessKey), nil
}

func (p *InMemoryDataStore) RetrieveConsumerKeys(consumerId string, next ba.NextToken, maxResults int) ([]*dm.AccessKey, ba.NextToken, error) {
    m := p.retrieveStringMapCollection(consumerId, _INMEMORY_ACCESS_KEYS_FOR_CONSUMER_ID_COLLECTION_NAME, consumerId)
    arr := make([]*dm.AccessKey, len(m))
    i := 0
    var err error
    for k := range m {
        arr[i], err = p.RetrieveAccessKey(k)
        i++
        if err != nil {
            break
        }
    }
    return arr, nil, err
}

func (p *InMemoryDataStore) RetrieveUserKeys(userId string, next ba.NextToken, maxResults int) ([]*dm.AccessKey, ba.NextToken, error) {
    m := p.retrieveStringMapCollection(userId, _INMEMORY_ACCESS_KEYS_FOR_USER_ID_COLLECTION_NAME, userId)
    arr := make([]*dm.AccessKey, len(m))
    i := 0
    var err error
    for k := range m {
        arr[i], err = p.RetrieveAccessKey(k)
        i++
        if err != nil {
            break
        }
    }
    return arr, nil, err
}

func (p *InMemoryDataStore) StoreUserPassword(password *dm.UserPassword) (*dm.UserPassword, error) {
    p.store(password.UserId, _INMEMORY_USER_PASSWORD_COLLECTION_NAME, password.UserId, password)
    return password, nil
}

func (p *InMemoryDataStore) StoreAccessKey(key *dm.AccessKey) (*dm.AccessKey, error) {
    if key.UserId == "" && key.ConsumerId == "" {
        return key, nil
    }
    if key.Id == "" {
        key.GenerateId()
    }
    if key.PrivateKey == "" {
        key.GeneratePrivateKey()
    }
    if err := key.BeforeCreate(); err != nil {
        return key, err
    }
    if err := key.BeforeSave(); err != nil {
        return key, err
    }
    p.store("", _INMEMORY_ACCESS_KEYS_COLLECTION_NAME, key.Id, key)
    if key.UserId != "" {
        p.addToStringMapCollection(key.UserId, _INMEMORY_ACCESS_KEYS_FOR_USER_ID_COLLECTION_NAME, key.UserId, key.Id, key.Id)
    }
    if key.ConsumerId != "" {
        p.addToStringMapCollection(key.ConsumerId, _INMEMORY_ACCESS_KEYS_FOR_CONSUMER_ID_COLLECTION_NAME, key.ConsumerId, key.Id, key.Id)
    }
    if err := key.AfterSave(); err != nil {
        return key, err
    }
    if err := key.AfterCreate(); err != nil {
        return key, err
    }
    return key, nil
}

func (p *InMemoryDataStore) DeleteUserPassword(userId string) (*dm.UserPassword, error) {
    oldValue, _ := p.delete(_INMEMORY_USER_PASSWORD_COLLECTION_NAME, userId)
    if oldValue != nil {
        pwd, _ := oldValue.(*dm.UserPassword)
        return pwd, nil
    }
    return nil, nil
}

func (p *InMemoryDataStore) DeleteAccessKey(accessKeyId string) (oldKey *dm.AccessKey, err error) {
    oldValue, _ := p.delete(_INMEMORY_ACCESS_KEYS_COLLECTION_NAME, accessKeyId)
    if oldValue != nil {
        if key, ok := oldValue.(*dm.AccessKey); ok {
            uid := key.UserId
            colName := _INMEMORY_ACCESS_KEYS_FOR_USER_ID_COLLECTION_NAME
            if uid == "" {
                uid = key.ConsumerId
                colName = _INMEMORY_ACCESS_KEYS_FOR_CONSUMER_ID_COLLECTION_NAME
            }
            p.removeFromStringMapCollection(uid, colName, uid, accessKeyId)
            oldKey = key
        }
    }
    return
}

func (p *InMemoryDataStore) BackendAuthenticationDataStore() ba.DataStore {
    return p
}
