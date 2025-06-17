package conn

import (
	"github.com/vivasoft-ltd/go-ems/config"
	"github.com/vivasoft-ltd/go-ems/worker"
)

var workerPool *worker.Pool

func ConnectWorker() {
	workerPool = worker.NewPool(config.App().NumberOfWorkers, 2*config.App().NumberOfWorkers)
	workerPool.Start()
}

func WorkerPool() *worker.Pool {
	return workerPool
}
