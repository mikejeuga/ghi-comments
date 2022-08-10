repo=$(shell basename "`pwd`")
gopher:
	@git init
	@touch .gitignore
	@go mod init github.com/mikejeuga/$(repo)
	@go mod tidy

run:
	@go run ./cmd/main.go

t: test
test:
	@go test ./... -v


ic: init
init:
	@git add .
	@git commit -m "Initial commit"
	@git remote add origin git@github.com:mikejeuga/${repo}.git
	@git branch -M main
	@git push -u origin main

c: commit
commit:
	@git add .
	@git commit -m "$m"
	@git pull --rebase
	git push

