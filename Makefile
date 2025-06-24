run:
	./build-n-serve.sh -c config.json

run-local:
	./build-n-serve.sh -c env/config.local.json

run-worker:
	./build-n-run-worker.sh -c config.json

run-worker-local:
	./build-n-run-worker.sh -c env/config.local.json
