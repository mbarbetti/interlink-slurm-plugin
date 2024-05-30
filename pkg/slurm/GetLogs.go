package slurm

import (
	"errors"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/containerd/containerd/log"
)

// GetLogsHandler reads Jobs' output file to return what's logged inside.
// What's returned is based on the provided parameters (Tail/LimitBytes/Timestamps/etc)
func (h *PluginHandler) GetLogsHandler(w http.ResponseWriter, r *http.Request) {
	log.G(h.Ctx).Info("Docker Sidecar: received GetLogs call")
	reqTime := time.Now()

	bodyBytes, _ := ReadRequestBody(r, w, h)
	logStruct, statusCode := ParseLogsJson(bodyBytes, w, h)

	path := h.Config.DataRootFolder + logStruct.Namespace + "-" + logStruct.PodUID
	var output []byte
	if logStruct.Opts.Timestamps {
		log.G(h.Ctx).Error(errors.New("Not Implemented"))
		statusCode = http.StatusInternalServerError
		w.WriteHeader(statusCode)
		return
	} else {
		log.G(h.Ctx).Info("Reading  " + path + "/" + logStruct.ContainerName + ".out")
		containerOutput, err1 := os.ReadFile(path + "/" + logStruct.ContainerName + ".out")
		if err1 != nil {
			log.G(h.Ctx).Error("Failed to read container logs.")
		}
		jobOutput, err2 := os.ReadFile(path + "/" + "job.out")
		if err2 != nil {
			log.G(h.Ctx).Error("Failed to read job logs.")
		}

		if err1 != nil && err2 != nil {
			log.G(h.Ctx).Error("Failed to retrieve logs.")
			statusCode = http.StatusInternalServerError
			w.WriteHeader(statusCode)
			return
		}

		output = append(output, jobOutput...)
		output = append(output, containerOutput...)

	}

	var returnedLogs string

	if logStruct.Opts.Tail != 0 {
		var lastLines []string

		splittedLines := strings.Split(string(output), "\n")

		if logStruct.Opts.Tail > len(splittedLines) {
			lastLines = splittedLines
		} else {
			lastLines = splittedLines[len(splittedLines)-logStruct.Opts.Tail-1:]
		}

		for _, line := range lastLines {
			returnedLogs += line + "\n"
		}
	} else if logStruct.Opts.LimitBytes != 0 {
		var lastBytes []byte
		if logStruct.Opts.LimitBytes > len(output) {
			lastBytes = output
		} else {
			lastBytes = output[len(output)-logStruct.Opts.LimitBytes-1:]
		}

		returnedLogs = string(lastBytes)
	} else {
		returnedLogs = string(output)
	}

	if logStruct.Opts.Timestamps && (logStruct.Opts.SinceSeconds != 0 || !logStruct.Opts.SinceTime.IsZero()) {
		temp := returnedLogs
		returnedLogs = ""
		splittedLogs := strings.Split(temp, "\n")
		timestampFormat := "2006-01-02T15:04:05.999999999Z"

		for _, Log := range splittedLogs {
			part := strings.SplitN(Log, " ", 2)
			timestampString := part[0]
			timestamp, err := time.Parse(timestampFormat, timestampString)
			if err != nil {
				continue
			}
			if logStruct.Opts.SinceSeconds != 0 {
				if reqTime.Sub(timestamp).Seconds() > float64(logStruct.Opts.SinceSeconds) {
					returnedLogs += Log + "\n"
				}
			} else {
				if timestamp.Sub(logStruct.Opts.SinceTime).Seconds() >= 0 {
					returnedLogs += Log + "\n"
				}
			}
		}
	}

	if statusCode != http.StatusOK {
		w.Write([]byte("Some errors occurred while checking container status. Check Slurm Plugin's logs"))
	} else {
		w.WriteHeader(statusCode)
		w.Write([]byte(returnedLogs))
	}
}
