package slurm

import (
	"net/http"
	"os"

	"github.com/containerd/containerd/log"
)

// StopHandler runs a scancel command, updating JIDs and cached statuses
func (h *PluginHandler) StopHandler(w http.ResponseWriter, r *http.Request) {
	log.G(h.Ctx).Info("Slurm Sidecar: received Stop call")

	bodyBytes, _ := ReadRequestBody(r, w, h)
	pod, statusCode := ParseDeleteJson(bodyBytes, w, h)

	filesPath := h.Config.DataRootFolder + pod.Namespace + "-" + string(pod.UID)

	err := deleteContainer(h.Ctx, h.Config, string(pod.UID), h.JIDs, filesPath+"/"+pod.Namespace)
	if err != nil {
		statusCode := http.StatusInternalServerError
		w.WriteHeader(statusCode)
		w.Write([]byte("Error deleting containers. Check Slurm Plugin's logs"))
		log.G(h.Ctx).Error(err)
		return
	}
	if os.Getenv("SHARED_FS") != "true" {
		err = os.RemoveAll(filesPath)
		if err != nil {
			statusCode = http.StatusInternalServerError
			w.WriteHeader(statusCode)
			w.Write([]byte("Error deleting containers. Check Slurm Plugin's logs"))
			log.G(h.Ctx).Error(err)
			return
		}
	}

	w.WriteHeader(statusCode)
	if statusCode != http.StatusOK {
		w.Write([]byte("Some errors occurred deleting containers. Check Slurm Plugin's logs"))
	} else {

		w.Write([]byte("All containers for submitted Pods have been deleted"))
	}
}
