run: build
	@./bin/memi

build:
	# tailwindcss -i view/css/app.css -o public/style.css
	@go tool templ generate view
	@go build -o bin/memi main.go


