package main

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"os"
)

type User struct {
	Login        string `json:"login"`
	FirstName    string `json:"firstName"`
	LastName     string `json:"lastName"`
	Email        string `json:"email"`
	Title        string `json:"title"`
	DisplayName  string `json:"displayName"`
	NickName     string `json:"nickName"`
	UserType     string `json:"userType"`
	Organization string `json:"organization"`
	Department   string `json:"department"`
	Division     string `json:"division"`
	StartDate    string `json:"startDate"`
}

func csvToHcl(csvFile, hclFile string) error {
	file, err := os.Open(csvFile)
	if err != nil {
		return err
	}
	defer file.Close()

	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		return err
	}

	if len(records) < 1 {
		return fmt.Errorf("no records found in CSV file")
	}

	headers := records[0]
	users := make(map[string]User)

	for _, record := range records[1:] {
		user := User{}
		for i, header := range headers {
			switch header {
			case "login":
				user.Login = record[i]
			case "firstName":
				user.FirstName = record[i]
			case "lastName":
				user.LastName = record[i]
			case "email":
				user.Email = record[i]
			case "title":
				user.Title = record[i]
			case "displayName":
				user.DisplayName = record[i]
			case "nickName":
				user.NickName = record[i]
			case "userType":
				user.UserType = record[i]
			case "organization":
				user.Organization = record[i]
			case "department":
				user.Department = record[i]
			case "division":
				user.Division = record[i]
			case "startDate":
				user.StartDate = record[i]
			}
		}
		users[user.Login] = user
	}

	usersJson, err := json.MarshalIndent(users, "", "  ")
	if err != nil {
		return err
	}

	hclContent := fmt.Sprintf(`
users = %s
`, string(usersJson))

	return os.WriteFile(hclFile, []byte(hclContent), 0644)
}

func main() {
	csvFile := "users.csv"
	hclFile := "variables.auto.tfvars"

	err := csvToHcl(csvFile, hclFile)
	if err != nil {
		fmt.Println("Error:", err)
	} else {
		fmt.Println("HCL file generated successfully.")
	}
}
