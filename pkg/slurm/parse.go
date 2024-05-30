package slurm

import (
	"encoding/json"
	"fmt" // debug
	"net/http"

	"github.com/containerd/containerd/log"
	commonIL "github.com/intertwin-eu/interlink/pkg/interlink"
	v1 "k8s.io/api/core/v1"
)

func ParseSubmitJson(bytes []byte, w http.ResponseWriter, h *PluginHandler) (podData commonIL.RetrievedPodData, statusCode int) {
	err := json.Unmarshal(bytes, &podData)
	if err != nil {
		statusCode = http.StatusInternalServerError
		w.WriteHeader(statusCode)
		w.Write([]byte("Some errors occurred while creating container. Check Slurm Plugin's logs"))
		log.G(h.Ctx).Error(err)
		fmt.Printf("%s", bytes) // debug
		return
	}
	return
}

func ParseDeleteJson(bytes []byte, w http.ResponseWriter, h *PluginHandler) (pod *v1.Pod, statusCode int) {
	err := json.Unmarshal(bytes, &pod)
	if err != nil {
		statusCode = http.StatusInternalServerError
		w.WriteHeader(statusCode)
		w.Write([]byte("Some errors occurred while deleting container. Check Slurm Plugin's logs"))
		log.G(h.Ctx).Error(err)
		return
	}
	return
}

func ParseStatusJson(bytes []byte, w http.ResponseWriter, h *PluginHandler) (podSlice []*v1.Pod, statusCode int) {
	err := json.Unmarshal(bytes, &podSlice)
	if err != nil {
		statusCode = http.StatusInternalServerError
		w.WriteHeader(statusCode)
		w.Write([]byte("Some errors occurred while retrieving container statusCode. Check Slurm Plugin's logs"))
		log.G(h.Ctx).Error(err)
		return
	}
	return
}

func ParseLogsJson(bytes []byte, w http.ResponseWriter, h *PluginHandler) (logStruct commonIL.LogStruct, statusCode int) {
	err := json.Unmarshal(bytes, &logStruct)
	if err != nil {
		statusCode = http.StatusInternalServerError
		w.WriteHeader(statusCode)
		w.Write([]byte("Some errors occurred while unmarshalling log request. Check Slurm Plugin's logs"))
		log.G(h.Ctx).Error(err)
		return
	}
	return
}
