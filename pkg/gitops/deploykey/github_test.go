package deploykey

import "testing"

func TestGetGitHubOwnerRepoFromRepoURL(t *testing.T) {
	testcases := []struct {
		title                string
		repoURL, owner, repo string
		ok                   bool
	}{
		{
			title:   "git@github.com",
			repoURL: "git@github.com:myorg/configrepo.git",
			owner:   "myorg",
			repo:    "configrepo",
			ok:      true,
		},
		{
			title:   "ssh://git@github.com",
			repoURL: "ssh://git@github.com/myorg/configrepo.git",
			owner:   "myorg",
			repo:    "configrepo",
			ok:      true,
		},
		{
			title:   "non-gh url",
			repoURL: "git@gitlab.com:gitlab-org/gitlab.git",
			owner:   "",
			repo:    "",
			ok:      false,
		},
	}

	for i := range testcases {
		tc := testcases[i]

		t.Run(tc.title, func(t *testing.T) {
			owner, repo, ok := getGitHubOwnerRepoFromRepoURL(tc.repoURL)

			if owner != tc.owner {
				t.Errorf("unexpected owner: want %s, got %s", tc.owner, owner)
			}

			if repo != tc.repo {
				t.Errorf("unexpected repo: want %s, got %s", tc.repo, repo)
			}

			if ok != tc.ok {
				t.Errorf("unexpected ok: want %v, got %v", tc.ok, ok)
			}
		})
	}
}
