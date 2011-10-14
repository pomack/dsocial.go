package dsocial_test

import (
    . "github.com/pomack/contacts.go/dsocial"
    "testing"
)

func testParsePhoneNumber(t *testing.T, formatted, country, area, local, ext string) {
    ph := new(PhoneNumber)
    v := formatted
    ParsePhoneNumber(v, ph)
    if ph.FormattedNumber != formatted || ph.CountryCode != country || ph.AreaCode != area || ph.LocalPhoneNumber != local || ph.ExtensionNumber != ext {
        t.Fatalf("Unable to parse phone number \"%s\", parsed as: %#v", formatted, ph)
    }
}

func TestParsePhoneNumberPreferredUS(t *testing.T) {
    testParsePhoneNumber(t, "+1-234-567-8900", "1", "234", "567-8900", "")
}

func TestParsePhoneNumberPreferredInternational(t *testing.T) {
    testParsePhoneNumber(t, "+1340-234-567-8900", "1340", "234", "567-8900", "")
}

func TestParsePhoneNumberPreferredUSWithExtension(t *testing.T) {
    testParsePhoneNumber(t, "+1-234-567-8900x32-418", "1", "234", "567-8900", "32-418")
}

func TestParsePhoneNumberPreferredInternationalWithExtension(t *testing.T) {
    testParsePhoneNumber(t, "+1340-234-567-8900 ext 17 145", "1340", "234", "567-8900", "17 145")
}

func TestParsePhoneNumberUSParentheses(t *testing.T) {
    testParsePhoneNumber(t, "(234) 567-8900", "", "234", "567-8900", "")
}

func TestParsePhoneNumberUSDashes(t *testing.T) {
    testParsePhoneNumber(t, "234-567-8900", "", "234", "567-8900", "")
}

func TestParsePhoneNumberUSDots(t *testing.T) {
    testParsePhoneNumber(t, "234.567.8900", "", "234", "567-8900", "")
}

func TestParsePhoneNumberUSDotsWithCountryCode(t *testing.T) {
    testParsePhoneNumber(t, "1.234.567.8900", "1", "234", "567-8900", "")
}

func TestParsePhoneNumberGermanyDIN(t *testing.T) {
    testParsePhoneNumber(t, "0AAAA BBBBBB", "", "0AAAA", "BBBBBB", "")
}

func TestParsePhoneNumberGermanyDINWithExtension(t *testing.T) {
    testParsePhoneNumber(t, "0AAAA BBBBBB-EE", "", "0AAAA", "BBBBBB-EE", "")
}

func TestParsePhoneNumberGermanyDINWithInternational(t *testing.T) {
    testParsePhoneNumber(t, "+49 0AAAA BBBBBB", "49", "0AAAA", "BBBBBB", "")
}

func TestParsePhoneNumberGermanyE123Local(t *testing.T) {
    testParsePhoneNumber(t, "(0AAAA) BBBBBB", "", "0AAAA", "BBBBBB", "")
}

func TestParsePhoneNumberGermanyE123International(t *testing.T) {
    testParsePhoneNumber(t, "+49 AAAA BBBBBB", "49", "AAAA", "BBBBBB", "")
}

func TestParsePhoneNumberGermanyMicrosoft(t *testing.T) {
    testParsePhoneNumber(t, "+49 (AAAA) BBBBBB", "49", "AAAA", "BBBBBB", "")
}

func TestParsePhoneNumberGermanyOld(t *testing.T) {
    testParsePhoneNumber(t, "0AAAA-BBBBBB", "", "0AAAA", "BBBBBB", "")
}

func TestParsePhoneNumberUKMobile(t *testing.T) {
    testParsePhoneNumber(t, "07AAA BBBBBB", "", "07AAA", "BBBBBB", "")
}

func TestParsePhoneNumberUKStd(t *testing.T) {
    testParsePhoneNumber(t, "(02x) AAAA AAAA", "", "02x", "AAAA-AAAA", "")
    testParsePhoneNumber(t, "(01xxx) AAAAAA", "", "01xxx", "AAAAAA", "")
    testParsePhoneNumber(t, "(01AAA) BBBBB", "", "01AAA", "BBBBB", "")
    testParsePhoneNumber(t, "(01AA AA) BBBBB", "", "01AA AA", "BBBBB", "")
    testParsePhoneNumber(t, "(01AA AA) BBBB", "", "01AA AA", "BBBB", "")
    testParsePhoneNumber(t, "0AAA BBB BBBB", "", "0AAA", "BBB-BBBB", "")
}

func TestParsePhoneNumberIndiaLandline(t *testing.T) {
    testParsePhoneNumber(t, "0AAA-BBBBBBB", "", "0AAA", "BBBBBBB", "")
}

func TestParsePhoneNumberIndiaMobile(t *testing.T) {
    testParsePhoneNumber(t, "AAAAA-BBBBB", "", "AAAAA", "BBBBB", "")
    testParsePhoneNumber(t, "+91-AAAAA-BBBBB", "91", "AAAAA", "BBBBB", "")
}

func TestParsePhoneNumberChinaLandlines(t *testing.T) {
    testParsePhoneNumber(t, "(0XXX) YYYY YYYY", "", "0XXX", "YYYY-YYYY", "")
    testParsePhoneNumber(t, "+86 755 AAAA YYYY", "86", "755", "AAAA-YYYY", "")
}

func TestParsePhoneNumberChinaMobile(t *testing.T) {
    testParsePhoneNumber(t, "13A YYYY ZZZZ", "", "13A", "YYYY-ZZZZ", "")
}

func TestParsePhoneNumberHongKong(t *testing.T) {
    testParsePhoneNumber(t, "AAAA BBBB", "", "AAAA", "BBBB", "")
}
