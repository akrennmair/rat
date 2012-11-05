TARGET=rat
GO_SRC=$(wildcard *.go)

all: $(TARGET)

$(TARGET): $(GO_SRC)
	go build -o $@

clean:
	go clean

.PHONY: clean
