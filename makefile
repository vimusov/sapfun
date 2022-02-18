BINDIR ?= /usr/bin

GOPATH := $(PWD)

TARGET := sapfun

$(TARGET):
	GOPATH=$(GOPATH) go build -o $@ $@.go

all: $(TARGET)

install:
	install -D --mode=0755 $(TARGET) $(DESTDIR)$(BINDIR)/$(TARGET)
