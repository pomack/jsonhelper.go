include $(GOROOT)/src/Make.inc

TARG=github.com/pomack/jsonhelper
GOFILES=\
    encode.go\
    helper.go\
    jsonarray.go\
    jsonobject.go\
    tags.go


include $(GOROOT)/src/Make.pkg
