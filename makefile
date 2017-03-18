
.PHONY: fmt

ifeq ($(OS),Windows_NT)
  TARGET = wandbox.exe
else
  TARGET = wandbox
endif

SOURCES=src/main.go

$(TARGET): $(SOURCES)
	go build -o $(TARGET) $(SOURCES)

fmt: $(SOURCES)
	gofmt -w $(SOURCES)
