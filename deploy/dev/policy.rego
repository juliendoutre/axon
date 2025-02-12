package authz

import rego.v1

default allowed := false

allowed if {
    input.iss = "http://localhost:7080/realms/master"
}
