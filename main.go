// The main package is the entry point for the application
package main

// Import necessary libraries
import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"time"
)

// APIClient holds the HTTP client and authentication tokens
type APIClient struct {
	client    *http.Client
	authToken string
	userID    string
	baseURL   string
	outputDir string
}

// Activity represents a MyWhoosh activity
type Activity struct {
	ID             string    `json:"id"`
	Title          string    `json:"title"`
	ActivityFileID string    `json:"activityFileId"`
	Sport          string    `json:"sportType"`
	StartTime      time.Time `json:"startDatetime"`
	Distance       float64   `json:"distance"`
	TotalWorkout   int       `json:"totalWorkout"`
	Heartrate      int       `json:"heartrate"`
	RideDuration   string    `json:"rideDuration"`
}

// ActivitiesResponse represents the API response from the activities endpoint
type ActivitiesResponse struct {
	Error   bool   `json:"error"`
	Code    int    `json:"code"`
	Message string `json:"message"`
	Display bool   `json:"display"`
	Data    struct {
		Results []Activity `json:"results"`
	} `json:"data"`
}

// NewAPIClient creates a new API client
func NewAPIClient(authToken, userID, outputDir string) *APIClient {
	return &APIClient{
		client:    &http.Client{Timeout: 30 * time.Second},
		authToken: authToken,
		userID:    userID,
		baseURL:   "https://service14.mywhoosh.com/v2/rider/profile",
		outputDir: outputDir,
	}
}

// GetActivities fetches all activities from the API
func (c *APIClient) GetActivities() ([]Activity, error) {
	var allActivities []Activity

	fmt.Println("📋 Fetching activities from MyWhoosh...")

	// Create request body with required parameters
	requestBody := map[string]interface{}{
		"page":     1,
		"sortDate": "DESC",
	}

	bodyData, err := json.Marshal(requestBody)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request body: %w", err)
	}

	req, err := http.NewRequest("POST", fmt.Sprintf("%s/activities", c.baseURL), bytes.NewReader(bodyData))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.authToken))
	req.Header.Set("User-ID", c.userID)
	req.Header.Set("Content-Type", "application/json")
	req.ContentLength = int64(len(bodyData))

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch activities: %w", err)
	}
	defer func() {
		_ = resp.Body.Close()
	}()

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusNoContent {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("API error: status %d, body: %s", resp.StatusCode, string(body))
	}

	if resp.StatusCode == http.StatusOK {
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return nil, err
		}

		var response ActivitiesResponse
		if err := json.Unmarshal(body, &response); err != nil {
			return nil, fmt.Errorf("failed to unmarshal activities: %w", err)
		}

		allActivities = response.Data.Results
	}

	fmt.Printf("✅ Found %d activities\n", len(allActivities))
	return allActivities, nil
}

// DownloadActivity downloads a single activity FIT file
func (c *APIClient) DownloadActivity(activity Activity) error {
	outputPath := filepath.Join(c.outputDir, fmt.Sprintf("%s_%s.fit",
		activity.StartTime.Format("20060102_150405"),
		activity.ActivityFileID))

	// Check if file already exists
	if _, err := os.Stat(outputPath); err == nil {
		fmt.Printf("⏭️  Skipping %s (already exists)\n", filepath.Base(outputPath))
		return nil
	}

	// Get presigned URL from MyWhoosh API
	presignedURL, err := c.GetPresignedURL(activity.ActivityFileID)
	if err != nil {
		return fmt.Errorf("failed to get presigned URL: %w", err)
	}

	fmt.Printf("⬇️  Downloading %s...\n", activity.ActivityFileID)

	resp, err := c.client.Get(presignedURL)
	if err != nil {
		return fmt.Errorf("failed to download from S3: %w", err)
	}
	defer func() {
		_ = resp.Body.Close()
	}()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("S3 download failed: status %d, body: %s", resp.StatusCode, string(body))
	}

	file, err := os.Create(outputPath)
	if err != nil {
		return fmt.Errorf("failed to create output file: %w", err)
	}
	defer func() {
		_ = file.Close()
	}()

	written, err := io.Copy(file, resp.Body)
	if err != nil {
		_ = os.Remove(outputPath)
		return fmt.Errorf("failed to write file: %w", err)
	}

	fmt.Printf("✅ Downloaded %s (%d bytes)\n", filepath.Base(outputPath), written)
	return nil
}

// GetPresignedURL gets a presigned S3 URL for downloading an activity file
func (c *APIClient) GetPresignedURL(activityFileID string) (string, error) {
	// Create request body with file ID
	requestBody := map[string]interface{}{
		"fileId": activityFileID,
	}

	bodyData, err := json.Marshal(requestBody)
	if err != nil {
		return "", fmt.Errorf("failed to marshal request body: %w", err)
	}

	req, err := http.NewRequest("POST", fmt.Sprintf("%s/download-activity-file", c.baseURL), bytes.NewReader(bodyData))
	if err != nil {
		return "", fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.authToken))
	req.Header.Set("User-ID", c.userID)
	req.Header.Set("Content-Type", "application/json")
	req.ContentLength = int64(len(bodyData))

	resp, err := c.client.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to fetch presigned URL: %w", err)
	}
	defer func() {
		_ = resp.Body.Close()
	}()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("API error: status %d, body: %s", resp.StatusCode, string(body))
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	var response struct {
		Error   bool   `json:"error"`
		Code    int    `json:"code"`
		Message string `json:"message"`
		Data    string `json:"data"`
	}

	if err := json.Unmarshal(body, &response); err != nil {
		return "", fmt.Errorf("failed to unmarshal response: %w", err)
	}

	if response.Error {
		return "", fmt.Errorf("API error: %s", response.Message)
	}

	return response.Data, nil
}

// SyncAllActivities downloads all activities in parallel
func (c *APIClient) SyncAllActivities(activities []Activity) error {
	if len(activities) == 0 {
		fmt.Println("ℹ️  No activities to sync")
		return nil
	}

	// Create output directory if it doesn't exist
	if err := os.MkdirAll(c.outputDir, 0755); err != nil {
		return fmt.Errorf("failed to create output directory: %w", err)
	}

	// Download activities sequentially
	var failedCount int
	for _, activity := range activities {
		if err := c.DownloadActivity(activity); err != nil {
			fmt.Printf("❌ %v\n", err)
			failedCount++
		}
	}

	successCount := len(activities) - failedCount
	fmt.Printf("\n📊 Sync complete: %d/%d successful\n", successCount, len(activities))

	if failedCount > 0 {
		return fmt.Errorf("%d downloads failed", failedCount)
	}

	return nil
}

// ListLocalFiles lists all FIT files in the output directory
func (c *APIClient) ListLocalFiles() ([]string, error) {
	var files []string

	entries, err := os.ReadDir(c.outputDir)
	if err != nil {
		if os.IsNotExist(err) {
			return files, nil
		}
		return nil, err
	}

	for _, entry := range entries {
		if !entry.IsDir() && filepath.Ext(entry.Name()) == ".fit" {
			files = append(files, entry.Name())
		}
	}

	return files, nil
}

func main() {
	authToken := flag.String("token", os.Getenv("MYWHOOSH_TOKEN"), "Authentication token (or MYWHOOSH_TOKEN env var)")
	userID := flag.String("user-id", os.Getenv("MYWHOOSH_USER_ID"), "User ID (or MYWHOOSH_USER_ID env var)")
	outputDir := flag.String("output", "./fit-files", "Output directory for FIT files")
	listOnly := flag.Bool("list", false, "List activities without downloading")

	flag.Parse()

	fmt.Println("🚴 MyWhoosh FIT File Exporter")
	fmt.Println("=========================================")

	// Validate authentication
	if *authToken == "" || *userID == "" {
		fmt.Println("❌ Error: authentication required via environment variables:")
		fmt.Println("   export MYWHOOSH_TOKEN=\"your-jwt-token\"")
		fmt.Println("   export MYWHOOSH_USER_ID=\"your-user-id\"")
		fmt.Println("")
		fmt.Println("Or use command-line flags:")
		fmt.Println("   ./mywhoosh-exporter -token <token> -user-id <id> -output <dir>")
		os.Exit(1)
	}

	fmt.Printf("\n👤 User ID: %s\n", *userID)

	// Create API client
	client := NewAPIClient(*authToken, *userID, *outputDir)

	// Fetch activities
	activities, err := client.GetActivities()
	if err != nil {
		fmt.Printf("❌ Error fetching activities: %v\n", err)
		os.Exit(1)
	}

	if *listOnly {
		fmt.Println("\n📋 Activities:")
		for i, a := range activities {
			fmt.Printf("%d. %s (%s) - %s\n", i+1, a.Title, a.ActivityFileID, a.StartTime.Format("2006-01-02"))
		}
	} else {
		// Sync all activities
		if err := client.SyncAllActivities(activities); err != nil {
			fmt.Printf("⚠️  Sync completed with errors: %v\n", err)
		}

		// List downloaded files
		fmt.Println("\n📁 Local FIT files:")
		localFiles, _ := client.ListLocalFiles()
		if len(localFiles) == 0 {
			fmt.Println("   (no files)")
		} else {
			for i, f := range localFiles {
				fpath := filepath.Join(*outputDir, f)
				info, _ := os.Stat(fpath)
				size := ""
				if info != nil {
					size = fmt.Sprintf(" (%d bytes)", info.Size())
				}
				fmt.Printf("%d. %s%s\n", i+1, f, size)
			}
		}
	}
}
