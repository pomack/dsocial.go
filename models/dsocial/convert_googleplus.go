package dsocial

import (
    "container/list"
    "github.com/pomack/contacts.go/googleplus"
)

func GooglePlusPersonToDsocial(g *googleplus.Person, o *Contact, dsocialUserId string) *Contact {
    if g == nil {
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
    c.Birthday = googleDateStringToDsocial(g.Birthday)
    c.Biography = g.AboutMe
    c.FavoriteQuotes = g.Tagline
    c.Nickname = g.Nickname
    name := join(" ", g.Name.HonorificPrefix, g.Name.GivenName, g.Name.MiddleName, g.Name.FamilyName, g.Name.HonorificSuffix)
    if name == "" {
        ParseName(g.DisplayName, c)
    } else {
        c.DisplayName = g.DisplayName
        c.Prefix = g.Name.HonorificPrefix
        c.GivenName = g.Name.GivenName
        c.MiddleName = g.Name.MiddleName
        c.Surname = g.Name.FamilyName
        c.Suffix = g.Name.HonorificSuffix
    }
    switch g.Gender {
    case googleplus.GENDER_MALE:
        c.Gender = REL_GENDER_MALE
    case googleplus.GENDER_FEMALE:
        c.Gender = REL_GENDER_FEMALE
    case googleplus.GENDER_OTHER:
        c.Gender = REL_GENDER_OTHER
    }
    c.Uris = make([]*Uri, len(g.Urls))
    for i, uri := range g.Urls {
        c.Uris[i] = googleplusUriToDsocial(&uri, o.Uris, dsocialUserId)
    }
    c.PostalAddresses = make([]*PostalAddress, len(g.PlacesLived))
    for i, addr := range g.PlacesLived {
        c.PostalAddresses[i] = &PostalAddress{
            Address:   addr.Value,
            Rel:       REL_ADDRESS_OTHER,
            IsPrimary: addr.Primary,
        }
    }
    educations := list.New()
    workhistories := list.New()
    for _, org := range g.Organizations {
        switch org.Type {
        case googleplus.ORGANIZATION_TYPE_SCHOOL:
            from := googleDateStringToDsocial(org.StartDate)
            to := googleDateStringToDsocial(org.EndDate)
            isCurrent := (from != nil && to == nil) || org.Primary
            educations.PushBack(&Education{
                Label:        org.Title,
                Institution:  org.Name,
                AttendedFrom: from,
                AttendedTill: to,
                IsCurrent:    isCurrent,
            })
        case googleplus.ORGANIZATION_TYPE_WORK:
            from := googleDateStringToDsocial(org.StartDate)
            to := googleDateStringToDsocial(org.EndDate)
            isCurrent := (from != nil && to == nil) || org.Primary
            workhistories.PushBack(&WorkHistory{
                Positions: []*WorkPosition{&WorkPosition{
                    Title:       org.Title,
                    Description: org.Description,
                    From:        from,
                    To:          to,
                    IsCurrent:   isCurrent,
                    Location:    org.Location,
                    Department:  org.Department,
                }},
                Company: org.Name,
            })
        }
    }
    c.Educations = make([]*Education, educations.Len())
    for i, e := 0, educations.Front(); e != nil; i, e = i+1, e.Next() {
        c.Educations[i] = e.Value.(*Education)
    }
    c.WorkHistories = make([]*WorkHistory, workhistories.Len())
    for i, e := 0, workhistories.Front(); e != nil; i, e = i+1, e.Next() {
        c.WorkHistories[i] = e.Value.(*WorkHistory)
    }
    c.Languages = make([]*Language, len(g.LanguagesSpoken))
    for i, lang := range g.LanguagesSpoken {
        c.Languages[i] = &Language{Name: lang}
    }
    c.EmailAddresses = make([]*Email, len(g.Emails))
    for i, email := range g.Emails {
        c.EmailAddresses[i] = googleplusEmailToDsocial(&email, o.EmailAddresses, dsocialUserId)
    }
    return c
}

func googleplusUriToDsocial(g *googleplus.Url, original []*Uri, dsocialUserId string) *Uri {
    if g == nil {
        return nil
    }
    var rel RelUri
    switch g.Type {
    case googleplus.URL_TYPE_HOME:
        rel = REL_URI_HOME
    case googleplus.URL_TYPE_WORK:
        rel = REL_URI_WORK
    case googleplus.URL_TYPE_BLOG:
        rel = REL_URI_BLOG
    case googleplus.URL_TYPE_PROFILE:
        rel = REL_URI_GOOGLE_PLUS
    case googleplus.URL_TYPE_OTHER:
        rel = REL_URI_OTHER
    default:
        rel = REL_URI_OTHER
    }
    u := &Uri{Rel: rel, Uri: g.Value, IsPrimary: g.Primary}
    id := ""
    if original != nil {
        for _, o := range original {
            if o != nil && o.Uri == g.Value {
                id = o.Id
            }
        }
    }
    u.Id = id
    u.Acl.OwnerId = dsocialUserId
    return u
}

func googleplusEmailToDsocial(g *googleplus.Email, original []*Email, dsocialUserId string) *Email {
    if g == nil {
        return nil
    }
    var rel RelEmail
    switch g.Type {
    case googleplus.EMAIL_TYPE_HOME:
        rel = REL_EMAIL_HOME
    case googleplus.EMAIL_TYPE_WORK:
        rel = REL_EMAIL_WORK
    case googleplus.EMAIL_TYPE_OTHER:
        rel = REL_EMAIL_OTHER
    default:
        rel = REL_EMAIL_OTHER
    }
    e := &Email{Rel: rel, EmailAddress: g.Value, IsPrimary: g.Primary}
    id := ""
    if original != nil {
        for _, o := range original {
            if e != nil && e.EmailAddress == g.Value {
                id = o.Id
            }
        }
    }
    e.Id = id
    e.Acl.OwnerId = dsocialUserId
    return e
}

func DsocialContactToGooglePlus(c *Contact, o *googleplus.Person) *googleplus.Person {
    if c == nil {
        return nil
    }
    g := new(googleplus.Person)
    if o != nil {
        g.Id = o.Id
    } else {
        o = new(googleplus.Person)
    }
    g.Birthday = dsocialDateToGoogleString(c.Birthday)
    g.AboutMe = c.Biography
    g.Tagline = c.FavoriteQuotes
    g.Nickname = c.Nickname
    g.DisplayName = c.DisplayName
    g.Name.HonorificPrefix, g.Name.GivenName, g.Name.MiddleName, g.Name.FamilyName, g.Name.HonorificSuffix = c.Prefix, c.GivenName, c.MiddleName, c.Surname, c.Suffix
    switch c.Gender {
    case REL_GENDER_MALE:
        g.Gender = googleplus.GENDER_MALE
    case REL_GENDER_FEMALE:
        g.Gender = googleplus.GENDER_FEMALE
    case REL_GENDER_OTHER:
        g.Gender = googleplus.GENDER_OTHER
    }

    g.Urls = make([]googleplus.Url, len(c.Uris))
    for i, uri := range c.Uris {
        dsocialUriToGooglePlus(uri, &g.Urls[i])
    }
    g.PlacesLived = make([]googleplus.PlaceLived, len(c.PostalAddresses))
    for i, addr := range c.PostalAddresses {
        dsocialPostalAddressToGooglePlus(addr, &g.PlacesLived[i])
    }
    numPositions := 0
    for _, wh := range c.WorkHistories {
        if wh.Positions == nil || len(wh.Positions) == 0 {
            numPositions++
        } else {
            numPositions += len(wh.Positions)
        }
    }
    g.Organizations = make([]googleplus.Organization, len(c.Educations)+numPositions)
    orgHasPrimary := false
    for i, ed := range c.Educations {
        orgHasPrimary = dsocialEducationToGooglePlus(ed, &g.Organizations[i], orgHasPrimary)
    }
    offset := len(c.Educations)
    for _, wh := range c.WorkHistories {
        if wh.Positions == nil || len(wh.Positions) == 0 {
            a := &g.Organizations[offset]
            orgHasPrimary = dsocialWorkHistoryToGooglePlus(wh, a, orgHasPrimary)
            offset++
        } else {
            for _, pos := range wh.Positions {
                a := &g.Organizations[offset]
                dsocialWorkHistoryToGooglePlus(wh, a, orgHasPrimary)
                orgHasPrimary = dsocialWorkPositionToGooglePlus(pos, a, orgHasPrimary)
                offset++
            }
        }
    }
    g.LanguagesSpoken = make([]string, len(c.Languages))
    for i, lang := range c.Languages {
        g.LanguagesSpoken[i] = lang.Name
    }
    g.Emails = make([]googleplus.Email, len(c.EmailAddresses))
    for i, email := range c.EmailAddresses {
        dsocialEmailToGooglePlus(email, &g.Emails[i])
    }
    return g
}

func dsocialPostalAddressToGooglePlus(c *PostalAddress, g *googleplus.PlaceLived) {
    g.Value = c.Address
    g.Primary = c.IsPrimary
}

func dsocialEducationToGooglePlus(c *Education, g *googleplus.Organization, orgHasPrimary bool) bool {
    g.Title = c.Label
    g.Name = c.Institution
    g.StartDate = dsocialDateToGoogleString(c.AttendedFrom)
    g.EndDate = dsocialDateToGoogleString(c.AttendedTill)
    if c.IsCurrent && !orgHasPrimary {
        g.Primary = true
        orgHasPrimary = true
    } else {
        g.Primary = false
    }
    return orgHasPrimary
}

func dsocialWorkHistoryToGooglePlus(c *WorkHistory, g *googleplus.Organization, orgHasPrimary bool) bool {
    g.Name = c.Company
    if c.IsCurrent && !orgHasPrimary {
        g.Primary = true
        orgHasPrimary = true
    } else {
        g.Primary = false
    }
    return orgHasPrimary
}

func dsocialWorkPositionToGooglePlus(c *WorkPosition, g *googleplus.Organization, orgHasPrimary bool) bool {
    g.Title = c.Title
    g.Description = c.Description
    g.StartDate = dsocialDateToGoogleString(c.From)
    g.EndDate = dsocialDateToGoogleString(c.To)
    g.Location = c.Location
    g.Department = c.Department
    if c.IsCurrent && !orgHasPrimary {
        g.Primary = true
        orgHasPrimary = true
    } else {
        g.Primary = false
    }
    return orgHasPrimary
}

func dsocialUriToGooglePlus(c *Uri, g *googleplus.Url) {
    if c == nil {
        return
    }
    switch c.Rel {
    case REL_URI_HOME:
        g.Type = googleplus.URL_TYPE_HOME
    case REL_URI_WORK:
        g.Type = googleplus.URL_TYPE_WORK
    case REL_URI_BLOG:
        g.Type = googleplus.URL_TYPE_BLOG
    case REL_URI_GOOGLE_PLUS:
        g.Type = googleplus.URL_TYPE_PROFILE
    case REL_URI_OTHER:
        g.Type = googleplus.URL_TYPE_OTHER
    default:
        g.Type = googleplus.URL_TYPE_OTHER
    }
    g.Value = c.Uri
    g.Primary = c.IsPrimary
    return
}

func dsocialEmailToGooglePlus(c *Email, g *googleplus.Email) {
    if c == nil {
        return
    }
    switch c.Rel {
    case REL_EMAIL_HOME:
        g.Type = googleplus.EMAIL_TYPE_HOME
    case REL_EMAIL_WORK:
        g.Type = googleplus.EMAIL_TYPE_WORK
    case REL_EMAIL_OTHER:
        g.Type = googleplus.EMAIL_TYPE_OTHER
    default:
        g.Type = googleplus.EMAIL_TYPE_OTHER
    }
    g.Value = c.EmailAddress
    g.Primary = c.IsPrimary
    return
}
