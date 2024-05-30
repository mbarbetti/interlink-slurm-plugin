package slurm

import (
	"io"
	"net/http"

	"github.com/containerd/containerd/log"
)

func ReadRequestBody(req *http.Request, w http.ResponseWriter, h *PluginHandler) (bytes []byte, status int) {
	status = http.StatusOK
	bytes, err := io.ReadAll(req.Body)
	if err != nil {
		status = http.StatusInternalServerError
		w.WriteHeader(status)
		w.Write([]byte("Some errors occurred while creating container. Check Slurm Plugin's logs"))
		log.G(h.Ctx).Error(err)
		return
	}
	return
}
