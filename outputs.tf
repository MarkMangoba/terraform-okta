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