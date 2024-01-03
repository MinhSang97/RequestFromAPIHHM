package main

import (
	"bufio"
	"database/sql"
	"fmt"
	"os"

	"demoAPI/dbutil"

	"demoAPI/callapi" // Update this line based on the actual directory structure
)

func main() {
	db, err := dbutil.ConnectDB()
	if err != nil {
		fmt.Println("Failed to connect to the database:", err)
		return
	}
	fmt.Println("Connected: ", db)

	rows, err := db.Raw("SELECT METER_ASSET_NO congTo FROM A_Data_catalogue WHERE STATUS ='1'").Rows()
	if err != nil {
		fmt.Println("Query error:", err)
		return
	}
	defer rows.Close()

	// var congToValues []string

	var congToValues []string

	for rows.Next() {
		var congTo sql.NullString
		if err := rows.Scan(&congTo); err != nil {
			fmt.Println("Scan error:", err)
			return
		}
		if congTo.Valid {
			fmt.Println(congTo.String)
			congToValues = append(congToValues, congTo.String) // Append the value to the slice
		} else {
			fmt.Println("(NULL)")
		}
	}

	if err := rows.Err(); err != nil {
		fmt.Println("Rows iteration error:", err)
	}

	err = saveToFile("output.txt", congToValues)
	if err != nil {
		fmt.Println("Error saving to file:", err)
	}

	//Gọi hàm trong callapi để lấy dữ liệu từ API
	callapi.GetDataFromAPI("output.txt", "admin", "hhm@1997")
	callapi.GetDataFromAPIWithRetry("output.txt", "admin", "hhm@1997")
}

func saveToFile(filename string, content []string) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	writer := bufio.NewWriter(file)
	for _, value := range content {
		_, err = writer.WriteString(value + "\n")
		if err != nil {
			return err
		}
	}

	writer.Flush()

	fmt.Printf("Data saved to %s\n", filename)
	return nil
}
