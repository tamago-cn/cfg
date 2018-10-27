all:
	cd web-cfg; npm run build
	rm -rf ../static/ ../templates
	cp -rf web-cfg/dist/static ../
	mkdir -p ../templates
	cp -rf web-cfg/dist/index.html ../templates/index.tmpl
	cd ..; go generate
	go install

clean:
	rm -rf ../static/ ../templates
