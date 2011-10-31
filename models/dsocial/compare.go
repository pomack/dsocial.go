package dsocial

import (
    "container/list"
    //"fmt"
    "sort"
    "strings"
)

type DsocialChanger interface {
    IsSimilarOrUpdated(original, latest DsocialChanger) (similar bool, same bool)
    GenerateChanges(original, latest DsocialChanger, basePath []*PathComponent, l *list.List)
    IsEmpty() bool
}

func (p *ContactReference) IsSimilarOrUpdated(original, latest DsocialChanger) (similar bool, same bool) {
    if original == nil && latest == nil {
        return true, true
    }
    if original == nil || latest == nil {
        return
    }
    o, ok := original.(*ContactReference)
    if !ok {
        return
    }
    m, ok := latest.(*ContactReference)
    if !ok {
        return
    }
    if o.ReferenceContactName != m.ReferenceContactName {
        return
    }
    if o.ReferenceContactId == "" || m.ReferenceContactId == "" {
        similar = true
    } else if p.ReferenceContactId == m.ReferenceContactId {
        similar = true
    } else {
        return
    }
    same = o.Text == m.Text
    return
}

func (p *ContactReference) GenerateChanges(original, latest DsocialChanger, basePath []*PathComponent, l *list.List) {
    if original == nil {
        ch := &Change{
            Path:       basePath,
            ChangeType: CHANGE_TYPE_ADD,
            NewValue:   latest,
        }
        l.PushBack(ch)
    } else if latest == nil {
        old := original.(*ContactReference)
        ch := &Change{
            Path:          NewPathComponentFromExisting(basePath, NewPathComponentId(old.Id)),
            ChangeType:    CHANGE_TYPE_DELETE,
            OriginalValue: original,
        }
        l.PushBack(ch)
    } else {
        old := original.(*ContactReference)
        now := latest.(*ContactReference)
        if ch := compareStrings(old.Text, now.Text, l); ch != nil {
            ch.Path = NewPathComponentFromExisting(basePath, NewPathComponentId(old.Id), NewPathComponentKey("text"))
        }
        if now.UserId != "" {
            if ch := compareStrings(old.UserId, now.UserId, l); ch != nil {
                ch.Path = NewPathComponentFromExisting(basePath, NewPathComponentId(old.Id), NewPathComponentKey("user_id"))
            }
        }
        if now.ContactId != "" {
            if ch := compareStrings(old.ContactId, now.ContactId, l); ch != nil {
                ch.Path = NewPathComponentFromExisting(basePath, NewPathComponentId(old.Id), NewPathComponentKey("contact_id"))
            }
        }
        if now.ReferenceContactId != "" {
            if ch := compareStrings(old.ReferenceContactId, now.ReferenceContactId, l); ch != nil {
                ch.Path = NewPathComponentFromExisting(basePath, NewPathComponentId(old.Id), NewPathComponentKey("reference_contact_id"))
            }
        }
    }
}

func (p *ContactReference) IsEmpty() bool {
    return p.Text == "" && p.ReferenceContactName == "" && p.ReferenceContactId == ""
}

func (p *ContactReference) GenerateArrayChanges(original, latest []*ContactReference, l *list.List, basePath []*PathComponent) {
    oarr, larr := p.ConvertArraysToDsocialChangers(original, latest)
    compareDsocialChangersArrays(oarr, larr, basePath, l)
}

func (p *ContactReference) ConvertArraysToDsocialChangers(original, latest []*ContactReference) ([]DsocialChanger, []DsocialChanger) {
    oarr := make([]DsocialChanger, len(original))
    larr := make([]DsocialChanger, len(latest))
    for i, obj := range original {
        oarr[i] = obj
    }
    for i, obj := range latest {
        larr[i] = obj
    }
    return oarr, larr
}

func (p *ContactReference) ArraysAreSame(original, latest []*ContactReference) bool {
    return arraysAreSame(p.ConvertArraysToDsocialChangers(original, latest))
}

func (p *ContactRef) IsSimilarOrUpdated(original, latest DsocialChanger) (similar bool, same bool) {
    if original == nil && latest == nil {
        return true, true
    }
    if original == nil || latest == nil {
        return
    }
    o, ok := original.(*ContactRef)
    if !ok {
        return
    }
    m, ok := latest.(*ContactRef)
    if !ok {
        return
    }
    if o.Id == m.Id && o.Id != "" {
        similar = true
    }
    if o.Name != m.Name {
        return
    }
    if o.Id == "" || m.Id == "" {
        similar = true
    } else if p.Id == m.Id {
        similar = true
        same = true
    }
    return
}

func (p *ContactRef) GenerateChanges(original, latest DsocialChanger, basePath []*PathComponent, l *list.List) {
    if original == nil {
        ch := &Change{
            Path:       basePath,
            ChangeType: CHANGE_TYPE_ADD,
            NewValue:   latest,
        }
        l.PushBack(ch)
    } else if latest == nil {
        old := original.(*ContactRef)
        ch := &Change{
            Path:          NewPathComponentFromExisting(basePath, NewPathComponentId(old.Id)),
            ChangeType:    CHANGE_TYPE_DELETE,
            OriginalValue: original,
        }
        l.PushBack(ch)
    } else {
        old := original.(*ContactRef)
        now := latest.(*ContactRef)
        if now.Id != "" {
            if ch := compareStrings(old.Id, now.Id, l); ch != nil {
                ch.Path = NewPathComponentFromExisting(basePath, NewPathComponentId(old.Id), NewPathComponentKey("id"))
            }
        }
        if now.Name != "" {
            if ch := compareStrings(old.Name, now.Name, l); ch != nil {
                ch.Path = NewPathComponentFromExisting(basePath, NewPathComponentId(old.Id), NewPathComponentKey("name"))
            }
        }
    }
}

func (p *ContactRef) IsEmpty() bool {
    return p.Id == "" && p.Name == ""
}

func (p *ContactRef) GenerateArrayChanges(original, latest []*ContactRef, l *list.List, basePath []*PathComponent) {
    oarr, larr := p.ConvertArraysToDsocialChangers(original, latest)
    compareDsocialChangersArrays(oarr, larr, basePath, l)
}

func (p *ContactRef) ConvertArraysToDsocialChangers(original, latest []*ContactRef) ([]DsocialChanger, []DsocialChanger) {
    oarr := make([]DsocialChanger, len(original))
    larr := make([]DsocialChanger, len(latest))
    for i, obj := range original {
        oarr[i] = obj
    }
    for i, obj := range latest {
        larr[i] = obj
    }
    return oarr, larr
}

func (p *ContactRef) ArraysAreSame(original, latest []*ContactRef) bool {
    return arraysAreSame(p.ConvertArraysToDsocialChangers(original, latest))
}

func (p *GroupRef) IsSimilarOrUpdated(original, latest DsocialChanger) (similar bool, same bool) {
    if original == nil && latest == nil {
        return true, true
    }
    if original == nil || latest == nil {
        return
    }
    o, ok := original.(*GroupRef)
    if !ok {
        return
    }
    m, ok := latest.(*GroupRef)
    if !ok {
        return
    }
    if o.Id == m.Id && o.Id != "" {
        similar = true
    }
    if o.Name != m.Name {
        return
    }
    if o.Id == "" || m.Id == "" {
        similar = true
    } else if p.Id == m.Id {
        similar = true
        same = true
    }
    return
}

func (p *GroupRef) GenerateChanges(original, latest DsocialChanger, basePath []*PathComponent, l *list.List) {
    if original == nil {
        ch := &Change{
            Path:       basePath,
            ChangeType: CHANGE_TYPE_ADD,
            NewValue:   latest,
        }
        l.PushBack(ch)
    } else if latest == nil {
        old := original.(*GroupRef)
        ch := &Change{
            Path:          NewPathComponentFromExisting(basePath, NewPathComponentId(old.Id)),
            ChangeType:    CHANGE_TYPE_DELETE,
            OriginalValue: original,
        }
        l.PushBack(ch)
    } else {
        old := original.(*GroupRef)
        now := latest.(*GroupRef)
        if now.Id != "" {
            if ch := compareStrings(old.Id, now.Id, l); ch != nil {
                ch.Path = NewPathComponentFromExisting(basePath, NewPathComponentId(old.Id), NewPathComponentKey("id"))
            }
        }
        if now.Name != "" {
            if ch := compareStrings(old.Name, now.Name, l); ch != nil {
                ch.Path = NewPathComponentFromExisting(basePath, NewPathComponentId(old.Id), NewPathComponentKey("name"))
            }
        }
    }
}

func (p *GroupRef) IsEmpty() bool {
    return p.Id == "" && p.Name == ""
}

func (p *GroupRef) GenerateArrayChanges(original, latest []*GroupRef, l *list.List, basePath []*PathComponent) {
    oarr, larr := p.ConvertArraysToDsocialChangers(original, latest)
    compareDsocialChangersArrays(oarr, larr, basePath, l)
}

func (p *GroupRef) ConvertArraysToDsocialChangers(original, latest []*GroupRef) ([]DsocialChanger, []DsocialChanger) {
    oarr := make([]DsocialChanger, len(original))
    larr := make([]DsocialChanger, len(latest))
    for i, obj := range original {
        oarr[i] = obj
    }
    for i, obj := range latest {
        larr[i] = obj
    }
    return oarr, larr
}

func (p *GroupRef) ArraysAreSame(original, latest []*GroupRef) bool {
    return arraysAreSame(p.ConvertArraysToDsocialChangers(original, latest))
}

func (p *PostalAddress) IsSimilarOrUpdated(original, latest DsocialChanger) (similar bool, same bool) {
    if original == nil && latest == nil {
        return true, true
    }
    if original == nil || latest == nil {
        return
    }
    o, ok := original.(*PostalAddress)
    if !ok {
        return
    }
    m, ok := latest.(*PostalAddress)
    if !ok {
        return
    }
    if o.Id == m.Id && o.Id != "" {
        similar = true
    } else if o.Id != "" && m.Id != "" {
        return false, false
    } else if checkSimilarStrings(o.Address, m.Address) {
        similar = true
    }
    if (checkSimilarStrings(o.Label, m.Label) ||
        checkSimilarStrings(string(o.Rel), string(m.Rel))) &&
        (checkSimilarStrings(o.Address, m.Address)) || (o.StreetAddress == m.StreetAddress &&
        o.OtherAddress == m.OtherAddress &&
        o.Municipality == m.Municipality &&
        o.Region == m.Region &&
        o.PostalCode == m.PostalCode &&
        o.Country == m.Country &&
        (len(o.StreetAddress)+len(o.OtherAddress)+len(o.Municipality)+len(o.Region)+len(o.PostalCode)+len(o.Country) != 0)) {
        similar = true
        same = o.Address == m.Address &&
            o.StreetAddress == m.StreetAddress &&
            o.OtherAddress == m.OtherAddress &&
            o.Municipality == m.Municipality &&
            o.Region == m.Region &&
            o.PostalCode == m.PostalCode &&
            o.Country == m.Country &&
            isSameDate(o.LocatedFrom, m.LocatedFrom) &&
            isSameDate(o.LocatedTill, m.LocatedTill) &&
            o.IsCurrent == m.IsCurrent &&
            o.IsPrimary == m.IsPrimary &&
            new(ContactReference).ArraysAreSame(o.References, m.References)
    }
    return
}

func (p *PostalAddress) GenerateChanges(original, latest DsocialChanger, basePath []*PathComponent, l *list.List) {
    if original == nil {
        ch := &Change{
            Path:       basePath,
            ChangeType: CHANGE_TYPE_ADD,
            NewValue:   latest,
        }
        l.PushBack(ch)
    } else if latest == nil {
        old := original.(*PostalAddress)
        ch := &Change{
            Path:          NewPathComponentFromExisting(basePath, NewPathComponentId(old.Id)),
            ChangeType:    CHANGE_TYPE_DELETE,
            OriginalValue: original,
        }
        l.PushBack(ch)
    } else {
        old := original.(*PostalAddress)
        now := latest.(*PostalAddress)
        if ch := compareStrings(old.Address, now.Address, l); ch != nil {
            ch.Path = NewPathComponentFromExisting(basePath, NewPathComponentId(old.Id), NewPathComponentKey("address"))
        }
        if ch := compareStrings(old.StreetAddress, now.StreetAddress, l); ch != nil {
            ch.Path = NewPathComponentFromExisting(basePath, NewPathComponentId(old.Id), NewPathComponentKey("street_address"))
        }
        if ch := compareStrings(old.OtherAddress, now.OtherAddress, l); ch != nil {
            ch.Path = NewPathComponentFromExisting(basePath, NewPathComponentId(old.Id), NewPathComponentKey("other_address"))
        }
        if ch := compareStrings(old.Municipality, now.Municipality, l); ch != nil {
            ch.Path = NewPathComponentFromExisting(basePath, NewPathComponentId(old.Id), NewPathComponentKey("municipality"))
        }
        if ch := compareStrings(old.Region, now.Region, l); ch != nil {
            ch.Path = NewPathComponentFromExisting(basePath, NewPathComponentId(old.Id), NewPathComponentKey("region"))
        }
        if ch := compareStrings(old.PostalCode, now.PostalCode, l); ch != nil {
            ch.Path = NewPathComponentFromExisting(basePath, NewPathComponentId(old.Id), NewPathComponentKey("postal_code"))
        }
        if ch := compareStrings(old.Country, now.Country, l); ch != nil {
            ch.Path = NewPathComponentFromExisting(basePath, NewPathComponentId(old.Id), NewPathComponentKey("country"))
        }
        if ch := compareDates(old.LocatedFrom, now.LocatedFrom, l); ch != nil {
            ch.Path = NewPathComponentFromExisting(basePath, NewPathComponentId(old.Id), NewPathComponentKey("located_from"))
        }
        if ch := compareDates(old.LocatedTill, now.LocatedTill, l); ch != nil {
            ch.Path = NewPathComponentFromExisting(basePath, NewPathComponentId(old.Id), NewPathComponentKey("located_till"))
        }
        if ch := compareBools(old.IsCurrent, now.IsCurrent, l); ch != nil {
            ch.Path = NewPathComponentFromExisting(basePath, NewPathComponentId(old.Id), NewPathComponentKey("is_current"))
        }
        if ch := compareBools(old.IsPrimary, now.IsPrimary, l); ch != nil {
            ch.Path = NewPathComponentFromExisting(basePath, NewPathComponentId(old.Id), NewPathComponentKey("is_primary"))
        }
        new(ContactReference).GenerateArrayChanges(old.References, now.References, l, NewPathComponentFromExisting(basePath, NewPathComponentId(old.Id), NewPathComponentKey("references")))
    }
}

func (p *PostalAddress) IsEmpty() bool {
    return p.Address == "" && p.StreetAddress == "" && p.OtherAddress == "" && p.Municipality == "" && p.Region == "" && p.PostalCode == "" && p.Country == ""
}

func (p *PostalAddress) GenerateArrayChanges(original, latest []*PostalAddress, l *list.List, basePath []*PathComponent) {
    oarr, larr := p.ConvertArraysToDsocialChangers(original, latest)
    compareDsocialChangersArrays(oarr, larr, basePath, l)
}

func (p *PostalAddress) ConvertArraysToDsocialChangers(original, latest []*PostalAddress) ([]DsocialChanger, []DsocialChanger) {
    oarr := make([]DsocialChanger, len(original))
    larr := make([]DsocialChanger, len(latest))
    for i, obj := range original {
        oarr[i] = obj
    }
    for i, obj := range latest {
        larr[i] = obj
    }
    return oarr, larr
}

func (p *PostalAddress) ArraysAreSame(original, latest []*PostalAddress) bool {
    return arraysAreSame(p.ConvertArraysToDsocialChangers(original, latest))
}

func (p *Education) IsSimilarOrUpdated(original, latest DsocialChanger) (similar bool, same bool) {
    if original == nil && latest == nil {
        return true, true
    }
    if original == nil || latest == nil {
        return
    }
    o, ok := original.(*Education)
    if !ok {
        return
    }
    m, ok := latest.(*Education)
    if !ok {
        return
    }
    if o.Id == m.Id && o.Id != "" {
        similar = true
    } else if o.Id != "" && m.Id != "" {
        return false, false
    } else if (checkSimilarStrings(string(o.Rel), string(m.Rel)) || o.GraduationYear == m.GraduationYear) &&
        checkSimilarStrings(o.Institution, m.Institution) &&
        (o.GraduationYear != 0) {
        similar = true
    }
    if similar {
        same = o.Label == m.Label &&
            o.Rel == m.Rel &&
            o.Institution == m.Institution &&
            o.GraduationYear == m.GraduationYear &&
            new(Degree).ArraysAreSame(o.Degrees, m.Degrees) &&
            unorderedStringArraysAreSame(o.Minors, m.Minors) &&
            isSameDate(o.AttendedFrom, m.AttendedFrom) &&
            isSameDate(o.AttendedTill, m.AttendedTill) &&
            o.Gpa == m.Gpa &&
            o.MajorGpa == m.MajorGpa &&
            o.IsCurrent == m.IsCurrent &&
            o.Graduated == m.Graduated &&
            new(ContactReference).ArraysAreSame(o.References, m.References) &&
            o.Notes == m.Notes &&
            unorderedStringArraysAreSame(o.Activities, m.Activities)
    }
    return
}

func (p *Education) GenerateChanges(original, latest DsocialChanger, basePath []*PathComponent, l *list.List) {
    if original == nil {
        ch := &Change{
            Path:       basePath,
            ChangeType: CHANGE_TYPE_ADD,
            NewValue:   latest,
        }
        l.PushBack(ch)
    } else if latest == nil {
        old := original.(*Education)
        ch := &Change{
            Path:          NewPathComponentFromExisting(basePath, NewPathComponentId(old.Id)),
            ChangeType:    CHANGE_TYPE_DELETE,
            OriginalValue: original,
        }
        l.PushBack(ch)
    } else {
        old := original.(*Education)
        now := latest.(*Education)
        if ch := compareStrings(old.Label, now.Label, l); ch != nil {
            ch.Path = NewPathComponentFromExisting(basePath, NewPathComponentId(old.Id), NewPathComponentKey("label"))
        }
        if ch := compareStrings(string(old.Rel), string(now.Rel), l); ch != nil {
            ch.Path = NewPathComponentFromExisting(basePath, NewPathComponentId(old.Id), NewPathComponentKey("rel"))
        }
        new(Degree).GenerateArrayChanges(old.Degrees, now.Degrees, l, NewPathComponentFromExisting(basePath, NewPathComponentId(old.Id), NewPathComponentKey("degrees")))
        compareUnorderedStringArrays(old.Minors, now.Minors, NewPathComponentFromExisting(basePath, NewPathComponentId(old.Id), NewPathComponentKey("minors")), l)
        if ch := compareInt16s(old.GraduationYear, now.GraduationYear, l); ch != nil {
            ch.Path = NewPathComponentFromExisting(basePath, NewPathComponentId(old.Id), NewPathComponentKey("graduation_year"))
        }
        if ch := compareStrings(old.Institution, now.Institution, l); ch != nil {
            ch.Path = NewPathComponentFromExisting(basePath, NewPathComponentId(old.Id), NewPathComponentKey("institution"))
        }
        if ch := compareDates(old.AttendedFrom, now.AttendedFrom, l); ch != nil {
            ch.Path = NewPathComponentFromExisting(basePath, NewPathComponentId(old.Id), NewPathComponentKey("attended_from"))
        }
        if ch := compareDates(old.AttendedTill, now.AttendedTill, l); ch != nil {
            ch.Path = NewPathComponentFromExisting(basePath, NewPathComponentId(old.Id), NewPathComponentKey("attended_till"))
        }
        if ch := compareFloat64s(old.Gpa, now.Gpa, l); ch != nil {
            ch.Path = NewPathComponentFromExisting(basePath, NewPathComponentId(old.Id), NewPathComponentKey("gpa"))
        }
        if ch := compareFloat64s(old.MajorGpa, now.MajorGpa, l); ch != nil {
            ch.Path = NewPathComponentFromExisting(basePath, NewPathComponentId(old.Id), NewPathComponentKey("major_gpa"))
        }
        if ch := compareBools(old.IsCurrent, now.IsCurrent, l); ch != nil {
            ch.Path = NewPathComponentFromExisting(basePath, NewPathComponentId(old.Id), NewPathComponentKey("is_current"))
        }
        if ch := compareBools(old.Graduated, now.Graduated, l); ch != nil {
            ch.Path = NewPathComponentFromExisting(basePath, NewPathComponentId(old.Id), NewPathComponentKey("graduated"))
        }
        new(ContactReference).GenerateArrayChanges(old.References, now.References, l, NewPathComponentFromExisting(basePath, NewPathComponentId(old.Id), NewPathComponentKey("references")))
        if ch := compareStrings(old.Notes, now.Notes, l); ch != nil {
            ch.Path = NewPathComponentFromExisting(basePath, NewPathComponentId(old.Id), NewPathComponentKey("notes"))
        }
        compareUnorderedStringArrays(old.Activities, now.Activities, NewPathComponentFromExisting(basePath, NewPathComponentId(old.Id), NewPathComponentKey("activities")), l)
    }
}

func (p *Education) IsEmpty() bool {
    return (p.Degrees == nil || len(p.Degrees) == 0) && (p.Minors == nil || len(p.Minors) == 0) && p.GraduationYear == 0 && p.Institution == "" && (p.AttendedFrom == nil || p.AttendedFrom.IsEmpty()) && (p.AttendedTill == nil || p.AttendedTill.IsEmpty()) && p.Gpa == 0.0 && p.MajorGpa == 0.0 && (p.References == nil || len(p.References) == 0) && p.Notes == "" && (p.Activities == nil || len(p.Activities) == 0)
}

func (p *Education) GenerateArrayChanges(original, latest []*Education, l *list.List, basePath []*PathComponent) {
    oarr, larr := p.ConvertArraysToDsocialChangers(original, latest)
    compareDsocialChangersArrays(oarr, larr, basePath, l)
}

func (p *Education) ConvertArraysToDsocialChangers(original, latest []*Education) ([]DsocialChanger, []DsocialChanger) {
    oarr := make([]DsocialChanger, len(original))
    larr := make([]DsocialChanger, len(latest))
    for i, obj := range original {
        oarr[i] = obj
    }
    for i, obj := range latest {
        larr[i] = obj
    }
    return oarr, larr
}

func (p *Education) ArraysAreSame(original, latest []*Education) bool {
    return arraysAreSame(p.ConvertArraysToDsocialChangers(original, latest))
}

func (p *Degree) IsSimilarOrUpdated(original, latest DsocialChanger) (similar bool, same bool) {
    if original == nil && latest == nil {
        return true, true
    }
    if original == nil || latest == nil {
        return
    }
    o, ok := original.(*Degree)
    if !ok {
        return
    }
    m, ok := latest.(*Degree)
    if !ok {
        return
    }
    if o.Id == m.Id && o.Id != "" {
        similar = true
    } else if o.Id != "" && m.Id != "" {
        return false, false
    } else if checkSimilarStrings(o.Degree, m.Degree) && checkSimilarStrings(o.Major, m.Major) {
        similar = true
    }
    if similar {
        same = o.Degree == m.Degree &&
            o.Major == m.Major
    }
    return
}

func (p *Degree) GenerateChanges(original, latest DsocialChanger, basePath []*PathComponent, l *list.List) {
    if original == nil {
        ch := &Change{
            Path:       basePath,
            ChangeType: CHANGE_TYPE_ADD,
            NewValue:   latest,
        }
        l.PushBack(ch)
    } else if latest == nil {
        old := original.(*Degree)
        ch := &Change{
            Path:          NewPathComponentFromExisting(basePath, NewPathComponentId(old.Id)),
            ChangeType:    CHANGE_TYPE_DELETE,
            OriginalValue: original,
        }
        l.PushBack(ch)
    } else {
        old := original.(*Degree)
        now := latest.(*Degree)
        if ch := compareStrings(old.Degree, now.Degree, l); ch != nil {
            ch.Path = NewPathComponentFromExisting(basePath, NewPathComponentId(old.Id), NewPathComponentKey("degree"))
        }
        if ch := compareStrings(old.Major, now.Major, l); ch != nil {
            ch.Path = NewPathComponentFromExisting(basePath, NewPathComponentId(old.Id), NewPathComponentKey("major"))
        }
    }
}

func (p *Degree) IsEmpty() bool {
    return p.Degree == "" && p.Major == ""
}

func (p *Degree) GenerateArrayChanges(original, latest []*Degree, l *list.List, basePath []*PathComponent) {
    oarr, larr := p.ConvertArraysToDsocialChangers(original, latest)
    compareDsocialChangersArrays(oarr, larr, basePath, l)
}

func (p *Degree) ConvertArraysToDsocialChangers(original, latest []*Degree) ([]DsocialChanger, []DsocialChanger) {
    oarr := make([]DsocialChanger, len(original))
    larr := make([]DsocialChanger, len(latest))
    for i, obj := range original {
        oarr[i] = obj
    }
    for i, obj := range latest {
        larr[i] = obj
    }
    return oarr, larr
}

func (p *Degree) ArraysAreSame(original, latest []*Degree) bool {
    return arraysAreSame(p.ConvertArraysToDsocialChangers(original, latest))
}

func (p *WorkPosition) IsSimilarOrUpdated(original, latest DsocialChanger) (similar bool, same bool) {
    if original == nil && latest == nil {
        return true, true
    }
    if original == nil || latest == nil {
        return
    }
    o, ok := original.(*WorkPosition)
    if !ok {
        return
    }
    m, ok := latest.(*WorkPosition)
    if !ok {
        return
    }
    if o.Id == m.Id && o.Id != "" {
        similar = true
    } else if o.Id != "" && m.Id != "" {
        return false, false
    } else if checkSimilarStrings(o.Title, m.Title) && checkSimilarStrings(o.Department, m.Department) {
        similar = true
    }
    if similar {
        same = o.Title == m.Title &&
            o.Department == m.Department &&
            o.Location == m.Location &&
            isSameDate(o.From, m.From) &&
            isSameDate(o.To, m.To) &&
            o.IsCurrent == m.IsCurrent &&
            o.Description == m.Description &&
            new(ContactReference).ArraysAreSame(o.References, m.References)
    }
    return
}

func (p *WorkPosition) GenerateChanges(original, latest DsocialChanger, basePath []*PathComponent, l *list.List) {
    if original == nil {
        ch := &Change{
            Path:       basePath,
            ChangeType: CHANGE_TYPE_ADD,
            NewValue:   latest,
        }
        l.PushBack(ch)
    } else if latest == nil {
        old := original.(*WorkPosition)
        ch := &Change{
            Path:          NewPathComponentFromExisting(basePath, NewPathComponentId(old.Id)),
            ChangeType:    CHANGE_TYPE_DELETE,
            OriginalValue: original,
        }
        l.PushBack(ch)
    } else {
        old := original.(*WorkPosition)
        now := latest.(*WorkPosition)
        if ch := compareStrings(old.Title, now.Title, l); ch != nil {
            ch.Path = NewPathComponentFromExisting(basePath, NewPathComponentId(old.Id), NewPathComponentKey("title"))
        }
        if ch := compareStrings(old.Department, now.Department, l); ch != nil {
            ch.Path = NewPathComponentFromExisting(basePath, NewPathComponentId(old.Id), NewPathComponentKey("department"))
        }
        if ch := compareStrings(old.Location, now.Location, l); ch != nil {
            ch.Path = NewPathComponentFromExisting(basePath, NewPathComponentId(old.Id), NewPathComponentKey("location"))
        }
        if ch := compareDates(old.From, now.From, l); ch != nil {
            ch.Path = NewPathComponentFromExisting(basePath, NewPathComponentId(old.Id), NewPathComponentKey("from"))
        }
        if ch := compareDates(old.To, now.To, l); ch != nil {
            ch.Path = NewPathComponentFromExisting(basePath, NewPathComponentId(old.Id), NewPathComponentKey("to"))
        }
        if ch := compareBools(old.IsCurrent, now.IsCurrent, l); ch != nil {
            ch.Path = NewPathComponentFromExisting(basePath, NewPathComponentId(old.Id), NewPathComponentKey("is_current"))
        }
        if ch := compareStrings(old.Description, now.Description, l); ch != nil {
            ch.Path = NewPathComponentFromExisting(basePath, NewPathComponentId(old.Id), NewPathComponentKey("description"))
        }
        new(ContactReference).GenerateArrayChanges(old.References, now.References, l, NewPathComponentFromExisting(basePath, NewPathComponentId(old.Id), NewPathComponentKey("references")))
    }
}

func (p *WorkPosition) IsEmpty() bool {
    return p.Title == "" && p.Department == "" && p.Location == "" && p.Description == ""
}

func (p *WorkPosition) GenerateArrayChanges(original, latest []*WorkPosition, l *list.List, basePath []*PathComponent) {
    oarr, larr := p.ConvertArraysToDsocialChangers(original, latest)
    compareDsocialChangersArrays(oarr, larr, basePath, l)
}

func (p *WorkPosition) ConvertArraysToDsocialChangers(original, latest []*WorkPosition) ([]DsocialChanger, []DsocialChanger) {
    oarr := make([]DsocialChanger, len(original))
    larr := make([]DsocialChanger, len(latest))
    for i, obj := range original {
        oarr[i] = obj
    }
    for i, obj := range latest {
        larr[i] = obj
    }
    return oarr, larr
}

func (p *WorkPosition) ArraysAreSame(original, latest []*WorkPosition) bool {
    return arraysAreSame(p.ConvertArraysToDsocialChangers(original, latest))
}

func (p *WorkHistory) IsSimilarOrUpdated(original, latest DsocialChanger) (similar bool, same bool) {
    if original == nil && latest == nil {
        return true, true
    }
    if original == nil || latest == nil {
        return
    }
    o, ok := original.(*WorkHistory)
    if !ok {
        return
    }
    m, ok := latest.(*WorkHistory)
    if !ok {
        return
    }
    if o.Id == m.Id && o.Id != "" {
        similar = true
    } else if o.Id != "" && m.Id != "" {
        return false, false
    } else if checkSimilarStrings(o.Company, m.Company) || checkSimilarStrings(o.Description, m.Description) {
        similar = true
    }
    if similar {
        same = o.Company == o.Company &&
            isSameDate(o.From, m.From) &&
            isSameDate(o.To, m.To) &&
            o.IsCurrent == m.IsCurrent &&
            o.Description == m.Description &&
            new(WorkPosition).ArraysAreSame(o.Positions, m.Positions)
    }
    return
}

func (p *WorkHistory) GenerateChanges(original, latest DsocialChanger, basePath []*PathComponent, l *list.List) {
    if original == nil {
        ch := &Change{
            Path:       basePath,
            ChangeType: CHANGE_TYPE_ADD,
            NewValue:   latest,
        }
        l.PushBack(ch)
    } else if latest == nil {
        old := original.(*WorkHistory)
        ch := &Change{
            Path:          NewPathComponentFromExisting(basePath, NewPathComponentId(old.Id)),
            ChangeType:    CHANGE_TYPE_DELETE,
            OriginalValue: original,
        }
        l.PushBack(ch)
    } else {
        old := original.(*WorkHistory)
        now := latest.(*WorkHistory)
        if ch := compareStrings(old.Company, now.Company, l); ch != nil {
            ch.Path = NewPathComponentFromExisting(basePath, NewPathComponentId(old.Id), NewPathComponentKey("company"))
        }
        if ch := compareDates(old.From, now.From, l); ch != nil {
            ch.Path = NewPathComponentFromExisting(basePath, NewPathComponentId(old.Id), NewPathComponentKey("from"))
        }
        if ch := compareDates(old.To, now.To, l); ch != nil {
            ch.Path = NewPathComponentFromExisting(basePath, NewPathComponentId(old.Id), NewPathComponentKey("to"))
        }
        if ch := compareBools(old.IsCurrent, now.IsCurrent, l); ch != nil {
            ch.Path = NewPathComponentFromExisting(basePath, NewPathComponentId(old.Id), NewPathComponentKey("is_current"))
        }
        if ch := compareStrings(old.Description, now.Description, l); ch != nil {
            ch.Path = NewPathComponentFromExisting(basePath, NewPathComponentId(old.Id), NewPathComponentKey("description"))
        }
        new(WorkPosition).GenerateArrayChanges(old.Positions, now.Positions, l, NewPathComponentFromExisting(basePath, NewPathComponentId(old.Id), NewPathComponentKey("positions")))
    }
}

func (p *WorkHistory) IsEmpty() bool {
    return p.Company == "" && p.Description == "" && (p.Positions == nil || len(p.Positions) == 0 || (len(p.Positions) == 1 && p.Positions[0].IsEmpty()))
}

func (p *WorkHistory) GenerateArrayChanges(original, latest []*WorkHistory, l *list.List, basePath []*PathComponent) {
    oarr, larr := p.ConvertArraysToDsocialChangers(original, latest)
    compareDsocialChangersArrays(oarr, larr, basePath, l)
}

func (p *WorkHistory) ConvertArraysToDsocialChangers(original, latest []*WorkHistory) ([]DsocialChanger, []DsocialChanger) {
    oarr := make([]DsocialChanger, len(original))
    larr := make([]DsocialChanger, len(latest))
    for i, obj := range original {
        oarr[i] = obj
    }
    for i, obj := range latest {
        larr[i] = obj
    }
    return oarr, larr
}

func (p *WorkHistory) ArraysAreSame(original, latest []*WorkHistory) bool {
    return arraysAreSame(p.ConvertArraysToDsocialChangers(original, latest))
}

func (p *PhoneNumber) IsSimilarOrUpdated(original, latest DsocialChanger) (similar bool, same bool) {
    if original == nil && latest == nil {
        return true, true
    }
    if original == nil || latest == nil {
        return
    }
    o, ok := original.(*PhoneNumber)
    if !ok {
        return
    }
    m, ok := latest.(*PhoneNumber)
    if !ok {
        return
    }
    if o.Id == m.Id && o.Id != "" {
        similar = true
    } else if o.Id != "" && m.Id != "" {
        return false, false
    } else if !o.IsEmpty() &&
        (checkSimilarStrings(o.FormattedNumber, m.FormattedNumber) ||
            (checkSimilarStrings(o.AreaCode, m.AreaCode) &&
                checkSimilarStrings(o.LocalPhoneNumber, m.LocalPhoneNumber) &&
                checkSimilarStrings(o.ExtensionNumber, m.ExtensionNumber))) {
        similar = true
    }
    if similar {
        same = o.Label == m.Label &&
            o.Rel == m.Rel &&
            o.FormattedNumber == m.FormattedNumber &&
            o.CountryCode == m.CountryCode &&
            o.AreaCode == m.AreaCode &&
            o.LocalPhoneNumber == m.LocalPhoneNumber &&
            o.ExtensionNumber == m.ExtensionNumber &&
            o.IsPrimary == m.IsPrimary
    }
    return
}

func (p *PhoneNumber) GenerateChanges(original, latest DsocialChanger, basePath []*PathComponent, l *list.List) {
    if original == nil {
        ch := &Change{
            Path:       basePath,
            ChangeType: CHANGE_TYPE_ADD,
            NewValue:   latest,
        }
        l.PushBack(ch)
    } else if latest == nil {
        old := original.(*PhoneNumber)
        ch := &Change{
            Path:          NewPathComponentFromExisting(basePath, NewPathComponentId(old.Id)),
            ChangeType:    CHANGE_TYPE_DELETE,
            OriginalValue: original,
        }
        l.PushBack(ch)
    } else {
        old := original.(*PhoneNumber)
        now := latest.(*PhoneNumber)
        if ch := compareStrings(old.Label, now.Label, l); ch != nil {
            ch.Path = NewPathComponentFromExisting(basePath, NewPathComponentId(old.Id), NewPathComponentKey("label"))
        }
        if ch := compareStrings(string(old.Rel), string(now.Rel), l); ch != nil {
            ch.Path = NewPathComponentFromExisting(basePath, NewPathComponentId(old.Id), NewPathComponentKey("rel"))
        }
        if ch := compareStrings(old.FormattedNumber, now.FormattedNumber, l); ch != nil {
            ch.Path = NewPathComponentFromExisting(basePath, NewPathComponentId(old.Id), NewPathComponentKey("formatted_number"))
        }
        if ch := compareStrings(old.CountryCode, now.CountryCode, l); ch != nil {
            ch.Path = NewPathComponentFromExisting(basePath, NewPathComponentId(old.Id), NewPathComponentKey("country_code"))
        }
        if ch := compareStrings(old.AreaCode, now.AreaCode, l); ch != nil {
            ch.Path = NewPathComponentFromExisting(basePath, NewPathComponentId(old.Id), NewPathComponentKey("area_code"))
        }
        if ch := compareStrings(old.LocalPhoneNumber, now.LocalPhoneNumber, l); ch != nil {
            ch.Path = NewPathComponentFromExisting(basePath, NewPathComponentId(old.Id), NewPathComponentKey("local_phone_number"))
        }
        if ch := compareStrings(old.ExtensionNumber, now.ExtensionNumber, l); ch != nil {
            ch.Path = NewPathComponentFromExisting(basePath, NewPathComponentId(old.Id), NewPathComponentKey("extension_number"))
        }
        if ch := compareBools(old.IsPrimary, now.IsPrimary, l); ch != nil {
            ch.Path = NewPathComponentFromExisting(basePath, NewPathComponentId(old.Id), NewPathComponentKey("is_primary"))
        }
    }
}

func (p *PhoneNumber) IsEmpty() bool {
    return p.FormattedNumber == "" && p.LocalPhoneNumber == "" && p.ExtensionNumber == ""
}

func (p *PhoneNumber) GenerateArrayChanges(original, latest []*PhoneNumber, l *list.List, basePath []*PathComponent) {
    oarr, larr := p.ConvertArraysToDsocialChangers(original, latest)
    compareDsocialChangersArrays(oarr, larr, basePath, l)
}

func (p *PhoneNumber) ConvertArraysToDsocialChangers(original, latest []*PhoneNumber) ([]DsocialChanger, []DsocialChanger) {
    oarr := make([]DsocialChanger, len(original))
    larr := make([]DsocialChanger, len(latest))
    for i, obj := range original {
        oarr[i] = obj
    }
    for i, obj := range latest {
        larr[i] = obj
    }
    return oarr, larr
}

func (p *PhoneNumber) ArraysAreSame(original, latest []*PhoneNumber) bool {
    return arraysAreSame(p.ConvertArraysToDsocialChangers(original, latest))
}

func (p *Email) IsSimilarOrUpdated(original, latest DsocialChanger) (similar bool, same bool) {
    if original == nil && latest == nil {
        return true, true
    }
    if original == nil || latest == nil {
        return
    }
    o, ok := original.(*Email)
    if !ok {
        return
    }
    m, ok := latest.(*Email)
    if !ok {
        return
    }
    if o.Id == m.Id && o.Id != "" {
        similar = true
    } else if o.Id != "" && m.Id != "" {
        return false, false
    } else if checkSimilarStrings(o.EmailAddress, m.EmailAddress) {
        similar = true
    }
    if similar {
        same = o.Label == m.Label &&
            o.Rel == m.Rel &&
            o.IsPrimary == m.IsPrimary
    }
    return
}

func (p *Email) GenerateChanges(original, latest DsocialChanger, basePath []*PathComponent, l *list.List) {
    if original == nil {
        ch := &Change{
            Path:       basePath,
            ChangeType: CHANGE_TYPE_ADD,
            NewValue:   latest,
        }
        l.PushBack(ch)
    } else if latest == nil {
        old := original.(*Email)
        ch := &Change{
            Path:          NewPathComponentFromExisting(basePath, NewPathComponentId(old.Id)),
            ChangeType:    CHANGE_TYPE_DELETE,
            OriginalValue: original,
        }
        l.PushBack(ch)
    } else {
        old := original.(*Email)
        now := latest.(*Email)
        if ch := compareStrings(old.Label, now.Label, l); ch != nil {
            ch.Path = NewPathComponentFromExisting(basePath, NewPathComponentId(old.Id), NewPathComponentKey("label"))
        }
        if ch := compareStrings(string(old.Rel), string(now.Rel), l); ch != nil {
            ch.Path = NewPathComponentFromExisting(basePath, NewPathComponentId(old.Id), NewPathComponentKey("rel"))
        }
        if ch := compareStrings(old.EmailAddress, now.EmailAddress, l); ch != nil {
            ch.Path = NewPathComponentFromExisting(basePath, NewPathComponentId(old.Id), NewPathComponentKey("email_address"))
        }
        if ch := compareBools(old.IsPrimary, now.IsPrimary, l); ch != nil {
            ch.Path = NewPathComponentFromExisting(basePath, NewPathComponentId(old.Id), NewPathComponentKey("is_primary"))
        }
    }
}

func (p *Email) IsEmpty() bool {
    return p.EmailAddress == ""
}

func (p *Email) GenerateArrayChanges(original, latest []*Email, l *list.List, basePath []*PathComponent) {
    oarr, larr := p.ConvertArraysToDsocialChangers(original, latest)
    compareDsocialChangersArrays(oarr, larr, basePath, l)
}

func (p *Email) ConvertArraysToDsocialChangers(original, latest []*Email) ([]DsocialChanger, []DsocialChanger) {
    oarr := make([]DsocialChanger, len(original))
    larr := make([]DsocialChanger, len(latest))
    for i, obj := range original {
        oarr[i] = obj
    }
    for i, obj := range latest {
        larr[i] = obj
    }
    return oarr, larr
}

func (p *Email) ArraysAreSame(original, latest []*Email) bool {
    return arraysAreSame(p.ConvertArraysToDsocialChangers(original, latest))
}

func (p *Uri) IsSimilarOrUpdated(original, latest DsocialChanger) (similar bool, same bool) {
    if original == nil && latest == nil {
        return true, true
    }
    if original == nil || latest == nil {
        return
    }
    o, ok := original.(*Uri)
    if !ok {
        return
    }
    m, ok := latest.(*Uri)
    if !ok {
        return
    }
    if o.Id == m.Id && o.Id != "" {
        similar = true
    } else if o.Id != "" && m.Id != "" {
        return false, false
    } else if checkSimilarStrings(o.Uri, m.Uri) {
        similar = true
    }
    if similar {
        same = o.Label == m.Label &&
            o.Rel == m.Rel &&
            o.IsPrimary == m.IsPrimary
    }
    return
}

func (p *Uri) GenerateChanges(original, latest DsocialChanger, basePath []*PathComponent, l *list.List) {
    if original == nil {
        ch := &Change{
            Path:       basePath,
            ChangeType: CHANGE_TYPE_ADD,
            NewValue:   latest,
        }
        l.PushBack(ch)
    } else if latest == nil {
        old := original.(*Uri)
        ch := &Change{
            Path:          NewPathComponentFromExisting(basePath, NewPathComponentId(old.Id)),
            ChangeType:    CHANGE_TYPE_DELETE,
            OriginalValue: original,
        }
        l.PushBack(ch)
    } else {
        old := original.(*Uri)
        now := latest.(*Uri)
        if ch := compareStrings(old.Label, now.Label, l); ch != nil {
            ch.Path = NewPathComponentFromExisting(basePath, NewPathComponentId(old.Id), NewPathComponentKey("label"))
        }
        if ch := compareStrings(string(old.Rel), string(now.Rel), l); ch != nil {
            ch.Path = NewPathComponentFromExisting(basePath, NewPathComponentId(old.Id), NewPathComponentKey("rel"))
        }
        if ch := compareStrings(old.Uri, now.Uri, l); ch != nil {
            ch.Path = NewPathComponentFromExisting(basePath, NewPathComponentId(old.Id), NewPathComponentKey("uri"))
        }
        if ch := compareBools(old.IsPrimary, now.IsPrimary, l); ch != nil {
            ch.Path = NewPathComponentFromExisting(basePath, NewPathComponentId(old.Id), NewPathComponentKey("is_primary"))
        }
    }
}

func (p *Uri) IsEmpty() bool {
    return p.Uri == ""
}

func (p *Uri) GenerateArrayChanges(original, latest []*Uri, l *list.List, basePath []*PathComponent) {
    oarr, larr := p.ConvertArraysToDsocialChangers(original, latest)
    compareDsocialChangersArrays(oarr, larr, basePath, l)
}

func (p *Uri) ConvertArraysToDsocialChangers(original, latest []*Uri) ([]DsocialChanger, []DsocialChanger) {
    oarr := make([]DsocialChanger, len(original))
    larr := make([]DsocialChanger, len(latest))
    for i, obj := range original {
        oarr[i] = obj
    }
    for i, obj := range latest {
        larr[i] = obj
    }
    return oarr, larr
}

func (p *Uri) ArraysAreSame(original, latest []*Uri) bool {
    return arraysAreSame(p.ConvertArraysToDsocialChangers(original, latest))
}

func (p *IM) IsSimilarOrUpdated(original, latest DsocialChanger) (similar bool, same bool) {
    if original == nil && latest == nil {
        return true, true
    }
    if original == nil || latest == nil {
        return
    }
    o, ok := original.(*IM)
    if !ok {
        return
    }
    m, ok := latest.(*IM)
    if !ok {
        return
    }
    if o.Id == m.Id && o.Id != "" {
        similar = true
    } else if o.Id != "" && m.Id != "" {
        return false, false
    } else if checkSimilarStrings(o.Handle, m.Handle) && checkSimilarStrings(string(o.Protocol), string(m.Protocol)) {
        similar = true
    }
    if similar {
        same = o.Label == m.Label &&
            o.Handle == m.Handle &&
            o.Protocol == m.Protocol &&
            o.Rel == m.Rel &&
            o.IsPrimary == m.IsPrimary
    }
    return
}

func (p *IM) GenerateChanges(original, latest DsocialChanger, basePath []*PathComponent, l *list.List) {
    if original == nil {
        ch := &Change{
            Path:       basePath,
            ChangeType: CHANGE_TYPE_ADD,
            NewValue:   latest,
        }
        l.PushBack(ch)
    } else if latest == nil {
        old := original.(*IM)
        ch := &Change{
            Path:          NewPathComponentFromExisting(basePath, NewPathComponentId(old.Id)),
            ChangeType:    CHANGE_TYPE_DELETE,
            OriginalValue: original,
        }
        l.PushBack(ch)
    } else {
        old := original.(*IM)
        now := latest.(*IM)
        if ch := compareStrings(old.Label, now.Label, l); ch != nil {
            ch.Path = NewPathComponentFromExisting(basePath, NewPathComponentId(old.Id), NewPathComponentKey("label"))
        }
        if ch := compareStrings(string(old.Rel), string(now.Rel), l); ch != nil {
            ch.Path = NewPathComponentFromExisting(basePath, NewPathComponentId(old.Id), NewPathComponentKey("rel"))
        }
        if ch := compareStrings(string(old.Protocol), string(now.Protocol), l); ch != nil {
            ch.Path = NewPathComponentFromExisting(basePath, NewPathComponentId(old.Id), NewPathComponentKey("protocol"))
        }
        if ch := compareStrings(old.Handle, now.Handle, l); ch != nil {
            ch.Path = NewPathComponentFromExisting(basePath, NewPathComponentId(old.Id), NewPathComponentKey("handle"))
        }
        if ch := compareBools(old.IsPrimary, now.IsPrimary, l); ch != nil {
            ch.Path = NewPathComponentFromExisting(basePath, NewPathComponentId(old.Id), NewPathComponentKey("is_primary"))
        }
    }
}

func (p *IM) IsEmpty() bool {
    return p.Handle == "" && p.Protocol == ""
}

func (p *IM) GenerateArrayChanges(original, latest []*IM, l *list.List, basePath []*PathComponent) {
    oarr, larr := p.ConvertArraysToDsocialChangers(original, latest)
    compareDsocialChangersArrays(oarr, larr, basePath, l)
}

func (p *IM) ConvertArraysToDsocialChangers(original, latest []*IM) ([]DsocialChanger, []DsocialChanger) {
    oarr := make([]DsocialChanger, len(original))
    larr := make([]DsocialChanger, len(latest))
    for i, obj := range original {
        oarr[i] = obj
    }
    for i, obj := range latest {
        larr[i] = obj
    }
    return oarr, larr
}

func (p *IM) ArraysAreSame(original, latest []*IM) bool {
    return arraysAreSame(p.ConvertArraysToDsocialChangers(original, latest))
}

func (p *Relationship) IsSimilarOrUpdated(original, latest DsocialChanger) (similar bool, same bool) {
    if original == nil && latest == nil {
        return true, true
    }
    if original == nil || latest == nil {
        return
    }
    o, ok := original.(*Relationship)
    if !ok {
        return
    }
    m, ok := latest.(*Relationship)
    if !ok {
        return
    }
    if o.Id == m.Id && o.Id != "" {
        similar = true
    } else if o.Id != "" && m.Id != "" {
        return false, false
    } else if checkSimilarStrings(o.ContactReferenceId, m.ContactReferenceId) || checkSimilarStrings(o.ContactReferenceName, m.ContactReferenceName) {
        similar = true
    }
    if similar {
        same = o.Label == m.Label &&
            o.Rel == m.Rel &&
            o.ContactReferenceId == m.ContactReferenceId &&
            o.ContactReferenceName == m.ContactReferenceName
    }
    return
}

func (p *Relationship) GenerateChanges(original, latest DsocialChanger, basePath []*PathComponent, l *list.List) {
    if original == nil {
        ch := &Change{
            Path:       basePath,
            ChangeType: CHANGE_TYPE_ADD,
            NewValue:   latest,
        }
        l.PushBack(ch)
    } else if latest == nil {
        old := original.(*Relationship)
        ch := &Change{
            Path:          NewPathComponentFromExisting(basePath, NewPathComponentId(old.Id)),
            ChangeType:    CHANGE_TYPE_DELETE,
            OriginalValue: original,
        }
        l.PushBack(ch)
    } else {
        old := original.(*Relationship)
        now := latest.(*Relationship)
        if ch := compareStrings(old.Label, now.Label, l); ch != nil {
            ch.Path = NewPathComponentFromExisting(basePath, NewPathComponentId(old.Id), NewPathComponentKey("label"))
        }
        if ch := compareStrings(string(old.Rel), string(now.Rel), l); ch != nil {
            ch.Path = NewPathComponentFromExisting(basePath, NewPathComponentId(old.Id), NewPathComponentKey("rel"))
        }
        if ch := compareStrings(old.ContactReferenceId, now.ContactReferenceId, l); ch != nil {
            ch.Path = NewPathComponentFromExisting(basePath, NewPathComponentId(old.Id), NewPathComponentKey("contact_reference_id"))
        }
        if ch := compareStrings(old.ContactReferenceName, now.ContactReferenceName, l); ch != nil {
            ch.Path = NewPathComponentFromExisting(basePath, NewPathComponentId(old.Id), NewPathComponentKey("contact_reference_name"))
        }
    }
}

func (p *Relationship) IsEmpty() bool {
    return p.ContactReferenceId == "" && p.ContactReferenceName == ""
}

func (p *Relationship) GenerateArrayChanges(original, latest []*Relationship, l *list.List, basePath []*PathComponent) {
    oarr, larr := p.ConvertArraysToDsocialChangers(original, latest)
    compareDsocialChangersArrays(oarr, larr, basePath, l)
}

func (p *Relationship) ConvertArraysToDsocialChangers(original, latest []*Relationship) ([]DsocialChanger, []DsocialChanger) {
    oarr := make([]DsocialChanger, len(original))
    larr := make([]DsocialChanger, len(latest))
    for i, obj := range original {
        oarr[i] = obj
    }
    for i, obj := range latest {
        larr[i] = obj
    }
    return oarr, larr
}

func (p *Relationship) ArraysAreSame(original, latest []*Relationship) bool {
    return arraysAreSame(p.ConvertArraysToDsocialChangers(original, latest))
}

func (p *ContactDate) IsSimilarOrUpdated(original, latest DsocialChanger) (similar bool, same bool) {
    if original == nil && latest == nil {
        return true, true
    }
    if original == nil || latest == nil {
        return
    }
    o, ok := original.(*ContactDate)
    if !ok {
        return
    }
    m, ok := latest.(*ContactDate)
    if !ok {
        return
    }
    if o.Id == m.Id && o.Id != "" {
        similar = true
    } else if o.Id != "" && m.Id != "" {
        return false, false
    } else if isSameDate(o.Value, m.Value) || checkSimilarStrings(o.Label, m.Label) || (checkSimilarStrings(string(o.Rel), string(m.Rel)) && o.Rel != REL_DATE_OTHER) {
        similar = true
    }
    if similar {
        same = o.Label == m.Label &&
            o.Rel == m.Rel &&
            o.Value == m.Value &&
            o.IsPrimary == m.IsPrimary
    }
    return
}

func (p *ContactDate) GenerateChanges(original, latest DsocialChanger, basePath []*PathComponent, l *list.List) {
    if original == nil {
        ch := &Change{
            Path:       basePath,
            ChangeType: CHANGE_TYPE_ADD,
            NewValue:   latest,
        }
        l.PushBack(ch)
    } else if latest == nil {
        old := original.(*ContactDate)
        ch := &Change{
            Path:          NewPathComponentFromExisting(basePath, NewPathComponentId(old.Id)),
            ChangeType:    CHANGE_TYPE_DELETE,
            OriginalValue: original,
        }
        l.PushBack(ch)
    } else {
        old := original.(*ContactDate)
        now := latest.(*ContactDate)
        if ch := compareStrings(old.Label, now.Label, l); ch != nil {
            ch.Path = NewPathComponentFromExisting(basePath, NewPathComponentId(old.Id), NewPathComponentKey("label"))
        }
        if ch := compareStrings(string(old.Rel), string(now.Rel), l); ch != nil {
            ch.Path = NewPathComponentFromExisting(basePath, NewPathComponentId(old.Id), NewPathComponentKey("rel"))
        }
        if ch := compareDates(old.Value, now.Value, l); ch != nil {
            ch.Path = NewPathComponentFromExisting(basePath, NewPathComponentId(old.Id), NewPathComponentKey("value"))
        }
        if ch := compareBools(old.IsPrimary, now.IsPrimary, l); ch != nil {
            ch.Path = NewPathComponentFromExisting(basePath, NewPathComponentId(old.Id), NewPathComponentKey("is_primary"))
        }
    }
}

func (p *ContactDate) IsEmpty() bool {
    return p.Value == nil || p.Value.IsEmpty()
}

func (p *ContactDate) GenerateArrayChanges(original, latest []*ContactDate, l *list.List, basePath []*PathComponent) {
    oarr, larr := p.ConvertArraysToDsocialChangers(original, latest)
    compareDsocialChangersArrays(oarr, larr, basePath, l)
}

func (p *ContactDate) ConvertArraysToDsocialChangers(original, latest []*ContactDate) ([]DsocialChanger, []DsocialChanger) {
    oarr := make([]DsocialChanger, len(original))
    larr := make([]DsocialChanger, len(latest))
    for i, obj := range original {
        oarr[i] = obj
    }
    for i, obj := range latest {
        larr[i] = obj
    }
    return oarr, larr
}

func (p *ContactDate) ArraysAreSame(original, latest []*ContactDate) bool {
    return arraysAreSame(p.ConvertArraysToDsocialChangers(original, latest))
}

func (p *ContactDateTime) IsSimilarOrUpdated(original, latest DsocialChanger) (similar bool, same bool) {
    if original == nil && latest == nil {
        return true, true
    }
    if original == nil || latest == nil {
        return
    }
    o, ok := original.(*ContactDateTime)
    if !ok {
        return
    }
    m, ok := latest.(*ContactDateTime)
    if !ok {
        return
    }
    if o.Id == m.Id && o.Id != "" {
        similar = true
    } else if o.Id != "" && m.Id != "" {
        return false, false
    } else if isSameDateTime(o.Value, m.Value) || checkSimilarStrings(o.Label, m.Label) || (checkSimilarStrings(string(o.Rel), string(m.Rel)) && o.Rel != REL_DATETIME_OTHER) {
        similar = true
    }
    if similar {
        same = o.Label == m.Label &&
            o.Rel == m.Rel &&
            o.Value == m.Value &&
            o.IsPrimary == m.IsPrimary
    }
    return
}

func (p *ContactDateTime) GenerateChanges(original, latest DsocialChanger, basePath []*PathComponent, l *list.List) {
    if original == nil {
        ch := &Change{
            Path:       basePath,
            ChangeType: CHANGE_TYPE_ADD,
            NewValue:   latest,
        }
        l.PushBack(ch)
    } else if latest == nil {
        old := original.(*ContactDateTime)
        ch := &Change{
            Path:          NewPathComponentFromExisting(basePath, NewPathComponentId(old.Id)),
            ChangeType:    CHANGE_TYPE_DELETE,
            OriginalValue: original,
        }
        l.PushBack(ch)
    } else {
        old := original.(*ContactDateTime)
        now := latest.(*ContactDateTime)
        if ch := compareStrings(old.Label, now.Label, l); ch != nil {
            ch.Path = NewPathComponentFromExisting(basePath, NewPathComponentId(old.Id), NewPathComponentKey("label"))
        }
        if ch := compareStrings(string(old.Rel), string(now.Rel), l); ch != nil {
            ch.Path = NewPathComponentFromExisting(basePath, NewPathComponentId(old.Id), NewPathComponentKey("rel"))
        }
        if ch := compareDateTimes(old.Value, now.Value, l); ch != nil {
            ch.Path = NewPathComponentFromExisting(basePath, NewPathComponentId(old.Id), NewPathComponentKey("value"))
        }
        if ch := compareBools(old.IsPrimary, now.IsPrimary, l); ch != nil {
            ch.Path = NewPathComponentFromExisting(basePath, NewPathComponentId(old.Id), NewPathComponentKey("is_primary"))
        }
    }
}

func (p *ContactDateTime) IsEmpty() bool {
    return p.Value == nil || p.Value.IsEmpty()
}

func (p *ContactDateTime) GenerateArrayChanges(original, latest []*ContactDateTime, l *list.List, basePath []*PathComponent) {
    oarr, larr := p.ConvertArraysToDsocialChangers(original, latest)
    compareDsocialChangersArrays(oarr, larr, basePath, l)
}

func (p *ContactDateTime) ConvertArraysToDsocialChangers(original, latest []*ContactDateTime) ([]DsocialChanger, []DsocialChanger) {
    oarr := make([]DsocialChanger, len(original))
    larr := make([]DsocialChanger, len(latest))
    for i, obj := range original {
        oarr[i] = obj
    }
    for i, obj := range latest {
        larr[i] = obj
    }
    return oarr, larr
}

func (p *ContactDateTime) ArraysAreSame(original, latest []*ContactDateTime) bool {
    return arraysAreSame(p.ConvertArraysToDsocialChangers(original, latest))
}

func (p *Group) IsSimilarOrUpdated(original, latest DsocialChanger) (similar bool, same bool) {
    if original == nil && latest == nil {
        return true, true
    }
    if original == nil || latest == nil {
        return
    }
    o, ok := original.(*Group)
    if !ok {
        return
    }
    m, ok := latest.(*Group)
    if !ok {
        return
    }
    if o == nil && m == nil {
        return true, true
    }
    if o == nil || m == nil {
        return
    }
    if o.Id == m.Id && o.Id != "" {
        similar = true
    } else if o.Id != "" && m.Id != "" {
        return false, false
    } else if checkSimilarStrings(o.Name, m.Name) {
        similar = true
    }
    if similar {
        same = o.Name == m.Name &&
            o.Description == m.Description &&
            unorderedStringArraysAreSame(o.ContactIds, m.ContactIds) &&
            unorderedStringArraysAreSame(o.ContactNames, m.ContactNames)
    }
    return
}

func (p *Group) GenerateChanges(original, latest DsocialChanger, basePath []*PathComponent, l *list.List) {
    old := original.(*Group)
    now := latest.(*Group)
    if old == nil {
        ch := &Change{
            Path:       basePath,
            ChangeType: CHANGE_TYPE_ADD,
            NewValue:   now,
        }
        l.PushBack(ch)
    } else if now == nil {
        ch := &Change{
            Path:          NewPathComponentFromExisting(basePath, NewPathComponentId(old.Id)),
            ChangeType:    CHANGE_TYPE_DELETE,
            OriginalValue: old,
        }
        l.PushBack(ch)
    } else {
        if ch := compareStrings(old.Name, now.Name, l); ch != nil {
            ch.Path = NewPathComponentFromExisting(basePath, NewPathComponentId(old.Id), NewPathComponentKey("name"))
        }
        if ch := compareStrings(old.Description, now.Description, l); ch != nil {
            ch.Path = NewPathComponentFromExisting(basePath, NewPathComponentId(old.Id), NewPathComponentKey("description"))
        }
        compareUnorderedStringArrays(old.ContactIds, now.ContactIds, NewPathComponentFromExisting(basePath, NewPathComponentId(old.Id), NewPathComponentKey("contact_ids")), l)
        compareUnorderedStringArrays(old.ContactNames, now.ContactNames, NewPathComponentFromExisting(basePath, NewPathComponentId(old.Id), NewPathComponentKey("contact_names")), l)
    }
}

func (p *Group) IsEmpty() bool {
    return p.Name == ""
}

func (p *Group) GenerateArrayChanges(original, latest []*Group, l *list.List, basePath []*PathComponent) {
    oarr, larr := p.ConvertArraysToDsocialChangers(original, latest)
    compareDsocialChangersArrays(oarr, larr, basePath, l)
}

func (p *Group) ConvertArraysToDsocialChangers(original, latest []*Group) ([]DsocialChanger, []DsocialChanger) {
    oarr := make([]DsocialChanger, len(original))
    larr := make([]DsocialChanger, len(latest))
    for i, obj := range original {
        oarr[i] = obj
    }
    for i, obj := range latest {
        larr[i] = obj
    }
    return oarr, larr
}

func (p *Group) ArraysAreSame(original, latest []*Group) bool {
    return arraysAreSame(p.ConvertArraysToDsocialChangers(original, latest))
}

func (p *Certification) IsSimilarOrUpdated(original, latest DsocialChanger) (similar bool, same bool) {
    if original == nil && latest == nil {
        return true, true
    }
    if original == nil || latest == nil {
        return
    }
    o, ok := original.(*Certification)
    if !ok {
        return
    }
    m, ok := latest.(*Certification)
    if !ok {
        return
    }
    if o.Id == m.Id && o.Id != "" {
        similar = true
    } else if o.Id != "" && m.Id != "" {
        return false, false
    } else if checkSimilarStrings(o.Name, m.Name) {
        similar = true
    }
    if similar {
        same = o.Name == m.Name &&
            o.Authority == m.Authority &&
            o.Number == m.Number &&
            isSameDate(o.AsOf, m.AsOf) &&
            isSameDate(o.ValidTill, m.ValidTill)
    }
    return
}

func (p *Certification) GenerateChanges(original, latest DsocialChanger, basePath []*PathComponent, l *list.List) {
    if original == nil {
        ch := &Change{
            Path:       basePath,
            ChangeType: CHANGE_TYPE_ADD,
            NewValue:   latest,
        }
        l.PushBack(ch)
    } else if latest == nil {
        old := original.(*Certification)
        ch := &Change{
            Path:          NewPathComponentFromExisting(basePath, NewPathComponentId(old.Id)),
            ChangeType:    CHANGE_TYPE_DELETE,
            OriginalValue: original,
        }
        l.PushBack(ch)
    } else {
        old := original.(*Certification)
        now := latest.(*Certification)
        if ch := compareStrings(old.Name, now.Name, l); ch != nil {
            ch.Path = NewPathComponentFromExisting(basePath, NewPathComponentId(old.Id), NewPathComponentKey("name"))
        }
        if ch := compareStrings(old.Authority, now.Authority, l); ch != nil {
            ch.Path = NewPathComponentFromExisting(basePath, NewPathComponentId(old.Id), NewPathComponentKey("authority"))
        }
        if ch := compareStrings(old.Number, now.Number, l); ch != nil {
            ch.Path = NewPathComponentFromExisting(basePath, NewPathComponentId(old.Id), NewPathComponentKey("number"))
        }
        if ch := compareDates(old.AsOf, now.AsOf, l); ch != nil {
            ch.Path = NewPathComponentFromExisting(basePath, NewPathComponentId(old.Id), NewPathComponentKey("as_of"))
        }
        if ch := compareDates(old.ValidTill, now.ValidTill, l); ch != nil {
            ch.Path = NewPathComponentFromExisting(basePath, NewPathComponentId(old.Id), NewPathComponentKey("valid_till"))
        }
    }
}

func (p *Certification) IsEmpty() bool {
    return p.Name == ""
}

func (p *Certification) GenerateArrayChanges(original, latest []*Certification, l *list.List, basePath []*PathComponent) {
    oarr, larr := p.ConvertArraysToDsocialChangers(original, latest)
    compareDsocialChangersArrays(oarr, larr, basePath, l)
}

func (p *Certification) ConvertArraysToDsocialChangers(original, latest []*Certification) ([]DsocialChanger, []DsocialChanger) {
    oarr := make([]DsocialChanger, len(original))
    larr := make([]DsocialChanger, len(latest))
    for i, obj := range original {
        oarr[i] = obj
    }
    for i, obj := range latest {
        larr[i] = obj
    }
    return oarr, larr
}

func (p *Certification) ArraysAreSame(original, latest []*Certification) bool {
    return arraysAreSame(p.ConvertArraysToDsocialChangers(original, latest))
}

func (p *Skill) IsSimilarOrUpdated(original, latest DsocialChanger) (similar bool, same bool) {
    if original == nil && latest == nil {
        return true, true
    }
    if original == nil || latest == nil {
        return
    }
    o, ok := original.(*Skill)
    if !ok {
        return
    }
    m, ok := latest.(*Skill)
    if !ok {
        return
    }
    if o.Id == m.Id && o.Id != "" {
        similar = true
    } else if o.Id != "" && m.Id != "" {
        return false, false
    } else if checkSimilarStrings(o.Name, m.Name) {
        similar = true
    }
    if similar {
        same = o.Name == m.Name &&
            o.Proficiency == m.Proficiency
    }
    return
}

func (p *Skill) GenerateChanges(original, latest DsocialChanger, basePath []*PathComponent, l *list.List) {
    if original == nil {
        ch := &Change{
            Path:       basePath,
            ChangeType: CHANGE_TYPE_ADD,
            NewValue:   latest,
        }
        l.PushBack(ch)
    } else if latest == nil {
        old := original.(*Skill)
        ch := &Change{
            Path:          NewPathComponentFromExisting(basePath, NewPathComponentId(old.Id)),
            ChangeType:    CHANGE_TYPE_DELETE,
            OriginalValue: original,
        }
        l.PushBack(ch)
    } else {
        old := original.(*Skill)
        now := latest.(*Skill)
        if ch := compareStrings(old.Name, now.Name, l); ch != nil {
            ch.Path = NewPathComponentFromExisting(basePath, NewPathComponentId(old.Id), NewPathComponentKey("name"))
        }
        if ch := compareStrings(old.Proficiency, now.Proficiency, l); ch != nil {
            ch.Path = NewPathComponentFromExisting(basePath, NewPathComponentId(old.Id), NewPathComponentKey("proficiency"))
        }
    }
}

func (p *Skill) IsEmpty() bool {
    return p.Name == ""
}

func (p *Skill) GenerateArrayChanges(original, latest []*Skill, l *list.List, basePath []*PathComponent) {
    oarr, larr := p.ConvertArraysToDsocialChangers(original, latest)
    compareDsocialChangersArrays(oarr, larr, basePath, l)
}

func (p *Skill) ConvertArraysToDsocialChangers(original, latest []*Skill) ([]DsocialChanger, []DsocialChanger) {
    oarr := make([]DsocialChanger, len(original))
    larr := make([]DsocialChanger, len(latest))
    for i, obj := range original {
        oarr[i] = obj
    }
    for i, obj := range latest {
        larr[i] = obj
    }
    return oarr, larr
}

func (p *Skill) ArraysAreSame(original, latest []*Skill) bool {
    return arraysAreSame(p.ConvertArraysToDsocialChangers(original, latest))
}

func (p *Language) IsSimilarOrUpdated(original, latest DsocialChanger) (similar bool, same bool) {
    if original == nil && latest == nil {
        return true, true
    }
    if original == nil || latest == nil {
        return
    }
    o, ok := original.(*Language)
    if !ok {
        return
    }
    m, ok := latest.(*Language)
    if !ok {
        return
    }
    if o.Id == m.Id && o.Id != "" {
        similar = true
    } else if o.Id != "" && m.Id != "" {
        return false, false
    } else if checkSimilarStrings(o.Name, m.Name) {
        similar = true
    }
    if similar {
        same = o.Name == m.Name &&
            o.ReadGradeLevel == m.ReadGradeLevel &&
            o.WriteGradeLevel == m.WriteGradeLevel
    }
    return
}

func (p *Language) GenerateChanges(original, latest DsocialChanger, basePath []*PathComponent, l *list.List) {
    if original == nil {
        ch := &Change{
            Path:       basePath,
            ChangeType: CHANGE_TYPE_ADD,
            NewValue:   latest,
        }
        l.PushBack(ch)
    } else if latest == nil {
        old := original.(*Language)
        ch := &Change{
            Path:          NewPathComponentFromExisting(basePath, NewPathComponentId(old.Id)),
            ChangeType:    CHANGE_TYPE_DELETE,
            OriginalValue: original,
        }
        l.PushBack(ch)
    } else {
        old := original.(*Language)
        now := latest.(*Language)
        if ch := compareStrings(old.Name, now.Name, l); ch != nil {
            ch.Path = NewPathComponentFromExisting(basePath, NewPathComponentId(old.Id), NewPathComponentKey("name"))
        }
        if ch := compareInts(old.ReadGradeLevel, now.ReadGradeLevel, l); ch != nil {
            ch.Path = NewPathComponentFromExisting(basePath, NewPathComponentId(old.Id), NewPathComponentKey("read_grade_level"))
        }
        if ch := compareInts(old.WriteGradeLevel, now.WriteGradeLevel, l); ch != nil {
            ch.Path = NewPathComponentFromExisting(basePath, NewPathComponentId(old.Id), NewPathComponentKey("write_grade_level"))
        }
    }
}

func (p *Language) IsEmpty() bool {
    return p.Name == ""
}

func (p *Language) GenerateArrayChanges(original, latest []*Language, l *list.List, basePath []*PathComponent) {
    oarr, larr := p.ConvertArraysToDsocialChangers(original, latest)
    compareDsocialChangersArrays(oarr, larr, basePath, l)
}

func (p *Language) ConvertArraysToDsocialChangers(original, latest []*Language) ([]DsocialChanger, []DsocialChanger) {
    oarr := make([]DsocialChanger, len(original))
    larr := make([]DsocialChanger, len(latest))
    for i, obj := range original {
        oarr[i] = obj
    }
    for i, obj := range latest {
        larr[i] = obj
    }
    return oarr, larr
}

func (p *Language) ArraysAreSame(original, latest []*Language) bool {
    return arraysAreSame(p.ConvertArraysToDsocialChangers(original, latest))
}

func (p *Contact) IsSimilarOrUpdated(original, latest DsocialChanger) (similar bool, same bool) {
    if original == nil && latest == nil {
        return true, true
    }
    if original == nil {
        return
    }
    if latest == nil {
        return
    }
    //fmt.Printf("original %#v\n\n", original)
    //fmt.Printf("latest %#v\n\n", latest)
    o, ok := original.(*Contact)
    if !ok {
        return
    }
    m, ok := latest.(*Contact)
    if !ok {
        return
    }
    if o == nil && m == nil {
        return true, true
    }
    if o == nil || m == nil {
        //fmt.Printf("o %#v\n\n", o)
        //fmt.Printf("m %#v\n\n", m)
        return
    }
    if o.Id == m.Id && o.Id != "" {
        similar = true
    } else {
        hasSimilarName := checkSimilarStrings(o.DisplayName, m.DisplayName) ||
            (checkSimilarStrings(o.GivenName, m.GivenName) &&
                (checkSimilarStrings(o.Surname, m.Surname) ||
                    checkSimilarStrings(o.MaidenName, m.MaidenName) ||
                    checkSimilarStrings(o.Surname, m.MaidenName) ||
                    checkSimilarStrings(o.MaidenName, m.Surname)))
        if hasSimilarName && checkSimilarStrings(o.PrimaryAddress, m.PrimaryAddress) {
            similar = true
        } else if hasSimilarName && checkSimilarStrings(o.PrimaryPhoneNumber, m.PrimaryPhoneNumber) {
            similar = true
        } else if hasSimilarName && checkSimilarStrings(o.PrimaryEmail, m.PrimaryEmail) {
            similar = true
        } else if hasSimilarName && checkSimilarStrings(o.PrimaryIm, m.PrimaryIm) {
            similar = true
        } else if hasSimilarName && checkSimilarStrings(o.Title, m.Title) && (checkSimilarStrings(o.Company, m.Company) || checkSimilarStrings(o.Department, m.Department)) {
            similar = true
        } else if hasSimilarName || o.DisplayName == "" || m.DisplayName == "" {
            if !similar && o.EmailAddresses != nil && len(o.EmailAddresses) > 0 && m.EmailAddresses != nil && len(m.EmailAddresses) > 0 {
                e := new(Email)
                for _, old := range o.EmailAddresses {
                    for _, now := range m.EmailAddresses {
                        if similar, _ = e.IsSimilarOrUpdated(old, now); similar {
                            if old.Rel == "" || (old.Rel != REL_EMAIL_OTHER && old.Rel != REL_EMAIL_WORK) || hasSimilarName {
                                break
                            } else {
                                similar = false
                            }
                        }
                    }
                }
            }
            if !similar && o.PhoneNumbers != nil && len(o.PhoneNumbers) > 0 && m.PhoneNumbers != nil && len(m.PhoneNumbers) > 0 {
                ph := new(PhoneNumber)
                for _, old := range o.PhoneNumbers {
                    for _, now := range m.PhoneNumbers {
                        if similar, _ = ph.IsSimilarOrUpdated(old, now); similar {
                            if old.Rel == "" || (old.Rel != REL_PHONE_OTHER && old.Rel != REL_PHONE_WORK) || hasSimilarName {
                                break
                            } else {
                                similar = false
                            }
                        }
                    }
                }
            }
            if !similar && o.PostalAddresses != nil && len(o.PostalAddresses) > 0 && m.PostalAddresses != nil && len(m.PostalAddresses) > 0 {
                ph := new(PostalAddress)
                for _, old := range o.PostalAddresses {
                    for _, now := range m.PostalAddresses {
                        if similar, _ = ph.IsSimilarOrUpdated(old, now); similar {
                            if old.Rel == "" || (old.Rel != REL_ADDRESS_OTHER && old.Rel != REL_ADDRESS_WORK) || hasSimilarName {
                                break
                            } else {
                                similar = false
                            }
                        }
                    }
                }
            }
            if !similar && o.Uris != nil && len(o.Uris) > 0 && m.Uris != nil && len(m.Uris) > 0 {
                u := new(Uri)
                for _, old := range o.Uris {
                    for _, now := range m.Uris {
                        if similar, _ = u.IsSimilarOrUpdated(old, now); similar {
                            if old.Rel == "" || (old.Rel != REL_URI_OTHER && old.Rel != REL_URI_WORK) || hasSimilarName {
                                break
                            } else {
                                similar = false
                            }
                        }
                    }
                }
            }
            if !similar && o.Ims != nil && len(o.Ims) > 0 && m.Ims != nil && len(m.Ims) > 0 {
                im := new(IM)
                for _, old := range o.Ims {
                    for _, now := range m.Ims {
                        if similar, _ = im.IsSimilarOrUpdated(old, now); similar {
                            break
                        }
                    }
                }
            }
        }
    }
    return
}

func (p *Contact) GenerateChanges(original, latest DsocialChanger, basePath []*PathComponent, l *list.List) {
    old := original.(*Contact)
    now := latest.(*Contact)
    if old == nil {
        ch := &Change{
            Path:       basePath,
            ChangeType: CHANGE_TYPE_ADD,
            NewValue:   latest,
        }
        l.PushBack(ch)
    } else if now == nil {
        old := original.(*Contact)
        ch := &Change{
            Path:          NewPathComponentFromExisting(basePath, NewPathComponentId(old.Id)),
            ChangeType:    CHANGE_TYPE_DELETE,
            OriginalValue: original,
        }
        l.PushBack(ch)
    } else {
        compareContactDetails(old, now, l)
        compareContactExternalIds(old, now, l)
        compareContactInternalIds(old, now, l)
        new(GroupRef).GenerateArrayChanges(old.GroupReferences, now.GroupReferences, l, NewPathComponentFromExisting(basePath, NewPathComponentKey("group_references")))
        new(PostalAddress).GenerateArrayChanges(old.PostalAddresses, now.PostalAddresses, l, NewPathComponentFromExisting(basePath, NewPathComponentKey("postal_addresses")))
        new(Education).GenerateArrayChanges(old.Educations, now.Educations, l, NewPathComponentFromExisting(basePath, NewPathComponentKey("educations")))
        new(WorkHistory).GenerateArrayChanges(old.WorkHistories, now.WorkHistories, l, NewPathComponentFromExisting(basePath, NewPathComponentKey("work_histories")))
        new(PhoneNumber).GenerateArrayChanges(old.PhoneNumbers, now.PhoneNumbers, l, NewPathComponentFromExisting(basePath, NewPathComponentKey("phone_numbers")))
        new(Email).GenerateArrayChanges(old.EmailAddresses, now.EmailAddresses, l, NewPathComponentFromExisting(basePath, NewPathComponentKey("email_addresses")))
        new(Uri).GenerateArrayChanges(old.Uris, now.Uris, l, NewPathComponentFromExisting(basePath, NewPathComponentKey("uris")))
        new(IM).GenerateArrayChanges(old.Ims, now.Ims, l, NewPathComponentFromExisting(basePath, NewPathComponentKey("ims")))
        new(Relationship).GenerateArrayChanges(old.Relationships, now.Relationships, l, NewPathComponentFromExisting(basePath, NewPathComponentKey("relationships")))
        new(ContactDate).GenerateArrayChanges(old.Dates, now.Dates, l, NewPathComponentFromExisting(basePath, NewPathComponentKey("dates")))
        new(ContactDateTime).GenerateArrayChanges(old.DateTimes, now.DateTimes, l, NewPathComponentFromExisting(basePath, NewPathComponentKey("datetimes")))
        new(Certification).GenerateArrayChanges(old.Certifications, now.Certifications, l, NewPathComponentFromExisting(basePath, NewPathComponentKey("certifications")))
        new(Skill).GenerateArrayChanges(old.Skills, now.Skills, l, NewPathComponentFromExisting(basePath, NewPathComponentKey("skills")))
        new(Language).GenerateArrayChanges(old.Languages, now.Languages, l, NewPathComponentFromExisting(basePath, NewPathComponentKey("languages")))
    }
}

func (p *Contact) IsEmpty() bool {
    return p.DisplayName == "" &&
        p.GivenName == "" && p.MiddleName == "" && p.Surname == "" && p.Prefix == "" && p.Suffix == "" &&
        p.Company == "" && p.Title == "" &&
        (p.PostalAddresses == nil || len(p.PostalAddresses) == 0) &&
        (p.Educations == nil || len(p.Educations) == 0) &&
        (p.WorkHistories == nil || len(p.WorkHistories) == 0) &&
        (p.PhoneNumbers == nil || len(p.PhoneNumbers) == 0) &&
        (p.EmailAddresses == nil || len(p.EmailAddresses) == 0) &&
        (p.Uris == nil || len(p.Uris) == 0) &&
        (p.Ims == nil || len(p.Ims) == 0) &&
        (p.Relationships == nil || len(p.Relationships) == 0)
}

func orderedStringArraysAreSame(original, latest []string) (same bool) {
    if (original == nil || len(original) == 0) && (latest == nil || len(latest) == 0) {
        same = true
    } else if original != nil && latest != nil && len(original) == len(latest) {
        same = true
        for i, oref := range original {
            if oref != latest[i] {
                same = false
                break
            }
        }
    }
    return
}

func unorderedStringArraysAreSame(old, now []string) (same bool) {
    if old == nil || len(old) == 0 {
        if now == nil || len(now) == 0 {
            return true
        }
        return false
    }
    if now == nil || len(now) == 0 {
        return false
    }
    sort.Strings(old)
    sort.Strings(now)
    lnow := len(now)
    lold := len(old)
    for _, o := range old {
        index := sort.SearchStrings(now, o)
        if index < 0 || index >= lnow || now[index] != o {
            return false
        }
    }
    for _, o := range now {
        index := sort.SearchStrings(old, o)
        if index < 0 || index >= lold || old[index] != o {
            return false
        }
    }
    return true
}

func arraysAreSame(original, latest []DsocialChanger) bool {
    same := true
    if (original == nil || len(original) == 0) && (latest == nil || len(latest) == 0) {
    } else if original != nil && latest != nil && len(original) == len(latest) {
        for _, oref := range original {
            found := false
            for _, mref := range latest {
                if _, refIsSame := oref.IsSimilarOrUpdated(oref, mref); refIsSame {
                    found = true
                    break
                }
            }
            if !found {
                same = false
                break
            }
        }
    } else {
        same = false
    }
    return same
}

func compareStrings(old, now string, l *list.List) (ch *Change) {
    if old == now {
        ch = nil
    } else if old == "" {
        ch = &Change{
            ChangeType: CHANGE_TYPE_ADD,
            NewValue:   now,
        }
    } else if now == "" {
        ch = &Change{
            ChangeType:    CHANGE_TYPE_DELETE,
            OriginalValue: old,
        }
    } else {
        ch = &Change{
            ChangeType:    CHANGE_TYPE_UPDATE,
            OriginalValue: old,
            NewValue:      now,
        }
    }
    if ch != nil {
        l.PushBack(ch)
    }
    return ch
}

func compareInts(old, now int, l *list.List) (ch *Change) {
    if old == now {
    } else if old == 0 {
        ch = &Change{
            ChangeType: CHANGE_TYPE_ADD,
            NewValue:   now,
        }
    } else if now == 0 {
        ch = &Change{
            ChangeType:    CHANGE_TYPE_DELETE,
            OriginalValue: old,
        }
    } else {
        ch = &Change{
            ChangeType:    CHANGE_TYPE_UPDATE,
            OriginalValue: old,
            NewValue:      now,
        }
    }
    if ch != nil {
        l.PushBack(ch)
    }
    return ch
}

func compareInt16s(old, now int16, l *list.List) (ch *Change) {
    if old == now {
    } else if old == 0 {
        ch = &Change{
            ChangeType: CHANGE_TYPE_ADD,
            NewValue:   now,
        }
    } else if now == 0 {
        ch = &Change{
            ChangeType:    CHANGE_TYPE_DELETE,
            OriginalValue: old,
        }
    } else {
        ch = &Change{
            ChangeType:    CHANGE_TYPE_UPDATE,
            OriginalValue: old,
            NewValue:      now,
        }
    }
    if ch != nil {
        l.PushBack(ch)
    }
    return ch
}

func compareFloat64s(old, now float64, l *list.List) (ch *Change) {
    if old == now {
    } else if old == 0 {
        ch = &Change{
            ChangeType: CHANGE_TYPE_ADD,
            NewValue:   now,
        }
    } else if now == 0 {
        ch = &Change{
            ChangeType:    CHANGE_TYPE_DELETE,
            OriginalValue: old,
        }
    } else {
        ch = &Change{
            ChangeType:    CHANGE_TYPE_UPDATE,
            OriginalValue: old,
            NewValue:      now,
        }
    }
    if ch != nil {
        l.PushBack(ch)
    }
    return ch
}

func compareDates(old, now *Date, l *list.List) (ch *Change) {
    if old == now || (old != nil && old.Equals(now)) {
    } else if old == nil {
        ch = &Change{
            ChangeType: CHANGE_TYPE_ADD,
            NewValue:   now,
        }
    } else if now == nil {
        ch = &Change{
            ChangeType:    CHANGE_TYPE_DELETE,
            OriginalValue: old,
        }
    } else {
        ch = &Change{
            ChangeType:    CHANGE_TYPE_UPDATE,
            OriginalValue: old,
            NewValue:      now,
        }
    }
    if ch != nil {
        l.PushBack(ch)
    }
    return ch
}

func compareDateTimes(old, now *DateTime, l *list.List) (ch *Change) {
    if old == now || (old != nil && old.Equals(now)) {
    } else if old == nil {
        ch = &Change{
            ChangeType: CHANGE_TYPE_ADD,
            NewValue:   now,
        }
    } else if now == nil {
        ch = &Change{
            ChangeType:    CHANGE_TYPE_DELETE,
            OriginalValue: old,
        }
    } else {
        ch = &Change{
            ChangeType:    CHANGE_TYPE_UPDATE,
            OriginalValue: old,
            NewValue:      now,
        }
    }
    if ch != nil {
        l.PushBack(ch)
    }
    return ch
}

func compareBools(old, now bool, l *list.List) (ch *Change) {
    if old != now {
        ch = &Change{
            ChangeType:    CHANGE_TYPE_UPDATE,
            OriginalValue: old,
            NewValue:      now,
        }
        l.PushBack(ch)
    }
    return ch
}

func generateDiffContactReferences(old, now *ContactReference, l *list.List, basePath []*PathComponent) (ch *Change) {
    if old == now ||
        (old != nil &&
            now != nil &&
            old.UserId == now.UserId &&
            old.ContactId == now.ContactId &&
            old.ReferenceContactId == now.ReferenceContactId &&
            old.Text == now.Text) {
    } else if old == nil {
        ch = &Change{
            Path:       basePath,
            ChangeType: CHANGE_TYPE_ADD,
            NewValue:   now,
        }
    } else if now == nil {
        ch = &Change{
            Path:          NewPathComponentFromExisting(basePath, NewPathComponentId(old.Id)),
            ChangeType:    CHANGE_TYPE_DELETE,
            OriginalValue: old,
        }
    } else {
        ch = &Change{
            Path:          NewPathComponentFromExisting(basePath, NewPathComponentId(old.Id)),
            ChangeType:    CHANGE_TYPE_UPDATE,
            OriginalValue: old,
            NewValue:      now,
        }
    }
    if ch != nil {
        l.PushBack(ch)
    }
    return ch
}

func compareContactDetails(old, now *Contact, l *list.List) {
    if ch := compareStrings(old.Prefix, now.Prefix, l); ch != nil {
        ch.Path = []*PathComponent{NewPathComponentKey("prefix")}
    }
    if ch := compareStrings(old.GivenName, now.GivenName, l); ch != nil {
        ch.Path = []*PathComponent{NewPathComponentKey("given_name")}
    }
    if ch := compareStrings(old.MiddleName, now.MiddleName, l); ch != nil {
        ch.Path = []*PathComponent{NewPathComponentKey("middle_name")}
    }
    if ch := compareStrings(old.Surname, now.Surname, l); ch != nil {
        ch.Path = []*PathComponent{NewPathComponentKey("surname")}
    }
    if ch := compareStrings(old.Suffix, now.Suffix, l); ch != nil {
        ch.Path = []*PathComponent{NewPathComponentKey("suffix")}
    }
    if ch := compareStrings(old.MaidenName, now.MaidenName, l); ch != nil {
        ch.Path = []*PathComponent{NewPathComponentKey("maiden_name")}
    }
    if ch := compareStrings(old.DisplayName, now.DisplayName, l); ch != nil {
        ch.Path = []*PathComponent{NewPathComponentKey("display_name")}
    }
    if ch := compareStrings(old.Nickname, now.Nickname, l); ch != nil {
        ch.Path = []*PathComponent{NewPathComponentKey("nickname")}
    }
    if ch := compareStrings(string(old.DisplayNameOrdering), string(now.DisplayNameOrdering), l); ch != nil {
        ch.Path = []*PathComponent{NewPathComponentKey("display_name_ordering")}
    }
    if ch := compareStrings(string(old.SortNameOrdering), string(now.SortNameOrdering), l); ch != nil {
        ch.Path = []*PathComponent{NewPathComponentKey("sort_name_ordering")}
    }
    if ch := compareStrings(old.Hometown, now.Hometown, l); ch != nil {
        ch.Path = []*PathComponent{NewPathComponentKey("hometown")}
    }
    if ch := compareStrings(string(old.Gender), string(now.Gender), l); ch != nil {
        ch.Path = []*PathComponent{NewPathComponentKey("gender")}
    }
    if ch := compareStrings(old.Biography, now.Biography, l); ch != nil {
        ch.Path = []*PathComponent{NewPathComponentKey("biography")}
    }
    if ch := compareStrings(old.FavoriteQuotes, now.FavoriteQuotes, l); ch != nil {
        ch.Path = []*PathComponent{NewPathComponentKey("favorite_quotes")}
    }
    if ch := compareStrings(string(old.RelationshipStatus), string(now.RelationshipStatus), l); ch != nil {
        ch.Path = []*PathComponent{NewPathComponentKey("relationship_status")}
    }
    if ch := compareStrings(old.Title, now.Title, l); ch != nil {
        ch.Path = []*PathComponent{NewPathComponentKey("title")}
    }
    if ch := compareStrings(old.Company, now.Company, l); ch != nil {
        ch.Path = []*PathComponent{NewPathComponentKey("company")}
    }
    if ch := compareStrings(old.Department, now.Department, l); ch != nil {
        ch.Path = []*PathComponent{NewPathComponentKey("department")}
    }
    if ch := compareStrings(old.Municipality, now.Municipality, l); ch != nil {
        ch.Path = []*PathComponent{NewPathComponentKey("municipality")}
    }
    if ch := compareStrings(old.Region, now.Region, l); ch != nil {
        ch.Path = []*PathComponent{NewPathComponentKey("region")}
    }
    if ch := compareStrings(old.PostalCode, now.PostalCode, l); ch != nil {
        ch.Path = []*PathComponent{NewPathComponentKey("postal_code")}
    }
    if ch := compareStrings(old.CountryCode, now.CountryCode, l); ch != nil {
        ch.Path = []*PathComponent{NewPathComponentKey("country_code")}
    }
    if ch := compareStrings(old.PrimaryAddress, now.PrimaryAddress, l); ch != nil {
        ch.Path = []*PathComponent{NewPathComponentKey("primary_address")}
    }
    if ch := compareStrings(old.PrimaryEmail, now.PrimaryEmail, l); ch != nil {
        ch.Path = []*PathComponent{NewPathComponentKey("primary_email")}
    }
    if ch := compareStrings(old.PrimaryPhoneNumber, now.PrimaryPhoneNumber, l); ch != nil {
        ch.Path = []*PathComponent{NewPathComponentKey("primary_phone_number")}
    }
    if ch := compareStrings(old.PrimaryUri, now.PrimaryUri, l); ch != nil {
        ch.Path = []*PathComponent{NewPathComponentKey("primary_uri")}
    }
    if ch := compareStrings(old.PrimaryIm, now.PrimaryIm, l); ch != nil {
        ch.Path = []*PathComponent{NewPathComponentKey("primary_im")}
    }
    if ch := compareStrings(old.Notes, now.Notes, l); ch != nil {
        ch.Path = []*PathComponent{NewPathComponentKey("notes")}
    }
    if ch := compareStrings(old.ThumbnailUrl, now.ThumbnailUrl, l); ch != nil {
        ch.Path = []*PathComponent{NewPathComponentKey("thumbnail_url")}
    }
}

func compareUnorderedStringArrays(old, now []string, path []*PathComponent, l *list.List) {
    if old == nil || len(old) == 0 {
        if now == nil || len(now) == 0 {
            return
        }
        for _, o := range now {
            l.PushBack(&Change{
                Path:       path,
                ChangeType: CHANGE_TYPE_ADD,
                NewValue:   o,
            })
        }
        return
    }
    if now == nil || len(now) == 0 {
        for _, o := range now {
            l.PushBack(&Change{
                Path:          path,
                ChangeType:    CHANGE_TYPE_DELETE,
                OriginalValue: o,
            })
        }
        return
    }
    sort.Strings(old)
    sort.Strings(now)
    lnow := len(now)
    lold := len(old)
    for _, o := range old {
        index := sort.SearchStrings(now, o)
        if index < 0 || index >= lnow || now[index] != o {
            l.PushBack(&Change{
                Path:          path,
                ChangeType:    CHANGE_TYPE_DELETE,
                OriginalValue: o,
            })
        }
    }
    for _, o := range now {
        index := sort.SearchStrings(old, o)
        if index < 0 || index >= lold || old[index] != o {
            l.PushBack(&Change{
                Path:       path,
                ChangeType: CHANGE_TYPE_ADD,
                NewValue:   o,
            })
        }
    }
}

func compareContactInternalIds(old, now *Contact, l *list.List) {
    compareUnorderedStringArrays(old.InternalUserIds, now.InternalUserIds, []*PathComponent{NewPathComponentKey("internal_ids")}, l)
}

func compareContactExternalIds(old, now *Contact, l *list.List) {
    compareUnorderedStringArrays(old.ExternalUserIds, now.ExternalUserIds, []*PathComponent{NewPathComponentKey("external_ids")}, l)
}

func isSameDate(old, now *Date) bool {
    if old == nil {
        if now == nil {
            return true
        }
        return now.IsEmpty()
    } else if now == nil {
        return old.IsEmpty()
    }
    return old.Equals(now)
}

func isSameDateTime(old, now *DateTime) bool {
    if old == nil {
        if now == nil {
            return true
        }
        return now.IsEmpty()
    } else if now == nil {
        return old.IsEmpty()
    }
    return old.Equals(now)
}

func isSimilarPostalAddress(old, now *PostalAddress) bool {
    if old == nil || now == nil {
        return false
    }
    if old.Address != "" && old.Address == now.Address {
        return true
    }
    if join("", now.StreetAddress, now.OtherAddress, now.Municipality, now.Region, now.PostalCode, now.Country) != "" &&
        (old.StreetAddress == "" || old.StreetAddress == now.StreetAddress) &&
        (old.OtherAddress == "" || old.OtherAddress == now.OtherAddress) &&
        (old.Municipality == "" || old.Municipality == now.Municipality) &&
        (old.Region == "" || old.Region == now.Region) &&
        (old.PostalCode == "" || old.PostalCode == now.PostalCode) &&
        (old.Country == "" || old.Country == now.Country) {
        return true
    }
    if old.Label == "" {
        if now.Label != "" {
            return false
        }
        if old.Rel != now.Rel {
            return false
        }
    } else if now.Label != old.Label {
        return false
    }
    if old.LocatedFrom != nil && now.LocatedFrom != nil {
        if old.LocatedFrom.Year == now.LocatedFrom.Year &&
            old.LocatedFrom.Month == now.LocatedFrom.Month &&
            old.LocatedFrom.Day == now.LocatedFrom.Day {
            return true
        }
    }
    if old.LocatedTill != nil && now.LocatedTill != nil {
        if old.LocatedTill.Year == now.LocatedTill.Year &&
            old.LocatedTill.Month == now.LocatedTill.Month &&
            old.LocatedTill.Day == now.LocatedTill.Day {
            return true
        }
    }
    return false
}

func compareDsocialChangersArrays(old, now []DsocialChanger, basepath []*PathComponent, l *list.List) {
    if old == nil || len(old) == 0 {
        if now == nil || len(now) == 0 {
            return
        }
        for _, n := range now {
            if n != nil {
                n.GenerateChanges(nil, n, basepath, l)
            }
        }
        return
    }
    if now == nil || len(now) == 0 {
        for _, n := range old {
            if n != nil {
                n.GenerateChanges(n, nil, basepath, l)
            }
        }
        return
    }
    loFound := make([]bool, len(old))
    lnFound := make([]bool, len(now))
    for i, o := range old {
        if o != nil {
            for j, n := range now {
                isSimilar, isSame := o.IsSimilarOrUpdated(o, n)
                if !lnFound[j] && isSimilar {
                    if !isSame {
                        o.GenerateChanges(o, n, basepath, l)
                    }
                    loFound[i] = true
                    lnFound[j] = true
                    break
                }
            }
        }
    }
    for i, o := range old {
        if !loFound[i] && o != nil {
            o.GenerateChanges(o, nil, basepath, l)
        }
    }
    for i, n := range now {
        if !lnFound[i] && n != nil {
            n.GenerateChanges(nil, n, basepath, l)
        }
    }
}

func checkSimilarStrings(a, b string) bool {
    return a != "" && b != "" && (a == b || strings.ToLower(strings.TrimSpace(a)) == strings.ToLower(strings.TrimSpace(b)))
}
