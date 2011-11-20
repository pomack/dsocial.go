package contacts

var (
    contactServicesMap            map[string]ContactsService
    selfUpdateableContactServices []ContactsService
    updateableContactServices     []ContactsService
)

func init() {
    contactServicesMap = make(map[string]ContactsService)
    for _, service := range []ContactsService{
        NewFacebookContactService(),
        NewGoogleContactService(),
        NewGooglePlusContactService(),
        NewLinkedInContactService(),
        NewSmugMugContactService(),
        NewTwitterContactService(),
        NewYahooContactService(),
    } {
        contactServicesMap[service.ServiceId()] = service
    }
    l := len(contactServicesMap)
    selfUpdateableContactServices = make([]ContactsService, l, l)
    updateableContactServices = make([]ContactsService, l, l)
    i, j := 0, 0
    for _, service := range contactServicesMap {
        if service.CanUpdateGroup(true) || service.CanUpdateContact(true) {
            selfUpdateableContactServices[i] = service
            i++
        }
        if service.CanUpdateGroup(false) || service.CanUpdateContact(false) {
            updateableContactServices[j] = service
            j++
        }
    }
    selfUpdateableContactServices = selfUpdateableContactServices[0:i]
    updateableContactServices = updateableContactServices[0:j]
}

func ContactServiceForName(name string) ContactsService {
    cs, _ := contactServicesMap[name]
    return cs
}

func SelfUpdateableContactServices() []ContactsService {
    return selfUpdateableContactServices
}

func UpdateableContactServices() []ContactsService {
    return updateableContactServices
}
