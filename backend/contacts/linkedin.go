package contacts

import (
    "github.com/pomack/jsonhelper.go/jsonhelper"
    "github.com/pomack/oauth2_client.go/oauth2_client"
    "github.com/pomack/contacts.go/linkedin"
    dm "github.com/pomack/dsocial.go/models/dsocial"
    "container/list"
    "os"
    "strconv"
    "url"
)

const (
    _LINKEDIN_MAX_CONNECTIONS_PER_CALL = 500
    _LINKEDIN_MAX_CONNECTIONS_PER_CALL_STRING = "500"
)

type LinkedInContactService struct {

}

func NewLinkedInContactService() *LinkedInContactService {
    return new(LinkedInContactService)
}

func (p *LinkedInContactService) ServiceId() string {
    return "www.linkedin.com"
}

func (p *LinkedInContactService) CreateOAuth2Client(settings jsonhelper.JSONObject) (client oauth2_client.OAuth2Client, err os.Error) {
    client = oauth2_client.NewLinkedInClient()
    client.Initialize(settings)
    return
}

func (p *LinkedInContactService) ConvertToDsocialContact(externalContact interface{}, originalDsocialContact *dm.Contact, dsocialUserId string) (dsocialContact *dm.Contact) {
    if externalContact == nil {
        return
    }
    if extContact, ok := externalContact.(*linkedin.Contact); ok && extContact != nil {
        dsocialContact = dm.LinkedInContactToDsocial(extContact, originalDsocialContact, dsocialUserId)
    }
    return
}

func (p *LinkedInContactService) ConvertToExternalContact(dsocialContact *dm.Contact, originalExternalContact interface{}, dsocialUserId string) (externalContact interface{}) {
    var origLinkedInContact *linkedin.Contact = nil
    if originalExternalContact != nil {
        origLinkedInContact, _ = originalExternalContact.(*linkedin.Contact)
    }
    externalContact = dm.DsocialContactToLinkedIn(dsocialContact, origLinkedInContact)
    return
}

func (p *LinkedInContactService) ConvertToDsocialGroup(externalGroup interface{}, originalDsocialGroup *dm.Group, dsocialUserId string) (dsocialGroup *dm.Group) {
    return
}

func (p *LinkedInContactService) ConvertToExternalGroup(dsocialGroup *dm.Group, originalExternalGroup interface{}, dsocialUserId string) (externalGroup interface{}) {
    return
}


func (p *LinkedInContactService) CanRetrieveAllContacts() bool {
    return false
}

func (p *LinkedInContactService) CanRetrieveAllConnections() bool {
    return false
}

func (p *LinkedInContactService) CanRetrieveAllGroups() bool {
    return false
}

func (p *LinkedInContactService) CanRetrieveContacts() bool {
    return true
}

func (p *LinkedInContactService) CanRetrieveConnections() bool {
    return false
}

func (p *LinkedInContactService) CanRetrieveGroups() bool {
    return false
}

func (p *LinkedInContactService) CanRetrieveContact(selfContact bool) bool {
    return true
}

func (p *LinkedInContactService) CanCreateContact(selfContact bool) bool {
    return false
}

func (p *LinkedInContactService) CanUpdateContact(selfContact bool) bool {
    return false
}

func (p *LinkedInContactService) CanDeleteContact(selfContact bool) bool {
    return false
}

func (p *LinkedInContactService) CanRetrieveGroup(selfContact bool) bool {
    return false
}

func (p *LinkedInContactService) CanCreateGroup(selfContact bool) bool {
    return false
}

func (p *LinkedInContactService) CanUpdateGroup(selfContact bool) bool {
    return false
}

func (p *LinkedInContactService) CanDeleteGroup(selfContact bool) bool {
    return false
}

func (p *LinkedInContactService) GroupListIncludesContactIds() bool {
    return false
}

func (p *LinkedInContactService) GroupInfoIncludesContactIds() bool {
    return false
}

func (p *LinkedInContactService) ContactInfoIncludesGroups() bool {
    return false
}

func (p *LinkedInContactService) RetrieveAllContacts(client oauth2_client.OAuth2Client, ds DataStoreService, dsocialUserId string) ([]*Contact, os.Error) {
    contacts, nextToken, err := p.RetrieveContacts(client, ds, dsocialUserId, nil)
    if contacts == nil || len(contacts) == 0 || nextToken == nil || err != nil {
        return contacts, err
    }
    l := list.New()
    for contacts != nil && len(contacts) > 0 {
        for _, c := range contacts {
            l.PushBack(c)
        }
        if err != nil {
            break
        }
        contacts, nextToken, err = p.RetrieveContacts(client, ds, dsocialUserId, nextToken)
    }
    contacts = make([]*Contact, l.Len())
    for i, e := 0, l.Front(); e != nil; i, e = i + 1, e.Next() {
        contacts[i] = e.Value.(*Contact)
    }
    return contacts, err
}

func (p *LinkedInContactService) RetrieveAllConnections(client oauth2_client.OAuth2Client, ds DataStoreService, dsocialUserId string) ([]*Contact, os.Error) {
    return make([]*Contact, 0), nil
}

func (p *LinkedInContactService) RetrieveAllGroups(client oauth2_client.OAuth2Client, ds DataStoreService, dsocialUserId string) ([]*Group, os.Error) {
    return make([]*Group, 0), nil
}

func (p *LinkedInContactService) RetrieveContacts(client oauth2_client.OAuth2Client, ds DataStoreService, dsocialUserId string, next NextToken) ([]*Contact, NextToken, os.Error) {
    m := make(url.Values)
    start := 0
    if next != nil {
        var ok bool
        if start, ok = next.(int); ok && start > 0 {
            m.Add("start", strconv.Itoa(start))
        } else {
            start = 0
        }
    }
    m.Add("count", _LINKEDIN_MAX_CONNECTIONS_PER_CALL_STRING)
    var outputNextToken NextToken = nil
    cl, err := linkedin.RetrieveConnections(client, nil, m)
    if cl == nil || cl.Values == nil || len(cl.Values) == 0 {
        return make([]*Contact, 0), nil, err
    }
    contacts := make([]*Contact, len(cl.Values))
    for i, extContact := range cl.Values {
        contacts[i], err = p.handleRetrievedContact(client, ds, dsocialUserId, extContact.Id, &extContact)
        if err != nil {
            break
        }
    }
    if len(contacts) == _LINKEDIN_MAX_CONNECTIONS_PER_CALL {
        outputNextToken = start + _LINKEDIN_MAX_CONNECTIONS_PER_CALL
    }
    return contacts, outputNextToken, err
}

func (p *LinkedInContactService) RetrieveConnections(client oauth2_client.OAuth2Client, ds DataStoreService, dsocialUserId string, next NextToken) ([]*Contact, NextToken, os.Error) {
    return make([]*Contact, 0), nil, nil
}

func (p *LinkedInContactService) RetrieveGroups(client oauth2_client.OAuth2Client, ds DataStoreService, dsocialUserId string, next NextToken) ([]*Group, NextToken, os.Error) {
    return make([]*Group, 0), nil, nil
}

func (p *LinkedInContactService) RetrieveContact(client oauth2_client.OAuth2Client, ds DataStoreService, dsocialUserId string, contactId string) (*Contact, os.Error) {
    extContact, err := linkedin.RetrieveProfile(client, contactId, nil, nil)
    if extContact == nil || err != nil {
        return nil, err
    }
    return p.handleRetrievedContact(client, ds, dsocialUserId, extContact.Id, extContact)
}

func (p *LinkedInContactService) handleRetrievedContact(client oauth2_client.OAuth2Client, ds DataStoreService, dsocialUserId string, contactId string, extContact *linkedin.Contact) (contact *Contact, err os.Error) {
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
    dsocialContact := dm.LinkedInContactToDsocial(extContact, origDsocialContact, dsocialUserId)
    contact = &Contact{
        ExternalServiceId: p.ServiceId(),
        ExternalUserId: externalUserId,
        ExternalContactId: extContact.Id,
        DsocialUserId: dsocialUserId,
        DsocialContactId: dsocialContactId,
        Value: dsocialContact,
    }
    return contact, useErr
}

func (p *LinkedInContactService) RetrieveGroup(client oauth2_client.OAuth2Client, ds DataStoreService, dsocialUserId string, groupId string) (*Group, os.Error) {
    return nil, nil
}

func (p *LinkedInContactService) CreateContact(client oauth2_client.OAuth2Client, ds DataStoreService, dsocialUserId string, contact *dm.Contact) (*Contact, os.Error) {
    return nil, nil
}

func (p *LinkedInContactService) CreateGroup(client oauth2_client.OAuth2Client, ds DataStoreService, dsocialUserId string, group *dm.Group) (*Group, os.Error) {
    return nil, nil
}

func (p *LinkedInContactService) UpdateContact(client oauth2_client.OAuth2Client, ds DataStoreService, dsocialUserId string, originalContact, contact *dm.Contact) (*Contact, os.Error) {
    return nil, nil
}

func (p *LinkedInContactService) UpdateGroup(client oauth2_client.OAuth2Client, ds DataStoreService, dsocialUserId string, originalGroup, group *dm.Group) (*Group, os.Error) {
    return nil, nil
}

func (p *LinkedInContactService) DeleteContact(client oauth2_client.OAuth2Client, ds DataStoreService, dsocialUserId, dsocialContactId string) (bool, os.Error) {
    return false, nil
}

func (p *LinkedInContactService) DeleteGroup(client oauth2_client.OAuth2Client, ds DataStoreService, dsocialUserId, dsocialGroupId string) (bool, os.Error) {
    return false, nil
}

func (p *LinkedInContactService) ContactsService() ContactsService {
    return p
}
