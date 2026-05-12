package suid

import (
	"fmt"
	"strings"

	"github.com/ctrsploit/sploit-spec/pkg/printer"
	"github.com/ctrsploit/sploit-spec/pkg/result"
	"github.com/ctrsploit/sploit-spec/pkg/result/item"
)

type FileItem struct {
	SubTitle  result.SubTitle `json:"-"`
	Path      item.Short      `json:"path"`
	Mode      item.Short      `json:"mode"`
	Owner     item.Short      `json:"owner"`
	Size      item.Short      `json:"size"`
	Usability item.Short      `json:"usability"`
}

type Result struct {
	Name  result.Title `json:"name"`
	Count item.Short   `json:"count"`
	Files []FileItem   `json:"files"`
}

func Human(files []File) Result {
	human := Result{
		Name: result.Title{
			Name: "SUID Files",
		},
		Count: item.Short{
			Name:   "total",
			Result: fmt.Sprintf("%d", len(files)),
		},
	}

	for _, file := range files {
		usability := "non-root-owner"
		if file.RootOwned && file.Executable {
			usability = "root-owned executable"
		} else if file.RootOwned {
			usability = "root-owned non-executable"
		} else if file.Executable {
			usability = "non-root-owned executable"
		}

		human.Files = append(human.Files, FileItem{
			SubTitle: result.SubTitle{
				Name: file.Path,
			},
			Path: item.Short{
				Name:   "path",
				Result: file.Path,
			},
			Mode: item.Short{
				Name:   "mode",
				Result: file.Mode,
			},
			Owner: item.Short{
				Name:   "owner",
				Result: fmt.Sprintf("%d:%d", file.UID, file.GID),
			},
			Size: item.Short{
				Name:   "size",
				Result: fmt.Sprintf("%d", file.Size),
			},
			Usability: item.Short{
				Name:   "usability",
				Result: usability,
			},
		})
	}
	return human
}

func Print(opts Options) error {
	files, err := Find(opts)
	if err != nil {
		return err
	}
	u := result.Union{
		Machine: files,
		Human:   Human(files),
	}
	fmt.Println(strings.TrimRight(printer.Printer.Print(u), "\n"))
	return nil
}
