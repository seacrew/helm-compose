HELM_HOME := $(shell bash -c 'eval $$(helm env); echo $$HELM_PLUGINS')

.PHONY: install
install: build
	echo $(HELM_HOME)
	mkdir -p $(HELM_HOME)/plugins/helm-compose/bin
	cp bin/compose $(HELM_HOME)/plugins/helm-compose/bin
	cp plugin.yaml $(HELM_HOME)/plugins/helm-compose/

.PHONY: build
build:
	mkdir -p bin/
	go build -v -o bin/compose

.PHONY: dist
dist: export COPYFILE_DISABLE=1 #teach OSX tar to not put ._* files in tar archive
dist: export CGO_ENABLED=0
dist:
	rm -rf build/compose/* release/*
	mkdir -p build/compose/bin release/
	cp README.md LICENSE plugin.yaml build/compose
	GOOS=linux GOARCH=amd64 go build -o build/compose/bin/compose -trimpath -ldflags="$(LDFLAGS)"
	tar -C build/ -zcvf $(CURDIR)/release/helm-compose-linux-amd64.tgz compose/
	GOOS=linux GOARCH=arm64 go build -o build/compose/bin/compose -trimpath -ldflags="$(LDFLAGS)"
	tar -C build/ -zcvf $(CURDIR)/release/helm-compose-linux-arm64.tgz compose/
	GOOS=freebsd GOARCH=amd64 go build -o build/compose/bin/compose -trimpath -ldflags="$(LDFLAGS)"
	tar -C build/ -zcvf $(CURDIR)/release/helm-compose-freebsd-amd64.tgz compose/
	GOOS=darwin GOARCH=amd64 go build -o build/compose/bin/compose -trimpath -ldflags="$(LDFLAGS)"
	tar -C build/ -zcvf $(CURDIR)/release/helm-compose-macos-amd64.tgz compose/
	GOOS=darwin GOARCH=arm64 go build -o build/compose/bin/compose -trimpath -ldflags="$(LDFLAGS)"
	tar -C build/ -zcvf $(CURDIR)/release/helm-compose-macos-arm64.tgz compose/
	rm build/compose/bin/compose
	GOOS=windows GOARCH=amd64 go build -o build/compose/bin/compose.exe -trimpath -ldflags="$(LDFLAGS)"
	tar -C build/ -zcvf $(CURDIR)/release/helm-compose-windows-amd64.tgz compose/