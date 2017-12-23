all: minify
	go install

init: update

update: update-js update-go all

update-go:
	go get -u

update-js:
	wget -O static/extensions/katex/katex.min.js https://cdnjs.cloudflare.com/ajax/libs/KaTeX/0.9.0-alpha2/katex.min.js
	wget -O static/extensions/katex/katex.min.css https://cdnjs.cloudflare.com/ajax/libs/KaTeX/0.9.0-alpha2/katex.min.css
	wget -O static/extensions/katex-autorender/katex-autorender.min.js https://cdnjs.cloudflare.com/ajax/libs/KaTeX/0.9.0-alpha2/contrib/auto-render.min.js

minify:
	minify -r --match .+\.js static/js/quest static/extensions/katex static/js static/extensions/katex-autorender -o static/js/quest.min.js