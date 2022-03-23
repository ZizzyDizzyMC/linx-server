package cleanup

import (
	"log"
	"time"

	"github.com/andreimarcu/linx-server/backends/localfs"
	"github.com/andreimarcu/linx-server/expiry"
)

func Cleanup(filesDir string, metaDir string, locksDir string, noLogs bool) {
	fileBackend := localfs.NewLocalfsBackend(metaDir, filesDir, locksDir)

	files, err := fileBackend.List()
	if err != nil {
		panic(err)
	}

	for _, filename := range files {
		locked, err := fileBackend.CheckLock(filename)
		if err != nil {
			log.Printf("Error checking if %s is locked: %s", filename, err)
		}
		if locked {
			log.Printf("%s is locked, it will be ignored", filename)
			continue
		}

		metadata, err := fileBackend.Head(filename)
		if err != nil {
			if !noLogs {
				log.Printf("Failed to find metadata for %s", filename)
			}
		}

		if expiry.IsTsExpired(metadata.Expiry) {
			if !noLogs {
				log.Printf("Delete %s", filename)
			}
			fileBackend.Delete(filename)
		}
	}
}

func PeriodicCleanup(minutes time.Duration, filesDir string, metaDir string, locksDir string, noLogs bool) {
	c := time.Tick(minutes)
	for range c {
		log.Printf("Running periodic cleanup")
		Cleanup(filesDir, metaDir, locksDir, noLogs)
		log.Printf("Finished periodic cleanup")
	}

}
