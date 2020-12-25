.PHONY: all clean

OUTPUT_PUTGET=putget

#reg: clean
#	go build -o ./bin/${OUTPUT_REG} cmd/reg/main.go

all: clean
	go build -o ./bin/${OUTPUT_PUTGET} cmd/putget/main.go

clean:
#	rm -f bin/${OUTPUT_REG}
	rm -f bin/${OUTPUT_PUTGET}