package dsocial

import (
    "container/list"
    "fmt"
    "github.com/pomack/contacts.go/google"
    "strconv"
    "strings"
    "time"
)

func googleHandleUnknownRel(rel string) (value string) {
    if index := strings.Index(rel, "#"); index != -1 {
        value = rel[index:]
    } else {
        value = rel
    }
    return
}

func googleDateTimeStringToDsocial(s string) (d *DateTime) {
    if len(s) == 10 {
        // YYYY-MM-DD
        d = new(DateTime)
        year, _ := strconv.Atoi(s[0:4])
        month, _ := strconv.Atoi(s[5:7])
        day, _ := strconv.Atoi(s[8:10])
        d.Year = int16(year)
        d.Month = int8(month)
        d.Day = int8(day)
    } else if len(s) == 7 {
        t, err := time.Parse("--01-02", s)
        if err == nil && !t.IsZero() {
            d = new(DateTime)
            d.Month = int8(t.Month())
            d.Day = int8(t.Day())
        }
    } else if len(s) == 25 {
        // YYYY-MM-DDTHH:MM:SS-05:00
        t, err := time.Parse("2006-01-02T15:04:05-07:00", s)
        if err == nil && !t.IsZero() {
            t2 := t.UTC()
            d = new(DateTime)
            d.Year = int16(t2.Year())
            d.Month = int8(t2.Month())
            d.Day = int8(t2.Day())
            d.Hour = int8(t2.Hour())
            d.Minute = int8(t2.Minute())
            d.Second = int8(t2.Second())
        }
    }
    return
}

func googleDateStringToDsocial(s string) (d *Date) {
    if len(s) == 10 {
        // YYYY-MM-DD
        d = new(Date)
        year, _ := strconv.Atoi(s[0:4])
        month, _ := strconv.Atoi(s[5:7])
        day, _ := strconv.Atoi(s[8:10])
        d.Year = int16(year)
        d.Month = int8(month)
        d.Day = int8(day)
    } else if len(s) == 7 {
        // --MM-DD
        d = new(Date)
        month, _ := strconv.Atoi(s[2:4])
        day, _ := strconv.Atoi(s[5:7])
        d.Month = int8(month)
        d.Day = int8(day)
    }
    return
}

func googleEmailToDsocial(g *google.Email, original []*Email, dsocialUserId string) *Email {
    var rel RelEmail
    var label string
    switch g.Rel {
    case "":
        label = g.Label
    case google.REL_HOME:
        rel = REL_EMAIL_HOME
    case google.REL_WORK:
        rel = REL_EMAIL_WORK
    case google.REL_OTHER:
        rel = REL_EMAIL_OTHER
    default:
        label = googleHandleUnknownRel(g.Rel)
    }
    if label == "" && rel == "" {
        rel = REL_EMAIL_OTHER
    }
    id := ""
    if original != nil {
        for _, e := range original {
            if e != nil && e.EmailAddress == g.Address {
                id = e.Id
            }
        }
    }
    e := &Email{EmailAddress: g.Address, Label: label, Rel: rel, IsPrimary: g.Primary == "true"}
    e.Id = id
    e.Acl.OwnerId = dsocialUserId
    return e
}

func findPrimaryEmail(arr []*Email) string {
    for _, addr := range arr {
        if addr.IsPrimary {
            return addr.EmailAddress
        }
    }
    if len(arr) > 0 {
        return arr[0].EmailAddress
    }
    return ""
}

func googlePhoneNumberToDsocial(g *google.PhoneNumber, original []*PhoneNumber, dsocialUserId string) *PhoneNumber {
    var rel RelPhoneNumber
    var label string
    switch g.Rel {
    case "":
        label = g.Label
    case google.REL_HOME:
        rel = REL_PHONE_HOME
    case google.REL_WORK:
        rel = REL_PHONE_WORK
    case google.REL_OTHER:
        rel = REL_PHONE_OTHER
    case google.REL_ASSISTANT:
        rel = REL_PHONE_ASSISTANT
    case google.REL_CALLBACK:
        rel = REL_PHONE_CALLBACK
    case google.REL_CAR:
        rel = REL_PHONE_CAR
    case google.REL_COMPANY_MAIN:
        rel = REL_PHONE_COMPANY_MAIN
    case google.REL_FAX:
        rel = REL_PHONE_FAX
    case google.REL_HOME_FAX:
        rel = REL_PHONE_HOME_FAX
    case google.REL_ISDN:
        rel = REL_PHONE_ISDN
    case google.REL_MAIN:
        rel = REL_PHONE_MAIN
    case google.REL_MOBILE:
        rel = REL_PHONE_MOBILE
    case google.REL_OTHER_FAX:
        rel = REL_PHONE_OTHER_FAX
    case google.REL_PAGER:
        rel = REL_PHONE_PAGER
    case google.REL_RADIO:
        rel = REL_PHONE_RADIO
    case google.REL_TELEX:
        rel = REL_PHONE_TELEX
    case google.REL_TTY_TDD:
        rel = REL_PHONE_TTY_TDD
    case google.REL_WORK_FAX:
        rel = REL_PHONE_WORK_FAX
    case google.REL_WORK_MOBILE:
        rel = REL_PHONE_WORK_MOBILE
    case google.REL_WORK_PAGER:
        rel = REL_PHONE_WORK_PAGER
    default:
        label = googleHandleUnknownRel(g.Rel)
    }
    if label == "" && rel == "" {
        rel = REL_PHONE_OTHER
    }
    ph := &PhoneNumber{Label: label, Rel: rel, IsPrimary: g.Primary == "true"}
    ParsePhoneNumber(g.Value, ph)
    id := ""
    if original != nil {
        for _, o := range original {
            if o != nil && o.FormattedNumber == g.Value {
                id = o.Id
            }
        }
    }
    ph.Id = id
    ph.Acl.OwnerId = dsocialUserId
    return ph
}

func findPrimaryPhoneNumber(arr []*PhoneNumber) string {
    for _, phone := range arr {
        if phone.IsPrimary {
            return phone.FormattedNumber
        }
    }
    if len(arr) > 0 {
        return arr[0].FormattedNumber
    }
    return ""
}

func googleImToDsocial(g *google.Im, original []*IM, dsocialUserId string) *IM {
    var rel RelIM
    var label string
    switch g.Rel {
    case "":
        label = g.Label
    case google.REL_HOME:
        rel = REL_IM_HOME
    case google.REL_NETMEETING:
        rel = REL_IM_NETMEETING
    case google.REL_WORK:
        rel = REL_IM_WORK
    case google.REL_OTHER:
        rel = REL_IM_OTHER
    default:
        label = googleHandleUnknownRel(g.Rel)
    }
    if label == "" && rel == "" {
        rel = REL_IM_OTHER
    }

    var protocol RelIMProtocol
    switch g.Protocol {
    case google.IM_PROTOCOL_AIM:
        protocol = REL_IM_PROT_AIM
    case google.IM_PROTOCOL_MSN:
        protocol = REL_IM_PROT_MSN
    case google.IM_PROTOCOL_YAHOO:
        protocol = REL_IM_PROT_YAHOO_MESSENGER
    case google.IM_PROTOCOL_SKYPE:
        protocol = REL_IM_PROT_SKYPE
    case google.IM_PROTOCOL_QQ:
        protocol = REL_IM_PROT_QQ
    case google.IM_PROTOCOL_GOOGLE_TALK:
        protocol = REL_IM_PROT_GOOGLE_TALK
    case google.IM_PROTOCOL_ICQ:
        protocol = REL_IM_PROT_ICQ
    case google.IM_PROTOCOL_JABBER:
        protocol = REL_IM_PROT_JABBER
    default:
        protocol = REL_IM_PROT_OTHER
    }
    id := ""
    if original != nil {
        for _, o := range original {
            if o != nil && o.Protocol == protocol && o.Handle == g.Address {
                id = o.Id
            }
        }
    }
    im := &IM{Handle: g.Address, Protocol: protocol, Label: label, Rel: rel, IsPrimary: g.Primary == "true"}
    im.Id = id
    im.Acl.OwnerId = dsocialUserId
    return im
}

func findPrimaryIM(arr []*IM) string {
    for _, addr := range arr {
        if addr.IsPrimary {
            return join(":", string(addr.Protocol), addr.Handle)
        }
    }
    if len(arr) > 0 {
        return join(":", string(arr[0].Protocol), arr[0].Handle)
    }
    return ""
}

func googlePostalAddressToDsocial(g *google.PostalAddress, original []*PostalAddress, dsocialUserId string) *PostalAddress {
    var rel RelPostalAddress
    var label string
    switch g.Rel {
    case "":
        label = g.Label
    case google.REL_HOME:
        rel = REL_ADDRESS_HOME
    case google.REL_WORK:
        rel = REL_ADDRESS_WORK
    case google.REL_OTHER:
        rel = REL_ADDRESS_OTHER
    default:
        label = googleHandleUnknownRel(g.Rel)
    }
    if label == "" && rel == "" {
        rel = REL_ADDRESS_OTHER
    }
    id := ""
    if original != nil {
        for _, o := range original {
            if o != nil && o.Address == g.Value {
                id = o.Id
            }
        }
    }
    pa := &PostalAddress{Address: g.Value, Label: label, Rel: rel, IsPrimary: g.Primary == "true", IsCurrent: true}
    pa.Id = id
    pa.Acl.OwnerId = dsocialUserId
    return pa
}

func googleStructuredPostalAddressToDsocial(g *google.StructuredPostalAddress, original []*PostalAddress, dsocialUserId string) *PostalAddress {
    var rel RelPostalAddress
    var label string
    switch g.Rel {
    case "":
        label = g.Label
    case google.REL_HOME:
        rel = REL_ADDRESS_HOME
    case google.REL_WORK:
        rel = REL_ADDRESS_WORK
    case google.REL_OTHER:
        rel = REL_ADDRESS_OTHER
    default:
        label = googleHandleUnknownRel(g.Rel)
    }
    if label == "" && rel == "" {
        rel = REL_ADDRESS_OTHER
    }
    other := ""
    if g.POBox.Value != "" {
        other = "P.O. Box " + g.POBox.Value
    } else if g.Neighborhood.Value != "" {
        other = g.Neighborhood.Value
    } else if g.Agent.Value != "" {
        other = g.Agent.Value
    } else if g.HouseName.Value != "" {
        other = g.HouseName.Value
    }
    pa := &PostalAddress{Address: g.FormattedAddress.Value,
        StreetAddress: g.Street.Value,
        OtherAddress:  other,
        Municipality:  g.City.Value,
        Region:        g.Region.Value,
        PostalCode:    g.Postcode.Value,
        Country:       g.Country.Value,
        Label:         label,
        Rel:           rel,
        IsPrimary:     g.Primary == "true",
        IsCurrent:     true,
    }
    id := ""
    if original != nil {
        for _, o := range original {
            if o != nil && o.Address == g.FormattedAddress.Value {
                id = o.Id
            }
        }
    }
    pa.Id = id
    pa.Acl.OwnerId = dsocialUserId
    return pa
}

func findPrimaryAddress(arr []*PostalAddress) string {
    for _, addr := range arr {
        if addr.IsPrimary {
            return addr.Address
        }
    }
    if len(arr) > 0 {
        return arr[0].Address
    }
    return ""
}

func googleEventToDsocial(g *google.Event, dsocialUserId string) (*ContactDate, *ContactDateTime) {
    dt := googleDateTimeStringToDsocial(g.When.StartTime)
    if dt == nil {
        return nil, nil
    }
    var label string
    var rel string
    switch g.Rel {
    case "":
        label = g.Label
    case google.REL_EVENT_ANNIVERSARY:
        rel = string(REL_DATE_ANNIVERSARY)
    case google.REL_EVENT_OTHER:
        rel = string(REL_DATE_OTHER)
    default:
        label = googleHandleUnknownRel(g.Rel)
    }
    if label == "" && rel == "" {
        rel = REL_OTHER
    }
    if dt.Hour == 0 && dt.Minute == 0 && dt.Second == 0 {
        value := &ContactDate{
            Rel:   RelDate(rel),
            Label: label,
            Value: &Date{Year: dt.Year, Month: dt.Month, Day: dt.Day},
        }
        value.Acl.OwnerId = dsocialUserId
        return value, nil
    }
    value := &ContactDateTime{
        Rel:   RelDateTime(rel),
        Label: label,
        Value: dt,
    }
    value.Acl.OwnerId = dsocialUserId
    return nil, value
}

func googleRelationToDsocial(g *google.Relation, original []*Relationship, dsocialUserId string) *Relationship {
    var rel RelRelationship
    var label string
    switch g.Rel {
    case "":
        label = g.Label
    case google.REL_RELATION_ASSISTANT:
        rel = REL_RELATIONSHIP_ASSISTANT
    case google.REL_RELATION_BROTHER:
        rel = REL_RELATIONSHIP_BROTHER
    case google.REL_RELATION_CHILD:
        rel = REL_RELATIONSHIP_CHILD
    case google.REL_RELATION_DOMESTIC_PARTNER:
        rel = REL_RELATIONSHIP_DOMESTIC_PARTNER
    case google.REL_RELATION_FATHER:
        rel = REL_RELATIONSHIP_FATHER
    case google.REL_RELATION_FRIEND:
        rel = REL_RELATIONSHIP_FRIEND
    case google.REL_RELATION_MANAGER:
        rel = REL_RELATIONSHIP_MANAGER
    case google.REL_RELATION_MOTHER:
        rel = REL_RELATIONSHIP_MOTHER
    case google.REL_RELATION_PARENT:
        rel = REL_RELATIONSHIP_PARENT
    case google.REL_RELATION_PARTNER:
        rel = REL_RELATIONSHIP_PARTNER
    case google.REL_RELATION_REFERRED_BY:
        rel = REL_RELATIONSHIP_REFERRED_BY
    case google.REL_RELATION_RELATIVE:
        rel = REL_RELATIONSHIP_RELATIVE
    case google.REL_RELATION_SISTER:
        rel = REL_RELATIONSHIP_SISTER
    case google.REL_RELATION_SPOUSE:
        rel = REL_RELATIONSHIP_SPOUSE
    default:
        label = googleHandleUnknownRel(g.Rel)
    }
    if label == "" && rel == "" {
        rel = REL_RELATIONSHIP_OTHER
    }
    r := &Relationship{ContactReferenceName: g.Value, Label: label, Rel: rel}
    id := ""
    if original != nil {
        for _, o := range original {
            if o != nil && o.ContactReferenceName == g.Value {
                id = o.Id
            }
        }
    }
    r.Id = id
    r.Acl.OwnerId = dsocialUserId
    return r
}

func googleWebsiteToDsocial(g *google.Website, original []*Uri, dsocialUserId string) *Uri {
    var rel RelUri
    var label string
    switch g.Rel {
    case "":
        label = g.Label
    case google.REL_WEBSITE_HOME_PAGE:
        rel = REL_URI_HOMEPAGE
    case google.REL_WEBSITE_BLOG:
        rel = REL_URI_BLOG
    case google.REL_WEBSITE_PROFILE:
        rel = REL_URI_GOOGLE_PROFILE
    case google.REL_WEBSITE_HOME:
        rel = REL_URI_HOME
    case google.REL_WEBSITE_WORK:
        rel = REL_URI_WORK
    case google.REL_WEBSITE_OTHER:
        rel = REL_URI_OTHER
    case google.REL_WEBSITE_FTP:
        rel = REL_URI_FTP
    default:
        label = googleHandleUnknownRel(g.Rel)
    }
    if label == "" && rel == "" {
        rel = REL_URI_OTHER
    }
    u := &Uri{Uri: g.Href, Label: label, Rel: rel, IsPrimary: g.Primary == "true"}
    id := ""
    if original != nil {
        for _, o := range original {
            if o != nil && o.Uri == g.Href {
                id = o.Id
            }
        }
    }
    u.Id = id
    u.Acl.OwnerId = dsocialUserId
    return u
}

func findPrimaryUri(arr []*Uri) string {
    for _, addr := range arr {
        if addr.IsPrimary {
            return addr.Uri
        }
    }
    if len(arr) > 0 {
        return arr[0].Uri
    }
    return ""
}

func GoogleContactToDsocial(g *google.Contact, o *Contact, dsocialUserId string) *Contact {
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
    c.DisplayName = g.Title.Value
    c.Notes = g.Content.Value
    c.Prefix = g.Name.NamePrefix.Value
    c.GivenName = g.Name.GivenName.Value
    c.MiddleName = g.Name.AdditionalName.Value
    c.Surname = g.Name.FamilyName.Value
    c.Suffix = g.Name.NameSuffix.Value
    c.Nickname = g.Nickname.Value
    c.MaidenName = g.MaidenName.Value
    if len(g.Name.FullName.Value) > 0 {
        c.DisplayName = g.Name.FullName.Value
    }
    switch g.Gender.Value {
    case google.GENDER_MALE:
        c.Gender = REL_GENDER_MALE
    case google.GENDER_FEMALE:
        c.Gender = REL_GENDER_FEMALE
    }
    for _, org := range g.Organizations {
        isPrimary := org.Primary == "true"
        if isPrimary || org.OrgTitle.Value != "" {
            c.Title = org.OrgTitle.Value
        }
        if isPrimary || org.OrgName.Value != "" {
            c.Company = org.OrgName.Value
        }
        if isPrimary || org.OrgDepartment.Value != "" {
            c.Department = org.OrgDepartment.Value
        }
        if isPrimary {
            break
        }
    }
    c.EmailAddresses = make([]*Email, len(g.Emails))
    for i, email := range g.Emails {
        c.EmailAddresses[i] = googleEmailToDsocial(&email, o.EmailAddresses, dsocialUserId)
    }
    c.PrimaryEmail = findPrimaryEmail(c.EmailAddresses)

    c.PhoneNumbers = make([]*PhoneNumber, len(g.PhoneNumbers))
    for i, phoneNumber := range g.PhoneNumbers {
        c.PhoneNumbers[i] = googlePhoneNumberToDsocial(&phoneNumber, o.PhoneNumbers, dsocialUserId)
    }
    c.PrimaryPhoneNumber = findPrimaryPhoneNumber(c.PhoneNumbers)

    c.Ims = make([]*IM, len(g.Ims))
    for i, im := range g.Ims {
        c.Ims[i] = googleImToDsocial(&im, o.Ims, dsocialUserId)
    }
    c.PrimaryIm = findPrimaryIM(c.Ims)

    if len(g.PostalAddresses) > 0 {
        c.PostalAddresses = make([]*PostalAddress, len(g.PostalAddresses))
        for i, addr := range g.PostalAddresses {
            c.PostalAddresses[i] = googlePostalAddressToDsocial(&addr, o.PostalAddresses, dsocialUserId)
        }
    } else if len(g.StructuredPostalAddresses) > 0 {
        c.PostalAddresses = make([]*PostalAddress, len(g.StructuredPostalAddresses))
        for i, addr := range g.StructuredPostalAddresses {
            c.PostalAddresses[i] = googleStructuredPostalAddressToDsocial(&addr, o.PostalAddresses, dsocialUserId)
        }
    }
    c.PrimaryAddress = findPrimaryAddress(c.PostalAddresses)

    c.Birthday = googleDateStringToDsocial(g.Birthday.When)
    dates := list.New()
    datetimes := list.New()
    for _, event := range g.Events {
        thedate, thedatetime := googleEventToDsocial(&event, dsocialUserId)
        if thedate != nil {
            dates.PushBack(thedate)
            if event.Rel == google.REL_EVENT_ANNIVERSARY {
                c.Anniversary = thedate.Value
            }
        }
        if thedatetime != nil {
            datetimes.PushBack(thedatetime)
            if event.Rel == google.REL_EVENT_ANNIVERSARY {
                c.Anniversary = &Date{
                    Year:  thedatetime.Value.Year,
                    Month: thedatetime.Value.Month,
                    Day:   thedatetime.Value.Day,
                }
            }
        }
    }
    n := dates.Len()
    c.Dates = make([]*ContactDate, n)
    for i, e := 0, dates.Front(); e != nil; i, e = i+1, e.Next() {
        c.Dates[i] = e.Value.(*ContactDate)
    }
    n = datetimes.Len()
    c.DateTimes = make([]*ContactDateTime, n)
    for i, e := 0, datetimes.Front(); e != nil; i, e = i+1, e.Next() {
        c.DateTimes[i] = e.Value.(*ContactDateTime)
    }
    // TODO group memberships
    c.Relationships = make([]*Relationship, len(g.Relationships))
    for i, relation := range g.Relationships {
        c.Relationships[i] = googleRelationToDsocial(&relation, o.Relationships, dsocialUserId)
    }
    c.Uris = make([]*Uri, len(g.Websites))
    for i, website := range g.Websites {
        c.Uris[i] = googleWebsiteToDsocial(&website, o.Uris, dsocialUserId)
    }
    return c
}

func GoogleGroupToDsocial(g *google.ContactGroup, o *Group, dsocialUserId string) *Group {
    if g == nil {
        return nil
    }
    c := new(Group)
    if o != nil {
        c.Id = o.Id
        c.UserId = o.UserId
    }
    c.UserId = dsocialUserId
    c.Acl.OwnerId = dsocialUserId
    c.Name = g.Title.Value
    c.Description = g.Content.Value
    return c
}

func DsocialGroupToGoogle(g *Group, o *google.ContactGroup) *google.ContactGroup {
    if g == nil {
        return nil
    }
    c := new(google.ContactGroup)
    if o != nil {
        c.Id.Value = o.Id.Value
        c.Etag = o.Etag
        c.Categories = o.Categories
        c.Updated = o.Updated
    }
    c.Title.Value = g.Name
    c.Content.Value = g.Description
    c.Xmlns = google.XMLNS_ATOM
    c.XmlnsGcontact = google.XMLNS_GCONTACT
    c.XmlnsBatch = google.XMLNS_GDATA_BATCH
    c.XmlnsGd = google.XMLNS_GD
    return c
}

///
/// Begin conversions from dsocial Contact to Google
///

func dsocialDateTimeToGoogleString(d *DateTime) (s string) {
    if d.Month != 0 && d.Day != 0 {
        if d.Year != 0 {
            if d.Hour != 0 || d.Minute != 0 || d.Second != 0 {
                // YYYY-MM-DDTHH:MM:SS-05:00
                s = fmt.Sprintf("%04d-%02d-%02dT%02d:%02d:%02d-00:00", int(d.Year), int(d.Month), int(d.Day), int(d.Hour), int(d.Minute), int(d.Second))
            } else {
                // YYYY-MM-DD
                s = fmt.Sprintf("%04d-%02d-%02d", int(d.Year), int(d.Month), int(d.Day))
            }
        } else {
            // --MM-DD
            s = fmt.Sprintf("--%02d-%02d", int(d.Month), int(d.Day))
        }
    }
    return
}

func dsocialDateToGoogleString(d *Date) (s string) {
    if d.Month != 0 && d.Day != 0 {
        if d.Year == 0 {
            s = fmt.Sprintf("--%02d-%02d", int(d.Month), int(d.Day))
        } else {
            s = fmt.Sprintf("%04d-%02d-%02d", int(d.Year), int(d.Month), int(d.Day))
        }
    }
    return
}

func dsocialEmailToGoogle(c *Email, g *google.Email, primaryEmail string) {
    var rel string
    var label string
    switch c.Rel {
    case "":
        label = g.Label
    case REL_EMAIL_HOME:
        rel = google.REL_HOME
    case REL_EMAIL_WORK:
        rel = google.REL_WORK
    case REL_EMAIL_OTHER:
        rel = google.REL_OTHER
    default:
        label = googleHandleUnknownRel(string(c.Rel))
    }
    if label == "" && rel == "" {
        rel = google.REL_OTHER
    }
    g.Address = c.EmailAddress
    g.Label = label
    g.Rel = rel
    if c.IsPrimary || c.EmailAddress == primaryEmail {
        g.Primary = "true"
    } else {
        g.Primary = ""
    }
    return
}

func dsocialPhoneNumberToGoogle(c *PhoneNumber, g *google.PhoneNumber, primaryNumber string) {
    var rel string
    var label string
    switch c.Rel {
    case "":
        label = c.Label
    case REL_PHONE_HOME:
        rel = google.REL_HOME
    case REL_PHONE_WORK:
        rel = google.REL_WORK
    case REL_PHONE_OTHER:
        rel = google.REL_OTHER
    case REL_PHONE_ASSISTANT:
        rel = google.REL_ASSISTANT
    case REL_PHONE_CALLBACK:
        rel = google.REL_CALLBACK
    case REL_PHONE_CAR:
        rel = google.REL_CAR
    case REL_PHONE_COMPANY_MAIN:
        rel = google.REL_COMPANY_MAIN
    case REL_PHONE_FAX:
        rel = google.REL_FAX
    case REL_PHONE_HOME_FAX:
        rel = google.REL_HOME_FAX
    case REL_PHONE_ISDN:
        rel = google.REL_ISDN
    case REL_PHONE_MAIN:
        rel = google.REL_MAIN
    case REL_PHONE_MOBILE:
        rel = google.REL_MOBILE
    case REL_PHONE_OTHER_FAX:
        rel = google.REL_OTHER_FAX
    case REL_PHONE_PAGER:
        rel = google.REL_PAGER
    case REL_PHONE_RADIO:
        rel = google.REL_RADIO
    case REL_PHONE_TELEX:
        rel = google.REL_TELEX
    case REL_PHONE_TTY_TDD:
        rel = google.REL_TTY_TDD
    case REL_PHONE_WORK_FAX:
        rel = google.REL_WORK_FAX
    case REL_PHONE_WORK_MOBILE:
        rel = google.REL_WORK_MOBILE
    case REL_PHONE_WORK_PAGER:
        rel = google.REL_WORK_PAGER
    default:
        label = googleHandleUnknownRel(string(c.Rel))
    }
    if label == "" && rel == "" {
        rel = google.REL_OTHER
    }
    g.Value = c.FormattedNumber
    g.Label = label
    g.Rel = rel
    g.Uri = ""
    if c.IsPrimary || c.FormattedNumber == primaryNumber {
        g.Primary = "true"
    } else {
        g.Primary = ""
    }
    return
}

func dsocialImToGoogle(c *IM, g *google.Im, primaryIm string) {
    var rel string
    var label string
    switch c.Rel {
    case "":
        label = c.Label
    case REL_IM_HOME:
        rel = google.REL_HOME
    case REL_IM_NETMEETING:
        rel = google.REL_NETMEETING
    case REL_IM_WORK:
        rel = google.REL_WORK
    case REL_IM_OTHER:
        rel = google.REL_OTHER
    default:
        label = googleHandleUnknownRel(string(c.Rel))
    }
    if label == "" && rel == "" {
        rel = google.REL_OTHER
    }

    var protocol string
    switch c.Protocol {
    case REL_IM_PROT_AIM:
        protocol = google.IM_PROTOCOL_AIM
    case REL_IM_PROT_MSN:
        protocol = google.IM_PROTOCOL_MSN
    case REL_IM_PROT_YAHOO_MESSENGER:
        protocol = google.IM_PROTOCOL_YAHOO
    case REL_IM_PROT_SKYPE:
        protocol = google.IM_PROTOCOL_SKYPE
    case REL_IM_PROT_QQ:
        protocol = google.IM_PROTOCOL_QQ
    case REL_IM_PROT_GOOGLE_TALK:
        protocol = google.IM_PROTOCOL_GOOGLE_TALK
    case REL_IM_PROT_ICQ:
        protocol = google.IM_PROTOCOL_ICQ
    case REL_IM_PROT_JABBER:
        protocol = google.IM_PROTOCOL_JABBER
    default:
        protocol = google.REL_OTHER
    }
    g.Label = label
    g.Rel = rel
    g.Address = c.Handle
    g.Protocol = protocol
    if c.IsPrimary || protocol+":"+c.Handle == primaryIm {
        g.Primary = "true"
    } else {
        g.Primary = ""
    }
    return
}

func dsocialPostalAddressToGoogle(c *PostalAddress, g *google.StructuredPostalAddress, primaryAddr string) {
    var rel string
    var label string
    switch c.Rel {
    case "":
        label = c.Label
    case REL_ADDRESS_HOME:
        rel = google.REL_HOME
    case REL_ADDRESS_WORK:
        rel = google.REL_WORK
    case REL_ADDRESS_OTHER:
        rel = google.REL_OTHER
    default:
        label = googleHandleUnknownRel(string(c.Rel))
    }
    if label == "" && rel == "" {
        rel = google.REL_OTHER
    }
    // TODO figure out how to handle all these types
    /*
       other := ""
       if g.POBox.Value != "" {
           other = "P.O. Box " + g.POBox.Value
       } else if g.Neighborhood.Value != "" {
           other = g.Neighborhood.Value
       } else if g.Agent.Value != "" {
           other = g.Agent.Value
       } else if g.HouseName.Value != "" {
           other = g.HouseName.Value
       }
    */
    g.Rel = rel
    g.Label = label
    g.FormattedAddress.Value = c.Address
    g.Street.Value = c.StreetAddress
    g.City.Value = c.Municipality
    g.Region.Value = c.Region
    g.Postcode.Value = c.PostalCode
    g.Country.Value = c.Country
    if c.IsPrimary || c.Address == primaryAddr {
        g.Primary = "true"
    } else {
        g.Primary = ""
    }
    lowerOther := strings.ToLower(c.OtherAddress)
    if strings.HasPrefix(lowerOther, "p.o. box") || strings.HasPrefix(lowerOther, "po box") {
        g.POBox.Value = c.OtherAddress
    } else {
        g.Neighborhood.Value = c.OtherAddress
    }
    return
}

func dsocialDateTimeToGoogle(c *ContactDateTime, g *google.Event) {
    g.When.StartTime = dsocialDateTimeToGoogleString(c.Value)
    var label string
    var rel string
    switch c.Rel {
    case "":
        label = c.Label
    case REL_DATETIME_OTHER:
        rel = google.REL_EVENT_OTHER
    default:
        label = googleHandleUnknownRel(string(c.Rel))
    }
    if label == "" && rel == "" {
        rel = google.REL_SHORT_OTHER
    }
    g.Label = label
    g.Rel = rel
    return
}

func dsocialDateToGoogle(c *ContactDate, g *google.Event) {
    g.When.StartTime = dsocialDateToGoogleString(c.Value)
    var label string
    var rel string
    switch c.Rel {
    case "":
        label = c.Label
    case REL_DATE_ANNIVERSARY:
        rel = google.REL_EVENT_ANNIVERSARY
    case REL_DATE_OTHER:
        rel = google.REL_EVENT_OTHER
    default:
        label = googleHandleUnknownRel(string(c.Rel))
    }
    if label == "" && rel == "" {
        rel = google.REL_SHORT_OTHER
    }
    g.Label = label
    g.Rel = rel
    return
}

func dsocialRelationshipToGoogle(c *Relationship, g *google.Relation) {
    var rel string
    var label string
    switch c.Rel {
    case "":
        label = c.Label
    case REL_RELATIONSHIP_ASSISTANT:
        rel = google.REL_RELATION_ASSISTANT
    case REL_RELATIONSHIP_BROTHER:
        rel = google.REL_RELATION_BROTHER
    case REL_RELATIONSHIP_CHILD:
        rel = google.REL_RELATION_CHILD
    case REL_RELATIONSHIP_DOMESTIC_PARTNER:
        rel = google.REL_RELATION_DOMESTIC_PARTNER
    case REL_RELATIONSHIP_FATHER:
        rel = google.REL_RELATION_FATHER
    case REL_RELATIONSHIP_FRIEND:
        rel = google.REL_RELATION_FRIEND
    case REL_RELATIONSHIP_MANAGER:
        rel = google.REL_RELATION_MANAGER
    case REL_RELATIONSHIP_MOTHER:
        rel = google.REL_RELATION_MOTHER
    case REL_RELATIONSHIP_PARENT:
        rel = google.REL_RELATION_PARENT
    case REL_RELATIONSHIP_PARTNER:
        rel = google.REL_RELATION_PARTNER
    case REL_RELATIONSHIP_REFERRED_BY:
        rel = google.REL_RELATION_REFERRED_BY
    case REL_RELATIONSHIP_RELATIVE:
        rel = google.REL_RELATION_RELATIVE
    case REL_RELATIONSHIP_SISTER:
        rel = google.REL_RELATION_SISTER
    case REL_RELATIONSHIP_SPOUSE:
        rel = google.REL_RELATION_SPOUSE
    default:
        label = googleHandleUnknownRel(string(c.Rel))
    }
    if label == "" && rel == "" {
        rel = google.REL_SHORT_OTHER
    }
    g.Label = label
    g.Rel = rel
    g.Value = c.ContactReferenceName
    return
}

func dsocialUriToGoogle(c *Uri, g *google.Website) {
    var rel string
    var label string
    switch c.Rel {
    case "":
        label = c.Label
    case REL_URI_HOMEPAGE:
        rel = google.REL_WEBSITE_HOME_PAGE
    case REL_URI_BLOG:
        rel = google.REL_WEBSITE_BLOG
    case REL_URI_GOOGLE_PROFILE:
        rel = google.REL_WEBSITE_PROFILE
    case REL_URI_HOME:
        rel = google.REL_WEBSITE_HOME
    case REL_URI_WORK:
        rel = google.REL_WEBSITE_WORK
    case REL_URI_OTHER:
        rel = google.REL_WEBSITE_OTHER
    case REL_URI_FTP:
        rel = google.REL_WEBSITE_FTP
    default:
        label = googleHandleUnknownRel(string(c.Rel))
    }
    if label == "" && rel == "" {
        rel = google.REL_SHORT_OTHER
    }
    g.Href = c.Uri
    g.Label = label
    g.Rel = rel
    if c.IsPrimary {
        g.Primary = "true"
    } else {
        g.Primary = ""
    }
    return
}

func DsocialContactToGoogle(c *Contact, o *google.Contact) *google.Contact {
    if c == nil {
        return nil
    }
    g := new(google.Contact)
    if o != nil {
        g.Id.Value = o.Id.Value
        g.Etag = o.Etag
    }
    // not supposed to specify title on output
    //g.Title.Value = c.DisplayName
    g.Content.Value = c.Notes
    g.Name.NamePrefix.Value, g.Name.GivenName.Value, g.Name.AdditionalName.Value, g.Name.FamilyName.Value, g.Name.NameSuffix.Value = c.Prefix, c.GivenName, c.MiddleName, c.Surname, c.Suffix
    g.Nickname.Value = c.Nickname
    g.MaidenName.Value = c.MaidenName
    g.Name.FullName.Value = c.DisplayName
    if c.Title != "" || c.Company != "" || c.Department != "" {
        g.Organizations = make([]google.Organization, 1)
        org := &g.Organizations[0]
        org.OrgTitle.Value = c.Title
        org.OrgName.Value = c.Company
        org.OrgDepartment.Value = c.Department
        org.Primary = "true"
    }
    switch c.Gender {
    case REL_GENDER_MALE:
        g.Gender.Value = google.GENDER_MALE
    case REL_GENDER_FEMALE:
        g.Gender.Value = google.GENDER_FEMALE
    }
    // skip for now
    currentWorkHistories := list.New()
    if c.WorkHistories != nil {
        for _, workhist := range c.WorkHistories {
            if workhist != nil && workhist.IsCurrent {
                currentWorkHistories.PushBack(workhist)
            }
        }
    }
    if currentWorkHistories.Len() > 0 {
        g.Organizations = make([]google.Organization, currentWorkHistories.Len())
        for i, iter := 0, currentWorkHistories.Front(); iter != nil; i, iter = i+1, iter.Next() {
            workhist := iter.Value.(*WorkHistory)
            org := &g.Organizations[i]
            org.OrgName.Value = workhist.Company
            if workhist.Positions != nil {
                for _, position := range workhist.Positions {
                    if position != nil && position.IsCurrent {
                        org.OrgTitle.Value = position.Title
                        org.OrgDepartment.Value = position.Department
                        break
                    }
                }
            }
        }
    }
    g.Emails = make([]google.Email, len(c.EmailAddresses))
    for i, email := range c.EmailAddresses {
        dsocialEmailToGoogle(email, &g.Emails[i], c.PrimaryEmail)
    }
    g.PhoneNumbers = make([]google.PhoneNumber, len(c.PhoneNumbers))
    for i, ph := range c.PhoneNumbers {
        dsocialPhoneNumberToGoogle(ph, &g.PhoneNumbers[i], c.PrimaryPhoneNumber)
    }
    g.StructuredPostalAddresses = make([]google.StructuredPostalAddress, len(c.PostalAddresses))
    for i, addr := range c.PostalAddresses {
        dsocialPostalAddressToGoogle(addr, &g.StructuredPostalAddresses[i], c.PrimaryAddress)
    }
    g.Ims = make([]google.Im, len(c.Ims))
    for i, im := range c.Ims {
        dsocialImToGoogle(im, &g.Ims[i], c.PrimaryIm)
    }
    g.Birthday.When = dsocialDateToGoogleString(c.Birthday)

    events := list.New()
    if c.Anniversary != nil && !c.Anniversary.IsEmpty() {
        event := new(google.Event)
        event.Rel = google.REL_EVENT_ANNIVERSARY
        event.When.StartTime = dsocialDateToGoogleString(c.Anniversary)
        events.PushBack(event)
    }
    for _, thedate := range c.Dates {
        if thedate != nil {
            event := new(google.Event)
            dsocialDateToGoogle(thedate, event)
            events.PushBack(event)
        }
    }
    for _, thedatetime := range c.DateTimes {
        if thedatetime != nil {
            event := new(google.Event)
            dsocialDateTimeToGoogle(thedatetime, event)
            events.PushBack(event)
        }
    }
    n := events.Len()
    g.Events = make([]google.Event, n)
    for i, e := 0, events.Front(); e != nil; i, e = i+1, e.Next() {
        g.Events[i] = *(e.Value.(*google.Event))
    }
    // TODO group memberships
    g.Relationships = make([]google.Relation, len(c.Relationships))
    for i, relation := range c.Relationships {
        dsocialRelationshipToGoogle(relation, &g.Relationships[i])
    }
    g.Websites = make([]google.Website, len(c.Uris))
    for i, uri := range c.Uris {
        dsocialUriToGoogle(uri, &g.Websites[i])
    }
    return g
}
