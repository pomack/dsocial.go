package contacts

import (
    "container/list"
    "fmt"
    dm "github.com/pomack/dsocial.go/models/dsocial"
    "github.com/pomack/oauth2_client.go/oauth2_client"
    "time"
)

type Pipeline struct {
}

func NewPipeline() *Pipeline {
    return new(Pipeline)
}

func (p *Pipeline) InitialSync(client oauth2_client.OAuth2Client, ds DataStoreService, cs ContactsService, csSettings ContactsServiceSettings, dsocialUserId, meContactId string) error {
    return p.Sync(client, ds, cs, csSettings, dsocialUserId, meContactId, true, false, false)
}

func (p *Pipeline) IncrementalSync(client oauth2_client.OAuth2Client, ds DataStoreService, cs ContactsService, csSettings ContactsServiceSettings, dsocialUserId, meContactId string) error {
    return p.Sync(client, ds, cs, csSettings, dsocialUserId, meContactId, true, true, true)
}

func (p *Pipeline) RemoveUnacceptedChanges(l *list.List, allowAdd, allowDelete, allowUpdate bool) *list.List {
    if allowAdd && allowDelete && allowUpdate {
        return l
    }
    n := list.New()
    for iter := l.Front(); iter != nil; iter = iter.Next() {
        if iter.Value == nil {
            continue
        }
        ch, _ := iter.Value.(*dm.Change)
        if ch == nil {
            continue
        }
        if !allowAdd && ch.ChangeType == dm.CHANGE_TYPE_ADD {
            fmt.Printf("[PIPELINE]: Removing add change for path of %#v\n", ch.Path)
        } else if !allowDelete && ch.ChangeType == dm.CHANGE_TYPE_DELETE {
            fmt.Printf("[PIPELINE]: Removing delete change for path of %#v\n", ch.Path)
        } else if !allowUpdate && ch.ChangeType == dm.CHANGE_TYPE_UPDATE {
            if !allowAdd || ch.Path == nil || len(ch.Path) <= 1 {
                fmt.Printf("[PIPELINE]: Removing update for path of %#v and allowAdd %v\n", ch.Path, allowAdd)
            } else {
                fmt.Printf("[PIPELINE]: Changing update to add for path of %#v and allowAdd %v\n", ch.Path, allowAdd)
                ch.ChangeType = dm.CHANGE_TYPE_ADD
                ch.OriginalValue = nil
                n.PushBack(ch)
            }
        } else {
            n.PushBack(ch)
        }
    }
    return n
}

func (p *Pipeline) findMatchingDsocialContact(ds DataStoreService, dsocialUserId string, contact *Contact) (extDsocialContact *dm.Contact, isSame bool, err error) {
    emptyContact := new(dm.Contact)
    if contact.DsocialContactId != "" {
        extDsocialContact, _, _ = ds.RetrieveDsocialContact(dsocialUserId, contact.DsocialContactId)
        if extDsocialContact != nil {
            _, isSame = emptyContact.IsSimilarOrUpdated(extDsocialContact, contact.Value)
            fmt.Printf("[PIPELINE]: findMatchingDsocialContact for %s with based on existing contact id will use %s and isSame %v\n", contact.Value.DisplayName, extDsocialContact.DisplayName, isSame)
        }
    }
    if extDsocialContact == nil {
        // this is a new contact from an existing service
        potentialMatches, err := ds.SearchForDsocialContacts(dsocialUserId, contact.Value)
        if err != nil {
            return nil, false, err
        }
        for _, potentialMatch := range potentialMatches {
            var isSimilar bool
            if isSimilar, isSame = emptyContact.IsSimilarOrUpdated(potentialMatch, contact.Value); isSimilar {
                extDsocialContact = potentialMatch
                break
            }
        }
        if extDsocialContact != nil {
            fmt.Printf("[PIPELINE]: findMatchingDsocialContact for %s was %s and isSame %v\n\tStoring mapping: %s/%s/%s -> %s\n", contact.Value.DisplayName, extDsocialContact.DisplayName, isSame, contact.ExternalServiceId, contact.ExternalUserId, contact.ExternalContactId, extDsocialContact.Id)
            _, _, err = ds.StoreDsocialExternalContactMapping(contact.ExternalServiceId, contact.ExternalUserId, contact.ExternalContactId, dsocialUserId, extDsocialContact.Id)
            contact.DsocialContactId = extDsocialContact.Id
        } else {
            fmt.Printf("[PIPELINE]: findMatchingDsocialContact cannot find similar for %s\n", contact.Value.DisplayName)
        }
    }
    return extDsocialContact, isSame, err
}

func (p *Pipeline) contactImport(cs ContactsService, ds DataStoreService, dsocialUserId string, contact *Contact, allowAdd, allowDelete, allowUpdate bool) (*dm.Contact, string, error) {
    emptyContact := new(dm.Contact)
    if contact == nil || contact.Value == nil {
        return nil, "", nil
    }
    fmt.Printf("[PIPELINE]: Importing contact with ExternalServiceId = %v, ExternalUserId = %v, ExternalContactId = %v, DsocialUserId = %v, DsocialContactId = %v\n", contact.ExternalServiceId, contact.ExternalUserId, contact.ExternalContactId, contact.DsocialUserId, contact.DsocialContactId)
    extDsocialContact, _, err := ds.RetrieveDsocialContactForExternalContact(contact.ExternalServiceId, contact.ExternalUserId, contact.ExternalContactId, dsocialUserId)
    if err != nil {
        return nil, "", err
    }
    var matchingContact *dm.Contact
    var isSame bool
    if extDsocialContact == nil {
        // We don't have a mapping for this external contact to an internal contact mapping
        // meaning we've never imported this contact previously from THIS service, but we may
        // already have the contact in our system, so let's see if we can find it
        matchingContact, isSame, err = p.findMatchingDsocialContact(ds, dsocialUserId, contact)
        if err != nil {
            return matchingContact, "", err
        }
        if matchingContact != nil {
            contact.DsocialContactId = matchingContact.Id
            extContact := cs.ConvertToExternalContact(matchingContact, nil, dsocialUserId)
            ds.StoreExternalContact(contact.ExternalServiceId, contact.ExternalUserId, dsocialUserId, contact.ExternalContactId, extContact)
            extDsocialContact = cs.ConvertToDsocialContact(extContact, matchingContact, dsocialUserId)
            if extDsocialContact != nil {
                AddIdsForDsocialContact(extDsocialContact, ds, dsocialUserId)
                //contact.ExternalContactId = extDsocialContact.Id
                extDsocialContact, err = ds.StoreDsocialContactForExternalContact(contact.ExternalServiceId, contact.ExternalUserId, contact.ExternalContactId, dsocialUserId, extDsocialContact)
                if extDsocialContact == nil || err != nil {
                    return matchingContact, "", err
                }
                if contact.DsocialContactId != "" {
                    if _, _, err = ds.StoreDsocialExternalContactMapping(contact.ExternalServiceId, contact.ExternalUserId, contact.ExternalContactId, dsocialUserId, contact.DsocialContactId); err != nil {
                        return matchingContact, "", err
                    }
                }
            }
        }
    } else {
        // we have a mapping for this external contact to an internal contact mapping
        // from THIS service, therefore let's use it
        if contact.DsocialContactId == "" {
            contact.DsocialContactId, err = ds.DsocialIdForExternalContactId(contact.ExternalServiceId, contact.ExternalUserId, dsocialUserId, contact.ExternalContactId)
            if err != nil {
                return nil, "", err
            }
        }
        if contact.DsocialContactId != "" {
            matchingContact, _, err = ds.RetrieveDsocialContact(dsocialUserId, contact.DsocialContactId)
            if err != nil {
                return nil, "", err
            }
        }
    }
    // ensure we have a contact Id
    if contact.DsocialContactId == "" {
        if matchingContact != nil {
            contact.DsocialContactId = matchingContact.Id
            fmt.Printf("[PIPELINE]: Will be using matchingContact Id: %v\n", matchingContact.Id)
        }
        if contact.DsocialContactId == "" {
            newContact := &dm.Contact{UserId: dsocialUserId}
            AddIdsForDsocialContact(newContact, ds, dsocialUserId)
            thecontact, err := ds.StoreDsocialContact(dsocialUserId, newContact.Id, newContact)
            if err != nil {
                return nil, "", err
            }
            contact.DsocialContactId = thecontact.Id
        }
    }
    if _, isSame = emptyContact.IsSimilarOrUpdated(extDsocialContact, contact.Value); isSame {
        return matchingContact, "", nil
    }
    l := new(list.List)
    emptyContact.GenerateChanges(extDsocialContact, contact.Value, nil, l)
    l = p.RemoveUnacceptedChanges(l, allowAdd, allowDelete, allowUpdate)
    changes := make([]*dm.Change, l.Len())
    for i, iter := 0, l.Front(); iter != nil; i, iter = i+1, iter.Next() {
        changes[i] = iter.Value.(*dm.Change)
    }
    changeset := &dm.ChangeSet{
        CreatedAt:      time.Now().UTC().Format(dm.UTC_DATETIME_FORMAT),
        ChangedBy:      contact.ExternalServiceId,
        ChangeImportId: contact.ExternalContactId,
        RecordId:       contact.DsocialContactId,
        Changes:        changes,
    }
    _, err = ds.StoreContactChangeSet(dsocialUserId, changeset)
    if err != nil {
        return matchingContact, changeset.Id, err
    }
    if extDsocialContact == nil {
        fmt.Printf("[PIPELINE]: OriginalExternalContact is nil and contact.DsocialContactId is %v and contact.Value.Id was %v\n", contact.DsocialContactId, contact.Value.Id)
        contact.Value.Id = contact.DsocialContactId
        AddIdsForDsocialContact(contact.Value, ds, dsocialUserId)
        contact.Value, err = ds.StoreDsocialContact(dsocialUserId, contact.DsocialContactId, contact.Value)
        fmt.Printf("[PIPELINE]: After storing contact.Value, contact.Value.Id is %v\n", contact.Value.Id)
        if err != nil {
            return matchingContact, changeset.Id, err
        }
        storedDsocialContact, err := ds.StoreDsocialContactForExternalContact(contact.ExternalServiceId, contact.ExternalUserId, contact.ExternalContactId, dsocialUserId, contact.Value)
        fmt.Printf("[PIPELINE]: After storing external contact, contact.Value.Id is %v\n", contact.Value.Id)
        _, _, err2 := ds.StoreDsocialExternalContactMapping(contact.ExternalServiceId, contact.ExternalUserId, contact.ExternalContactId, dsocialUserId, contact.DsocialContactId)
        if err == nil {
            err = err2
        }
        return storedDsocialContact, changeset.Id, err
    }
    var storedDsocialContact *dm.Contact = nil
    if contact.DsocialContactId != "" {
        if storedDsocialContact, _, err = ds.RetrieveDsocialContact(dsocialUserId, contact.DsocialContactId); err != nil {
            return matchingContact, changeset.Id, err
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
    AddIdsForDsocialContact(storedDsocialContact, ds, dsocialUserId)
    _, err = ds.StoreDsocialContact(dsocialUserId, contact.DsocialContactId, storedDsocialContact)
    return storedDsocialContact, changeset.Id, err
}

func (p *Pipeline) findMatchingDsocialGroup(ds DataStoreService, dsocialUserId string, group *Group) (extDsocialGroup *dm.Group, isSame bool, err error) {
    emptyGroup := new(dm.Group)
    if group.DsocialGroupId != "" {
        extDsocialGroup, _, _ = ds.RetrieveDsocialGroup(dsocialUserId, group.DsocialGroupId)
        if extDsocialGroup != nil {
            _, isSame = emptyGroup.IsSimilarOrUpdated(extDsocialGroup, group.Value)
            fmt.Printf("[PIPELINE]: findMatchingDsocialGroup for %s with based on existing group id will use %s and isSame %v\n", group.Value.Name, extDsocialGroup.Name, isSame)
        }
    }
    if extDsocialGroup == nil {
        // this is a new group from an existing service
        potentialMatches, err := ds.SearchForDsocialGroups(dsocialUserId, group.Value.Name)
        if err != nil {
            return nil, false, err
        }
        for _, potentialMatch := range potentialMatches {
            var isSimilar bool
            if isSimilar, isSame = emptyGroup.IsSimilarOrUpdated(potentialMatch, group.Value); isSimilar {
                extDsocialGroup = potentialMatch
                break
            }
        }
        if extDsocialGroup != nil {
            fmt.Printf("[PIPELINE]: findMatchingDsocialGroup for %s was %s and isSame %v\n\tStoring mapping: %s/%s/%s -> %s\n", group.Value.Name, extDsocialGroup.Name, isSame, group.ExternalServiceId, group.ExternalUserId, group.ExternalGroupId, extDsocialGroup.Id)
            _, _, err = ds.StoreDsocialExternalGroupMapping(group.ExternalServiceId, group.ExternalUserId, group.ExternalGroupId, dsocialUserId, extDsocialGroup.Id)
            group.DsocialGroupId = extDsocialGroup.Id
        } else {
            fmt.Printf("[PIPELINE]: findMatchingDsocialGroup cannot find similar for %s\n", group.Value.Name)
        }
    }
    return extDsocialGroup, isSame, err
}

func (p *Pipeline) groupImport(cs ContactsService, ds DataStoreService, dsocialUserId string, group *Group, minimumIncludes *list.List, allowAdd, allowDelete, allowUpdate bool) (*dm.Group, string, error) {
    emptyGroup := new(dm.Group)
    if group == nil || group.Value == nil {
        return nil, "", nil
    }
    //fmt.Printf("[PIPELINE]: Syncing group: %s\n", group.Value.Name)
    if group.Value.ContactIds == nil {
        group.Value.ContactIds = make([]string, 0, 10)
    }
    if group.Value.ContactNames == nil {
        group.Value.ContactNames = make([]string, 0, 10)
    }
    if len(group.Value.ContactIds) == 0 && len(group.Value.ContactNames) == 0 && minimumIncludes != nil {
        sv1 := make([]string, len(group.Value.ContactIds), len(group.Value.ContactIds)+minimumIncludes.Len())
        sv2 := make([]string, len(group.Value.ContactNames), len(group.Value.ContactNames)+minimumIncludes.Len())
        copy(sv1, group.Value.ContactIds)
        copy(sv2, group.Value.ContactNames)
        for iter := minimumIncludes.Front(); iter != nil; iter = iter.Next() {
            contactRef := iter.Value.(*dm.ContactRef)
            sv1 = append(sv1, contactRef.Id)
            sv2 = append(sv2, contactRef.Name)
        }
    } else if minimumIncludes != nil {
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
                sv1 := make([]string, len(group.Value.ContactIds), len(group.Value.ContactIds)+1)
                sv2 := make([]string, len(group.Value.ContactNames), len(group.Value.ContactNames)+1)
                copy(sv1, group.Value.ContactIds)
                copy(sv2, group.Value.ContactNames)
                sv1 = append(sv1, contactRef.Id)
                sv2 = append(sv2, contactRef.Name)
            } else {
                group.Value.ContactIds[refLocation] = contactRef.Id
                group.Value.ContactNames[refLocation] = contactRef.Name
            }
        }
    }
    extDsocialGroup, _, err := ds.RetrieveDsocialGroupForExternalGroup(group.ExternalServiceId, group.ExternalUserId, group.ExternalGroupId, dsocialUserId)
    if err != nil {
        return nil, "", err
    }
    var matchingGroup *dm.Group
    var isSame bool
    if extDsocialGroup == nil {
        // We don't have a mapping for this external group to an internal group mapping
        // meaning we've never imported this group previously from THIS service, but we may
        // already have the group in our system, so let's see if we can find it
        matchingGroup, isSame, err = p.findMatchingDsocialGroup(ds, dsocialUserId, group)
        if err != nil {
            return matchingGroup, "", err
        }
        if matchingGroup != nil {
            group.DsocialGroupId = matchingGroup.Id
            extGroup := cs.ConvertToExternalGroup(matchingGroup, nil, dsocialUserId)
            ds.StoreExternalGroup(group.ExternalServiceId, group.ExternalUserId, dsocialUserId, group.ExternalGroupId, extGroup)
            extDsocialGroup = cs.ConvertToDsocialGroup(extGroup, matchingGroup, dsocialUserId)
            if extDsocialGroup != nil {
                AddIdsForDsocialGroup(extDsocialGroup, ds, dsocialUserId)
                fmt.Printf("[PIPELINE]: groupImport() before store dsoc group ExternalGroupId: %v and extDsocialGroup.Id %v matchingGroup.Id %v\n", group.ExternalGroupId, extDsocialGroup.Id, matchingGroup.Id)
                //group.ExternalGroupId = extDsocialGroup.Id
                //extDsocialGroup.Id = group.DsocialGroupId
                extDsocialGroup, err = ds.StoreDsocialGroupForExternalGroup(group.ExternalServiceId, group.ExternalUserId, group.ExternalGroupId, dsocialUserId, extDsocialGroup)
                if extDsocialGroup == nil || err != nil {
                    return matchingGroup, "", err
                }
                //extDsocialGroup.Id = group.DsocialGroupId
                fmt.Printf("[PIPELINE]: groupImport() before store mapping ExternalGroupId: %v and DsocialGroupId %v\n", group.ExternalGroupId, group.DsocialGroupId)
                if _, _, err = ds.StoreDsocialExternalGroupMapping(group.ExternalServiceId, group.ExternalUserId, group.ExternalGroupId, dsocialUserId, group.DsocialGroupId); err != nil {
                    return matchingGroup, "", err
                }
            }
        }
    } else {
        // we have a mapping for this external group to an internal group mapping
        // from THIS service, therefore let's use it
        if group.DsocialGroupId == "" {
            group.DsocialGroupId, err = ds.DsocialIdForExternalGroupId(group.ExternalServiceId, group.ExternalUserId, dsocialUserId, group.ExternalGroupId)
            if err != nil {
                return nil, "", err
            }
        }
        if group.DsocialGroupId != "" {
            matchingGroup, _, err = ds.RetrieveDsocialGroup(dsocialUserId, group.DsocialGroupId)
            if err != nil {
                return nil, "", err
            }
        }
    }
    // ensure we have a contact Id
    if group.DsocialGroupId == "" {
        if matchingGroup != nil {
            group.DsocialGroupId = matchingGroup.Id
            fmt.Printf("[PIPELINE]: Will be using matchingGroup Id: %v\n", matchingGroup.Id)
        }
        if group.DsocialGroupId == "" {
            newGroup := &dm.Group{UserId: dsocialUserId}
            AddIdsForDsocialGroup(newGroup, ds, dsocialUserId)
            thegroup, err := ds.StoreDsocialGroup(dsocialUserId, newGroup.Id, newGroup)
            if err != nil {
                return nil, "", err
            }
            group.DsocialGroupId = thegroup.Id
        }
    }
    if _, isSame = emptyGroup.IsSimilarOrUpdated(extDsocialGroup, group.Value); isSame {
        return matchingGroup, "", nil
    }
    l := new(list.List)
    emptyGroup.GenerateChanges(extDsocialGroup, group.Value, nil, l)
    l = p.RemoveUnacceptedChanges(l, allowAdd, allowDelete, allowUpdate)
    changes := make([]*dm.Change, l.Len())
    for i, iter := 0, l.Front(); iter != nil; i, iter = i+1, iter.Next() {
        changes[i] = iter.Value.(*dm.Change)
    }
    changeset := &dm.ChangeSet{
        CreatedAt:      time.Now().UTC().Format(dm.UTC_DATETIME_FORMAT),
        ChangedBy:      group.ExternalServiceId,
        ChangeImportId: group.ExternalGroupId,
        RecordId:       group.DsocialGroupId,
        Changes:        changes,
    }
    _, err = ds.StoreGroupChangeSet(dsocialUserId, changeset)
    if err != nil {
        return matchingGroup, changeset.Id, nil
    }
    if extDsocialGroup == nil {
        fmt.Printf("[PIPELINE]: OriginalExternalGroup is nil and group.DsocialGroupId is %v and group.Value.Id was %v\n", group.DsocialGroupId, group.Value.Id)
        group.Value.Id = group.DsocialGroupId
        AddIdsForDsocialGroup(group.Value, ds, dsocialUserId)
        group.Value, err = ds.StoreDsocialGroup(dsocialUserId, group.DsocialGroupId, group.Value)
        if err != nil {
            return matchingGroup, changeset.Id, err
        }
        storedDsocialGroup, err := ds.StoreDsocialGroupForExternalGroup(group.ExternalServiceId, group.ExternalUserId, group.ExternalGroupId, dsocialUserId, group.Value)
        _, _, err2 := ds.StoreDsocialExternalGroupMapping(group.ExternalServiceId, group.ExternalUserId, group.ExternalGroupId, dsocialUserId, group.DsocialGroupId)
        if err == nil {
            err = err2
        }
        return storedDsocialGroup, changeset.Id, err
    }
    var storedDsocialGroup *dm.Group = nil
    if group.DsocialGroupId != "" {
        if storedDsocialGroup, _, err = ds.RetrieveDsocialGroup(dsocialUserId, group.DsocialGroupId); err != nil {
            return matchingGroup, changeset.Id, err
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
    AddIdsForDsocialGroup(storedDsocialGroup, ds, dsocialUserId)
    _, err = ds.StoreDsocialGroup(dsocialUserId, group.DsocialGroupId, storedDsocialGroup)
    return storedDsocialGroup, changeset.Id, err
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
                Id:   contact.Id,
                Name: contact.DisplayName,
            })
        }
    }
}

func (p *Pipeline) Sync(client oauth2_client.OAuth2Client, ds DataStoreService, cs ContactsService, csSettings ContactsServiceSettings, dsocialUserId, meContactId string, allowAdd, allowDelete, allowUpdate bool) (err error) {
    err = p.Import(client, ds, cs, csSettings, dsocialUserId, meContactId, allowAdd, allowDelete, allowUpdate)
    if err != nil {
        return err
    }
    err = p.Export(client, ds, cs, csSettings, dsocialUserId, meContactId)
    return
}

func (p *Pipeline) importContacts(client oauth2_client.OAuth2Client, ds DataStoreService, cs ContactsService, dsocialUserId string, allowAdd, allowDelete, allowUpdate bool, groupMappings map[string]*list.List, contactChangesetIds []string) (err error) {
    checkGroupsInContacts := cs.ContactInfoIncludesGroups()
    var nextToken NextToken = "blah"
    for contacts, useNextToken, err := cs.RetrieveContacts(client, ds, dsocialUserId, nil); (len(contacts) > 0 && nextToken != nil) || err != nil; contacts, useNextToken, err = cs.RetrieveContacts(client, ds, dsocialUserId, nextToken) {
        if err != nil {
            break
        }
        for _, contact := range contacts {
            finalContact, changesetId, err := p.contactImport(cs, ds, dsocialUserId, contact, allowAdd, allowDelete, allowUpdate)
            if changesetId != "" {
                contactChangesetIds = append(contactChangesetIds, changesetId)
            }
            if checkGroupsInContacts && finalContact != nil && finalContact.GroupReferences != nil && len(finalContact.GroupReferences) > 0 {
                p.addContactToGroupMappings(groupMappings, finalContact)
            }
            if err != nil {
                break
            }
        }
        nextToken = useNextToken
        if err != nil {
            break
        }
    }
    return
}

func (p *Pipeline) importConnections(client oauth2_client.OAuth2Client, ds DataStoreService, cs ContactsService, dsocialUserId string, allowAdd, allowDelete, allowUpdate bool, groupMappings map[string]*list.List, contactChangesetIds []string) (err error) {
    checkGroupsInContacts := cs.ContactInfoIncludesGroups()
    var nextToken NextToken = "blah"
    for connections, useNextToken, err := cs.RetrieveConnections(client, ds, dsocialUserId, nil); (len(connections) > 0 && nextToken != nil) || err != nil; connections, useNextToken, err = cs.RetrieveConnections(client, ds, dsocialUserId, nextToken) {
        if err != nil {
            break
        }
        for _, connection := range connections {
            contact, err := cs.RetrieveContact(client, ds, dsocialUserId, connection.ExternalContactId)
            if err != nil {
                break
            }
            finalContact, changesetId, err := p.contactImport(cs, ds, dsocialUserId, contact, allowAdd, allowDelete, allowUpdate)
            if changesetId != "" {
                contactChangesetIds = append(contactChangesetIds, changesetId)
            }
            if checkGroupsInContacts && finalContact != nil && finalContact != nil && finalContact.GroupReferences != nil && len(finalContact.GroupReferences) > 0 {
                p.addContactToGroupMappings(groupMappings, finalContact)
            }
            if err != nil {
                break
            }
        }
        nextToken = useNextToken
        if err != nil {
            break
        }
    }
    return
}

func (p *Pipeline) importGroups(client oauth2_client.OAuth2Client, ds DataStoreService, cs ContactsService, dsocialUserId string, allowAdd, allowDelete, allowUpdate bool, groupMappings map[string]*list.List, groupChangesetIds []string) (err error) {
    var nextToken NextToken = "blah"
    for groups, useNextToken, err := cs.RetrieveGroups(client, ds, dsocialUserId, nil); (len(groups) > 0 && nextToken != nil) || err != nil; groups, useNextToken, err = cs.RetrieveGroups(client, ds, dsocialUserId, nextToken) {
        if err != nil {
            break
        }
        for _, group := range groups {
            var changesetId string
            _, changesetId, err = p.groupImport(cs, ds, dsocialUserId, group, groupMappings[group.Value.Name], allowAdd, allowDelete, allowUpdate)
            if changesetId != "" {
                groupChangesetIds = append(groupChangesetIds, changesetId)
            }
            if err != nil {
                break
            }
        }
        nextToken = useNextToken
        if err != nil {
            break
        }
    }
    return
}

func (p *Pipeline) queueContactChangeSetsToApply(ds DataStoreService, cs ContactsService, csSettings ContactsServiceSettings, settings []ContactsServiceSettings, dsocialUserId string, changesetIds []string) (err error) {
    if changesetIds == nil || len(changesetIds) == 0 || settings == nil || len(settings) == 0 {
        return
    }
    thisServiceName := cs.ServiceId()
    thisServiceId := csSettings.Id()
    ids := []string(changesetIds)
    for _, setting := range settings {
        if setting.Id() == thisServiceId && thisServiceName == setting.ContactsServiceId() {
            // if we import from this service, don't export to it
            continue
        }
        if _, err2 := ds.AddContactChangeSetsToApply(dsocialUserId, setting.Id(), setting.ContactsServiceId(), ids); err2 != nil {
            if err == nil {
                err = err2
            }
            break
        }
    }
    return
}

func (p *Pipeline) queueGroupChangeSetsToApply(ds DataStoreService, cs ContactsService, csSettings ContactsServiceSettings, settings []ContactsServiceSettings, dsocialUserId string, changesetIds []string) (err error) {
    if changesetIds == nil || len(changesetIds) == 0 || settings == nil || len(settings) == 0 {
        return
    }
    thisServiceName := cs.ServiceId()
    thisServiceId := csSettings.Id()
    ids := []string(changesetIds)
    for _, setting := range settings {
        if setting.Id() == thisServiceId && thisServiceName == setting.ContactsServiceId() {
            // if we import from this service, don't export to it
            continue
        }
        if _, err2 := ds.AddGroupChangeSetsToApply(dsocialUserId, setting.Id(), setting.ContactsServiceId(), ids); err2 != nil {
            if err == nil {
                err = err2
            }
            break
        }
    }
    return
}

func (p *Pipeline) Import(client oauth2_client.OAuth2Client, ds DataStoreService, cs ContactsService, csSettings ContactsServiceSettings, dsocialUserId, meContactId string, allowAdd, allowDelete, allowUpdate bool) (err error) {
    fmt.Printf("[PIPELINE]: Starting importing...\n")
    if !csSettings.AllowRetrieveContactInfo() || !cs.CanImportContactsOrGroups() {
        // do nothing
    } else {
        groupMappings := make(map[string]*list.List)
        contactChangesetIds := make([]string, 0, 8)
        groupChangesetIds := make([]string, 0, 8)
        if cs.CanRetrieveContacts() {
            err = p.importContacts(client, ds, cs, dsocialUserId, allowAdd, allowDelete, allowUpdate, groupMappings, contactChangesetIds)
        } else if cs.CanRetrieveConnections() {
            err = p.importConnections(client, ds, cs, dsocialUserId, allowAdd, allowDelete, allowUpdate, groupMappings, contactChangesetIds)
        }
        if err == nil && cs.CanRetrieveGroups() {
            err = p.importGroups(client, ds, cs, dsocialUserId, allowAdd, allowDelete, allowUpdate, groupMappings, groupChangesetIds)
        }
        if len(contactChangesetIds) > 0 || len(groupChangesetIds) > 0 {
            settings, _ := ds.RetrieveAllContactsServiceSettingsForUser(dsocialUserId)
            err = p.queueContactChangeSetsToApply(ds, cs, csSettings, settings, dsocialUserId, contactChangesetIds)
            if err == nil {
                err = p.queueGroupChangeSetsToApply(ds, cs, csSettings, settings, dsocialUserId, groupChangesetIds)
            }
        }
    }
    fmt.Printf("[PIPELINE]: Done importing.\n")
    return
}

func (p *Pipeline) extractAllChangeSetIds(applyable []*dm.ChangeSetsToApply, changesets map[string]*dm.ChangeSet) []string {
    arr := make([]string, 0, len(changesets))
    for _, toApply := range applyable {
        for _, changesetId := range toApply.ChangeSetIds {
            changeset, _ := changesets[changesetId]
            if changeset != nil {
                arr = append(arr, changeset.Id)
            }
        }
    }
    return arr
}

func (p *Pipeline) markAllContactChangeSetsNotApplyable(ds DataStoreService, dsocialUserId, externalServiceId, externalServiceName string) (err error) {
    applyable, changesets, err := ds.RetrieveContactChangeSetsToApply(dsocialUserId, externalServiceId, externalServiceName)
    if err != nil {
        return
    }
    changesetIdsNotApplyable := p.extractAllChangeSetIds(applyable, changesets)
    err = p.markContactChangeSetsNotApplyable(ds, dsocialUserId, externalServiceId, externalServiceName, changesetIdsNotApplyable)
    return
}

func (p *Pipeline) markAllGroupChangeSetsNotApplyable(ds DataStoreService, dsocialUserId, externalServiceId, externalServiceName string) (err error) {
    applyable, changesets, err := ds.RetrieveGroupChangeSetsToApply(dsocialUserId, externalServiceId, externalServiceName)
    if err != nil {
        return
    }
    changesetIdsNotApplyable := p.extractAllChangeSetIds(applyable, changesets)
    err = p.markGroupChangeSetsNotApplyable(ds, dsocialUserId, externalServiceId, externalServiceName, changesetIdsNotApplyable)
    return
}

func (p *Pipeline) markContactChangeSetsNotApplyable(ds DataStoreService, dsocialUserId, externalServiceId, externalServiceName string, changesetIdsNotApplyable []string) (err error) {
    if changesetIdsNotApplyable == nil || len(changesetIdsNotApplyable) == 0 {
        return
    }
    if _, err = ds.AddContactChangeSetsNotCurrentlyApplyable(dsocialUserId, externalServiceId, externalServiceName, changesetIdsNotApplyable); err != nil {
        return
    }
    err = ds.RemoveContactChangeSetsToApply(dsocialUserId, externalServiceId, externalServiceName, changesetIdsNotApplyable)
    return
}

func (p *Pipeline) markGroupChangeSetsNotApplyable(ds DataStoreService, dsocialUserId, externalServiceId, externalServiceName string, changesetIdsNotApplyable []string) (err error) {
    if changesetIdsNotApplyable == nil || len(changesetIdsNotApplyable) == 0 {
        return
    }
    if _, err = ds.AddGroupChangeSetsNotCurrentlyApplyable(dsocialUserId, externalServiceId, externalServiceName, changesetIdsNotApplyable); err != nil {
        return
    }
    err = ds.RemoveGroupChangeSetsToApply(dsocialUserId, externalServiceId, externalServiceName, changesetIdsNotApplyable)
    return
}

func (p *Pipeline) handleDeleteContact(client oauth2_client.OAuth2Client, ds DataStoreService, cs ContactsService, dsocialUserId, externalServiceId, externalUserId, dsocialContactId string) (err error) {
    fmt.Printf("[PIPELINE]: Handling Delete Contact...\n")
    externalContactId, err := ds.ExternalContactIdForDsocialId(externalServiceId, externalUserId, dsocialUserId, dsocialContactId)
    if err != nil {
        return
    }
    if externalContactId == "" {
        // nothing to delete
        return
    }
    _, err = DeleteContactOnExternalService(client, cs, ds, dsocialUserId, dsocialContactId)
    return
}

func (p *Pipeline) handleCreateContact(client oauth2_client.OAuth2Client, ds DataStoreService, cs ContactsService, dsocialUserId, dsocialContactId string) (err error) {
    fmt.Printf("[PIPELINE]: Handling Create Contact...\n")
    // don't have it locally
    dsocContact, _, err := ds.RetrieveDsocialContact(dsocialUserId, dsocialContactId)
    if err != nil {
        return
    }
    if dsocContact == nil {
        // can't create what doesn't exist anymore...ignore
        return
    }
    _, err = CreateContactOnExternalService(client, cs, ds, dsocialUserId, dsocContact)
    return
}

func (p *Pipeline) handleUpdateContact(client oauth2_client.OAuth2Client, ds DataStoreService, cs ContactsService, dsocialUserId, externalServiceId, externalUserId, dsocialContactId string, changeset *dm.ChangeSet) error {
    fmt.Printf("[PIPELINE]: Handling Update Contact...\n")
    externalContactId, err := ds.ExternalContactIdForDsocialId(externalServiceId, externalUserId, dsocialUserId, dsocialContactId)
    if err != nil {
        return err
    }
    if externalContactId != "" {
        dsocExternalContact, _, err := ds.RetrieveDsocialContactForExternalContact(externalServiceId, externalUserId, externalContactId, dsocialUserId)
        if err != nil {
            return err
        }
        if dsocExternalContact != nil {
            for _, change := range changeset.Changes {
                dm.ApplyChange(dsocExternalContact, change)
            }
        }
        origDsocExternalContact, _, err := ds.RetrieveDsocialContactForExternalContact(externalServiceId, externalUserId, externalContactId, dsocialUserId)
        if err != nil {
            return err
        }
        _, err = UpdateContactOnExternalService(client, cs, ds, dsocialUserId, origDsocExternalContact, dsocExternalContact)
    } else {
        err = p.handleCreateContact(client, ds, cs, dsocialUserId, dsocialContactId)
    }
    return err
}

func (p *Pipeline) handleDeleteGroup(client oauth2_client.OAuth2Client, ds DataStoreService, cs ContactsService, dsocialUserId, externalServiceId, externalUserId, dsocialGroupId string) (err error) {
    fmt.Printf("[PIPELINE]: Handling Delete Group...\n")
    externalGroupId, err := ds.ExternalGroupIdForDsocialId(externalServiceId, externalUserId, dsocialUserId, dsocialGroupId)
    if err != nil {
        return
    }
    if externalGroupId == "" {
        // nothing to delete
        return
    }
    _, err = DeleteGroupOnExternalService(client, cs, ds, dsocialUserId, dsocialGroupId)
    return
}

func (p *Pipeline) handleCreateGroup(client oauth2_client.OAuth2Client, ds DataStoreService, cs ContactsService, dsocialUserId, dsocialGroupId string) (err error) {
    fmt.Printf("[PIPELINE]: Handling Create Group...\n")
    // don't have it locally
    dsocGroup, _, err := ds.RetrieveDsocialGroup(dsocialUserId, dsocialGroupId)
    if err != nil {
        return
    }
    if dsocGroup == nil {
        // can't create what doesn't exist anymore...ignore
        return
    }
    _, err = CreateGroupOnExternalService(client, cs, ds, dsocialUserId, dsocGroup)
    return
}

func (p *Pipeline) handleUpdateGroup(client oauth2_client.OAuth2Client, ds DataStoreService, cs ContactsService, dsocialUserId, externalServiceId, externalUserId, dsocialGroupId string, changeset *dm.ChangeSet) error {
    fmt.Printf("[PIPELINE]: Handling Update Group...\n")
    externalGroupId, err := ds.ExternalGroupIdForDsocialId(externalServiceId, externalUserId, dsocialUserId, dsocialGroupId)
    if err != nil {
        return err
    }
    if externalGroupId != "" {
        dsocExternalGroup, _, err := ds.RetrieveDsocialGroupForExternalGroup(externalServiceId, externalUserId, externalGroupId, dsocialUserId)
        if err != nil {
            return err
        }
        if dsocExternalGroup != nil {
            for _, change := range changeset.Changes {
                dm.ApplyChange(dsocExternalGroup, change)
            }
        }
        origDsocExternalGroup, _, err := ds.RetrieveDsocialGroupForExternalGroup(externalServiceId, externalUserId, externalGroupId, dsocialUserId)
        if err != nil {
            return err
        }
        _, err = UpdateGroupOnExternalService(client, cs, ds, dsocialUserId, origDsocExternalGroup, dsocExternalGroup)
    } else {
        err = p.handleCreateGroup(client, ds, cs, dsocialUserId, dsocialGroupId)
    }
    return err
}

func (p *Pipeline) applyContactChangeSets(client oauth2_client.OAuth2Client, ds DataStoreService, cs ContactsService, csSettings ContactsServiceSettings, dsocialUserId, meContactId string) (err error) {
    externalServiceId := csSettings.Id()
    externalServiceName := csSettings.ContactsServiceId()
    externalUserId := csSettings.ExternalUserId()
    applyable, changesets, err := ds.RetrieveContactChangeSetsToApply(dsocialUserId, externalServiceId, externalServiceName)
    if err != nil {
        return
    }
    fmt.Printf("[PIPELINE]: Will be applying up to %d changesets from %d applyable for %s, %s, %s\n", len(changesets), len(applyable), externalServiceId, externalServiceName, externalUserId)
    changesetIdsNotApplyable := make([]string, 0, 8)
    changesetIdsApplied := make([]string, 0, 8)
    for _, toApply := range applyable {
        if err != nil {
            break
        }
        for _, changesetId := range toApply.ChangeSetIds {
            if err != nil {
                break
            }
            changeset, _ := changesets[changesetId]
            if changeset != nil {
                isMe := changeset.RecordId == meContactId
                canCreate, canUpdate, canDelete := cs.CanCreateContact(isMe), cs.CanUpdateContact(isMe), cs.CanDeleteContact(isMe)
                if !canCreate && !canUpdate && !canDelete {
                    changesetIdsNotApplyable = append(changesetIdsNotApplyable, changeset.Id)
                    continue
                }
                fmt.Printf("[PIPELINE]: handling changeset id %s\n", changeset.Id)
                isCreate, isUpdate, isDelete := false, false, false
                if len(changeset.Changes) > 1 {
                    isUpdate = true
                } else if len(changeset.Changes) == 1 {
                    switch changeset.Changes[0].ChangeType {
                    case dm.CHANGE_TYPE_CREATE:
                        isCreate = true
                    case dm.CHANGE_TYPE_ADD:
                        if changeset.Changes[0].Path == nil || len(changeset.Changes[0].Path) == 0 {
                            isCreate = true
                        } else {
                            isUpdate = true
                        }
                    case dm.CHANGE_TYPE_UPDATE:
                        isUpdate = true
                    case dm.CHANGE_TYPE_DELETE:
                        if changeset.Changes[0].Path == nil || len(changeset.Changes[0].Path) == 0 {
                            isDelete = true
                        } else {
                            isUpdate = true
                        }
                    }
                }
                if !(isCreate && canCreate) && !(isUpdate && canUpdate) && !(isDelete && canDelete) {
                    changesetIdsNotApplyable = append(changesetIdsNotApplyable, changeset.Id)
                    continue
                }
                dsocialContactId := changeset.RecordId
                if isDelete {
                    err = p.handleDeleteContact(client, ds, cs, dsocialUserId, externalServiceId, externalUserId, dsocialContactId)
                } else if isCreate {
                    err = p.handleCreateContact(client, ds, cs, dsocialUserId, dsocialContactId)
                } else {
                    // must be update
                    err = p.handleUpdateContact(client, ds, cs, dsocialUserId, externalServiceId, externalUserId, dsocialContactId, changeset)
                }
                changesetIdsApplied = append(changesetIdsApplied, changeset.Id)
            } else {
                fmt.Printf("[PIPELINE]: Unable to apply nil changeset\n")
            }
        }
    }
    if err == nil {
        fmt.Printf("[PIPELINE]: Skipped applying %d changesets\n", len(changesetIdsNotApplyable))
        err = p.markContactChangeSetsNotApplyable(ds, dsocialUserId, externalServiceId, externalServiceName, changesetIdsNotApplyable)
        if err == nil {
            err = ds.RemoveContactChangeSetsToApply(dsocialUserId, externalServiceId, externalServiceName, changesetIdsApplied)
        }
    }
    fmt.Printf("[PIPELINE]: Done applying contact changesets\n")
    return
}

func (p *Pipeline) applyGroupChangeSets(client oauth2_client.OAuth2Client, ds DataStoreService, cs ContactsService, csSettings ContactsServiceSettings, dsocialUserId, meContactId string) (err error) {
    externalServiceId := csSettings.Id()
    externalServiceName := csSettings.ContactsServiceId()
    externalUserId := csSettings.ExternalUserId()
    applyable, changesets, err := ds.RetrieveGroupChangeSetsToApply(dsocialUserId, externalServiceId, externalServiceName)
    if err != nil {
        return
    }
    fmt.Printf("[PIPELINE]: Will be applying up to %d changesets from %d applyable\n", len(changesets), len(applyable))
    changesetIdsNotApplyable := make([]string, 0, 8)
    changesetIdsApplied := make([]string, 0, 8)
    for _, toApply := range applyable {
        if err != nil {
            break
        }
        for _, changesetId := range toApply.ChangeSetIds {
            if err != nil {
                break
            }
            changeset, _ := changesets[changesetId]
            if changeset != nil {
                isMe := changeset.RecordId == meContactId
                canCreate, canUpdate, canDelete := cs.CanCreateGroup(isMe), cs.CanUpdateGroup(isMe), cs.CanDeleteGroup(isMe)
                if !canCreate && !canUpdate && !canDelete {
                    changesetIdsNotApplyable = append(changesetIdsNotApplyable, changeset.Id)
                    continue
                }
                fmt.Printf("[PIPELINE]: handling changeset id %s\n", changeset.Id)
                isCreate, isUpdate, isDelete := false, false, false
                if len(changeset.Changes) > 1 {
                    isUpdate = true
                } else if len(changeset.Changes) == 1 {
                    switch changeset.Changes[0].ChangeType {
                    case dm.CHANGE_TYPE_CREATE:
                        isCreate = true
                    case dm.CHANGE_TYPE_ADD:
                        if changeset.Changes[0].Path == nil || len(changeset.Changes[0].Path) == 0 {
                            isCreate = true
                        } else {
                            isUpdate = true
                        }
                    case dm.CHANGE_TYPE_UPDATE:
                        isUpdate = true
                    case dm.CHANGE_TYPE_DELETE:
                        if changeset.Changes[0].Path == nil || len(changeset.Changes[0].Path) == 0 {
                            isDelete = true
                        } else {
                            isUpdate = true
                        }
                    }
                }
                if !(isCreate && canCreate) && !(isUpdate && canUpdate) && !(isDelete && canDelete) {
                    changesetIdsNotApplyable = append(changesetIdsNotApplyable, changeset.Id)
                    continue
                }
                dsocialGroupId := changeset.RecordId
                if isDelete {
                    err = p.handleDeleteGroup(client, ds, cs, dsocialUserId, externalServiceId, externalUserId, dsocialGroupId)
                } else if isCreate {
                    err = p.handleCreateGroup(client, ds, cs, dsocialUserId, dsocialGroupId)
                } else {
                    // must be update
                    err = p.handleUpdateGroup(client, ds, cs, dsocialUserId, externalServiceId, externalUserId, dsocialGroupId, changeset)
                }
                changesetIdsApplied = append(changesetIdsApplied, changeset.Id)
            }
        }
    }
    if err == nil {
        fmt.Printf("[PIPELINE]: Skipped applying %d changesets\n", len(changesetIdsNotApplyable))
        err = p.markGroupChangeSetsNotApplyable(ds, dsocialUserId, externalServiceId, externalServiceName, changesetIdsNotApplyable)
        if err == nil {
            err = ds.RemoveContactChangeSetsToApply(dsocialUserId, externalServiceId, externalServiceName, changesetIdsApplied)
        }
    }
    fmt.Printf("[PIPELINE]: Done applying group changesets\n")
    return
}

func (p *Pipeline) Export(client oauth2_client.OAuth2Client, ds DataStoreService, cs ContactsService, csSettings ContactsServiceSettings, dsocialUserId, meContactId string) (err error) {
    fmt.Printf("[PIPELINE]: Starting exporting...\n")
    externalServiceId := csSettings.Id()
    externalServiceName := csSettings.ContactsServiceId()
    if !csSettings.AllowModifyContactInfo() || !cs.CanExportContactsOrGroups() {
        err = p.markAllContactChangeSetsNotApplyable(ds, dsocialUserId, externalServiceId, externalServiceName)
        if err == nil {
            p.markAllGroupChangeSetsNotApplyable(ds, dsocialUserId, externalServiceId, externalServiceName)
        }
    } else {
        err = p.applyContactChangeSets(client, ds, cs, csSettings, dsocialUserId, meContactId)
        if err == nil {
            err = p.applyGroupChangeSets(client, ds, cs, csSettings, dsocialUserId, meContactId)
        }
    }
    fmt.Printf("[PIPELINE]: Done exporting.\n")
    return
}
