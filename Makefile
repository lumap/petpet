.PHONY: clean build run

build:
	@echo "Building..."
	@go build -o bin/petpet main.go
	@echo "Finished script. Check bin directory."
clean:
	@echo "Cleaning up..."
	@rm bin/petpet
	@echo "Done!"
run: build
	@bin/petpet