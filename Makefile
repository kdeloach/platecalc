.PHONY: calc-wasm
calc-wasm:
	cd ./cmd/calc-wasm/ && \
		GOOS=js GOARCH=wasm go build  -o ../../bin/platecalc.wasm

clean:
	rm -f ./bin/platecalc.wasm
