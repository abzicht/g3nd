
run:
	cd ../engine/renderer/shaders && go run github.com/g3n/engine/tools/g3nshaders
	go run ./main.go physics-experimental.particles

.PHONY=run
