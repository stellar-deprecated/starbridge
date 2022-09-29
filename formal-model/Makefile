MODEL_DIR=/home/user/model/
OUTPUT_DIR=/home/user/output/
N_THREADS ?= 2

.PHONY: build dependency-graph check-proof clean

all: build

clean:
	rm -r output

build:
	docker build -t starbridge-ivy docker/

output:
	mkdir -p output

# NOTE: this can sometimes fail (it seems there is a race condition somewhere...); just retry in that case.
dependency-graph: | output starbridge.ivy
	docker run --init --rm -t -v "$(CURDIR):$(MODEL_DIR)" -v "$(CURDIR)/output:$(OUTPUT_DIR)" starbridge-ivy plot_dependencies.sh $(MODEL_DIR)/starbridge.ivy  $(OUTPUT_DIR) $(N_THREADS)

check-proof:
	docker run --init --rm -t -v $(CURDIR):$(MODEL_DIR) starbridge-ivy check_invariants.sh $(MODEL_DIR)/starbridge.ivy
