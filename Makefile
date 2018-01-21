# you need the following tools installed and in your path variable
# - github.com/tdewolff/minify/cmd/minify
# - lessc
# - wget

all: init minify generate-css
	go install

init: update

update: update-js update-go

update-go:
	go get -u

update-js:
	wget -O static/extensions/katex/katex.min.js https://cdnjs.cloudflare.com/ajax/libs/KaTeX/0.9.0-alpha2/katex.min.js
	wget -O static/extensions/katex/katex.min.css https://cdnjs.cloudflare.com/ajax/libs/KaTeX/0.9.0-alpha2/katex.min.css
	wget -O static/extensions/katex-autorender/katex-autorender.min.js https://cdnjs.cloudflare.com/ajax/libs/KaTeX/0.9.0-alpha2/contrib/auto-render.min.js

minify:
	minify -r --match .+\.js static/js/quest static/extensions/katex static/js static/extensions/katex-autorender -o static/js/quest.min.js

generate-css:
	lessc static/style/about.less static/style/about.css
	lessc static/style/account.less static/style/account.css
	lessc static/style/calendar.less static/style/calendar.css
	lessc static/style/day.less static/style/day.css
	lessc static/style/home.less static/style/home.css
	lessc static/style/register.less static/style/register.css
	lessc static/style/login.less static/style/login.css