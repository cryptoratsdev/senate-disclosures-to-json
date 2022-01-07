run:
	go run *.go

deps:
	sudo apt install -y libreoffice tesseract-ocr

ifdef GH_TOKEN
PUBLIC_URL = https://$(GH_TOKEN)@github.com/cryptoratsdev/senate-disclosures.git
else
PUBLIC_URL = git@github.com:cryptoratsdev/senate-disclosures.git
endif
output:
	git clone $(PUBLIC_URL) output

setup-git:
	git config --global user.email "$(GIT_EMAIL)"
	git config --global user.name "$(GIT_NAME)"

TS := $(date)
commit: output
	cd output \
	&& git add  . \
	&& (git commit -a -m "Data updated at $(shell date)"  || echo "Nothing to commit") \
	&& git push \
	&& cd ..
