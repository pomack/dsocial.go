package dsocial

import (
    "fmt"
    "github.com/pomack/jsonhelper.go/jsonhelper"
    "strings"
    "time"
)

type Date struct {
    Year  int16 `json:"year,omitempty"`
    Month int8  `json:"month,omitempty"`
    Day   int8  `json:"day,omitempty"`
}

func (p *Date) String() string {
    return fmt.Sprintf("%04d-%02d-%02d", p.Year, p.Month, p.Day)
}
func (p *Date) Equals(d *Date) bool {
    return d != nil && p.Year == d.Year && p.Month == d.Month && p.Day == d.Day
}
func (p *Date) IsEmpty() bool {
    return p == nil || (p.Year == 0 && p.Month == 0 && p.Day == 0)
}

type DateTime struct {
    Year   int16 `json:"year,omitempty"`
    Month  int8  `json:"month,omitempty"`
    Day    int8  `json:"day,omitempty"`
    Hour   int8  `json:"hour,omitempty"`
    Minute int8  `json:"minute,omitempty"`
    Second int8  `json:"second,omitempty"`
}

func (p *DateTime) String() string {
    return fmt.Sprintf("%04d-%02d-%02dT%02d:%02d:%02dZ", p.Year, p.Month, p.Day, p.Hour, p.Minute, p.Second)
}
func (p *DateTime) Equals(d *DateTime) bool {
    return d != nil && p.Year == d.Year && p.Month == d.Month && p.Day == d.Day && p.Hour == d.Hour && p.Minute == d.Minute && p.Second == d.Second
}
func (p *DateTime) IsEmpty() bool {
    return p == nil || (p.Year == 0 && p.Month == 0 && p.Day == 0 && p.Hour == 0 && p.Minute == 0 && p.Second == 0)
}

type Location struct {
    Latitude  float64 `json:"latitude,omitempty"`
    Longitude float64 `json:"longitude,omitempty"`
    Elevation float64 `json:"elevation,omitempty"`
}

type PersistableModel struct {
    Id         string `json:"id,omitempty"`
    Etag       string `json:"etag,omitempty"`
    CreatedAt  int64  `json:"created_at,omitempty"`
    ModifiedAt int64  `json:"modified_at,omitempty"`
}

func (p *PersistableModel) InitFromJSONObject(obj jsonhelper.JSONObject) {
    p.Id = obj.GetAsString("id")
    p.Etag = obj.GetAsString("etag")
    p.CreatedAt = obj.GetAsInt64("created_at")
    p.ModifiedAt = obj.GetAsInt64("modified_at")
}

func (p *PersistableModel) CleanFromUser(user *User, original *PersistableModel) {
    if original == nil {
        p.CreatedAt = 0
        p.ModifiedAt = 0
    } else {
        p.Id = original.Id
        p.CreatedAt = original.CreatedAt
        p.ModifiedAt = original.ModifiedAt
    }
}

func (p *PersistableModel) Validate(createNew bool, errors map[string][]error) (isValid bool) {
    if errors == nil {
        errors = make(map[string][]error)
    }
    p.Id, _ = validateId(p.Id, createNew, "id", errors)
    isValid = len(errors) == 0
    return
}

func (p *PersistableModel) BeforeCreate() error {
    p.Etag = GenerateEtag()
    p.CreatedAt = time.Now().Unix()
    p.ModifiedAt = p.CreatedAt
    return nil
}

func (p *PersistableModel) BeforeUpdate() error {
    p.Etag = GenerateEtag()
    p.ModifiedAt = time.Now().Unix()
    return nil
}

func (p *PersistableModel) BeforeSave() error {
    return nil
}

func (p *PersistableModel) BeforeDelete() error {
    return nil
}

func (p *PersistableModel) AfterCreate() error {
    return nil
}

func (p *PersistableModel) AfterUpdate() error {
    return nil
}

func (p *PersistableModel) AfterSave() error {
    return nil
}

func (p *PersistableModel) AfterDelete() error {
    return nil
}

func removeEmptyStrings(arr []string) []string {
    sv := make([]string, 0, len(arr))
    for _, s := range arr {
        if s != "" {
            sv = append(sv, s)
        }
    }
    return sv
}

func join(sep string, values ...string) string {
    return strings.Join(removeEmptyStrings(values), sep)
}

func addIfNonSpaces(sv []string, s string) []string {
    p := strings.TrimSpace(s)
    if s != "" {
        sv = append(sv, p)
    }
    return sv
}

func ParsePhoneNumber(s string, number *PhoneNumber) {
    slen := len(s)
    if slen == 0 {
        return
    }
    number.FormattedNumber = s
    sv := make([]string, 0, 2*slen)
    start := 0
    for i := 0; i < slen; i++ {
        b := s[i]
        switch b {
        case '-', ' ', '.', ')', '_':
            if start < i {
                sv = addIfNonSpaces(sv, s[start:i])
            }
            start = i + 1
        case '(':
            for j := i + 1; j < slen; j++ {
                if s[j] == ')' {
                    if start < i {
                        sv = addIfNonSpaces(sv, s[start:i])
                    }
                    sv = addIfNonSpaces(sv, s[i+1:j])
                    i = j
                    start = j + 1
                    break
                }
            }
        case 'x', 'X':
            if start < i {
                lastChar := s[i-1]
                if lastChar == 'e' || lastChar == 'E' {
                    continue
                }
                sv = addIfNonSpaces(sv, s[start:i])
            }
            number.ExtensionNumber = strings.TrimSpace(s[i+1:])
            start = slen
            i = slen
            break
        case 't', 'T':
            l := i - start
            if l > 2 && strings.ToLower(s[i-2:i+1]) == "ext" {
                if start < i-2 {
                    sv = addIfNonSpaces(sv, s[start:i-2])
                }
            }
            number.ExtensionNumber = strings.TrimSpace(s[i+1:])
            start = slen
            i = slen
            break
        }
    }
    if start < slen {
        sv = addIfNonSpaces(sv, s[start:])
    }
    parts := sv
    if len(sv) == 0 {
        return
    }
    at := 0
    l := len(parts)
    if parts[0][0] == '+' {
        // preferred format is +<country code>-<area code>-<local code>
        number.CountryCode = parts[0][1:]
        at++
    } else if len(parts[0]) < 3 {
        // make an assumption that this is a country code
        number.CountryCode = parts[0]
        at++
    }
    if l > at+1 {
        number.AreaCode = parts[at]
        at++
    }
    if l > at {
        number.LocalPhoneNumber = strings.Join(parts[at:], "-")
    }
}

func ParseName(s string, c *Contact) {
    c.DisplayName = s
    nameParts := removeEmptyStrings(strings.Split(s, " "))
    switch len(nameParts) {
    case 0:
    case 1:
        c.GivenName = nameParts[0]
    case 2:
        c.GivenName = nameParts[0]
        c.Surname = nameParts[1]
    case 3:
        c.GivenName = nameParts[0]
        c.MiddleName = nameParts[1]
        c.Surname = nameParts[2]
    default:
        c.GivenName = nameParts[0]
        c.MiddleName = strings.Join(nameParts[1:len(nameParts)-1], " ")
        c.Surname = nameParts[len(nameParts)-1]
    }
}

func GenerateEtag() string {
    return generateSalt(24)
}
