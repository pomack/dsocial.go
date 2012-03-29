package contacts

import (
    "github.com/pomack/contacts.go/google"
    dm "github.com/pomack/dsocial.go/models/dsocial"
    "github.com/pomack/jsonhelper.go/jsonhelper"
    "github.com/pomack/oauth2_client.go/oauth2_client"
    //"github.com/pomack/jsonhelper.go/jsonhelper"

    "net/url"
    "strconv"
    "strings"
)

type GoogleContactService struct {
}

type GoogleContactServiceSettings struct {
    StandardContactServiceSettings `json:"settings,omitempty,collapse"`
}

func NewGoogleContactServiceSettings() *GoogleContactServiceSettings {
    s := new(GoogleContactServiceSettings)
    s.SetAllowRetrieveContactInfo(true)
    s.SetAllowModifyContactInfo(true)
    return s
}

func (p *GoogleContactServiceSettings) ContactsServiceId() string {
    return GOOGLE_CONTACT_SERVICE_ID
}

func NewGoogleContactService() *GoogleContactService {
    return new(GoogleContactService)
}

func (p *GoogleContactService) ServiceId() string {
    return GOOGLE_CONTACT_SERVICE_ID
}

func (p *GoogleContactService) CreateOAuth2Client(settings jsonhelper.JSONObject) (client oauth2_client.OAuth2Client, err error) {
    client = oauth2_client.NewGoogleClient()
    client.Initialize(settings)
    return
}

func (p *GoogleContactService) ConvertToDsocialContact(externalContact interface{}, originalDsocialContact *dm.Contact, dsocialUserId string) (dsocialContact *dm.Contact) {
    if externalContact == nil {
        return
    }
    if extContact, ok := externalContact.(*google.Contact); ok && extContact != nil {
        dsocialContact = dm.GoogleContactToDsocial(extContact, originalDsocialContact, dsocialUserId)
    }
    return
}

func (p *GoogleContactService) ConvertToExternalContact(dsocialContact *dm.Contact, originalExternalContact interface{}, dsocialUserId string) (externalContact interface{}) {
    var origGoogleContact *google.Contact = nil
    if originalExternalContact != nil {
        origGoogleContact, _ = originalExternalContact.(*google.Contact)
    }
    externalContact = dm.DsocialContactToGoogle(dsocialContact, origGoogleContact)
    return
}

func (p *GoogleContactService) ConvertToDsocialGroup(externalGroup interface{}, originalDsocialGroup *dm.Group, dsocialUserId string) (dsocialGroup *dm.Group) {
    if externalGroup == nil {
        return
    }
    if extGroup, ok := externalGroup.(*google.ContactGroup); ok && extGroup != nil {
        dsocialGroup = dm.GoogleGroupToDsocial(extGroup, originalDsocialGroup, dsocialUserId)
    }
    return
}

func (p *GoogleContactService) ConvertToExternalGroup(dsocialGroup *dm.Group, originalExternalGroup interface{}, dsocialUserId string) (externalGroup interface{}) {
    var origGoogleGroup *google.ContactGroup = nil
    if originalExternalGroup != nil {
        origGoogleGroup, _ = originalExternalGroup.(*google.ContactGroup)
    }
    externalGroup = dm.DsocialGroupToGoogle(dsocialGroup, origGoogleGroup)
    return
}

func (p *GoogleContactService) CanRetrieveAllContacts() bool {
    return false
}

func (p *GoogleContactService) CanRetrieveAllConnections() bool {
    return false
}

func (p *GoogleContactService) CanRetrieveAllGroups() bool {
    return false
}

func (p *GoogleContactService) CanRetrieveContacts() bool {
    return true
}

func (p *GoogleContactService) CanRetrieveConnections() bool {
    return false
}

func (p *GoogleContactService) CanRetrieveGroups() bool {
    return true
}

func (p *GoogleContactService) CanImportContactsOrGroups() bool {
    return true
}

func (p *GoogleContactService) CanExportContactsOrGroups() bool {
    return true
}

func (p *GoogleContactService) CanRetrieveContact(selfContact bool) bool {
    return true
}

func (p *GoogleContactService) CanCreateContact(selfContact bool) bool {
    return true
}

func (p *GoogleContactService) CanUpdateContact(selfContact bool) bool {
    return true
}

func (p *GoogleContactService) CanDeleteContact(selfContact bool) bool {
    return true
}

func (p *GoogleContactService) CanRetrieveGroup(selfContact bool) bool {
    return true
}

func (p *GoogleContactService) CanCreateGroup(selfContact bool) bool {
    return true
}

func (p *GoogleContactService) CanUpdateGroup(selfContact bool) bool {
    return true
}

func (p *GoogleContactService) CanDeleteGroup(selfContact bool) bool {
    return true
}

func (p *GoogleContactService) GroupListIncludesContactIds() bool {
    return false
}

func (p *GoogleContactService) GroupInfoIncludesContactIds() bool {
    return false
}

func (p *GoogleContactService) ContactInfoIncludesGroups() bool {
    return true
}

func (p *GoogleContactService) RetrieveAllContacts(client oauth2_client.OAuth2Client, ds DataStoreService, dsocialUserId string) ([]*Contact, error) {
    contacts, _, err := p.RetrieveContacts(client, ds, dsocialUserId, nil)
    return contacts, err
}

func (p *GoogleContactService) RetrieveAllConnections(client oauth2_client.OAuth2Client, ds DataStoreService, dsocialUserId string) ([]*Contact, error) {
    return make([]*Contact, 0), nil
}

func (p *GoogleContactService) RetrieveAllGroups(client oauth2_client.OAuth2Client, ds DataStoreService, dsocialUserId string) ([]*Group, error) {
    groups, _, err := p.RetrieveGroups(client, ds, dsocialUserId, nil)
    return groups, err
}

func (p *GoogleContactService) RetrieveContacts(client oauth2_client.OAuth2Client, ds DataStoreService, dsocialUserId string, next NextToken) ([]*Contact, NextToken, error) {
    var m url.Values
    if next == nil {
    } else if s, ok := next.(string); ok {
        if s != "" {
            if strings.HasPrefix(s, "https://www.google.com/") || strings.HasPrefix(s, "http://www.google.com/") {
                uri, err := url.Parse(s)
                if err == nil {
                    q, err := url.ParseQuery(uri.RawQuery)
                    if err == nil {
                        m = q
                    }
                }
            }
            if m == nil {
                m = make(url.Values)
                m.Add("q", s)
            }
        }
    } else if maxResults, ok := next.(int); ok {
        m = make(url.Values)
        m.Add("max-results", strconv.Itoa(maxResults))
    } else if maxResults, ok := next.(int64); ok {
        m = make(url.Values)
        m.Add("max-results", strconv.FormatInt(maxResults, 10))
    } else if cq, ok := next.(*google.ContactQuery); ok {
        m = make(url.Values)
        if cq.Alt != "" {
            m.Add("alt", cq.Alt)
        }
        if cq.Q != "" {
            m.Add("q", cq.Q)
        }
        if cq.MaxResults > 0 {
            m.Add("max-results", strconv.FormatInt(cq.MaxResults, 10))
        }
        if cq.StartIndex > 0 {
            m.Add("start-index", strconv.FormatInt(cq.StartIndex, 10))
        }
        if cq.UpdatedMin != "" {
            m.Add("updated-min", cq.UpdatedMin)
        }
        if cq.OrderBy != "" {
            m.Add("orderby", cq.OrderBy)
        }
        if cq.ShowDeleted {
            m.Add("showdeleted", "true")
        }
        if cq.RequireAllDeleted {
            m.Add("requirealldeleted", "true")
        }
        if cq.SortOrder != "" {
            m.Add("sortorder", cq.SortOrder)
        }
        if cq.Group != "" {
            m.Add("group", cq.Group)
        }
    }
    feed, err := google.RetrieveContacts(client, m)
    var theNextToken NextToken = nil
    if feed != nil && feed.Links != nil && len(feed.Links) > 0 {
        for _, link := range feed.Links {
            if link.Rel == "next" {
                theNextToken = link.Href
                break
            }
        }
    }
    if feed == nil || feed.Entries == nil || len(feed.Entries) == 0 || err != nil {
        return make([]*Contact, 0), theNextToken, err
    }
    contacts := make([]*Contact, len(feed.Entries))
    externalServiceId := p.ServiceId()
    userInfo, err := client.RetrieveUserInfo()
    externalUserId := userInfo.Guid()
    var useErr error = nil
    for i, googleContact := range feed.Entries {
        externalContactId := googleContact.ContactId()
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
                ds.StoreExternalContact(externalServiceId, externalUserId, dsocialUserId, externalContactId, &googleContact)
            }
        }
        dsocialContact := dm.GoogleContactToDsocial(&googleContact, origDsocialContact, dsocialUserId)
        contacts[i] = &Contact{
            ExternalServiceId: p.ServiceId(),
            ExternalUserId:    googleContact.ContactUserId(),
            ExternalContactId: externalContactId,
            DsocialUserId:     dsocialUserId,
            DsocialContactId:  dsocialContactId,
            Value:             dsocialContact,
        }
    }
    return contacts, theNextToken, useErr
}

func (p *GoogleContactService) RetrieveConnections(client oauth2_client.OAuth2Client, ds DataStoreService, dsocialUserId string, next NextToken) ([]*Contact, NextToken, error) {
    return make([]*Contact, 0), nil, nil
}

func (p *GoogleContactService) RetrieveGroups(client oauth2_client.OAuth2Client, ds DataStoreService, dsocialUserId string, next NextToken) ([]*Group, NextToken, error) {
    var m url.Values
    if next == nil {
    } else if s, ok := next.(string); ok {
        if s != "" {
            if strings.HasPrefix(s, "https://www.google.com/") {
                uri, err := url.Parse(s)
                if err == nil {
                    q, err := url.ParseQuery(uri.RawQuery)
                    if err == nil {
                        m = q
                    }
                }
            }
            if m == nil {
                m = make(url.Values)
                m.Add("q", s)
            }
        }
    } else if maxResults, ok := next.(int); ok {
        m = make(url.Values)
        m.Add("max-results", strconv.Itoa(maxResults))
    } else if maxResults, ok := next.(int64); ok {
        m = make(url.Values)
        m.Add("max-results", strconv.FormatInt(maxResults, 10))
    } else if gq, ok := next.(*google.GroupQuery); ok {
        m = make(url.Values)
        if gq.Alt != "" {
            m.Add("alt", gq.Alt)
        }
        if gq.Q != "" {
            m.Add("q", gq.Q)
        }
        if gq.MaxResults > 0 {
            m.Add("max-results", strconv.FormatInt(gq.MaxResults, 10))
        }
        if gq.StartIndex > 0 {
            m.Add("start-index", strconv.FormatInt(gq.StartIndex, 10))
        }
        if gq.UpdatedMin != "" {
            m.Add("updated-min", gq.UpdatedMin)
        }
        if gq.OrderBy != "" {
            m.Add("orderby", gq.OrderBy)
        }
        if gq.ShowDeleted {
            m.Add("showdeleted", "true")
        }
        if gq.RequireAllDeleted {
            m.Add("requirealldeleted", "true")
        }
        if gq.SortOrder != "" {
            m.Add("sortorder", gq.SortOrder)
        }
    }
    resp, err := google.RetrieveGroups(client, m)
    var theNextToken NextToken = nil
    if resp != nil && resp.Feed != nil && resp.Feed.Links != nil && len(resp.Feed.Links) > 0 {
        for _, link := range resp.Feed.Links {
            if link.Rel == "next" {
                theNextToken = link.Href
            }
        }
    }
    if resp == nil || resp.Feed == nil || resp.Feed.Entries == nil || len(resp.Feed.Entries) == 0 || err != nil {
        return make([]*Group, 0), theNextToken, err
    }
    groups := make([]*Group, len(resp.Feed.Entries))
    externalServiceId := p.ServiceId()
    userInfo, err := client.RetrieveUserInfo()
    externalUserId := userInfo.Guid()
    var useErr error = nil
    for i, googleGroup := range resp.Feed.Entries {
        externalGroupId := googleGroup.GroupId()
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
                ds.StoreExternalGroup(externalServiceId, externalUserId, dsocialUserId, externalGroupId, &googleGroup)
            }
        }
        var dsocialGroup *dm.Group = dm.GoogleGroupToDsocial(&googleGroup, origDsocialGroup, dsocialUserId)
        groups[i] = &Group{
            ExternalServiceId: p.ServiceId(),
            ExternalUserId:    googleGroup.GroupUserId(),
            ExternalGroupId:   googleGroup.GroupId(),
            DsocialUserId:     dsocialUserId,
            DsocialGroupId:    dsocialGroupId,
            Value:             dsocialGroup,
        }
    }
    return groups, theNextToken, useErr
}

func (p *GoogleContactService) RetrieveContact(client oauth2_client.OAuth2Client, ds DataStoreService, dsocialUserId string, contactId string) (*Contact, error) {
    googleContact, err := google.RetrieveContact(client, contactId, nil)
    if googleContact == nil || err != nil {
        return nil, err
    }
    externalServiceId := p.ServiceId()
    userInfo, err := client.RetrieveUserInfo()
    externalUserId := userInfo.Guid()
    useErr := err
    dsocialContactId := ""
    var origDsocialContact *dm.Contact = nil
    externalContactId := googleContact.ContactId()
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
            ds.StoreExternalContact(externalServiceId, externalUserId, dsocialUserId, externalContactId, googleContact)
        }
    }
    dsocialContact := dm.GoogleContactToDsocial(googleContact, origDsocialContact, dsocialUserId)
    contact := &Contact{
        ExternalServiceId: p.ServiceId(),
        ExternalUserId:    googleContact.ContactUserId(),
        ExternalContactId: googleContact.ContactId(),
        DsocialUserId:     dsocialUserId,
        DsocialContactId:  dsocialContactId,
        Value:             dsocialContact,
    }
    return contact, useErr
}

func (p *GoogleContactService) RetrieveGroup(client oauth2_client.OAuth2Client, ds DataStoreService, dsocialUserId string, groupId string) (*Group, error) {
    resp, err := google.RetrieveGroup(client, groupId, nil)
    if resp == nil || resp.Entry == nil || err != nil {
        return nil, err
    }
    googleGroup := resp.Entry
    externalServiceId := p.ServiceId()
    userInfo, err := client.RetrieveUserInfo()
    externalUserId := userInfo.Guid()
    useErr := err
    externalGroupId := googleGroup.GroupId()
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
            ds.StoreExternalGroup(externalServiceId, externalUserId, dsocialUserId, externalGroupId, googleGroup)
        }
    }
    var dsocialGroup *dm.Group = dm.GoogleGroupToDsocial(googleGroup, origDsocialGroup, dsocialUserId)
    group := &Group{
        ExternalServiceId: p.ServiceId(),
        ExternalUserId:    externalUserId,
        ExternalGroupId:   externalGroupId,
        DsocialUserId:     dsocialUserId,
        DsocialGroupId:    dsocialGroupId,
        Value:             dsocialGroup,
    }
    return group, useErr
}

func (p *GoogleContactService) CreateContactOnExternalService(client oauth2_client.OAuth2Client, externalContact interface{}) (externalContactResult interface{}, externalContactId string, err error) {
    if externalContact == nil {
        return nil, "", nil
    }
    gContact := externalContact.(*google.Contact)
    resp, err := google.CreateContact(client, "", "", gContact)
    if resp == nil || resp.Entry == nil {
        return nil, "", err
    }
    return resp.Entry, resp.Entry.ContactId(), err
}

func (p *GoogleContactService) CreateGroupOnExternalService(client oauth2_client.OAuth2Client, externalGroup interface{}) (externalGroupResult interface{}, externalGroupId string, err error) {
    if externalGroup == nil {
        return nil, "", nil
    }
    gGroup := externalGroup.(*google.ContactGroup)
    resp, err := google.CreateGroup(client, "", "", gGroup)
    if resp == nil || resp.Entry == nil {
        return nil, "", err
    }
    return resp.Entry, resp.Entry.GroupId(), err
}

func (p *GoogleContactService) UpdateContactOnExternalService(client oauth2_client.OAuth2Client, originalContact, latestContact interface{}) (externalContactResult interface{}, externalContactId string, err error) {
    if originalContact == nil || latestContact == nil {
        return nil, "", nil
    }
    originalGContact := originalContact.(*google.Contact)
    latestGContact := latestContact.(*google.Contact)
    latestGContact.SetContactId(originalGContact.ContactId())
    resp, err := google.UpdateContact(client, "", "", originalGContact, latestGContact)
    if resp == nil || resp.Entry == nil {
        return nil, "", err
    }
    return resp.Entry, resp.Entry.ContactId(), err
}

func (p *GoogleContactService) UpdateGroupOnExternalService(client oauth2_client.OAuth2Client, originalGroup, latestGroup interface{}) (externalGroupResult interface{}, externalGroupId string, err error) {
    if originalGroup == nil || latestGroup == nil {
        return nil, "", nil
    }
    originalGGroup := originalGroup.(*google.ContactGroup)
    latestGGroup := latestGroup.(*google.ContactGroup)
    latestGGroup.SetGroupId(originalGGroup.GroupId())
    resp, err := google.UpdateGroup(client, "", "", originalGGroup, latestGGroup)
    if resp == nil || resp.Entry == nil {
        return nil, "", err
    }
    return resp.Entry, resp.Entry.GroupId(), err
}

func (p *GoogleContactService) DeleteContactOnExternalService(client oauth2_client.OAuth2Client, originalContact interface{}) (bool, error) {
    if originalContact == nil {
        return false, nil
    }
    originalGContact := originalContact.(*google.Contact)
    err := google.DeleteContact(client, "", "", originalGContact)
    return true, err
}

func (p *GoogleContactService) DeleteGroupOnExternalService(client oauth2_client.OAuth2Client, originalGroup interface{}) (bool, error) {
    if originalGroup == nil {
        return false, nil
    }
    originalGGroup := originalGroup.(*google.ContactGroup)
    err := google.DeleteGroup(client, "", "", originalGGroup)
    return true, err
}

func (p *GoogleContactService) ContactsService() ContactsService {
    return p
}
