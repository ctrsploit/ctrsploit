package pkg

import (
	"context"
	"fmt"
	"github.com/google/go-containerregistry/pkg/name"
	"github.com/google/go-containerregistry/pkg/v1/remote"
)

func list(ctx context.Context, src string) error {
	repo, err := name.NewRepository(src)
	if err != nil {
		return fmt.Errorf("parsing repo %q: %w", src, err)
	}

	puller, err := remote.NewPuller()
	if err != nil {
		return err
	}

	lister, err := puller.Lister(ctx, repo)
	if err != nil {
		return fmt.Errorf("reading tags for %s: %w", repo, err)
	}

	for lister.HasNext() {
		tags, err := lister.Next(ctx)
		if err != nil {
			return err
		}
		for _, tag := range tags.Tags {
			fmt.Println(tag)
		}
	}
	return nil
}
