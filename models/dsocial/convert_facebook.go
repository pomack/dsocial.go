package dsocial

import (
    "github.com/pomack/contacts.go/facebook"
    "container/list"
    "fmt"
    "strings"
    "strconv"
)

func fbDateToDsocialDate(fbDate string) *Date {
    // known formats
    // YYY-MM
    // MM/DD/YYYY
    // MM/DD
    var theDate *Date
    separator := "/"
    if strings.Contains(fbDate, "-") {
        separator = "-"
    } else if strings.Contains(fbDate, ".") {
        separator = "."
    } else if strings.Contains(fbDate, "_") {
        separator = "_"
    }
    parts := strings.Split(fbDate, separator)
    year, month, day := 0, 0, 0
    switch len(parts) {
    case 3:
        month, _ = strconv.Atoi(parts[0])
        day, _ = strconv.Atoi(parts[1])
        year, _ = strconv.Atoi(parts[2])
    case 2:
        month, _ = strconv.Atoi(parts[0])
        day, _ = strconv.Atoi(parts[1])
        if month > 12 || month < 0 {
            year = month
            month = day
        } else if day > 31 || day < 0 {
            year = day
            day = 0
        }
    case 1:
        month, _ = strconv.Atoi(parts[0])
        if month > 12 || month < 0 {
            year = month
            month = 0
        }
    }
    if month != 0 || day != 0 || year != 0 {
        theDate = new(Date)
        theDate.Month = int8(month)
        theDate.Day = int8(day)
        theDate.Year = int16(year)
    }
    return theDate
}

func fbEducationToDsocialEducation(e *facebook.Education, original []*Education, dsocialUserId string) *Education {
    if e == nil || len(e.Type) == 0 {
        return nil
    }
    ed := new(Education)
    ed.Rel, _ = facebookToEducationRelMap[strings.ToLower(e.Type)]
    if len(ed.Rel) == 0 {
        ed.Label = e.Type
    }
    ed.Institution = e.School.Name
    if len(e.Year.Name) > 0 {
        year, _ := strconv.Atoi(e.Year.Name)
        ed.GraduationYear = int16(year)
    }
    if len(e.Concentrations) > 0 {
        degrees := make([]*Degree, len(e.Concentrations))
        for i, c := range e.Concentrations {
            degree := new(Degree)
            degree.Degree = e.Degree.Name
            degree.Major = c.Name
            degrees[i] = degree
        }
        ed.Degrees = degrees
    } else if len(e.Degree.Name) > 0 {
        degrees := make([]*Degree, 1)
        degree := new(Degree)
        degree.Degree = e.Degree.Name
        degrees[0] = degree
    }
    id := ""
    if original != nil {
        for _, o := range original {
            if isSimilar, _ := ed.IsSimilarOrUpdated(ed, o); isSimilar {
                id = o.Id
            }
        }
    }
    ed.Id = id
    ed.Acl.OwnerId = dsocialUserId
    return ed
}

func fbWorkPlaceToDsocialWorkHistory(w *facebook.WorkPlace, original []*WorkHistory, dsocialUserId string) *WorkHistory {
    if w == nil || len(w.Employer.Name) <= 0 {
        return nil
    }
    wh := new(WorkHistory)
    wp := new(WorkPosition)
    wh.Company = w.Employer.Name
    wp.Title = w.Position.Name
    wp.Description = w.Description
    from := fbDateToDsocialDate(w.StartDate)
    to := fbDateToDsocialDate(w.EndDate)
    wp.From = from
    wp.To = to
    if from != nil && to == nil {
        wp.IsCurrent = true
    }
    wp.Location = w.Location.Name
    references := make([]*ContactReference, len(w.With))
    for i, r := range w.With {
        ref := new(ContactReference)
        ref.ReferenceContactName = r.Name
        references[i] = ref
    }
    wp.References = references
    wh.Positions = []*WorkPosition{wp}
    id := ""
    if original != nil {
        for _, o := range original {
            if isSimilar, _ := wh.IsSimilarOrUpdated(wh, o); isSimilar {
                id = o.Id
                if o.Positions != nil {
                    for _, o2 := range o.Positions {
                        if isSimilar2, _ := wp.IsSimilarOrUpdated(wp, o2); isSimilar2 {
                            wp.Id = o2.Id
                            break
                        }
                    }
                }
            }
        }
    }
    wh.Id = id
    wh.Acl.OwnerId = dsocialUserId
    wp.Acl.OwnerId = dsocialUserId
    return wh
}

func fbUriToDsocial(link string, rel RelUri, original []*Uri, dsocialUserId string, uriList *list.List) {
    if len(link) > 0 {
        fbUri := new(Uri)
        fbUri.Uri = link
        fbUri.Rel = rel
        if original != nil && len(original) > 0 {
            for _, u := range original {
                if u.Uri == link {
                    fbUri.Id = u.Id
                    break
                }
            }
        }
        fbUri.Acl.OwnerId = dsocialUserId
        uriList.PushBack(fbUri)
    }
}

func fbRelationshipToDsocial(name string, rel RelRelationship, original []*Relationship, dsocialUserId string, relationshipList *list.List) {
    if len(name) > 0 {
        fbRelationship := new(Relationship)
        fbRelationship.ContactReferenceName = name
        fbRelationship.Rel = rel
        if original != nil && len(original) > 0 {
            for _, r := range original {
                if r.ContactReferenceName == name {
                    fbRelationship.Id = r.Id
                    break
                }
            }
        }
        fbRelationship.Acl.OwnerId = dsocialUserId
        relationshipList.PushBack(fbRelationship)
    }
}

func FacebookContactToDsocial(fbContact *facebook.Contact, o *Contact, dsocialUserId string) *Contact {
    if fbContact == nil {
        return nil
    }
    c := new(Contact)
    if o != nil {
        c.Id = o.Id
        c.UserId = o.UserId
    } else {
        o = new(Contact)
    }
    c.UserId = dsocialUserId
    c.Acl.OwnerId = dsocialUserId
    c.DisplayName = fbContact.Name
    c.GivenName = fbContact.FirstName
    c.MiddleName = fbContact.MiddleName
    c.Surname = fbContact.LastName
    c.DisplayNameOrdering = GIVEN_MIDDLE_SURNAME
    c.SortNameOrdering = SURNAME_GIVEN_MIDDLE
    urlList := list.New()
    fbUriToDsocial(fbContact.Link, REL_URI_FACEBOOK, o.Uris, dsocialUserId, urlList)
    fbUriToDsocial(fbContact.Website, REL_URI_WEBSITE, o.Uris, dsocialUserId, urlList)
    uris := make([]*Uri, urlList.Len())
    uriIndex := 0
    for e := urlList.Front(); e != nil; e = e.Next() {
        uris[uriIndex] = e.Value.(*Uri)
        uriIndex++
    }
    c.Uris = uris
    c.Biography = fbContact.Bio
    c.FavoriteQuotes = fbContact.Quotes
    c.Birthday = fbDateToDsocialDate(fbContact.Birthday)
    c.Hometown = fbContact.Hometown.Name
    gender, _ := facebookToGenderMap[strings.ToLower(fbContact.Gender)]
    if len(gender) > 0 {
        c.Gender = gender
    }
    educations := make([]*Education, len(fbContact.Educations))
    for i, e := range fbContact.Educations {
        educations[i] = fbEducationToDsocialEducation(&e, o.Educations, dsocialUserId)
    }
    c.Educations = educations
    workplaces := make([]*WorkHistory, len(fbContact.WorkPlaces))
    for i, w := range fbContact.WorkPlaces {
        workplaces[i] = fbWorkPlaceToDsocialWorkHistory(&w, o.WorkHistories, dsocialUserId)
    }
    c.WorkHistories = workplaces
    relationshipStatus, _ := facebookToRelationshipStatusMap[strings.ToLower(fbContact.RelationshipStatus)]
    if len(relationshipStatus) > 0 {
        c.RelationshipStatus = relationshipStatus
    }
    relationshipList := list.New()
    fbRelationshipToDsocial(fbContact.SignificantOther.Name, REL_RELATIONSHIP_SPOUSE, o.Relationships, dsocialUserId, relationshipList)
    c.Relationships = make([]*Relationship, relationshipList.Len())
    for i, iter := 0, relationshipList.Front(); iter != nil; i, iter = i + 1, iter.Next() {
        c.Relationships[i] = iter.Value.(*Relationship)
    }
    return c
}




func dsocialDateToFacebook(thedate *Date) string {
    // known formats
    // YYY-MM
    // MM/DD/YYYY
    // MM/DD
    if thedate == nil {
        return ""
    }
    var s string
    if thedate.Year > 0 {
        if thedate.Month > 0 {
            if thedate.Day > 0 {
                s = fmt.Sprintf("%02d/%02d/%04d", thedate.Month, thedate.Day, thedate.Year)
            } else {
                s = fmt.Sprintf("%04d-%02d", thedate.Year, thedate.Month)
            }
        } else {
            s = fmt.Sprintf("%04d", thedate.Year)
        }
    } else if thedate.Month > 0 && thedate.Day > 0 {
        s = fmt.Sprintf("%02d/%02d", thedate.Month, thedate.Day)
    }
    return s
}

func dsocialEducationToFacebook(c *Education, f *facebook.Education) {
    if c == nil || f == nil {
        return
    }
    f.Type, _ = educationRelToFacebookMap[c.Rel]
    if len(f.Type) == 0 {
        f.Type = c.Label
    }
    f.School.Name = c.Institution
    if c.GraduationYear != 0 {
        f.Year.Name = strconv.Itoa(int(c.GraduationYear))
    }
    if c.Degrees != nil {
        f.Concentrations = make([]facebook.Concentration, len(c.Degrees))
        for i, d := range c.Degrees {
            f.Degree.Name = d.Degree
            f.Concentrations[i].Name = d.Major
        }
    }
    return
}

func dsocialWorkPositionToFacebook(wh *WorkHistory, wp *WorkPosition, w *facebook.WorkPlace) {
    if wh == nil || w == nil {
        return
    }
    w.Employer.Name = wh.Company
    w.Description = wh.Description
    if wp != nil {
        w.Position.Name = wp.Title
        w.Description = wp.Description
        w.StartDate = dsocialDateToFacebook(wp.From)
        if !wp.IsCurrent {
            w.EndDate = dsocialDateToFacebook(wp.To)
        }
        w.Location.Name = wp.Location
        if wp.References != nil {
            w.With = make([]facebook.ContactReference, len(wp.References))
            for i, r := range wp.References {
                w.With[i].Name = r.ReferenceContactName
            }
        }
    }
    return
}



func DsocialContactToFacebook(c *Contact, o *facebook.Contact) *facebook.Contact {
    if c == nil {
        return nil
    }
    f := new(facebook.Contact)
    if o != nil {
        f.Id = o.Id
    } else {
        o = new(facebook.Contact)
    }
    f.Name = c.DisplayName
    f.FirstName = c.GivenName
    f.MiddleName = c.MiddleName
    f.LastName = c.Surname
    f.Bio = c.Biography
    f.Quotes = c.FavoriteQuotes
    f.Birthday = dsocialDateToFacebook(c.Birthday)
    f.Hometown.Name = c.Hometown
    f.Gender, _ = genderToFacebookMap[c.Gender]
    f.RelationshipStatus, _ = relationshipStatusToFacebookMap[c.RelationshipStatus]
    if c.Relationships != nil {
        for _, r := range c.Relationships {
            switch r.Rel {
            case REL_RELATIONSHIP_SPOUSE:
                f.SignificantOther.Name = r.ContactReferenceName
            }
        }
    }
    if c.Uris != nil {
        for _, u := range c.Uris {
            switch u.Rel {
            case REL_URI_FACEBOOK:
                f.Link = u.Uri
            case REL_URI_WEBSITE:
                f.Website = u.Uri
            }
        }
    }
    if c.Educations != nil {
        f.Educations = make([]facebook.Education, len(c.Educations))
        for i, e := range c.Educations {
            dsocialEducationToFacebook(e, &f.Educations[i])
        }
    }
    if c.WorkHistories != nil {
        numWorkPlaces := 0
        for _, wh := range c.WorkHistories {
            if wh.Positions == nil || len(wh.Positions) == 0 {
                numWorkPlaces += 1
            } else {
                numWorkPlaces += len(wh.Positions)
            }
        }
        f.WorkPlaces = make([]facebook.WorkPlace, numWorkPlaces)
        i := 0
        for _, wh := range c.WorkHistories {
            if wh.Positions == nil || len(wh.Positions) == 0 {
                dsocialWorkPositionToFacebook(wh, nil, &f.WorkPlaces[i])
                i++
            } else {
                for _, wp := range wh.Positions {
                    dsocialWorkPositionToFacebook(wh, wp, &f.WorkPlaces[i])
                    i++
                }
            }
        }
    }
    return f
}
