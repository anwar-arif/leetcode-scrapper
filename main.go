package main

import (
	"fmt"
	"leetcode-scrapper/scrapper"
	"leetcode-scrapper/utils"
	"time"
)

var chunkSize = 20

func downloadCompanyProblems(favoriteSlug string) {
	// Example 1: Scrape problems from a favorite list
	fmt.Println(fmt.Sprintf("Scraping %s list...", favoriteSlug))
	favoriteResponse, err := leetcodeScrapper.GetFavoriteQuestionList(favoriteSlug, 0, chunkSize)
	if err != nil {
		fmt.Printf("Error scraping favorite list: %v\n", err)
		return
	}

	totalQuestions := favoriteResponse.Data.FavoriteQuestionList.TotalLength

	var iteration = 1
	for len(favoriteResponse.Data.FavoriteQuestionList.Questions) < totalQuestions {
		res, resErr := leetcodeScrapper.GetFavoriteQuestionList(favoriteSlug, iteration*chunkSize, chunkSize)
		if resErr != nil {
			fmt.Printf("Error scraping favorite list: %v\n", resErr)
		}
		favoriteResponse.Data.FavoriteQuestionList.Questions = append(favoriteResponse.Data.FavoriteQuestionList.Questions, res.Data.FavoriteQuestionList.Questions...)
		iteration++
		fmt.Println(fmt.Sprintf("%s: fetched %d questions out of %d questions", favoriteSlug, len(favoriteResponse.Data.FavoriteQuestionList.Questions), totalQuestions))
		time.Sleep(1 * time.Second)
	}

	fmt.Printf("Found %d problems in favorite list\n", len(favoriteResponse.Data.FavoriteQuestionList.Questions))

	// Save favorite list problems
	if err := utils.SaveToFile(favoriteResponse, fmt.Sprintf("data/%s.json", favoriteSlug)); err != nil {
		fmt.Printf("Error saving favorite list: %v\n", err)
		return
	}
}

var leetcodeScrapper *scrapper.LeetCodeScraper

func init() {
	chunkSize = 10
	leetcodeScrapper = scrapper.NewLeetCodeScraper()
}

func main() {
	leetcodeScrapper.GetRandomQuestion()
	//slugLists := []string{
	// facebook
	//"facebook-thirty-days",
	//"facebook-three-months",
	//"facebook-six-months",
	// google
	//"google-thirty-days",
	//"google-three-months",
	//"google-six-months",
	// amazon
	//"amazon-thirty-days",
	//"amazon-three-months",
	//"amazon-six-months",
	// microsoft
	//"microsoft-thirty-days",
	//"microsoft-three-months",
	//"microsoft-six-months",
	// uber
	//"uber-thirty-days",
	//"uber-three-months",
	//"uber-six-months",
	// apple
	//"apple-thirty-days",
	//"apple-three-months",
	//"apple-six-months",
	// netflix
	//"netflix-thirty-days",
	//"netflix-three-months",
	//"netflix-six-months",
	// bloomberg
	//"bloomberg-thirty-days",
	//"bloomberg-three-months",
	//"bloomberg-six-months",
	// tiktok
	//"tiktok-thirty-days",
	//"tiktok-three-months",
	//"tiktok-six-months",
	//}
	// downloadCompanyProblems(slugLists[0])
	//for _, slug := range slugLists {
	//	downloadCompanyProblems(slug)
	//	fmt.Printf("Scraping %s completed\n", slug)
	//	time.Sleep(2 * time.Second)
	//}

	//// Example 2: Scrape all problems (first batch)
	//fmt.Println("Scraping all problems (first 50)...")
	//allProblems, err := scraper.GetAllProblems(0, 50)
	//if err != nil {
	//	fmt.Printf("Error scraping all problems: %v\n", err)
	//	return
	//}
	//
	//fmt.Printf("Found %d problems\n", len(allProblems))
	//
	//// Save all problems
	//if err := utils.SaveToFile(allProblems, "data/all_problems_batch_1.json"); err != nil {
	//	fmt.Printf("Error saving all problems: %v\n", err)
	//	return
	//}
	//
	//// Example 3: Get detailed information for specific problems
	//fmt.Println("Getting detailed information for first few problems...")
	//detailsDir := "data/problem_details"
	//
	//for i, problem := range allProblems[:5] { // Get details for first 5 problems
	//	if problem.PaidOnly {
	//		fmt.Printf("Skipping paid problem: %s\n", problem.Title)
	//		continue
	//	}
	//
	//	fmt.Printf("Getting details for: %s\n", problem.Title)
	//	detail, err := scraper.GetProblemDetail(problem.TitleSlug)
	//	if err != nil {
	//		fmt.Printf("Error getting details for %s: %v\n", problem.Title, err)
	//		continue
	//	}
	//
	//	filename := fmt.Sprintf("%s/%d_%s.json", detailsDir, i+1,
	//		strings.ReplaceAll(problem.TitleSlug, "-", "_"))
	//
	//	if err := utils.SaveToFile(detail, filename); err != nil {
	//		fmt.Printf("Error saving details for %s: %v\n", problem.Title, err)
	//		continue
	//	}
	//
	//	// Be respectful to the server
	//	time.Sleep(1 * time.Second)
	//}

	//fmt.Println("Scraping completed! Check the 'data' directory for results.")
}
