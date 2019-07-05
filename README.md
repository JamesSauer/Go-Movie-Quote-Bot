# Movie Quote Bot
* Finds all movie entries on [WikiQuote.org](https://en.wikiquote.org)
* Scrapes quotes off a random page
* Picks one at random and prints it to the console

# Usage
```
$ go get github.com/JamesSauer/Go-Movie-Quote-Bot
$ go install github.com/JamesSauer/Go-Movie-Quote-Bot
$ mqbot
> That's part of your problem: you haven't seen enough movies. All of life's riddles are answered in the movies.
>     - Davis, Grand Canyon (1991 film)
```

By default, mqbot attempts to retrieve a quote from the database before scraping Wikiquote.
It checks the MQBOT_POSTGRES environment variable for a connection string.

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

To test the database connection:
```
$ mqbot testdb
> Successfully connected to DB!
```

To scrape a random movie entry and persist all its quotes to the database:
```
$ mqbot scrape1page
> Successfully scraped and saved the entry for the movie "Angels with Dirty Faces"!
```

To scrape ALL movie entries:
**Note: The bot currently doesn't check whether or not a quote already exists within the database.**
**Using this command more than once will result in duplicate database entries.**
```
$ mqbot scrapeall
> This command will attempt to scrape the entirety of wikiquote.org's movie quotes.
> This will take more than 10 minutes.
>
> Do you want to proceed? (yes/y/no/n)
$ yes
> Scraped 2378 pages in 12m10.7716143s!
```