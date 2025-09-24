package graph

import "github.com/tabed23/cloudmarket-auth/graph/repos"

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.

type Resolver struct{
	repos.Repository
}
