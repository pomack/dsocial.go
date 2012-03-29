package inmemory

import (
    dm "github.com/pomack/dsocial.go/models/dsocial"
    "github.com/pomack/jsonhelper.go/jsonhelper"
    "io"
    "json"
    "os"
    "strconv"
)

type inMemoryObject interface{}

type inMemoryCollection struct {
    Data map[string]inMemoryObject `json:"data"`
    Name string                    `json:"name"`
}

type InMemoryDataStore struct {
    Collections map[string]*inMemoryCollection `json:"collections"`
    NextId      int64                          `json:"next_id"`
}

func NewInMemoryDataStore() *InMemoryDataStore {
    return &InMemoryDataStore{
        Collections: make(map[string]*inMemoryCollection),
        NextId:      1,
    }
}

func (p *InMemoryDataStore) GenerateId(dsocialUserId string, collectionName string) string {
    nextId := dsocialUserId + "/" + collectionName + "/" + strconv.Itoa64(p.NextId)
    p.NextId++
    return nextId
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
    // These three have different types based on the service
    //m[_INMEMORY_EXTERNAL_CONTACT_COLLECTION_NAME] = "external_contacts"
    //m[_INMEMORY_EXTERNAL_GROUP_COLLECTION_NAME] = "external_group"
    //m[_INMEMORY_USER_TO_CONTACT_SETTINGS_COLLECTION_NAME] = "user_to_contact_settings"
    // These four are all strings
    //m[_INMEMORY_EXTERNAL_TO_INTERNAL_CONTACT_MAPPING_COLLECTION_NAME] = ""
    //m[_INMEMORY_INTERNAL_TO_EXTERNAL_CONTACT_MAPPING_COLLECTION_NAME] = ""
    //m[_INMEMORY_EXTERNAL_TO_INTERNAL_GROUP_MAPPING_COLLECTION_NAME] = ""
    //m[_INMEMORY_INTERNAL_TO_EXTERNAL_GROUP_MAPPING_COLLECTION_NAME] = ""
    m[_INMEMORY_CHANGESETS_TO_APPLY_COLLECTION_NAME] = new(dm.ChangeSetsToApply)
    m[_INMEMORY_CHANGESETS_NOT_CURRENTLY_APPLYABLE_COLLECTION_NAME] = new(dm.ChangeSetsToApply)

    for k, collection := range p.Collections {
        for k1, v1 := range collection.Data {
            m1, _ := jsonhelper.MarshalWithOptions(v1, dm.UTC_DATETIME_FORMAT)
            b1, _ := json.Marshal(m1)
            if obj, ok := m[k]; ok {
                err = json.Unmarshal(b1, obj)
                if err != nil {
                    return err
                }
                collection.Data[k1] = obj
            } else {
                collection.Data[k1] = b1
            }
        }
    }
    return nil
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

func (p *InMemoryDataStore) store(dsocialUserId, collectionName, id string, value interface{}) string {
    if id == "" {
        id = p.GenerateId(dsocialUserId, collectionName)
    }
    p.retrieveCollection(collectionName).Data[id] = inMemoryObject(value)
    return id
}

func (p *InMemoryDataStore) delete(collectionName, id string) (oldValue interface{}, existed bool) {
    if id != "" {
        m := p.retrieveCollection(collectionName).Data
        oldValue, existed = m[id]
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

func (p *InMemoryDataStore) retrieveString(collectionName, id string) (string, bool) {
    if id == "" {
        return "", false
    }
    v, ok := p.retrieveCollection(collectionName).Data[id]
    if ok {
        value, _ := v.(string)
        return value, ok
    }
    return "", ok
}

func (p *InMemoryDataStore) retrieveStringMapCollection(userId, collectionName, id string) map[string]string {
    if len(collectionName) == 0 || len(id) == 0 {
        return make(map[string]string)
    }
    var m map[string]string
    if names, ok := p.retrieve(collectionName, id); ok {
        m = names.(map[string]string)
    } else {
        m = make(map[string]string)
        if userId != "" && collectionName != "" && id != "" {
            p.store(userId, collectionName, id, m)
        }
    }
    return m
}

func (p *InMemoryDataStore) retrieveFromStringMapCollection(userId, collectionName, colKey, id string) (value string, found bool) {
    if len(collectionName) == 0 || len(colKey) == 0 || len(id) == 0 {
        return
    }
    var m map[string]string
    if names, ok := p.retrieve(collectionName, id); ok {
        m = names.(map[string]string)
        value, found = m[id]
    }
    return
}

func (p *InMemoryDataStore) addToStringMapCollection(userId, collectionName, colKey, key, value string) {
    if len(collectionName) == 0 || len(colKey) == 0 || len(key) == 0 {
        return
    }
    var m map[string]string
    if names, ok := p.retrieve(collectionName, colKey); ok {
        m = names.(map[string]string)
    } else {
        m = make(map[string]string)
    }
    m[key] = value
    p.store(userId, collectionName, colKey, m)
}

func (p *InMemoryDataStore) removeFromStringMapCollection(userId, collectionName, colKey, key string) {
    if len(collectionName) == 0 || len(colKey) == 0 || len(key) == 0 {
        return
    }
    var m map[string]string
    if names, ok := p.retrieve(collectionName, colKey); ok {
        m = names.(map[string]string)
        m[key] = "", false
        if len(m) == 0 {
            p.delete(collectionName, colKey)
        } else {
            p.store(userId, collectionName, colKey, m)
        }
    }
}
