# Movie Quote Bot
* Finds all movie entries on [WikiQuote.org](https://en.wikiquote.org)
* Scrapes quotes off the pages
* Optionally persists quotes in a Postgres database
* Picks one at random and prints it to the console

# Usage
```
$ go get github.com/JamesSauer/Go-Movie-Quote-Bot
$ cd $GOPATH/src/github.com/JamesSauer/Go-Movie-Quote-Bot
$ go install
$ mqbot
> That's part of your problem: you haven't seen enough movies. All of life's riddles are answered in the movies.
>     - Davis, Grand Canyon (1991 film)
```

By default, mqbot attempts to retrieve a quote from the database before scraping Wikiquote.   
It checks the MQBOT_POSTGRES environment variable for a connection string.

To test the database connection:
```
$ mqbot testdb
> Successfully connected to DB!
```

To set up the schema:
```
$ mqbot initdb
> Successfully set up database schema.
```

To force scraping a fresh quote, use the --fresh or -f flag:
```
$ mqbot -f
> ...
```

To avoid using scraping as a fallback, use the --database or -db flag:
```
$ mqbot -db
> ...
```

To scrape a random movie entry and persist all its quotes to the database:
```
$ mqbot scrape1
> Successfully scraped and saved the entry for the movie "Angels with Dirty Faces"!
```

To scrape ALL movie entries:
```
$ mqbot scrapeall
> This command will attempt to scrape the entirety of wikiquote.org's movie quotes.
> This might take more than 10 minutes. Do you want to proceed? (yes/y/no/n)
$ yes
> Scraped 2378 pages in 12m10.7716143s!
```