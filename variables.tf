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