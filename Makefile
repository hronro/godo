.PHONY: build run templ tailwindcss watch watch-run watch-templ watch-tailwind clean clean-bin clean-templ clean-tailwind

build: templ tailwind
	go build .

run: templ tailwind
	go run .

templ:
	go tool templ generate

tailwind:
	tailwindcss -i app.css -o static/style.css --minify

watch:
	make -j8 watch-run watch-templ watch-tailwind

watch-run:
	go tool wgo run -file static/style.css .

watch-templ:
	go tool wgo -file .templ make templ

watch-tailwind:
	tailwindcss -i app.css -o static/style.css --watch --minify

clean: clean-bin clean-templ clean-tailwind

clean-bin:
	rm -f ./godo ./godo.exe

clean-templ:
	rm -f ./templates/*_templ.go
	rm -f ./layout/*_templ.go

clean-tailwind:
	rm -f static/style.css
