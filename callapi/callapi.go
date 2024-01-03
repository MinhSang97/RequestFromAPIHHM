package callapi

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// GetDataFromAPI lấy giá trị congTo và thoiGian từ file và gọi API để lấy dữ liệu
func GetDataFromAPI(filename, user, password string) {
	// Đọc giá trị congTo từ file
	congToValues, err := readFromFile(filename)
	if err != nil {
		fmt.Println("Error reading from file:", err)
		return
	}

	for _, congTo := range congToValues {
		startTime := time.Date(2024, time.January, 1, 0, 0, 0, 0, time.UTC)
		endTime := time.Date(2024, time.January, 3, 0, 0, 0, 0, time.UTC)

		for !startTime.After(endTime) {
			// Delay for 10 seconds before making the API call
			time.Sleep(2 * time.Second)

			// Gọi hàm trong callAPI để lấy dữ liệu từ API
			data, err := getData(congTo, startTime.Format("02-01-06"), user, password)
			if err != nil {
				fmt.Println("Lỗi khi lấy dữ liệu từ API:", err)
				return
			}

			// In kết quả
			fmt.Printf("Data %s và thoiGian %s: %s\n", congTo, startTime.Format("02-01-06"), string(data))

			// Lưu data vào MongoDB
			err = saveDataToMongoDB(data)
			if err != nil {
				fmt.Println("Lỗi khi lưu dữ liệu vào MongoDB:", err)
				return
			}

			// Tăng thời gian lên 1 ngày
			startTime = startTime.Add(24 * time.Hour)
		}
	}
}

const (
	mongoURI       = "mongodb://localhost:27017"
	databaseName   = "demoAPI"
	collectionName = "demoAPI"
)

func saveDataToMongoDB(data []byte) error {
	// Parse JSON data
	var jsonData []map[string]interface{}
	if err := json.Unmarshal(data, &jsonData); err != nil {
		return err
	}

	// Connect to MongoDB
	client, err := mongo.Connect(context.Background(), options.Client().ApplyURI(mongoURI))
	if err != nil {
		return err
	}
	defer client.Disconnect(context.Background())

	// Access the specified database and collection
	db := client.Database(databaseName)
	collection := db.Collection(collectionName)

	// Insert data into MongoDB
	for _, document := range jsonData {
		_, err := collection.InsertOne(context.Background(), document)
		if err != nil {
			return err
		}
	}

	return nil
}

func readFromFile(filename string) ([]string, error) {
	file, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	content := string(file)
	congToValues := strings.Split(content, "\n")

	return congToValues, nil
}

// Hàm để gọi API data với token
func getData(congTo, thoiGian, user, password string) ([]byte, error) {
	loginURL := "http://118.69.35.119:62000/login"
	dataURL := "http://118.69.35.119:62000/data"

	// Thực hiện logic lấy token tương tự như trong main.go
	token, err := getAuthToken(loginURL, user, password)
	if err != nil {
		fmt.Println("Error refreshing token:", err)
		return nil, err
	}

	// Gọi API data với token được nhận
	data, err := getDataFromAPI(dataURL, token, congTo, thoiGian)
	if err != nil {
		return nil, err
	}

	return data, nil
}

func getDataFromAPI(url, token, congTo, thoiGian string) ([]byte, error) {
	// Tạo URL với các tham số
	apiURLWithParams := fmt.Sprintf("%s?cong_to=%s&thoi_gian=%s", url, congTo, thoiGian)

	// Tạo request
	req, err := http.NewRequest("GET", apiURLWithParams, nil)
	if err != nil {
		return nil, err
	}

	// Thêm Authorization header với key
	req.Header.Add("Authorization", token)

	// Gửi request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// Đọc dữ liệu
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return body, nil
}

// Hàm để gọi API login và nhận token

const tokenFilePath = "token.txt"

const retryDelay = 2 * time.Minute

// GetDataFromAPIWithRetry lấy giá trị congTo và thoiGian từ file và gọi API với việc thử lại khi có lỗi
func GetDataFromAPIWithRetry(filename, user, password string) {
	for {
		// Thực hiện lấy dữ liệu từ API
		err := getDataFromAPIWithRetry(filename, user, password)
		if err != nil {
			fmt.Println("Error:", err)
			fmt.Printf("Retry after %v...\n", retryDelay)
			time.Sleep(retryDelay)
		}
	}
}

// getDataFromAPIWithRetry thực hiện lấy dữ liệu từ API với việc thử lại khi có lỗi
func getDataFromAPIWithRetry(filename, user, password string) error {
	// Đọc giá trị congTo từ file
	congToValues, err := readFromFile(filename)
	if err != nil {
		return fmt.Errorf("Error reading from file: %v", err)
	}

	for _, congTo := range congToValues {
		startTime := time.Date(2024, time.January, 1, 0, 0, 0, 0, time.UTC)
		endTime := time.Date(2024, time.January, 3, 0, 0, 0, 0, time.UTC)

		for !startTime.After(endTime) {
			// Delay for 10 seconds before making the API call
			time.Sleep(2 * time.Second)

			// Gọi hàm trong callAPI để lấy dữ liệu từ API
			data, err := getData(congTo, startTime.Format("02-01-06"), user, password)
			if err != nil {
				return fmt.Errorf("Lỗi khi lấy dữ liệu từ API: %v", err)
			}

			// In kết quả
			fmt.Printf("Data %s và thoiGian %s: %s\n", congTo, startTime.Format("02-01-06"), string(data))

			// Lưu data vào MongoDB
			err = saveDataToMongoDB(data)
			if err != nil {
				return fmt.Errorf("Lỗi khi lưu dữ liệu vào MongoDB: %v", err)
			}

			// Tăng thời gian lên 1 ngày
			startTime = startTime.Add(24 * time.Hour)
		}
	}

	return nil
}

// Hàm để gọi API login và nhận token
func getAuthToken(url, user, password string) (string, error) {
	// Thực hiện logic đăng nhập mới
	apiURLWithParams := fmt.Sprintf("%s?user=%s&password=%s", url, user, password)
	response, err := http.Get(apiURLWithParams)
	if err != nil {
		return "", err
	}
	defer response.Body.Close()

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return "", err
	}

	var result map[string]interface{}
	if err := json.Unmarshal(body, &result); err != nil {
		return "", err
	}

	token, ok := result["token"].(string)
	if !ok {
		return "", fmt.Errorf("Không thể lấy token từ kết quả API login: %v", result)
	}

	// Lưu token vào file
	err = writeTokenToFile(token)
	if err != nil {
		fmt.Println("Error writing token to file:", err)
	}

	return token, nil
}

// Hàm kiểm tra xem token có hết hạn hay không
func isTokenExpired(token string) bool {
	// Thực hiện logic kiểm tra thời hạn của token
	// Đây là một giả định, bạn cần thay đổi nó dựa trên cách API của bạn xác định thời hạn token
	// Ví dụ: nếu token có thời hạn 8 tiếng, bạn có thể kiểm tra thời hạn bằng cách so sánh thời điểm hiện tại với thời điểm tạo token
	// Đoạn mã dưới đây chỉ là một ví dụ giả định và cần được điều chỉnh cho đúng với API của bạn.
	// Thời hạn token được giả sử là 8 tiếng.
	tokenCreationTime, err := readTokenCreationTime()
	if err != nil {
		return true // Trả về true để đảm bảo đăng nhập lại nếu không thể đọc thời gian tạo token
	}
	tokenExpirationTime := tokenCreationTime.Add(8 * time.Hour)

	return time.Now().After(tokenExpirationTime)
}

// Hàm đọc token từ file
func readTokenFromFile() (string, error) {
	file, err := ioutil.ReadFile(tokenFilePath)
	if err != nil {
		return "", err
	}

	return string(file), nil
}

// Hàm ghi token vào file
func writeTokenToFile(token string) error {
	return ioutil.WriteFile(tokenFilePath, []byte(token), 0644)
}

// Hàm đọc thời gian tạo token từ file (giả sử)
func readTokenCreationTime() (time.Time, error) {
	// Đọc dữ liệu từ file và chuyển đổi thành thời gian
	file, err := ioutil.ReadFile("token_creation_time.txt")
	if err != nil {
		return time.Time{}, err
	}

	strTime := string(file)
	parsedTime, err := time.Parse(time.RFC3339, strTime)
	if err != nil {
		return time.Time{}, err
	}

	return parsedTime, nil
}

// Hàm ghi thời gian tạo token vào file (giả sử)
func writeTokenAndCreationTimeToFile(token string, creationTime time.Time) error {
	// Ghi token vào file
	err := ioutil.WriteFile(tokenFilePath, []byte(token), 0644)
	if err != nil {
		return err
	}

	// Ghi thời gian tạo token vào file
	err = ioutil.WriteFile("token_creation_time.txt", []byte(creationTime.Format(time.RFC3339)), 0644)
	if err != nil {
		return err
	}

	return nil
}
