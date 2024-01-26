build-css:
	tailwindcss -i ./assets/tailwind.css -o ./assets/styles.css

build-server:
	go run .

build: build-css build-server

watch-css:
	tailwindcss -i ./assets/tailwind.css -o ./assets/styles.css --watch

watch-server:
	gow -e=go,html run .