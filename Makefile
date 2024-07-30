build-app:
	@go build -o bin/app ./app/

css:
	@tailwind -i ./assets/app.css -o ./public/assets/app.css --watch

run: build-app
	@./bin/app

dev:
	@./bin/air


clean:
	@rm -rf bin
