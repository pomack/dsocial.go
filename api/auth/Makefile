include $(GOROOT)/src/Make.inc

TARG=github.com/pomack/dsocial.go/api/auth
GOFILES=\
	authorize.go\
	deauthorize.go\
	generate_private_key.go\
	login.go\
	logout.go\
	set_password.go\


include $(GOROOT)/src/Make.pkg

check:
	gomake test


