@default: run

@run:
  go run .

test:
  go test -v ./...

clean:
  rm -rf build/

build:
  @echo "Building reddittui..."

  @echo "Creating build directory at build/..."
  mkdir -p build

  @echo "Installing dependencies..."
  go mod tidy

  @echo "Building reddittui application..."
  go build -o build/reddittui main.go

  @echo "Build complete."

install: build
  @echo "Installing reddittui..."
  ./install.sh
  @echo "Installation complete."

uninstall: clean
  @echo "Cleaning reddittui..."
  sudo rm -f /usr/local/bin/reddittui
  @echo "Clean complete"

setupIntegTests:
  mkdir -p ~/.cache/reddittui
  cp -r testData/* ~/.cache/reddittui
