.PHONY: all clean

OUTPUT=registry

all: clean
	go build -o ./bin/${OUTPUT} main.go

clean:
	rm -f bin/${OUTPUT}