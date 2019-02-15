package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"strconv"
)

func infocmd(opt Options) error {
	user := nvls(opt.Info.User, EnvUser)
	repo := nvls(opt.Info.Repo, EnvRepo)
	token := nvls(opt.Info.Token, EnvToken)
	tag := opt.Info.Tag

	if user == "" || repo == "" {
		return fmt.Errorf("user and repo need to be passed as arguments")
	}

	/* find regular git tags */
	allTags, err := Tags(user, repo, token)
	if err != nil {
		return fmt.Errorf("could not fetch tags, %v", err)
	}
	if len(allTags) == 0 {
		return fmt.Errorf("no tags available for %v/%v", user, repo)
	}

	/* list all tags */
	tags := make([]Tag, 0, len(allTags))
	for _, t := range allTags {
		/* if the user only requested see one tag, skip the ones that
		 * don't match */
		if tag != "" && t.Name != tag {
			continue
		}
		tags = append(tags, t)
	}

	/* if no tags conformed to the users' request, exit */
	if len(tags) == 0 {
		return fmt.Errorf("no tag '%v' was found for %v/%v", tag, user, repo)
	}

	fmt.Println("git tags:")
	for _, t := range tags {
		fmt.Println("-", t.String())
	}

	/* list releases + assets */
	var releases []Release
	if tag == "" {
		/* get all releases */
		vprintf("%v/%v: getting information for all releases\n", user, repo)
		releases, err = Releases(user, repo, token)
		if err != nil {
			return err
		}
	} else {
		/* get only one release */
		vprintf("%v/%v/%v: getting information for the release\n", user, repo, tag)
		release, err := ReleaseOfTag(user, repo, tag, token)
		if err != nil {
			return err
		}
		releases = []Release{*release}
	}

	/* if no tags conformed to the users' request, exit */
	if len(releases) == 0 {
		return fmt.Errorf("no release(s) were found for %v/%v (%v)", user, repo, tag)
	}

	fmt.Println("releases:")
	for _, release := range releases {
		fmt.Println("-", release.String())
	}

	return nil
}

func uploadcmd(opt Options) error {
	user := nvls(opt.Upload.User, EnvUser)
	repo := nvls(opt.Upload.Repo, EnvRepo)
	token := nvls(opt.Upload.Token, EnvToken)
	tag := opt.Upload.Tag
	name := opt.Upload.Name
	label := opt.Upload.Label
	file := opt.Upload.File

	vprintln("uploading...")

	if file == nil {
		return fmt.Errorf("provided file was not valid")
	}
	defer file.Close()

	if err := ValidateCredentials(user, repo, token, tag); err != nil {
		return err
	}

	/* find the release corresponding to the entered tag, if any */
	rel, err := ReleaseOfTag(user, repo, tag, token)
	if err != nil {
		return err
	}

	v := url.Values{}
	v.Set("name", name)
	if label != "" {
		v.Set("label", label)
	}

	url := rel.CleanUploadUrl() + "?" + v.Encode()

	resp, err := DoAuthRequest("POST", url, "application/octet-stream",
		token, nil, file)
	if resp != nil {
		defer resp.Body.Close()
	}
	if err != nil {
		return fmt.Errorf("can't create upload request to %v, %v", url, err)
	}

	vprintln("RESPONSE:", resp)
	if resp.StatusCode != http.StatusCreated {
		if msg, err := ToMessage(resp.Body); err == nil {
			return fmt.Errorf("could not upload, status code (%v), %v",
				resp.Status, msg)
		} else {
			return fmt.Errorf("could not upload, status code (%v)", resp.Status)
		}
	}

	if VERBOSITY != 0 {
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return fmt.Errorf("error while reading response, %v", err)
		}
		vprintln("BODY:", string(body))
	}

	return nil
}

func downloadcmd(opt Options) error {
	user := nvls(opt.Download.User, EnvUser)
	repo := nvls(opt.Download.Repo, EnvRepo)
	token := nvls(opt.Download.Token, EnvToken)
	tag := opt.Download.Tag
	name := opt.Download.Name
	latest := opt.Download.Latest

	vprintln("downloading...")

	if err := ValidateTarget(user, repo, tag, latest); err != nil {
		return err
	}

	// Find the release corresponding to the entered tag, if any.
	var rel *Release
	var err error
	if latest {
		rel, err = LatestRelease(user, repo, token)
	} else {
		rel, err = ReleaseOfTag(user, repo, tag, token)
	}
	if err != nil {
		return err
	}

	assetId := 0
	for _, asset := range rel.Assets {
		if asset.Name == name {
			assetId = asset.Id
		}
	}

	if assetId == 0 {
		return fmt.Errorf("coud not find asset named %s", name)
	}

	var resp *http.Response
	var url string
	if token == "" {
		url = GH_URL + fmt.Sprintf("/%s/%s/releases/download/%s/%s", user, repo, tag, name)
		resp, err = http.Get(url)
	} else {
		url = ApiURL() + fmt.Sprintf(ASSET_DOWNLOAD_URI, user, repo, assetId)
		resp, err = DoAuthRequest("GET", url, "", token, map[string]string{
			"Accept": "application/octet-stream",
		}, nil)
	}
	if resp != nil {
		defer resp.Body.Close()
	}
	if err != nil {
		return fmt.Errorf("could not fetch releases, %v", err)
	}

	vprintln("GET", url, "->", resp)

	contentLength, err := strconv.ParseInt(resp.Header.Get("Content-Length"), 10, 64)
	if err != nil {
		return err
	}

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("github did not respond with 200 OK but with %v", resp.Status)
	}

	out, err := os.Create(name)
	if err != nil {
		return fmt.Errorf("could not create file %s", name)
	}
	defer out.Close()

	n, err := io.Copy(out, resp.Body)
	if n != contentLength {
		return fmt.Errorf("downloaded data did not match content length %d != %d", contentLength, n)
	}
	return err
}

func ValidateTarget(user, repo, tag string, latest bool) error {
	if user == "" {
		return fmt.Errorf("empty user")
	}
	if repo == "" {
		return fmt.Errorf("empty repo")
	}
	if tag == "" && !latest {
		return fmt.Errorf("empty tag")
	}
	return nil
}

func ValidateCredentials(user, repo, token, tag string) error {
	if err := ValidateTarget(user, repo, tag, false); err != nil {
		return err
	}
	if token == "" {
		return fmt.Errorf("empty token")
	}
	return nil
}

func releasecmd(opt Options) error {
	cmdopt := opt.Release
	user := nvls(cmdopt.User, EnvUser)
	repo := nvls(cmdopt.Repo, EnvRepo)
	token := nvls(cmdopt.Token, EnvToken)
	tag := cmdopt.Tag
	name := nvls(cmdopt.Name, tag)
	desc := nvls(cmdopt.Desc, tag)
	target := nvls(cmdopt.Target)
	draft := cmdopt.Draft
	prerelease := cmdopt.Prerelease

	vprintln("releasing...")

	if err := ValidateCredentials(user, repo, token, tag); err != nil {
		return err
	}

	params := ReleaseCreate{
		TagName:         tag,
		TargetCommitish: target,
		Name:            name,
		Body:            desc,
		Draft:           draft,
		Prerelease:      prerelease,
	}

	/* encode params as json */
	payload, err := json.Marshal(params)
	if err != nil {
		return fmt.Errorf("can't encode release creation params, %v", err)
	}
	reader := bytes.NewReader(payload)

	uri := fmt.Sprintf("/repos/%s/%s/releases", user, repo)
	resp, err := DoAuthRequest("POST", ApiURL()+uri, "application/json",
		token, nil, reader)
	if resp != nil {
		defer resp.Body.Close()
	}
	if err != nil {
		return fmt.Errorf("while submitting %v, %v", string(payload), err)
	}

	vprintln("RESPONSE:", resp)
	if resp.StatusCode != http.StatusCreated {
		if resp.StatusCode == 422 {
			return fmt.Errorf("github returned %v (this is probably because the release already exists)",
				resp.Status)
		}
		return fmt.Errorf("github returned %v", resp.Status)
	}

	if VERBOSITY != 0 {
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return fmt.Errorf("error while reading response, %v", err)
		}
		vprintln("BODY:", string(body))
	}

	return nil
}

func editcmd(opt Options) error {
	cmdopt := opt.Edit
	user := nvls(cmdopt.User, EnvUser)
	repo := nvls(cmdopt.Repo, EnvRepo)
	token := nvls(cmdopt.Token, EnvToken)
	tag := cmdopt.Tag
	name := nvls(cmdopt.Name, tag)
	desc := nvls(cmdopt.Desc, tag)
	draft := cmdopt.Draft
	prerelease := cmdopt.Prerelease

	vprintln("editing...")

	if err := ValidateCredentials(user, repo, token, tag); err != nil {
		return err
	}

	id, err := IdOfTag(user, repo, tag, token)
	if err != nil {
		return err
	}

	vprintf("release %v has id %v\n", tag, id)

	/* the release create struct works for editing releases as well */
	params := ReleaseCreate{
		TagName:    tag,
		Name:       name,
		Body:       desc,
		Draft:      draft,
		Prerelease: prerelease,
	}

	/* encode the parameters as JSON, as required by the github API */
	payload, err := json.Marshal(params)
	if err != nil {
		return fmt.Errorf("can't encode release creation params, %v", err)
	}

	uri := fmt.Sprintf("/repos/%s/%s/releases/%d", user, repo, id)
	resp, err := DoAuthRequest("PATCH", ApiURL()+uri, "application/json",
		token, nil, bytes.NewReader(payload))
	if resp != nil {
		defer resp.Body.Close()
	}
	if err != nil {
		return fmt.Errorf("while submitting %v, %v", string(payload), err)
	}

	vprintln("RESPONSE:", resp)
	if resp.StatusCode != http.StatusOK {
		if resp.StatusCode == 422 {
			return fmt.Errorf("github returned %v (this is probably because the release already exists)",
				resp.Status)
		}
		return fmt.Errorf("github returned unexpected status code %v", resp.Status)
	}

	if VERBOSITY != 0 {
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return fmt.Errorf("error while reading response, %v", err)
		}
		vprintln("BODY:", string(body))
	}

	return nil
}

func deletecmd(opt Options) error {
	user, repo, token, tag := nvls(opt.Delete.User, EnvUser),
		nvls(opt.Delete.Repo, EnvRepo),
		nvls(opt.Delete.Token, EnvToken),
		opt.Delete.Tag
	vprintln("deleting...")

	id, err := IdOfTag(user, repo, tag, token)
	if err != nil {
		return err
	}

	vprintf("release %v has id %v\n", tag, id)

	resp, err := httpDelete(ApiURL()+fmt.Sprintf("/repos/%s/%s/releases/%d",
		user, repo, id), token)
	if resp != nil {
		defer resp.Body.Close()
	}
	if err != nil {
		return fmt.Errorf("release deletion unsuccesful, %v", err)
	}

	if resp.StatusCode != http.StatusNoContent {
		return fmt.Errorf("could not delete the release corresponding to tag %s on repo %s/%s",
			tag, user, repo)
	}

	return nil
}

func httpDelete(url, token string) (*http.Response, error) {
	resp, err := DoAuthRequest("DELETE", url, "application/json", token, nil, nil)
	if err != nil {
		return nil, err
	}

	return resp, nil
}
