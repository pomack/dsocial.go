package datastore

import (
    "github.com/pomack/jsonhelper.go/jsonhelper"
    dm "github.com/pomack/dsocial.go/models/dsocial"
    bc "github.com/pomack/dsocial.go/backend/contacts"
    "container/list"
    "fmt"
    "io"
    "json"
    "os"
    "strconv"
    "strings"
    "time"
)

const (
    _INMEMORY_CONTACT_COLLECTION_NAME = "contacts"
    _INMEMORY_CONNECTION_COLLECTION_NAME = "connections"
    _INMEMORY_GROUP_COLLECTION_NAME = "group"
    _INMEMORY_CHANGESET_COLLECTION_NAME = "changesets"
    _INMEMORY_EXTERNAL_CONTACT_COLLECTION_NAME = "external_contacts"
    _INMEMORY_EXTERNAL_GROUP_COLLECTION_NAME = "external_group"
    _INMEMORY_CONTACT_FOR_EXTERNAL_CONTACT_COLLECTION_NAME = "contacts_for_external_contacts"
    _INMEMORY_GROUP_FOR_EXTERNAL_GROUP_COLLECTION_NAME = "groups_for_external_group"
    _INMEMORY_EXTERNAL_TO_INTERNAL_CONTACT_MAPPING_COLLECTION_NAME = "external_to_internal_contact_mappings"
    _INMEMORY_INTERNAL_TO_EXTERNAL_CONTACT_MAPPING_COLLECTION_NAME = "internal_to_external_contact_mappings"
    _INMEMORY_EXTERNAL_TO_INTERNAL_GROUP_MAPPING_COLLECTION_NAME = "external_to_internal_group_mappings"
    _INMEMORY_INTERNAL_TO_EXTERNAL_GROUP_MAPPING_COLLECTION_NAME = "internal_to_external_group_mappings"
    _INMEMORY_USER_TO_CONTACT_SETTINGS_COLLECTION_NAME = "user_to_contact_settings"
)

type inMemoryObject interface{}

type inMemoryCollection struct {
    Data map[string]inMemoryObject  `json:"data"`
    Name string                     `json:"name"`
}

type InMemoryDataStore struct {
    Collections map[string]*inMemoryCollection  `json:"collections"`
    NextId int64                                `json:"next_id"`
}

func NewInMemoryDataStore() *InMemoryDataStore {
    return &InMemoryDataStore{
        Collections: make(map[string]*inMemoryCollection),
        NextId: 1,
    }
}

func (p *InMemoryDataStore) retrieveCollection(name string) (m *inMemoryCollection) {
    var ok bool
    if m, ok = p.Collections[name]; !ok {
        m = &inMemoryCollection{
            Data: make(map[string]inMemoryObject),
            Name: name,
        }
        p.Collections[name] = m
    }
    return m
}

func (p *InMemoryDataStore) retrieveContactCollection() (m *inMemoryCollection) {
    return p.retrieveCollection(_INMEMORY_CONTACT_COLLECTION_NAME)
}

func (p *InMemoryDataStore) retrieveConnectionCollection() (m *inMemoryCollection) {
    return p.retrieveCollection(_INMEMORY_CONNECTION_COLLECTION_NAME)
}

func (p *InMemoryDataStore) retrieveGroupCollection() (m *inMemoryCollection) {
    return p.retrieveCollection(_INMEMORY_GROUP_COLLECTION_NAME)
}

func (p *InMemoryDataStore) retrieveChangesetCollection() (m *inMemoryCollection) {
    return p.retrieveCollection(_INMEMORY_CHANGESET_COLLECTION_NAME)
}

func (p *InMemoryDataStore) GenerateId(dsocialUserId string, collectionName string) string {
    nextId := dsocialUserId + "/" + collectionName + "/" + strconv.Itoa64(p.NextId)
    p.NextId++
    return nextId
}

func (p *InMemoryDataStore) store(dsocialUserId, collectionName, id string, value interface{}) string {
    if id == "" {
        id = p.GenerateId(dsocialUserId, collectionName)
    }
    p.retrieveCollection(collectionName).Data[id] = inMemoryObject(value)
    return id
}

func (p *InMemoryDataStore) delete(dsocialUserId, collectionName, id string) (existed bool) {
    if id != "" {
        m := p.retrieveCollection(collectionName).Data
        _, existed = m[id]
        m[id] = nil, false
    }
    return
}

func (p *InMemoryDataStore) retrieve(collectionName, id string) (interface{}, bool) {
    if id == "" {
        return nil, false
    }
    v, ok := p.retrieveCollection(collectionName).Data[id]
    return v, ok
}


    
func (p *InMemoryDataStore) RetrieveAllContactsServiceSettingsForUser(dsocialUserId string) (settings []bc.ContactsServiceSettings, err os.Error) {
    v, ok := p.retrieve(_INMEMORY_USER_TO_CONTACT_SETTINGS_COLLECTION_NAME, dsocialUserId)
    if v == nil || !ok {
        settings = make([]bc.ContactsServiceSettings, 0)
    } else {
        settings = v.([]bc.ContactsServiceSettings)
    }
    return
}

func (p *InMemoryDataStore) RetrieveContactsServiceSettingsForService(dsocialUserId, contactsServiceId string) (settings []bc.ContactsServiceSettings, err os.Error) {
    allSettings, _ := p.RetrieveAllContactsServiceSettingsForUser(dsocialUserId)
    arr := make([]bc.ContactsServiceSettings, len(allSettings))
    i := 0
    for _, s := range allSettings {
        if s.ContactsServiceId() == contactsServiceId {
            arr[i] = s
            i++
        }
    }
    settings = arr[0:i]
    return
}

func (p *InMemoryDataStore) RetrieveContactsServiceSettings(dsocialUserId, contactsServiceId, id string) (settings bc.ContactsServiceSettings, err os.Error) {
    allSettings, _ := p.RetrieveAllContactsServiceSettingsForUser(dsocialUserId)
    for _, s := range allSettings {
        if s.Id() == id && s.ContactsServiceId() == contactsServiceId {
            settings = s
            break
        }
    }
    return
}

func (p *InMemoryDataStore) SetContactsServiceSettings(settings bc.ContactsServiceSettings) (id string, err os.Error) {
    if settings == nil {
        return
    }
    dsocialUserId := settings.DsocialUserId()
    id = settings.Id()
    if id == "" {
        id = p.GenerateId(dsocialUserId, _INMEMORY_USER_TO_CONTACT_SETTINGS_COLLECTION_NAME)
        settings.SetId(id)
    }
    v, _ := p.retrieve(_INMEMORY_USER_TO_CONTACT_SETTINGS_COLLECTION_NAME, dsocialUserId)
    var arr []bc.ContactsServiceSettings
    if v == nil {
        arr = []bc.ContactsServiceSettings{settings}
    } else {
        arr = v.([]bc.ContactsServiceSettings)
        found := false
        for i, s := range arr {
            if s.Id() == id {
                found = true
                arr[i] = settings
            }
        }
        if !found {
            arr2 := make([]bc.ContactsServiceSettings, len(arr) + 1)
            copy(arr2, arr)
            arr2[len(arr)] = settings
            arr = arr2
        }
    }
    p.store(dsocialUserId, _INMEMORY_USER_TO_CONTACT_SETTINGS_COLLECTION_NAME, dsocialUserId, arr)
    return
}

func (p *InMemoryDataStore) DeleteContactsServiceSettings(dsocialUserId, contactsServiceId, id string) (err os.Error) {
    allSettings, _ := p.RetrieveAllContactsServiceSettingsForUser(dsocialUserId)
    l := len(allSettings)
    for i, s := range allSettings {
        if s.Id() == id && s.ContactsServiceId() == contactsServiceId {
            copy(allSettings[i:l], allSettings[i+1:l])
            p.store(dsocialUserId, _INMEMORY_USER_TO_CONTACT_SETTINGS_COLLECTION_NAME, dsocialUserId, allSettings[0:l-1])
            break
        }
    }
    return
}



func (p *InMemoryDataStore) SearchForDsocialContacts(dsocialUserId string, contact *dm.Contact) (contacts []*dm.Contact, err os.Error) {
    if contact == nil {
        return make([]*dm.Contact, 0), nil
    }
    collection := p.retrieveContactCollection()
    l := list.New()
    for _, v := range collection.Data {
        if c, ok := v.(*dm.Contact); ok && c != nil {
            if isSimilar, _ := contact.IsSimilarOrUpdated(contact, c); isSimilar {
                c2 := new(dm.Contact)
                *c2 = *c
                l.PushBack(c2)
            }
        }
    }
    rc := make([]*dm.Contact, l.Len())
    for i, iter := 0, l.Front(); iter != nil; i, iter = i+1, iter.Next() {
        if iter.Value != nil {
            rc[i] = iter.Value.(*dm.Contact)
        }
    }
    return rc, nil
}

func (p *InMemoryDataStore) SearchForDsocialGroups(dsocialUserId string, groupName string) (groups []*dm.Group, err os.Error) {
    if groupName == "" {
        return make([]*dm.Group, 0), nil
    }
    collection := p.retrieveGroupCollection()
    l := list.New()
    for _, v := range collection.Data {
        if g, ok := v.(*dm.Group); ok && g != nil {
            if g.Name == groupName {
                g2 := new(dm.Group)
                *g2 = *g
                l.PushBack(g2)
            }
        }
    }
    rc := make([]*dm.Group, l.Len())
    for i, iter := 0, l.Front(); iter != nil; i, iter = i+1, iter.Next() {
        if iter.Value != nil {
            rc[i] = iter.Value.(*dm.Group)
        }
    }
    return rc, nil
}

func (p *InMemoryDataStore) storeChangeSet(dsocialUserId string, changeset *dm.ChangeSet) (*dm.ChangeSet, os.Error) {
    if changeset == nil {
        return nil, nil
    }
    if changeset.Id == "" {
        changeset.Id = p.GenerateId(dsocialUserId, _INMEMORY_CHANGESET_COLLECTION_NAME)
    }
    if changeset.CreatedAt == "" {
        changeset.CreatedAt = time.UTC().Format(dm.UTC_DATETIME_FORMAT)
    }
    obj := new(dm.ChangeSet)
    *obj = *changeset
    p.store(dsocialUserId, _INMEMORY_CHANGESET_COLLECTION_NAME, changeset.Id, obj)
    return changeset, nil
}

func (p *InMemoryDataStore) retrieveChangeSets(dsocialId string, after *time.Time) ([]*dm.ChangeSet, bc.NextToken, os.Error) {
    l := list.New()
    var afterString string
    if after != nil {
        afterString = after.Format(dm.UTC_DATETIME_FORMAT)
    }
    for _, v := range p.retrieveChangesetCollection().Data {
        if cs, ok := v.(*dm.ChangeSet); ok {
            if cs.RecordId == dsocialId {
                if after == nil || cs.CreatedAt > afterString {
                    cs2 := new(dm.ChangeSet)
                    *cs2 = *cs
                    l.PushBack(cs2)
                }
            }
        }
    }
    rc := make([]*dm.ChangeSet, l.Len())
    for i, iter := 0, l.Front(); iter != nil; i, iter = i+1, iter.Next() {
        if iter.Value != nil {
            rc[i] = iter.Value.(*dm.ChangeSet)
        }
    }
    return rc, nil, nil
}

func (p *InMemoryDataStore) StoreContactChangeSet(dsocialUserId string, changeset *dm.ChangeSet) (*dm.ChangeSet, os.Error) {
    return p.storeChangeSet(dsocialUserId, changeset)
}

func (p *InMemoryDataStore) RetrieveContactChangeSets(dsocialId string, after *time.Time) ([]*dm.ChangeSet, bc.NextToken, os.Error) {
    return p.retrieveChangeSets(dsocialId, after)
}

func (p *InMemoryDataStore) StoreGroupChangeSet(dsocialUserId string, changeset *dm.ChangeSet) (*dm.ChangeSet, os.Error) {
    return p.storeChangeSet(dsocialUserId, changeset)
}

func (p *InMemoryDataStore) RetrieveGroupChangeSets(dsocialId string, after *time.Time) ([]*dm.ChangeSet, bc.NextToken, os.Error) {
    return p.retrieveChangeSets(dsocialId, after)
}
    
    // Retrieve the dsocial contact id for the specified external service/user id/contact id combo
    // Returns:
    //   dsocialContactId : the dsocial contact id if it exists or empty if not found
    //   err : error or nil
func (p *InMemoryDataStore) DsocialIdForExternalContactId(externalServiceId, externalUserId, dsocialUserId, externalContactId string) (dsocialContactId string, err os.Error) {
    k := strings.Join([]string{externalServiceId, externalUserId, externalContactId}, "|")
    id := dsocialUserId + "/" + k
    v, ok := p.retrieve(_INMEMORY_EXTERNAL_TO_INTERNAL_CONTACT_MAPPING_COLLECTION_NAME, id)
    if ok {
        dsocialContactId, _ = v.(string)
    }
    return
}
    // Retrieve the dsocial group id for the specified external service/user id/group id combo
    // Returns:
    //   dsocialGroupId : the dsocial group id if it exists or empty if not found
    //   err : error or nil
func (p *InMemoryDataStore) DsocialIdForExternalGroupId(externalServiceId, externalUserId, dsocialUserId, externalGroupId string) (dsocialGroupId string, err os.Error) {
    k := strings.Join([]string{externalServiceId, externalUserId, externalGroupId}, "|")
    id := dsocialUserId + "/" + k
    v, ok := p.retrieve(_INMEMORY_EXTERNAL_TO_INTERNAL_GROUP_MAPPING_COLLECTION_NAME, id)
    if ok {
        dsocialGroupId, _ = v.(string)
    }
    return
}
    // Retrieve the external contact id for the specified external service/external user id/dsocial user id/dsocial contact id combo
    // Returns:
    //   externalContactId : the dsocial contact id if it exists or empty if not found
    //   err : error or nil
func (p *InMemoryDataStore) ExternalContactIdForDsocialId(externalServiceId, externalUserId, dsocialUserId, dsocialContactId string) (externalContactId string, err os.Error) {
    id := strings.Join([]string{externalServiceId, externalUserId}, "|")
    v, ok := p.retrieve(_INMEMORY_INTERNAL_TO_EXTERNAL_CONTACT_MAPPING_COLLECTION_NAME, dsocialContactId)
    if ok {
        m, _ := v.(map[string]string)
        if m != nil {
            externalContactId, _ = m[id]
        }
    }
    return
}
    // Retrieve the external group id for the specified external service/external user id/dsocial user id/dsocial group id combo
    // Returns:
    //   externalGroupId : the dsocial group id if it exists or empty if not found
    //   err : error or nil
func (p *InMemoryDataStore) ExternalGroupIdForDsocialId(externalServiceId, externalUserId, dsocialUserId, dsocialGroupId string) (externalGroupId string, err os.Error) {
    id := strings.Join([]string{externalServiceId, externalUserId}, "|")
    v, ok := p.retrieve(_INMEMORY_INTERNAL_TO_EXTERNAL_GROUP_MAPPING_COLLECTION_NAME, dsocialGroupId)
    if ok {
        m, _ := v.(map[string]string)
        if m != nil {
            externalGroupId, _ = m[id]
        }
    }
    return
}
    // Stores the dsocial contact id <-> external contact id mapping
    // Returns:
    //   externalExisted : whether the external contact id mapping already existed and was overwritten
    //   dsocialExisted : whether the dsocial contact id mapping already existed and was overwritten
    //   err : error or nil
func (p *InMemoryDataStore) StoreDsocialExternalContactMapping(externalServiceId, externalUserId, externalContactId, dsocialUserId, dsocialContactId string) (externalExisted, dsocialExisted bool, err os.Error) {
    k1 := strings.Join([]string{externalServiceId, externalUserId, externalContactId}, "|")
    k2 := strings.Join([]string{externalServiceId, externalUserId}, "|")
    id1 := dsocialUserId + "/" + k1
    v1, externalExisted := p.retrieve(_INMEMORY_EXTERNAL_TO_INTERNAL_CONTACT_MAPPING_COLLECTION_NAME, id1)
    v2, dsocialExisted := p.retrieve(_INMEMORY_INTERNAL_TO_EXTERNAL_CONTACT_MAPPING_COLLECTION_NAME, dsocialContactId)
    if v1 != dsocialContactId {
        //fmt.Printf("[DS]: Storing %s %v -> %v\n", _INMEMORY_EXTERNAL_TO_INTERNAL_CONTACT_MAPPING_COLLECTION_NAME, id1, dsocialContactId)
        p.store(dsocialUserId, _INMEMORY_EXTERNAL_TO_INTERNAL_CONTACT_MAPPING_COLLECTION_NAME, id1, dsocialContactId)
    }
    if v2 == nil {
        v2 = make(map[string]string)
    }
    if m, ok := v2.(map[string]string); ok {
        currentExternalContactId, _ := m[k2]
        if currentExternalContactId != externalContactId {
            //fmt.Printf("[DS]: Storing %s %v -> %v -> %v\n", _INMEMORY_INTERNAL_TO_EXTERNAL_CONTACT_MAPPING_COLLECTION_NAME, dsocialContactId, k2, externalContactId)
            m[k2] = externalContactId
            p.store(dsocialUserId, _INMEMORY_INTERNAL_TO_EXTERNAL_CONTACT_MAPPING_COLLECTION_NAME, dsocialContactId, m)
        }
    }
    if strings.HasPrefix(externalContactId, "testname/contact/") {
        panic(fmt.Sprintf("Invalid externalContactId: %v for key: %v", externalContactId, k1))
    }
    if !strings.HasPrefix(dsocialContactId, "testname/contact/") {
        panic(fmt.Sprintf("Invalid dsocialContactId: %v for key: %v", dsocialContactId, k1))
    }
    return
}
    // Stores the dsocial contact id <-> external group id mapping
    // Returns:
    //   externalExisted : whether the external group id mapping already existed and was overwritten
    //   dsocialExisted : whether the dsocial group id mapping already existed and was overwritten
    //   err : error or nil
func (p *InMemoryDataStore) StoreDsocialExternalGroupMapping(externalServiceId, externalUserId, externalGroupId, dsocialUserId, dsocialGroupId string) (externalExisted, dsocialExisted bool, err os.Error) {
    k1 := strings.Join([]string{externalServiceId, externalUserId, externalGroupId}, "|")
    k2 := strings.Join([]string{externalServiceId, externalUserId}, "|")
    id1 := dsocialUserId + "/" + k1
    v1, externalExisted := p.retrieve(_INMEMORY_EXTERNAL_TO_INTERNAL_GROUP_MAPPING_COLLECTION_NAME, id1)
    v2, dsocialExisted := p.retrieve(_INMEMORY_INTERNAL_TO_EXTERNAL_GROUP_MAPPING_COLLECTION_NAME, dsocialGroupId)
    if v1 != dsocialGroupId {
        //fmt.Printf("[DS]: Storing %s %v -> %v\n", _INMEMORY_EXTERNAL_TO_INTERNAL_GROUP_MAPPING_COLLECTION_NAME, id1, dsocialGroupId)
        p.store(dsocialUserId, _INMEMORY_EXTERNAL_TO_INTERNAL_GROUP_MAPPING_COLLECTION_NAME, id1, dsocialGroupId)
    }
    if v2 == nil {
        v2 = make(map[string]string)
    }
    if m, ok := v2.(map[string]string); ok {
        currentExternalGroupId, _ := m[k2]
        if currentExternalGroupId != externalGroupId {
            //fmt.Printf("[DS]: Storing %s %v -> %v -> %v\n", _INMEMORY_INTERNAL_TO_EXTERNAL_GROUP_MAPPING_COLLECTION_NAME, dsocialGroupId, k2, externalGroupId)
            m[k2] = externalGroupId
            p.store(dsocialUserId, _INMEMORY_INTERNAL_TO_EXTERNAL_GROUP_MAPPING_COLLECTION_NAME, dsocialGroupId, m)
        }
    }
    if !strings.HasPrefix(dsocialGroupId, "testname/group/") {
        panic(fmt.Sprintf("Invalid dsocialGroupId: %v for key: %v", dsocialGroupId, k1))
    }
    return
}

    // Retrieve external contact
    // Returns:
    //   externalContact : the contact as stored into the service using StoreExternalContact or nil if not found
    //   id : the internal id used to store the external contact
    //   err : error or nil
func (p *InMemoryDataStore) RetrieveExternalContact(externalServiceId, externalUserId, dsocialUserId, externalContactId string) (externalContact interface{}, id string, err os.Error) {
    k := strings.Join([]string{externalServiceId, externalUserId, externalContactId}, "|")
    id = dsocialUserId + "/" + k
    externalContact, _ = p.retrieve(_INMEMORY_EXTERNAL_CONTACT_COLLECTION_NAME, id)
    return
}
    // Retrieve external group
    // Returns:
    //   externalGroup : the group as stored into the service using StoreExternalGroup or nil if not found
    //   id : the internal id used to store the external group
    //   err : error or nil
func (p *InMemoryDataStore) RetrieveExternalGroup(externalServiceId, externalUserId, dsocialUserId, externalGroupId string) (externalGroup interface{}, id string, err os.Error) {
    k := strings.Join([]string{externalServiceId, externalUserId, externalGroupId}, "|")
    id = dsocialUserId + "/" + k
    externalGroup, _ = p.retrieve(_INMEMORY_EXTERNAL_GROUP_COLLECTION_NAME, id)
    return
}
    // Stores external contact
    // Returns:
    //   id : the internal id used to store the external contact
    //   err : error or nil
func (p *InMemoryDataStore) StoreExternalContact(externalServiceId, externalUserId, dsocialUserId, externalContactId string, contact interface{}) (id string, err os.Error) {
    //id = strings.Join([]string{dsocialUserId, externalServiceId, externalUserId, externalContactId}, "/")
    k := strings.Join([]string{externalServiceId, externalUserId, externalContactId}, "|")
    id = dsocialUserId + "/" + k
    p.store(dsocialUserId, _INMEMORY_EXTERNAL_CONTACT_COLLECTION_NAME, id, contact)
    if strings.HasPrefix(externalContactId, "testname/contact/") {
        panic(fmt.Sprintf("Storing external contact with invalid externalContactId: %v\n", externalContactId))
    }
    return
}
    // Stores external group
    // Returns:
    //   id : the internal id used to store the external group
    //   err : error or nil
func (p *InMemoryDataStore) StoreExternalGroup(externalServiceId, externalUserId, dsocialUserId, externalGroupId string, group interface{}) (id string, err os.Error) {
    //id = strings.Join([]string{dsocialUserId, externalServiceId, externalUserId, externalGroupId}, "/")
    k := strings.Join([]string{externalServiceId, externalUserId, externalGroupId}, "|")
    id = dsocialUserId + "/" + k
    p.store(dsocialUserId, _INMEMORY_EXTERNAL_GROUP_COLLECTION_NAME, id, group)
    if strings.HasPrefix(externalGroupId, "testname/group/") {
        panic(fmt.Sprintf("Storing external group with invalid externalGroupId: %v\n", externalGroupId))
    }
    return
}
    // Deletes external contact
    // Returns:
    //   existed : whether the contact existed upon deletion
    //   err : error or nil
func (p *InMemoryDataStore) DeleteExternalContact(externalServiceId, externalUserId, dsocialUserId, externalContactId string) (existed bool, err os.Error) {
    //k := strings.Join([]string{dsocialUserId, externalServiceId, externalUserId, externalContactId}, "/")
    k := strings.Join([]string{externalServiceId, externalUserId, externalContactId}, "|")
    id := dsocialUserId + "/" + k
    existed = p.delete(dsocialUserId, _INMEMORY_EXTERNAL_CONTACT_COLLECTION_NAME, id)
    return
}
    // Deletes external group
    // Returns:
    //   existed : whether the group existed upon deletion
    //   err : error or nil
func (p *InMemoryDataStore) DeleteExternalGroup(externalServiceId, externalUserId, dsocialUserId, externalGroupId string) (existed bool, err os.Error) {
    //k := strings.Join([]string{dsocialUserId, externalServiceId, externalUserId, externalGroupId}, "/")
    k := strings.Join([]string{externalServiceId, externalUserId, externalGroupId}, "|")
    id := dsocialUserId + "/" + k
    existed = p.delete(dsocialUserId, _INMEMORY_EXTERNAL_GROUP_COLLECTION_NAME, id)
    return
}
    
    
    // Retrieve dsocial contact
    // Returns:
    //   dsocialContact : the contact as stored into the service using StoreDsocialContact or nil if not found
    //   id : the internal id used to store the dsocial contact
    //   err : error or nil
func (p *InMemoryDataStore) RetrieveDsocialContactForExternalContact(externalServiceId, externalUserId, externalContactId, dsocialUserId string) (dsocialContact *dm.Contact, id string, err os.Error) {
    //k := strings.Join([]string{dsocialUserId, externalServiceId, externalUserId, externalContactId}, "/")
    k := strings.Join([]string{externalServiceId, externalUserId, externalContactId}, "|")
    extid := dsocialUserId + "/" + k
    v, _ := p.retrieve(_INMEMORY_CONTACT_FOR_EXTERNAL_CONTACT_COLLECTION_NAME, extid)
    id, _ = p.DsocialIdForExternalContactId(externalServiceId, externalUserId, dsocialUserId, externalContactId)
    c, _ := v.(*dm.Contact)
    if c != nil {
        dsocialContact = new(dm.Contact)
        *dsocialContact = *c
    }
    return
}
    // Retrieve dsocial group
    // Returns:
    //   dsocialGroup : the group as stored into the service using StoreDsocialGroup or nil if not found
    //   id : the internal id used to store the dsocial group
    //   err : error or nil
func (p *InMemoryDataStore) RetrieveDsocialGroupForExternalGroup(externalServiceId, externalUserId, externalGroupId, dsocialUserId string) (dsocialGroup *dm.Group, id string, err os.Error) {
    //k := strings.Join([]string{dsocialUserId, externalServiceId, externalUserId, externalGroupId}, "/")
    k := strings.Join([]string{externalServiceId, externalUserId, externalGroupId}, "|")
    extid := dsocialUserId + "/" + k
    v, _ := p.retrieve(_INMEMORY_GROUP_FOR_EXTERNAL_GROUP_COLLECTION_NAME, extid)
    id, _ = p.DsocialIdForExternalGroupId(externalServiceId, externalUserId, dsocialUserId, externalGroupId)
    g, _ := v.(*dm.Group)
    if g != nil {
        dsocialGroup = new(dm.Group)
        *dsocialGroup = *g
    }
    return
}
    // Stores dsocial contact
    // Returns:
    //   dsocialContact : the contact, modified to include items like Id and LastModified/Created
    //   err : error or nil
func (p *InMemoryDataStore) StoreDsocialContactForExternalContact(externalServiceId, externalUserId, externalContactId, dsocialUserId string, contact *dm.Contact) (dsocialContact *dm.Contact, err os.Error) {
    //k := strings.Join([]string{dsocialUserId, externalServiceId, externalUserId, externalContactId}, "/")
    k := strings.Join([]string{externalServiceId, externalUserId, externalContactId}, "|")
    extid := dsocialUserId + "/" + k
    bc.AddIdsForDsocialContact(contact, p, dsocialUserId)
    c := new(dm.Contact)
    *c = *contact
    c.Id = extid
    p.store(dsocialUserId, _INMEMORY_CONTACT_FOR_EXTERNAL_CONTACT_COLLECTION_NAME, extid, c)
    dsocialContact = contact
    return
}
    // Stores dsocial group
    // Returns:
    //   dsocialGroup : the group, modified to include items like Id and LastModified/Created
    //   err : error or nil
func (p *InMemoryDataStore) StoreDsocialGroupForExternalGroup(externalServiceId, externalUserId, externalGroupId, dsocialUserId string, group *dm.Group) (dsocialGroup *dm.Group, err os.Error) {
    //k := strings.Join([]string{dsocialUserId, externalServiceId, externalUserId, externalGroupId}, "/")
    k := strings.Join([]string{externalServiceId, externalUserId, externalGroupId}, "|")
    extid := dsocialUserId + "/" + k
    bc.AddIdsForDsocialGroup(group, p, dsocialUserId)
    g := new(dm.Group)
    *g = *group
    g.Id = extid
    fmt.Printf("[DS]: Storing %s with id %v for %s\n", _INMEMORY_GROUP_FOR_EXTERNAL_GROUP_COLLECTION_NAME, extid, g.Name)
    p.store(dsocialUserId, _INMEMORY_GROUP_FOR_EXTERNAL_GROUP_COLLECTION_NAME, extid, g)
    dsocialGroup = group
    return
}
    // Deletes dsocial contact
    // Returns:
    //   existed : whether the contact existed upon deletion
    //   err : error or nil
func (p *InMemoryDataStore) DeleteDsocialContactForExternalContact(externalServiceId, externalUserId, externalContactId, dsocialUserId string) (existed bool, err os.Error) {
    //k := strings.Join([]string{dsocialUserId, externalServiceId, externalUserId, externalContactId}, "|")
    k := strings.Join([]string{externalServiceId, externalUserId, externalContactId}, "|")
    extid := dsocialUserId + "/" + k
    existed = p.delete(dsocialUserId, _INMEMORY_CONTACT_FOR_EXTERNAL_CONTACT_COLLECTION_NAME, extid)
    return
}
    // Deletes dsocial group
    // Returns:
    //   existed : whether the group existed upon deletion
    //   err : error or nil
func (p *InMemoryDataStore) DeleteDsocialGroupForExternalGroup(externalServiceId, externalUserId, externalGroupId, dsocialUserId string) (existed bool, err os.Error) {
    //k := strings.Join([]string{dsocialUserId, externalServiceId, externalUserId, externalGroupId}, "/")
    k := strings.Join([]string{externalServiceId, externalUserId, externalGroupId}, "|")
    extid := dsocialUserId + "/" + k
    existed = p.delete(dsocialUserId, _INMEMORY_GROUP_FOR_EXTERNAL_GROUP_COLLECTION_NAME, extid)
    return
}


    // Retrieve dsocial contact
    // Returns:
    //   dsocialContact : the contact as stored into the service using StoreDsocialContact or nil if not found
    //   id : the internal id used to store the dsocial contact
    //   err : error or nil
func (p *InMemoryDataStore) RetrieveDsocialContact(dsocialUserId, dsocialContactId string) (dsocialContact *dm.Contact, id string, err os.Error) {
    //k := strings.Join([]string{dsocialUserId, dsocialContactId}, "/")
    v, _ := p.retrieve(_INMEMORY_CONTACT_COLLECTION_NAME, dsocialContactId)
    if v != nil {
        if contact, ok := v.(*dm.Contact); ok {
            dsocialContact = new(dm.Contact)
            *dsocialContact = *contact
            id = dsocialContact.Id
        }
    }
    return 
}
    // Retrieve dsocial group
    // Returns:
    //   dsocialGroup : the group as stored into the service using StoreDsocialGroup or nil if not found
    //   id : the internal id used to store the dsocial group
    //   err : error or nil
func (p *InMemoryDataStore) RetrieveDsocialGroup(dsocialUserId, dsocialGroupId string) (dsocialGroup *dm.Group, id string, err os.Error) {
    //k := strings.Join([]string{dsocialUserId, dsocialGroupId}, "/")
    v, _ := p.retrieve(_INMEMORY_GROUP_COLLECTION_NAME, dsocialGroupId)
    if v != nil {
        if group, ok := v.(*dm.Group); ok {
            dsocialGroup = new(dm.Group)
            *dsocialGroup = *group
            id = dsocialGroup.Id
        }
    }
    return 
}
    // Stores dsocial contact
    // Returns:
    //   dsocialContact : the contact, modified to include items like Id and LastModified/Created
    //   err : error or nil
func (p *InMemoryDataStore) StoreDsocialContact(dsocialUserId, dsocialContactId string, contact *dm.Contact) (dsocialContact *dm.Contact, err os.Error) {
    if dsocialContactId == "" {
        dsocialContactId = p.GenerateId(dsocialUserId, "contact")
        fmt.Printf("[DS]: Generated Id for storing dsocial contact: %v\n", dsocialContactId)
        contact.Id = dsocialContactId
    } else {
        fmt.Printf("[DS]: Using existing contact id: %v\n", dsocialContactId)
    }
    bc.AddIdsForDsocialContact(contact, p, dsocialUserId)
    //k := strings.Join([]string{dsocialUserId, dsocialContactId}, "/")
    c := new(dm.Contact)
    *c = *contact
    p.store(dsocialUserId, _INMEMORY_CONTACT_COLLECTION_NAME, dsocialContactId, c)
    dsocialContact = contact
    return 
}
    // Stores dsocial group
    // Returns:
    //   dsocialGroup : the group, modified to include items like Id and LastModified/Created
    //   err : error or nil
func (p *InMemoryDataStore) StoreDsocialGroup(dsocialUserId, dsocialGroupId string, group *dm.Group) (dsocialGroup *dm.Group, err os.Error) {
    if dsocialGroupId == "" {
        dsocialGroupId = p.GenerateId(dsocialUserId, "group")
        fmt.Printf("[DS]: Generated Id for storing dsocial group: %v\n", dsocialGroupId)
        group.Id = dsocialGroupId
    } else {
        fmt.Printf("[DS]: Using existing group id: %v\n", dsocialGroupId)
    }
    //k := strings.Join([]string{dsocialUserId, dsocialGroupId}, "/")
    g := new(dm.Group)
    *g = *group
    p.store(dsocialUserId, _INMEMORY_GROUP_COLLECTION_NAME, dsocialGroupId, g)
    dsocialGroup = group
    return 
}
    // Deletes dsocial contact
    // Returns:
    //   existed : whether the contact existed upon deletion
    //   err : error or nil
func (p *InMemoryDataStore) DeleteDsocialContact(dsocialUserId, dsocialContactId string) (existed bool, err os.Error) {
    //k := strings.Join([]string{dsocialUserId, dsocialContactId}, "/")
    existed = p.delete(dsocialUserId, _INMEMORY_CONTACT_COLLECTION_NAME, dsocialContactId)
    return 
}
    // Deletes dsocial group
    // Returns:
    //   existed : whether the group existed upon deletion
    //   err : error or nil
func (p *InMemoryDataStore) DeleteDsocialGroup(dsocialUserId, dsocialGroupId string) (existed bool, err os.Error) {
    //k := strings.Join([]string{dsocialUserId, dsocialGroupId}, "/")
    existed = p.delete(dsocialUserId, _INMEMORY_GROUP_COLLECTION_NAME, dsocialGroupId)
    return 
}

func (p *InMemoryDataStore) BackendContactsDataStoreService() (bc.DataStoreService) {
    return p
}

func (p *InMemoryDataStore) Encode(w io.Writer) os.Error {
    v, err := jsonhelper.MarshalWithOptions(p, dm.UTC_DATETIME_FORMAT)
    if err != nil {
        return err
    }
    return json.NewEncoder(w).Encode(v)
}

func (p *InMemoryDataStore) Decode(r io.Reader) os.Error {
    err := json.NewDecoder(r).Decode(p)
    if err != nil {
        return err
    }
    m := make(map[string]interface{})
    m[_INMEMORY_CONTACT_COLLECTION_NAME] = new(dm.Contact)
    m[_INMEMORY_CONNECTION_COLLECTION_NAME] = new(dm.Contact)
    m[_INMEMORY_GROUP_COLLECTION_NAME] = new(dm.Group)
    m[_INMEMORY_CHANGESET_COLLECTION_NAME] = new(dm.ChangeSet)
    m[_INMEMORY_CONTACT_FOR_EXTERNAL_CONTACT_COLLECTION_NAME] = new(dm.Contact)
    m[_INMEMORY_GROUP_FOR_EXTERNAL_GROUP_COLLECTION_NAME] = new(dm.Group)
    
    for k, collection := range p.Collections {
        if obj, ok := m[k]; ok {
            for k1, v1 := range collection.Data {
                m1, _ := jsonhelper.MarshalWithOptions(v1, dm.UTC_DATETIME_FORMAT)
                b1, _ := json.Marshal(m1)
                err = json.Unmarshal(b1, obj)
                if err != nil {
                    return err
                }
                collection.Data[k1] = obj
            }
        }
    }
    return nil
}

