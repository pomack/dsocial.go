package dsocial

import (
    "container/list"
    "container/vector"
    "github.com/pomack/contacts.go/yahoo"
)

func YahooContactToDsocial(l *yahoo.Contact, o *Contact, dsocialUserId string) *Contact {
    if l == nil {
        return nil
    }
    c := new(Contact)
    if o != nil {
        c.Id = o.Id
    }
    c.UserId = dsocialUserId
    c.Acl.OwnerId = dsocialUserId
    emails := list.New()
    phones := list.New()
    uris := list.New()
    ims := list.New()
    addresses := list.New()
    for _, field := range l.Fields {
        switch field.Type {
        case yahoo.YAHOO_FIELD_TYPE_GUID:
            continue
        case yahoo.YAHOO_FIELD_TYPE_NICKNAME:
            c.Nickname = field.Value.(string)
        case yahoo.YAHOO_FIELD_TYPE_EMAIL:
            rel := REL_EMAIL_OTHER
            for _, flag := range field.Flags {
                switch flag {
                case yahoo.YAHOO_FIELD_FLAG_HOME, yahoo.YAHOO_FIELD_FLAG_PERSONAL:
                    rel = REL_EMAIL_HOME
                case yahoo.YAHOO_FIELD_FLAG_WORK:
                    rel = REL_EMAIL_WORK
                }
            }
            emails.PushBack(&Email{Rel: rel, EmailAddress: field.Value.(string)})
        case yahoo.YAHOO_FIELD_TYPE_YAHOOID:
            rel := REL_IM_OTHER
            for _, flag := range field.Flags {
                switch flag {
                case yahoo.YAHOO_FIELD_FLAG_HOME, yahoo.YAHOO_FIELD_FLAG_PERSONAL:
                    rel = REL_IM_HOME
                case yahoo.YAHOO_FIELD_FLAG_WORK:
                    rel = REL_IM_WORK
                }
            }
            ims.PushBack(&IM{Rel: rel, Protocol: REL_IM_PROT_YAHOO_MESSENGER, Handle: field.Value.(string)})
        case yahoo.YAHOO_FIELD_TYPE_OTHERID:
            rel := REL_IM_OTHER
            var prot RelIMProtocol
            for _, flag := range field.Flags {
                switch flag {
                case yahoo.YAHOO_FIELD_FLAG_AOL:
                    prot = REL_IM_PROT_AIM
                case yahoo.YAHOO_FIELD_FLAG_BLOG:
                    continue
                case yahoo.YAHOO_FIELD_FLAG_DOTMAC:
                    prot = REL_IM_PROT_OTHER
                case yahoo.YAHOO_FIELD_FLAG_EXTERNAL:
                    continue
                case yahoo.YAHOO_FIELD_FLAG_FAX:
                    continue
                case yahoo.YAHOO_FIELD_FLAG_GOOGLE:
                    prot = REL_IM_PROT_GOOGLE_TALK
                case yahoo.YAHOO_FIELD_FLAG_HOME:
                    rel = REL_IM_HOME
                case yahoo.YAHOO_FIELD_FLAG_IBM:
                    prot = REL_IM_PROT_SAMETIME
                case yahoo.YAHOO_FIELD_FLAG_ICQ:
                    prot = REL_IM_PROT_ICQ
                case yahoo.YAHOO_FIELD_FLAG_IRC:
                    prot = REL_IM_PROT_IRC
                case yahoo.YAHOO_FIELD_FLAG_JABBER:
                    prot = REL_IM_PROT_JABBER
                case yahoo.YAHOO_FIELD_FLAG_LCS:
                    continue
                case yahoo.YAHOO_FIELD_FLAG_MOBILE:
                    rel = REL_IM_OTHER
                case yahoo.YAHOO_FIELD_FLAG_MSN:
                    prot = REL_IM_PROT_MSN
                case yahoo.YAHOO_FIELD_FLAG_PAGER:
                    continue
                case yahoo.YAHOO_FIELD_FLAG_PERSONAL:
                    rel = REL_IM_HOME
                case yahoo.YAHOO_FIELD_FLAG_PHOTO:
                    continue
                case yahoo.YAHOO_FIELD_FLAG_QQ:
                    prot = REL_IM_PROT_QQ
                case yahoo.YAHOO_FIELD_FLAG_SKYPE:
                    prot = REL_IM_PROT_SKYPE
                case yahoo.YAHOO_FIELD_FLAG_WORK:
                    rel = REL_IM_WORK
                case yahoo.YAHOO_FIELD_FLAG_YAHOOPHONE, yahoo.YAHOO_FIELD_FLAG_YJP, yahoo.YAHOO_FIELD_FLAG_Y360:
                    continue
                }
            }
            ims.PushBack(&IM{Rel: rel, Protocol: prot, Handle: field.Value.(string)})
        case yahoo.YAHOO_FIELD_TYPE_PHONE:
            rel := REL_PHONE_OTHER
            for _, flag := range field.Flags {
                switch flag {
                case yahoo.YAHOO_FIELD_FLAG_EXTERNAL:
                    rel = REL_PHONE_EXTERNAL
                case yahoo.YAHOO_FIELD_FLAG_FAX:
                    rel = REL_PHONE_FAX
                case yahoo.YAHOO_FIELD_FLAG_GOOGLE:
                    rel = REL_PHONE_GOOGLE_VOICE
                case yahoo.YAHOO_FIELD_FLAG_HOME:
                    rel = REL_PHONE_HOME
                case yahoo.YAHOO_FIELD_FLAG_MOBILE:
                    rel = REL_PHONE_MOBILE
                case yahoo.YAHOO_FIELD_FLAG_PAGER:
                    rel = REL_PHONE_PAGER
                case yahoo.YAHOO_FIELD_FLAG_PERSONAL:
                    rel = REL_PHONE_HOME
                case yahoo.YAHOO_FIELD_FLAG_SKYPE:
                    rel = REL_PHONE_SKYPE
                case yahoo.YAHOO_FIELD_FLAG_WORK:
                    rel = REL_PHONE_WORK
                }
            }
            ph := &PhoneNumber{Rel: rel}
            ParsePhoneNumber(field.Value.(string), ph)
            phones.PushBack(ph)
        case yahoo.YAHOO_FIELD_TYPE_JOBTITLE:
            c.Title = field.Value.(string)
        case yahoo.YAHOO_FIELD_TYPE_COMPANY:
            c.Company = field.Value.(string)
        case yahoo.YAHOO_FIELD_TYPE_NOTES:
            c.Notes = field.Value.(string)
        case yahoo.YAHOO_FIELD_TYPE_LINK:
            rel := REL_URI_OTHER
            for _, flag := range field.Flags {
                switch flag {
                case yahoo.YAHOO_FIELD_FLAG_BLOG:
                    rel = REL_URI_BLOG
                case yahoo.YAHOO_FIELD_FLAG_GOOGLE:
                    rel = REL_URI_GOOGLE_PROFILE
                }
            }
            uris.PushBack(Uri{Rel: rel, Uri: field.Value.(string)})
        case yahoo.YAHOO_FIELD_TYPE_CUSTOM:
            continue
        case yahoo.YAHOO_FIELD_TYPE_NAME:
            name := field.Value.(yahoo.Name)
            c.GivenName = name.GivenName
            c.MiddleName = name.MiddleName
            c.Surname = name.FamilyName
            c.Prefix = name.Prefix
            c.Suffix = name.Suffix
            c.DisplayName = join(" ", c.Prefix, c.GivenName, c.MiddleName, c.Surname, c.Suffix)
        case yahoo.YAHOO_FIELD_TYPE_ADDRESS:
            addr := field.Value.(yahoo.Address)
            a := new(PostalAddress)
            a.StreetAddress = addr.Street
            a.Municipality = addr.City
            a.Region = addr.StateOrProvince
            a.PostalCode = addr.PostalCode
            if addr.CountryCode != "" {
                a.Country = addr.CountryCode
            } else {
                a.Country = addr.Country
            }
            a.Address = join("\n", addr.Street, join(" ", addr.City, addr.StateOrProvince, addr.PostalCode), addr.CountryCode)
            if c.PrimaryAddress == "" {
                c.PrimaryAddress = a.Address
            }
            addresses.PushBack(a)
        case yahoo.YAHOO_FIELD_TYPE_BIRTHDAY:
            birthday := field.Value.(yahoo.Date)
            d := new(Date)
            d.Year = int16(birthday.Year)
            d.Month = int8(birthday.Month)
            d.Day = int8(birthday.Day)
            c.Birthday = d
        case yahoo.YAHOO_FIELD_TYPE_ANNIVERSARY:
            anniversary := field.Value.(yahoo.Date)
            d := new(Date)
            d.Year = int16(anniversary.Year)
            d.Month = int8(anniversary.Month)
            d.Day = int8(anniversary.Day)
            c.Anniversary = d
        default:
            //panic("Unknown field type: " + field.Type)
            continue
        }
    }
    c.EmailAddresses = make([]*Email, emails.Len())
    for i, iter := 0, emails.Front(); iter != nil; i, iter = i+1, iter.Next() {
        c.EmailAddresses[i] = iter.Value.(*Email)
    }
    c.PhoneNumbers = make([]*PhoneNumber, phones.Len())
    for i, iter := 0, phones.Front(); iter != nil; i, iter = i+1, iter.Next() {
        c.PhoneNumbers[i] = iter.Value.(*PhoneNumber)
    }
    c.Uris = make([]*Uri, uris.Len())
    for i, iter := 0, uris.Front(); iter != nil; i, iter = i+1, iter.Next() {
        c.Uris[i] = iter.Value.(*Uri)
    }
    c.Ims = make([]*IM, ims.Len())
    for i, iter := 0, ims.Front(); iter != nil; i, iter = i+1, iter.Next() {
        c.Ims[i] = iter.Value.(*IM)
    }
    c.PostalAddresses = make([]*PostalAddress, addresses.Len())
    for i, iter := 0, addresses.Front(); iter != nil; i, iter = i+1, iter.Next() {
        c.PostalAddresses[i] = iter.Value.(*PostalAddress)
    }
    c.GroupReferences = make([]*GroupRef, len(l.Categories))
    for i, category := range l.Categories {
        c.GroupReferences[i] = &GroupRef{
            Name: category.Name,
        }
    }
    return c
}

func YahooCategoryToDsocial(g *yahoo.Category, o *Group, dsocialUserId string) *Group {
    if g == nil {
        return nil
    }
    c := new(Group)
    if o != nil {
        c.Id = o.Id
    }
    c.UserId = dsocialUserId
    c.Acl.OwnerId = dsocialUserId
    c.Name = g.Name
    return c
}

func DsocialGroupToYahoo(g *Group, o *yahoo.Category) *yahoo.Category {
    if g == nil {
        return nil
    }
    c := new(yahoo.Category)
    if o != nil {
        c.Id = o.Id
        c.Created = o.Created
        c.Updated = o.Updated
        c.Uri = o.Uri
    }
    c.Name = g.Name
    return c
}

func DsocialContactToYahoo(c *Contact, o *yahoo.Contact) *yahoo.Contact {
    if c == nil {
        return nil
    }
    y := new(yahoo.Contact)
    if o != nil {
        y.Id = o.Id
    }
    fields := list.New()
    if c.Nickname != "" {
        fields.PushBack(&yahoo.ContactField{
            Type:  yahoo.YAHOO_FIELD_TYPE_NICKNAME,
            Value: c.Nickname,
        })
    }
    for _, email := range c.EmailAddresses {
        flag := ""
        switch email.Rel {
        case REL_EMAIL_HOME:
            flag = yahoo.YAHOO_FIELD_FLAG_HOME
        case REL_EMAIL_WORK:
            flag = yahoo.YAHOO_FIELD_FLAG_WORK
        }
        var flags []string = nil
        if flag != "" {
            flags = make([]string, 1)
            flags[0] = flag
        }
        fields.PushBack(&yahoo.ContactField{
            Type:  yahoo.YAHOO_FIELD_TYPE_EMAIL,
            Value: email.EmailAddress,
            Flags: flags,
        })
    }
    for _, im := range c.Ims {
        var flags vector.StringVector
        thetype := yahoo.YAHOO_FIELD_TYPE_OTHERID
        switch im.Rel {
        case REL_IM_HOME:
            flags.Push(yahoo.YAHOO_FIELD_FLAG_HOME)
        case REL_IM_WORK:
            flags.Push(yahoo.YAHOO_FIELD_FLAG_WORK)
        }
        switch im.Protocol {
        case REL_IM_PROT_YAHOO_MESSENGER:
            thetype = yahoo.YAHOO_FIELD_TYPE_YAHOOID
        case REL_IM_PROT_AIM:
            flags.Push(yahoo.YAHOO_FIELD_FLAG_AOL)
        case REL_IM_PROT_DOTMAC:
            flags.Push(yahoo.YAHOO_FIELD_FLAG_DOTMAC)
        case REL_IM_PROT_GOOGLE_TALK:
            flags.Push(yahoo.YAHOO_FIELD_FLAG_GOOGLE)
        case REL_IM_PROT_SAMETIME:
            flags.Push(yahoo.YAHOO_FIELD_FLAG_IBM)
        case REL_IM_PROT_ICQ:
            flags.Push(yahoo.YAHOO_FIELD_FLAG_ICQ)
        case REL_IM_PROT_IRC:
            flags.Push(yahoo.YAHOO_FIELD_FLAG_IRC)
        case REL_IM_PROT_JABBER:
            flags.Push(yahoo.YAHOO_FIELD_FLAG_JABBER)
        case REL_IM_PROT_MSN:
            flags.Push(yahoo.YAHOO_FIELD_FLAG_MSN)
        case REL_IM_PROT_QQ:
            flags.Push(yahoo.YAHOO_FIELD_FLAG_QQ)
        case REL_IM_PROT_SKYPE:
            flags.Push(yahoo.YAHOO_FIELD_FLAG_SKYPE)
        default:
            continue
        }
        fields.PushBack(&yahoo.ContactField{
            Type:  thetype,
            Value: im.Handle,
            Flags: []string(flags),
        })
    }
    for _, phone := range c.PhoneNumbers {
        var flags vector.StringVector
        switch phone.Rel {
        case REL_PHONE_EXTERNAL:
            flags.Push(yahoo.YAHOO_FIELD_FLAG_EXTERNAL)
        case REL_PHONE_FAX:
            flags.Push(yahoo.YAHOO_FIELD_FLAG_FAX)
        case REL_PHONE_GOOGLE_VOICE:
            flags.Push(yahoo.YAHOO_FIELD_FLAG_GOOGLE)
        case REL_PHONE_HOME:
            flags.Push(yahoo.YAHOO_FIELD_FLAG_HOME)
        case REL_PHONE_MOBILE:
            flags.Push(yahoo.YAHOO_FIELD_FLAG_MOBILE)
        case REL_PHONE_PAGER:
            flags.Push(yahoo.YAHOO_FIELD_FLAG_PAGER)
        case REL_PHONE_SKYPE:
            flags.Push(yahoo.YAHOO_FIELD_FLAG_SKYPE)
        case REL_PHONE_WORK:
            flags.Push(yahoo.YAHOO_FIELD_FLAG_WORK)
        }
        fields.PushBack(&yahoo.ContactField{
            Type:  yahoo.YAHOO_FIELD_TYPE_PHONE,
            Value: phone.FormattedNumber,
            Flags: []string(flags),
        })
    }
    if c.Title != "" {
        fields.PushBack(&yahoo.ContactField{
            Type:  yahoo.YAHOO_FIELD_TYPE_JOBTITLE,
            Value: c.Title,
        })
    }
    if c.Company != "" {
        fields.PushBack(&yahoo.ContactField{
            Type:  yahoo.YAHOO_FIELD_TYPE_COMPANY,
            Value: c.Company,
        })
    }
    if c.Notes != "" {
        fields.PushBack(&yahoo.ContactField{
            Type:  yahoo.YAHOO_FIELD_TYPE_NOTES,
            Value: c.Notes,
        })
    }
    if c.Company != "" {
        fields.PushBack(&yahoo.ContactField{
            Type:  yahoo.YAHOO_FIELD_TYPE_COMPANY,
            Value: c.Company,
        })
    }
    for _, uri := range c.Uris {
        var flags vector.StringVector
        switch uri.Rel {
        case REL_URI_BLOG:
            flags.Push(yahoo.YAHOO_FIELD_FLAG_BLOG)
        case REL_URI_GOOGLE_PROFILE:
            flags.Push(yahoo.YAHOO_FIELD_FLAG_GOOGLE)
        }
        fields.PushBack(&yahoo.ContactField{
            Type:  yahoo.YAHOO_FIELD_TYPE_LINK,
            Value: uri.Uri,
            Flags: []string(flags),
        })
    }
    if c.Prefix != "" || c.GivenName != "" || c.MiddleName != "" || c.Surname != "" || c.Suffix != "" {
        fields.PushBack(&yahoo.ContactField{
            Type: yahoo.YAHOO_FIELD_TYPE_LINK,
            Value: yahoo.Name{
                Prefix:     c.Prefix,
                GivenName:  c.GivenName,
                MiddleName: c.MiddleName,
                FamilyName: c.Surname,
                Suffix:     c.Suffix,
            },
        })
    }
    for _, addr := range c.PostalAddresses {
        var flags vector.StringVector
        switch addr.Rel {
        case REL_ADDRESS_HOME:
            flags.Push(yahoo.YAHOO_FIELD_FLAG_HOME)
        case REL_ADDRESS_WORK:
            flags.Push(yahoo.YAHOO_FIELD_FLAG_WORK)
        }
        country := ""
        countryCode := ""
        if len(addr.Country) <= 3 {
            countryCode = addr.Country
        } else {
            country = addr.Country
        }
        fields.PushBack(&yahoo.ContactField{
            Type: yahoo.YAHOO_FIELD_TYPE_ADDRESS,
            Value: yahoo.Address{
                Street:          addr.StreetAddress,
                City:            addr.Municipality,
                StateOrProvince: addr.Region,
                PostalCode:      addr.PostalCode,
                Country:         country,
                CountryCode:     countryCode,
            },
            Flags: []string(flags),
        })
    }
    if c.Birthday != nil && !c.Birthday.IsEmpty() {
        fields.PushBack(&yahoo.ContactField{
            Type: yahoo.YAHOO_FIELD_TYPE_BIRTHDAY,
            Value: yahoo.Date{
                Year:  int(c.Birthday.Year),
                Month: int(c.Birthday.Month),
                Day:   int(c.Birthday.Day),
            },
        })
    }
    if c.Anniversary != nil && !c.Anniversary.IsEmpty() {
        fields.PushBack(&yahoo.ContactField{
            Type: yahoo.YAHOO_FIELD_TYPE_ANNIVERSARY,
            Value: yahoo.Date{
                Year:  int(c.Anniversary.Year),
                Month: int(c.Anniversary.Month),
                Day:   int(c.Anniversary.Day),
            },
        })
    }
    for _, thedate := range c.Dates {
        if thedate.Rel == REL_DATE_ANNIVERSARY {
            fields.PushBack(&yahoo.ContactField{
                Type: yahoo.YAHOO_FIELD_TYPE_ANNIVERSARY,
                Value: yahoo.Date{
                    Year:  int(thedate.Value.Year),
                    Month: int(thedate.Value.Month),
                    Day:   int(thedate.Value.Day),
                },
            })
        }
    }
    if fields.Len() > 0 {
        y.Fields = make([]yahoo.ContactField, fields.Len())
        for i, e := 0, fields.Front(); e != nil; i, e = i+1, e.Next() {
            y.Fields[i] = *(e.Value.(*yahoo.ContactField))
        }
    }
    if c.GroupReferences != nil && len(c.GroupReferences) > 0 {
        y.Categories = make([]yahoo.Category, len(c.GroupReferences))
        for i, groupReference := range c.GroupReferences {
            y.Categories[i].Name = groupReference.Name
        }
    }
    return y
}
