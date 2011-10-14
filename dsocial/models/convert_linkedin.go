package dsocial

import (
    "github.com/pomack/contacts.go/linkedin"
    "strings"
)

const (
    // im account types
    _LINKEDIN_IM_ACCOUNT_TYPE_AIM   = "aim"
    _LINKEDIN_IM_ACCOUNT_TYPE_GTALK = "gtalk"
    _LINKEDIN_IM_ACCOUNT_TYPE_ICQ   = "icq"
    _LINKEDIN_IM_ACCOUNT_TYPE_MSN   = "msn"
    _LINKEDIN_IM_ACCOUNT_TYPE_SKYPE = "skype"
    _LINKEDIN_IM_ACCOUNT_TYPE_YAHOO = "yahoo"

    // phone types
    _LINKEDIN_PHONE_TYPE_HOME   = "home"
    _LINKEDIN_PHONE_TYPE_WORK   = "work"
    _LINKEDIN_PHONE_TYPE_MOBILE = "mobile"
)

func LinkedInContactToDsocial(l *linkedin.Contact) *Contact {
    if l == nil {
        return nil
    }
    c := new(Contact)
    c.DisplayName = join(" ", l.FirstName, l.LastName)
    c.GivenName = l.FirstName
    c.Surname = l.LastName
    c.Certifications = make([]*Certification, len(l.Certifications.Values))
    for i, cert := range l.Certifications.Values {
        c.Certifications[i] = linkedinCertificationToDsocial(&cert)
    }
    c.Birthday = dateFromLinkedinDate(&l.DateOfBirth)
    c.Educations = make([]*Education, len(l.Educations.Values))
    for i, ed := range l.Educations.Values {
        c.Educations[i] = linkedinEducationToDsocial(&ed)
    }
    c.Title = l.Headline
    c.Ims = make([]*IM, len(l.ImAccounts.Values)+len(l.TwitterAccounts.Values))
    for i, im := range l.ImAccounts.Values {
        c.Ims[i] = linkedinImToDsocial(&im)
    }
    for i, twitter := range l.TwitterAccounts.Values {
        c.Ims[len(l.ImAccounts.Values)+i] = &IM{
            Rel:      REL_IM_OTHER,
            Protocol: REL_IM_PROT_TWITTER,
            Handle:   twitter.ProviderAccountName,
        }
    }
    c.PrimaryAddress = l.MainAddress
    if l.MainAddress != "" || l.Location.Country.Code != "" || l.Location.Name != "" {
        c.PostalAddresses = []*PostalAddress{&PostalAddress{
            Address:          l.MainAddress,
            Municipality:     l.Location.Name,
            Country:          l.Location.Country.Code,
            Rel:              REL_ADDRESS_HOME,
        }}
    }
    c.Skills = make([]*Skill, len(l.Skills.Values))
    for i, skill := range l.Skills.Values {
        c.Skills[i] = linkedinSkillToDsocial(&skill)
    }
    c.Languages = make([]*Language, len(l.Languages.Values))
    for i, language := range l.Languages.Values {
        c.Languages[i] = linkedinLanguageToDsocial(&language)
    }
    c.Uris = make([]*Uri, len(l.Urls.Values) + 1)
    c.Uris[0] = &Uri{Rel: REL_URI_LINKEDIN, Uri: l.PublicProfileUrl}
    for i, resource := range l.Urls.Values {
        c.Uris[i+1] = linkedinUrlResourceToDsocial(&resource)
    }
    c.PhoneNumbers = make([]*PhoneNumber, len(l.PhoneNumbers.Values))
    for i, phoneNumber := range l.PhoneNumbers.Values {
        c.PhoneNumbers[i] = linkedinPhoneNumberToDsocial(&phoneNumber)
    }
    c.WorkHistories = make([]*WorkHistory, len(l.Positions.Values))
    for i, position := range l.Positions.Values {
        c.WorkHistories[i] = linkedinPositionToDsocial(&position)
    }
    c.Biography = l.Summary
    c.Notes = l.Specialties
    return c
}

func dateFromLinkedinDate(l *linkedin.Date) (d *Date) {
    if l != nil && (l.Year != 0 || l.Month > 0 || l.Day > 0) {
        d = &Date{
            Year:  int16(l.Year),
            Month: int8(l.Month),
            Day:   int8(l.Day),
        }
    }
    return d
}

func linkedinCertificationToDsocial(l *linkedin.Certification) *Certification {
    if l == nil {
        return nil
    }
    return &Certification{
        Name:      l.Name,
        Authority: l.Authority.Name,
        Number:    l.Number,
        AsOf:      dateFromLinkedinDate(&l.StartDate),
        ValidTill: dateFromLinkedinDate(&l.EndDate),
    }
}

func linkedinEducationToDsocial(l *linkedin.Education) *Education {
    if l == nil {
        return nil
    }
    attendedFrom := dateFromLinkedinDate(&l.StartDate)
    attendedTill := dateFromLinkedinDate(&l.EndDate)
    var degrees []*Degree
    if l.Degree != "" || l.FieldOfStudy != "" {
        degrees = []*Degree{&Degree{Degree: l.Degree, Major: l.FieldOfStudy}}
    }
    rel := REL_EDUCATION_OTHER
    lInstitution := strings.ToLower(l.SchoolName)
    if strings.Contains(lInstitution, "universi") || strings.Contains(lInstitution, "college") {
        rel = REL_EDUCATION_COLLEGE
    }
    return &Education{
        Institution:  l.SchoolName,
        AttendedFrom: attendedFrom,
        AttendedTill: attendedTill,
        IsCurrent:    attendedFrom != nil && attendedTill == nil,
        Notes:        l.Notes,
        Activities:   removeEmptyStrings(strings.Split(l.Activities, ",")),
        Degrees:      degrees,
        Rel:          rel,
    }
}

func linkedinImToDsocial(l *linkedin.ImAccount) *IM {
    if l == nil {
        return nil
    }
    var protocol RelIMProtocol
    switch l.Type {
    case _LINKEDIN_IM_ACCOUNT_TYPE_AIM:
        protocol = REL_IM_PROT_AIM
    case _LINKEDIN_IM_ACCOUNT_TYPE_GTALK:
        protocol = REL_IM_PROT_GOOGLE_TALK
    case _LINKEDIN_IM_ACCOUNT_TYPE_ICQ:
        protocol = REL_IM_PROT_ICQ
    case _LINKEDIN_IM_ACCOUNT_TYPE_MSN:
        protocol = REL_IM_PROT_MSN
    case _LINKEDIN_IM_ACCOUNT_TYPE_SKYPE:
        protocol = REL_IM_PROT_SKYPE
    case _LINKEDIN_IM_ACCOUNT_TYPE_YAHOO:
        protocol = REL_IM_PROT_YAHOO_MESSENGER
    default:
        protocol = REL_IM_PROT_OTHER
    }
    return &IM{
        Rel:      REL_IM_OTHER,
        Handle:   l.Name,
        Protocol: protocol,
    }
}

func linkedinUrlResourceToDsocial(l *linkedin.Url) *Uri {
    if l == nil {
        return nil
    }
    return &Uri{
        Label: l.Name,
        Uri:   l.Url,
    }
}

func linkedinPhoneNumberToDsocial(l *linkedin.PhoneNumber) *PhoneNumber {
    if l == nil || l.Number == "" {
        return nil
    }
    var rel RelPhoneNumber
    switch l.Type {
    case _LINKEDIN_PHONE_TYPE_HOME:
        rel = REL_PHONE_HOME
    case _LINKEDIN_PHONE_TYPE_WORK:
        rel = REL_PHONE_WORK
    case _LINKEDIN_PHONE_TYPE_MOBILE:
        rel = REL_PHONE_MOBILE
    default:
        rel = REL_PHONE_OTHER
    }
    rc := &PhoneNumber{Rel: rel}
    ParsePhoneNumber(l.Number, rc)
    return rc
}

func linkedinPositionToDsocial(l *linkedin.Position) *WorkHistory {
    if l == nil {
        return nil
    }
    from := dateFromLinkedinDate(&l.StartDate)
    to := dateFromLinkedinDate(&l.EndDate)
    return &WorkHistory{
        Company: l.Company.Name,
        From:    from,
        To:      to,
        Positions: []*WorkPosition{
            &WorkPosition{
                Title:       l.Title,
                From:        from,
                To:          to,
                Description: l.Summary,
                IsCurrent:   l.IsCurrent,
            },
        },
        IsCurrent: l.IsCurrent,
    }
}

func linkedinSkillToDsocial(l *linkedin.SkillWrapper) *Skill {
    if l == nil {
        return nil
    }
    var proficiency string
    if l.Proficiency.Name != "" {
        proficiency = l.Proficiency.Name
    } else if l.Years.Name != "" {
        proficiency = l.Years.Name
    }
    return &Skill{Name: l.Skill.Name, Proficiency: proficiency}
}

func linkedinLanguageToDsocial(l *linkedin.LanguageWrapper) *Language {
    if l == nil {
        return nil
    }
    return &Language{Name: l.Language.Name}
}
