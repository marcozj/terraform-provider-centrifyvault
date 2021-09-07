---
subcategory: "Access"
---

# centrifyvault_userpassword (Resource)

This resource allows you to update password for Centrify Directory User.

## Example Usage

```terraform
data "centrifyvault_user" "testuser" {
    username = "testuser@example.com"
}

resource "centrifyvault_userpassword" "testuser" {
  user_uuid = data.centrifyvault_user.testuser.id
  password  = "xxxxxxxxxxx"
}
```

## Argument Reference

### Required

- `user_uuid` - (String) The UUID of Centrify Directory User.
- `password` - (String, Sensitive) Password of the user.

**Limitation:** `userpassword` isn't supported in import process.
