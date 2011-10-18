package dsocial

import (
    "github.com/pomack/contacts.go/smugmug"
)

func SmugMugUserToDsocial(s *smugmug.PersonReference) *Contact {
    if s == nil {
        return nil
    }
    c := new(Contact)
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
