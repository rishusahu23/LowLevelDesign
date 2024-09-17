package main

import (
	"errors"
	"fmt"
)

// Review represents a review made by a user for a product.
type Review struct {
	ID        string
	UserID    string
	ProductID string
	Rating    int
	Content   string
}

// ReviewService interface defines the operations for managing reviews.
type ReviewService interface {
	CreateReview(userID, productID string, rating int, content string) error
	UpdateReview(reviewID string, rating int, content string) error
	DeleteReview(reviewID string) error
	FetchReviews(productID string) ([]Review, error)
	FetchReview(reviewID string) (Review, error)
}

// ReviewStorage interface defines the operations for storing reviews.
type ReviewStorage interface {
	SaveReview(review Review) error
	UpdateReview(review Review) error
	DeleteReview(reviewID string) error
	GetReviewByID(reviewID string) (Review, error)
	GetReviewsByProductID(productID string) ([]Review, error)
}

// ReviewModeration interface defines the operations for moderating reviews.
type ReviewModeration interface {
	ApproveReview(reviewID string) error
	RejectReview(reviewID string, reason string) error
	FlagReview(reviewID string, reason string) error
}

// AnalyticsService interface defines the operations for analyzing reviews.
type AnalyticsService interface {
	CalculateAverageRating(productID string) (float64, error)
	CountReviews(productID string) (int, error)
}

// DefaultReviewService is a basic implementation of ReviewService.
type DefaultReviewService struct {
	storage ReviewStorage
}

func (s *DefaultReviewService) CreateReview(userID, productID string, rating int, content string) error {
	review := Review{
		ID:        fmt.Sprintf("%s-%s", userID, productID),
		UserID:    userID,
		ProductID: productID,
		Rating:    rating,
		Content:   content,
	}
	return s.storage.SaveReview(review)
}

func (s *DefaultReviewService) UpdateReview(reviewID string, rating int, content string) error {
	review, err := s.storage.GetReviewByID(reviewID)
	if err != nil {
		return err
	}
	review.Rating = rating
	review.Content = content
	return s.storage.UpdateReview(review)
}

func (s *DefaultReviewService) DeleteReview(reviewID string) error {
	return s.storage.DeleteReview(reviewID)
}

func (s *DefaultReviewService) FetchReviews(productID string) ([]Review, error) {
	return s.storage.GetReviewsByProductID(productID)
}

func (s *DefaultReviewService) FetchReview(reviewID string) (Review, error) {
	return s.storage.GetReviewByID(reviewID)
}

// InMemoryReviewStorage is a simple in-memory implementation of ReviewStorage.
type InMemoryReviewStorage struct {
	reviews map[string]Review
}

func NewInMemoryReviewStorage() *InMemoryReviewStorage {
	return &InMemoryReviewStorage{reviews: make(map[string]Review)}
}

func (s *InMemoryReviewStorage) SaveReview(review Review) error {
	s.reviews[review.ID] = review
	return nil
}

func (s *InMemoryReviewStorage) UpdateReview(review Review) error {
	if _, exists := s.reviews[review.ID]; exists {
		s.reviews[review.ID] = review
		return nil
	}
	return errors.New("review not found")
}

func (s *InMemoryReviewStorage) DeleteReview(reviewID string) error {
	delete(s.reviews, reviewID)
	return nil
}

func (s *InMemoryReviewStorage) GetReviewByID(reviewID string) (Review, error) {
	if review, exists := s.reviews[reviewID]; exists {
		return review, nil
	}
	return Review{}, errors.New("review not found")
}

func (s *InMemoryReviewStorage) GetReviewsByProductID(productID string) ([]Review, error) {
	var result []Review
	for _, review := range s.reviews {
		if review.ProductID == productID {
			result = append(result, review)
		}
	}
	return result, nil
}

// DefaultReviewModeration is a basic implementation of ReviewModeration.
type DefaultReviewModeration struct{}

func (m *DefaultReviewModeration) ApproveReview(reviewID string) error {
	fmt.Printf("Review %s approved\n", reviewID)
	return nil
}

func (m *DefaultReviewModeration) RejectReview(reviewID string, reason string) error {
	fmt.Printf("Review %s rejected: %s\n", reviewID, reason)
	return nil
}

func (m *DefaultReviewModeration) FlagReview(reviewID string, reason string) error {
	fmt.Printf("Review %s flagged: %s\n", reviewID, reason)
	return nil
}

// DefaultAnalyticsService is a basic implementation of AnalyticsService.
type DefaultAnalyticsService struct{}

func (a *DefaultAnalyticsService) CalculateAverageRating(productID string) (float64, error) {
	// Mocking the calculation
	return 4.5, nil
}

func (a *DefaultAnalyticsService) CountReviews(productID string) (int, error) {
	// Mocking the count
	return 42, nil
}

// ModeratedReviewService is a decorator that adds moderation functionality to ReviewService.
type ModeratedReviewService struct {
	ReviewService
	moderation ReviewModeration
}

func NewModeratedReviewService(baseService ReviewService, moderation ReviewModeration) *ModeratedReviewService {
	return &ModeratedReviewService{
		ReviewService: baseService,
		moderation:    moderation,
	}
}

func (m *ModeratedReviewService) CreateReview(userID, productID string, rating int, content string) error {
	err := m.ReviewService.CreateReview(userID, productID, rating, content)
	if err != nil {
		return err
	}
	// Automatically flag reviews for moderation if certain conditions are met
	reviews, err := m.ReviewService.FetchReviews(productID)
	if err == nil && len(reviews) > 0 {
		latestReview := reviews[len(reviews)-1] // Get the latest review
		if latestReview.Rating < 3 {            // Example condition for moderation
			return m.moderation.FlagReview(latestReview.ID, "Low rating")
		}
	}
	return nil
}

// AnalyticsReviewService is a decorator that adds analytics functionality to ReviewService.
type AnalyticsReviewService struct {
	ReviewService
	analytics AnalyticsService
}

func NewAnalyticsReviewService(baseService ReviewService, analytics AnalyticsService) *AnalyticsReviewService {
	return &AnalyticsReviewService{
		ReviewService: baseService,
		analytics:     analytics,
	}
}

func (a *AnalyticsReviewService) FetchReviews(productID string) ([]Review, error) {
	reviews, err := a.ReviewService.FetchReviews(productID)
	if err != nil {
		return nil, err
	}
	// Log or update analytics after fetching reviews
	avgRating, err := a.analytics.CalculateAverageRating(productID)
	if err == nil {
		fmt.Printf("Average rating for product %s is %.2f\n", productID, avgRating)
	}
	return reviews, nil
}

// Main function to demonstrate the system
func main() {
	// Create the components using the factory pattern
	storage := NewInMemoryReviewStorage()
	moderation := &DefaultReviewModeration{}
	analytics := &DefaultAnalyticsService{}

	baseService := &DefaultReviewService{
		storage: storage,
	}

	// Decorate the base service with moderation and analytics
	moderatedService := NewModeratedReviewService(baseService, moderation)
	analyticsService := NewAnalyticsReviewService(moderatedService, analytics)

	// Create a review
	err := analyticsService.CreateReview("user1", "product1", 2, "Not as expected!")
	if err != nil {
		fmt.Println("Error creating review:", err)
		return
	}

	// Fetch reviews for a product
	reviews, err := analyticsService.FetchReviews("product1")
	if err != nil {
		fmt.Println("Error fetching reviews:", err)
		return
	}

	fmt.Println("Reviews for product1:", reviews)
}
