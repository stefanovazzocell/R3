.PHONY: build clean install

@:
	@echo "Preparing..."
	@rm -r ../www/*
	@touch ../www/placeholder
	@echo "Copying/Linking files..."
	@rsync -a * ../www/ --exclude Makefile --exclude res/style.css --exclude res/main.js --exclude index.html
	@ln index.html ../www/index.html
	@ln res/main.js ../www/res/main.js
	@ln res/style.css ../www/res/style.css
	@echo "Starting server..."
	cd ../www/; python3 -m http.server 8888 --bind 0.0.0.0
	@echo "Cleanup..."
	@rm -r ../www/*
	@touch ../www/placeholder

build:
	@echo "Preparing..."
	@rm -r ../www/*
	@touch ../www/placeholder
	@echo "Copying files..."
	@rsync -a * ../www/ --exclude Makefile --exclude res/style.css --exclude res/main.js --exclude index.html
	@echo "Compressing files [CSS]..."
	@csso res/style.css > ../www/res/style.css
	@echo "Compressing files [JS]..."
	@terser res/main.js -m -c --toplevel --mangle-props --name-cache ../www/res/main.js.map > ../www/res/main.js
	@echo "Compressing files [HTML]..."
	@html-minifier --collapse-whitespace --remove-comments --remove-optional-tags --remove-redundant-attributes --remove-script-type-attributes --remove-tag-whitespace --use-short-doctype --collapse-boolean-attributes index.html > ../www/index.html
	@echo "Starting server..."
	cd ../www/; python3 -m http.server 8888 --bind 0.0.0.0
	@echo "Cleanup..."
	@rm -r ../www/*
	@touch ../www/placeholder

clean:
	@rm -r ../www/*
	@touch ../www/placeholder

install:
	@echo "Installing required packages"
	@sudo apt install npm
	@sudo npm install terser csso-cli html-minifier -g