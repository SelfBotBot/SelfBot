package owo

import "fmt"

type PomfResponse struct {
	Success     bool   `json:"success"`
	Errorcode   int    `json:"errorcode"`
	Description string `json:"description"`
	Files       []struct {
		Name   string `json:"name"`
		RawUrl string `json:"url"`
		Hash   string `json:"hash"`
		Size   int    `json:"size"`
	} `json:"files"`
}

func (d PomfResponse) Err() error {
	if d.Success {
		return nil
	}

	return PomfError{
		ErrorCode:   d.Errorcode,
		Description: d.Description,
	}
}

type PomfError struct {
	ErrorCode   int    `json:"errorcode"`
	Description string `json:"description"`
}

func (e PomfError) Error() string {
	return fmt.Sprintf("pomf error (%d) %q", e.ErrorCode, e.Description)
}
