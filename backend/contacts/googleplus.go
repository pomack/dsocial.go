package contacts

import (
    "github.com/pomack/contacts.go/googleplus"
    dm "github.com/pomack/dsocial.go/models/dsocial"
    "github.com/pomack/jsonhelper.go/jsonhelper"
    "github.com/pomack/oauth2_client.go/oauth2_client"
)

type GooglePlusContactService struct {
}

type GooglePlusContactServiceSettings struct {
    StandardContactServiceSettings `json:"settings,omitempty,collapse"`
}

func NewGooglePlusContactServiceSettings() *GooglePlusContactServiceSettings {
    s := new(GooglePlusContactServiceSettings)
    s.SetAllowRetrieveContactInfo(true)
    s.SetAllowModifyContactInfo(true)
    return s
}

func (p *GooglePlusContactServiceSettings) ContactsServiceId() string {
    return GOOGLE_PLUS_CONTACT_SERVICE_ID
}

func NewGooglePlusContactService() *GooglePlusContactService {
    return new(GooglePlusContactService)
}

func (p *GooglePlusContactService) ServiceId() string {
    return GOOGLE_PLUS_CONTACT_SERVICE_ID
}

func (p *GooglePlusContactService) CreateOAuth2Client(settings jsonhelper.JSONObject) (client oauth2_client.OAuth2Client, err error) {
    client = oauth2_client.NewGooglePlusClient()
    client.Initialize(settings)
    return
}

func (p *GooglePlusContactService) ConvertToDsocialContact(externalContact interface{}, originalDsocialContact *dm.Contact, dsocialUserId string) (dsocialContact *dm.Contact) {
    if externalContact == nil {
        return
    }
    if extContact, ok := externalContact.(*googleplus.Person); ok && extContact != nil {
        dsocialContact = dm.GooglePlusPersonToDsocial(extContact, originalDsocialContact, dsocialUserId)
    }
    return
}

func (p *GooglePlusContactService) ConvertToExternalContact(dsocialContact *dm.Contact, originalExternalContact interface{}, dsocialUserId string) (externalContact interface{}) {
    var origGooglePlusContact *googleplus.Person = nil
    if originalExternalContact != nil {
        origGooglePlusContact, _ = originalExternalContact.(*googleplus.Person)
    }
    externalContact = dm.DsocialContactToGooglePlus(dsocialContact, origGooglePlusContact)
    return
}

func (p *GooglePlusContactService) ConvertToDsocialGroup(externalGroup interface{}, originalDsocialGroup *dm.Group, dsocialUserId string) (dsocialGroup *dm.Group) {
    return
}

func (p *GooglePlusContactService) ConvertToExternalGroup(dsocialGroup *dm.Group, originalExternalGroup interface{}, dsocialUserId string) (externalGroup interface{}) {
    return
}

func (p *GooglePlusContactService) CanRetrieveAllContacts() bool {
    return false
}

func (p *GooglePlusContactService) CanRetrieveAllConnections() bool {
    return false
}

func (p *GooglePlusContactService) CanRetrieveAllGroups() bool {
    return false
}

func (p *GooglePlusContactService) CanRetrieveContacts() bool {
    return false
}

func (p *GooglePlusContactService) CanRetrieveConnections() bool {
    return false
}

func (p *GooglePlusContactService) CanRetrieveGroups() bool {
    return false
}

func (p *GooglePlusContactService) CanRetrieveContact(selfContact bool) bool {
    return true
}

func (p *GooglePlusContactService) CanImportContactsOrGroups() bool {
    return true
}

func (p *GooglePlusContactService) CanExportContactsOrGroups() bool {
    return false
}

func (p *GooglePlusContactService) CanCreateContact(selfContact bool) bool {
    return false
}

func (p *GooglePlusContactService) CanUpdateContact(selfContact bool) bool {
    return false
}

func (p *GooglePlusContactService) CanDeleteContact(selfContact bool) bool {
    return false
}

func (p *GooglePlusContactService) CanRetrieveGroup(selfContact bool) bool {
    return false
}

func (p *GooglePlusContactService) CanCreateGroup(selfContact bool) bool {
    return false
}

func (p *GooglePlusContactService) CanUpdateGroup(selfContact bool) bool {
    return false
}

func (p *GooglePlusContactService) CanDeleteGroup(selfContact bool) bool {
    return false
}

func (p *GooglePlusContactService) GroupListIncludesContactIds() bool {
    return false
}

func (p *GooglePlusContactService) GroupInfoIncludesContactIds() bool {
    return false
}

func (p *GooglePlusContactService) ContactInfoIncludesGroups() bool {
    return false
}

func (p *GooglePlusContactService) RetrieveAllContacts(client oauth2_client.OAuth2Client, ds DataStoreService, dsocialUserId string) ([]*Contact, error) {
    contacts, _, err := p.RetrieveContacts(client, ds, dsocialUserId, nil)
    return contacts, err
}

func (p *GooglePlusContactService) RetrieveAllConnections(client oauth2_client.OAuth2Client, ds DataStoreService, dsocialUserId string) ([]*Contact, error) {
    return make([]*Contact, 0), nil
}

func (p *GooglePlusContactService) RetrieveAllGroups(client oauth2_client.OAuth2Client, ds DataStoreService, dsocialUserId string) ([]*Group, error) {
    return make([]*Group, 0), nil
}

func (p *GooglePlusContactService) RetrieveContacts(client oauth2_client.OAuth2Client, ds DataStoreService, dsocialUserId string, next NextToken) ([]*Contact, NextToken, error) {
    return make([]*Contact, 0), nil, nil
}

func (p *GooglePlusContactService) RetrieveConnections(client oauth2_client.OAuth2Client, ds DataStoreService, dsocialUserId string, next NextToken) ([]*Contact, NextToken, error) {
    return make([]*Contact, 0), nil, nil
}

func (p *GooglePlusContactService) RetrieveGroups(client oauth2_client.OAuth2Client, ds DataStoreService, dsocialUserId string, next NextToken) ([]*Group, NextToken, error) {
    return make([]*Group, 0), nil, nil
}

func (p *GooglePlusContactService) RetrieveContact(client oauth2_client.OAuth2Client, ds DataStoreService, dsocialUserId string, contactId string) (*Contact, error) {
    googleplusContact, err := googleplus.RetrieveContact(client, contactId, nil)
    if googleplusContact == nil || err != nil {
        return nil, err
    }
    externalServiceId := p.ServiceId()
    userInfo, err := client.RetrieveUserInfo()
    externalUserId := userInfo.Guid()
    useErr := err
    dsocialContactId := ""
    var origDsocialContact *dm.Contact = nil
    externalContactId := googleplusContact.Id
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
            ds.StoreExternalContact(externalServiceId, externalUserId, dsocialUserId, externalContactId, googleplusContact)
        }
    }
    dsocialContact := dm.GooglePlusPersonToDsocial(googleplusContact, origDsocialContact, dsocialUserId)
    contact := &Contact{
        ExternalServiceId: p.ServiceId(),
        ExternalUserId:    externalUserId,
        ExternalContactId: googleplusContact.Id,
        DsocialUserId:     dsocialUserId,
        DsocialContactId:  dsocialContactId,
        Value:             dsocialContact,
    }
    return contact, useErr
}

func (p *GooglePlusContactService) RetrieveGroup(client oauth2_client.OAuth2Client, ds DataStoreService, dsocialUserId string, groupId string) (*Group, error) {
    return nil, nil
}

func (p *GooglePlusContactService) CreateContactOnExternalService(client oauth2_client.OAuth2Client, contact interface{}) (interface{}, string, error) {
    return nil, "", nil
}

func (p *GooglePlusContactService) CreateGroupOnExternalService(client oauth2_client.OAuth2Client, group interface{}) (interface{}, string, error) {
    return nil, "", nil
}

func (p *GooglePlusContactService) UpdateContactOnExternalService(client oauth2_client.OAuth2Client, originalContact, contact interface{}) (interface{}, string, error) {
    return nil, "", nil
}

func (p *GooglePlusContactService) UpdateGroupOnExternalService(client oauth2_client.OAuth2Client, originalGroup, group interface{}) (interface{}, string, error) {
    return nil, "", nil
}

func (p *GooglePlusContactService) DeleteContactOnExternalService(client oauth2_client.OAuth2Client, contact interface{}) (bool, error) {
    return false, nil
}

func (p *GooglePlusContactService) DeleteGroupOnExternalService(client oauth2_client.OAuth2Client, group interface{}) (bool, error) {
    return false, nil
}

func (p *GooglePlusContactService) ContactsService() ContactsService {
    return p
}
