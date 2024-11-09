
sources := $(wildcard */*.go)

all: cribbl.msi

cribledge.exe: $(sources)
	go build .

installer/product.wixobj: installer/product.wxs cribledge.exe
	candle -arch x64 -out installer/product.wixobj installer/product.wxs

cribbl.msi: installer/product.wixobj
	light -out cribbl.msi -spdb installer/product.wixobj

test:
	go test -count 1 -v ./...

clean:
	go clean .
	rm -f cribbl.msi
	rm -f cribledge.exe
	rm -f installer/*.wixobj


	