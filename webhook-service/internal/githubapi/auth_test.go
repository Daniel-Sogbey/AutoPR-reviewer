package githubapi

import "testing"

func TestGithubAuth(t *testing.T) {

	_, err := Auth()
	if err != nil {
		t.Errorf("expecting successful request but got error: %v ", err)
	}
}
