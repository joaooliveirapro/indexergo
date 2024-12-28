# Indexer 🔎
Provided a set of HTML documents this app will compile a frequency count of:
- HTML tags in the HTML page 
- Text content tokens* 
- Allow for page ranking by query search

## Features
- Content search
- [TF-IDF Page ranking](https://wikipedia.org/wiki/Tf%E2%80%93idf) (Wikipedia link)
- Page information includes:
    - Response status code
    - Redirect history
    - HTML tags frequency
    - Text content tokens frequency
- Indexing is cached (in index.json) for performance


## Install
```sh
$ go get github.com/joaooliveirapro/indexergo # install
$ go mod tidy                                 # clean up dependencies
```

## How to use
```go
    ig := indexergo.Indexer{
	URLsFilePath: "",                                  // Provide a path to a .txt file
        URLsList: []string{"https://mysite.com"},          // OR list the URLs individually
        LookByQuerySelector: []string{".ats-description"}, // Optional (recommended for better results)
    }

    err := ig.IndexDocuments()
    if err != nil {
        log.Fatal(err.Error())
    }

    docs, err := ig.Search("some keywords")
    if err != nil {
        log.Fatal(err.Error())
    }

    for i, doc := range docs {
        fmt.Printf("%d - %s - Rank: %f", i, doc.URL, doc.Ranking)
    }
```
<div style="background-color: #fff3cd; border: 1px solid #ffecb5; padding: 10px;">
<strong>URLsFilePath</strong> - file must contain one URL per line. No comma at end of line.

<strong>*Tokens</strong> - are individual words. Punctuation is removed and all tokens are lowercase.
</div>

### Results
```sh
# Query: "some keywords"
1 - Doc_1 - Rank: 1.23
2 - Doc_2 - Rank: 0.98
...
```

### Cached index
```json
// index.json
[
  {
    "httpResponse": {
      "statusCode": 200,
      "url": "https://careers.adeccogroup.com/en/job/-/-/22630/72523943584",
      "redirected": false,
      "redirectsHistory": null
    },
    "htmlTags": {
      "a": 69,
      "body": 1,
      "br": 54,
      "button": 16,
      "div": 81,
      "form": 4,
      "h1": 1,
      "label": 10,
      "legend": 1,
      "li": 59,
      ...
    },
    "contentTokens": {
      "additional": 1,
      "adecco": 5,
      "advanced": 1,
      "alignment": 1,
      "all": 2,
      ...
    },
    "timestamp": "27-12-2024 21:49:14"
  }
]

```

### License
The MIT License (MIT)