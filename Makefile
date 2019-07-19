


.PHONY: proto

proto:
	prototool generate ./primitive
	prototool format ./primitive -w -f
	prototool lint ./primitive
