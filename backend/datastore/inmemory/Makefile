include $(GOROOT)/src/Make.inc

TARG=github.com/pomack/dsocial.go/backend/datastore/inmemory
GOFILES=\
	accounts_service.go\
	authentication_service.go\
	authorization_service.go\
	changesets.go\
	const.go\
	contacts_service.go\
	datastore.go\


include $(GOROOT)/src/Make.pkg

check:
	gomake test


