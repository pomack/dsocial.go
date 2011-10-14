package dsocial

const (
    
    // ChangeType
    CHANGE_TYPE_CREATE  ChangeType = "create"
    CHANGE_TYPE_ADD     ChangeType = "add"
    CHANGE_TYPE_UPDATE  ChangeType = "update"
    CHANGE_TYPE_DELETE  ChangeType = "delete"
    
    // Path Component Type
    PATH_COMPONENT_TYPE_ID = 1
    PATH_COMPONENT_TYPE_KEY = 2
    PATH_COMPONENT_TYPE_INDEX = 3
    PATH_COMPONENT_TYPE_MAP_INDEX = 4
    
    
    // Contact Name Ordering
    GIVEN_MIDDLE_SURNAME          ContactNameOrdering = "gms"
    SURNAME_GIVEN_MIDDLE          ContactNameOrdering = "sgm"
    GIVEN_SURNAME_MIDDLE          ContactNameOrdering = "gsm"
    SURNAME_MIDDLE_GIVEN          ContactNameOrdering = "smg"
    MIDDLE_GIVEN_SURNAME          ContactNameOrdering = "mgs"
    MIDDLE_SURNAME_GIVEN          ContactNameOrdering = "msg"
    DEFAULT_CONTACT_NAME_ORDERING ContactNameOrdering = GIVEN_MIDDLE_SURNAME

    // Common rels
    REL_HOME  = "home"
    REL_WORK  = "work"
    REL_OTHER = "other"

    // Postal Address Types
    REL_ADDRESS_HOME  RelPostalAddress = REL_HOME
    REL_ADDRESS_WORK  RelPostalAddress = REL_WORK
    REL_ADDRESS_OTHER RelPostalAddress = REL_OTHER
    REL_POBOX         RelPostalAddress = "po_box"

    // Education types
    REL_EDUCATION_ELEMENTARY_SCHOOL RelEducation = "elementary_school"
    REL_EDUCATION_MIDDLE_SCHOOL     RelEducation = "middle_school"
    REL_EDUCATION_HIGH_SCHOOL       RelEducation = "high_school"
    REL_EDUCATION_COLLEGE           RelEducation = "college"
    REL_EDUCATION_GRADUATE_SCHOOL   RelEducation = "graduate_school"
    REL_EDUCATION_VOCATIONAL        RelEducation = "vocational"
    REL_EDUCATION_OTHER   RelEducation = REL_OTHER

    // Phone rels
    REL_PHONE_HOME         RelPhoneNumber = REL_HOME
    REL_PHONE_WORK         RelPhoneNumber = REL_WORK
    REL_PHONE_OTHER        RelPhoneNumber = REL_OTHER
    REL_PHONE_ASSISTANT    RelPhoneNumber = "assistant"
    REL_PHONE_CALLBACK     RelPhoneNumber = "callback"
    REL_PHONE_CAR          RelPhoneNumber = "car"
    REL_PHONE_COMPANY_MAIN RelPhoneNumber = "company_main"
    REL_PHONE_EXTERNAL     RelPhoneNumber = "external"
    REL_PHONE_FAX          RelPhoneNumber = "fax"
    REL_PHONE_GOOGLE_VOICE RelPhoneNumber = "google_voice"
    REL_PHONE_HOME_FAX     RelPhoneNumber = "home_fax"
    REL_PHONE_ISDN         RelPhoneNumber = "isdn"
    REL_PHONE_MAIN         RelPhoneNumber = "main"
    REL_PHONE_MOBILE       RelPhoneNumber = "mobile"
    REL_PHONE_OTHER_FAX    RelPhoneNumber = "other_fax"
    REL_PHONE_PAGER        RelPhoneNumber = "pager"
    REL_PHONE_RADIO        RelPhoneNumber = "radio"
    REL_PHONE_SIP          RelPhoneNumber = "sip"
    REL_PHONE_SKYPE        RelPhoneNumber = "skype"
    REL_PHONE_TELEX        RelPhoneNumber = "telex"
    REL_PHONE_TTY_TDD      RelPhoneNumber = "tty_tdd"
    REL_PHONE_WORK_FAX     RelPhoneNumber = "work_fax"
    REL_PHONE_WORK_MOBILE  RelPhoneNumber = "work_mobile"
    REL_PHONE_WORK_PAGER   RelPhoneNumber = "work_pager"

    // Email rels
    REL_EMAIL_HOME  RelEmail = REL_HOME
    REL_EMAIL_WORK  RelEmail = REL_WORK
    REL_EMAIL_OTHER RelEmail = REL_OTHER

    // urls
    REL_URI_BLOG           RelUri = "blog"
    REL_URI_DSOCIAL        RelUri = "dsocial"
    REL_URI_FACEBOOK       RelUri = "facebook"
    REL_URI_FTP            RelUri = "ftp"
    REL_URI_GOOGLE_PROFILE RelUri = "google_profile"
    REL_URI_GOOGLE_PLUS    RelUri = "google_plus"
    REL_URI_HCARD          RelUri = "hcard"
    REL_URI_HOMEPAGE       RelUri = "homepage"
    REL_URI_LINKEDIN       RelUri = "linkedin"
    REL_URI_SMUGMUG        RelUri = "smugmug"
    REL_URI_TWITTER        RelUri = "twitter"
    REL_URI_VCARD          RelUri = "vcard"
    REL_URI_WEBSITE        RelUri = "website"
    REL_URI_YAHOO          RelUri = "yahoo"
    REL_URI_HOME           RelUri = REL_HOME
    REL_URI_WORK           RelUri = REL_WORK
    REL_URI_OTHER          RelUri = REL_OTHER

    // IM rels
    REL_IM_HOME       RelIM = REL_HOME
    REL_IM_WORK       RelIM = REL_WORK
    REL_IM_OTHER      RelIM = REL_OTHER
    REL_IM_NETMEETING RelIM = "netmeeting"

    // IM Protocol rels
    REL_IM_PROT_AIM             RelIMProtocol = "aim"
    REL_IM_PROT_BONJOUR         RelIMProtocol = "bonjour"
    REL_IM_PROT_DOTMAC          RelIMProtocol = "dotmac"
    REL_IM_PROT_FACEBOOK        RelIMProtocol = "facebook"
    REL_IM_PROT_GADU_GADU       RelIMProtocol = "gadu_gadu"
    REL_IM_PROT_GOOGLE_TALK     RelIMProtocol = "google_talk"
    REL_IM_PROT_GROUPWISE       RelIMProtocol = "groupwise"
    REL_IM_PROT_ICQ             RelIMProtocol = "icq"
    REL_IM_PROT_IRC             RelIMProtocol = "irc"
    REL_IM_PROT_JABBER          RelIMProtocol = "jabber"
    REL_IM_PROT_LIVEJOURNAL     RelIMProtocol = "livejournal"
    REL_IM_PROT_MOBILE_ME       RelIMProtocol = "mobile_me"
    REL_IM_PROT_MSN             RelIMProtocol = "msn"
    REL_IM_PROT_MYSPACE_IM      RelIMProtocol = "myspaceim"
    REL_IM_PROT_QQ              RelIMProtocol = "qq"
    REL_IM_PROT_SAMETIME        RelIMProtocol = "sametime"
    REL_IM_PROT_SIP             RelIMProtocol = "sip"
    REL_IM_PROT_SKYPE           RelIMProtocol = "skype"
    REL_IM_PROT_STATUSNET       RelIMProtocol = "statusnet"
    REL_IM_PROT_TWITTER         RelIMProtocol = "twitter"
    REL_IM_PROT_YAHOO_MESSENGER RelIMProtocol = "yahoo_messenger"
    REL_IM_PROT_OTHER           RelIMProtocol = REL_OTHER

    // relationship rels
    REL_RELATIONSHIP_ASSISTANT        RelRelationship = "assistant"
    REL_RELATIONSHIP_AUNT             RelRelationship = "aunt"
    REL_RELATIONSHIP_BROTHER          RelRelationship = "brother"
    REL_RELATIONSHIP_CHILD            RelRelationship = "child"
    REL_RELATIONSHIP_COUSIN           RelRelationship = "cousin"
    REL_RELATIONSHIP_DOMESTIC_PARTNER RelRelationship = "domestic_partner"
    REL_RELATIONSHIP_FATHER           RelRelationship = "father"
    REL_RELATIONSHIP_FRIEND           RelRelationship = "friend"
    REL_RELATIONSHIP_HUSBAND          RelRelationship = "husband"
    REL_RELATIONSHIP_MANAGER          RelRelationship = "manager"
    REL_RELATIONSHIP_MOTHER           RelRelationship = "mother"
    REL_RELATIONSHIP_PARENT           RelRelationship = "parent"
    REL_RELATIONSHIP_PARTNER          RelRelationship = "partner"
    REL_RELATIONSHIP_REFERRED_BY      RelRelationship = "referred_by"
    REL_RELATIONSHIP_RELATIVE         RelRelationship = "relative"
    REL_RELATIONSHIP_SIBLING          RelRelationship = "sibling"
    REL_RELATIONSHIP_SISTER           RelRelationship = "sister"
    REL_RELATIONSHIP_SPOUSE           RelRelationship = "spouse"
    REL_RELATIONSHIP_UNCLE            RelRelationship = "uncle"
    REL_RELATIONSHIP_WIFE             RelRelationship = "wife"
    REL_RELATIONSHIP_OTHER            RelRelationship = "other"

    // Date rels
    REL_DATE_OTHER       RelDate = REL_OTHER
    REL_DATE_ANNIVERSARY RelDate = "anniversary"

    // Datetime rels
    REL_DATETIME_OTHER     RelDateTime = REL_OTHER
    REL_DATETIME_BIRTHTIME RelDateTime = "birthtime"

    // Gender types
    REL_GENDER_MALE   RelGender = "male"
    REL_GENDER_FEMALE RelGender = "female"
    REL_GENDER_OTHER  RelGender = REL_OTHER

    // relationship status
    REL_SINGLE                    RelRelationshipStatus = "single"
    REL_DIVORCED                  RelRelationshipStatus = "divorced"
    REL_IN_A_RELATIONSHIP         RelRelationshipStatus = "in_a_relationship"
    REL_ENGAGED                   RelRelationshipStatus = "engaged"
    REL_SEPARATED                 RelRelationshipStatus = "separated"
    REL_MARRIED                   RelRelationshipStatus = "married"
    REL_ITS_COMPLICATED           RelRelationshipStatus = "it's_complicated"
    REL_OPEN_RELATIONSHIP         RelRelationshipStatus = "open_relationship"
    REL_WIDOWED                   RelRelationshipStatus = "widowed"
    REL_IN_DOMESTIC_PARTNERSHIP   RelRelationshipStatus = "in_domestic_partnership"
    REL_IN_CIVIL_UNION            RelRelationshipStatus = "in_civil_union"
    REL_RELATIONSHIP_STATUS_OTHER RelRelationshipStatus = REL_OTHER
)


