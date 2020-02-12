# Sensu Redist(ribute)

[![GoDoc](https://godoc.org/github.com/julian7/sensu-redist?status.svg)](https://godoc.org/github.com/julian7/sensu-redist)

Sensu Redist takes a [Sensu Asset file](https://docs.sensu.io/sensu-go/latest/reference/assets/) with all its builds, and uploads them to another place. This makes Sensu Assets available in air gapped systems.

## Usage

```text
  sensu-redist NAME [flags]

Flags:
  -h, --help             help for sensu-redist
  -s, --source string    Source Asset URL (required)
  -t, --target string    Target URL stub (required)
  -u, --upload string    Upload to (scp target, required)
      --version          version for sensu-redist
  -w, --workdir string   Work dir, not deleted (default: random temporary dir, deleted)
```

Example:

```shell
$ sensu-redist --source https://github.com/julian7/sensu-base-checks/releases/download/v0.2.2/sensu-base-checks-v0.2.2-asset.yml --target https://artifacts.example.com/sensu-base-checks/ --upload user@air-gapped.example.com:/var/www/public/sensu-base-checks
The asset is now available at https://artifacts.example.com/sensu-base-checks/sensu-base-checks-v0.2.2-asset.yml
```

The command downloads the asset file from `--source`, then all its referenced build targets to a work directory (in this case, an ephemeral temp directory), it changes URLs in the asset file to match with `--target` stub, and uploads them to `--upload`, via scp. It doesn't handle secret keys, usernames, passwords, but relies on your configuration for things like port settings, or agent usage.

## Legal

This project is licensed under [Blue Oak Model License v1.0.0](https://blueoakcouncil.org/license/1.0.0). It is not registered either at OSI or GNU, therefore GitHub and Google are widely looking at the other direction. However, this is the license I'm most happy with: you can read and understand it with no legal degree, and there are no hidden or cryptic meanings in it.

The project is also governed with [Contributor Covenant](https://contributor-covenant.org/)'s [Code of Conduct](https://www.contributor-covenant.org/version/1/4/) in mind. I'm not copying it here, as a pledge for taking the verbatim version by the word, and we are not going to modify it in any way.

## In case of trouble

Open a ticket, perhaps a pull request. We support [GitHub Flow](https://guides.github.com/introduction/flow/). You might want to [fork](https://guides.github.com/activities/forking/) this project first.
