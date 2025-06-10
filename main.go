package main

import (
	"fmt"
	"leetcode-scrapper/scrapper"
	"leetcode-scrapper/utils"
)

func downloadCompanyProblems(favoriteSlug string) {
	// Example 1: Scrape problems from a favorite list
	fmt.Println("Scraping Facebook 30 Days favorite list...")
	favoriteResponse, err := leetcodeScrapper.GetFavoriteQuestionList(favoriteSlug, 0, 10)
	if err != nil {
		fmt.Printf("Error scraping favorite list: %v\n", err)
		return
	}

	fmt.Printf("Found %d problems in favorite list\n", len(favoriteResponse.Data.FavoriteQuestionList.Questions))

	// Save favorite list problems
	if err := utils.SaveToFile(favoriteResponse, "data/facebook_30_days.json"); err != nil {
		fmt.Printf("Error saving favorite list: %v\n", err)
		return
	}
}

var leetcodeScrapper *scrapper.LeetCodeScraper

func init() {
	leetcodeScrapper = scrapper.NewLeetCodeScraper()
}

func main() {
	slugLists := []string{
		"facebook-thirty-days",
		"facebook-three-months",
		"facebook-six-months",
	}
	downloadCompanyProblems(slugLists[0])
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

	fmt.Println("Scraping completed! Check the 'data' directory for results.")
}
