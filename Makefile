bin/tailwindcss: 
	./bin/install-tailwind

assets/styles.css: bin/tailwindcss assets/*.js
	./bin/tailwindcss -i ./assets/tailwind.css -o ./assets/styles.css

build-assets: assets/styles.css assets/icons.svg assets/typer.js assets/manifest.json
	./bin/cachebust assets/styles.css
	./bin/cachebust assets/icons.svg
	./bin/cachebust assets/typer.js
	./bin/cachebust assets/manifest.json


watch-css:
	./bin/tailwindcss -i ./assets/tailwind.css -o ./assets/styles.css --watch

watch-server:
	gow -e=go,html,js,css run .

