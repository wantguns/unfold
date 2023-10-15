# Unfold

Unfold is an unofficial [Fold Money](https://fold.money) CLI client, which
covers the bare minimum API routes to fetch your transactions for a given
period and even write them to a sqlite database.

Fold's API is not publically available, I had to MITM their app to write this
tool, and so **there might be unforeseen consequences for your Fold account if
you use this tool**.

### Prerequisites

- You need a Fold Account, which is currently only available on an invite basis
- You need to connect to whichever banks you have using the Fold app first

### Usage

**Caution: For all I know, Fold has not publically released a client which can
support multiple sessions, which means when you use this CLI, you will be
automatically logged out on your Phone's app**

1. First, login to your account:
    ```bash
    $ unfold login
    ```

2. Then you can fetch your transactions:  

    a. In plaintext:
      ```bash
      $ unfold transactions
      ```

    b. In plaintext and also write to a db:
      ```bash
      # Write to a local file called `db.sqlite` by default
      $ unfold transactions -s 2023-09-20 --db
      ```

    c. For a complete glossary of available options:
      ```
      $ unfold transactions -h
      Prints the transactions from all of your accounts (default period: 1 month)

      Usage:
        unfold transactions [flags]

      Flags:
        -d, --db               Save the results in a sqlite db
        -D, --db-path string   Sets path for the database (default "db.sqlite")
        -h, --help             help for transactions
        -s, --since string     fetch transactions since in this format: YYYY-MM-DD (default "2023-09-15")
        -t, --till string      fetch transactions till in this format: YYYY-MM-DD (default "2023-10-16")

      Global Flags:
            --config string   config file (default is $HOME/.config/unfold/config.yaml)
        -v, --debug           Enable debug mode
      ```

There are a few more subcommands which Unfold provides and uses internally. You
can get a list by:
```
$ unfold
An unofficial cli client for fold.money

Usage:
  unfold [command]

Available Commands:
  availability Returns a range of dates for when your banking data is available
  completion   Generate the autocompletion script for the specified shell
  help         Help about any command
  login        Log in to your fold account
  refresh      Refresh your auth tokens
  transactions Prints the transactions from all of your accounts in the last x days
  user         Get your account details

Flags:
      --config string   config file (default is $HOME/.config/unfold/config.yaml)
  -h, --help            help for unfold

Use "unfold [command] --help" for more information about a command.
```

### Credits

[Fold Money](https://fold.money), for their Account Aggregator integration
