package scrapper

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"leetcode-scrapper/config"
	"net/http"
	"time"
)

// GraphQL request structures
type GraphQLRequest struct {
	Query         string      `json:"query"`
	Variables     interface{} `json:"variables"`
	OperationName string      `json:"operationName"`
}

// Problem list structures
type FavoriteQuestionListVariables struct {
	Skip          int       `json:"skip"`
	Limit         int       `json:"limit"`
	FavoriteSlug  string    `json:"favoriteSlug"`
	FiltersV2     FiltersV2 `json:"filtersV2"`
	SearchKeyword string    `json:"searchKeyword"`
	SortBy        SortBy    `json:"sortBy"`
}

type FiltersV2 struct {
	FilterCombineType   string                 `json:"filterCombineType"`
	StatusFilter        StatusFilter           `json:"statusFilter"`
	DifficultyFilter    DifficultyFilter       `json:"difficultyFilter"`
	LanguageFilter      LanguageFilter         `json:"languageFilter"`
	TopicFilter         TopicFilter            `json:"topicFilter"`
	AcceptanceFilter    map[string]interface{} `json:"acceptanceFilter"`
	FrequencyFilter     map[string]interface{} `json:"frequencyFilter"`
	FrontendIdFilter    map[string]interface{} `json:"frontendIdFilter"`
	LastSubmittedFilter map[string]interface{} `json:"lastSubmittedFilter"`
	PublishedFilter     map[string]interface{} `json:"publishedFilter"`
	CompanyFilter       CompanyFilter          `json:"companyFilter"`
	PositionFilter      PositionFilter         `json:"positionFilter"`
	PremiumFilter       PremiumFilter          `json:"premiumFilter"`
}

type StatusFilter struct {
	QuestionStatuses []string `json:"questionStatuses"`
	Operator         string   `json:"operator"`
}

type DifficultyFilter struct {
	Difficulties []string `json:"difficulties"`
	Operator     string   `json:"operator"`
}

type LanguageFilter struct {
	LanguageSlugs []string `json:"languageSlugs"`
	Operator      string   `json:"operator"`
}

type TopicFilter struct {
	TopicSlugs []string `json:"topicSlugs"`
	Operator   string   `json:"operator"`
}

type CompanyFilter struct {
	CompanySlugs []string `json:"companySlugs"`
	Operator     string   `json:"operator"`
}

type PositionFilter struct {
	PositionSlugs []string `json:"positionSlugs"`
	Operator      string   `json:"operator"`
}

type PremiumFilter struct {
	PremiumStatus []string `json:"premiumStatus"`
	Operator      string   `json:"operator"`
}

type SortBy struct {
	SortField string `json:"sortField"`
	SortOrder string `json:"sortOrder"`
}

// Response structures
type FavoriteQuestionListResponse struct {
	Data struct {
		FavoriteQuestionList struct {
			Questions   []Question `json:"questions"`
			TotalLength int        `json:"totalLength"`
			HasMore     bool       `json:"hasMore"`
		} `json:"favoriteQuestionList"`
	} `json:"data"`
}

type Question struct {
	Difficulty         string     `json:"difficulty"`
	ID                 int        `json:"id"`
	PaidOnly           bool       `json:"paidOnly"`
	QuestionFrontendID string     `json:"questionFrontendId"`
	Status             string     `json:"status"`
	Title              string     `json:"title"`
	TitleSlug          string     `json:"titleSlug"`
	TranslatedTitle    string     `json:"translatedTitle"`
	IsInMyFavorites    bool       `json:"isInMyFavorites"`
	Frequency          float64    `json:"frequency"`
	AcRate             float64    `json:"acRate"`
	TopicTags          []TopicTag `json:"topicTags"`
}

type TopicTag struct {
	Name string `json:"name"`
	Slug string `json:"slug"`
}

// Problem detail structures
type ProblemDetailResponse struct {
	Data struct {
		Question struct {
			QuestionID         string     `json:"questionId"`
			QuestionFrontendID string     `json:"questionFrontendId"`
			Title              string     `json:"title"`
			TitleSlug          string     `json:"titleSlug"`
			IsPaidOnly         bool       `json:"isPaidOnly"`
			Difficulty         string     `json:"difficulty"`
			SimilarQuestions   string     `json:"similarQuestions"`
			ExampleTestcases   string     `json:"exampleTestcases"`
			TopicTags          []TopicTag `json:"topicTags"`
			CompanyTagStats    string     `json:"companyTagStats"`
			Stats              string     `json:"stats"`
		} `json:"question"`
	} `json:"data"`
}

// LeetCodeScraper handles the scraping operations
type LeetCodeScraper struct {
	client  *http.Client
	baseURL string
	headers map[string]string
}

// NewLeetCodeScraper creates a new scraper instance
func NewLeetCodeScraper() *LeetCodeScraper {
	return &LeetCodeScraper{
		client:  &http.Client{Timeout: 30 * time.Second},
		baseURL: "https://leetcode.com/graphql/",
		headers: map[string]string{
			"Content-Type": "application/json",
			"Referer":      "https://leetcode.com/",
			"Cookie":       config.GetApp("config.yaml").Headers.Cookie,
		},
	}
}

// makeRequest makes a GraphQL request to LeetCode
func (s *LeetCodeScraper) makeRequest(query string, variables interface{}, operationName string) ([]byte, error) {
	reqBody := GraphQLRequest{
		Query:         query,
		Variables:     variables,
		OperationName: operationName,
	}

	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	req, err := http.NewRequest("POST", s.baseURL, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	for key, value := range s.headers {
		req.Header.Set(key, value)
	}

	resp, err := s.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to make request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("received non-200 status code: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	return body, nil
}

func buildQueryAndVariables(favoriteSlug string, skip, limit int) (string, interface{}) {
	query := `
	query favoriteQuestionList($favoriteSlug: String!, $filter: FavoriteQuestionFilterInput, $filtersV2: QuestionFilterInput, $searchKeyword: String, $sortBy: QuestionSortByInput, $limit: Int, $skip: Int, $version: String = "v2") {
		favoriteQuestionList(
			favoriteSlug: $favoriteSlug
			filter: $filter
			filtersV2: $filtersV2
			searchKeyword: $searchKeyword
			sortBy: $sortBy
			limit: $limit
			skip: $skip
			version: $version
		) {
			questions {
				difficulty
				id
				paidOnly
				questionFrontendId
				status
				title
				titleSlug
				translatedTitle
				isInMyFavorites
				frequency
				acRate
				topicTags {
					name
					nameTranslated
					slug
				}
			}
			totalLength
			hasMore
		}
	}`

	variables := FavoriteQuestionListVariables{
		Skip:         skip,
		Limit:        limit,
		FavoriteSlug: favoriteSlug,
		FiltersV2: FiltersV2{
			FilterCombineType:   "ALL",
			StatusFilter:        StatusFilter{QuestionStatuses: []string{}, Operator: "IS"},
			DifficultyFilter:    DifficultyFilter{Difficulties: []string{}, Operator: "IS"},
			LanguageFilter:      LanguageFilter{LanguageSlugs: []string{}, Operator: "IS"},
			TopicFilter:         TopicFilter{TopicSlugs: []string{}, Operator: "IS"},
			AcceptanceFilter:    map[string]interface{}{},
			FrequencyFilter:     map[string]interface{}{},
			FrontendIdFilter:    map[string]interface{}{},
			LastSubmittedFilter: map[string]interface{}{},
			PublishedFilter:     map[string]interface{}{},
			CompanyFilter:       CompanyFilter{CompanySlugs: []string{}, Operator: "IS"},
			PositionFilter:      PositionFilter{PositionSlugs: []string{}, Operator: "IS"},
			PremiumFilter:       PremiumFilter{PremiumStatus: []string{}, Operator: "IS"},
		},
		SearchKeyword: "",
		SortBy: SortBy{
			SortField: "CUSTOM",
			SortOrder: "ASCENDING",
		},
	}

	return query, variables
}

// GetFavoriteQuestionList fetches questions from a favorite list
func (s *LeetCodeScraper) GetFavoriteQuestionList(favoriteSlug string, skip, limit int) (*FavoriteQuestionListResponse, error) {

	query, variables := buildQueryAndVariables(favoriteSlug, skip, limit)
	body, err := s.makeRequest(query, variables, "favoriteQuestionList")
	if err != nil {
		return nil, err
	}

	var response FavoriteQuestionListResponse
	if err := json.Unmarshal(body, &response); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	return &response, nil
}

// GetAllProblems fetches all problems from LeetCode
func (s *LeetCodeScraper) GetAllProblems(skip, limit int) ([]Question, error) {
	query := `
	query problemsetQuestionList($categorySlug: String, $limit: Int, $skip: Int, $filters: QuestionListFilterInput) {
		problemsetQuestionList: questionList(
			categorySlug: $categorySlug
			limit: $limit
			skip: $skip
			filters: $filters
		) {
			questions: data {
				acRate
				difficulty
				freqBar
				frontendQuestionId: questionFrontendId
				isFavor
				paidOnly: isPaidOnly
				status
				title
				titleSlug
				topicTags {
					name
					id
					slug
				}
			}
		}
	}`

	variables := map[string]interface{}{
		"categorySlug": "",
		"skip":         skip,
		"limit":        limit,
		"filters":      map[string]interface{}{},
	}

	body, err := s.makeRequest(query, variables, "problemsetQuestionList")
	if err != nil {
		return nil, err
	}

	var response struct {
		Data struct {
			ProblemsetQuestionList struct {
				Questions []struct {
					AcRate             float64    `json:"acRate"`
					Difficulty         string     `json:"difficulty"`
					FreqBar            float64    `json:"freqBar"`
					FrontendQuestionID string     `json:"frontendQuestionId"`
					IsFavor            bool       `json:"isFavor"`
					PaidOnly           bool       `json:"paidOnly"`
					Status             string     `json:"status"`
					Title              string     `json:"title"`
					TitleSlug          string     `json:"titleSlug"`
					TopicTags          []TopicTag `json:"topicTags"`
				} `json:"questions"`
			} `json:"problemsetQuestionList"`
		} `json:"data"`
	}

	if err := json.Unmarshal(body, &response); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	// Convert to our Question structure
	var questions []Question
	for _, q := range response.Data.ProblemsetQuestionList.Questions {
		questions = append(questions, Question{
			Difficulty:         q.Difficulty,
			PaidOnly:           q.PaidOnly,
			QuestionFrontendID: q.FrontendQuestionID,
			Status:             q.Status,
			Title:              q.Title,
			TitleSlug:          q.TitleSlug,
			IsInMyFavorites:    q.IsFavor,
			Frequency:          q.FreqBar,
			AcRate:             q.AcRate,
			TopicTags:          q.TopicTags,
		})
	}

	return questions, nil
}

// GetProblemDetail fetches detailed information for a specific problem
func (s *LeetCodeScraper) GetProblemDetail(titleSlug string) (*ProblemDetailResponse, error) {
	query := `
	query questionData($titleSlug: String!) {
		question(titleSlug: $titleSlug) {
			questionId
			questionFrontendId
			title
			titleSlug
			translatedTitle
			isPaidOnly
			difficulty
			similarQuestions
			exampleTestcases
			contributors {
				username
				profileUrl
				avatarUrl
				__typename
			}
			topicTags {
				slug
				translatedName
				__typename
			}
			companyTagStats
			stats
		}
	}`

	variables := map[string]interface{}{
		"titleSlug": titleSlug,
	}

	body, err := s.makeRequest(query, variables, "questionData")
	if err != nil {
		return nil, err
	}

	var response ProblemDetailResponse
	if err := json.Unmarshal(body, &response); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	return &response, nil
}
