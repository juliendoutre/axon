package authz

import rego.v1

default allowed := false

allowed if {
    input.claims.iss = "https://gitlab.com"
}
