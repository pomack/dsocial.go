package contacts

import (
    "github.com/pomack/oauth2_client.go/oauth2_client"
    "github.com/pomack/contacts.go/google"
    dm "github.com/pomack/dsocial.go/models/dsocial"
    //"github.com/pomack/jsonhelper.go/jsonhelper"
    "os"
    "strconv"
    "strings"
    "url"
)

type GoogleContactService struct {

}

func NewGoogleContactService() *GoogleContactService {
    return new(GoogleContactService)
}

func (p *GoogleContactService) ConvertToDsocialContact(externalContact interface{}) (dsocialContact *dm.Contact) {
    if externalContact == nil {
        return
    }
    if extContact, ok := externalContact.(*google.Contact); ok && extContact != nil {
        dsocialContact = dm.GoogleContactToDsocial(extContact)
    }
    return
}

func (p *GoogleContactService) ConvertToExternalContact(dsocialContact *dm.Contact) (externalContact interface{}) {
    externalContact = dm.DsocialContactToGoogle(dsocialContact)
    return
}

func (p *GoogleContactService) ConvertToDsocialGroup(externalGroup interface{}) (dsocialGroup *dm.Group) {
    if externalGroup == nil {
        return
    }
    if extGroup, ok := externalGroup.(*google.ContactGroup); ok && extGroup != nil {
        dsocialGroup = dm.GoogleGroupToDsocial(extGroup)
    }
    return
}

func (p *GoogleContactService) ConvertToExternalGroup(dsocialGroup *dm.Group) (externalGroup interface{}) {
    externalGroup = dm.DsocialGroupToGoogle(dsocialGroup)
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

func (p *GoogleContactService) RetrieveAllContacts(client oauth2_client.OAuth2Client, ds DataStoreService, dsocialUserId string) ([]*Contact, os.Error) {
    contacts, _, err := p.RetrieveContacts(client, ds, dsocialUserId, nil)
    return contacts, err
}

func (p *GoogleContactService) RetrieveAllConnections(client oauth2_client.OAuth2Client, ds DataStoreService, dsocialUserId string) ([]*Contact, os.Error) {
    return make([]*Contact, 0), nil
}

func (p *GoogleContactService) RetrieveAllGroups(client oauth2_client.OAuth2Client, ds DataStoreService, dsocialUserId string) ([]*Group, os.Error) {
    groups, _, err := p.RetrieveGroups(client, ds, dsocialUserId, nil)
    return groups, err
}

func (p *GoogleContactService) RetrieveContacts(client oauth2_client.OAuth2Client, ds DataStoreService, dsocialUserId string, next NextToken) ([]*Contact, NextToken, os.Error) {
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
        m.Add("max-results", strconv.Itoa64(maxResults))
    } else if cq, ok := next.(*google.ContactQuery); ok {
        m = make(url.Values)
        if cq.Alt != "" {
            m.Add("alt", cq.Alt)
        }
        if cq.Q != "" {
            m.Add("q", cq.Q)
        }
        if cq.MaxResults > 0 {
            m.Add("max-results", strconv.Itoa64(cq.MaxResults))
        }
        if cq.StartIndex > 0 {
            m.Add("start-index", strconv.Itoa64(cq.StartIndex))
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
    externalServiceId := client.ServiceId()
    userInfo, err := client.RetrieveUserInfo()
    externalUserId := userInfo.Guid()
    var useErr os.Error = nil
    for i, googleContact := range feed.Entries {
        dsocialContact := dm.GoogleContactToDsocial(&googleContact)
        dsocialContact.UserId = dsocialUserId
        externalContactId := googleContact.ContactId()
        contacts[i] = &Contact{
            ExternalServiceId: client.ServiceId(),
            ExternalUserId: googleContact.ContactUserId(),
            ExternalContactId: externalContactId,
            DsocialUserId: dsocialUserId,
            Value: dsocialContact,
        }
        if len(externalContactId) > 0 {
            dsocialContactId, err := ds.DsocialIdForExternalContactId(externalServiceId, externalUserId, dsocialUserId, externalContactId)
            if err != nil {
                useErr = err
                continue
            }
            if dsocialContactId != "" {
                dsocialContact.Id = dsocialContactId
                contacts[i].DsocialContactId = dsocialContactId
            } else {
                ds.StoreExternalContact(externalServiceId, externalUserId, dsocialUserId, externalContactId, &googleContact)
            }
        }
    }
    return contacts, theNextToken, useErr
}

func (p *GoogleContactService) RetrieveConnections(client oauth2_client.OAuth2Client, ds DataStoreService, dsocialUserId string, next NextToken) ([]*Contact, NextToken, os.Error) {
    return make([]*Contact, 0), nil, nil
}

func (p *GoogleContactService) RetrieveGroups(client oauth2_client.OAuth2Client, ds DataStoreService, dsocialUserId string, next NextToken) ([]*Group, NextToken, os.Error) {
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
        m.Add("max-results", strconv.Itoa64(maxResults))
    } else if gq, ok := next.(*google.GroupQuery); ok {
        m = make(url.Values)
        if gq.Alt != "" {
            m.Add("alt", gq.Alt)
        }
        if gq.Q != "" {
            m.Add("q", gq.Q)
        }
        if gq.MaxResults > 0 {
            m.Add("max-results", strconv.Itoa64(gq.MaxResults))
        }
        if gq.StartIndex > 0 {
            m.Add("start-index", strconv.Itoa64(gq.StartIndex))
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
    externalServiceId := client.ServiceId()
    userInfo, err := client.RetrieveUserInfo()
    externalUserId := userInfo.Guid()
    var useErr os.Error = nil
    for i, googleGroup := range resp.Feed.Entries {
        var dsocialGroup *dm.Group = dm.GoogleGroupToDsocial(&googleGroup)
        dsocialGroup.UserId = dsocialUserId
        externalGroupId := googleGroup.GroupId()
        groups[i] = &Group{
            ExternalServiceId: client.ServiceId(),
            ExternalUserId: googleGroup.GroupUserId(),
            ExternalGroupId: googleGroup.GroupId(),
            DsocialUserId: dsocialUserId,
            Value: dsocialGroup,
        }
        if len(externalGroupId) > 0 {
            dsocialGroupId, err := ds.DsocialIdForExternalGroupId(externalServiceId, externalUserId, dsocialUserId, externalGroupId)
            if err != nil {
                if useErr == nil {
                    useErr = err
                }
                continue
            }
            if dsocialGroupId != "" {
                dsocialGroup.Id = dsocialGroupId
                groups[i].DsocialGroupId = dsocialGroupId
            } else {
                ds.StoreExternalGroup(externalServiceId, externalUserId, dsocialUserId, externalGroupId, &googleGroup)
            }
        }
    }
    return groups, theNextToken, useErr
}

func (p *GoogleContactService) RetrieveContact(client oauth2_client.OAuth2Client, ds DataStoreService, dsocialUserId string, contactId string) (*Contact, os.Error) {
    googleContact, err := google.RetrieveContact(client, contactId, nil)
    if googleContact == nil || err != nil {
        return nil, err
    }
    externalServiceId := client.ServiceId()
    userInfo, err := client.RetrieveUserInfo()
    externalUserId := userInfo.Guid()
    var useErr os.Error = nil
    dsocialContact := dm.GoogleContactToDsocial(googleContact)
    dsocialContact.UserId = dsocialUserId
    externalContactId := googleContact.ContactId()
    contact := &Contact{
        ExternalServiceId: client.ServiceId(),
        ExternalUserId: googleContact.ContactUserId(),
        ExternalContactId: googleContact.ContactId(),
        DsocialUserId: dsocialUserId,
        Value: dsocialContact,
    }
    if len(externalContactId) > 0 {
        dsocialContactId, err := ds.DsocialIdForExternalContactId(externalServiceId, externalUserId, dsocialUserId, contactId)
        if err != nil {
            useErr = err
        }
        if dsocialContactId != "" {
            dsocialContact.Id = dsocialContactId
            contact.DsocialContactId = dsocialContactId
        } else {
            ds.StoreExternalContact(externalServiceId, externalUserId, dsocialUserId, externalContactId, googleContact)
        }
    }
    return contact, useErr
}

func (p *GoogleContactService) RetrieveGroup(client oauth2_client.OAuth2Client, ds DataStoreService, dsocialUserId string, groupId string) (*Group, os.Error) {
    resp, err := google.RetrieveGroup(client, groupId, nil)
    if resp == nil || resp.Entry == nil || err != nil {
        return nil, err
    }
    googleGroup := resp.Entry
    externalServiceId := client.ServiceId()
    userInfo, err := client.RetrieveUserInfo()
    externalUserId := userInfo.Guid()
    var useErr os.Error = nil
    var dsocialGroup *dm.Group = dm.GoogleGroupToDsocial(googleGroup)
    dsocialGroup.UserId = dsocialUserId
    externalGroupId := googleGroup.GroupId()
    group := &Group{
        ExternalServiceId: client.ServiceId(),
        ExternalUserId: externalUserId,
        ExternalGroupId: externalGroupId,
        DsocialUserId: dsocialUserId,
        Value: dsocialGroup,
    }
    if len(externalGroupId) > 0 {
        dsocialGroupId, err := ds.DsocialIdForExternalGroupId(externalServiceId, externalUserId, dsocialUserId, externalGroupId)
        if err != nil {
            if useErr == nil {
                useErr = err
            }
        }
        if dsocialGroupId != "" {
            dsocialGroup.Id = dsocialGroupId
            group.DsocialGroupId = dsocialGroupId
        } else {
            ds.StoreExternalGroup(externalServiceId, externalUserId, dsocialUserId, externalGroupId, googleGroup)
        }
    }
    return group, useErr
}

func (p *GoogleContactService) CreateContact(client oauth2_client.OAuth2Client, ds DataStoreService, dsocialUserId string, contact *dm.Contact) (*Contact, os.Error) {
    if contact == nil {
        return nil, nil
    }
    userInfo, err := client.RetrieveUserInfo()
    if err != nil {
        return nil, err
    }
    gContactId, err := ds.ExternalContactIdForDsocialId(client.ServiceId(), userInfo.Guid(), dsocialUserId, contact.Id)
    if err != nil {
        return nil, err
    }
    if gContactId != "" {
        originalContact, _, err := ds.RetrieveDsocialContact(dsocialUserId, contact.Id)
        if err != nil {
            return nil, err
        }
        return p.UpdateContact(client, ds, dsocialUserId, originalContact, contact)
    }
    gContact := dm.DsocialContactToGoogle(contact)
    resp, err := google.CreateContact(client, "", "", gContact)
    if resp == nil || resp.Entry == nil || err != nil {
        return nil, err
    }
    dsocialContact := dm.GoogleContactToDsocial(resp.Entry)
    if contact.Id != "" {
        dsocialContact.Id = contact.Id
        _, _, err = ds.StoreDsocialExternalContactMapping(client.ServiceId(), userInfo.Guid(), resp.Entry.ContactId(), dsocialUserId, contact.Id)
    }
    outContact := &Contact{
        ExternalServiceId: client.ServiceId(),
        ExternalUserId: resp.Entry.ContactUserId(),
        ExternalContactId: resp.Entry.ContactId(),
        DsocialUserId: dsocialUserId,
        DsocialContactId: dsocialContact.Id,
        Value: dsocialContact,
    }
    return outContact, err
}

func (p *GoogleContactService) CreateGroup(client oauth2_client.OAuth2Client, ds DataStoreService, dsocialUserId string, group *dm.Group) (*Group, os.Error) {
    if group == nil {
        return nil, nil
    }
    userInfo, err := client.RetrieveUserInfo()
    if err != nil {
        return nil, err
    }
    gGroupId, err := ds.ExternalContactIdForDsocialId(client.ServiceId(), userInfo.Guid(), dsocialUserId, group.Id)
    if err != nil {
        return nil, err
    }
    if gGroupId != "" {
        originalGroup, _, err := ds.RetrieveDsocialGroup(dsocialUserId, group.Id)
        if err != nil {
            return nil, err
        }
        return p.UpdateGroup(client, ds, dsocialUserId, originalGroup, group)
    }
    var gGroup *google.ContactGroup = dm.DsocialGroupToGoogle(group)
    resp, err := google.CreateGroup(client, "", "", gGroup)
    if resp == nil || resp.Entry == nil || err != nil {
        return nil, err
    }
    var dsocialGroup *dm.Group = dm.GoogleGroupToDsocial(resp.Entry)
    if group.Id != "" {
        dsocialGroup.Id = group.Id
        _, _, err = ds.StoreDsocialExternalContactMapping(client.ServiceId(), userInfo.Guid(), resp.Entry.GroupId(), dsocialUserId, group.Id)
    }
    outGroup := &Group{
        ExternalServiceId: client.ServiceId(),
        ExternalUserId: resp.Entry.GroupUserId(),
        ExternalGroupId: resp.Entry.GroupId(),
        DsocialUserId: dsocialUserId,
        DsocialGroupId: dsocialGroup.Id,
        Value: dsocialGroup,
    }
    return outGroup, err
}

func (p *GoogleContactService) UpdateContact(client oauth2_client.OAuth2Client, ds DataStoreService, dsocialUserId string, originalContact, contact *dm.Contact) (*Contact, os.Error) {
    if contact == nil || originalContact == nil {
        return nil, nil
    }
    userInfo, err := client.RetrieveUserInfo()
    if err != nil {
        return nil, err
    }
    gContactId, err := ds.ExternalContactIdForDsocialId(client.ServiceId(), userInfo.Guid(), dsocialUserId, originalContact.Id)
    if err != nil {
        return nil, err
    }
    if gContactId == "" {
        return p.CreateContact(client, ds, dsocialUserId, contact)
    }
    originalGContact, _, err := ds.RetrieveExternalContact(client.ServiceId(), userInfo.Guid(), dsocialUserId, gContactId)
    if err != nil {
        return nil, err
    }
    gContact := dm.DsocialContactToGoogle(contact)
    gContact.SetContactId(gContactId)
    resp, err := google.UpdateContact(client, "", "", originalGContact.(*google.Contact), gContact)
    if resp == nil || resp.Entry == nil || err != nil {
        return nil, err
    }
    dsocialContact := dm.GoogleContactToDsocial(resp.Entry)
    dsocialContact.Id = originalContact.Id
    outContact := &Contact{
        ExternalServiceId: client.ServiceId(),
        ExternalUserId: resp.Entry.ContactUserId(),
        ExternalContactId: resp.Entry.ContactId(),
        DsocialUserId: dsocialUserId,
        DsocialContactId: dsocialContact.Id,
        Value: dsocialContact,
    }
    return outContact, err
}

func (p *GoogleContactService) UpdateGroup(client oauth2_client.OAuth2Client, ds DataStoreService, dsocialUserId string, originalGroup, group *dm.Group) (*Group, os.Error) {
    if group == nil {
        return nil, nil
    }
    userInfo, err := client.RetrieveUserInfo()
    if err != nil {
        return nil, err
    }
    gGroupId, err := ds.ExternalContactIdForDsocialId(client.ServiceId(), userInfo.Guid(), dsocialUserId, originalGroup.Id)
    if err != nil {
        return nil, err
    }
    if gGroupId == "" {
        return p.CreateGroup(client, ds, dsocialUserId, group)
    }
    originalGGroup, _, err := ds.RetrieveExternalGroup(client.ServiceId(), userInfo.Guid(), dsocialUserId, gGroupId)
    if err != nil {
        return nil, err
    }
    var gGroup *google.ContactGroup = dm.DsocialGroupToGoogle(group)
    gGroup.SetGroupId(gGroupId)
    resp, err := google.UpdateGroup(client, "", "", originalGGroup.(*google.ContactGroup), gGroup)
    if resp == nil || resp.Entry == nil || err != nil {
        return nil, err
    }
    dsocialGroup := dm.GoogleGroupToDsocial(resp.Entry)
    outGroup := &Group{
        ExternalServiceId: client.ServiceId(),
        ExternalUserId: resp.Entry.GroupUserId(),
        ExternalGroupId: resp.Entry.GroupId(),
        DsocialUserId: dsocialUserId,
        DsocialGroupId: dsocialGroup.Id,
        Value: dsocialGroup,
    }
    return outGroup, err
}

func (p *GoogleContactService) DeleteContact(client oauth2_client.OAuth2Client, ds DataStoreService, dsocialUserId, dsocialContactId string) (bool, os.Error) {
    if dsocialContactId == "" || dsocialUserId == "" {
        return false, nil
    }
    userInfo, err := client.RetrieveUserInfo()
    if err != nil {
        return true, err
    }
    gContactId, err := ds.ExternalContactIdForDsocialId(client.ServiceId(), userInfo.Guid(), dsocialUserId, dsocialContactId)
    if gContactId == "" || err != nil {
        return false, err
    }
    original, _, err := ds.RetrieveExternalContact(client.ServiceId(), userInfo.Guid(), dsocialUserId, gContactId)
    if err != nil {
        return true, err
    }
    err = google.DeleteContact(client, "", "", original.(*google.Contact))
    return true, err
}

func (p *GoogleContactService) DeleteGroup(client oauth2_client.OAuth2Client, ds DataStoreService, dsocialUserId, dsocialGroupId string) (bool, os.Error) {
    if dsocialGroupId == "" || dsocialUserId == "" {
        return false, nil
    }
    userInfo, err := client.RetrieveUserInfo()
    if err != nil {
        return true, err
    }
    gGroupId, err := ds.ExternalContactIdForDsocialId(client.ServiceId(), userInfo.Guid(), dsocialUserId, dsocialGroupId)
    if gGroupId == "" || err != nil {
        return false, err
    }
    original, _, err := ds.RetrieveExternalGroup(client.ServiceId(), userInfo.Guid(), dsocialUserId, gGroupId)
    if err != nil {
        return true, err
    }
    err = google.DeleteGroup(client, "", "", original.(*google.ContactGroup))
    return true, err
}

func (p *GoogleContactService) ContactsService() ContactsService {
    return p
}
