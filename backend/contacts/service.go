package contacts

import (
    "github.com/pomack/oauth2_client.go/oauth2_client"
    dm "github.com/pomack/dsocial.go/models/dsocial"
    "os"
    "time"
)

type Contact struct {
    ExternalServiceId string
    ExternalUserId string
    ExternalContactId string
    DsocialUserId string
    DsocialContactId string
    Value *dm.Contact
}

type Group struct {
    ExternalServiceId string
    ExternalUserId string
    ExternalGroupId string
    DsocialUserId string
    DsocialGroupId string
    Value *dm.Group
}


type NextToken interface{}

type DataStoreService interface {
    
    SearchForDsocialContacts(dsocialUserId string, contact *dm.Contact) (contacts []*dm.Contact, err os.Error)
    SearchForDsocialGroups(dsocialUserId string, groupName string) (groups []*dm.Group, err os.Error)
    
    StoreContactChangeSet(changeset *dm.ChangeSet) (*dm.ChangeSet, os.Error)
    RetrieveContactChangeSets(dsocialId string, after *time.Time) ([]*dm.ChangeSet, NextToken, os.Error)
    
    StoreGroupChangeSet(changeset *dm.ChangeSet) (*dm.ChangeSet, os.Error)
    RetrieveGroupChangeSets(dsocialId string, after *time.Time) ([]*dm.ChangeSet, NextToken, os.Error)
    
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
    
    // Convert the external contact for this Contacts Service to a dsocial contact or nil if not convertible or input is nil
    ConvertToDsocialContact(externalContact interface{}) (dsocialContact *dm.Contact)
    // Convert the dsocial contact to the external contact for this Contacts Service or nil if input is nil
    ConvertToExternalContact(dsocialContact *dm.Contact) (externalContact interface{})
    // Convert the external group for this Contacts Service to a dsocial group or nil if not convertible or input is nil
    ConvertToDsocialGroup(externalGroup interface{}) (dsocialGroup *dm.Group)
    // Convert the dsocial group to the external group for this Contacts Service or nil if input is nil
    ConvertToExternalGroup(dsocialGroup *dm.Group) (externalGroup interface{})
    
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
    CreateContact(client oauth2_client.OAuth2Client, ds DataStoreService, dsocialUserId string, contact *dm.Contact) (updatedContact *Contact, err os.Error)
    // Creates the specified group
    // Returns:
    //   updatedGroup : updated version of group with fields updated like Id and LastModified
    //   err : error or nil
    CreateGroup(client oauth2_client.OAuth2Client, ds DataStoreService, dsocialUserId string, group *dm.Group) (updatedGroup *Group, err os.Error)
    // Updates the specified contact
    // Returns:
    //   updatedContact : updated version of contact with fields updated like LastModified
    //   err : error or nil
    UpdateContact(client oauth2_client.OAuth2Client, ds DataStoreService, dsocialUserId string, originalContact, contact *dm.Contact) (updatedContact *Contact, err os.Error)
    // Updates the specified group
    // Returns:
    //   updatedGroup : updated version of group with fields updated like LastModified
    //   err : error or nil
    UpdateGroup(client oauth2_client.OAuth2Client, ds DataStoreService, dsocialUserId string, originalGroup, group *dm.Group) (updatedGroup *Group, err os.Error)
    // Deletes the specified contact
    // Returns:
    //   existed : whether the contact existed upon deletiong
    //   err : error or nil
    DeleteContact(client oauth2_client.OAuth2Client, ds DataStoreService, dsocialUserId, dsocialContactId string) (existed bool, err os.Error)
    // Deletes the specified group
    // Returns:
    //   existed : whether the group existed upon deletiong
    //   err : error or nil
    DeleteGroup(client oauth2_client.OAuth2Client, ds DataStoreService, dsocialUserId, dsocialGroupId string) (existed bool, err os.Error)
}