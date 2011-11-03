package dsocial

import (
    "github.com/pomack/contacts.go/smugmug"
)

func SmugMugUserToDsocial(s *smugmug.PersonReference, o *Contact, dsocialUserId string) *Contact {
    if s == nil {
        return nil
    }
    c := new(Contact)
    if o != nil {
        c.Id = o.Id
    }
    c.UserId = dsocialUserId
    c.Acl.OwnerId = dsocialUserId
    ParseName(s.Name, c)
    c.Nickname = s.NickName
    if s.Url != "" {
        c.Uris = []*Uri{&Uri{
            Rel: REL_URI_SMUGMUG,
            Uri: s.Url,
        }}
    }
    return c
}


func DsocialContactToSmugMug(c *Contact, o *smugmug.PersonReference) *smugmug.PersonReference {
    if c == nil {
        return nil
    }
    s := new(smugmug.PersonReference)
    s.Name = c.DisplayName
    s.NickName = c.Nickname
    if c.Uris != nil {
        for _, u := range c.Uris {
            if u.Rel == REL_URI_SMUGMUG {
                s.Url = u.Uri
                break
            }
        }
    }
    return s
}
