default:
  go run .

test:
  go test ./...

clean:
  rm -rf build/

build:
  @echo "Creating build directory..."
  mkdir -p build

  @echo "Building reddittui application..."
  go build -o build/reddittui main.go

install: build
  ./install.sh

uninstall: clean
  @echo "Removing binary from /usr/local/bin/reddittui..."
  sudo rm -f /usr/local/bin/reddittui
