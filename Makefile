
build:
	@go build -o mivog main.go


run: build
	@./mivog
