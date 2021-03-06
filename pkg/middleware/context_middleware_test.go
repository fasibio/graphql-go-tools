package middleware

import (
	"context"
	"fmt"
	"github.com/jensneuse/graphql-go-tools/pkg/testhelper"
	"testing"
)

func TestContextMiddleware(t *testing.T) {
	t.Run("simple", func(t *testing.T) {
		ctx := context.Background()
		// it's important to quote the value so the lexer will recognize it's a string value
		// we might push this including checks into the implementation
		ctx = context.WithValue(ctx, "user", []byte(`"jsmith@example.org"`))

		got := InvokeMiddleware(&ContextMiddleware{}, ctx, publicSchema, publicQuery)
		want := testhelper.UglifyRequestString(privateQuery)

		if want != got {
			panic(fmt.Errorf("\nwant:\n%s\ngot:\n%s", want, got))
		}
	})
}

const publicSchema = `
directive @addArgumentFromContext(
	name: String!
	contextKey: String!
) on FIELD_DEFINITION

scalar String

schema {
	query: Query
}

type Query {
	documents: [Document] @addArgumentFromContext(name: "user",contextKey: "user")
}

type Document implements Node {
	owner: String
	sensitiveInformation: String
}
`

/*

the public schema for reference

schema {
	query: Query
}

type Query {
	documents(user: String!): [Document]
}

type Document implements Node {
	owner: String
	sensitiveInformation: String
}
*/

const publicQuery = `
query myDocuments {
	documents {
		sensitiveInformation
	}
}
`

const privateQuery = `
query myDocuments {
	documents(user: "jsmith@example.org") {
		sensitiveInformation
	}
}
`
