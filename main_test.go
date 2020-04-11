package main

import (
	"io"
	"reflect"
	"strings"
	"testing"
	"testing/iotest"
)

func TestConfirm(t *testing.T) {
	tests := []struct {
		name    string
		reader  io.Reader
		want    bool
		wantErr bool
	}{
		{
			name:    "y",
			reader:  strings.NewReader("y"),
			want:    true,
			wantErr: false,
		},
		{
			name:    "y",
			reader:  strings.NewReader("no"),
			want:    false,
			wantErr: false,
		},
		{
			name:    "",
			reader:  iotest.DataErrReader(strings.NewReader("")),
			want:    false,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := confirm(tt.reader)
			if (err != nil) != tt.wantErr {
				t.Errorf("Confirm() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("Confirm() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_githubRepo(t *testing.T) {
	tests := []struct {
		name  string
		repo  string
		want  githubRepo
		want1 bool
	}{
		{
			"github.com/owner/repo",
			"github.com/owner/repo",
			githubRepo{
				path:  "github.com/owner/repo",
				owner: "owner",
				repo:  "repo",
			},
			true,
		},
		{
			"github.com/owner/repo.xxx",
			"github.com/owner/repo.xxx",
			githubRepo{
				path:  "github.com/owner/repo.xxx",
				owner: "owner",
				repo:  "repo.xxx",
			},
			true,
		},
		{
			"github.com/owner/repo_xxx",
			"github.com/owner/repo_xxx",
			githubRepo{
				path:  "github.com/owner/repo_xxx",
				owner: "owner",
				repo:  "repo_xxx",
			},
			true,
		},
		{
			"github.com/owner/repo-xxx",
			"github.com/owner/repo-xxx",
			githubRepo{
				path:  "github.com/owner/repo-xxx",
				owner: "owner",
				repo:  "repo-xxx",
			},
			true,
		},
		{
			"github.com/owner/repo/v2",
			"github.com/owner/repo/v2",
			githubRepo{
				path:  "github.com/owner/repo",
				owner: "owner",
				repo:  "repo",
			},
			true,
		},
		{
			"golang.org/x/lint",
			"golang.org/x/lint",
			githubRepo{
				path:  "github.com/golang/lint",
				owner: "golang",
				repo:  "lint",
			},
			true,
		},
		{
			"xxx.com/owner/repo",
			"xxx.com/owner/repo",
			githubRepo{},
			false,
		},
		{
			"empty",
			"",
			githubRepo{},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1 := isGithubRepo(tt.repo)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("githubRepo() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("githubRepo() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}
