# Okta User and Group Management with Terraform

This repository provides a Terraform configuration to create and manage Okta users and groups based on a CSV file. The setup includes the following:

- Creates Okta users with specified attributes.
- Creates Okta groups by department.
- Assigns users to their respective department groups.
- Outputs the list of created users and groups.

## Directory Structure

- `main.tf`: Main Terraform configuration file.
- `variables.tf`: Defines the variables used in the Terraform configuration.
- `versions.tf`: Specifies the required Terraform and provider versions.
- `terraform.tfvars`: Contains the values for the variables.
- `outputs.tf`: Defines the outputs for the Terraform configuration.
- `users.csv`: CSV file containing user data.
- `generate_hcl.go`: Go script to convert the CSV file to HCL format.

## Prerequisites

- [Terraform](https://www.terraform.io/downloads.html) installed.
- [Go](https://golang.org/dl/) installed.
- Okta account with API token.

## Setup Instructions

1. Place the `users.csv` file in the same directory as the Terraform files and the Go script.

   ```csv
   login,firstName,lastName,email,title,displayName,nickName,userType,organization,department,division,startDate
   jdoe,John,Doe,jdoe@example.com,Engineer,John Doe,Johnny,employee,Acme Corp,Engineering,Development,2023-01-15
   asmith,Alice,Smith,asmith@example.com,Manager,Alice Smith,AliceM,employee,Acme Corp,Engineering,Management,2023-01-20
   ```

2. Run the Go script to generate the `variables.auto.tfvars` file:

   ```sh
   go run generate_hcl.go
   ```

3. Initialize Terraform:

   ```sh
   terraform init
   ```

4. Apply the Terraform configuration:

   ```sh
   terraform apply
   ```

## Files

```hcl name=main.tf
provider "okta" {
  org_name  = var.okta_org_name
  api_token = var.okta_api_token
}

resource "okta_user" "users" {
  for_each = var.users

  login        = each.value.login
  email        = each.value.email
  first_name   = each.value.firstName
  last_name    = each.value.lastName
  display_name = each.value.displayName
  nick_name    = each.value.nickName
  title        = each.value.title
  user_type    = each.value.userType
  organization = each.value.organization
  department   = each.value.department
  division     = each.value.division

  custom_profile_attributes = jsonencode({
    start_date = each.value.startDate
  })

  lifecycle {
    ignore_changes = [password]
  }
}

resource "okta_group" "department_groups" {
  for_each = toset([for user in var.users : user.department])

  name = each.value
}

resource "okta_user_group_memberships" "group_membership" {
  for_each = var.users

  user_id = okta_user.users[each.key].id
  groups  = [okta_group.department_groups[each.value.department].id]
}
```

```hcl name=variables.tf
variable "okta_org_name" {
  description = "Okta organization name"
  type        = string
}

variable "okta_api_token" {
  description = "Okta API token"
  type        = string
  sensitive   = true
}

variable "csv_file" {
  description = "Path to the CSV file containing user data"
  type        = string
}

variable "users" {
  description = "List of users to be created in Okta"
  type = map(object({
    login        = string
    firstName    = string
    lastName     = string
    email        = string
    title        = string
    displayName  = string
    nickName     = string
    userType     = string
    organization = string
    department   = string
    division     = string
    startDate    = string
  }))
}
```

```hcl name=versions.tf
terraform {
  required_providers {
    okta = {
      source  = "okta/okta"
      version = "3.23.0"
    }
  }
  required_version = ">= 1.0.0"
}
```

```hcl name=terraform.tfvars.example
(rename to terraform.tfvars)
okta_org_name  = "your-okta-org-name" #omit okta.com
okta_api_token = "your-okta-api-token"
csv_file       = "path-to-your-csv-file/users.csv"
```

```hcl name=outputs.tf
output "created_users" {
  value = [for user in okta_user.users : {
    login       = user.login
    email       = user.email
    first_name  = user.first_name
    last_name   = user.last_name
    department  = user.department
    start_date  = jsondecode(user.custom_profile_attributes)["start_date"]
  }]
}

output "created_groups" {
  value = [for group in okta_group.department_groups : group.name]
}
```

```go name=generate_hcl.go
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