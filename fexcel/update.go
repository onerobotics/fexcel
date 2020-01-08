package fexcel

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/blang/semver"
)

type releaseResponse struct {
	Url     string `json:"url"`
	TagName string `json:"tag_name"`
	Assets  []struct {
		DownloadUrl string `json:"browser_download_url"`
	} `json:"assets"`
}

func CheckForUpdates(w io.Writer) error {
	fmt.Fprintf(w, "\nChecking for updates... ")

	githubClient := http.Client{
		Timeout: 2 * time.Second,
	}

	req, err := http.NewRequest("GET", "https://api.github.com/repos/onerobotics/fexcel/releases/latest", nil)
	if err != nil {
		return err
	}

	req.Header.Set("User-Agent", "fexcel-v"+Version)

	res, err := githubClient.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	var response releaseResponse
	err = json.NewDecoder(res.Body).Decode(&response)
	if err != nil {
		return err
	}

	if len(response.Assets) != 1 {
		return fmt.Errorf("Invalid # of assets: %d", len(response.Assets))
	}

	currentVersion, err := semver.Make(Version)
	if err != nil {
		return err
	}
	latestVersion, err := semver.Make(response.TagName[1:]) // tag has v prefix
	if err != nil {
		return err
	}

	if latestVersion.GT(currentVersion) {
		fmt.Fprintf(w, "A new version is available!\n")
		fmt.Fprintf(w, "  The latest version is: %s. You are on v%s\n", response.TagName, Version)
		fmt.Fprintf(w, "  You can view the latest release here:\n    %s\n", response.Url)
		fmt.Fprintf(w, "  You can download the latest release here:\n    %s\n", response.Assets[0].DownloadUrl)
	} else {
		fmt.Fprintf(w, "You are on the latest version: v%s\n", Version)
	}

	return nil
}
