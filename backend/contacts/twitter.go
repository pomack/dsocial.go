package contacts

import (
    "container/list"
    "github.com/pomack/contacts.go/twitter"
    dm "github.com/pomack/dsocial.go/models/dsocial"
    "github.com/pomack/jsonhelper.go/jsonhelper"
    "github.com/pomack/oauth2_client.go/oauth2_client"
)

type TwitterContactServiceSettings struct {
    StandardContactServiceSettings `json:"settings,omitempty,collapse"`
}

func NewTwitterContactServiceSettings() *TwitterContactServiceSettings {
    s := new(TwitterContactServiceSettings)
    s.SetAllowRetrieveContactInfo(true)
    s.SetAllowModifyContactInfo(true)
    return s
}

func (p *TwitterContactServiceSettings) ContactsServiceId() string {
    return TWITTER_CONTACT_SERVICE_ID
}

type TwitterContactService struct {
}

func NewTwitterContactService() *TwitterContactService {
    return new(TwitterContactService)
}

func (p *TwitterContactService) ServiceId() string {
    return TWITTER_CONTACT_SERVICE_ID
}

func (p *TwitterContactService) CreateOAuth2Client(settings jsonhelper.JSONObject) (client oauth2_client.OAuth2Client, err error) {
    client = oauth2_client.NewTwitterClient()
    client.Initialize(settings)
    return
}

func (p *TwitterContactService) ConvertToDsocialContact(externalContact interface{}, originalDsocialContact *dm.Contact, dsocialUserId string) (dsocialContact *dm.Contact) {
    if externalContact == nil {
        return
    }
    if extContact, ok := externalContact.(*twitter.User); ok && extContact != nil {
        dsocialContact = dm.TwitterUserToDsocial(extContact, originalDsocialContact, dsocialUserId)
    }
    return
}

func (p *TwitterContactService) ConvertToExternalContact(dsocialContact *dm.Contact, originalExternalContact interface{}, dsocialUserId string) (externalContact interface{}) {
    var origTwitterContact *twitter.User = nil
    if originalExternalContact != nil {
        origTwitterContact, _ = originalExternalContact.(*twitter.User)
    }
    externalContact = dm.DsocialContactToTwitter(dsocialContact, origTwitterContact)
    return
}

func (p *TwitterContactService) ConvertToDsocialGroup(externalGroup interface{}, originalDsocialGroup *dm.Group, dsocialUserId string) (dsocialGroup *dm.Group) {
    return
}

func (p *TwitterContactService) ConvertToExternalGroup(dsocialGroup *dm.Group, originalExternalGroup interface{}, dsocialUserId string) (externalGroup interface{}) {
    return
}

func (p *TwitterContactService) CanRetrieveAllContacts() bool {
    return false
}

func (p *TwitterContactService) CanRetrieveAllConnections() bool {
    return false
}

func (p *TwitterContactService) CanRetrieveAllGroups() bool {
    return false
}

func (p *TwitterContactService) CanRetrieveContacts() bool {
    return false
}

func (p *TwitterContactService) CanRetrieveConnections() bool {
    return false
}

func (p *TwitterContactService) CanRetrieveGroups() bool {
    return false
}

func (p *TwitterContactService) CanRetrieveContact(selfContact bool) bool {
    return true
}

func (p *TwitterContactService) CanImportContactsOrGroups() bool {
    return true
}

func (p *TwitterContactService) CanExportContactsOrGroups() bool {
    return false
}

func (p *TwitterContactService) CanCreateContact(selfContact bool) bool {
    return false
}

func (p *TwitterContactService) CanUpdateContact(selfContact bool) bool {
    return false
}

func (p *TwitterContactService) CanDeleteContact(selfContact bool) bool {
    return false
}

func (p *TwitterContactService) CanRetrieveGroup(selfContact bool) bool {
    return false
}

func (p *TwitterContactService) CanCreateGroup(selfContact bool) bool {
    return false
}

func (p *TwitterContactService) CanUpdateGroup(selfContact bool) bool {
    return false
}

func (p *TwitterContactService) CanDeleteGroup(selfContact bool) bool {
    return false
}

func (p *TwitterContactService) GroupListIncludesContactIds() bool {
    return false
}

func (p *TwitterContactService) GroupInfoIncludesContactIds() bool {
    return false
}

func (p *TwitterContactService) ContactInfoIncludesGroups() bool {
    return false
}

func (p *TwitterContactService) RetrieveAllContacts(client oauth2_client.OAuth2Client, ds DataStoreService, dsocialUserId string) ([]*Contact, error) {
    return make([]*Contact, 0), nil
}

func (p *TwitterContactService) RetrieveAllConnections(client oauth2_client.OAuth2Client, ds DataStoreService, dsocialUserId string) ([]*Contact, error) {
    return make([]*Contact, 0), nil
}

func (p *TwitterContactService) RetrieveAllGroups(client oauth2_client.OAuth2Client, ds DataStoreService, dsocialUserId string) ([]*Group, error) {
    return make([]*Group, 0), nil
}

func (p *TwitterContactService) listToContacts(l *list.List) []*Contact {
    if l == nil {
        return make([]*Contact, 0)
    }
    arr := make([]*Contact, l.Len())
    for i, e := 0, l.Front(); e != nil; i, e = i+1, e.Next() {
        if c, ok := e.Value.(*Contact); ok {
            arr[i] = c
        }
    }
    return arr
}

func (p *TwitterContactService) RetrieveContacts(client oauth2_client.OAuth2Client, ds DataStoreService, dsocialUserId string, next NextToken) ([]*Contact, NextToken, error) {
    return make([]*Contact, 0), nil, nil
}

func (p *TwitterContactService) RetrieveConnections(client oauth2_client.OAuth2Client, ds DataStoreService, dsocialUserId string, next NextToken) ([]*Contact, NextToken, error) {
    return make([]*Contact, 0), nil, nil
}

func (p *TwitterContactService) RetrieveGroups(client oauth2_client.OAuth2Client, ds DataStoreService, dsocialUserId string, next NextToken) ([]*Group, NextToken, error) {
    return make([]*Group, 0), nil, nil
}

func (p *TwitterContactService) RetrieveContact(client oauth2_client.OAuth2Client, ds DataStoreService, dsocialUserId string, contactId string) (*Contact, error) {
    extContact, err := twitter.RetrieveUser(client, contactId, 0, true, nil)
    if extContact == nil || err != nil {
        return nil, err
    }
    return p.handleRetrievedContact(client, ds, dsocialUserId, extContact.IdStr, extContact)
}

func (p *TwitterContactService) handleRetrievedContact(client oauth2_client.OAuth2Client, ds DataStoreService, dsocialUserId string, contactId string, extContact *twitter.User) (contact *Contact, err error) {
    if extContact == nil {
        return nil, nil
    }
    externalServiceId := p.ServiceId()
    userInfo, err := client.RetrieveUserInfo()
    externalUserId := userInfo.Guid()
    var useErr error = nil
    dsocialContactId := ""
    var origDsocialContact *dm.Contact = nil
    externalContactId := extContact.IdStr
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
    dsocialContact := dm.TwitterUserToDsocial(extContact, origDsocialContact, dsocialUserId)
    contact = &Contact{
        ExternalServiceId: p.ServiceId(),
        ExternalUserId:    externalUserId,
        ExternalContactId: externalContactId,
        DsocialUserId:     dsocialUserId,
        DsocialContactId:  dsocialContactId,
        Value:             dsocialContact,
    }
    return contact, useErr
}

func (p *TwitterContactService) RetrieveGroup(client oauth2_client.OAuth2Client, ds DataStoreService, dsocialUserId string, groupId string) (*Group, error) {
    return nil, nil
}

func (p *TwitterContactService) CreateContactOnExternalService(client oauth2_client.OAuth2Client, contact interface{}) (interface{}, string, error) {
    return nil, "", nil
}

func (p *TwitterContactService) CreateGroupOnExternalService(client oauth2_client.OAuth2Client, group interface{}) (interface{}, string, error) {
    return nil, "", nil
}

func (p *TwitterContactService) UpdateContactOnExternalService(client oauth2_client.OAuth2Client, originalContact, contact interface{}) (interface{}, string, error) {
    return nil, "", nil
}

func (p *TwitterContactService) UpdateGroupOnExternalService(client oauth2_client.OAuth2Client, originalGroup, group interface{}) (interface{}, string, error) {
    return nil, "", nil
}

func (p *TwitterContactService) DeleteContactOnExternalService(client oauth2_client.OAuth2Client, contact interface{}) (bool, error) {
    return false, nil
}

func (p *TwitterContactService) DeleteGroupOnExternalService(client oauth2_client.OAuth2Client, group interface{}) (bool, error) {
    return false, nil
}

func (p *TwitterContactService) ContactsService() ContactsService {
    return p
}
