# mastodon audit trail backup

This repository backs up my
[follower list](followers.txt),
[following list](following.txt),
[blocked accounts list](blocked_accounts.txt) and
[muted accounts list](mutes.txt) periodically using GitHub Actions.

## Set up

1. Fork this repository.
1. `git rm *.txt` and commit.
1. Create a Mastodon app.
1. Go to Repository Settings &rarr; Secrets and add secrets from your Mastodon app:

   - MASTODON_CLIENT_KEY
   - MASTODON_CLIENT_SECRET
   - MASTODON_ACCESS_TOKEN

1. Go to Repository Settings &rarr; Secrets and add the base url of your Mastodon instance:

   - MASTODON_INSTANCE

1. See [.github/workflows/update.yml](/.github/workflows/update.yml) and modify the cron schedule (in UTC) as you see fit.

1. Commit and push. Once the time arrives, the cron would work, and commit the lists into `.txt` files and push the changes.