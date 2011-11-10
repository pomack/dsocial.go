package contacts

import (
    "github.com/pomack/oauth2_client.go/oauth2_client"
    "github.com/pomack/contacts.go/facebook"
    dm "github.com/pomack/dsocial.go/models/dsocial"
    "container/list"
    "os"
)


type FacebookContactService struct {

}

func NewFacebookContactService() *FacebookContactService {
    return new(FacebookContactService)
}

func (p *FacebookContactService) ServiceId() string {
    return "www.facebook.com"
}

func (p *FacebookContactService) ConvertToDsocialContact(externalContact interface{}, originalDsocialContact *dm.Contact, dsocialUserId string) (dsocialContact *dm.Contact) {
    if externalContact == nil {
        return
    }
    if extContact, ok := externalContact.(*facebook.Contact); ok && extContact != nil {
        dsocialContact = dm.FacebookContactToDsocial(extContact, originalDsocialContact, dsocialUserId)
    }
    return
}

func (p *FacebookContactService) ConvertToExternalContact(dsocialContact *dm.Contact, originalExternalContact interface{}, dsocialUserId string) (externalContact interface{}) {
    var origFacebookContact *facebook.Contact = nil
    if originalExternalContact != nil {
        origFacebookContact, _ = originalExternalContact.(*facebook.Contact)
    }
    externalContact = dm.DsocialContactToFacebook(dsocialContact, origFacebookContact)
    return
}

func (p *FacebookContactService) ConvertToDsocialGroup(externalGroup interface{}, originalDsocialGroup *dm.Group, dsocialUserId string) (dsocialGroup *dm.Group) {
    return
}

func (p *FacebookContactService) ConvertToExternalGroup(dsocialGroup *dm.Group, originalExternalGroup interface{}, dsocialUserId string) (externalGroup interface{}) {
    return
}


func (p *FacebookContactService) CanRetrieveAllContacts() bool {
    return false
}

func (p *FacebookContactService) CanRetrieveAllConnections() bool {
    return false
}

func (p *FacebookContactService) CanRetrieveAllGroups() bool {
    return false
}

func (p *FacebookContactService) CanRetrieveContacts() bool {
    return false
}

func (p *FacebookContactService) CanRetrieveConnections() bool {
    return false
}

func (p *FacebookContactService) CanRetrieveGroups() bool {
    return false
}

func (p *FacebookContactService) CanRetrieveContact(selfContact bool) bool {
    return true
}

func (p *FacebookContactService) CanCreateContact(selfContact bool) bool {
    return false
}

func (p *FacebookContactService) CanUpdateContact(selfContact bool) bool {
    return false
}

func (p *FacebookContactService) CanDeleteContact(selfContact bool) bool {
    return false
}

func (p *FacebookContactService) CanRetrieveGroup(selfContact bool) bool {
    return false
}

func (p *FacebookContactService) CanCreateGroup(selfContact bool) bool {
    return false
}

func (p *FacebookContactService) CanUpdateGroup(selfContact bool) bool {
    return false
}

func (p *FacebookContactService) CanDeleteGroup(selfContact bool) bool {
    return false
}

func (p *FacebookContactService) GroupListIncludesContactIds() bool {
    return false
}

func (p *FacebookContactService) GroupInfoIncludesContactIds() bool {
    return false
}

func (p *FacebookContactService) ContactInfoIncludesGroups() bool {
    return false
}

func (p *FacebookContactService) RetrieveAllContacts(client oauth2_client.OAuth2Client, ds DataStoreService, dsocialUserId string) ([]*Contact, os.Error) {
    return make([]*Contact, 0), nil
}

func (p *FacebookContactService) RetrieveAllConnections(client oauth2_client.OAuth2Client, ds DataStoreService, dsocialUserId string) ([]*Contact, os.Error) {
    return make([]*Contact, 0), nil
}

func (p *FacebookContactService) RetrieveAllGroups(client oauth2_client.OAuth2Client, ds DataStoreService, dsocialUserId string) ([]*Group, os.Error) {
    return make([]*Group, 0), nil
}

func (p *FacebookContactService) listToContacts(l *list.List) []*Contact {
    if l == nil {
        return make([]*Contact, 0)
    }
    arr := make([]*Contact, l.Len())
    for i, e := 0, l.Front(); e != nil; i, e = i + 1, e.Next() {
        if c, ok := e.Value.(*Contact); ok {
            arr[i] = c
        }
    }
    return arr
}

func (p *FacebookContactService) RetrieveContacts(client oauth2_client.OAuth2Client, ds DataStoreService, dsocialUserId string, next NextToken) ([]*Contact, NextToken, os.Error) {
    return make([]*Contact, 0), nil, nil
}

func (p *FacebookContactService) RetrieveConnections(client oauth2_client.OAuth2Client, ds DataStoreService, dsocialUserId string, next NextToken) ([]*Contact, NextToken, os.Error) {
    return make([]*Contact, 0), nil, nil
}

func (p *FacebookContactService) RetrieveGroups(client oauth2_client.OAuth2Client, ds DataStoreService, dsocialUserId string, next NextToken) ([]*Group, NextToken, os.Error) {
    return make([]*Group, 0), nil, nil
}

func (p *FacebookContactService) RetrieveContact(client oauth2_client.OAuth2Client, ds DataStoreService, dsocialUserId string, contactId string) (*Contact, os.Error) {
    extContact, err := facebook.RetrieveContact(client, contactId)
    if extContact == nil || err != nil {
        return nil, err
    }
    return p.handleRetrievedContact(client, ds, dsocialUserId, extContact.Id, extContact)
}

func (p *FacebookContactService) handleRetrievedContact(client oauth2_client.OAuth2Client, ds DataStoreService, dsocialUserId string, contactId string, extContact *facebook.Contact) (contact *Contact, err os.Error) {
    if extContact == nil {
        return nil, nil
    }
    externalServiceId := p.ServiceId()
    userInfo, err := client.RetrieveUserInfo()
    externalUserId := userInfo.Guid()
    var useErr os.Error = nil
    dsocialContactId := ""
    var origDsocialContact *dm.Contact = nil
    externalContactId := extContact.Id
    if len(externalContactId) > 0 {
        dsocialContactId, err = ds.DsocialIdForExternalContactId(externalServiceId, externalUserId, dsocialUserId, contactId)
        if err != nil {
            useErr = err
        }
        if dsocialContactId != "" {
            origDsocialContact, _, err = ds.RetrieveDsocialContactForExternalContact(externalServiceId, externalUserId, externalContactId, dsocialUserId)
            if err != nil && useErr == nil {
                useErr = err
            }
        } else {
            ds.StoreExternalContact(externalServiceId, externalUserId, dsocialUserId, externalContactId, extContact)
        }
    }
    dsocialContact := dm.FacebookContactToDsocial(extContact, origDsocialContact, dsocialUserId)
    contact = &Contact{
        ExternalServiceId: p.ServiceId(),
        ExternalUserId: externalUserId,
        ExternalContactId: externalContactId,
        DsocialUserId: dsocialUserId,
        DsocialContactId: dsocialContactId,
        Value: dsocialContact,
    }
    return contact, useErr
}

func (p *FacebookContactService) RetrieveGroup(client oauth2_client.OAuth2Client, ds DataStoreService, dsocialUserId string, groupId string) (*Group, os.Error) {
    return nil, nil
}

func (p *FacebookContactService) CreateContact(client oauth2_client.OAuth2Client, ds DataStoreService, dsocialUserId string, contact *dm.Contact) (*Contact, os.Error) {
    return nil, nil
}

func (p *FacebookContactService) CreateGroup(client oauth2_client.OAuth2Client, ds DataStoreService, dsocialUserId string, group *dm.Group) (*Group, os.Error) {
    return nil, nil
}

func (p *FacebookContactService) UpdateContact(client oauth2_client.OAuth2Client, ds DataStoreService, dsocialUserId string, originalContact, contact *dm.Contact) (*Contact, os.Error) {
    return nil, nil
}

func (p *FacebookContactService) UpdateGroup(client oauth2_client.OAuth2Client, ds DataStoreService, dsocialUserId string, originalGroup, group *dm.Group) (*Group, os.Error) {
    return nil, nil
}

func (p *FacebookContactService) DeleteContact(client oauth2_client.OAuth2Client, ds DataStoreService, dsocialUserId, dsocialContactId string) (bool, os.Error) {
    return false, nil
}

func (p *FacebookContactService) DeleteGroup(client oauth2_client.OAuth2Client, ds DataStoreService, dsocialUserId, dsocialGroupId string) (bool, os.Error) {
    return false, nil
}

func (p *FacebookContactService) ContactsService() ContactsService {
    return p
}
