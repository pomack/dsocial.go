package dsocial_test

import (
    . "github.com/pomack/dsocial.go/models/dsocial"
    "testing"
)

func TestApplyAddChange(t *testing.T) {
    contact := &Contact{
        DisplayName: "Martha Parke Custis",
        GivenName:   "Martha",
        MiddleName:  "Parke",
        Surname:     "Custis",
    }

    nickname := "Patsy"
    ApplyAddChange(contact, nickname, []*PathComponent{NewPathComponentKey("nickname")})
    if contact.Nickname != nickname {
        t.Fatalf("Expected Nickname to be updated to %#v but was %#v", nickname, contact.Nickname)
    }

    d := &Date{Year: 1776, Month: 8, Day: 21}
    ApplyAddChange(contact, d, []*PathComponent{NewPathComponentKey("birthday")})
    if contact.Birthday == nil || contact.Birthday.Year != d.Year || contact.Birthday.Month != d.Month || contact.Birthday.Day != d.Day {
        t.Fatalf("Expected birthday to be updated to %#v but was %#v", d, contact.Birthday)
    }

    externalUserId1 := "blah1"
    externalUserId2 := "blah2"
    externalUserId3 := "blah3"
    externalUserId4 := "blah4"
    ApplyAddChange(contact, externalUserId1, []*PathComponent{NewPathComponentKey("external_user_ids")})
    ApplyAddChange(contact, externalUserId2, []*PathComponent{NewPathComponentKey("external_user_ids")})
    ApplyAddChange(contact, externalUserId3, []*PathComponent{NewPathComponentKey("external_user_ids")})
    ApplyAddChange(contact, externalUserId4, []*PathComponent{NewPathComponentKey("external_user_ids"), NewPathComponentIndex("1")})
    if contact.ExternalUserIds == nil || len(contact.ExternalUserIds) != 4 || contact.ExternalUserIds[0] != externalUserId1 || contact.ExternalUserIds[1] != externalUserId4 || contact.ExternalUserIds[2] != externalUserId2 || contact.ExternalUserIds[3] != externalUserId3 {
        t.Fatalf("Expected ExternalUserId to be updated to %#v but was %#v", []string{externalUserId1, externalUserId4, externalUserId2, externalUserId3}, contact.ExternalUserIds)
    }

    relationship1 := &Relationship{Rel: REL_RELATIONSHIP_HUSBAND, ContactReferenceName: "George Washington"}
    relationship2 := &Relationship{Rel: REL_RELATIONSHIP_CHILD, ContactReferenceName: "Daniel Custis"}
    relationship3 := &Relationship{Rel: REL_RELATIONSHIP_CHILD, ContactReferenceName: "Francis Custis"}
    relationship4 := &Relationship{Rel: REL_RELATIONSHIP_CHILD, ContactReferenceName: "John Parke Custis"}
    relationship5 := &Relationship{Rel: REL_RELATIONSHIP_CHILD, ContactReferenceName: "Martha Parke Custis"}
    ApplyAddChange(contact, relationship1, []*PathComponent{NewPathComponentKey("relationships")})
    ApplyAddChange(contact, relationship2, []*PathComponent{NewPathComponentKey("relationships")})
    ApplyAddChange(contact, relationship3, []*PathComponent{NewPathComponentKey("relationships")})
    ApplyAddChange(contact, relationship4, []*PathComponent{NewPathComponentKey("relationships")})
    ApplyAddChange(contact, relationship5, []*PathComponent{NewPathComponentKey("relationships")})
    if contact.Relationships == nil || len(contact.Relationships) != 5 || contact.Relationships[0] != relationship1 || contact.Relationships[1] != relationship2 || contact.Relationships[2] != relationship3 || contact.Relationships[3] != relationship4 || contact.Relationships[4] != relationship5 {
        t.Fatalf("Expected Relationship to be updated to %#v but was %#v", []*Relationship{relationship1, relationship2, relationship3, relationship4, relationship5}, contact.Relationships)
    }

    postalAddress1 := &PostalAddress{
        Address:       "White House Plantation, Pamunkey River, New Kent County, Virginia",
        StreetAddress: "White House Plantation",
        OtherAddress:  "Pamunkey River",
        Municipality:  "New Kent County",
        Region:        "Virginia",
        Country:       "us",
    }
    postalAddress1.Id = "asdfsdfsdfsf"

    postalAddress2 := &PostalAddress{
        Address:       "Mount Vernon Plantation, Alexandria, Virginia",
        StreetAddress: "Mount Vernon Plantation",
        Municipality:  "Alexandria",
        Region:        "Virginia",
        Country:       "us",
    }
    postalAddress2.Id = "ejhkhfsfjsjf"

    postalAddress3 := &PostalAddress{
        Address:       "1600 Pennsylvania Ave NW, Washington, DC 20500",
        StreetAddress: "1600 Pennsylvania Ave NW",
        Municipality:  "Washington",
        Region:        "DC",
        PostalCode:    "20500",
        Country:       "us",
    }
    postalAddress3.Id = "sfljsflsfjsfsjfjsjfjsf"

    addrRef1 := &ContactReference{ReferenceContactName: "John Adams"}
    addrRef2 := &ContactReference{ReferenceContactName: "John Hancock"}
    addrRef3 := &ContactReference{ReferenceContactName: "Thomas Jefferson"}
    addrRef4 := &ContactReference{ReferenceContactName: "Benjamin Franklin"}
    ApplyAddChange(contact, postalAddress1, []*PathComponent{NewPathComponentKey("postal_addresses")})
    ApplyAddChange(contact, postalAddress2, []*PathComponent{NewPathComponentKey("postal_addresses")})
    ApplyAddChange(contact, postalAddress3, []*PathComponent{NewPathComponentKey("postal_addresses")})
    ApplyAddChange(contact, addrRef1, []*PathComponent{NewPathComponentKey("postal_addresses"), NewPathComponentId(postalAddress3.Id), NewPathComponentKey("references")})
    ApplyAddChange(contact, addrRef2, []*PathComponent{NewPathComponentKey("postal_addresses"), NewPathComponentIndex("2"), NewPathComponentKey("references")})
    ApplyAddChange(contact, addrRef3, []*PathComponent{NewPathComponentKey("postal_addresses"), NewPathComponentId(postalAddress3.Id), NewPathComponentKey("references")})
    ApplyAddChange(contact, addrRef4, []*PathComponent{NewPathComponentKey("postal_addresses"), NewPathComponentId(postalAddress3.Id), NewPathComponentKey("references")})
    if contact.PostalAddresses == nil || len(contact.PostalAddresses) != 3 || contact.PostalAddresses[0].Id != postalAddress1.Id || contact.PostalAddresses[1].Id != postalAddress2.Id || contact.PostalAddresses[2].Id != postalAddress3.Id {
        t.Fatalf("Expected PostalAddresses to be updated to %#v but was %#v", []*PostalAddress{postalAddress1, postalAddress2, postalAddress3}, contact.PostalAddresses)
    }
    postalRefs := contact.PostalAddresses[2].References
    if postalRefs == nil || len(postalRefs) != 4 || postalRefs[0].ReferenceContactName != addrRef1.ReferenceContactName || postalRefs[1].ReferenceContactName != addrRef2.ReferenceContactName || postalRefs[2].ReferenceContactName != addrRef3.ReferenceContactName || postalRefs[3].ReferenceContactName != addrRef4.ReferenceContactName {
        t.Fatalf("Expected PostalAddresses[2].References to be updated to %#v but was %#v", []*ContactReference{addrRef1, addrRef2, addrRef3, addrRef4}, postalRefs)
    }
}

func TestApplyDeleteChange(t *testing.T) {
    contact := &Contact{
        DisplayName:     "Martha Parke Custis",
        GivenName:       "Martha",
        MiddleName:      "Parke",
        Surname:         "Custis",
        Nickname:        "Patsy",
        ExternalUserIds: []string{"blah1", "blah2", "blah3", "blah4"},
        Relationships: []*Relationship{
            &Relationship{Rel: REL_RELATIONSHIP_HUSBAND, ContactReferenceName: "George Washington"},
            &Relationship{Rel: REL_RELATIONSHIP_CHILD, ContactReferenceName: "Daniel Custis"},
            &Relationship{Rel: REL_RELATIONSHIP_CHILD, ContactReferenceName: "Francis Custis"},
            &Relationship{Rel: REL_RELATIONSHIP_CHILD, ContactReferenceName: "John Parke Custis"},
            &Relationship{Rel: REL_RELATIONSHIP_CHILD, ContactReferenceName: "Martha Parke Custis"},
        },
        PostalAddresses: []*PostalAddress{
            &PostalAddress{
                AclPersistableModel: AclPersistableModel{PersistableModel: PersistableModel{Id: "asdfsdfsdfsf"}},
                Address:             "White House Plantation, Pamunkey River, New Kent County, Virginia",
                StreetAddress:       "White House Plantation",
                OtherAddress:        "Pamunkey River",
                Municipality:        "New Kent County",
                Region:              "Virginia",
                Country:             "us",
            },
            &PostalAddress{
                AclPersistableModel: AclPersistableModel{PersistableModel: PersistableModel{Id: "ejhkhfsfjsjf"}},
                Address:             "Mount Vernon Plantation, Alexandria, Virginia",
                StreetAddress:       "Mount Vernon Plantation",
                Municipality:        "Alexandria",
                Region:              "Virginia",
                Country:             "us",
            },
            &PostalAddress{
                AclPersistableModel: AclPersistableModel{PersistableModel: PersistableModel{Id: "sfljsflsfjsfsjfjsjfjsf"}},
                Address:             "1600 Pennsylvania Ave NW, Washington, DC 20500",
                StreetAddress:       "1600 Pennsylvania Ave NW",
                Municipality:        "Washington",
                Region:              "DC",
                PostalCode:          "20500",
                Country:             "us",
                References: []*ContactReference{
                    &ContactReference{ReferenceContactName: "John Adams"},
                    &ContactReference{ReferenceContactName: "John Hancock"},
                    &ContactReference{ReferenceContactName: "Thomas Jefferson"},
                    &ContactReference{ReferenceContactName: "Benjamin Franklin"},
                },
            },
        },
    }

    ApplyDeleteChange(contact, nil, []*PathComponent{NewPathComponentKey("nickname")})
    if contact.Nickname != "" {
        t.Fatalf("Expected Nickname to be updated to %#v but was %#v", "", contact.Nickname)
    }

    ApplyDeleteChange(contact, nil, []*PathComponent{NewPathComponentKey("birthday")})
    if contact.Birthday != nil {
        t.Fatalf("Expected Birthday to be updated to %#v but was %#v", nil, contact.Birthday)
    }

    externalUserId1 := "blah1"
    externalUserId2 := "blah2"
    externalUserId3 := "blah3"
    externalUserId4 := "blah4"
    ApplyDeleteChange(contact, externalUserId3, []*PathComponent{NewPathComponentKey("external_user_ids")})
    if contact.ExternalUserIds == nil || len(contact.ExternalUserIds) != 3 || contact.ExternalUserIds[0] != externalUserId1 || contact.ExternalUserIds[1] != externalUserId2 || contact.ExternalUserIds[2] != externalUserId4 {
        t.Fatalf("Expected ExternalUserId to be updated to %#v but was %#v", []string{externalUserId1, externalUserId2, externalUserId4}, contact.ExternalUserIds)
    }
    ApplyDeleteChange(contact, externalUserId1, []*PathComponent{NewPathComponentKey("external_user_ids")})
    if contact.ExternalUserIds == nil || len(contact.ExternalUserIds) != 2 || contact.ExternalUserIds[0] != externalUserId2 || contact.ExternalUserIds[1] != externalUserId4 {
        t.Fatalf("Expected ExternalUserId to be updated to %#v but was %#v", []string{externalUserId2, externalUserId4}, contact.ExternalUserIds)
    }
    ApplyDeleteChange(contact, nil, []*PathComponent{NewPathComponentKey("external_user_ids"), NewPathComponentIndex("1")})
    if contact.ExternalUserIds == nil || len(contact.ExternalUserIds) != 1 || contact.ExternalUserIds[0] != externalUserId2 {
        t.Fatalf("Expected ExternalUserId to be updated to %#v but was %#v", []string{externalUserId2}, contact.ExternalUserIds)
    }
    ApplyDeleteChange(contact, externalUserId2, []*PathComponent{NewPathComponentKey("external_user_ids")})
    if contact.ExternalUserIds != nil && len(contact.ExternalUserIds) != 0 {
        t.Fatalf("Expected ExternalUserId to be updated to %#v but was %#v", []string{}, contact.ExternalUserIds)
    }

    relationship1 := &Relationship{Rel: REL_RELATIONSHIP_HUSBAND, ContactReferenceName: "George Washington"}
    relationship2 := &Relationship{Rel: REL_RELATIONSHIP_CHILD, ContactReferenceName: "Daniel Custis"}
    relationship3 := &Relationship{Rel: REL_RELATIONSHIP_CHILD, ContactReferenceName: "Francis Custis"}
    relationship4 := &Relationship{Rel: REL_RELATIONSHIP_CHILD, ContactReferenceName: "John Parke Custis"}
    relationship5 := &Relationship{Rel: REL_RELATIONSHIP_CHILD, ContactReferenceName: "Martha Parke Custis"}
    ApplyDeleteChange(contact, relationship1, []*PathComponent{NewPathComponentKey("relationships")})
    if contact.Relationships == nil || len(contact.Relationships) != 4 || contact.Relationships[0].ContactReferenceName != relationship2.ContactReferenceName || contact.Relationships[1].ContactReferenceName != relationship3.ContactReferenceName || contact.Relationships[2].ContactReferenceName != relationship4.ContactReferenceName || contact.Relationships[3].ContactReferenceName != relationship5.ContactReferenceName {
        t.Fatalf("Expected Relationship to be updated to %#v but was %#v", []*Relationship{relationship2, relationship3, relationship4, relationship5}, contact.Relationships)
    }
    ApplyDeleteChange(contact, relationship5, []*PathComponent{NewPathComponentKey("relationships"), NewPathComponentIndex("3")})
    if contact.Relationships == nil || len(contact.Relationships) != 3 || contact.Relationships[0].ContactReferenceName != relationship2.ContactReferenceName || contact.Relationships[1].ContactReferenceName != relationship3.ContactReferenceName || contact.Relationships[2].ContactReferenceName != relationship4.ContactReferenceName {
        t.Fatalf("Expected Relationship to be updated to %#v but was %#v", []*Relationship{relationship2, relationship3, relationship4}, contact.Relationships)
    }
    ApplyDeleteChange(contact, relationship2, []*PathComponent{NewPathComponentKey("relationships"), NewPathComponentIndex("0")})
    if contact.Relationships == nil || len(contact.Relationships) != 2 || contact.Relationships[0].ContactReferenceName != relationship3.ContactReferenceName || contact.Relationships[1].ContactReferenceName != relationship4.ContactReferenceName {
        t.Fatalf("Expected Relationship to be updated to %#v but was %#v", []*Relationship{relationship3, relationship4}, contact.Relationships)
    }
    ApplyDeleteChange(contact, relationship4, []*PathComponent{NewPathComponentKey("relationships")})
    if contact.Relationships == nil || len(contact.Relationships) != 1 || contact.Relationships[0].ContactReferenceName != relationship3.ContactReferenceName {
        t.Fatalf("Expected Relationship to be updated to %#v but was %#v", []*Relationship{relationship3}, contact.Relationships)
    }
    ApplyDeleteChange(contact, relationship3, []*PathComponent{NewPathComponentKey("relationships")})
    if contact.Relationships != nil && len(contact.Relationships) != 0 {
        t.Fatalf("Expected Relationship to be updated to %#v but was %#v", []*Relationship{}, contact.Relationships)
    }

    postalAddress1 := &PostalAddress{
        AclPersistableModel: AclPersistableModel{PersistableModel: PersistableModel{Id: "asdfsdfsdfsf"}},
        Address:             "White House Plantation, Pamunkey River, New Kent County, Virginia",
        StreetAddress:       "White House Plantation",
        OtherAddress:        "Pamunkey River",
        Municipality:        "New Kent County",
        Region:              "Virginia",
        Country:             "us",
    }

    postalAddress2 := &PostalAddress{
        AclPersistableModel: AclPersistableModel{PersistableModel: PersistableModel{Id: "ejhkhfsfjsjf"}},
        Address:             "Mount Vernon Plantation, Alexandria, Virginia",
        StreetAddress:       "Mount Vernon Plantation",
        Municipality:        "Alexandria",
        Region:              "Virginia",
        Country:             "us",
    }

    postalAddress3 := &PostalAddress{
        AclPersistableModel: AclPersistableModel{PersistableModel: PersistableModel{Id: "sfljsflsfjsfsjfjsjfjsf"}},
        Address:             "1600 Pennsylvania Ave NW, Washington, DC 20500",
        StreetAddress:       "1600 Pennsylvania Ave NW",
        Municipality:        "Washington",
        Region:              "DC",
        PostalCode:          "20500",
        Country:             "us",
    }

    addrRef1 := &ContactReference{ReferenceContactName: "John Adams"}
    addrRef2 := &ContactReference{ReferenceContactName: "John Hancock"}
    addrRef3 := &ContactReference{ReferenceContactName: "Thomas Jefferson"}
    addrRef4 := &ContactReference{ReferenceContactName: "Benjamin Franklin"}
    ApplyDeleteChange(contact, postalAddress2, []*PathComponent{NewPathComponentKey("postal_addresses"), NewPathComponentId(postalAddress2.Id)})
    if contact.PostalAddresses == nil || len(contact.PostalAddresses) != 2 || contact.PostalAddresses[0].Id != postalAddress1.Id || contact.PostalAddresses[1].Id != postalAddress3.Id {
        t.Fatalf("Expected PostalAddresses to be updated to %#v but was %#v", []*PostalAddress{postalAddress1, postalAddress3}, contact.PostalAddresses)
    }
    ApplyDeleteChange(contact, addrRef3, []*PathComponent{NewPathComponentKey("postal_addresses"), NewPathComponentId(postalAddress3.Id), NewPathComponentKey("references")})
    postalRefs := contact.PostalAddresses[1].References
    if postalRefs == nil || len(postalRefs) != 3 || postalRefs[0].ReferenceContactName != addrRef1.ReferenceContactName || postalRefs[1].ReferenceContactName != addrRef2.ReferenceContactName || postalRefs[2].ReferenceContactName != addrRef4.ReferenceContactName {
        t.Fatalf("1 Expected PostalAddresses[1].References to be updated to %#v but was %#v", []*ContactReference{addrRef1, addrRef2, addrRef4}, postalRefs)
    }
    ApplyDeleteChange(contact, addrRef4, []*PathComponent{NewPathComponentKey("postal_addresses"), NewPathComponentIndex("1"), NewPathComponentKey("references"), NewPathComponentIndex("2")})
    postalRefs = contact.PostalAddresses[1].References
    if postalRefs == nil || len(postalRefs) != 2 || postalRefs[0].ReferenceContactName != addrRef1.ReferenceContactName || postalRefs[1].ReferenceContactName != addrRef2.ReferenceContactName {
        t.Fatalf("2 Expected PostalAddresses[1].References to be updated to %#v but was %#v", []*ContactReference{addrRef1, addrRef2}, postalRefs)
    }
    ApplyDeleteChange(contact, addrRef1, []*PathComponent{NewPathComponentKey("postal_addresses"), NewPathComponentId(postalAddress3.Id), NewPathComponentKey("references")})
    postalRefs = contact.PostalAddresses[1].References
    if postalRefs == nil || len(postalRefs) != 1 || postalRefs[0].ReferenceContactName != addrRef2.ReferenceContactName {
        t.Fatalf("3 Expected PostalAddresses[1].References to be updated to %#v but was %#v", []*ContactReference{addrRef2}, postalRefs)
    }
    ApplyDeleteChange(contact, addrRef2, []*PathComponent{NewPathComponentKey("postal_addresses"), NewPathComponentId(postalAddress3.Id), NewPathComponentKey("references")})
    postalRefs = contact.PostalAddresses[1].References
    if postalRefs != nil && len(postalRefs) != 0 {
        t.Fatalf("4 Expected PostalAddresses[1].References to be updated to %#v but was %#v", []*ContactReference{}, postalRefs)
    }

    ApplyDeleteChange(contact, postalAddress3, []*PathComponent{NewPathComponentKey("postal_addresses")})
    if contact.PostalAddresses == nil || len(contact.PostalAddresses) != 1 || contact.PostalAddresses[0].Id != postalAddress1.Id {
        t.Fatalf("Expected PostalAddresses to be updated to %#v but was %#v", []*PostalAddress{postalAddress1}, contact.PostalAddresses)
    }
    ApplyDeleteChange(contact, postalAddress1, []*PathComponent{NewPathComponentKey("postal_addresses")})
    if contact.PostalAddresses != nil && len(contact.PostalAddresses) != 0 {
        t.Fatalf("Expected PostalAddresses to be updated to %#v but was %#v", []*PostalAddress{}, contact.PostalAddresses)
    }

}

func TestApplyUpdateChange(t *testing.T) {
    contact := &Contact{
        DisplayName:     "Martha Parke Custis",
        GivenName:       "Martha",
        MiddleName:      "Parke",
        Surname:         "Custis",
        Nickname:        "Patsie",
        Birthday:        &Date{Year: 2010, Month: 5, Day: 7},
        ExternalUserIds: []string{"blah1", "blah2", "blah3", "blah4"},
        Relationships: []*Relationship{
            &Relationship{Rel: REL_RELATIONSHIP_HUSBAND, ContactReferenceName: "Daniel Parke Custis"},
            &Relationship{Rel: REL_RELATIONSHIP_CHILD, ContactReferenceName: "Daniel Custis"},
            &Relationship{Rel: REL_RELATIONSHIP_CHILD, ContactReferenceName: "Francis Custis"},
            &Relationship{Rel: REL_RELATIONSHIP_CHILD, ContactReferenceName: "John Parke Custis"},
            &Relationship{Rel: REL_RELATIONSHIP_CHILD, ContactReferenceName: "Martha Parke Custis"},
        },
        PostalAddresses: []*PostalAddress{
            &PostalAddress{
                AclPersistableModel: AclPersistableModel{PersistableModel: PersistableModel{Id: "asdfsdfsdfsf"}},
                Address:             "White House Plantation",
                StreetAddress:       "White House Plantation",
                OtherAddress:        "Pamunkey River",
                Municipality:        "New Kent County",
                Region:              "Virginia",
                Country:             "us",
            },
            &PostalAddress{
                AclPersistableModel: AclPersistableModel{PersistableModel: PersistableModel{Id: "ejhkhfsfjsjf"}},
                Address:             "Mount Vernon Plantation, Alexandria, Virginia",
                StreetAddress:       "Mount Vernon PLANTATION",
                Municipality:        "Alexandria",
                Region:              "Virginia",
                Country:             "us",
            },
            &PostalAddress{
                AclPersistableModel: AclPersistableModel{PersistableModel: PersistableModel{Id: "sfljsflsfjsfsjfjsjfjsf"}},
                Address:             "1600 NW Pennsylvania Ave, Washington, DC 20500",
                StreetAddress:       "1600 NW Pennsylvania Ave",
                Municipality:        "Washington, DC",
                Region:              "D.C.",
                PostalCode:          "20500",
                Country:             "us",
                References: []*ContactReference{
                    &ContactReference{ReferenceContactName: "James Adams"},
                    &ContactReference{ReferenceContactName: "John Hancock"},
                    &ContactReference{ReferenceContactName: "Thomas Jefferson"},
                    &ContactReference{ReferenceContactName: "Benjamin Franklin"},
                },
            },
        },
    }

    ApplyUpdateChange(contact, contact.Nickname, "Patsy", []*PathComponent{NewPathComponentKey("nickname")})
    if contact.Nickname != "Patsy" {
        t.Fatalf("Expected Nickname to be updated to %#v but was %#v", "", contact.Nickname)
    }

    ApplyUpdateChange(contact, &Date{Year: 2010, Month: 5, Day: 7}, &Date{Year: 1731, Month: 6, Day: 2}, []*PathComponent{NewPathComponentKey("birthday")})
    if contact.Birthday == nil || contact.Birthday.Year != 1731 || contact.Birthday.Month != 6 || contact.Birthday.Day != 2 {
        t.Fatalf("Expected Birthday to be updated to %#v but was %#v", &Date{Year: 1731, Month: 6, Day: 2}, contact.Birthday)
    }

    externalUserId1 := "hiho1"
    externalUserId2 := "hiho2"
    externalUserId3 := "hiho3"
    externalUserId4 := "hiho4"
    ApplyUpdateChange(contact, contact.ExternalUserIds[2], externalUserId3, []*PathComponent{NewPathComponentKey("external_user_ids")})
    if contact.ExternalUserIds == nil || len(contact.ExternalUserIds) != 4 || contact.ExternalUserIds[2] != externalUserId3 {
        t.Fatalf("Expected ExternalUserId[2] to be updated to %#v but was %#v", externalUserId3, contact.ExternalUserIds[2])
    }
    ApplyUpdateChange(contact, contact.ExternalUserIds[0], externalUserId1, []*PathComponent{NewPathComponentKey("external_user_ids"), NewPathComponentIndex("0")})
    if contact.ExternalUserIds == nil || len(contact.ExternalUserIds) != 4 || contact.ExternalUserIds[0] != externalUserId1 {
        t.Fatalf("Expected ExternalUserId[0] to be updated to %#v but was %#v", externalUserId1, contact.ExternalUserIds[0])
    }
    ApplyUpdateChange(contact, contact.ExternalUserIds[1], externalUserId2, []*PathComponent{NewPathComponentKey("external_user_ids"), NewPathComponentIndex("1")})
    if contact.ExternalUserIds == nil || len(contact.ExternalUserIds) != 4 || contact.ExternalUserIds[1] != externalUserId2 {
        t.Fatalf("Expected ExternalUserId[1] to be updated to %#v but was %#v", externalUserId2, contact.ExternalUserIds[1])
    }
    ApplyUpdateChange(contact, contact.ExternalUserIds[3], externalUserId4, []*PathComponent{NewPathComponentKey("external_user_ids")})
    if contact.ExternalUserIds == nil || len(contact.ExternalUserIds) != 4 || contact.ExternalUserIds[3] != externalUserId4 {
        t.Fatalf("Expected ExternalUserId[3] to be updated to %#v but was %#v", externalUserId4, contact.ExternalUserIds[3])
    }

    relationship1 := &Relationship{Rel: REL_RELATIONSHIP_HUSBAND, ContactReferenceName: "George Washington"}
    ApplyUpdateChange(contact, contact.Relationships[0].ContactReferenceName, relationship1.ContactReferenceName, []*PathComponent{NewPathComponentKey("relationships"), NewPathComponentIndex("0"), NewPathComponentKey("contact_reference_name")})
    if contact.Relationships == nil || len(contact.Relationships) != 5 || contact.Relationships[0].ContactReferenceName != relationship1.ContactReferenceName {
        t.Fatalf("Expected Relationships[0].ContactReferenceName to be updated to %#v but was %#v with length %d", relationship1.ContactReferenceName, contact.Relationships[0].ContactReferenceName, len(contact.Relationships))
    }

    postalAddress1 := &PostalAddress{
        AclPersistableModel: AclPersistableModel{PersistableModel: PersistableModel{Id: "asdfsdfsdfsf"}},
        Address:             "White House Plantation, Pamunkey River, New Kent County, Virginia",
        StreetAddress:       "White House Plantation",
        OtherAddress:        "Pamunkey River",
        Municipality:        "New Kent County",
        Region:              "Virginia",
        Country:             "us",
    }

    postalAddress2 := &PostalAddress{
        AclPersistableModel: AclPersistableModel{PersistableModel: PersistableModel{Id: "ejhkhfsfjsjf"}},
        Address:             "Mount Vernon Plantation, Alexandria, Virginia",
        StreetAddress:       "Mount Vernon Plantation",
        Municipality:        "Alexandria",
        Region:              "Virginia",
        Country:             "us",
    }

    postalAddress3 := &PostalAddress{
        AclPersistableModel: AclPersistableModel{PersistableModel: PersistableModel{Id: "sfljsflsfjsfsjfjsjfjsf"}},
        Address:             "1600 Pennsylvania Ave NW, Washington, DC 20500",
        StreetAddress:       "1600 Pennsylvania Ave NW",
        Municipality:        "Washington",
        Region:              "DC",
        PostalCode:          "20500",
        Country:             "us",
        References: []*ContactReference{
            &ContactReference{ReferenceContactName: "John Adams"},
            &ContactReference{ReferenceContactName: "John Hancock"},
            &ContactReference{ReferenceContactName: "Thomas Jefferson"},
            &ContactReference{ReferenceContactName: "Benjamin Franklin"},
        },
    }

    ApplyUpdateChange(contact, contact.PostalAddresses[0].Address, postalAddress1.Address, []*PathComponent{NewPathComponentKey("postal_addresses"), NewPathComponentId(postalAddress1.Id), NewPathComponentKey("address")})
    if contact.PostalAddresses == nil || len(contact.PostalAddresses) != 3 || contact.PostalAddresses[0].Id != postalAddress1.Id || contact.PostalAddresses[1].Id != postalAddress2.Id || contact.PostalAddresses[2].Id != postalAddress3.Id {
        t.Fatalf("Expected PostalAddresses to be updated to %#v but was %#v", []*PostalAddress{postalAddress1, postalAddress2, postalAddress3}, contact.PostalAddresses)
    }
    if contact.PostalAddresses[0].Address != postalAddress1.Address {
        t.Fatalf("Expected PostalAddresses[0] to be updated to %#v but was %#v", postalAddress1, contact.PostalAddresses[0])
    }

    ApplyUpdateChange(contact, contact.PostalAddresses[1].StreetAddress, postalAddress2.StreetAddress, []*PathComponent{NewPathComponentKey("postal_addresses"), NewPathComponentIndex("1"), NewPathComponentKey("street_address")})
    if contact.PostalAddresses[0].StreetAddress != postalAddress1.StreetAddress {
        t.Fatalf("Expected PostalAddresses[0] to be updated to %#v but was %#v", postalAddress2, contact.PostalAddresses[1])
    }

    ApplyUpdateChange(contact, contact.PostalAddresses[2].Address, postalAddress3.Address, []*PathComponent{NewPathComponentKey("postal_addresses"), NewPathComponentId(postalAddress3.Id), NewPathComponentKey("address")})
    ApplyUpdateChange(contact, contact.PostalAddresses[2].StreetAddress, postalAddress3.StreetAddress, []*PathComponent{NewPathComponentKey("postal_addresses"), NewPathComponentIndex("2"), NewPathComponentKey("street_address")})
    ApplyUpdateChange(contact, contact.PostalAddresses[2].Municipality, postalAddress3.Municipality, []*PathComponent{NewPathComponentKey("postal_addresses"), NewPathComponentId(postalAddress3.Id), NewPathComponentKey("municipality")})
    ApplyUpdateChange(contact, contact.PostalAddresses[2].Region, postalAddress3.Region, []*PathComponent{NewPathComponentKey("postal_addresses"), NewPathComponentId(postalAddress3.Id), NewPathComponentKey("region")})
    ApplyUpdateChange(contact, contact.PostalAddresses[2].References[0].ReferenceContactName, postalAddress3.References[0].ReferenceContactName, []*PathComponent{NewPathComponentKey("postal_addresses"), NewPathComponentId(postalAddress3.Id), NewPathComponentKey("references"), NewPathComponentIndex("0"), NewPathComponentKey("reference_contact_name")})
    if contact.PostalAddresses[2].References[0].ReferenceContactName != postalAddress3.References[0].ReferenceContactName {
        t.Fatalf("Expected PostalAddresses[2].References[0].ReferenceContactName to be updated to %#v but was %#v", postalAddress3.References[0], contact.PostalAddresses[2].References[0])
    }
}
