package bitbucket

import (
	"github.com/ktrysmt/go-bitbucket"
)

type Bitbucket struct {
	HttpClient *bitbucket.Client
}

// NewBitbucketClient creates new http client to connect with bitbucket with provided credentials.
func NewBitbucketClient(key, secret string) *Bitbucket {
	var (
		b Bitbucket
	)
	b.HttpClient = bitbucket.NewOAuthClientCredentials(key, secret)
	return &b
}

// GetFileContent is a function to get content from text file in the repo
func (b *Bitbucket) GetFileContent(owner, repoName, path string) (*bitbucket.RepositoryBlob, error) {
	branch := b.GetMainBranch(owner, repoName)
	opt := &bitbucket.RepositoryBlobOptions{
		Owner:    owner,
		RepoSlug: repoName,
		Ref:      branch,
		Path:     path,
	}

	content, err := b.HttpClient.Repositories.Repository.GetFileBlob(opt)
	return content, err
}

func (b *Bitbucket) GetMainBranch(owner, repoName string) string {
	ro := &bitbucket.RepositoryOptions{
		Owner:    owner,
		RepoSlug: repoName,
	}
	repo, _ := b.HttpClient.Repositories.Repository.Get(ro)
	return repo.Mainbranch.Name
}