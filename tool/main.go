package main

import (
	"context"
	"flag"
	"fmt"
	"os"

	"github.com/mattn/go-mastodon"
)

type runFunc func(context.Context, *mastodon.Client, mastodon.ID) ([]*mastodon.Account, error)

var cmds = map[string]runFunc{
	"id":        nil, // exits early, since we're passing the result of this to every other function
	"following": following,
	"followers": followers,
	"mutes":     muted,
	"blocks":    blocked,
}

func init() { flag.Parse() }

func main() {
	if flag.NArg() != 1 {
		panic(fmt.Sprintf("requires 1 positional argument (cmd name); got %d args", flag.NArg()))
	}
	cmd, ok := cmds[flag.Arg(0)]
	if !ok {
		panic("unknown command: " + flag.Arg(0))
	}

	var (
		clientKey        = "MASTODON_CLIENT_KEY"
		clientSecret     = "MASTODON_CLIENT_SECRET"
		accessToken      = "MASTODON_ACCESS_TOKEN"
		mastodonInstance = "MASTODON_INSTANCE"
	)
	for _, v := range []*string{&clientKey, &clientSecret, &accessToken, &mastodonInstance} {
		if vv := os.Getenv(*v); vv == "" {
			panic(*v + " env var not set")
		} else {
			*v = vv
		}
	}
	ctx := context.Background()

	client := mastodon.NewClient(&mastodon.Config{
		Server:       mastodonInstance,
		ClientID:     clientKey,
		ClientSecret: clientSecret,
		AccessToken:  accessToken,
	})
	id, err := selfID(ctx, client)
	if err != nil {
		panic(err)
	}
	if flag.Arg(0) == "id" {
		fmt.Print(id)
		return
	}

	if err := print(ctx, client, id, cmd); err != nil {
		panic(err)
	}
}

func selfID(ctx context.Context, c *mastodon.Client) (mastodon.ID, error) {
	acc, err := c.GetAccountCurrentUser(ctx)
	if err != nil {
		return "", err
	}
	return acc.ID, nil
}

func print(ctx context.Context, c *mastodon.Client, id mastodon.ID, fnc runFunc) error {
	out, err := fnc(ctx, c, id)
	if err != nil {
		return err
	}

	for _, acc := range out {
		fmt.Printf("%s,%s\n", acc.ID, acc.Acct)
	}
	return nil
}

func following(ctx context.Context, c *mastodon.Client, id mastodon.ID) ([]*mastodon.Account, error) {
	var accounts []*mastodon.Account
	var pg mastodon.Pagination
	for {
		fs, err := c.GetAccountFollowing(ctx, id, &pg)
		if err != nil {
			return nil, err
		}
		accounts = append(accounts, fs...)
		if pg.MaxID == "" {
			break
		}
		pg.SinceID = ""
		pg.MinID = ""
	}
	return accounts, nil
}

func followers(ctx context.Context, c *mastodon.Client, id mastodon.ID) ([]*mastodon.Account, error) {
	var accounts []*mastodon.Account
	var pg mastodon.Pagination
	for {
		fs, err := c.GetAccountFollowers(ctx, id, &pg)
		if err != nil {
			return nil, err
		}
		accounts = append(accounts, fs...)
		if pg.MaxID == "" {
			break
		}
		pg.SinceID = ""
		pg.MinID = ""
	}
	return accounts, nil
}

func blocked(ctx context.Context, c *mastodon.Client, id mastodon.ID) ([]*mastodon.Account, error) {
	var accounts []*mastodon.Account
	var pg mastodon.Pagination
	for {
		fs, err := c.GetBlocks(ctx, &pg)
		if err != nil {
			return nil, err
		}
		accounts = append(accounts, fs...)
		if pg.MaxID == "" {
			break
		}
		pg.SinceID = ""
		pg.MinID = ""
	}
	return accounts, nil
}

func muted(ctx context.Context, c *mastodon.Client, id mastodon.ID) ([]*mastodon.Account, error) {
	var accounts []*mastodon.Account
	var pg mastodon.Pagination
	for {
		fs, err := c.GetMutes(ctx, &pg)
		if err != nil {
			return nil, err
		}
		accounts = append(accounts, fs...)
		if pg.MaxID == "" {
			break
		}
		pg.SinceID = ""
		pg.MinID = ""
	}
	return accounts, nil
}
