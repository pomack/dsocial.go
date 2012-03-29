package contacts

import (
    "container/list"
    "github.com/pomack/contacts.go/smugmug"
    dm "github.com/pomack/dsocial.go/models/dsocial"
    "github.com/pomack/jsonhelper.go/jsonhelper"
    "github.com/pomack/oauth2_client.go/oauth2_client"
)

type SmugMugContactServiceSettings struct {
    StandardContactServiceSettings `json:"settings,omitempty,collapse"`
}

func NewSmugMugContactServiceSettings() *SmugMugContactServiceSettings {
    s := new(SmugMugContactServiceSettings)
    s.SetAllowRetrieveContactInfo(true)
    s.SetAllowModifyContactInfo(true)
    return s
}

func (p *SmugMugContactServiceSettings) ContactsServiceId() string {
    return SMUGMUG_CONTACT_SERVICE_ID
}

type SmugMugContactService struct {
}

func NewSmugMugContactService() *SmugMugContactService {
    return new(SmugMugContactService)
}

func (p *SmugMugContactService) ServiceId() string {
    return SMUGMUG_CONTACT_SERVICE_ID
}

func (p *SmugMugContactService) CreateOAuth2Client(settings jsonhelper.JSONObject) (client oauth2_client.OAuth2Client, err error) {
    client = oauth2_client.NewSmugMugClient()
    client.Initialize(settings)
    return
}

func (p *SmugMugContactService) ConvertToDsocialContact(externalContact interface{}, originalDsocialContact *dm.Contact, dsocialUserId string) (dsocialContact *dm.Contact) {
    if externalContact == nil {
        return
    }
    if extContact, ok := externalContact.(*smugmug.PersonReference); ok && extContact != nil {
        dsocialContact = dm.SmugMugUserToDsocial(extContact, originalDsocialContact, dsocialUserId)
    }
    return
}

func (p *SmugMugContactService) ConvertToExternalContact(dsocialContact *dm.Contact, originalExternalContact interface{}, dsocialUserId string) (externalContact interface{}) {
    var origSmugMugContact *smugmug.PersonReference = nil
    if originalExternalContact != nil {
        origSmugMugContact, _ = originalExternalContact.(*smugmug.PersonReference)
    }
    externalContact = dm.DsocialContactToSmugMug(dsocialContact, origSmugMugContact)
    return
}

func (p *SmugMugContactService) ConvertToDsocialGroup(externalGroup interface{}, originalDsocialGroup *dm.Group, dsocialUserId string) (dsocialGroup *dm.Group) {
    return
}

func (p *SmugMugContactService) ConvertToExternalGroup(dsocialGroup *dm.Group, originalExternalGroup interface{}, dsocialUserId string) (externalGroup interface{}) {
    return
}

func (p *SmugMugContactService) CanRetrieveAllContacts() bool {
    return false
}

func (p *SmugMugContactService) CanRetrieveAllConnections() bool {
    return false
}

func (p *SmugMugContactService) CanRetrieveAllGroups() bool {
    return false
}

func (p *SmugMugContactService) CanRetrieveContacts() bool {
    return true
}

func (p *SmugMugContactService) CanRetrieveConnections() bool {
    return false
}

func (p *SmugMugContactService) CanRetrieveGroups() bool {
    return false
}

func (p *SmugMugContactService) CanRetrieveContact(selfContact bool) bool {
    return true
}

func (p *SmugMugContactService) CanImportContactsOrGroups() bool {
    return true
}

func (p *SmugMugContactService) CanExportContactsOrGroups() bool {
    return false
}

func (p *SmugMugContactService) CanCreateContact(selfContact bool) bool {
    return false
}

func (p *SmugMugContactService) CanUpdateContact(selfContact bool) bool {
    return false
}

func (p *SmugMugContactService) CanDeleteContact(selfContact bool) bool {
    return false
}

func (p *SmugMugContactService) CanRetrieveGroup(selfContact bool) bool {
    return false
}

func (p *SmugMugContactService) CanCreateGroup(selfContact bool) bool {
    return false
}

func (p *SmugMugContactService) CanUpdateGroup(selfContact bool) bool {
    return false
}

func (p *SmugMugContactService) CanDeleteGroup(selfContact bool) bool {
    return false
}

func (p *SmugMugContactService) GroupListIncludesContactIds() bool {
    return false
}

func (p *SmugMugContactService) GroupInfoIncludesContactIds() bool {
    return false
}

func (p *SmugMugContactService) ContactInfoIncludesGroups() bool {
    return false
}

func (p *SmugMugContactService) RetrieveAllContacts(client oauth2_client.OAuth2Client, ds DataStoreService, dsocialUserId string) ([]*Contact, error) {
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
    for i, e := 0, l.Front(); e != nil; i, e = i+1, e.Next() {
        contacts[i] = e.Value.(*Contact)
    }
    return contacts, err
}

func (p *SmugMugContactService) RetrieveAllConnections(client oauth2_client.OAuth2Client, ds DataStoreService, dsocialUserId string) ([]*Contact, error) {
    return make([]*Contact, 0), nil
}

func (p *SmugMugContactService) RetrieveAllGroups(client oauth2_client.OAuth2Client, ds DataStoreService, dsocialUserId string) ([]*Group, error) {
    return make([]*Group, 0), nil
}

func (p *SmugMugContactService) listToContacts(l *list.List) []*Contact {
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

func (p *SmugMugContactService) RetrieveContacts(client oauth2_client.OAuth2Client, ds DataStoreService, dsocialUserId string, next NextToken) ([]*Contact, NextToken, error) {
    l := list.New()
    var useErr error
    famResp, err := smugmug.RetrieveFamily(client, nil)
    if err != nil {
        useErr = err
    }
    if famResp != nil && famResp.Family != nil {
        for _, u := range famResp.Family {
            contact, err := p.handleRetrievedContact(client, ds, dsocialUserId, u.NickName, &u)
            if contact != nil {
                l.PushBack(contact)
            }
            if err != nil && useErr == nil {
                useErr = err
            }
        }
    }
    if useErr != nil {
        return p.listToContacts(l), nil, useErr
    }
    fansResp, err := smugmug.RetrieveFans(client, nil)
    if err != nil {
        useErr = err
    }
    if fansResp != nil && fansResp.Fans != nil {
        for _, u := range fansResp.Fans {
            contact, err := p.handleRetrievedContact(client, ds, dsocialUserId, u.NickName, &u)
            if contact != nil {
                l.PushBack(contact)
            }
            if err != nil && useErr == nil {
                useErr = err
            }
        }
    }
    if useErr != nil {
        return p.listToContacts(l), nil, useErr
    }
    userInfo, err := client.RetrieveUserInfo()
    if err != nil {
        return p.listToContacts(l), nil, err
    }
    userResp, err := smugmug.RetrieveUserInfo(client, userInfo.Username(), nil)
    if err != nil {
        useErr = err
    }
    if userResp != nil {
        contact, err := p.handleRetrievedContact(client, ds, dsocialUserId, userResp.User.NickName, &userResp.User)
        if contact != nil {
            l.PushBack(contact)
        }
        if err != nil && useErr == nil {
            useErr = err
        }
    }
    return p.listToContacts(l), nil, useErr
}

func (p *SmugMugContactService) RetrieveConnections(client oauth2_client.OAuth2Client, ds DataStoreService, dsocialUserId string, next NextToken) ([]*Contact, NextToken, error) {
    return make([]*Contact, 0), nil, nil
}

func (p *SmugMugContactService) RetrieveGroups(client oauth2_client.OAuth2Client, ds DataStoreService, dsocialUserId string, next NextToken) ([]*Group, NextToken, error) {
    return make([]*Group, 0), nil, nil
}

func (p *SmugMugContactService) RetrieveContact(client oauth2_client.OAuth2Client, ds DataStoreService, dsocialUserId string, contactId string) (*Contact, error) {
    extContact, err := smugmug.RetrieveUserInfo(client, contactId, nil)
    if extContact == nil || err != nil {
        return nil, err
    }
    return p.handleRetrievedContact(client, ds, dsocialUserId, contactId, &extContact.User)
}

func (p *SmugMugContactService) handleRetrievedContact(client oauth2_client.OAuth2Client, ds DataStoreService, dsocialUserId string, contactId string, extContact *smugmug.PersonReference) (contact *Contact, err error) {
    if extContact == nil {
        return nil, nil
    }
    externalServiceId := p.ServiceId()
    userInfo, err := client.RetrieveUserInfo()
    externalUserId := userInfo.Guid()
    var useErr error = nil
    dsocialContactId := ""
    var origDsocialContact *dm.Contact = nil
    externalContactId := extContact.NickName
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
    dsocialContact := dm.SmugMugUserToDsocial(extContact, origDsocialContact, dsocialUserId)
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

func (p *SmugMugContactService) RetrieveGroup(client oauth2_client.OAuth2Client, ds DataStoreService, dsocialUserId string, groupId string) (*Group, error) {
    return nil, nil
}

func (p *SmugMugContactService) CreateContactOnExternalService(client oauth2_client.OAuth2Client, contact interface{}) (interface{}, string, error) {
    return nil, "", nil
}

func (p *SmugMugContactService) CreateGroupOnExternalService(client oauth2_client.OAuth2Client, group interface{}) (interface{}, string, error) {
    return nil, "", nil
}

func (p *SmugMugContactService) UpdateContactOnExternalService(client oauth2_client.OAuth2Client, originalContact, contact interface{}) (interface{}, string, error) {
    return nil, "", nil
}

func (p *SmugMugContactService) UpdateGroupOnExternalService(client oauth2_client.OAuth2Client, originalGroup, group interface{}) (interface{}, string, error) {
    return nil, "", nil
}

func (p *SmugMugContactService) DeleteContactOnExternalService(client oauth2_client.OAuth2Client, contact interface{}) (bool, error) {
    return false, nil
}

func (p *SmugMugContactService) DeleteGroupOnExternalService(client oauth2_client.OAuth2Client, group interface{}) (bool, error) {
    return false, nil
}

func (p *SmugMugContactService) ContactsService() ContactsService {
    return p
}
