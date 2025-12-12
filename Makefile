#
# Build tasks
#
gen:
	go tool templ generate

build:
	go build -o .build/app
	
css:
	./node_modules/.bin/tailwindcss -i ./assets/site.css -o ./public/site.css

clean:
	rm -rf .build
	rm -rf .tmp
	rm -rf node_modules
	rm -f **/*_templ.go
	rm -f **/*_templ.txt


migrate:
	sqlc generate
	goose sqlite3 ./db.sqlite -dir sql/migrations up 


#
# Development tasks
#

live/templ:
	go tool templ generate --watch --proxy="http://localhost:3333" --open-browser=false

live/server:
	go tool air
	
live/sync_assets:
	go tool air \
		--build.cmd "templ generate --notify-proxy" \
		--build.bin "true" \
		--build.exclude_dir "" \
		--build.include_dir "public" \
		--build.include_ext "js,css"

live/css:
	./node_modules/.bin/tailwindcss -i ./assets/site.css -o ./public/site.css --watch

dev: 
	make -j5 live/templ live/server live/sync_assets #live/css
