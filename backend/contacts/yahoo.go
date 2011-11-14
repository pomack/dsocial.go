package contacts

import (
    "github.com/pomack/jsonhelper.go/jsonhelper"
    "github.com/pomack/oauth2_client.go/oauth2_client"
    "github.com/pomack/contacts.go/yahoo"
    dm "github.com/pomack/dsocial.go/models/dsocial"
    //"github.com/pomack/jsonhelper.go/jsonhelper"
    "os"
    "strconv"
    "url"
)

type YahooContactService struct {

}

func NewYahooContactService() *YahooContactService {
    return new(YahooContactService)
}

func (p *YahooContactService) ServiceId() string {
    return "www.yahoo.com"
}

func (p *YahooContactService) CreateOAuth2Client(settings jsonhelper.JSONObject) (client oauth2_client.OAuth2Client, err os.Error) {
    client = oauth2_client.NewYahooClient()
    client.Initialize(settings)
    return
}

func (p *YahooContactService) ConvertToDsocialContact(externalContact interface{}, originalDsocialContact *dm.Contact, dsocialUserId string) (dsocialContact *dm.Contact) {
    if externalContact == nil {
        return
    }
    if extContact, ok := externalContact.(*yahoo.Contact); ok && extContact != nil {
        dsocialContact = dm.YahooContactToDsocial(extContact, originalDsocialContact, dsocialUserId)
    }
    return
}

func (p *YahooContactService) ConvertToExternalContact(dsocialContact *dm.Contact, originalExternalContact interface{}, dsocialUserId string) (externalContact interface{}) {
    var origYahooContact *yahoo.Contact = nil
    if originalExternalContact != nil {
        origYahooContact, _ = originalExternalContact.(*yahoo.Contact)
    }
    externalContact = dm.DsocialContactToYahoo(dsocialContact, origYahooContact)
    return
}

func (p *YahooContactService) ConvertToDsocialGroup(externalGroup interface{}, originalDsocialGroup *dm.Group, dsocialUserId string) (dsocialGroup *dm.Group) {
    if externalGroup == nil {
        return
    }
    if extGroup, ok := externalGroup.(*yahoo.Category); ok && extGroup != nil {
        dsocialGroup = dm.YahooCategoryToDsocial(extGroup, originalDsocialGroup, dsocialUserId)
    }
    return
}

func (p *YahooContactService) ConvertToExternalGroup(dsocialGroup *dm.Group, originalExternalGroup interface{}, dsocialUserId string) (externalGroup interface{}) {
    var origYahooGroup *yahoo.Category = nil
    if originalExternalGroup != nil {
        origYahooGroup, _ = originalExternalGroup.(*yahoo.Category)
    }
    externalGroup = dm.DsocialGroupToYahoo(dsocialGroup, origYahooGroup)
    return
}


func (p *YahooContactService) CanRetrieveAllContacts() bool {
    return false
}

func (p *YahooContactService) CanRetrieveAllConnections() bool {
    return false
}

func (p *YahooContactService) CanRetrieveAllGroups() bool {
    return false
}

func (p *YahooContactService) CanRetrieveContacts() bool {
    return true
}

func (p *YahooContactService) CanRetrieveConnections() bool {
    return false
}

func (p *YahooContactService) CanRetrieveGroups() bool {
    return true
}

func (p *YahooContactService) CanRetrieveContact(selfContact bool) bool {
    return true
}

func (p *YahooContactService) CanImportContactsOrGroups() bool {
    return true
}

func (p *YahooContactService) CanExportContactsOrGroups() bool {
    return true
}

func (p *YahooContactService) CanCreateContact(selfContact bool) bool {
    return true
}

func (p *YahooContactService) CanUpdateContact(selfContact bool) bool {
    return true
}

func (p *YahooContactService) CanDeleteContact(selfContact bool) bool {
    return true
}

func (p *YahooContactService) CanRetrieveGroup(selfContact bool) bool {
    return true
}

func (p *YahooContactService) CanCreateGroup(selfContact bool) bool {
    return false
}

func (p *YahooContactService) CanUpdateGroup(selfContact bool) bool {
    return false
}

func (p *YahooContactService) CanDeleteGroup(selfContact bool) bool {
    return false
}

func (p *YahooContactService) GroupListIncludesContactIds() bool {
    return false
}

func (p *YahooContactService) GroupInfoIncludesContactIds() bool {
    return false
}

func (p *YahooContactService) ContactInfoIncludesGroups() bool {
    return true
}

func (p *YahooContactService) RetrieveAllContacts(client oauth2_client.OAuth2Client, ds DataStoreService, dsocialUserId string) ([]*Contact, os.Error) {
    contacts, _, err := p.RetrieveContacts(client, ds, dsocialUserId, nil)
    return contacts, err
}

func (p *YahooContactService) RetrieveAllConnections(client oauth2_client.OAuth2Client, ds DataStoreService, dsocialUserId string) ([]*Contact, os.Error) {
    return make([]*Contact, 0), nil
}

func (p *YahooContactService) RetrieveAllGroups(client oauth2_client.OAuth2Client, ds DataStoreService, dsocialUserId string) ([]*Group, os.Error) {
    groups, _, err := p.RetrieveGroups(client, ds, dsocialUserId, nil)
    return groups, err
}

func (p *YahooContactService) RetrieveContacts(client oauth2_client.OAuth2Client, ds DataStoreService, dsocialUserId string, next NextToken) ([]*Contact, NextToken, os.Error) {
    var m url.Values
    m = make(url.Values)
    m.Add("count", "max")
    if next == nil {
    } else if start, ok := next.(int); ok {
        m.Add("start", strconv.Itoa(start))
    }
    resp, err := yahoo.RetrieveContacts(client, m)
    if resp == nil || resp.Contacts.Contacts == nil || len(resp.Contacts.Contacts) == 0 || err != nil {
        return make([]*Contact, 0), nil, err
    }
    contacts := make([]*Contact, len(resp.Contacts.Contacts))
    externalServiceId := p.ServiceId()
    userInfo, err := client.RetrieveUserInfo()
    externalUserId := userInfo.Guid()
    var useErr os.Error = nil
    for i, yahooContact := range resp.Contacts.Contacts {
        externalContactId := strconv.Itoa64(yahooContact.Id)
        dsocialContactId := ""
        var origDsocialContact *dm.Contact = nil
        if len(externalContactId) > 0 {
            dsocialContactId, err = ds.DsocialIdForExternalContactId(externalServiceId, externalUserId, dsocialUserId, externalContactId)
            if err != nil {
                useErr = err
                continue
            }
            if dsocialContactId != "" {
                origDsocialContact, _, err = ds.RetrieveDsocialContactForExternalContact(externalServiceId, externalUserId, externalContactId, dsocialUserId)
                if err != nil {
                    useErr = err
                    continue
                }
            } else {
                ds.StoreExternalContact(externalServiceId, externalUserId, dsocialUserId, externalContactId, &yahooContact)
            }
        }
        dsocialContact := dm.YahooContactToDsocial(&yahooContact, origDsocialContact, dsocialUserId)
        contacts[i] = &Contact{
            ExternalServiceId: p.ServiceId(),
            ExternalUserId: externalUserId,
            ExternalContactId: externalContactId,
            DsocialUserId: dsocialUserId,
            DsocialContactId: dsocialContactId,
            Value: dsocialContact,
        }
    }
    return contacts, nil, useErr
}

func (p *YahooContactService) RetrieveConnections(client oauth2_client.OAuth2Client, ds DataStoreService, dsocialUserId string, next NextToken) ([]*Contact, NextToken, os.Error) {
    return make([]*Contact, 0), nil, nil
}

func (p *YahooContactService) RetrieveGroups(client oauth2_client.OAuth2Client, ds DataStoreService, dsocialUserId string, next NextToken) ([]*Group, NextToken, os.Error) {
    var m url.Values
    m = make(url.Values)
    m.Add("count", "max")
    if next == nil {
    } else if start, ok := next.(int); ok {
        m.Add("start", strconv.Itoa(start))
    }
    resp, err := yahoo.RetrieveCategories(client, m)
    if resp == nil || resp.Categories.Categories == nil || len(resp.Categories.Categories) == 0 || err != nil {
        return make([]*Group, 0), nil, err
    }
    groups := make([]*Group, len(resp.Categories.Categories))
    externalServiceId := p.ServiceId()
    userInfo, err := client.RetrieveUserInfo()
    externalUserId := userInfo.Guid()
    var useErr os.Error = nil
    for i, yahooGroup := range resp.Categories.Categories {
        var externalGroupId string
        if yahooGroup.Id > 0 {
            externalGroupId = strconv.Itoa64(yahooGroup.Id)
        }
        var origDsocialGroup *dm.Group = nil
        dsocialGroupId := ""
        if len(externalGroupId) > 0 {
            dsocialGroupId, err = ds.DsocialIdForExternalGroupId(externalServiceId, externalUserId, dsocialUserId, externalGroupId)
            if err != nil {
                if useErr == nil {
                    useErr = err
                }
                continue
            }
            if dsocialGroupId != "" {
                origDsocialGroup, _, err = ds.RetrieveDsocialGroupForExternalGroup(externalServiceId, externalUserId, externalGroupId, dsocialUserId)
                if err != nil {
                    if useErr == nil {
                        useErr = err
                    }
                    continue
                }
            } else {
                ds.StoreExternalGroup(externalServiceId, externalUserId, dsocialUserId, externalGroupId, &yahooGroup)
            }
        }
        var dsocialGroup *dm.Group = dm.YahooCategoryToDsocial(&yahooGroup, origDsocialGroup, dsocialUserId)
        groups[i] = &Group{
            ExternalServiceId: p.ServiceId(),
            ExternalUserId: externalUserId,
            ExternalGroupId: externalGroupId,
            DsocialUserId: dsocialUserId,
            DsocialGroupId: dsocialGroupId,
            Value: dsocialGroup,
        }
    }
    return groups, nil, useErr
}

func (p *YahooContactService) RetrieveContact(client oauth2_client.OAuth2Client, ds DataStoreService, dsocialUserId string, contactId string) (*Contact, os.Error) {
    resp, err := yahoo.RetrieveContact(client, contactId, nil)
    if resp == nil || err != nil {
        return nil, err
    }
    yahooContact := &resp.Contact
    externalServiceId := p.ServiceId()
    userInfo, err := client.RetrieveUserInfo()
    externalUserId := userInfo.Guid()
    useErr := err
    dsocialContactId := ""
    var origDsocialContact *dm.Contact = nil
    var externalContactId string
    if yahooContact.Id > 0 {
        externalContactId = strconv.Itoa64(yahooContact.Id)
    }
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
            ds.StoreExternalContact(externalServiceId, externalUserId, dsocialUserId, externalContactId, yahooContact)
        }
    }
    dsocialContact := dm.YahooContactToDsocial(yahooContact, origDsocialContact, dsocialUserId)
    contact := &Contact{
        ExternalServiceId: p.ServiceId(),
        ExternalUserId: externalUserId,
        ExternalContactId: externalContactId,
        DsocialUserId: dsocialUserId,
        DsocialContactId: dsocialContactId,
        Value: dsocialContact,
    }
    return contact, useErr
}

func (p *YahooContactService) RetrieveGroup(client oauth2_client.OAuth2Client, ds DataStoreService, dsocialUserId string, groupId string) (*Group, os.Error) {
    resp, err := yahoo.RetrieveCategory(client, groupId, nil)
    if resp == nil || err != nil {
        return nil, err
    }
    yahooGroup := &resp.Category
    externalServiceId := p.ServiceId()
    userInfo, err := client.RetrieveUserInfo()
    externalUserId := userInfo.Guid()
    useErr := err
    var externalGroupId string
    if yahooGroup.Id > 0 {
        externalGroupId = strconv.Itoa64(yahooGroup.Id)
    }
    dsocialGroupId := ""
    var origDsocialGroup *dm.Group = nil
    if len(externalGroupId) > 0 {
        dsocialGroupId, err = ds.DsocialIdForExternalGroupId(externalServiceId, externalUserId, dsocialUserId, externalGroupId)
        if err != nil {
            if useErr == nil {
                useErr = err
            }
        }
        if dsocialGroupId != "" {
            origDsocialGroup, _, err = ds.RetrieveDsocialGroupForExternalGroup(externalServiceId, externalUserId, externalGroupId, dsocialUserId)
            if err != nil && useErr == nil {
                useErr = err
            }
        } else {
            ds.StoreExternalGroup(externalServiceId, externalUserId, dsocialUserId, externalGroupId, yahooGroup)
        }
    }
    var dsocialGroup *dm.Group = dm.YahooCategoryToDsocial(yahooGroup, origDsocialGroup, dsocialUserId)
    group := &Group{
        ExternalServiceId: p.ServiceId(),
        ExternalUserId: externalUserId,
        ExternalGroupId: externalGroupId,
        DsocialUserId: dsocialUserId,
        DsocialGroupId: dsocialGroupId,
        Value: dsocialGroup,
    }
    return group, useErr
}

func (p *YahooContactService) CreateContactOnExternalService(client oauth2_client.OAuth2Client, contact interface{}) (interface{}, string, os.Error) {
    if contact == nil {
        return nil, "", nil
    }
    yContact := contact.(*yahoo.Contact)
    err := yahoo.CreateContact(client, "", yContact)
    return yContact, strconv.Itoa64(yContact.Id), err
}

func (p *YahooContactService) CreateGroupOnExternalService(client oauth2_client.OAuth2Client, group interface{}) (interface{}, string, os.Error) {
    return nil, "", nil
}

func (p *YahooContactService) UpdateContactOnExternalService(client oauth2_client.OAuth2Client, originalContact, contact interface{}) (interface{}, string, os.Error) {
    if contact == nil {
        return nil, "", nil
    }
    if originalContact == nil {
        return p.CreateContactOnExternalService(client, contact)
    }
    originalYContact := originalContact.(*yahoo.Contact)
    yContact := contact.(*yahoo.Contact)
    err := yahoo.UpdateContact(client, "", strconv.Itoa64(originalYContact.Id), yContact)
    return yContact, strconv.Itoa64(yContact.Id), err
}

func (p *YahooContactService) UpdateGroupOnExternalService(client oauth2_client.OAuth2Client, originalGroup, group interface{}) (interface{}, string, os.Error) {
    return nil, "", nil
}

func (p *YahooContactService) DeleteContactOnExternalService(client oauth2_client.OAuth2Client, contact interface{}) (bool, os.Error) {
    if contact == nil {
        return false, nil
    }
    yContact := contact.(*yahoo.Contact)
    err := yahoo.DeleteContact(client, "", strconv.Itoa64(yContact.Id))
    return true, err
}

func (p *YahooContactService) DeleteGroupOnExternalService(client oauth2_client.OAuth2Client, group interface{}) (bool, os.Error) {
    return false, nil
}




func (p *YahooContactService) ContactsService() ContactsService {
    return p
}
