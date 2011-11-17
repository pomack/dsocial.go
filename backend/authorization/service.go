package authorization

import (
    dm "github.com/pomack/dsocial.go/models/dsocial"
    "os"
)


type NextToken interface{}

type DataStore interface {
    SetSession(session *dm.Session) (*dm.Session, os.Error)
    DeleteSession(sessionId string) (os.Error)
    RetrieveSession(sessionId string) (*dm.Session, os.Error)
    RetrieveSessionsForUserId(userId string, next NextToken, maxResults int) ([]*dm.Session, NextToken, os.Error)
    RetrieveSessionsForConsumerId(consumerId string, next NextToken, maxResults int) ([]*dm.Session, NextToken, os.Error)
    RetrieveSessionsForExternalUserId(externalUserId string, next NextToken, maxResults int) ([]*dm.Session, NextToken, os.Error)
    
    SetAuthorizationToken(authToken *dm.AuthorizationToken) (*dm.AuthorizationToken, os.Error)
    DeleteAuthorizationToken(authTokenId string) (os.Error)
    RetrieveAuthorizationToken(authTokenId string) (*dm.AuthorizationToken, os.Error)
    RetrieveAuthorizationTokensForUserId(userId string, next NextToken, maxResults int) ([]*dm.AuthorizationToken, NextToken, os.Error)
    RetrieveAuthorizationTokensForConsumerId(consumerId string, next NextToken, maxResults int) ([]*dm.AuthorizationToken, NextToken, os.Error)
    RetrieveAuthorizationTokensForExternalUserId(externalUserId string, next NextToken, maxResults int) ([]*dm.AuthorizationToken, NextToken, os.Error)
}

