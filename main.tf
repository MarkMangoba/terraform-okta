provider "okta" {
  org_name  = var.okta_org_name
  api_token = var.okta_api_token
}

# Create users
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

# Create groups by department
resource "okta_group" "department_groups" {
  for_each = toset([for user in var.users : user.department])

  name = each.value
}

# Assign users to groups based on department
resource "okta_user_group_memberships" "group_membership" {
  for_each = var.users

  user_id = okta_user.users[each.key].id
  groups  = [okta_group.department_groups[each.value.department].id]
}