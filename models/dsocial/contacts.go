package dsocial

type ContactNameOrdering string
type RelPostalAddress string
type RelEducation string
type RelPhoneNumber string
type RelEmail string
type RelUri string
type RelIM string
type RelIMProtocol string
type RelRelationship string
type RelDate string
type RelDateTime string
type RelGender string
type RelRelationshipStatus string

type ContactReference struct {
    AclPersistableModel  `json:"model,omitempty,collapse"`
    UserId               string `json:"user_id,omitempty"`
    ContactId            string `json:"contact_id,omitempty"`
    ReferenceContactId   string `json:"reference_contact_id,omitempty"`
    ReferenceContactName string `json:"reference_contact_name,omitempty"`
    Text                 string `json:"text,omitempty"`
}

type PostalAddress struct {
    AclPersistableModel `json:"model,omitempty,collapse"`
    Address             string              `json:"address,omitempty"`
    Label               string              `json:"label,omitempty"`
    Rel                 RelPostalAddress    `json:"rel,omitempty"`
    StreetAddress       string              `json:"street_address,omitempty"`
    OtherAddress        string              `json:"other_address,omitempty"`
    Municipality        string              `json:"municipality,omitempty"`
    Region              string              `json:"region,omitempty"`
    PostalCode          string              `json:"postal_code,omitempty"`
    Country             string              `json:"country,omitempty"`
    LocatedFrom         *Date               `json:"located_from,omitempty"`
    LocatedTill         *Date               `json:"located_till,omitempty"`
    IsCurrent           bool                `json:"is_current,omitempty"`
    IsPrimary           bool                `json:"is_primary,omitempty"`
    References          []*ContactReference `json:"references,omitempty"`
}

type Degree struct {
    PersistableModel `json:"model,omitempty,collapse"`
    Degree           string `json:"degree,omitempty"`
    Major            string `json:"major,omitempty"`
}

type Education struct {
    AclPersistableModel `json:"model,omitempty,collapse"`
    Label               string              `json:"label,omitempty"`
    Rel                 RelEducation        `json:"rel,omitempty"`
    Degrees             []*Degree           `json:"degrees,omitempty"`
    Minors              []string            `json:"minors,omitempty"`
    GraduationYear      int16               `json:"graduation_year,omitempty"`
    Institution         string              `json:"institution,omitempty"`
    AttendedFrom        *Date               `json:"attended_from,omitempty"`
    AttendedTill        *Date               `json:"attended_till,omitempty"`
    Gpa                 float64             `json:"gpa,omitempty"`
    MajorGpa            float64             `json:"major_gpa,omitempty"`
    IsCurrent           bool                `json:"is_current,omitempty"`
    Graduated           bool                `json:"graduated,omitempty"`
    References          []*ContactReference `json:"references,omitempty"`
    Notes               string              `json:"notes,omitempty"`
    Activities          []string            `json:"activities,omitempty"`
}

type WorkPosition struct {
    AclPersistableModel `json:"model,omitempty,collapse"`
    Title               string              `json:"title,omitempty"`
    Department          string              `json:"department,omitempty"`
    Location            string              `json:"location,omitempty"`
    From                *Date               `json:"from,omitempty"`
    To                  *Date               `json:"to,omitempty"`
    IsCurrent           bool                `json:"is_current,omitempty"`
    Description         string              `json:"description,omitempty"`
    References          []*ContactReference `json:"references,omitempty"`
}

type WorkHistory struct {
    AclPersistableModel `json:"model,omitempty,collapse"`
    Company             string          `json:"company,omitempty"`
    From                *Date           `json:"from,omitempty"`
    To                  *Date           `json:"to,omitempty"`
    IsCurrent           bool            `json:"is_current,omitempty"`
    Description         string          `json:"description,omitempty"`
    Positions           []*WorkPosition `json:"positions,omitempty"`
}

type PhoneNumber struct {
    AclPersistableModel `json:"model,omitempty,collapse"`
    Label               string         `json:"label,omitempty"`
    Rel                 RelPhoneNumber `json:"rel,omitempty"`
    FormattedNumber     string         `json:"formatted_number,omitempty"`
    CountryCode         string         `json:"country_code,omitempty"`
    AreaCode            string         `json:"area_code,omitempty"`
    LocalPhoneNumber    string         `json:"local_phone_number,omitempty"`
    ExtensionNumber     string         `json:"extension_number,omitempty"`
    IsPrimary           bool           `json:"is_primary,omitempty"`
}

type Email struct {
    AclPersistableModel `json:"model,omitempty,collapse"`
    Label               string   `json:"label,omitempty"`
    Rel                 RelEmail `json:"rel,omitempty"`
    EmailAddress        string   `json:"email_address,omitempty"`
    IsPrimary           bool     `json:"is_primary,omitempty"`
}

type Uri struct {
    AclPersistableModel `json:"model,omitempty,collapse"`
    Label               string `json:"label,omitempty"`
    Rel                 RelUri `json:"rel,omitempty"`
    Uri                 string `json:"uri,omitempty"`
    IsPrimary           bool   `json:"is_primary,omitempty"`
}

type IM struct {
    AclPersistableModel `json:"model,omitempty,collapse"`
    Label               string        `json:"label,omitempty"`
    Rel                 RelIM         `json:"rel,omitempty"`
    Protocol            RelIMProtocol `json:"protocol,omitempty"`
    Handle              string        `json:"handle,omitempty"`
    IsPrimary           bool          `json:"is_primary,omitempty"`
}

type Relationship struct {
    AclPersistableModel  `json:"model,omitempty,collapse"`
    Label                string          `json:"label,omitempty"`
    Rel                  RelRelationship `json:"rel,omitempty"`
    ContactReferenceId   string          `json:"contact_reference_id,omitempty"`
    ContactReferenceName string          `json:"contact_reference_name,omitempty"`
}

type ContactDate struct {
    AclPersistableModel `json:"model,omitempty,collapse"`
    Label               string  `json:"label,omitempty"`
    Rel                 RelDate `json:"rel,omitempty"`
    Value               *Date   `json:"value,omitempty"`
    IsPrimary           bool    `json:"is_primary,omitempty"`
}

type ContactDateTime struct {
    AclPersistableModel `json:"model,omitempty,collapse"`
    Label               string      `json:"label,omitempty"`
    Rel                 RelDateTime `json:"rel,omitempty"`
    Value               *DateTime   `json:"value,omitempty"`
    IsPrimary           bool        `json:"is_primary,omitempty"`
}

type Group struct {
    AclPersistableModel `json:"model,omitempty,collapse"`
    UserId              string   `json:"user_id,omitempty"`
    Name                string   `json:"name,omitempty"`
    Description         string   `json:"description,omitempty"`
    ContactIds          []string `json:"contact_ids,omitempty"`
    ContactNames        []string `json:"contact_names,omitempty"`
}

type Certification struct {
    AclPersistableModel `json:"model,omitempty,collapse"`
    Name                string `json:"name,omitempty"`
    Authority           string `json:"authority,omitempty"`
    Number              string `json:"number,omitempty"`
    AsOf                *Date  `json:"as_of,omitempty"`
    ValidTill           *Date  `json:"valid_till,omitempty"`
}

type Skill struct {
    AclPersistableModel `json:"model,omitempty,collapse"`
    Name                string `json:"name,omitempty"`
    Proficiency         string `json:"proficiency,omitempty"`
}

type Language struct {
    AclPersistableModel `json:"model,omitempty,collapse"`
    Name                string `json:"name,omitempty"`
    ReadGradeLevel      int    `json:"read_grade_level,omitempty"`
    WriteGradeLevel     int    `json:"write_grade_level,omitempty"`
}

type Contact struct {
    AclPersistableModel `json:"model,omitempty,collapse"`
    Label               string                `json:"label,omitempty"`
    UserId              string                `json:"user_id,omitempty"`
    Prefix              string                `json:"prefix,omitempty"`
    GivenName           string                `json:"given_name,omitempty"`
    MiddleName          string                `json:"middle_name,omitempty"`
    Surname             string                `json:"surname,omitempty"`
    Suffix              string                `json:"suffix,omitempty"`
    MaidenName          string                `json:"maiden_name,omitempty"`
    DisplayName         string                `json:"display_name,omitempty"`
    Nickname            string                `json:"nickname,omitempty"`
    DisplayNameOrdering ContactNameOrdering   `json:"display_name_ordering,omitempty"`
    SortNameOrdering    ContactNameOrdering   `json:"sort_name_ordering,omitempty"`
    Hometown            string                `json:"hometown,omitempty"`
    Gender              RelGender             `json:"gender,omitempty"`
    Biography           string                `json:"biography,omitempty"`
    FavoriteQuotes      string                `json:"favorite_quotes,omitempty"`
    RelationshipStatus  RelRelationshipStatus `json:"relationship_status,omitempty"`
    IsOrganization      bool                  `json:"is_organization,omitempty"`
    Title               string                `json:"title,omitempty"`
    Company             string                `json:"company,omitempty"`
    Department          string                `json:"department,omitempty"`
    Municipality        string                `json:"municipality,omitempty"`
    Region              string                `json:"region,omitempty"`
    PostalCode          string                `json:"postal_code,omitempty"`
    CountryCode         string                `json:"country_code,omitempty"`
    Birthday            *Date                 `json:"birthday,omitempty"`
    Anniversary         *Date                 `json:"anniversary,omitempty"`
    Death               *Date                 `json:"death,omitempty"`
    PrimaryAddress      string                `json:"primary_address,omitempty"`
    PrimaryEmail        string                `json:"primary_email,omitempty"`
    PrimaryPhoneNumber  string                `json:"primary_phone_number,omitempty"`
    PrimaryUri          string                `json:"primary_uri,omitempty"`
    PrimaryIm           string                `json:"primary_im,omitempty"`
    Notes               string                `json:"notes,omitempty"`
    ThumbnailUrl        string                `json:"thumbnail_url,omitempty"`
    InternalUserIds     []string              `json:"internal_user_ids,omitempty"`
    ExternalUserIds     []string              `json:"external_user_ids,omitempty"`
    GroupNames          []string              `json:"group_names,omitempty"`
    PostalAddresses     []*PostalAddress      `json:"postal_addresses,omitempty"`
    Educations          []*Education          `json:"educations,omitempty"`
    WorkHistories       []*WorkHistory        `json:"work_histories,omitempty"`
    PhoneNumbers        []*PhoneNumber        `json:"phone_numbers,omitempty"`
    EmailAddresses      []*Email              `json:"email_addresses,omitempty"`
    Uris                []*Uri                `json:"uris,omitempty"`
    Ims                 []*IM                 `json:"ims,omitempty"`
    Relationships       []*Relationship       `json:"relationships,omitempty"`
    Dates               []*ContactDate        `json:"dates,omitempty"`
    DateTimes           []*ContactDateTime    `json:"datetimes,omitempty"`
    Certifications      []*Certification      `json:"certifications,omitempty"`
    Skills              []*Skill              `json:"skills,omitempty"`
    Languages           []*Language           `json:"languages,omitempty"`
}

func (p ContactNameOrdering) String() (s string) {
    switch p {
    case GIVEN_MIDDLE_SURNAME:
        s = `given,middle,surname`
    case SURNAME_GIVEN_MIDDLE:
        s = `surname,given,middle`
    case GIVEN_SURNAME_MIDDLE:
        s = `given,surname,middle`
    case SURNAME_MIDDLE_GIVEN:
        s = `surname,given,middle`
    case MIDDLE_GIVEN_SURNAME:
        s = `middle,given,surname`
    case MIDDLE_SURNAME_GIVEN:
        s = `middle,surname,given`
    default:
        s = `given,middle,surname`
    }
    return
}

func ToContactNameOrdering(s string) (p ContactNameOrdering) {
    switch s {
    case `given,middle,surname`:
        p = GIVEN_MIDDLE_SURNAME
    case `surname,given,middle`:
        p = SURNAME_GIVEN_MIDDLE
    case `given,surname,middle`:
        p = GIVEN_SURNAME_MIDDLE
    case `surname,middle,given`:
        p = SURNAME_MIDDLE_GIVEN
    case `middle,given,surname`:
        p = MIDDLE_GIVEN_SURNAME
    case `middle,surname,given`:
        p = MIDDLE_SURNAME_GIVEN
    default:
        p = GIVEN_MIDDLE_SURNAME
    }
    return
}

func (p ContactNameOrdering) MarshalJSON() (s string) {
    switch p {
    case GIVEN_MIDDLE_SURNAME:
        s = `given,middle,surname`
    case SURNAME_GIVEN_MIDDLE:
        s = `surname,given,middle`
    case GIVEN_SURNAME_MIDDLE:
        s = `given,surname,middle`
    case SURNAME_MIDDLE_GIVEN:
        s = `surname,given,middle`
    case MIDDLE_GIVEN_SURNAME:
        s = `middle,given,surname`
    case MIDDLE_SURNAME_GIVEN:
        s = `middle,surname,given`
    default:
        s = `given,middle,surname`
    }
    return
}
