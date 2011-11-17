package inmemory

import (
    "os"
)

const (
    
    _INMEMORY_USER_ACCOUNT_COLLECTION_NAME = "user_accounts"
    _INMEMORY_USER_ACCOUNT_ID_FOR_USERNAME_COLLECTION_NAME = "user_account_id_for_username"
    _INMEMORY_USER_ACCOUNT_ID_FOR_EMAIL_COLLECTION_NAME = "user_account_id_for_email"
    _INMEMORY_CONSUMER_ACCOUNT_COLLECTION_NAME = "consumer_accounts"
    _INMEMORY_EXTERNAL_USER_ACCOUNT_COLLECTION_NAME = "external_user_accounts"
    
    _INMEMORY_CONTACT_COLLECTION_NAME = "contacts"
    _INMEMORY_CONNECTION_COLLECTION_NAME = "connections"
    _INMEMORY_GROUP_COLLECTION_NAME = "group"
    _INMEMORY_CHANGESET_COLLECTION_NAME = "changesets"
    _INMEMORY_EXTERNAL_CONTACT_COLLECTION_NAME = "external_contacts"
    _INMEMORY_EXTERNAL_GROUP_COLLECTION_NAME = "external_group"
    _INMEMORY_CONTACT_FOR_EXTERNAL_CONTACT_COLLECTION_NAME = "contacts_for_external_contacts"
    _INMEMORY_GROUP_FOR_EXTERNAL_GROUP_COLLECTION_NAME = "groups_for_external_group"
    _INMEMORY_EXTERNAL_TO_INTERNAL_CONTACT_MAPPING_COLLECTION_NAME = "external_to_internal_contact_mappings"
    _INMEMORY_INTERNAL_TO_EXTERNAL_CONTACT_MAPPING_COLLECTION_NAME = "internal_to_external_contact_mappings"
    _INMEMORY_EXTERNAL_TO_INTERNAL_GROUP_MAPPING_COLLECTION_NAME = "external_to_internal_group_mappings"
    _INMEMORY_INTERNAL_TO_EXTERNAL_GROUP_MAPPING_COLLECTION_NAME = "internal_to_external_group_mappings"
    _INMEMORY_USER_TO_CONTACT_SETTINGS_COLLECTION_NAME = "user_to_contact_settings"
    _INMEMORY_CHANGESETS_TO_APPLY_COLLECTION_NAME = "changesets_to_apply"
    _INMEMORY_CHANGESETS_NOT_CURRENTLY_APPLYABLE_COLLECTION_NAME = "changesets_not_currently_applyable"
)

var (
    ERR_ACCOUNT_ALREADY_EXISTS_WITH_SPECIFIED_ID os.Error
    ERR_ACCOUNT_ALREADY_EXISTS_WITH_SPECIFIED_USERNAME os.Error
    ERR_ACCOUNT_ALREADY_EXISTS_WITH_SPECIFIED_SHORTNAME os.Error
    ERR_ACCOUNT_ALREADY_EXISTS_WITH_SPECIFIED_EMAIL os.Error
    ERR_ACCOUNT_MUST_SPECIFY_SHORTNAME os.Error
)

func init() {
    ERR_ACCOUNT_ALREADY_EXISTS_WITH_SPECIFIED_ID = os.NewError("Account already exists with specified id")
    ERR_ACCOUNT_ALREADY_EXISTS_WITH_SPECIFIED_USERNAME = os.NewError("Account already exists with specified username")
    ERR_ACCOUNT_ALREADY_EXISTS_WITH_SPECIFIED_SHORTNAME = os.NewError("Account already exists with specified short-name")
    ERR_ACCOUNT_ALREADY_EXISTS_WITH_SPECIFIED_EMAIL = os.NewError("Account already exists with specified email")
    ERR_ACCOUNT_MUST_SPECIFY_SHORTNAME = os.NewError("Must specify short-name")
}


