package dsocial

import (
    "github.com/pomack/contacts.go/googleplus"
    "container/list"
)

func GooglePlusPersonToDsocial(g *googleplus.Person) *Contact {
    if g == nil {
        return nil
    }
    c := new(Contact)
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
        c.Uris[i] = googleplusUriToDsocial(&uri)
    }
    c.PostalAddresses = make([]*PostalAddress, len(g.PlacesLived))
    for i, addr := range g.PlacesLived {
        c.PostalAddresses[i] = &PostalAddress{
            Address:addr.Value,
            Rel:REL_ADDRESS_OTHER,
            IsPrimary:addr.Primary,
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
                Label:org.Title,
                Institution:org.Name,
                AttendedFrom:from,
                AttendedTill:to,
                IsCurrent:isCurrent,
            })
        case googleplus.ORGANIZATION_TYPE_WORK:
            from := googleDateStringToDsocial(org.StartDate)
            to := googleDateStringToDsocial(org.EndDate)
            isCurrent := (from != nil && to == nil) || org.Primary
            workhistories.PushBack(&WorkHistory{
                Positions:[]*WorkPosition{&WorkPosition{
                   Title:org.Title, 
                   Description:org.Description,
                   From:from,
                   To:to,
                   IsCurrent:isCurrent,
                   Location:org.Location,
                   Department:org.Department,
                }},
                Company:org.Name,
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
        c.Languages[i] = &Language{Name:lang}
    }
    c.EmailAddresses = make([]*Email, len(g.Emails))
    for i, email := range g.Emails {
        c.EmailAddresses[i] = googleplusEmailToDsocial(&email)
    }
    return c
}

func googleplusUriToDsocial(g *googleplus.Url) *Uri {
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
    return &Uri{Rel:rel, Uri:g.Value, IsPrimary:g.Primary}
}

func googleplusEmailToDsocial(g *googleplus.Email) *Email {
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
    return &Email{Rel:rel, EmailAddress:g.Value, IsPrimary:g.Primary}
}




