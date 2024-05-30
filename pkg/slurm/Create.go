package slurm

import (
	"encoding/json"
	"net/http"
	"os"
	"strings"

	"github.com/containerd/containerd/log"
)

// SubmitHandler generates and submits a SLURM batch script according to provided data.
// 1 Pod = 1 Job. If a Pod has multiple containers, every container is a line with it's parameters in the SLURM script.
func (h *PluginHandler) SubmitHandler(w http.ResponseWriter, r *http.Request) {
	log.G(h.Ctx).Info("Slurm Plugin: received Submit call")

	bodyBytes, _ := ReadRequestBody(r, w, h)
	data, statusCode := ParseSubmitJson(bodyBytes, w, h)

	//data := data_chunk[0]
	containers := data.Pod.Spec.Containers
	metadata := data.Pod.ObjectMeta
	filesPath := h.Config.DataRootFolder + data.Pod.Namespace + "-" + string(data.Pod.UID)

	var singularity_command_pod []SingularityCommand

	for _, container := range containers {
		log.G(h.Ctx).Info("- Beginning script generation for container " + container.Name)
		singularityPrefix := SlurmConfigInst.SingularityPrefix
		if singularityAnnotation, ok := metadata.Annotations["slurm-job.vk.io/singularity-commands"]; ok {
			singularityPrefix += " " + singularityAnnotation
		}

		singularityMounts := ""
		if singMounts, ok := metadata.Annotations["slurm-job.vk.io/singularity-mounts"]; ok {
			singularityMounts = singMounts
		}

		singularityOptions := ""
		if singOpts, ok := metadata.Annotations["slurm-job.vk.io/singularity-options"]; ok {
			singularityOptions = singOpts
		}

		commstr1 := []string{"singularity", "exec", "--containall", "--nv", singularityMounts, singularityOptions}

		envs := prepareEnvs(h.Ctx, container)
		image := ""
		mounts, err := prepareMounts(h.Ctx, h.Config, data, container, filesPath)
		log.G(h.Ctx).Debug(mounts)
		if err != nil {
			statusCode := http.StatusInternalServerError
			w.WriteHeader(statusCode)
			w.Write([]byte("Error prepairing mounts. Check Slurm Plugin's logs"))
			log.G(h.Ctx).Error(err)
			os.RemoveAll(filesPath)
			return
		}

		image = container.Image
		if strings.HasPrefix(container.Image, "/") {
			if image_uri, ok := metadata.Annotations["slurm-job.vk.io/image-root"]; ok {
				image = image_uri + container.Image
			} else {
				log.G(h.Ctx).Info("- image-uri annotation not specified for path in remote filesystem")
			}
		} else {
			image = container.Image
		}

		log.G(h.Ctx).Debug("-- Appending all commands together...")
		singularity_command := append(commstr1, envs...)
		singularity_command = append(singularity_command, mounts...)
		singularity_command = append(singularity_command, image)
		singularity_command = append(singularity_command, container.Command...)
		singularity_command = append(singularity_command, container.Args...)

		singularity_command_pod = append(singularity_command_pod, SingularityCommand{command: singularity_command, containerName: container.Name})
	}

	path, err := produceSLURMScript(h.Ctx, h.Config, string(data.Pod.UID), filesPath, metadata, singularity_command_pod)
	if err != nil {
		statusCode := http.StatusInternalServerError
		w.WriteHeader(statusCode)
		w.Write([]byte("Error producing Slurm script. Check Slurm Plugin's logs"))
		log.G(h.Ctx).Error(err)
		os.RemoveAll(filesPath)
		return
	}
	out, err := SLURMBatchSubmit(h.Ctx, h.Config, path)
	if err != nil {
		statusCode := http.StatusInternalServerError
		w.WriteHeader(statusCode)
		w.Write([]byte("Error submitting Slurm script. Check Slurm Plugin's logs"))
		log.G(h.Ctx).Error(err)
		os.RemoveAll(filesPath)
		return
	}
	log.G(h.Ctx).Info(out)
	jid, err := handleJID(h.Ctx, data.Pod, h.JIDs, out, filesPath)
	if err != nil {
		statusCode := http.StatusInternalServerError
		w.WriteHeader(statusCode)
		w.Write([]byte("Error handling JID. Check Slurm Plugin's logs"))
		log.G(h.Ctx).Error(err)
		os.RemoveAll(filesPath)
		err = deleteContainer(h.Ctx, h.Config, string(data.Pod.UID), h.JIDs, filesPath)
		if err != nil {
			log.G(h.Ctx).Error(err)
		}
		return
	}

	//to be changed to commonIL.CreateStruct
	var returnedJID CreateStruct //returnValue
	var returnedJIDBytes []byte

	returnedJID = CreateStruct{PodUID: string(data.Pod.UID), PodJID: jid}

	returnedJIDBytes, err = json.Marshal(returnedJID)
	if err != nil {
		statusCode := http.StatusInternalServerError
		w.WriteHeader(statusCode)
		w.Write([]byte("Error marshaling JID. Check Slurm Plugin's logs"))
		log.G(h.Ctx).Error(err)
		return
	}

	w.WriteHeader(statusCode)

	if statusCode != http.StatusOK {
		w.Write([]byte("Some errors occurred while creating containers. Check Slurm Plugin's logs"))
	} else {
		w.Write(returnedJIDBytes)
	}
}
