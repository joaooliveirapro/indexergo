package main

func main() {
	/* Example */

	// ig := indexergo.Indexer{
	// 	URLsFilePath: "", // Provide a path to a .txt file (1 URL per line)
	// 	URLsList: []string{
	// 		"https://careers.adeccogroup.com/en/job/sofia/senior-finance-data-analyst/22630/72523943584",
	// 		"https://careers.adeccogroup.com/en/job/lucerne/sales-consultant-medical-luzern-80-100-w-m-d/37746/69337371728",
	// 		"https://careers.adeccogroup.com/en/job/prague/it-security-architect/37746/66582894672",
	// 		"https://careers.adeccogroup.com/en/job/dusseldorf/head-of-rewards-germany-m-w-d/37746/74720226000",
	// 	}, // OR list the URLs individually
	// 	LookByQuerySelector: []string{".ats-description"}, // Optional (recommended for better results)
	// }

	// err := ig.IndexDocuments()
	// if err != nil {
	// 	log.Fatal(err.Error())
	// }

	// docs, err := ig.Search("adecco finance")
	// if err != nil {
	// 	log.Fatal(err.Error())
	// }

	// for i, doc := range docs {
	// 	fmt.Printf("%d - [%.5f] - %s\n %+v\n", i, doc.Ranking, doc.URL, doc.QueryWeight)
	// }
}
