# RSS

> All the thoughts fit to bot.

An RSS feed combining our blog's and podcasts' RSS feeds into
one feed for the past week's worth of content.

Used as the data source for our weekly newsletter.

## Developing

First you'll need a working [go installation],
and project cloned into your [go work environment]
(that is, `$GOPATH/src/github.com/thoughtbot/rss`).

  [go installation]: http://golang.org/doc/install
  [go work environment]: http://golang.org/doc/code.html

Run `bin/setup` to install the project's dependencies.

If you add or update a dependency,
run `godep save ./...` to vendor the changes.

## Testing

To test the `rss` package, run `go test ./...`.

## Deployment

The master branch on GitHub is automatically deployed
to the `thoughtbot-rss` app on Heroku
after the CI build passes.
>>>>>>> 678d7fc... Set up CircleCI with Bernerd's best practices
