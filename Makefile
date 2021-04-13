server: build
	go build -o build/server main.go


build:
	mkdir -p build

clean:
	rm -rf build
