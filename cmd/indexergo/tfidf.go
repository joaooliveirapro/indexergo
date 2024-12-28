package indexergo

import "math"

type DF map[string]int
type IDF map[string]float64

func CalculateDF(docs *[]PageInfo, searchTokens []string) DF {
	// Step 1: Calculate TF for each query word
	// In how many document does each token appear?
	df := DF{}
	// Initialise DF map with each word
	for _, word := range searchTokens {
		df[word] = 0
	}
	for _, doc := range *docs {
		for _, word := range searchTokens {
			_, exists := doc.ContentTokens[word]
			if exists {
				df[word] += 1
			}
		}
	}
	return df
}

func CalculateIDF(docs *[]PageInfo, searchTokens []string, df DF) IDF {
	// Step 2: Calculate IDF for each query word
	corpus_len := len(*docs) // Total documents
	idf := IDF{}
	for _, word := range searchTokens {
		idf[word] = math.Log(float64(corpus_len) / float64(df[word]))
	}
	return idf
}

func CalculateTFIDF(docs *[]PageInfo, searchTokens []string) []Document {
	df := CalculateDF(docs, searchTokens)
	idf := CalculateIDF(docs, searchTokens, df)
	tfidfResults := []Document{}
	for _, page := range *docs {
		doc := Document{URL: page.HTTPResponse.URL, Ranking: 0, QueryWeight: map[string]float64{}}
		total_words := 0
		for _, v := range page.ContentTokens {
			total_words += v
		}
		tfidfScore := 0.0
		for _, word := range searchTokens {
			tf := float64(page.ContentTokens[word]) / float64(total_words)
			doc.QueryWeight[word] = tf
			tfidfScore += tf * idf[word]
		}
		doc.Ranking = tfidfScore
		tfidfResults = append(tfidfResults, doc)
	}
	return tfidfResults
}
