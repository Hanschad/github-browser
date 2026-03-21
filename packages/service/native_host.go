package main

import (
	"bufio"
	"encoding/binary"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
)

const NativeHostName = "com.github.browser"

type NativeRequest struct {
	Action string  `json:"action"`
	URL    string  `json:"url,omitempty"`
	IDE    string  `json:"ide,omitempty"`
	Config *Config `json:"config,omitempty"`
}

type NativeResponse struct {
	Status  string  `json:"status"`
	Message string  `json:"message,omitempty"`
	Version string  `json:"version,omitempty"`
	Path    string  `json:"path,omitempty"`
	Mode    string  `json:"mode,omitempty"`
	Config  *Config `json:"config,omitempty"`
}

func RunNativeHost(service *Service) error {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)

	for {
		request, err := readNativeRequest(reader)
		if errors.Is(err, io.EOF) {
			return nil
		}
		if err != nil {
			return err
		}

		response := service.handleNativeRequest(request)
		if err := writeNativeResponse(writer, response); err != nil {
			return err
		}
	}
}

func (s *Service) handleNativeRequest(request *NativeRequest) *NativeResponse {
	if request == nil {
		return &NativeResponse{
			Status:  "error",
			Message: "empty native request",
		}
	}

	switch request.Action {
	case "health":
		health := s.Health("native")
		return &NativeResponse{
			Status:  health.Status,
			Message: "Native host ready",
			Version: health.Version,
			Mode:    health.Mode,
		}
	case "open":
		if request.URL == "" {
			return &NativeResponse{
				Status:  "error",
				Message: "url is required",
			}
		}

		response, err := s.Open(OpenRequest{
			URL: request.URL,
			IDE: request.IDE,
		})
		if err != nil {
			return &NativeResponse{
				Status:  "error",
				Message: err.Error(),
			}
		}

		return &NativeResponse{
			Status:  response.Status,
			Message: response.Message,
			Path:    response.Path,
			Mode:    "native",
		}
	case "getConfig":
		return &NativeResponse{
			Status: "ok",
			Config: s.config,
			Mode:   "native",
		}
	case "updateConfig":
		if request.Config == nil {
			return &NativeResponse{
				Status:  "error",
				Message: "config is required",
			}
		}

		if err := s.UpdateConfig(request.Config); err != nil {
			return &NativeResponse{
				Status:  "error",
				Message: err.Error(),
			}
		}

		return &NativeResponse{
			Status:  "ok",
			Message: "Config updated",
			Config:  s.config,
			Mode:    "native",
		}
	default:
		return &NativeResponse{
			Status:  "error",
			Message: fmt.Sprintf("unsupported action: %s", request.Action),
		}
	}
}

func readNativeRequest(reader io.Reader) (*NativeRequest, error) {
	header := make([]byte, 4)
	if _, err := io.ReadFull(reader, header); err != nil {
		return nil, err
	}

	length := binary.LittleEndian.Uint32(header)
	if length == 0 {
		return nil, fmt.Errorf("invalid native message length: %d", length)
	}

	payload := make([]byte, length)
	if _, err := io.ReadFull(reader, payload); err != nil {
		return nil, err
	}

	var request NativeRequest
	if err := json.Unmarshal(payload, &request); err != nil {
		return nil, err
	}

	return &request, nil
}

func writeNativeResponse(writer *bufio.Writer, response *NativeResponse) error {
	payload, err := json.Marshal(response)
	if err != nil {
		return err
	}

	header := make([]byte, 4)
	binary.LittleEndian.PutUint32(header, uint32(len(payload)))

	if _, err := writer.Write(header); err != nil {
		return err
	}
	if _, err := writer.Write(payload); err != nil {
		return err
	}

	return writer.Flush()
}
