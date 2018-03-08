# you need the following tools installed and in your path variable
# - github.com/tdewolff/minify/cmd/minify
# - lessc
# - wget
# - gzip
# - zip
# - tar

all: init minify generate-css
	go install

init: update

update: update-js update-go update-fonts

update-go:
	go get -u

update-js:
	wget -O - --header="Accept-Encoding: gzip" https://cdnjs.cloudflare.com/ajax/libs/KaTeX/0.9.0-alpha2/katex.min.js | gunzip > static/extensions/katex/katex.min.js
	wget -O - --header="Accept-Encoding: gzip" https://cdnjs.cloudflare.com/ajax/libs/KaTeX/0.9.0-alpha2/katex.min.css | gunzip > static/extensions/katex/katex-fonts/katex.min.css
	wget -O - --header="Accept-Encoding: gzip" https://cdnjs.cloudflare.com/ajax/libs/KaTeX/0.9.0-alpha2/contrib/auto-render.min.js | gunzip > static/extensions/katex-autorender/katex-autorender.min.js

update-fonts:
	wget -O static/fonts/zillaslab.ttf https://github.com/google/fonts/raw/master/ofl/zillaslab/ZillaSlab-Regular.ttf # gzipped version not available
	wget -O - --header="Accept-Encoding: gzip" https://raw.githubusercontent.com/jung-kurt/gofpdf/master/font/cp1252.map | gunzip > static/fonts/cp1252.map
minify:
	minify -r --match .+\.js static/js/quest static/extensions/katex static/extensions/katex-autorender -o static/js/quest.min.js

generate-css:
	rm -f static/style/*.css
	lessc static/style/about.less static/style/about.css
	lessc static/style/account.less static/style/account.css
	lessc static/style/calendar.less static/style/calendar.css
	lessc static/style/day.less static/style/day.css
	lessc static/style/home.less static/style/home.css
	lessc static/style/register.less static/style/register.css
	lessc static/style/login.less static/style/login.css
	lessc static/style/tos.less static/style/tos.css
	lessc static/style/certificate.less static/style/certificate.css


release: config = config.ini
release: git_commit=$(shell git rev-list -1 HEAD)
release: version = development

release:
ifeq ($(config),)
	@echo "Assuming you want to use config.ini ! You can specify a different file using the following syntax: make release config=<your_config.ini>"
endif
ifeq ($(version),)
	@echo "Assuming this is a development build ! You can specify a version using the following syntax: make release version=<version_string>"
endif
ifeq ($(os),)
	@echo "Assuming you want to compile for your GOOS ! You can specify a different os using the following syntax: make release os=<target_os>"
	@echo "You must export your GOOS before running this command !"
endif
ifeq ($(arc),)
	@echo "Assuming you want to compile for your GOARCH ! You can specify a different architecture using the following syntax: make release arc=<target_arc>"
	@echo "You must export your GOARCH before running this command !"
endif
	@read -r -p "Also, make sure your repository is complete and up to date ! Continue ? [Y/n]: " continue; \
	([ "$(continue)" != "n" ] && [ "$(continue)" != "N" ] && [ "$(continue)" != "no" ]) || (echo "Exiting."; exit 1;);

ifeq ($(os),)
ifeq ($(arc),)
	@go build -ldflags "-X main.GitCommit=$(git_commit) -X main.Version=$(version) -X main.DefaultConfigPath=$(config)"
ifeq ($(GOOS), linux)
	@rm -f "infoimadvent($(GOOS)|$(GOARCH)).tar.gz"
	@tar cfz "infoimadvent($(GOOS)|$(GOARCH)).tar.gz" infoimadvent LICENSE README.md $(config) templates locales data/sessions/keep.me data/documents/keep.me data/assets/keep.me static/style/*.css static/fonts static/js/*.min.js static/favicon.ico static/style/img/* static/extensions/katex/katex-fonts/fonts static/extensions/katex/katex-fonts/katex.min.css
else
	@rm -f "infoimadvent($(GOOS)|$(GOARCH)).zip"
	@zip -qr "infoimadvent($(GOOS)|$(GOARCH))" infoimadvent LICENSE README.md $(config) templates locales data/sessions/keep.me data/documents/keep.me data/assets/keep.me static/style/*.css static/fonts static/js/*.min.js static/favicon.ico static/style/img/* static/extensions/katex/katex-fonts/fonts static/extensions/katex/katex-fonts/katex.min.css
endif
else
	@GOARCH=$(arc) go build -ldflags "-X main.GitCommit=$(git_commit) -X main.Version=$(version) -X main.DefaultConfigPath=$(config)"
ifeq ($(GOOS), linux)
	@rm -f "infoimadvent($(GOOS)|$(arc)).tar.gz"
	@tar cfz "infoimadvent($(GOOS)|$(arc)).tar.gz" infoimadvent LICENSE README.md $(config) templates locales data/sessions/keep.me data/documents/keep.me data/assets/keep.me static/style/*.css static/fonts static/js/*.min.js static/favicon.ico static/style/img/* static/extensions/katex/katex-fonts/fonts static/extensions/katex/katex-fonts/katex.min.css
else
	@rm -f "infoimadvent($(GOOS)|$(arc)).zip"
	@zip -qr "infoimadvent($(GOOS)|$(arc))" infoimadvent LICENSE README.md $(config) templates locales data/sessions/keep.me data/documents/keep.me data/assets/keep.me static/style/*.css static/fonts static/js/*.min.js static/favicon.ico static/style/img/* static/extensions/katex/katex-fonts/fonts static/extensions/katex/katex-fonts/katex.min.css
endif
endif
else
ifeq ($(arc),)
	@GOOS=$(os) go build -ldflags "-X main.GitCommit=$(git_commit) -X main.Version=$(version) -X main.DefaultConfigPath=$(config)"
ifeq ($(os), linux)
	@rm -f "infoimadvent($(os)|$(GOARCH)).tar.gz"
	@tar cfz "infoimadvent($(os)|$(GOARCH)).tar.gz" infoimadvent LICENSE README.md $(config) templates locales data/sessions/keep.me data/documents/keep.me data/assets/keep.me static/style/*.css static/fonts static/js/*.min.js static/favicon.ico static/style/img/* static/extensions/katex/katex-fonts/fonts static/extensions/katex/katex-fonts/katex.min.css
else
	@rm -f "infoimadvent($(os)|$(GOARCH)).zip"
	@zip -qr "infoimadvent($(os)|$(GOARCH))" infoimadvent LICENSE README.md $(config) templates locales data/sessions/keep.me data/documents/keep.me data/assets/keep.me static/style/*.css static/fonts static/js/*.min.js static/favicon.ico static/style/img/* static/extensions/katex/katex-fonts/fonts static/extensions/katex/katex-fonts/katex.min.css
endif
else
	@GOOS=$(os) GOARCH=$(arc) go build -ldflags "-X main.GitCommit=$(git_commit) -X main.Version=$(version) -X main.DefaultConfigPath=$(config)"
ifeq ($(os), linux)
	@rm -f "infoimadvent($(os)|$(arc)).tar.gz"
	@tar cfz "infoimadvent($(os)|$(arc)).tar.gz" infoimadvent LICENSE README.md $(config) templates locales data/sessions/keep.me data/documents/keep.me data/assets/keep.me
else
	@rm -f "infoimadvent($(os)|$(arc)).zip"
	@zip -qr "infoimadvent($(os)|$(arc))" infoimadvent LICENSE README.md $(config) templates locales data/sessions/keep.me data/documents/keep.me data/assets/keep.me static/style/*.css static/fonts static/js/*.min.js static/favicon.ico static/style/img/* static/extensions/katex/katex-fonts/fonts static/extensions/katex/katex-fonts/katex.min.css
endif
endif
endif

