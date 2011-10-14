package dsocial

import (
    "fmt"
    "reflect"
    "strconv"
    "strings"
)

type ChangeSet struct {
    Id              string      `json:"id"`
    RecordId        string      `json:"record_id"`
    CreatedAt       string      `json:"created_at"`
    ChangedBy       string      `json:"changed_by"`
    ChangeImportId  string      `json:"change_import_id,omitempty"`
    Changes         []*Change   `json:"changes"`
}

type Change struct {
    Id              string              `json:"id"`
    RecordId        string              `json:"record_id"`
    Path            []*PathComponent    `json:"path"`
    ChangeType      ChangeType          `json:"change_type"`
    OriginalValue   interface{}         `json:"original_value,omitempty"`
    NewValue        interface{}         `json:"new_value,omitempty"`
}

type ChangeType string

type PathComponent struct {
    Type    int             `json:"type"`
    Value   string          `json:"value"`
}

func (p *PathComponent) IsId() bool { return p.Type == PATH_COMPONENT_TYPE_ID }
func (p *PathComponent) IsKey() bool { return p.Type == PATH_COMPONENT_TYPE_KEY }
func (p *PathComponent) IsIndex() bool { return p.Type == PATH_COMPONENT_TYPE_INDEX }
func (p *PathComponent) IsMapIndex() bool { return p.Type == PATH_COMPONENT_TYPE_MAP_INDEX }

func NewPathComponentId(value string) *PathComponent{
    return &PathComponent{Type:PATH_COMPONENT_TYPE_ID, Value:value}
}

func NewPathComponentKey(value string) *PathComponent{
    return &PathComponent{Type:PATH_COMPONENT_TYPE_KEY, Value:value}
}

func NewPathComponentIndex(value string) *PathComponent{
    return &PathComponent{Type:PATH_COMPONENT_TYPE_INDEX, Value:value}
}

func NewPathComponentMapIndex(value string) *PathComponent{
    return &PathComponent{Type:PATH_COMPONENT_TYPE_MAP_INDEX, Value:value}
}

func NewPathComponentFromExisting(pathComponent []*PathComponent, parts ...*PathComponent) []*PathComponent {
    if pathComponent == nil {
        return parts
    }
    l := len(pathComponent)
    p := make([]*PathComponent, l + len(parts))
    for i, path := range pathComponent {
        p[i] = path
    }
    for i, path := range parts {
        p[i + l] = path
    }
    return p
}

func ApplyAddChange(itemToModify interface{}, value interface{}, path []*PathComponent) {
    if value == nil { return }
    if itemToModify == nil { return }
    v := reflect.Indirect(reflect.ValueOf(itemToModify))
    lpath := len(path)
    valueType := reflect.TypeOf(value)
    valueValue := reflect.ValueOf(value)
    for i, component := range path {
        isLast := i + 1 == lpath
        switch component.Type {
        case PATH_COMPONENT_TYPE_ID:
            l := v.Len()
            var nextV reflect.Value
            for k := 0; k < l && !nextV.IsValid(); k++ {
                pelem := v.Index(k)
                if pelem.IsNil() { continue }
                elem := reflect.Indirect(pelem)
                idValue := elem.FieldByName("Id")
                if idValue.IsValid() {
                    if idValue.String() == component.Value {
                        nextV = elem
                        break
                    }
                }
            }
            if nextV.IsValid() {
                v = nextV
            } else {
                // TODO figure out what to do when cannot find component to add to
                return
            }
        case PATH_COMPONENT_TYPE_KEY:
            var nextV reflect.Value
            stype := v.Type()
            l := stype.NumField()
            for k := 0; k < l; k++ {
                sfield := stype.Field(k)
                tagparts := strings.Split(sfield.Tag.Get("json"), ",")
                fieldName := sfield.Name
                if len(tagparts) > 0 && len(tagparts[0]) > 0 {
                    fieldName = tagparts[0]
                }
                if fieldName == component.Value {
                    vField := v.Field(k)
                    nextV = reflect.Indirect(vField)
                    if isLast {
                        fieldKind := sfield.Type.Kind()
                        if sfield.Type.AssignableTo(valueType) {
                            vField.Set(valueValue)
                        } else if fieldKind == reflect.Array || fieldKind == reflect.Slice {
                            if nextV.IsNil() {
                                slice := reflect.MakeSlice(sfield.Type, 1, 10)
                                slice.Index(0).Set(valueValue)
                                vField.Set(slice)
                            } else {
                                vField.Set(reflect.Append(nextV.Slice(0, nextV.Len()), valueValue))
                            }
                        } else {
                            // TODO figure out what to do when cannot find component to add to
                            panic("1 cannot figure out how to add " + fieldName)
                        }
                        return
                    }
                    break
                }
            }
            if nextV.IsValid() {
                v = nextV
            } else {
                // TODO figure out what to do when cannot find component to add to
                panic(fmt.Sprintf("2 cannot figure out how to add %#v to %s with value %#v", component.Value, v.Type().String(), value))
                return
            }
        case PATH_COMPONENT_TYPE_INDEX:
            atIndex, _ := strconv.Atoi(component.Value)
            if isLast {
                if atIndex < 0 {
                    slice := reflect.MakeSlice(v.Type(), 1, v.Cap()+1)
                    slice.Index(0).Set(valueValue)
                    reflect.AppendSlice(slice, v.Slice(0, v.Len()))
                } else if v.Len() <= atIndex {
                    reflect.Append(v.Slice(0, v.Len()), valueValue)
                } else {
                    slice := reflect.MakeSlice(v.Type(), 1, v.Cap()+1)
                    slice.Index(0).Set(valueValue)
                    v.Set(reflect.AppendSlice(v.Slice(0, atIndex), reflect.AppendSlice(slice, v.Slice(atIndex, v.Len()))))
                }
                return
            } else {
                if atIndex >= 0 && atIndex < v.Len() {
                    v = reflect.Indirect(v.Index(atIndex))
                } else {
                    // TODO figure out what to do when cannot find component to add to
                    return
                }
            }
        case PATH_COMPONENT_TYPE_MAP_INDEX:
            mapindex := reflect.ValueOf(component.Value)
            nextV := v.MapIndex(mapindex)
            if isLast {
                v.SetMapIndex(mapindex, valueValue)
                return
            }
            if !nextV.IsValid() || nextV.IsNil() {
                // TODO figure out what to do when cannot find component to add to
                return
            }
            v = nextV
        }
    }
}


func ApplyDeleteChange(itemToModify interface{}, valueToDelete interface{}, path []*PathComponent) {
    if itemToModify == nil { return }
    v := reflect.Indirect(reflect.ValueOf(itemToModify))
    lpath := len(path)
    for i, component := range path {
        isLast := i + 1 == lpath
        switch component.Type {
        case PATH_COMPONENT_TYPE_ID:
            l := v.Len()
            var nextV reflect.Value
            for k := 0; k < l && !nextV.IsValid(); k++ {
                pelem := v.Index(k)
                if pelem.IsNil() { continue }
                elem := reflect.Indirect(pelem)
                idValue := elem.FieldByName("Id")
                if idValue.IsValid() {
                    if idValue.String() == component.Value {
                        if isLast {
                            v.Set(reflect.AppendSlice(v.Slice(0, k), v.Slice(k+1, v.Len())))
                            return
                        }
                        nextV = elem
                        break
                    }
                }
            }
            if nextV.IsValid() {
                v = nextV
            } else {
                return
            }
        case PATH_COMPONENT_TYPE_KEY:
            l := v.NumField()
            var nextV reflect.Value
            stype := v.Type()
            for k := 0; k < l && !nextV.IsValid(); k++ {
                sfield := stype.Field(k)
                tagparts := strings.Split(sfield.Tag.Get("json"), ",")
                fieldName := sfield.Name
                if len(tagparts) > 0 && len(tagparts[0]) > 0 {
                    fieldName = tagparts[0]
                }
                if fieldName == component.Value {
                    vField := v.Field(k)
                    nextV = reflect.Indirect(vField)
                    if isLast {
                        fieldKind := sfield.Type.Kind()
                        if fieldKind == reflect.Array || fieldKind == reflect.Slice {
                            if !nextV.IsValid() || nextV.IsNil() {
                                // nothing to do
                            } else {
                                if valueToDelete == nil {
                                    vField.Set(reflect.Zero(sfield.Type))
                                } else {
                                    l1 := nextV.Len()
                                    for i1 := 0; i1 < l1; i1++ {
                                        if reflect.DeepEqual(nextV.Index(i1).Interface(), valueToDelete) {
                                            if i1 + 1 >= l1 {
                                                vField.Set(nextV.Slice(0, i1))
                                            } else {
                                                vField.Set(reflect.AppendSlice(nextV.Slice(0, i1), nextV.Slice(i1+1, l1)))
                                            }
                                            break
                                        }
                                    }
                                }
                            }
                        } else {
                            vField.Set(reflect.Zero(sfield.Type))
                        }
                        return
                    }
                }
            }
            if nextV.IsValid() {
                v = nextV
            } else {
                return
            }
        case PATH_COMPONENT_TYPE_INDEX:
            atIndex, _ := strconv.Atoi(component.Value)
            if isLast {
                if atIndex >= 0 && atIndex < v.Len() {
                    if atIndex + 1 == v.Len() {
                        v.Set(v.Slice(0, atIndex))
                    } else {
                        v.Set(reflect.AppendSlice(v.Slice(0, atIndex), v.Slice(atIndex+1, v.Len())))
                    }
                    return
                }
                return
            } else {
                if atIndex >= 0 && atIndex < v.Len() {
                    v = reflect.Indirect(v.Index(atIndex))
                } else {
                    return
                }
            }
        case PATH_COMPONENT_TYPE_MAP_INDEX:
            mapindex := reflect.ValueOf(component.Value)
            nextV := reflect.Indirect(v.MapIndex(mapindex))
            if isLast {
                v.SetMapIndex(mapindex, reflect.Zero(v.Type().Elem()))
                return
            }
            if !nextV.IsValid() || nextV.IsNil() {
                return
            }
            v = nextV
        }
    }
}


func ApplyUpdateChange(itemToModify interface{}, original interface{}, latest interface{}, path []*PathComponent) {
    if latest == nil { return }
    if itemToModify == nil { return }
    v := reflect.Indirect(reflect.ValueOf(itemToModify))
    lpath := len(path)
    latestType := reflect.TypeOf(latest)
    latestValue := reflect.ValueOf(latest)
    for i, component := range path {
        isLast := i + 1 == lpath
        switch component.Type {
        case PATH_COMPONENT_TYPE_ID:
            l := v.Len()
            var nextV reflect.Value
            for k := 0; k < l && !nextV.IsValid(); k++ {
                pelem := v.Index(k)
                if pelem.IsNil() { continue }
                elem := reflect.Indirect(pelem)
                idValue := elem.FieldByName("Id")
                if idValue.IsValid() {
                    if idValue.String() == component.Value {
                        nextV = elem
                        break
                    }
                }
            }
            if nextV.IsValid() {
                v = nextV
            } else {
                // TODO figure out what to do when cannot find component to add to
                panic("Could not find id " + component.Value)
                return
            }
        case PATH_COMPONENT_TYPE_KEY:
            var nextV reflect.Value
            stype := v.Type()
            l := stype.NumField()
            for k := 0; k < l; k++ {
                sfield := stype.Field(k)
                tagparts := strings.Split(sfield.Tag.Get("json"), ",")
                fieldName := sfield.Name
                if len(tagparts) > 0 && len(tagparts[0]) > 0 {
                    fieldName = tagparts[0]
                }
                if fieldName == component.Value {
                    vField := v.Field(k)
                    nextV = reflect.Indirect(vField)
                    if isLast {
                        fieldKind := sfield.Type.Kind()
                        if sfield.Type.AssignableTo(latestType) {
                            if reflect.DeepEqual(vField.Interface(), original) {
                                fmt.Sprintf("Changing value from %#v to %#v", vField.Interface(), latestValue)
                                vField.Set(latestValue)
                            } else {
                                panic(fmt.Sprintf("1 cannot update because field %s is %#v not %#v which was expected to update to %#v", sfield.Name, nextV.Interface(), original, latest))
                            }
                        } else if fieldKind == reflect.Array || fieldKind == reflect.Slice {
                            found := false
                            if nextV.IsNil() {
                                // TODO figure out what to do when cannot find component to update to
                            } else {
                                l1 := nextV.Len()
                                for i1 := 0; i1 < l1; i1++ {
                                    if reflect.DeepEqual(reflect.Indirect(nextV.Index(i1)).Interface(), original) {
                                        found = true
                                        fmt.Sprintf("Changing value from %#v to %#v", reflect.Indirect(nextV.Index(i1)).Interface(), latestValue)
                                        nextV.Index(i1).Set(latestValue)
                                        break
                                    }
                                }
                            }
                            if !found {
                                panic(fmt.Sprintf("2 cannot update because field %s is %#v which does not contain %#v which was expected to update to %#v", sfield.Name, nextV, original, latest))
                            }
                        } else {
                            // TODO figure out what to do when cannot find component to update to
                            panic("3 cannot figure out how to update " + fieldName)
                        }
                        return
                    }
                    break
                }
            }
            if nextV.IsValid() {
                v = nextV
            } else {
                // TODO figure out what to do when cannot find component to update to
                panic(fmt.Sprintf("4 cannot figure out how to update %#v to %s with value %#v", component.Value, v.Type().String(), latest))
                return
            }
        case PATH_COMPONENT_TYPE_INDEX:
            atIndex, _ := strconv.Atoi(component.Value)
            if isLast {
                // treat update for index as insert at location
                // and remove existing
                l1 := v.Len()
                found := false
                curValue := v.Slice(0, l1)
                for i1 := 0; i1 < l1; i1++ {
                    if reflect.DeepEqual(reflect.Indirect(v.Index(i1)).Interface(), original) {
                        found = true
                        if i1 + 1 >= l1 {
                            curValue = v.Slice(0, i1)
                        } else {
                            curValue = reflect.AppendSlice(v.Slice(0, i1), v.Slice(i1+1, l1))
                        }
                        break
                    }
                }
                if !found {
                    panic(fmt.Sprintf("5 cannot update index %s because original value %#v not found", component.Value, original))
                }
                if atIndex < 0 {
                    slice := reflect.MakeSlice(curValue.Type(), 1, curValue.Cap()+1)
                    slice.Index(0).Set(latestValue)
                    reflect.AppendSlice(slice, curValue.Slice(0, curValue.Len()))
                } else if v.Len() <= atIndex {
                    reflect.Append(curValue.Slice(0, curValue.Len()), latestValue)
                } else {
                    slice := reflect.MakeSlice(curValue.Type(), 1, curValue.Cap()+1)
                    slice.Index(0).Set(latestValue)
                    v.Set(reflect.AppendSlice(curValue.Slice(0, atIndex), reflect.AppendSlice(slice, curValue.Slice(atIndex, curValue.Len()))))
                }
                return
            } else {
                if atIndex >= 0 && atIndex < v.Len() {
                    v = reflect.Indirect(v.Index(atIndex))
                } else {
                    // TODO figure out what to do when cannot find component to update to
                    panic("Could not find index " + component.Value)
                    return
                }
            }
        case PATH_COMPONENT_TYPE_MAP_INDEX:
            mapindex := reflect.ValueOf(component.Value)
            nextV := v.MapIndex(mapindex)
            if isLast {
                if nextV.IsValid() && !nextV.IsNil() && reflect.DeepEqual(reflect.Indirect(nextV).Interface(), original) {
                    v.SetMapIndex(mapindex, latestValue)
                } else {
                    panic(fmt.Sprintf("6 cannot update map index %s because original value %#v not found", component.Value, original))
                }
                return
            }
            if !nextV.IsValid() || nextV.IsNil() {
                // TODO figure out what to do when cannot find component to update to
                return
            }
            v = nextV
        }
    }
}

func ApplyChange(itemToModify interface{}, ch *Change) {
    if ch == nil { return }
    switch ch.ChangeType {
    case CHANGE_TYPE_ADD:
        ApplyAddChange(itemToModify, ch.NewValue, ch.Path)
    case CHANGE_TYPE_DELETE:
        ApplyDeleteChange(itemToModify, ch.OriginalValue, ch.Path)
    case CHANGE_TYPE_UPDATE:
        ApplyUpdateChange(itemToModify, ch.OriginalValue, ch.NewValue, ch.Path)
    }
}
