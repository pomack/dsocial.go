package authorization

import dm "github.com/pomack/dsocial.go/models/dsocial"

type NextToken interface{}

type DataStore interface {
    RetrieveSession(sessionId string) (*dm.Session, error)
    StoreSession(session *dm.Session) (*dm.Session, error)
    DeleteSession(sessionId string) error
    RetrieveSessionsForUserId(userId string, next NextToken, maxResults int) ([]*dm.Session, NextToken, error)
    RetrieveSessionsForConsumerId(consumerId string, next NextToken, maxResults int) ([]*dm.Session, NextToken, error)
    RetrieveSessionsForExternalUserId(externalUserId string, next NextToken, maxResults int) ([]*dm.Session, NextToken, error)

    RetrieveAuthorizationToken(authTokenId string) (*dm.AuthorizationToken, error)
    StoreAuthorizationToken(authToken *dm.AuthorizationToken) (*dm.AuthorizationToken, error)
    DeleteAuthorizationToken(authTokenId string) error
    RetrieveAuthorizationTokensForUserId(userId string, next NextToken, maxResults int) ([]*dm.AuthorizationToken, NextToken, error)
    RetrieveAuthorizationTokensForConsumerId(consumerId string, next NextToken, maxResults int) ([]*dm.AuthorizationToken, NextToken, error)
    RetrieveAuthorizationTokensForExternalUserId(externalUserId string, next NextToken, maxResults int) ([]*dm.AuthorizationToken, NextToken, error)
}
