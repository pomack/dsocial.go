package dsocial

import (
    "github.com/pomack/contacts.go/twitter"
)

func TwitterUserToDsocial(l *twitter.User, o *Contact, dsocialUserId string) *Contact {
    if l == nil {
        return nil
    }
    c := new(Contact)
    if o != nil {
        c.Id = o.Id
    }
    c.UserId = dsocialUserId
    c.Acl.OwnerId = dsocialUserId
    ParseName(l.Name, c)
    if l.ScreenName != "" {
        c.Ims = []*IM{&IM{
            Rel:      REL_IM_OTHER,
            Protocol: REL_IM_PROT_TWITTER,
            Handle:   l.ScreenName,
        }}
    }
    if l.Lang != "" {
        c.Languages = []*Language{&Language{
            Name: l.Lang,
        }}
    }
    c.Biography = l.Description
    c.PrimaryAddress = l.Location
    if l.Url != nil && *l.Url != "" {
        c.Uris = []*Uri{&Uri{
            Rel: REL_URI_TWITTER,
            Uri: *l.Url,
        }}
    }
    return c
}

func DsocialContactToTwitter(c *Contact, o *twitter.User) *twitter.User {
    if c == nil {
        return nil
    }
    t := new(twitter.User)
    if o != nil {
        t.Id = o.Id
        t.IdStr = o.IdStr
    }
    t.Name = c.DisplayName
    if c.Ims != nil {
        for _, im := range c.Ims {
            if im.Protocol == REL_IM_PROT_TWITTER {
                t.ScreenName = im.Handle
                break
            }
        }
    }
    if c.Languages != nil && len(c.Languages) > 0 {
        t.Lang = c.Languages[0].Name
    }
    t.Description = c.Biography
    t.Location = c.PrimaryAddress
    if c.Uris != nil {
        for _, uri := range c.Uris {
            if uri.Rel == REL_URI_TWITTER {
                t.Url = &uri.Uri
                break
            }
        }
    }
    return t
}
