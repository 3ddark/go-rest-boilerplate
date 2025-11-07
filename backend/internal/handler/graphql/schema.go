package graphql

import (
	"fmt"

	"ths-erp.com/internal/auth"
	"ths-erp.com/internal/dto"
	"ths-erp.com/internal/service"

	"github.com/graphql-go/graphql"
)

// buildGraphQLSchema, servisleri kullanarak GraphQL şemasını oluşturur.
func buildGraphQLSchema(userService service.IUserService, permService service.IPermissionService) (graphql.Schema, error) {
	// User Tipi
	userType := graphql.NewObject(graphql.ObjectConfig{
		Name: "User",
		Fields: graphql.Fields{
			"id":    &graphql.Field{Type: graphql.Int},
			"name":  &graphql.Field{Type: graphql.String},
			"email": &graphql.Field{Type: graphql.String},
		},
	})

	// Login Yanıt Tipi
	loginResponseType := graphql.NewObject(graphql.ObjectConfig{
		Name: "LoginResponse",
		Fields: graphql.Fields{
			"token": &graphql.Field{Type: graphql.String},
			"user":  &graphql.Field{Type: userType},
		},
	})

	// Root Query
	rootQuery := graphql.NewObject(graphql.ObjectConfig{
		Name:   "Query",
		Fields: buildQueryFields(userService, permService, userType),
	})

	// Root Mutation
	rootMutation := graphql.NewObject(graphql.ObjectConfig{
		Name:   "Mutation",
		Fields: buildMutationFields(userService, permService, userType, loginResponseType),
	})

	// Şemayı Oluştur
	return graphql.NewSchema(graphql.SchemaConfig{
		Query:    rootQuery,
		Mutation: rootMutation,
	})
}

// buildQueryFields, GraphQL Query alanlarını oluşturur.
func buildQueryFields(userService service.IUserService, permService service.IPermissionService, userType *graphql.Object) graphql.Fields {
	return graphql.Fields{
		"user": &graphql.Field{
			Type: userType,
			Args: graphql.FieldConfigArgument{
				"id": &graphql.ArgumentConfig{Type: graphql.NewNonNull(graphql.Int)},
			},
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				authUser, err := auth.GetUserFromContext(p.Context)
				if err != nil {
					return nil, err
				}

				if allowed, err := permService.CheckPermission(p.Context, authUser.UserID, "users", "select"); !allowed || err != nil {
					if err != nil {
						return nil, fmt.Errorf("permission check failed: %v", err)
					}
					return nil, fmt.Errorf("permission denied")
				}

				id := p.Args["id"].(int)
				return userService.GetUser(p.Context, id)
			},
		},
		"users": &graphql.Field{
			Type: graphql.NewList(userType),
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				authUser, err := auth.GetUserFromContext(p.Context)
				if err != nil {
					return nil, err
				}

				if allowed, err := permService.CheckPermission(p.Context, authUser.UserID, "users", "select"); !allowed || err != nil {
					if err != nil {
						return nil, fmt.Errorf("permission check failed: %v", err)
					}
					return nil, fmt.Errorf("permission denied")
				}

				return userService.GetAllUsers(p.Context)
			},
		},
		"me": &graphql.Field{
			Type: userType,
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				authUser, err := auth.GetUserFromContext(p.Context)
				if err != nil {
					return nil, err
				}

				return userService.GetUser(p.Context, authUser.UserID)
			},
		},
	}
}

// buildMutationFields, GraphQL Mutation alanlarını oluşturur.
func buildMutationFields(userService service.IUserService, permService service.IPermissionService, userType, loginResponseType *graphql.Object) graphql.Fields {
	return graphql.Fields{
		"login": &graphql.Field{
			Type: loginResponseType,
			Args: graphql.FieldConfigArgument{
				"email":    &graphql.ArgumentConfig{Type: graphql.NewNonNull(graphql.String)},
				"password": &graphql.ArgumentConfig{Type: graphql.NewNonNull(graphql.String)},
			},
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				email := p.Args["email"].(string)
				password := p.Args["password"].(string)

				user, err := userService.Authenticate(p.Context, email, password)
				if err != nil {
					return nil, fmt.Errorf("invalid credentials")
				}
				token, err := auth.GenerateJWT(user.ID, user.Email)
				if err != nil {
					return nil, fmt.Errorf("could not generate token")
				}

				userResponse, err := userService.GetUser(p.Context, user.ID)
				if err != nil {
					return nil, err
				}

				return dto.LoginResponse{Token: token, User: userResponse}, nil
			},
		},
		"createUser": &graphql.Field{
			Type: userType,
			Args: graphql.FieldConfigArgument{
				"name":     &graphql.ArgumentConfig{Type: graphql.NewNonNull(graphql.String)},
				"email":    &graphql.ArgumentConfig{Type: graphql.NewNonNull(graphql.String)},
				"password": &graphql.ArgumentConfig{Type: graphql.NewNonNull(graphql.String)},
			},
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				authUser, err := auth.GetUserFromContext(p.Context)
				if err != nil {
					return nil, err
				}

				if allowed, err := permService.CheckPermission(p.Context, authUser.UserID, "users", "add"); !allowed || err != nil {
					if err != nil {
						return nil, fmt.Errorf("permission check failed: %v", err)
					}
					return nil, fmt.Errorf("permission denied")
				}

				req := &dto.CreateUserRequest{
					Name:     p.Args["name"].(string),
					Email:    p.Args["email"].(string),
					Password: p.Args["password"].(string),
				}
				return userService.CreateUser(p.Context, req)
			},
		},
		"updateUser": &graphql.Field{
			Type: userType,
			Args: graphql.FieldConfigArgument{
				"id":    &graphql.ArgumentConfig{Type: graphql.NewNonNull(graphql.Int)},
				"name":  &graphql.ArgumentConfig{Type: graphql.String},
				"email": &graphql.ArgumentConfig{Type: graphql.String},
			},
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				authUser, err := auth.GetUserFromContext(p.Context)
				if err != nil {
					return nil, err
				}

				if allowed, err := permService.CheckPermission(p.Context, authUser.UserID, "users", "update"); !allowed || err != nil {
					if err != nil {
						return nil, fmt.Errorf("permission check failed: %v", err)
					}
					return nil, fmt.Errorf("permission denied")
				}

				id := p.Args["id"].(int)
				req := &dto.UpdateUserRequest{}
				if name, ok := p.Args["name"].(string); ok {
					req.Name = name
				}
				if email, ok := p.Args["email"].(string); ok {
					req.Email = email
				}
				return userService.UpdateUser(p.Context, id, req)
			},
		},
		"deleteUser": &graphql.Field{
			Type: graphql.Boolean,
			Args: graphql.FieldConfigArgument{
				"id": &graphql.ArgumentConfig{Type: graphql.NewNonNull(graphql.Int)},
			},
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				authUser, err := auth.GetUserFromContext(p.Context)
				if err != nil {
					return nil, err
				}

				if allowed, err := permService.CheckPermission(p.Context, authUser.UserID, "users", "delete"); !allowed || err != nil {
					if err != nil {
						return nil, fmt.Errorf("permission check failed: %v", err)
					}
					return nil, fmt.Errorf("permission denied")
				}

				id := p.Args["id"].(int)
				err = userService.DeleteUser(p.Context, id)
				return err == nil, err
			},
		},
	}
}
