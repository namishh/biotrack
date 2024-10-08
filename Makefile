build:
	@bunx tailwind -i ./assets/app.css -o ./public/app.css
	@templ generate
	@go build -o bin/app ./app/

css:
	@bunx tailwind -i ./assets/app.css -o ./public/app.css --watch

run: build
	@./bin/app

templ:
	@templ generate --watch

dev:
	@./bin/air


clean:
	@rm -rf bin
