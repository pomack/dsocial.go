
all: install

DIRS=\
	apiutil/\
	account/\
	auth/\
	contacts/\
	external_services/\
	groups/\

TEST=\
	$(filter-out $(NOTEST),$(DIRS))


clean.dirs: $(addsuffix .clean, $(DIRS))
install.dirs: $(addsuffix .install, $(DIRS))
nuke.dirs: $(addsuffix .nuke, $(DIRS))
test.dirs: $(addsuffix .test, $(TEST))

%.clean:
	+cd $* && make clean

%.install:
	+cd $* && make install

%.nuke:
	+cd $* && make nuke

%.test:
	+cd $* && make test

%.check:
	+cd $* && make check

clean: clean.dirs

install: install.dirs

test:   test.dirs

check:	check.dirs

nuke: nuke.dirs

echo-dirs:
	@echo $(DIRS)
