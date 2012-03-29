package contacts

import (
    "fmt"
    dm "github.com/pomack/dsocial.go/models/dsocial"
    "github.com/pomack/jsonhelper.go/jsonhelper"
    "github.com/pomack/oauth2_client.go/oauth2_client"
    "os"
    "time"
)

type Contact struct {
    ExternalServiceId string
    ExternalUserId    string
    ExternalContactId string
    DsocialUserId     string
    DsocialContactId  string
    Value             *dm.Contact
}

type Group struct {
    ExternalServiceId string
    ExternalUserId    string
    ExternalGroupId   string
    DsocialUserId     string
    DsocialGroupId    string
    Value             *dm.Group
}

type ContactsServiceSettings interface {
    Id() string
    DsocialUserId() string
    ContactsServiceId() string
    ExternalUserId() string
    ClientProperties() jsonhelper.JSONObject
    AllowRetrieveContactInfo() bool
    AllowModifyContactInfo() bool
    SetId(id string)
    SetDsocialUserId(dsocialUserId string)
    SetClientProperties(obj jsonhelper.JSONObject)
    SetExternalUserId(externalUserId string)
    SetAllowRetrieveContactInfo(allow bool)
    SetAllowModifyContactInfo(allow bool)
}

type NextToken interface{}

type DataStoreService interface {
    RetrieveAllContactsServiceSettingsForUser(dsocialUserId string) (settings []ContactsServiceSettings, err os.Error)
    RetrieveContactsServiceSettingsForService(dsocialUserId, contactsServiceId string) (settings []ContactsServiceSettings, err os.Error)
    RetrieveContactsServiceSettings(dsocialUserId, contactsServiceId, id string) (settings ContactsServiceSettings, err os.Error)
    SetContactsServiceSettings(settings ContactsServiceSettings) (id string, err os.Error)
    DeleteContactsServiceSettings(dsocialUserId, contactsServiceId, id string) (err os.Error)

    SearchForDsocialContacts(dsocialUserId string, contact *dm.Contact) (contacts []*dm.Contact, err os.Error)
    SearchForDsocialGroups(dsocialUserId string, groupName string) (groups []*dm.Group, err os.Error)

    StoreContactChangeSet(dsocialUserId string, changeset *dm.ChangeSet) (*dm.ChangeSet, os.Error)
    RetrieveContactChangeSets(dsocialId string, after *time.Time) ([]*dm.ChangeSet, NextToken, os.Error)

    StoreGroupChangeSet(dsocialUserId string, changeset *dm.ChangeSet) (*dm.ChangeSet, os.Error)
    RetrieveGroupChangeSets(dsocialId string, after *time.Time) ([]*dm.ChangeSet, NextToken, os.Error)

    AddContactChangeSetsToApply(dsocialUserId, serviceId, serviceName string, changesetIds []string) (id string, err os.Error)
    AddGroupChangeSetsToApply(dsocialUserId, serviceId, serviceName string, changesetIds []string) (id string, err os.Error)
    AddContactChangeSetsNotCurrentlyApplyable(dsocialUserId, serviceId, serviceName string, changesetIds []string) (id string, err os.Error)
    AddGroupChangeSetsNotCurrentlyApplyable(dsocialUserId, serviceId, serviceName string, changesetIds []string) (id string, err os.Error)

    RetrieveContactChangeSetsToApply(dsocialUserId, serviceId, serviceName string) ([]*dm.ChangeSetsToApply, map[string]*dm.ChangeSet, os.Error)
    RetrieveGroupChangeSetsToApply(dsocialUserId, serviceId, serviceName string) ([]*dm.ChangeSetsToApply, map[string]*dm.ChangeSet, os.Error)
    RetrieveContactChangeSetsNotCurrentlyApplyable(dsocialUserId, serviceId, serviceName string) ([]*dm.ChangeSetsToApply, map[string]*dm.ChangeSet, os.Error)
    RetrieveGroupChangeSetsNotCurrentlyApplyable(dsocialUserId, serviceId, serviceName string) ([]*dm.ChangeSetsToApply, map[string]*dm.ChangeSet, os.Error)

    RemoveContactChangeSetsToApply(dsocialUserId, serviceId, serviceName string, changeSetIdsToApply []string) os.Error
    RemoveGroupChangeSetsToApply(dsocialUserId, serviceId, serviceName string, changeSetIdsToApply []string) (err os.Error)
    RemoveContactChangeSetsNotCurrentlyApplyable(dsocialUserId, serviceId, serviceName string, changeSetIdsToApply []string) (err os.Error)
    RemoveGroupChangeSetsNotCurrentlyApplyable(dsocialUserId, serviceId, serviceName string, changeSetIdsToApply []string) (err os.Error)

    // Generates a new unique id for the specified collection name
    GenerateId(dsocialUserId, collectionName string) string

    // Retrieve the dsocial contact id for the specified external service/user id/contact id combo
    // Returns:
    //   dsocialContactId : the dsocial contact id if it exists or empty if not found
    //   err : error or nil
    DsocialIdForExternalContactId(externalServiceId, externalUserId, dsocialUserId, externalContactId string) (dsocialContactId string, err os.Error)
    // Retrieve the dsocial group id for the specified external service/user id/group id combo
    // Returns:
    //   dsocialGroupId : the dsocial group id if it exists or empty if not found
    //   err : error or nil
    DsocialIdForExternalGroupId(externalServiceId, externalUserId, dsocialUserId, externalGroupId string) (dsocialGroupId string, err os.Error)
    // Retrieve the external contact id for the specified external service/external user id/dsocial user id/dsocial contact id combo
    // Returns:
    //   externalContactId : the dsocial contact id if it exists or empty if not found
    //   err : error or nil
    ExternalContactIdForDsocialId(externalServiceId, externalUserId, dsocialUserId, dsocialContactId string) (externalContactId string, err os.Error)
    // Retrieve the external group id for the specified external service/external user id/dsocial user id/dsocial group id combo
    // Returns:
    //   externalGroupId : the dsocial group id if it exists or empty if not found
    //   err : error or nil
    ExternalGroupIdForDsocialId(externalServiceId, externalUserId, dsocialUserId, dsocialGroupId string) (externalGroupId string, err os.Error)
    // Stores the dsocial contact id <-> external contact id mapping
    // Returns:
    //   externalExisted : whether the external contact id mapping already existed and was overwritten
    //   dsocialExisted : whether the dsocial contact id mapping already existed and was overwritten
    //   err : error or nil
    StoreDsocialExternalContactMapping(externalServiceId, externalUserId, externalContactId, dsocialUserId, dsocialContactId string) (externalExisted, dsocialExisted bool, err os.Error)
    // Stores the dsocial contact id <-> external group id mapping
    // Returns:
    //   externalExisted : whether the external group id mapping already existed and was overwritten
    //   dsocialExisted : whether the dsocial group id mapping already existed and was overwritten
    //   err : error or nil
    StoreDsocialExternalGroupMapping(externalServiceId, externalUserId, externalGroupId, dsocialUserId, dsocialGroupId string) (externalExisted, dsocialExisted bool, err os.Error)

    // Retrieve external contact
    // Returns:
    //   externalContact : the contact as stored into the service using StoreExternalContact or nil if not found
    //   id : the internal id used to store the external contact
    //   err : error or nil
    RetrieveExternalContact(externalServiceId, externalUserId, dsocialUserId, externalContactId string) (externalContact interface{}, id string, err os.Error)
    // Retrieve external group
    // Returns:
    //   externalGroup : the group as stored into the service using StoreExternalGroup or nil if not found
    //   id : the internal id used to store the external group
    //   err : error or nil
    RetrieveExternalGroup(externalServiceId, externalUserId, dsocialUserId, externalGroupId string) (externalGroup interface{}, id string, err os.Error)
    // Stores external contact
    // Returns:
    //   id : the internal id used to store the external contact
    //   err : error or nil
    StoreExternalContact(externalServiceId, externalUserId, dsocialUserId, externalContactId string, contact interface{}) (id string, err os.Error)
    // Stores external group
    // Returns:
    //   id : the internal id used to store the external group
    //   err : error or nil
    StoreExternalGroup(externalServiceId, externalUserId, dsocialUserId, externalGroupId string, group interface{}) (id string, err os.Error)
    // Deletes external contact
    // Returns:
    //   existed : whether the contact existed upon deletion
    //   err : error or nil
    DeleteExternalContact(externalServiceId, externalUserId, dsocialUserId, externalContactId string) (existed bool, err os.Error)
    // Deletes external group
    // Returns:
    //   existed : whether the group existed upon deletion
    //   err : error or nil
    DeleteExternalGroup(externalServiceId, externalUserId, dsocialUserId, externalGroupId string) (existed bool, err os.Error)

    // Retrieve dsocial contact
    // Returns:
    //   dsocialContact : the contact as stored into the service using StoreDsocialContact or nil if not found
    //   id : the internal id used to store the dsocial contact
    //   err : error or nil
    RetrieveDsocialContactForExternalContact(externalServiceId, externalUserId, externalContactId, dsocialUserId string) (dsocialContact *dm.Contact, id string, err os.Error)
    // Retrieve dsocial group
    // Returns:
    //   dsocialGroup : the group as stored into the service using StoreDsocialGroup or nil if not found
    //   id : the internal id used to store the dsocial group
    //   err : error or nil
    RetrieveDsocialGroupForExternalGroup(externalServiceId, externalUserId, externalGroupId, dsocialUserId string) (dsocialGroup *dm.Group, id string, err os.Error)
    // Stores dsocial contact
    // Returns:
    //   dsocialContact : the contact, modified to include items like Id and LastModified/Created
    //   err : error or nil
    StoreDsocialContactForExternalContact(externalServiceId, externalUserId, externalContactId, dsocialUserId string, contact *dm.Contact) (dsocialContact *dm.Contact, err os.Error)
    // Stores dsocial group
    // Returns:
    //   dsocialGroup : the group, modified to include items like Id and LastModified/Created
    //   err : error or nil
    StoreDsocialGroupForExternalGroup(externalServiceId, externalUserId, externalGroupId, dsocialUserId string, group *dm.Group) (dsocialGroup *dm.Group, err os.Error)
    // Deletes dsocial contact
    // Returns:
    //   existed : whether the contact existed upon deletion
    //   err : error or nil
    DeleteDsocialContactForExternalContact(externalServiceId, externalUserId, externalContactId, dsocialUserId string) (existed bool, err os.Error)
    // Deletes dsocial group
    // Returns:
    //   existed : whether the group existed upon deletion
    //   err : error or nil
    DeleteDsocialGroupForExternalGroup(externalServiceId, externalUserId, externalGroupId, dsocialUserId string) (existed bool, err os.Error)

    // List dsocial contacts
    // Returns:
    //   dsocialContacts : list of contacts stored in the service
    //   next : token for the next list of contacts
    //   err : error or nil
    ListDsocialContacts(dsocialUserId string, from NextToken, maxCount int) (dsocialContacts []*dm.Contact, next NextToken, err os.Error)

    // Retrieve dsocial contact
    // Returns:
    //   dsocialContact : the contact as stored into the service using StoreDsocialContact or nil if not found
    //   id : the internal id used to store the dsocial contact
    //   err : error or nil
    RetrieveDsocialContact(dsocialUserId, dsocialContactId string) (dsocialContact *dm.Contact, id string, err os.Error)
    // Retrieve dsocial group
    // Returns:
    //   dsocialGroup : the group as stored into the service using StoreDsocialGroup or nil if not found
    //   id : the internal id used to store the dsocial group
    //   err : error or nil
    RetrieveDsocialGroup(dsocialUserId, dsocialGroupId string) (dsocialGroup *dm.Group, id string, err os.Error)
    // Stores dsocial contact
    // Returns:
    //   dsocialContact : the contact, modified to include items like Id and LastModified/Created
    //   err : error or nil
    StoreDsocialContact(dsocialUserId, dsocialContactId string, contact *dm.Contact) (dsocialContact *dm.Contact, err os.Error)
    // Stores dsocial group
    // Returns:
    //   dsocialGroup : the group, modified to include items like Id and LastModified/Created
    //   err : error or nil
    StoreDsocialGroup(dsocialUserId, dsocialGroupId string, group *dm.Group) (dsocialGroup *dm.Group, err os.Error)
    // Deletes dsocial contact
    // Returns:
    //   existed : whether the contact existed upon deletion
    //   err : error or nil
    DeleteDsocialContact(dsocialUserId, dsocialContactId string) (existed bool, err os.Error)
    // Deletes dsocial group
    // Returns:
    //   existed : whether the group existed upon deletion
    //   err : error or nil
    DeleteDsocialGroup(dsocialUserId, dsocialGroupId string) (existed bool, err os.Error)
}

type ContactsService interface {
    ServiceId() string
    // Create an OAuth2Client based on the specified settings for this contacts service
    CreateOAuth2Client(settings jsonhelper.JSONObject) (client oauth2_client.OAuth2Client, err os.Error)
    // Convert the external contact for this Contacts Service to a dsocial contact or nil if not convertible or input is nil
    ConvertToDsocialContact(externalContact interface{}, originalDsocialContact *dm.Contact, dsocialUserId string) (dsocialContact *dm.Contact)
    // Convert the dsocial contact to the external contact for this Contacts Service or nil if input is nil
    ConvertToExternalContact(dsocialContact *dm.Contact, originalExternalContact interface{}, dsocialUserId string) (externalContact interface{})
    // Convert the external group for this Contacts Service to a dsocial group or nil if not convertible or input is nil
    ConvertToDsocialGroup(externalGroup interface{}, originalDsocialGroup *dm.Group, dsocialUserId string) (dsocialGroup *dm.Group)
    // Convert the dsocial group to the external group for this Contacts Service or nil if input is nil
    ConvertToExternalGroup(dsocialGroup *dm.Group, originalExternalGroup interface{}, dsocialUserId string) (externalGroup interface{})

    // Whether this service can retrieve all contacts at once
    CanRetrieveAllContacts() bool
    // Whether this service can retrieve all connections (partial contact info or even just ids) at once
    CanRetrieveAllConnections() bool
    // Whether this service can retrieve all groups at once
    CanRetrieveAllGroups() bool
    // Whether this service can retrieve all contacts using paging
    CanRetrieveContacts() bool
    // Whether this service can retrieve all connections (partial contact info or even just ids) using paging
    CanRetrieveConnections() bool
    // Whether this service can retrieve all groups using paging
    CanRetrieveGroups() bool
    // Whether this service supports retrieve contact or group information from the service
    CanImportContactsOrGroups() bool
    // Whether this service supports modifications to contact or group information to the service
    CanExportContactsOrGroups() bool
    // Whether this service can retrieve contact information, either for self or for others
    CanRetrieveContact(selfContact bool) bool
    // Whether this service can create contact information, either for self or for others
    CanCreateContact(selfContact bool) bool
    // Whether this service can update contact information, either for self or for others
    CanUpdateContact(selfContact bool) bool
    // Whether this service can delete contact information, either for self or for others
    CanDeleteContact(selfContact bool) bool
    // Whether this service can retrieve group information, either for self or for others
    CanRetrieveGroup(selfContact bool) bool
    // Whether this service can create group information, either for self or for others
    CanCreateGroup(selfContact bool) bool
    // Whether this service can update group information, either for self or for others
    CanUpdateGroup(selfContact bool) bool
    // Whether this service can delete group information, either for self or for others
    CanDeleteGroup(selfContact bool) bool
    // Whether this service shows group memberships when retrieving a list of groups
    GroupListIncludesContactIds() bool
    // Whether this service shows group memberships when retrieving a single group
    GroupInfoIncludesContactIds() bool
    // Whether this service shows group memberships when retrieving a single contact
    ContactInfoIncludesGroups() bool

    // Retrieve all contacts using the specified client
    // Returns:
    //   contacts : all contacts
    //   err : error or nil
    RetrieveAllContacts(client oauth2_client.OAuth2Client, ds DataStoreService, dsocialUserId string) (contacts []*Contact, err os.Error)
    // Retrieve all connections using the specified client
    // Returns:
    //   connections : all connections
    //   err : error or nil
    RetrieveAllConnections(client oauth2_client.OAuth2Client, ds DataStoreService, dsocialUserId string) (connections []*Contact, err os.Error)
    // Retrieve all groups using the specified client
    // Returns:
    //   groups : all groups
    //   err : error or nil
    RetrieveAllGroups(client oauth2_client.OAuth2Client, ds DataStoreService, dsocialUserId string) (groups []*Group, err os.Error)
    // Retrieve contacts using next as an opaque pointer for where to start listing from using the specified client
    // Returns:
    //   contacts : contacts
    //   nextToken : token to the next page, if contacts are empty then no more exist and nextToken is irrelevant
    //   err : error or nil
    RetrieveContacts(client oauth2_client.OAuth2Client, ds DataStoreService, dsocialUserId string, next NextToken) (contacts []*Contact, nextToken NextToken, err os.Error)
    // Retrieve connections using next as an opaque pointer for where to start listing from using the specified client
    // Returns:
    //   connections : connections
    //   nextToken : token to the next page, if connections are empty then no more exist and nextToken is irrelevant
    //   err : error or nil
    RetrieveConnections(client oauth2_client.OAuth2Client, ds DataStoreService, dsocialUserId string, next NextToken) (connections []*Contact, nextToken NextToken, err os.Error)
    // Retrieve groups using next as an opaque pointer for where to start listing from using the specified client
    // Returns:
    //   groups : groups
    //   nextToken : token to the next page, if groups are empty then no more exist and nextToken is irrelevant
    //   err : error or nil
    RetrieveGroups(client oauth2_client.OAuth2Client, ds DataStoreService, dsocialUserId string, next NextToken) (groups []*Group, nextToken NextToken, err os.Error)
    // Retrieve the specified contact for the contactId or self-contact if contactId is empty
    // Returns:
    //   contact : contact or nil if not found
    //   err : error or nil
    RetrieveContact(client oauth2_client.OAuth2Client, ds DataStoreService, dsocialUserId string, contactId string) (contact *Contact, err os.Error)
    // Retrieve the specified group for the groupId
    // Returns:
    //   group : group or nil if not found
    //   err : error or nil
    RetrieveGroup(client oauth2_client.OAuth2Client, ds DataStoreService, dsocialUserId string, groupId string) (group *Group, err os.Error)
    // Creates the specified contact
    // Returns:
    //   updatedContact : updated version of contact with fields updated like Id and LastModified
    //   err : error or nil
    CreateContactOnExternalService(client oauth2_client.OAuth2Client, externalContact interface{}) (externalContactResult interface{}, externalContactId string, err os.Error)
    //CreateContact(client oauth2_client.OAuth2Client, ds DataStoreService, dsocialUserId string, contact *dm.Contact) (updatedContact *Contact, err os.Error)
    // Creates the specified group
    // Returns:
    //   updatedGroup : updated version of group with fields updated like Id and LastModified
    //   err : error or nil
    CreateGroupOnExternalService(client oauth2_client.OAuth2Client, externalGroup interface{}) (externalGroupResult interface{}, externalGroupId string, err os.Error)
    //CreateGroup(client oauth2_client.OAuth2Client, ds DataStoreService, dsocialUserId string, group *dm.Group) (updatedGroup *Group, err os.Error)
    // Updates the specified contact
    // Returns:
    //   updatedContact : updated version of contact with fields updated like LastModified
    //   err : error or nil
    UpdateContactOnExternalService(client oauth2_client.OAuth2Client, originalExternalContact, latestExternalContact interface{}) (externalContactResult interface{}, externalContactId string, err os.Error)
    //UpdateContact(client oauth2_client.OAuth2Client, ds DataStoreService, dsocialUserId string, originalContact, contact *dm.Contact) (updatedContact *Contact, err os.Error)
    // Updates the specified group
    // Returns:
    //   updatedGroup : updated version of group with fields updated like LastModified
    //   err : error or nil
    UpdateGroupOnExternalService(client oauth2_client.OAuth2Client, originalExternalGroup, latestExternalGroup interface{}) (externalGroupResult interface{}, externalGroupId string, err os.Error)
    //UpdateGroup(client oauth2_client.OAuth2Client, ds DataStoreService, dsocialUserId string, originalGroup, group *dm.Group) (updatedGroup *Group, err os.Error)
    // Deletes the specified contact
    // Returns:
    //   existed : whether the contact existed upon deletiong
    //   err : error or nil
    DeleteContactOnExternalService(client oauth2_client.OAuth2Client, originalExternalContact interface{}) (existed bool, err os.Error)
    //DeleteContact(client oauth2_client.OAuth2Client, ds DataStoreService, dsocialUserId, dsocialContactId string) (existed bool, err os.Error)
    // Deletes the specified group
    // Returns:
    //   existed : whether the group existed upon deletiong
    //   err : error or nil
    DeleteGroupOnExternalService(client oauth2_client.OAuth2Client, originalExternalGroup interface{}) (existed bool, err os.Error)
    //DeleteGroup(client oauth2_client.OAuth2Client, ds DataStoreService, dsocialUserId, dsocialGroupId string) (existed bool, err os.Error)
}

type StandardContactServiceSettings struct {
    FieldId                       string                `json:"id"`
    FieldDsocialUserId            string                `json:"dsocial_user_id"`
    FieldExternalUserId           string                `json:"external_user_id"`
    FieldClientProperties         jsonhelper.JSONObject `json:"client_properties"`
    FieldAllowRetrieveContactInfo bool                  `json:"allow_retrieve_contact_info"`
    FieldAllowModifyContactInfo   bool                  `json:"allow_modify_contact_info"`
}

func (p *StandardContactServiceSettings) Id() string {
    return p.FieldId
}

func (p *StandardContactServiceSettings) DsocialUserId() string {
    return p.FieldDsocialUserId
}

func (p *StandardContactServiceSettings) ExternalUserId() string {
    return p.FieldExternalUserId
}

func (p *StandardContactServiceSettings) ClientProperties() jsonhelper.JSONObject {
    return p.FieldClientProperties
}

func (p *StandardContactServiceSettings) AllowRetrieveContactInfo() bool {
    return p.FieldAllowRetrieveContactInfo
}

func (p *StandardContactServiceSettings) AllowModifyContactInfo() bool {
    return p.FieldAllowModifyContactInfo
}

func (p *StandardContactServiceSettings) SetId(id string) {
    p.FieldId = id
}

func (p *StandardContactServiceSettings) SetDsocialUserId(dsocialUserId string) {
    p.FieldDsocialUserId = dsocialUserId
}

func (p *StandardContactServiceSettings) SetClientProperties(obj jsonhelper.JSONObject) {
    p.FieldClientProperties = obj
}

func (p *StandardContactServiceSettings) SetExternalUserId(externalUserId string) {
    p.FieldExternalUserId = externalUserId
}

func (p *StandardContactServiceSettings) SetAllowRetrieveContactInfo(allow bool) {
    p.FieldAllowRetrieveContactInfo = allow
}

func (p *StandardContactServiceSettings) SetAllowModifyContactInfo(allow bool) {
    p.FieldAllowModifyContactInfo = allow
}

func addIdForAclPersistableModel(m *dm.AclPersistableModel, ds DataStoreService, collectionName string, ownerId string) {
    if m == nil {
        return
    }
    if m.Acl.OwnerId == "" {
        m.Acl.OwnerId = ownerId
    }
    if m.Id == "" {
        m.Id = ds.GenerateId(ownerId, collectionName)
    }
}

func AddIdsForDsocialContact(c *dm.Contact, ds DataStoreService, dsocialUserId string) (err os.Error) {
    if c == nil {
        return
    }
    if c.UserId == "" {
        c.UserId = dsocialUserId
    }
    if c.Acl.OwnerId == "" {
        c.Acl.OwnerId = dsocialUserId
    }
    if c.Id == "" {
        c.Id = ds.GenerateId(dsocialUserId, "contact")
    }
    if c.PostalAddresses != nil {
        for _, addr := range c.PostalAddresses {
            if addr.Acl.OwnerId == "" {
                addr.Acl.OwnerId = dsocialUserId
            }
            if addr.Id == "" {
                addr.Id = ds.GenerateId(dsocialUserId, "address")
            }
        }
    }
    if c.Educations != nil {
        for _, ed := range c.Educations {
            if ed.Acl.OwnerId == "" {
                ed.Acl.OwnerId = dsocialUserId
            }
            if ed.Id == "" {
                ed.Id = ds.GenerateId(dsocialUserId, "education")
            }
        }
    }
    if c.WorkHistories != nil {
        for _, wh := range c.WorkHistories {
            if wh.Acl.OwnerId == "" {
                wh.Acl.OwnerId = dsocialUserId
            }
            if wh.Id == "" {
                wh.Id = ds.GenerateId(dsocialUserId, "workhistory")
            }
        }
    }
    if c.PhoneNumbers != nil {
        for _, p := range c.PhoneNumbers {
            if p.Acl.OwnerId == "" {
                p.Acl.OwnerId = dsocialUserId
            }
            if p.Id == "" {
                p.Id = ds.GenerateId(dsocialUserId, "phone")
            }
        }
    }
    if c.EmailAddresses != nil {
        for _, e := range c.EmailAddresses {
            if e.Acl.OwnerId == "" {
                e.Acl.OwnerId = dsocialUserId
            }
            if e.Id == "" {
                e.Id = ds.GenerateId(dsocialUserId, "email")
            }
        }
    }
    if c.Uris != nil {
        for _, u := range c.Uris {
            if u.Acl.OwnerId == "" {
                u.Acl.OwnerId = dsocialUserId
            }
            if u.Id == "" {
                u.Id = ds.GenerateId(dsocialUserId, "uri")
            }
        }
    }
    if c.Ims != nil {
        for _, im := range c.Ims {
            if im.Acl.OwnerId == "" {
                im.Acl.OwnerId = dsocialUserId
            }
            if im.Id == "" {
                im.Id = ds.GenerateId(dsocialUserId, "im")
            }
        }
    }
    if c.Relationships != nil {
        for _, r := range c.Relationships {
            if r.Acl.OwnerId == "" {
                r.Acl.OwnerId = dsocialUserId
            }
            if r.Id == "" {
                r.Id = ds.GenerateId(dsocialUserId, "relationship")
            }
        }
    }
    if c.Dates != nil {
        for _, d := range c.Dates {
            if d.Acl.OwnerId == "" {
                d.Acl.OwnerId = dsocialUserId
            }
            if d.Id == "" {
                d.Id = ds.GenerateId(dsocialUserId, "date")
            }
        }
    }
    if c.DateTimes != nil {
        for _, d := range c.DateTimes {
            if d.Acl.OwnerId == "" {
                d.Acl.OwnerId = dsocialUserId
            }
            if d.Id == "" {
                d.Id = ds.GenerateId(dsocialUserId, "datetime")
            }
        }
    }
    if c.Certifications != nil {
        for _, cert := range c.Certifications {
            if cert.Acl.OwnerId == "" {
                cert.Acl.OwnerId = dsocialUserId
            }
            if cert.Id == "" {
                cert.Id = ds.GenerateId(dsocialUserId, "certification")
            }
        }
    }
    if c.Skills != nil {
        for _, s := range c.Skills {
            if s.Acl.OwnerId == "" {
                s.Acl.OwnerId = dsocialUserId
            }
            if s.Id == "" {
                s.Id = ds.GenerateId(dsocialUserId, "skill")
            }
        }
    }
    if c.Languages != nil {
        for _, l := range c.Languages {
            if l.Acl.OwnerId == "" {
                l.Acl.OwnerId = dsocialUserId
            }
            if l.Id == "" {
                l.Id = ds.GenerateId(dsocialUserId, "language")
            }
        }
    }
    return
}

func AddIdsForDsocialGroup(g *dm.Group, ds DataStoreService, dsocialUserId string) (err os.Error) {
    if g == nil {
        return
    }
    if g.UserId == "" {
        g.UserId = dsocialUserId
    }
    if g.Acl.OwnerId == "" {
        g.Acl.OwnerId = dsocialUserId
    }
    if g.Id == "" {
        g.Id = ds.GenerateId(dsocialUserId, "group")
    }
    return
}

func CreateContactOnExternalService(client oauth2_client.OAuth2Client, cs ContactsService, ds DataStoreService, dsocialUserId string, contact *dm.Contact) (*Contact, os.Error) {
    if contact == nil {
        return nil, nil
    }
    userInfo, err := client.RetrieveUserInfo()
    if err != nil {
        return nil, err
    }
    externalServiceId := cs.ServiceId()
    externalUserId := userInfo.Guid()
    externalContactId, err := ds.ExternalContactIdForDsocialId(externalServiceId, externalUserId, dsocialUserId, contact.Id)
    if err != nil {
        return nil, err
    }
    if externalContactId != "" {
        originalContact, _, err := ds.RetrieveDsocialContact(dsocialUserId, contact.Id)
        if err != nil {
            return nil, err
        }
        return UpdateContactOnExternalService(client, cs, ds, dsocialUserId, originalContact, contact)
    }
    externalContact := cs.ConvertToExternalContact(contact, nil, dsocialUserId)
    externalContact, externalContactId, err = cs.CreateContactOnExternalService(client, externalContact)
    if err != nil {
        return nil, err
    }
    externalContactId2, err := ds.StoreExternalContact(externalServiceId, externalUserId, dsocialUserId, externalContactId, externalContact)
    if err != nil {
        return nil, err
    }
    fmt.Printf("[SERVICE]: extContactId: %s, extContactId2: %s\n", externalContactId, externalContactId2)
    if externalContactId2 != "" {
        externalContactId = externalContactId2
    }
    dsocialContactForExternal := cs.ConvertToDsocialContact(externalContact, contact, dsocialUserId)
    dsocialContactForExternal, err = ds.StoreDsocialContactForExternalContact(externalServiceId, externalUserId, externalContactId, dsocialUserId, dsocialContactForExternal)
    if err != nil {
        return nil, err
    }
    _, _, err = ds.StoreDsocialExternalContactMapping(externalServiceId, externalUserId, externalContactId, dsocialUserId, dsocialContactForExternal.Id)
    outContact := &Contact{
        ExternalServiceId: externalServiceId,
        ExternalUserId:    externalUserId,
        ExternalContactId: externalContactId,
        DsocialUserId:     dsocialUserId,
        DsocialContactId:  dsocialContactForExternal.Id,
        Value:             dsocialContactForExternal,
    }
    return outContact, err
}

func CreateGroupOnExternalService(client oauth2_client.OAuth2Client, cs ContactsService, ds DataStoreService, dsocialUserId string, group *dm.Group) (*Group, os.Error) {
    if group == nil {
        return nil, nil
    }
    userInfo, err := client.RetrieveUserInfo()
    if err != nil {
        return nil, err
    }
    externalServiceId := cs.ServiceId()
    externalUserId := userInfo.Guid()
    externalGroupId, err := ds.ExternalGroupIdForDsocialId(externalServiceId, externalUserId, dsocialUserId, group.Id)
    if err != nil {
        return nil, err
    }
    if externalGroupId != "" {
        originalGroup, _, err := ds.RetrieveDsocialGroup(dsocialUserId, group.Id)
        if err != nil {
            return nil, err
        }
        return UpdateGroupOnExternalService(client, cs, ds, dsocialUserId, originalGroup, group)
    }
    externalGroup := cs.ConvertToExternalGroup(group, nil, dsocialUserId)
    externalGroup, externalGroupId, err = cs.CreateGroupOnExternalService(client, externalGroup)
    if err != nil {
        return nil, err
    }
    externalGroupId, err = ds.StoreExternalGroup(externalServiceId, externalUserId, dsocialUserId, externalGroupId, externalGroup)
    if err != nil {
        return nil, err
    }
    dsocialGroupForExternal := cs.ConvertToDsocialGroup(externalGroup, group, dsocialUserId)
    dsocialGroupForExternal, err = ds.StoreDsocialGroupForExternalGroup(externalServiceId, externalUserId, externalGroupId, dsocialUserId, dsocialGroupForExternal)
    if err != nil {
        return nil, err
    }
    _, _, err = ds.StoreDsocialExternalGroupMapping(externalServiceId, externalUserId, externalGroupId, dsocialUserId, dsocialGroupForExternal.Id)
    outGroup := &Group{
        ExternalServiceId: externalServiceId,
        ExternalUserId:    externalUserId,
        ExternalGroupId:   externalGroupId,
        DsocialUserId:     dsocialUserId,
        DsocialGroupId:    dsocialGroupForExternal.Id,
        Value:             dsocialGroupForExternal,
    }
    return outGroup, err
}

func UpdateContactOnExternalService(client oauth2_client.OAuth2Client, cs ContactsService, ds DataStoreService, dsocialUserId string, originalContact, contact *dm.Contact) (*Contact, os.Error) {
    if contact == nil || originalContact == nil {
        return nil, nil
    }
    userInfo, err := client.RetrieveUserInfo()
    if err != nil {
        return nil, err
    }
    externalServiceId := cs.ServiceId()
    externalUserId := userInfo.Guid()
    externalContactId, err := ds.ExternalContactIdForDsocialId(externalServiceId, externalUserId, dsocialUserId, originalContact.Id)
    if err != nil {
        return nil, err
    }
    if externalContactId == "" {
        return CreateContactOnExternalService(client, cs, ds, dsocialUserId, contact)
    }
    originalExternalContact, _, err := ds.RetrieveExternalContact(externalServiceId, externalUserId, dsocialUserId, externalContactId)
    if err != nil {
        return nil, err
    }
    latestExternalContact := cs.ConvertToExternalContact(contact, originalExternalContact, dsocialUserId)
    latestExternalContact2, externalContactId2, err := cs.UpdateContactOnExternalService(client, originalExternalContact, latestExternalContact)
    if err != nil {
        return nil, err
    }
    if latestExternalContact2 != nil {
        latestExternalContact = latestExternalContact2
    }
    if externalContactId2 != "" {
        externalContactId = externalContactId2
    }
    latestDsocialContactForExternal := cs.ConvertToDsocialContact(latestExternalContact, originalContact, dsocialUserId)
    _, err = ds.StoreExternalContact(externalServiceId, externalUserId, dsocialUserId, externalContactId, latestExternalContact)
    if err != nil {
        return nil, err
    }
    latestDsocialContactForExternal, err = ds.StoreDsocialContactForExternalContact(externalServiceId, externalUserId, externalContactId, dsocialUserId, latestDsocialContactForExternal)
    outContact := &Contact{
        ExternalServiceId: externalServiceId,
        ExternalUserId:    externalUserId,
        ExternalContactId: externalContactId,
        DsocialUserId:     dsocialUserId,
        DsocialContactId:  latestDsocialContactForExternal.Id,
        Value:             latestDsocialContactForExternal,
    }
    return outContact, err
}

func UpdateGroupOnExternalService(client oauth2_client.OAuth2Client, cs ContactsService, ds DataStoreService, dsocialUserId string, originalGroup, group *dm.Group) (*Group, os.Error) {
    if group == nil || originalGroup == nil {
        return nil, nil
    }
    userInfo, err := client.RetrieveUserInfo()
    if err != nil {
        return nil, err
    }
    externalServiceId := cs.ServiceId()
    externalUserId := userInfo.Guid()
    externalGroupId, err := ds.ExternalGroupIdForDsocialId(externalServiceId, externalUserId, dsocialUserId, originalGroup.Id)
    if err != nil {
        return nil, err
    }
    if externalGroupId == "" {
        return CreateGroupOnExternalService(client, cs, ds, dsocialUserId, group)
    }
    originalExternalGroup, _, err := ds.RetrieveExternalGroup(externalServiceId, externalUserId, dsocialUserId, externalGroupId)
    if err != nil {
        return nil, err
    }
    latestExternalGroup := cs.ConvertToExternalGroup(group, originalExternalGroup, dsocialUserId)
    latestExternalGroup2, externalGroupId2, err := cs.UpdateGroupOnExternalService(client, originalExternalGroup, latestExternalGroup)
    if err != nil {
        return nil, err
    }
    if latestExternalGroup2 != nil {
        latestExternalGroup = latestExternalGroup2
    }
    if externalGroupId2 != "" {
        externalGroupId = externalGroupId2
    }
    latestDsocialGroupForExternal := cs.ConvertToDsocialGroup(latestExternalGroup, originalGroup, dsocialUserId)
    _, err = ds.StoreExternalGroup(externalServiceId, externalUserId, dsocialUserId, externalGroupId, latestExternalGroup)
    if err != nil {
        return nil, err
    }
    latestDsocialGroupForExternal, err = ds.StoreDsocialGroupForExternalGroup(externalServiceId, externalUserId, externalGroupId, dsocialUserId, latestDsocialGroupForExternal)
    outGroup := &Group{
        ExternalServiceId: externalServiceId,
        ExternalUserId:    externalUserId,
        ExternalGroupId:   externalGroupId,
        DsocialUserId:     dsocialUserId,
        DsocialGroupId:    latestDsocialGroupForExternal.Id,
        Value:             latestDsocialGroupForExternal,
    }
    return outGroup, err
}

func DeleteContactOnExternalService(client oauth2_client.OAuth2Client, cs ContactsService, ds DataStoreService, dsocialUserId, dsocialContactId string) (bool, os.Error) {
    if dsocialContactId == "" || dsocialUserId == "" {
        return false, nil
    }
    userInfo, err := client.RetrieveUserInfo()
    if err != nil {
        return true, err
    }
    externalServiceId := cs.ServiceId()
    externalUserId := userInfo.Guid()
    externalContactId, err := ds.ExternalContactIdForDsocialId(externalServiceId, externalUserId, dsocialUserId, dsocialContactId)
    if externalContactId == "" || err != nil {
        return externalContactId == "", err
    }
    externalContact, _, err := ds.RetrieveExternalContact(externalServiceId, externalUserId, dsocialUserId, externalContactId)
    if err != nil {
        return true, err
    }
    if externalContact == nil {
        return false, err
    }
    _, err = cs.DeleteContactOnExternalService(client, externalContact)
    if err != nil {
        return true, err
    }
    _, err = ds.DeleteDsocialContactForExternalContact(externalServiceId, externalUserId, externalContactId, dsocialUserId)
    if err != nil {
        return true, err
    }
    _, err = ds.DeleteExternalContact(externalServiceId, externalUserId, dsocialUserId, externalContactId)
    return true, err
}

func DeleteGroupOnExternalService(client oauth2_client.OAuth2Client, cs ContactsService, ds DataStoreService, dsocialUserId, dsocialGroupId string) (bool, os.Error) {
    if dsocialGroupId == "" || dsocialUserId == "" {
        return false, nil
    }
    userInfo, err := client.RetrieveUserInfo()
    if err != nil {
        return true, err
    }
    externalServiceId := cs.ServiceId()
    externalUserId := userInfo.Guid()
    externalGroupId, err := ds.ExternalGroupIdForDsocialId(externalServiceId, externalUserId, dsocialUserId, dsocialGroupId)
    if externalGroupId == "" || err != nil {
        return externalGroupId == "", err
    }
    externalGroup, _, err := ds.RetrieveExternalGroup(externalServiceId, externalUserId, dsocialUserId, externalGroupId)
    if err != nil {
        return true, err
    }
    if externalGroup == nil {
        return false, err
    }
    _, err = cs.DeleteGroupOnExternalService(client, externalGroup)
    if err != nil {
        return true, err
    }
    _, err = ds.DeleteDsocialGroupForExternalGroup(externalServiceId, externalUserId, externalGroupId, dsocialUserId)
    if err != nil {
        return true, err
    }
    _, err = ds.DeleteExternalGroup(externalServiceId, externalUserId, dsocialUserId, externalGroupId)
    return true, err
}
