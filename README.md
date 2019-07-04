# Movie Quote Bot
* Crawls [WikiQuote.org](https://en.wikiquote.org) for a random movie entry
* Scrapes all quotes off that page
* Picks one at random and prints it to the console

```
$ go get github.com/JamesSauer/Go-Movie-Quote-Bot
$ go install github.com/JamesSauer/Go-Movie-Quote-Bot
$ mqbot
That's part of your problem: you haven't seen enough movies. All of life's riddles are answered in the movies.
    - Davis, Grand Canyon (1991 film)
```

By default, mqbot attempts to retrieve a quote from the database before scraping Wikiquote.
To force scraping a fresh quote, use the --fresh or -f flag:
```
$mqbot -f
```

To avoid using scraping as a fallback, use the --database or -db flag:
```
$mqbot -db
```

To test the database connection:
```
$mqbot testdb
```