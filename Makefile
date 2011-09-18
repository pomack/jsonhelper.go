include $(GOROOT)/src/Make.inc

TARG=github.com/pomack/jsonhelper
GOFILES=\
    helper.go\
    jsonarray.go\
    jsonobject.go\


include $(GOROOT)/src/Make.pkg
