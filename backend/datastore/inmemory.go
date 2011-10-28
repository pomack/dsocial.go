package datastore

import (
    "github.com/pomack/jsonhelper.go/jsonhelper"
    dm "github.com/pomack/dsocial.go/models/dsocial"
    bc "github.com/pomack/dsocial.go/backend/contacts"
    "container/list"
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

func (p *InMemoryDataStore) GenerateId(collectionName string) string {
    nextId := collectionName + "/" + strconv.Itoa64(p.NextId)
    p.NextId++
    return nextId
}

func (p *InMemoryDataStore) store(collectionName, id string, value interface{}) string {
    if id == "" {
        id = p.GenerateId(collectionName)
    }
    p.retrieveCollection(collectionName).Data[id] = inMemoryObject(value)
    return id
}

func (p *InMemoryDataStore) delete(collectionName, id string) (existed bool) {
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

func (p *InMemoryDataStore) storeChangeSet(changeset *dm.ChangeSet) (*dm.ChangeSet, os.Error) {
    if changeset == nil {
        return nil, nil
    }
    if changeset.Id == "" {
        changeset.Id = p.GenerateId(_INMEMORY_CHANGESET_COLLECTION_NAME)
    }
    if changeset.CreatedAt == "" {
        changeset.CreatedAt = time.UTC().Format(dm.UTC_DATETIME_FORMAT)
    }
    obj := new(dm.ChangeSet)
    *obj = *changeset
    p.store(_INMEMORY_CHANGESET_COLLECTION_NAME, changeset.Id, obj)
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

func (p *InMemoryDataStore) StoreContactChangeSet(changeset *dm.ChangeSet) (*dm.ChangeSet, os.Error) {
    return p.storeChangeSet(changeset)
}

func (p *InMemoryDataStore) RetrieveContactChangeSets(dsocialId string, after *time.Time) ([]*dm.ChangeSet, bc.NextToken, os.Error) {
    return p.retrieveChangeSets(dsocialId, after)
}

func (p *InMemoryDataStore) StoreGroupChangeSet(changeset *dm.ChangeSet) (*dm.ChangeSet, os.Error) {
    return p.storeChangeSet(changeset)
}

func (p *InMemoryDataStore) RetrieveGroupChangeSets(dsocialId string, after *time.Time) ([]*dm.ChangeSet, bc.NextToken, os.Error) {
    return p.retrieveChangeSets(dsocialId, after)
}
    
    // Retrieve the dsocial contact id for the specified external service/user id/contact id combo
    // Returns:
    //   dsocialContactId : the dsocial contact id if it exists or empty if not found
    //   err : error or nil
func (p *InMemoryDataStore) DsocialIdForExternalContactId(externalServiceId, externalUserId, dsocialUserId, externalContactId string) (dsocialContactId string, err os.Error) {
    v, ok := p.retrieve(_INMEMORY_EXTERNAL_TO_INTERNAL_CONTACT_MAPPING_COLLECTION_NAME, strings.Join([]string{dsocialUserId, externalServiceId, externalUserId, externalContactId}, "/"))
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
    v, ok := p.retrieve(_INMEMORY_EXTERNAL_TO_INTERNAL_GROUP_MAPPING_COLLECTION_NAME, strings.Join([]string{dsocialUserId, externalServiceId, externalUserId, externalGroupId}, "/"))
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
    v, ok := p.retrieve(_INMEMORY_INTERNAL_TO_EXTERNAL_CONTACT_MAPPING_COLLECTION_NAME, strings.Join([]string{dsocialUserId, externalServiceId, externalUserId, dsocialContactId}, "/"))
    if ok {
        externalContactId, _ = v.(string)
    }
    return
}
    // Retrieve the external group id for the specified external service/external user id/dsocial user id/dsocial group id combo
    // Returns:
    //   externalGroupId : the dsocial group id if it exists or empty if not found
    //   err : error or nil
func (p *InMemoryDataStore) ExternalGroupIdForDsocialId(externalServiceId, externalUserId, dsocialUserId, dsocialGroupId string) (externalGroupId string, err os.Error) {
    v, ok := p.retrieve(_INMEMORY_INTERNAL_TO_EXTERNAL_GROUP_MAPPING_COLLECTION_NAME, strings.Join([]string{dsocialUserId, externalServiceId, externalUserId, dsocialGroupId}, "/"))
    if ok {
        externalGroupId, _ = v.(string)
    }
    return
}
    // Stores the dsocial contact id <-> external contact id mapping
    // Returns:
    //   externalExisted : whether the external contact id mapping already existed and was overwritten
    //   dsocialExisted : whether the dsocial contact id mapping already existed and was overwritten
    //   err : error or nil
func (p *InMemoryDataStore) StoreDsocialExternalContactMapping(externalServiceId, externalUserId, externalContactId, dsocialUserId, dsocialContactId string) (externalExisted, dsocialExisted bool, err os.Error) {
    k1 := strings.Join([]string{dsocialUserId, externalServiceId, externalUserId, externalContactId}, "/")
    k2 := strings.Join([]string{dsocialUserId, externalServiceId, externalUserId, dsocialContactId}, "/")
    v1, externalExisted := p.retrieve(_INMEMORY_EXTERNAL_TO_INTERNAL_CONTACT_MAPPING_COLLECTION_NAME, k1)
    v2, dsocialExisted := p.retrieve(_INMEMORY_INTERNAL_TO_EXTERNAL_CONTACT_MAPPING_COLLECTION_NAME, k2)
    if v1 != dsocialContactId {
        p.store(_INMEMORY_EXTERNAL_TO_INTERNAL_CONTACT_MAPPING_COLLECTION_NAME, k1, dsocialContactId)
    }
    if v2 != externalContactId {
        p.store(_INMEMORY_INTERNAL_TO_EXTERNAL_CONTACT_MAPPING_COLLECTION_NAME, k2, externalContactId)
    }
    return
}
    // Stores the dsocial contact id <-> external group id mapping
    // Returns:
    //   externalExisted : whether the external group id mapping already existed and was overwritten
    //   dsocialExisted : whether the dsocial group id mapping already existed and was overwritten
    //   err : error or nil
func (p *InMemoryDataStore) StoreDsocialExternalGroupMapping(externalServiceId, externalUserId, externalGroupId, dsocialUserId, dsocialGroupId string) (externalExisted, dsocialExisted bool, err os.Error) {
    k1 := strings.Join([]string{dsocialUserId, externalServiceId, externalUserId, externalGroupId}, "/")
    k2 := strings.Join([]string{dsocialUserId, externalServiceId, externalUserId, dsocialGroupId}, "/")
    v1, externalExisted := p.retrieve(_INMEMORY_EXTERNAL_TO_INTERNAL_GROUP_MAPPING_COLLECTION_NAME, k1)
    v2, dsocialExisted := p.retrieve(_INMEMORY_INTERNAL_TO_EXTERNAL_GROUP_MAPPING_COLLECTION_NAME, k2)
    if v1 != dsocialGroupId {
        p.store(_INMEMORY_EXTERNAL_TO_INTERNAL_GROUP_MAPPING_COLLECTION_NAME, k1, dsocialGroupId)
    }
    if v2 != externalGroupId {
        p.store(_INMEMORY_INTERNAL_TO_EXTERNAL_GROUP_MAPPING_COLLECTION_NAME, k2, externalGroupId)
    }
    return
}

    // Retrieve external contact
    // Returns:
    //   externalContact : the contact as stored into the service using StoreExternalContact or nil if not found
    //   id : the internal id used to store the external contact
    //   err : error or nil
func (p *InMemoryDataStore) RetrieveExternalContact(externalServiceId, externalUserId, dsocialUserId, externalContactId string) (externalContact interface{}, id string, err os.Error) {
    id = strings.Join([]string{dsocialUserId, externalServiceId, externalUserId, externalContactId}, "/")
    externalContact, _ = p.retrieve(_INMEMORY_EXTERNAL_CONTACT_COLLECTION_NAME, id)
    return
}
    // Retrieve external group
    // Returns:
    //   externalGroup : the group as stored into the service using StoreExternalGroup or nil if not found
    //   id : the internal id used to store the external group
    //   err : error or nil
func (p *InMemoryDataStore) RetrieveExternalGroup(externalServiceId, externalUserId, dsocialUserId, externalGroupId string) (externalGroup interface{}, id string, err os.Error) {
    id = strings.Join([]string{dsocialUserId, externalServiceId, externalUserId, externalGroupId}, "/")
    externalGroup, _ = p.retrieve(_INMEMORY_EXTERNAL_GROUP_COLLECTION_NAME, id)
    return
}
    // Stores external contact
    // Returns:
    //   id : the internal id used to store the external contact
    //   err : error or nil
func (p *InMemoryDataStore) StoreExternalContact(externalServiceId, externalUserId, dsocialUserId, externalContactId string, contact interface{}) (id string, err os.Error) {
    id = strings.Join([]string{dsocialUserId, externalServiceId, externalUserId, externalContactId}, "/")
    p.store(_INMEMORY_EXTERNAL_CONTACT_COLLECTION_NAME, id, contact)
    return
}
    // Stores external group
    // Returns:
    //   id : the internal id used to store the external group
    //   err : error or nil
func (p *InMemoryDataStore) StoreExternalGroup(externalServiceId, externalUserId, dsocialUserId, externalGroupId string, group interface{}) (id string, err os.Error) {
    id = strings.Join([]string{dsocialUserId, externalServiceId, externalUserId, externalGroupId}, "/")
    p.store(_INMEMORY_EXTERNAL_GROUP_COLLECTION_NAME, id, group)
    return
}
    // Deletes external contact
    // Returns:
    //   existed : whether the contact existed upon deletion
    //   err : error or nil
func (p *InMemoryDataStore) DeleteExternalContact(externalServiceId, externalUserId, dsocialUserId, externalContactId string) (existed bool, err os.Error) {
    k := strings.Join([]string{dsocialUserId, externalServiceId, externalUserId, externalContactId}, "/")
    existed = p.delete(_INMEMORY_EXTERNAL_CONTACT_COLLECTION_NAME, k)
    return
}
    // Deletes external group
    // Returns:
    //   existed : whether the group existed upon deletion
    //   err : error or nil
func (p *InMemoryDataStore) DeleteExternalGroup(externalServiceId, externalUserId, dsocialUserId, externalGroupId string) (existed bool, err os.Error) {
    k := strings.Join([]string{dsocialUserId, externalServiceId, externalUserId, externalGroupId}, "/")
    existed = p.delete(_INMEMORY_EXTERNAL_GROUP_COLLECTION_NAME, k)
    return
}
    
    
    // Retrieve dsocial contact
    // Returns:
    //   dsocialContact : the contact as stored into the service using StoreDsocialContact or nil if not found
    //   id : the internal id used to store the dsocial contact
    //   err : error or nil
func (p *InMemoryDataStore) RetrieveDsocialContactForExternalContact(externalServiceId, externalUserId, externalContactId, dsocialUserId string) (dsocialContact *dm.Contact, id string, err os.Error) {
    k := strings.Join([]string{dsocialUserId, externalServiceId, externalUserId, externalContactId}, "/")
    v, _ := p.retrieve(_INMEMORY_CONTACT_FOR_EXTERNAL_CONTACT_COLLECTION_NAME, k)
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
    k := strings.Join([]string{dsocialUserId, externalServiceId, externalUserId, externalGroupId}, "/")
    v, _ := p.retrieve(_INMEMORY_GROUP_FOR_EXTERNAL_GROUP_COLLECTION_NAME, k)
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
    k := strings.Join([]string{dsocialUserId, externalServiceId, externalUserId, externalContactId}, "/")
    c := new(dm.Contact)
    *c = *contact
    p.store(_INMEMORY_CONTACT_FOR_EXTERNAL_CONTACT_COLLECTION_NAME, k, c)
    dsocialContact = contact
    return
}
    // Stores dsocial group
    // Returns:
    //   dsocialGroup : the group, modified to include items like Id and LastModified/Created
    //   err : error or nil
func (p *InMemoryDataStore) StoreDsocialGroupForExternalGroup(externalServiceId, externalUserId, externalGroupId, dsocialUserId string, group *dm.Group) (dsocialGroup *dm.Group, err os.Error) {
    k := strings.Join([]string{dsocialUserId, externalServiceId, externalUserId, externalGroupId}, "/")
    g := new(dm.Group)
    *g = *group
    p.store(_INMEMORY_GROUP_FOR_EXTERNAL_GROUP_COLLECTION_NAME, k, g)
    dsocialGroup = group
    return
}
    // Deletes dsocial contact
    // Returns:
    //   existed : whether the contact existed upon deletion
    //   err : error or nil
func (p *InMemoryDataStore) DeleteDsocialContactForExternalContact(externalServiceId, externalUserId, externalContactId, dsocialUserId string) (existed bool, err os.Error) {
    k := strings.Join([]string{dsocialUserId, externalServiceId, externalUserId, externalContactId}, "/")
    existed = p.delete(_INMEMORY_CONTACT_FOR_EXTERNAL_CONTACT_COLLECTION_NAME, k)
    return
}
    // Deletes dsocial group
    // Returns:
    //   existed : whether the group existed upon deletion
    //   err : error or nil
func (p *InMemoryDataStore) DeleteDsocialGroupForExternalGroup(externalServiceId, externalUserId, externalGroupId, dsocialUserId string) (existed bool, err os.Error) {
    k := strings.Join([]string{dsocialUserId, externalServiceId, externalUserId, externalGroupId}, "/")
    existed = p.delete(_INMEMORY_GROUP_FOR_EXTERNAL_GROUP_COLLECTION_NAME, k)
    return
}


    // Retrieve dsocial contact
    // Returns:
    //   dsocialContact : the contact as stored into the service using StoreDsocialContact or nil if not found
    //   id : the internal id used to store the dsocial contact
    //   err : error or nil
func (p *InMemoryDataStore) RetrieveDsocialContact(dsocialUserId, dsocialContactId string) (dsocialContact *dm.Contact, id string, err os.Error) {
    k := strings.Join([]string{dsocialUserId, dsocialContactId}, "/")
    v, _ := p.retrieve(_INMEMORY_CONTACT_COLLECTION_NAME, k)
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
    k := strings.Join([]string{dsocialUserId, dsocialGroupId}, "/")
    v, _ := p.retrieve(_INMEMORY_GROUP_COLLECTION_NAME, k)
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
        dsocialContactId = p.GenerateId("contact")
        contact.Id = dsocialContactId
    }
    k := strings.Join([]string{dsocialUserId, dsocialContactId}, "/")
    c := new(dm.Contact)
    *c = *contact
    p.store(_INMEMORY_CONTACT_COLLECTION_NAME, k, c)
    dsocialContact = contact
    return 
}
    // Stores dsocial group
    // Returns:
    //   dsocialGroup : the group, modified to include items like Id and LastModified/Created
    //   err : error or nil
func (p *InMemoryDataStore) StoreDsocialGroup(dsocialUserId, dsocialGroupId string, group *dm.Group) (dsocialGroup *dm.Group, err os.Error) {
    if dsocialGroupId == "" {
        dsocialGroupId = p.GenerateId("group")
        group.Id = dsocialGroupId
    }
    k := strings.Join([]string{dsocialUserId, dsocialGroupId}, "/")
    g := new(dm.Group)
    *g = *group
    p.store(_INMEMORY_GROUP_COLLECTION_NAME, k, g)
    dsocialGroup = group
    return 
}
    // Deletes dsocial contact
    // Returns:
    //   existed : whether the contact existed upon deletion
    //   err : error or nil
func (p *InMemoryDataStore) DeleteDsocialContact(dsocialUserId, dsocialContactId string) (existed bool, err os.Error) {
    k := strings.Join([]string{dsocialUserId, dsocialContactId}, "/")
    existed = p.delete(_INMEMORY_CONTACT_COLLECTION_NAME, k)
    return 
}
    // Deletes dsocial group
    // Returns:
    //   existed : whether the group existed upon deletion
    //   err : error or nil
func (p *InMemoryDataStore) DeleteDsocialGroup(dsocialUserId, dsocialGroupId string) (existed bool, err os.Error) {
    k := strings.Join([]string{dsocialUserId, dsocialGroupId}, "/")
    existed = p.delete(_INMEMORY_GROUP_COLLECTION_NAME, k)
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

