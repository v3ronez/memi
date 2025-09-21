run: build
	@./bin/memi

build:
	npx tailwindcss -i view/css/app.css -o public/styles.css
	@go tool templ generate view
	@go build -o bin/memi main.go

live/air:
	@go tool air

live/templ:
	@go tool templ generate --watch --proxy="http://localhost:8080" --open-browser=false -v

live/tailwind:
	npx tailwindcss -i view/css/app.css -o public/styles.css --minify --watch

live:
	make -j2 live/air live/templ

