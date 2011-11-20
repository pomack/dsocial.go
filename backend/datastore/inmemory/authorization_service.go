package inmemory

import (
    dm "github.com/pomack/dsocial.go/models/dsocial"
    ba "github.com/pomack/dsocial.go/backend/authorization"
    "os"
)

func (p *InMemoryDataStore) RetrieveSession(sessionId string) (session *dm.Session, err os.Error) {
    v, _ := p.retrieve(_INMEMORY_SESSIONS_COLLECTION_NAME, sessionId)
    if v != nil {
        session, _ = v.(*dm.Session)
    }
    return
}

func (p *InMemoryDataStore) StoreSession(session *dm.Session) (*dm.Session, os.Error) {
    if session == nil {
        return nil, nil
    }
    uid := session.UID()
    if session.Id == "" {
        session.Id = p.GenerateId(uid, _INMEMORY_SESSIONS_COLLECTION_NAME)
    }
    p.store(uid, _INMEMORY_SESSIONS_COLLECTION_NAME, session.Id, session)
    p.addToStringMapCollection(uid, _INMEMORY_SESSION_IDS_FOR_USER_ID_COLLECTION_NAME, session.UserId, session.Id, session.Id)
    p.addToStringMapCollection(uid, _INMEMORY_SESSION_IDS_FOR_CONSUMER_ID_COLLECTION_NAME, session.ConsumerId, session.Id, session.Id)
    p.addToStringMapCollection(uid, _INMEMORY_SESSION_IDS_FOR_EXTERNAL_USER_ID_COLLECTION_NAME, session.ExternalUserId, session.Id, session.Id)
    return session, nil
}

func (p *InMemoryDataStore) DeleteSession(sessionId string) os.Error {
    oldValue, _ := p.delete(_INMEMORY_SESSIONS_COLLECTION_NAME, sessionId)
    if oldValue != nil {
        session, _ := oldValue.(*dm.Session)
        if session != nil {
            uid := session.UID()
            p.removeFromStringMapCollection(uid, _INMEMORY_SESSION_IDS_FOR_USER_ID_COLLECTION_NAME, session.UserId, session.Id)
            p.removeFromStringMapCollection(uid, _INMEMORY_SESSION_IDS_FOR_CONSUMER_ID_COLLECTION_NAME, session.ConsumerId, session.Id)
            p.removeFromStringMapCollection(uid, _INMEMORY_SESSION_IDS_FOR_EXTERNAL_USER_ID_COLLECTION_NAME, session.ExternalUserId, session.Id)
        }
    }
    return nil
}

func (p *InMemoryDataStore) RetrieveSessionsForUserId(userId string, next ba.NextToken, maxResults int) ([]*dm.Session, ba.NextToken, os.Error) {
    m := p.retrieveStringMapCollection(userId, _INMEMORY_SESSION_IDS_FOR_USER_ID_COLLECTION_NAME, userId)
    sessions := make([]*dm.Session, len(m))
    i := 0
    var err os.Error
    for k := range m {
        sessions[i], err = p.RetrieveSession(k)
        i++
        if err != nil {
            break
        }
    }
    return sessions, nil, err
}

func (p *InMemoryDataStore) RetrieveSessionsForConsumerId(consumerId string, next ba.NextToken, maxResults int) ([]*dm.Session, ba.NextToken, os.Error) {
    m := p.retrieveStringMapCollection(consumerId, _INMEMORY_SESSION_IDS_FOR_CONSUMER_ID_COLLECTION_NAME, consumerId)
    sessions := make([]*dm.Session, len(m))
    i := 0
    var err os.Error
    for k := range m {
        sessions[i], err = p.RetrieveSession(k)
        i++
        if err != nil {
            break
        }
    }
    return sessions, nil, err
}

func (p *InMemoryDataStore) RetrieveSessionsForExternalUserId(externalUserId string, next ba.NextToken, maxResults int) ([]*dm.Session, ba.NextToken, os.Error) {
    m := p.retrieveStringMapCollection(externalUserId, _INMEMORY_SESSION_IDS_FOR_EXTERNAL_USER_ID_COLLECTION_NAME, externalUserId)
    sessions := make([]*dm.Session, len(m))
    i := 0
    var err os.Error
    for k := range m {
        sessions[i], err = p.RetrieveSession(k)
        i++
        if err != nil {
            break
        }
    }
    return sessions, nil, err
}

func (p *InMemoryDataStore) RetrieveAuthorizationToken(authTokenId string) (authToken *dm.AuthorizationToken, err os.Error) {
    v, _ := p.retrieve(_INMEMORY_AUTH_TOKENS_COLLECTION_NAME, authTokenId)
    if v != nil {
        authToken, _ = v.(*dm.AuthorizationToken)
    }
    return
}

func (p *InMemoryDataStore) StoreAuthorizationToken(authToken *dm.AuthorizationToken) (*dm.AuthorizationToken, os.Error) {
    if authToken == nil {
        return nil, nil
    }
    uid := authToken.UID()
    if authToken.Id == "" {
        authToken.Id = p.GenerateId(uid, _INMEMORY_AUTH_TOKENS_COLLECTION_NAME)
    }
    p.store(uid, _INMEMORY_AUTH_TOKENS_COLLECTION_NAME, authToken.Id, authToken)
    p.addToStringMapCollection(uid, _INMEMORY_AUTH_TOKEN_IDS_FOR_USER_ID_COLLECTION_NAME, authToken.UserId, authToken.Id, authToken.Id)
    p.addToStringMapCollection(uid, _INMEMORY_AUTH_TOKEN_IDS_FOR_CONSUMER_ID_COLLECTION_NAME, authToken.ConsumerId, authToken.Id, authToken.Id)
    p.addToStringMapCollection(uid, _INMEMORY_AUTH_TOKEN_IDS_FOR_EXTERNAL_USER_ID_COLLECTION_NAME, authToken.ExternalUserId, authToken.Id, authToken.Id)
    return authToken, nil
}

func (p *InMemoryDataStore) DeleteAuthorizationToken(authTokenId string) os.Error {
    oldValue, _ := p.delete(_INMEMORY_AUTH_TOKENS_COLLECTION_NAME, authTokenId)
    if oldValue != nil {
        authToken, _ := oldValue.(*dm.AuthorizationToken)
        if authToken != nil {
            uid := authToken.UID()
            p.removeFromStringMapCollection(uid, _INMEMORY_AUTH_TOKEN_IDS_FOR_USER_ID_COLLECTION_NAME, authToken.UserId, authToken.Id)
            p.removeFromStringMapCollection(uid, _INMEMORY_AUTH_TOKEN_IDS_FOR_CONSUMER_ID_COLLECTION_NAME, authToken.ConsumerId, authToken.Id)
            p.removeFromStringMapCollection(uid, _INMEMORY_AUTH_TOKEN_IDS_FOR_EXTERNAL_USER_ID_COLLECTION_NAME, authToken.ExternalUserId, authToken.Id)
        }
    }
    return nil
}

func (p *InMemoryDataStore) RetrieveAuthorizationTokensForUserId(userId string, next ba.NextToken, maxResults int) ([]*dm.AuthorizationToken, ba.NextToken, os.Error) {
    m := p.retrieveStringMapCollection(userId, _INMEMORY_AUTH_TOKEN_IDS_FOR_USER_ID_COLLECTION_NAME, userId)
    authTokens := make([]*dm.AuthorizationToken, len(m))
    i := 0
    var err os.Error
    for k := range m {
        authTokens[i], err = p.RetrieveAuthorizationToken(k)
        i++
        if err != nil {
            break
        }
    }
    return authTokens, nil, err
}

func (p *InMemoryDataStore) RetrieveAuthorizationTokensForConsumerId(consumerId string, next ba.NextToken, maxResults int) ([]*dm.AuthorizationToken, ba.NextToken, os.Error) {
    m := p.retrieveStringMapCollection(consumerId, _INMEMORY_AUTH_TOKEN_IDS_FOR_CONSUMER_ID_COLLECTION_NAME, consumerId)
    authTokens := make([]*dm.AuthorizationToken, len(m))
    i := 0
    var err os.Error
    for k := range m {
        authTokens[i], err = p.RetrieveAuthorizationToken(k)
        i++
        if err != nil {
            break
        }
    }
    return authTokens, nil, err
}

func (p *InMemoryDataStore) RetrieveAuthorizationTokensForExternalUserId(externalUserId string, next ba.NextToken, maxResults int) ([]*dm.AuthorizationToken, ba.NextToken, os.Error) {
    m := p.retrieveStringMapCollection(externalUserId, _INMEMORY_AUTH_TOKEN_IDS_FOR_EXTERNAL_USER_ID_COLLECTION_NAME, externalUserId)
    authTokens := make([]*dm.AuthorizationToken, len(m))
    i := 0
    var err os.Error
    for k := range m {
        authTokens[i], err = p.RetrieveAuthorizationToken(k)
        i++
        if err != nil {
            break
        }
    }
    return authTokens, nil, err
}

func (p *InMemoryDataStore) BackendAuthorizationDataStore() ba.DataStore {
    return p
}
