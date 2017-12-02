#
.PHONY:	dummy


TARGET=gw




all:	$(TARGET)




$(TARGET):	main.go
	go build
	strip $@




tidy:	dummy
	$(RM) *~


clean:	dummy
	$(RM) $(TARGET)
