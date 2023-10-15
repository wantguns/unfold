# Unfold

Unfold is an unofficial [Fold Money](https://fold.money) CLI client, which
covers the bare minimum API routes to fetch your transactions for a given
period and even write them to a sqlite database. It also provides options for
running an internal cron job through which this CLI acts as a daemon and
fetches transactions every time the cron job's timer is met.

Fold's API is not publically available, I had to MITM their app to write this
tool, and so **there might be unforeseen consequences for your Fold account if
you use this tool**.

### Prerequisites

- You need a Fold Account, which is currently only available on an invite basis
- You need to connect to whichever banks you have using the Fold app first

### Installation

- Using golang's build system:
  ```bash
  $ go install github.com/wantguns/unfold
  ```

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

    c. Create an internal cron job to fetch transactions every 20 seconds and save them to a db: 
      ```bash
      # Note: You need to enable the `-d` or `--db` flag to ensure that the changes are written to a database
      $ unfold transactions -s 2023-09-20 --db -w '@every 20s'
      12:19AM INF Cron job set for fetching transactions, going into daemon mode
      12:19AM INF Fetched transactions till 2023-10-17
      12:20AM INF Fetched transactions till 2023-10-16
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
        -s, --since string     fetch transactions since in this format: YYYY-MM-DD (default "XXXX-XX-XX")
        -t, --till string      fetch transactions till in this format: YYYY-MM-DD (default "XXXX-XX-XX")
        -w, --watch string     Set an internal cron job to trigger this command. You can use non-standard cron expressions like '@every 6h'. This will disable plaintext mode, so add a '-d' flag if you want to write to db

      Global Flags:
            --config string   config file (default is $HOME/.config/unfold/config.yaml)
        -v, --debug           Enable debug mode
            Prints the transactions from all of your accounts (default period: 1 month)
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
  transactions Prints the transactions from all of your accounts (default period: 1 month)
  user         Get your account details

Flags:
      --config string   config file (default is $HOME/.config/unfold/config.yaml)
  -v, --debug           Enable debug mode  
  -h, --help            help for unfold

Use "unfold [command] --help" for more information about a command.
```

### Credits

[Fold Money](https://fold.money), for their Account Aggregator integration
