package inmemory

import (
    "container/list"
    bc "github.com/pomack/dsocial.go/backend/contacts"
    dm "github.com/pomack/dsocial.go/models/dsocial"
    "os"
    "time"
)

func (p *InMemoryDataStore) retrieveChangesetCollection() (m *inMemoryCollection) {
    return p.retrieveCollection(_INMEMORY_CHANGESET_COLLECTION_NAME)
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

func (p *InMemoryDataStore) retrieveChangeSetsById(ids []string, m map[string]*dm.ChangeSet) map[string]*dm.ChangeSet {
    if m == nil {
        m = make(map[string]*dm.ChangeSet)
    }
    // make a set from ids to make it faster to query
    idmap := make(map[string]bool, len(ids))
    for _, id := range ids {
        idmap[id] = true
    }
    for _, v := range p.retrieveChangesetCollection().Data {
        if cs, ok := v.(*dm.ChangeSet); ok {
            if _, ok := idmap[cs.Id]; ok {
                cs2 := new(dm.ChangeSet)
                *cs2 = *cs
                m[cs.Id] = cs2
            }
        }
    }
    return m
}

func (p *InMemoryDataStore) addChangeSetsToApply(dsocialUserId, collectionName, recordType, serviceId, serviceName string, changesetIds []string) (id string, err os.Error) {
    if len(dsocialUserId) == 0 || len(recordType) == 0 || len(serviceId) == 0 || len(serviceName) == 0 || changesetIds == nil || len(changesetIds) == 0 {
        return
    }
    v, ok := p.retrieve(collectionName, dsocialUserId)
    var l *list.List
    if !ok {
        l = list.New()
        p.store(dsocialUserId, collectionName, dsocialUserId, l)
    } else {
        l = v.(*list.List)
    }
    id = p.GenerateId(dsocialUserId, collectionName)
    l.PushBack(&dm.ChangeSetsToApply{
        Id:           id,
        UserId:       dsocialUserId,
        RecordType:   recordType,
        ServiceId:    serviceId,
        ServiceName:  serviceName,
        ChangeSetIds: changesetIds,
    })
    return
}

func (p *InMemoryDataStore) retrieveChangeSetsToApply(dsocialUserId, collectionName, recordType, serviceId, serviceName string) (arr []*dm.ChangeSetsToApply, m map[string]*dm.ChangeSet, err os.Error) {
    m = make(map[string]*dm.ChangeSet)
    if dsocialUserId == "" || recordType == "" || serviceId == "" || serviceName == "" {
        arr = make([]*dm.ChangeSetsToApply, 0)
        return
    }
    v, ok := p.retrieve(collectionName, dsocialUserId)
    if !ok || v == nil {
        arr = make([]*dm.ChangeSetsToApply, 0)
    } else {
        l := v.(*list.List)
        arr = make([]*dm.ChangeSetsToApply, l.Len())
        i := 0
        for e := l.Front(); e != nil; e = e.Next() {
            ch := e.Value.(*dm.ChangeSetsToApply)
            if ch.RecordType == recordType && ch.ServiceName == serviceName && ch.ServiceId == serviceId {
                arr[i] = ch
                i++
                p.retrieveChangeSetsById(ch.ChangeSetIds, m)
            }
        }
        arr = arr[0:i]
    }
    return
}

func (p *InMemoryDataStore) removeChangeSetsToApply(dsocialUserId, collectionName, recordType string, serviceId, serviceName string, ids []string) (err os.Error) {
    if dsocialUserId == "" || collectionName == "" || recordType == "" || ids == nil || len(ids) == 0 {
        return
    }
    // make a set from ids to make it faster to query
    idmap := make(map[string]bool, len(ids))
    for _, id := range ids {
        idmap[id] = true
    }
    v, ok := p.retrieve(collectionName, dsocialUserId)
    if ok && v != nil {
        l := v.(*list.List)
        for e := l.Front(); e != nil; e = e.Next() {
            ch := e.Value.(*dm.ChangeSetsToApply)
            if ch.RecordType == recordType && ch.ServiceName == serviceName && ch.ServiceId == serviceId {
                if _, ok := idmap[ch.Id]; ok {
                    l.Remove(e)
                    break
                }
            }
        }
    }
    return
}
