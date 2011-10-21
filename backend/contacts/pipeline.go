package contacts

import (
    "github.com/pomack/oauth2_client.go/oauth2_client"
    dm "github.com/pomack/dsocial.go/models/dsocial"
    "container/list"
    "container/vector"
    "os"
    "time"
)

type Pipeline struct {
}

func (p *Pipeline) InitialSync(client oauth2_client.OAuth2Client, ds DataStoreService, cs ContactsService, dsocialUserId string) os.Error {
    return p.IncrementalSync(client, ds, cs, dsocialUserId)
}

func (p *Pipeline) findMatchingDsocialContact(ds DataStoreService, dsocialUserId string, contact *Contact) (originalExternalContact *dm.Contact, isSame bool, err os.Error) {
    emptyContact := new(dm.Contact)
    if contact.DsocialContactId == "" {
        // this is a new contact from an existing service
        potentialMatches, err := ds.SearchForDsocialContacts(dsocialUserId, contact.Value)
        if err != nil {
            return nil, false, err
        }
        for _, potentialMatch := range potentialMatches {
            var isSimilar bool
            if isSimilar, isSame = emptyContact.IsSimilarOrUpdated(potentialMatch, contact.Value); isSimilar {
                originalExternalContact = potentialMatch
                break
            }
        }
        if originalExternalContact != nil {
            _, _, err = ds.StoreDsocialExternalContactMapping(contact.ExternalServiceId, contact.ExternalUserId, contact.ExternalContactId, dsocialUserId, originalExternalContact.Id)
            contact.DsocialContactId = originalExternalContact.Id
        }
    }
    return originalExternalContact, isSame, err
}

func (p *Pipeline) contactSync(ds DataStoreService, dsocialUserId string, contact *Contact) (*dm.Contact, os.Error) {
    emptyContact := new(dm.Contact)
    if contact == nil || contact.Value == nil {
        return nil, nil
    }
    matchingContact, isSame, err := p.findMatchingDsocialContact(ds, dsocialUserId, contact)
    if isSame || err != nil {
        return matchingContact, err
    }
    originalExternalContact, _, err := ds.RetrieveDsocialContactForExternalContact(contact.ExternalServiceId, contact.ExternalUserId, contact.ExternalContactId, dsocialUserId)
    if err != nil {
        return matchingContact, err
    }
    if _, isSame = emptyContact.IsSimilarOrUpdated(originalExternalContact, contact.Value); isSame {
        return matchingContact, nil
    }
    l := new(list.List)
    emptyContact.GenerateChanges(originalExternalContact, contact.Value, nil, l)
    changes := make([]*dm.Change, l.Len())
    changeset := &dm.ChangeSet{
        CreatedAt: time.UTC().Format(dm.UTC_DATETIME_FORMAT),
        ChangedBy: contact.ExternalServiceId,
        ChangeImportId: contact.ExternalContactId,
        RecordId: contact.DsocialContactId,
        Changes: changes,
    }
    _, err = ds.StoreContactChangeSet(changeset)
    if err != nil {
        return matchingContact, err
    }
    if originalExternalContact == nil {
        contact.Value, err = ds.StoreDsocialContact(dsocialUserId, "", contact.Value)
        if err != nil {
            return matchingContact, err
        }
        _, err = ds.StoreDsocialContactForExternalContact(contact.ExternalServiceId, contact.ExternalUserId, contact.ExternalContactId, dsocialUserId, contact.Value)
        return matchingContact, err
    }
    var storedDsocialContact *dm.Contact = nil
    if contact.DsocialContactId != "" {
        if storedDsocialContact, _, err = ds.RetrieveDsocialContact(dsocialUserId, contact.DsocialContactId); err != nil {
            return matchingContact, err
        }
    }
    if storedDsocialContact == nil {
        storedDsocialContact = new(dm.Contact)
    }
    for j, iter := 0, l.Front(); iter != nil; j, iter = j+1, iter.Next() {
        change := iter.Value.(*dm.Change)
        dm.ApplyChange(storedDsocialContact, change)
        changes[j] = change
    }
    _, err = ds.StoreDsocialContact(dsocialUserId, contact.DsocialContactId, storedDsocialContact)
    return storedDsocialContact, err
}

func (p *Pipeline) findMatchingDsocialGroup(ds DataStoreService, dsocialUserId string, group *Group) (originalExternalGroup *dm.Group, isSame bool, err os.Error) {
    emptyGroup := new(dm.Group)
    if group.DsocialGroupId == "" {
        // this is a new group from an existing service
        potentialMatches, err := ds.SearchForDsocialGroups(dsocialUserId, group.Value.Name)
        if err != nil {
            return nil, false, err
        }
        for _, potentialMatch := range potentialMatches {
            var isSimilar bool
            if isSimilar, isSame = emptyGroup.IsSimilarOrUpdated(potentialMatch, group.Value); isSimilar {
                originalExternalGroup = potentialMatch
                break
            }
        }
        if originalExternalGroup != nil {
            _, _, err = ds.StoreDsocialExternalGroupMapping(group.ExternalServiceId, group.ExternalUserId, group.ExternalGroupId, dsocialUserId, originalExternalGroup.Id)
            group.DsocialGroupId = originalExternalGroup.Id
        }
    }
    return originalExternalGroup, isSame, err
}

func (p *Pipeline) groupSync(ds DataStoreService, dsocialUserId string, group *Group, minimumIncludes *list.List) (*dm.Group, os.Error) {
    emptyGroup := new(dm.Group)
    if group == nil || group.Value == nil {
        return nil, nil
    }
    if group.Value.ContactIds == nil {
        group.Value.ContactIds = make([]string, 0, 10)
    }
    if group.Value.ContactNames == nil {
        group.Value.ContactNames = make([]string, 0, 10)
    }
    if len(group.Value.ContactIds) == 0 || len(group.Value.ContactIds) == 0 && minimumIncludes != nil {
        sv1 := vector.StringVector(group.Value.ContactIds)
        sv2 := vector.StringVector(group.Value.ContactNames)
        sv1.Resize(sv1.Len(), sv1.Len() + minimumIncludes.Len())
        sv2.Resize(sv2.Len(), sv2.Len() + minimumIncludes.Len())
        for iter := minimumIncludes.Front(); iter != nil; iter = iter.Next() {
            contactRef := iter.Value.(*dm.ContactRef)
            sv1.Push(contactRef.Id)
            sv2.Push(contactRef.Name)
        }
    } else if minimumIncludes == nil {
        for iter := minimumIncludes.Front(); iter != nil; iter = iter.Next() {
            contactRef := iter.Value.(*dm.ContactRef)
            refLocation := -1
            if contactRef.Id != "" {
                for i, id := range group.Value.ContactIds {
                    if id == contactRef.Id {
                        refLocation = i
                        break
                    }
                }
            }
            if refLocation == -1 && contactRef.Name != "" {
                for i, name := range group.Value.ContactNames {
                    if name == contactRef.Name {
                        refLocation = i
                        break
                    }
                }
            }
            if refLocation == -1 {
                sv1 := vector.StringVector(group.Value.ContactIds)
                sv2 := vector.StringVector(group.Value.ContactNames)
                sv1.Push(contactRef.Id)
                sv2.Push(contactRef.Name)
            } else {
                group.Value.ContactIds[refLocation] = contactRef.Id
                group.Value.ContactNames[refLocation] = contactRef.Name
            }
        }
    }
    matchingGroup, isSame, err := p.findMatchingDsocialGroup(ds, dsocialUserId, group)
    if isSame || err != nil {
        return matchingGroup, err
    }
    originalExternalGroup, _, err := ds.RetrieveDsocialGroupForExternalGroup(group.ExternalServiceId, group.ExternalUserId, group.ExternalGroupId, dsocialUserId)
    if err != nil {
        return matchingGroup, err
    }
    if _, isSame = emptyGroup.IsSimilarOrUpdated(originalExternalGroup, group.Value); isSame {
        return matchingGroup, nil
    }
    l := new(list.List)
    emptyGroup.GenerateChanges(originalExternalGroup, group.Value, nil, l)
    changes := make([]*dm.Change, l.Len())
    changeset := &dm.ChangeSet{
        CreatedAt: time.UTC().Format(dm.UTC_DATETIME_FORMAT),
        ChangedBy: group.ExternalServiceId,
        ChangeImportId: group.ExternalGroupId,
        RecordId: group.DsocialGroupId,
        Changes: changes,
    }
    _, err = ds.StoreContactChangeSet(changeset)
    if err != nil {
        return matchingGroup, nil
    }
    if originalExternalGroup == nil {
        group.Value, err = ds.StoreDsocialGroup(dsocialUserId, "", group.Value)
        if err != nil {
            return matchingGroup, err
        }
        storedDsocialGroup, err := ds.StoreDsocialGroupForExternalGroup(group.ExternalServiceId, group.ExternalUserId, group.ExternalGroupId, dsocialUserId, group.Value)
        return storedDsocialGroup, err
    }
    var storedDsocialGroup *dm.Group = nil
    if group.DsocialGroupId != "" {
        if storedDsocialGroup, _, err = ds.RetrieveDsocialGroup(dsocialUserId, group.DsocialGroupId); err != nil {
            return matchingGroup, err
        }
    }
    if storedDsocialGroup == nil {
        storedDsocialGroup = new(dm.Group)
    }
    for j, iter := 0, l.Front(); iter != nil; j, iter = j+1, iter.Next() {
        change := iter.Value.(*dm.Change)
        dm.ApplyChange(storedDsocialGroup, change)
        changes[j] = change
    }
    _, err = ds.StoreDsocialGroup(dsocialUserId, group.DsocialGroupId, storedDsocialGroup)
    return storedDsocialGroup, err
}

func (p *Pipeline) addContactToGroupMappings(m map[string]*list.List, contact *dm.Contact) {
    for _, groupRef := range contact.GroupReferences {
        if groupRef.Name != "" {
            var l *list.List
            var ok bool
            if l, ok = m[groupRef.Name]; !ok {
                l = list.New()
                m[groupRef.Name] = l
            }
            l.PushBack(&dm.ContactRef{
                Id: contact.Id,
                Name: contact.DisplayName,
            })
        }
    }
}

func (p *Pipeline) IncrementalSync(client oauth2_client.OAuth2Client, ds DataStoreService, cs ContactsService, dsocialUserId string) os.Error {
    groupMappings := make(map[string]*list.List)
    checkGroupsInContacts := cs.ContactInfoIncludesGroups()
    if cs.CanRetrieveContacts() {
        for contacts, nextToken, err := cs.RetrieveContacts(client, ds, dsocialUserId, nil); len(contacts) > 0 || nextToken != nil; contacts, nextToken, err = cs.RetrieveContacts(client, ds, dsocialUserId, nextToken) {
            if err != nil {
                return err
            }
            for _, contact := range contacts {
                finalContact, err := p.contactSync(ds, dsocialUserId, contact)
                if checkGroupsInContacts && finalContact != nil && finalContact.GroupReferences != nil && len(finalContact.GroupReferences) > 0 {
                    p.addContactToGroupMappings(groupMappings, finalContact)
                } 
                if err != nil {
                    return err
                }
            }
        }
    } else if cs.CanRetrieveConnections() {
        for connections, nextToken, err := cs.RetrieveConnections(client, ds, dsocialUserId, nil); len(connections) > 0 || nextToken != nil; connections, nextToken, err = cs.RetrieveConnections(client, ds, dsocialUserId, nextToken) {
            if err != nil {
                return err
            }
            for _, connection := range connections {
                contact, err := cs.RetrieveContact(client, ds, dsocialUserId, connection.ExternalContactId)
                if err != nil {
                    return err
                }
                finalContact, err := p.contactSync(ds, dsocialUserId, contact)
                if checkGroupsInContacts && finalContact != nil && finalContact != nil && finalContact.GroupReferences != nil && len(finalContact.GroupReferences) > 0 {
                    p.addContactToGroupMappings(groupMappings, finalContact)
                } 
                if err != nil {
                    return err
                }
            }
        }
        
    }
    if cs.CanRetrieveGroups() {
        for groups, nextToken, err := cs.RetrieveGroups(client, ds, dsocialUserId, nil); len(groups) > 0 || nextToken != nil; groups, nextToken, err = cs.RetrieveGroups(client, ds, dsocialUserId, nextToken) {
            if err != nil {
                return err
            }
            for _, group := range groups {
                _, err = p.groupSync(ds, dsocialUserId, group, groupMappings[group.Value.Name])
                if err != nil {
                    return err
                }
            }
        }
    } 
    return nil
}
