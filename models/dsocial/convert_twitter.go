package dsocial

import (
    "github.com/pomack/contacts.go/twitter"
)

func TwitterUserToDsocial(l *twitter.User) *Contact {
    if l == nil {
        return nil
    }
    c := new(Contact)
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
