Contributing
============

We love pull requests from everyone.
By participating in this project,
you agree to abide by the thoughtbot [code of conduct].

  [code of conduct]: https://thoughtbot.com/open-source-code-of-conduct

We expect everyone to follow the code of conduct
anywhere in thoughtbot's project codebases,
issue trackers, chatrooms, and mailing lists.

Fork the repo.

Get a working [Go installation] and clone the project:

  [Go installation]: http://golang.org/doc/install

Run `./bin/setup` to install the project's dependencies.

To test the `rss` package, run `go test -v ./...`.

Make your change, with new passing tests.

Run `go run main.go` to see the change at `localhost:8080` in a web browser.

Push to your fork. Write a [good commit message][commit]. Submit a pull request.

  [commit]: http://tbaggery.com/2008/04/19/a-note-about-git-commit-messages.html

Others will give constructive feedback.
This is a time for discussion and improvements,
and making the necessary changes will be required before we can
merge the contribution.

The master branch on GitHub is automatically deployed
to the `thoughtbot-rss` app on Heroku
after the CI build passes.
