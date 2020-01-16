package fexcel

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/blang/semver"
	wordwrap "github.com/mitchellh/go-wordwrap"
)

type GitHubRelease struct {
	HTMLUrl     string    `json:"html_url"`
	TagName     string    `json:"tag_name"`
	Name        string    `json:"name"`
	Body        string    `json:"body"`
	PublishedAt time.Time `json:"published_at"`
	Assets      []struct {
		DownloadUrl string `json:"browser_download_url"`
	} `json:"assets"`
}

type GitHubUpdateChecker struct {
}

func (g *GitHubUpdateChecker) UpdateCheck(w io.Writer) error {
	fmt.Fprintf(w, "\nChecking for updates... ")

	client := http.Client{
		Timeout: 2 * time.Second,
	}

	req, err := http.NewRequest("GET", "https://api.github.com/repos/onerobotics/fexcel/releases/latest", nil)
	if err != nil {
		return err
	}

	req.Header.Set("User-Agent", "fexcel-v"+Version)

	res, err := client.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	var release GitHubRelease
	err = json.NewDecoder(res.Body).Decode(&release)
	if err != nil {
		return err
	}

	currentVersion, err := semver.Make(Version)
	if err != nil {
		return err
	}
	latestVersion, err := semver.Make(release.TagName[1:]) // tag has v prefix
	if err != nil {
		return err
	}

	if latestVersion.GT(currentVersion) {
		fmt.Fprintf(w, "an update is available!\n")
		fmt.Fprintf(w, "\nfexcel %s was released on %s (you are on v%s).\n", release.TagName, release.PublishedAt.Format("January 01, 2006"), Version)
		fmt.Fprintf(w, "Download here: %s\n\n", release.HTMLUrl)
		fmt.Fprintf(w, "%s\n", release.Name)
		fmt.Fprintf(w, "%s\n\n", wordwrap.WrapString(release.Body, 80))
	} else {
		fmt.Fprintf(w, "You are on the latest version: v%s\n", Version)
	}

	return nil
}
