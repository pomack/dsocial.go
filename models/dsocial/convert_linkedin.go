package dsocial

import (
    "github.com/pomack/contacts.go/linkedin"
    "container/list"
    "strconv"
    "strings"
)

func LinkedInContactToDsocial(l *linkedin.Contact, o *Contact, dsocialUserId string) *Contact {
    if l == nil {
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
    c.DisplayName = join(" ", l.FirstName, l.LastName)
    c.GivenName = l.FirstName
    c.Surname = l.LastName
    c.Certifications = make([]*Certification, len(l.Certifications.Values))
    for i, cert := range l.Certifications.Values {
        c.Certifications[i] = linkedinCertificationToDsocial(&cert, o.Certifications, dsocialUserId)
    }
    c.Birthday = dateFromLinkedinDate(&l.DateOfBirth)
    c.Educations = make([]*Education, len(l.Educations.Values))
    for i, ed := range l.Educations.Values {
        c.Educations[i] = linkedinEducationToDsocial(&ed, o.Educations, dsocialUserId)
    }
    c.Title = l.Headline
    c.Ims = make([]*IM, len(l.ImAccounts.Values)+len(l.TwitterAccounts.Values))
    for i, im := range l.ImAccounts.Values {
        c.Ims[i] = linkedinImToDsocial(&im, o.Ims, dsocialUserId)
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
            Address:      l.MainAddress,
            Municipality: l.Location.Name,
            Country:      l.Location.Country.Code,
            Rel:          REL_ADDRESS_HOME,
        }}
    }
    c.Skills = make([]*Skill, len(l.Skills.Values))
    for i, skill := range l.Skills.Values {
        c.Skills[i] = linkedinSkillToDsocial(&skill, o.Skills, dsocialUserId)
    }
    c.Languages = make([]*Language, len(l.Languages.Values))
    for i, language := range l.Languages.Values {
        c.Languages[i] = linkedinLanguageToDsocial(&language, o.Languages, dsocialUserId)
    }
    c.Uris = make([]*Uri, len(l.Urls.Values)+1)
    c.Uris[0] = &Uri{Rel: REL_URI_LINKEDIN, Uri: l.PublicProfileUrl}
    for i, resource := range l.Urls.Values {
        c.Uris[i+1] = linkedinUrlResourceToDsocial(&resource, o.Uris, dsocialUserId)
    }
    c.PhoneNumbers = make([]*PhoneNumber, len(l.PhoneNumbers.Values))
    for i, phoneNumber := range l.PhoneNumbers.Values {
        c.PhoneNumbers[i] = linkedinPhoneNumberToDsocial(&phoneNumber, o.PhoneNumbers, dsocialUserId)
    }
    c.WorkHistories = make([]*WorkHistory, len(l.Positions.Values))
    for i, position := range l.Positions.Values {
        c.WorkHistories[i] = linkedinPositionToDsocial(&position, o.WorkHistories, dsocialUserId)
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

func linkedinCertificationToDsocial(l *linkedin.Certification, original []*Certification, dsocialUserId string) *Certification {
    if l == nil {
        return nil
    }
    c := &Certification{
        Name:      l.Name,
        Authority: l.Authority.Name,
        Number:    l.Number,
        AsOf:      dateFromLinkedinDate(&l.StartDate),
        ValidTill: dateFromLinkedinDate(&l.EndDate),
    }
    id := ""
    if original != nil {
        for _, o := range original {
            if c != nil && c.Name == l.Name {
                id = o.Id
            }
        }
    }
    c.Id = id
    c.Acl.OwnerId = dsocialUserId
    return c
}

func linkedinEducationToDsocial(l *linkedin.Education, original []*Education, dsocialUserId string) *Education {
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
    e := &Education{
        Institution:  l.SchoolName,
        AttendedFrom: attendedFrom,
        AttendedTill: attendedTill,
        IsCurrent:    attendedFrom != nil && attendedTill == nil,
        Notes:        l.Notes,
        Activities:   removeEmptyStrings(strings.Split(l.Activities, ",")),
        Degrees:      degrees,
        Rel:          rel,
    }
    id := ""
    if original != nil {
        for _, o := range original {
            if isSimilar, _ := e.IsSimilarOrUpdated(e, o); isSimilar {
                id = o.Id
            }
        }
    }
    e.Id = id
    e.Acl.OwnerId = dsocialUserId
    return e
}

func linkedinImToDsocial(l *linkedin.ImAccount, original []*IM, dsocialUserId string) *IM {
    if l == nil {
        return nil
    }
    var protocol RelIMProtocol
    switch l.Type {
    case linkedin.IM_ACCOUNT_TYPE_AIM:
        protocol = REL_IM_PROT_AIM
    case linkedin.IM_ACCOUNT_TYPE_GTALK:
        protocol = REL_IM_PROT_GOOGLE_TALK
    case linkedin.IM_ACCOUNT_TYPE_ICQ:
        protocol = REL_IM_PROT_ICQ
    case linkedin.IM_ACCOUNT_TYPE_MSN:
        protocol = REL_IM_PROT_MSN
    case linkedin.IM_ACCOUNT_TYPE_SKYPE:
        protocol = REL_IM_PROT_SKYPE
    case linkedin.IM_ACCOUNT_TYPE_YAHOO:
        protocol = REL_IM_PROT_YAHOO_MESSENGER
    default:
        protocol = REL_IM_PROT_OTHER
    }
    i := &IM{
        Rel:      REL_IM_OTHER,
        Handle:   l.Name,
        Protocol: protocol,
    }
    id := ""
    if original != nil {
        for _, o := range original {
            if isSimilar, _ := i.IsSimilarOrUpdated(i, o); isSimilar {
                id = o.Id
            }
        }
    }
    i.Id = id
    i.Acl.OwnerId = dsocialUserId
    return i
}

func linkedinUrlResourceToDsocial(l *linkedin.Url, original []*Uri, dsocialUserId string) *Uri {
    if l == nil {
        return nil
    }
    u := &Uri{
        Label: l.Name,
        Uri:   l.Url,
    }
    id := ""
    if original != nil {
        for _, o := range original {
            if isSimilar, _ := u.IsSimilarOrUpdated(u, o); isSimilar {
                id = o.Id
            }
        }
    }
    u.Id = id
    u.Acl.OwnerId = dsocialUserId
    return u
}

func linkedinPhoneNumberToDsocial(l *linkedin.PhoneNumber, original []*PhoneNumber, dsocialUserId string) *PhoneNumber {
    if l == nil || l.Number == "" {
        return nil
    }
    var rel RelPhoneNumber
    switch l.Type {
    case linkedin.PHONE_TYPE_HOME:
        rel = REL_PHONE_HOME
    case linkedin.PHONE_TYPE_WORK:
        rel = REL_PHONE_WORK
    case linkedin.PHONE_TYPE_MOBILE:
        rel = REL_PHONE_MOBILE
    default:
        rel = REL_PHONE_OTHER
    }
    ph := &PhoneNumber{Rel: rel}
    ParsePhoneNumber(l.Number, ph)
    id := ""
    if original != nil {
        for _, o := range original {
            if isSimilar, _ := ph.IsSimilarOrUpdated(ph, o); isSimilar {
                id = o.Id
            }
        }
    }
    ph.Id = id
    ph.Acl.OwnerId = dsocialUserId
    return ph
}

func linkedinPositionToDsocial(l *linkedin.Position, original []*WorkHistory, dsocialUserId string) *WorkHistory {
    if l == nil {
        return nil
    }
    from := dateFromLinkedinDate(&l.StartDate)
    to := dateFromLinkedinDate(&l.EndDate)
    w := &WorkHistory{
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
    id := ""
    if original != nil {
        for _, o := range original {
            if isSimilar, _ := w.IsSimilarOrUpdated(w, o); isSimilar {
                id = o.Id
            }
        }
    }
    w.Id = id
    w.Acl.OwnerId = dsocialUserId
    return w
}

func linkedinSkillToDsocial(l *linkedin.SkillWrapper, original []*Skill, dsocialUserId string) *Skill {
    if l == nil {
        return nil
    }
    var proficiency string
    if l.Proficiency.Name != "" {
        proficiency = l.Proficiency.Name
    } else if l.Years.Name != "" {
        proficiency = l.Years.Name
    }
    s := &Skill{Name: l.Skill.Name, Proficiency: proficiency}
    id := ""
    if original != nil {
        for _, o := range original {
            if isSimilar, _ := s.IsSimilarOrUpdated(s, o); isSimilar {
                id = o.Id
            }
        }
    }
    s.Id = id
    s.Acl.OwnerId = dsocialUserId
    return s
}

func linkedinLanguageToDsocial(l *linkedin.LanguageWrapper, original []*Language, dsocialUserId string) *Language {
    if l == nil {
        return nil
    }
    lang := &Language{Name: l.Language.Name}
    id := ""
    if original != nil {
        for _, o := range original {
            if isSimilar, _ := lang.IsSimilarOrUpdated(lang, o); isSimilar {
                id = o.Id
            }
        }
    }
    lang.Id = id
    lang.Acl.OwnerId = dsocialUserId
    return lang
}


func DsocialContactToLinkedIn(c *Contact, o *linkedin.Contact) (l *linkedin.Contact) {
    if c == nil {
        return
    }
    l = new(linkedin.Contact)
    if o != nil {
        l.Id = o.Id
    }
    l.FirstName = c.GivenName
    l.LastName = c.Surname
    size := len(c.Certifications)
    l.Certifications.Total = size
    l.Certifications.Values = make([]linkedin.Certification, size)
    for i, cert := range c.Certifications {
        dsocialCertificationToLinkedIn(cert, &l.Certifications.Values[i])
    }
    dsocialDateToLinkedIn(c.Birthday, &l.DateOfBirth)
    size = len(c.Educations)
    l.Educations.Total = size
    l.Educations.Values = make([]linkedin.Education, size)
    for i, ed := range c.Educations {
        dsocialEducationToLinkedIn(ed, &l.Educations.Values[i])
    }
    l.Headline = c.Title
    twitterAccounts := list.New()
    imAccounts := list.New()
    for _, im := range c.Ims {
        if im.Protocol == REL_IM_PROT_TWITTER {
            twitterAccounts.PushBack(im)
        } else {
            imAccounts.PushBack(im)
        }
    }
    size = imAccounts.Len()
    l.ImAccounts.Total = size
    l.ImAccounts.Values = make([]linkedin.ImAccount, size)
    for i, e := 0, imAccounts.Front(); e != nil; i, e = i + 1, e.Next() {
        dsocialImToLinkedIn(e.Value.(*IM), &l.ImAccounts.Values[i])
    }
    size = twitterAccounts.Len()
    l.TwitterAccounts.Total = size
    l.TwitterAccounts.Values = make([]linkedin.TwitterAccount, size)
    for i, e := 0, twitterAccounts.Front(); e != nil; i, e = i + 1, e.Next() {
        l.TwitterAccounts.Values[i].ProviderAccountName = e.Value.(*IM).Handle
    }
    l.MainAddress = c.PrimaryAddress
    if c.PostalAddresses != nil && len(c.PostalAddresses) > 0 {
        var primaryPostalAddress *PostalAddress = nil
        for _, addr := range c.PostalAddresses {
            if addr.IsPrimary {
                primaryPostalAddress = addr
                break
            }
        }
        if primaryPostalAddress == nil {
            primaryPostalAddress = c.PostalAddresses[0]
        }
        l.Location.Name = primaryPostalAddress.Municipality
        l.Location.Country.Code = primaryPostalAddress.Country
    }
    size = len(c.Skills)
    l.Skills.Total = size
    l.Skills.Values = make([]linkedin.SkillWrapper, size)
    for i, skill := range c.Skills {
        dsocialSkillToLinkedIn(skill, &l.Skills.Values[i])
    }
    size = len(c.Languages)
    l.Languages.Total = size
    l.Languages.Values = make([]linkedin.LanguageWrapper, size)
    for i, skill := range c.Languages {
        dsocialLanguageToLinkedIn(skill, &l.Languages.Values[i])
    }
    size = len(c.Uris)
    for _, uri := range c.Uris {
        if uri.Rel == REL_URI_LINKEDIN {
            l.PublicProfileUrl = uri.Uri
            size -= 1
        }
    }
    l.Urls.Total = size
    l.Urls.Values = make([]linkedin.Url, size)
    offset := 0
    for i, uri := range c.Uris {
        if uri.Rel == REL_URI_LINKEDIN && l.PublicProfileUrl == uri.Uri {
            offset -= 1
            continue
        }
        dsocialUriToLinkedIn(uri, &l.Urls.Values[i + offset])
    }
    size = len(c.PhoneNumbers)
    l.PhoneNumbers.Total = size
    l.PhoneNumbers.Values = make([]linkedin.PhoneNumber, size)
    for i, phone := range c.PhoneNumbers {
        dsocialPhoneNumberToLinkedIn(phone, &l.PhoneNumbers.Values[i])
    }
    size = 0
    for _, wh := range c.WorkHistories {
        if wh.Positions == nil || len(wh.Positions) == 0 {
            size++
        } else {
            size += len(wh.Positions)
        }
    }
    l.Positions.Total = size
    l.Positions.Values = make([]linkedin.Position, size)
    offset = 0
    for _, wh := range c.WorkHistories {
        if wh.Positions == nil || len(wh.Positions) == 0 {
            dsocialWorkHistoryToLinkedIn(wh, nil, &l.Positions.Values[offset])
            offset++
        } else {
            for _, pos := range wh.Positions {
                dsocialWorkHistoryToLinkedIn(wh, pos, &l.Positions.Values[offset])
                offset++
            }
        }
    }
    l.Summary = c.Biography
    l.Specialties = c.Notes
    return
}

func dsocialDateToLinkedIn(d *Date, l *linkedin.Date) {
    if d != nil && (d.Year != 0 || d.Month > 0 || d.Day > 0) {
        if d.Year > 0 { l.Year = int(d.Year) }
        if d.Month > 0 { l.Month = int(d.Month) }
        if d.Day > 0 { l.Day = int(d.Day) }
    }
}

func dsocialCertificationToLinkedIn(c *Certification, l *linkedin.Certification) {
    if c == nil {
        return
    }
    l.Name = c.Name
    l.Authority.Name = c.Authority
    l.Number = c.Number
    dsocialDateToLinkedIn(c.AsOf, &l.StartDate)
    dsocialDateToLinkedIn(c.ValidTill, &l.EndDate)
}

func dsocialEducationToLinkedIn(c *Education, l *linkedin.Education) {
    if c == nil {
        return
    }
    l.SchoolName = c.Institution
    dsocialDateToLinkedIn(c.AttendedFrom, &l.StartDate)
    dsocialDateToLinkedIn(c.AttendedTill, &l.EndDate)
    l.Notes = c.Notes
    l.Activities = strings.Join(c.Activities, ", ")
    if c.Degrees != nil && len(c.Degrees) > 0 {
        l.Degree = c.Degrees[0].Degree
        l.FieldOfStudy = c.Degrees[0].Major
    }
}

func dsocialImToLinkedIn(c *IM, l *linkedin.ImAccount) {
    if c == nil {
        return
    }
    switch c.Protocol {
    case REL_IM_PROT_AIM:
        l.Type = linkedin.IM_ACCOUNT_TYPE_AIM
    case REL_IM_PROT_GOOGLE_TALK:
        l.Type = linkedin.IM_ACCOUNT_TYPE_GTALK
    case REL_IM_PROT_ICQ:
        l.Type = linkedin.IM_ACCOUNT_TYPE_ICQ
    case REL_IM_PROT_MSN:
        l.Type = linkedin.IM_ACCOUNT_TYPE_MSN
    case REL_IM_PROT_SKYPE:
        l.Type = linkedin.IM_ACCOUNT_TYPE_SKYPE
    case REL_IM_PROT_YAHOO_MESSENGER:
        l.Type = linkedin.IM_ACCOUNT_TYPE_YAHOO
    default:
        l.Type = ""
    }
    l.Name = c.Handle
}

func dsocialUriToLinkedIn(c *Uri, l *linkedin.Url) {
    if c == nil {
        return
    }
    l.Name = c.Label
    l.Url = c.Uri
}

func dsocialPhoneNumberToLinkedIn(c *PhoneNumber, l *linkedin.PhoneNumber) {
    if c == nil || c.FormattedNumber == "" {
        return
    }
    switch c.Rel {
    case REL_PHONE_HOME:
        l.Type = linkedin.PHONE_TYPE_HOME
    case REL_PHONE_WORK:
        l.Type = linkedin.PHONE_TYPE_WORK
    case REL_PHONE_MOBILE:
        l.Type = linkedin.PHONE_TYPE_MOBILE
    default:
        l.Type = ""
    }
    l.Number = c.FormattedNumber
}

func dsocialWorkHistoryToLinkedIn(wh *WorkHistory, pos *WorkPosition, l *linkedin.Position) {
    if wh == nil && pos == nil {
        return
    }
    if wh != nil {
        l.Company.Name = wh.Company
        dsocialDateToLinkedIn(pos.From, &l.StartDate)
        dsocialDateToLinkedIn(pos.To, &l.EndDate)
        l.IsCurrent = wh.IsCurrent
    }
    if pos != nil {
        dsocialDateToLinkedIn(pos.From, &l.StartDate)
        dsocialDateToLinkedIn(pos.To, &l.EndDate)
        l.Title = pos.Title
        l.Summary = pos.Description
        l.IsCurrent = pos.IsCurrent
    }
}

func dsocialSkillToLinkedIn(c *Skill, l *linkedin.SkillWrapper) {
    if l == nil {
        return
    }
    l.Skill.Name = c.Name
    years, _ := strconv.Atoi(c.Proficiency)
    if years > 0 {
        l.Years.Name = c.Proficiency
    } else {
        l.Proficiency.Name = c.Proficiency
    }
}

func dsocialLanguageToLinkedIn(c *Language, l *linkedin.LanguageWrapper) {
    if l == nil {
        return
    }
    l.Language.Name = c.Name
}
