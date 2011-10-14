package dsocial

import (
    "github.com/pomack/contacts.go/facebook"
    "container/list"
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

func fbEducationToDsocialEducation(e *facebook.Education) *Education {
    if e == nil || len(e.Type) == 0 {
        return nil
    }
    ed := new(Education)
    sType := strings.ToLower(e.Type)
    switch sType {
    case "elementary school":
        ed.Rel = REL_EDUCATION_ELEMENTARY_SCHOOL
    case "middle school":
        ed.Rel = REL_EDUCATION_MIDDLE_SCHOOL
    case "high school":
        ed.Rel = REL_EDUCATION_HIGH_SCHOOL
    case "college":
        ed.Rel = REL_EDUCATION_COLLEGE
    case "graduate school":
        ed.Rel = REL_EDUCATION_GRADUATE_SCHOOL
    case "vocational":
        ed.Rel = REL_EDUCATION_VOCATIONAL
    default:
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
    return ed
}

func fbWorkPlaceToDsocialWorkHistory(w *facebook.WorkPlace) *WorkHistory {
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
    return wh
}

func FbContactToDsocial(fbContact *facebook.Contact) *Contact {
    if fbContact == nil {
        return nil
    }
    c := new(Contact)
    c.DisplayName = fbContact.Name
    c.GivenName = fbContact.FirstName
    c.MiddleName = fbContact.MiddleName
    c.Surname = fbContact.LastName
    c.DisplayNameOrdering = GIVEN_MIDDLE_SURNAME
    c.SortNameOrdering = SURNAME_GIVEN_MIDDLE
    urlList := list.New()
    if len(fbContact.Link) > 0 {
        fbUri := new(Uri)
        fbUri.Uri = fbContact.Link
        fbUri.Rel = REL_URI_FACEBOOK
        urlList.PushBack(fbUri)
    }
    if len(fbContact.Website) > 0 {
        ws := new(Uri)
        ws.Rel = REL_URI_WEBSITE
        ws.Uri = fbContact.Website
        urlList.PushBack(ws)
    }
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
    gender := RelGender(strings.ToLower(fbContact.Gender))
    switch gender {
    case "male":
        c.Gender = REL_GENDER_MALE
    case "female":
        c.Gender = REL_GENDER_FEMALE
    }
    educations := make([]*Education, len(fbContact.Educations))
    for i, e := range fbContact.Educations {
        educations[i] = fbEducationToDsocialEducation(&e)
    }
    c.Educations = educations
    workplaces := make([]*WorkHistory, len(fbContact.WorkPlaces))
    for i, w := range fbContact.WorkPlaces {
        workplaces[i] = fbWorkPlaceToDsocialWorkHistory(&w)
    }
    c.WorkHistories = workplaces
    relationshipStatus := strings.ToLower(fbContact.RelationshipStatus)
    var useRelationshipStatus RelRelationshipStatus
    switch relationshipStatus {
    case "single":
        useRelationshipStatus = REL_SINGLE
    case "in a relationship":
        useRelationshipStatus = REL_IN_A_RELATIONSHIP
    case "engaged":
        useRelationshipStatus = REL_ENGAGED
    case "married":
        useRelationshipStatus = REL_MARRIED
    case "it's complicated":
        useRelationshipStatus = REL_ITS_COMPLICATED
    case "in an open relationship":
        useRelationshipStatus = REL_OPEN_RELATIONSHIP
    case "widowed":
        useRelationshipStatus = REL_WIDOWED
    case "separated":
        useRelationshipStatus = REL_SEPARATED
    case "divorced":
        useRelationshipStatus = REL_DIVORCED
    case "in a civil union":
        useRelationshipStatus = REL_IN_CIVIL_UNION
    case "in a domestic partnership":
        useRelationshipStatus = REL_IN_DOMESTIC_PARTNERSHIP
    }
    if len(useRelationshipStatus) > 0 {
        c.RelationshipStatus = useRelationshipStatus
    }
    if len(fbContact.SignificantOther.Name) > 0 {
        relationship := new(Relationship)
        relationship.Rel = REL_RELATIONSHIP_SPOUSE
        relationship.ContactReferenceName = fbContact.SignificantOther.Name
        c.Relationships = []*Relationship{relationship}
    }
    return c
}
