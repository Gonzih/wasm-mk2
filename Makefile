SUBDIRS := ./ast ./parser ./component ./walker ./registry ./core ./scope ./tree ./event
autotest:
	find . -iname '*.go' | entr -r bash -c "echo && echo && echo && go test -v --cover $(SUBDIRS)"
