package psql_test

import (
	"bytes"
	"encoding/json"
	"testing"
)

func withComplexArgs(t *testing.T) {
	gql := `query {
		proDUcts(
			# returns only 30 items
			limit: 30,

			# starts from item 10, commented out for now
			# offset: 10,

			# orders the response items by highest price
			order_by: { price: desc },

			# no duplicate prices returned
			distinct: [ price ]

			# only items with an id >= 20 and < 28 are returned
			where: { id: { and: { greater_or_equals: 20, lt: 28 } } }) {
			id
			NAME
			price
		}
	}`

	compileGQLToPSQL(t, gql, nil, "user")
}

func withWhereIn(t *testing.T) {
	gql := `query {
		products(where: { id: { in: $list } }) {
			id
		}
	}`

	vars := map[string]json.RawMessage{
		"list": json.RawMessage(`[1,2,3]`),
	}

	compileGQLToPSQL(t, gql, vars, "user")
}

func withWhereAndList(t *testing.T) {
	gql := `query {
		products(
			where: {
				and: [
					{ not: { id: { is_null: true } } },
					{ price: { gt: 10 } },
				] } ) {
			id
			name
			price
		}
	}`

	compileGQLToPSQL(t, gql, nil, "user")
}

func withWhereIsNull(t *testing.T) {
	gql := `query {
		products(
			where: {
				and: {
					not: { id: { is_null: true } },
					price: { gt: 10 }
				}}) {
			id
			name
			price
		}
	}`

	compileGQLToPSQL(t, gql, nil, "user")
}

func withWhereMultiOr(t *testing.T) {
	gql := `query {
		products(
			where: {
				or: {
					not: { id: { is_null: true } },
					price: { gt: 10 },
					price: { lt: 20 }
				} }
			) {
			id
			name
			price
		}
	}`

	compileGQLToPSQL(t, gql, nil, "user")
}

func fetchByID(t *testing.T) {
	gql := `query {
		product(id: $id) {
			id
			name
		}
	}`

	compileGQLToPSQL(t, gql, nil, "user")
}

func searchQuery(t *testing.T) {
	gql := `query {
		products(search: $query) {
			id
			name
			search_rank
			search_headline_description
		}
	}`

	compileGQLToPSQL(t, gql, nil, "admin")
}

func oneToMany(t *testing.T) {
	gql := `query {
		users {
			email
			products {
				name
				price
			}
		}
	}`

	compileGQLToPSQL(t, gql, nil, "user")
}

func oneToManyReverse(t *testing.T) {
	gql := `query {
		products {
			name
			price
			users {
				email
			}
		}
	}`

	compileGQLToPSQL(t, gql, nil, "user")
}

func oneToManyArray(t *testing.T) {
	gql := `
	query {
		product {
			name
			price
			tags {
				id
				name
			}
		}
		tags {
			name
			product {
				name
			}
		}
	}`

	compileGQLToPSQL(t, gql, nil, "admin")
}

func manyToMany(t *testing.T) {
	gql := `query {
		products {
			name
			customers {
				email
				full_name
			}
		}
	}`

	compileGQLToPSQL(t, gql, nil, "user")
}

func manyToManyReverse(t *testing.T) {
	gql := `query {
		customers {
			email
			full_name
			products {
				name
			}
		}
	}`

	compileGQLToPSQL(t, gql, nil, "user")
}

func aggFunction(t *testing.T) {
	gql := `query {
		products {
			name
			count_price
		}
	}`

	compileGQLToPSQL(t, gql, nil, "user")
}

func aggFunctionBlockedByCol(t *testing.T) {
	gql := `query {
		products {
			name
			count_price
		}
	}`

	compileGQLToPSQL(t, gql, nil, "anon")
}

func aggFunctionDisabled(t *testing.T) {
	gql := `query {
		products {
			name
			count_price
		}
	}`

	compileGQLToPSQL(t, gql, nil, "anon1")
}

func aggFunctionWithFilter(t *testing.T) {
	gql := `query {
		products(where: { id: { gt: 10 } }) {
			id
			max_price
		}
	}`

	compileGQLToPSQL(t, gql, nil, "user")
}

func syntheticTables(t *testing.T) {
	gql := `query {
		me {
			email
		}
	}`

	compileGQLToPSQL(t, gql, nil, "user")
}

func queryWithVariables(t *testing.T) {
	gql := `query {
		product(id: $PRODUCT_ID, where: { price: { eq: $PRODUCT_PRICE } }) {
			id
			name
		}
	}`

	compileGQLToPSQL(t, gql, nil, "user")
}

func withWhereOnRelations(t *testing.T) {
	gql := `query {
		users(where: { 
				not: { 
					products: { 
						price: { gt: 3 }
					} 
				} 
			}) {
			id
			email
		}
	}`

	compileGQLToPSQL(t, gql, nil, "user")
}

func multiRoot(t *testing.T) {
	gql := `query {
		product {
			id
			name
			customer {
				email
			}
			customers {
				email
			}
		}
		user {
			id
			email
		}
		customer {
			id
		}
	}`

	compileGQLToPSQL(t, gql, nil, "user")
}

func withFragment1(t *testing.T) {
	gql := `
	fragment userFields1 on user {
		id
		email
	}

	query {
		users {
			...userFields2
	
			created_at
			...userFields1
		}
	}
	
	fragment userFields2 on user {
		first_name
		last_name
	}`

	compileGQLToPSQL(t, gql, nil, "anon")
}

func withFragment2(t *testing.T) {
	gql := `
	query {
		users {
			...userFields2
	
			created_at
			...userFields1
		}
	}

	fragment userFields1 on user {
		id
		email
	}
	
	fragment userFields2 on user {
		first_name
		last_name
	}`

	compileGQLToPSQL(t, gql, nil, "anon")
}

func withFragment3(t *testing.T) {
	gql := `

	fragment userFields1 on user {
		id
		email
	}
	
	fragment userFields2 on user {
		first_name
		last_name
	}

	query {
		users {
			...userFields2
	
			created_at
			...userFields1
		}
	}
`

	compileGQLToPSQL(t, gql, nil, "anon")
}

// func withInlineFragment(t *testing.T) {
// 	gql := `
// 	query {
// 		users {
// 			... on users {
// 				id
// 				email
// 			}
// 			created_at
// 			... on user {
// 				first_name
// 				last_name
// 			}
// 		}
// 	}
// `

// 	compileGQLToPSQL(t, gql, nil, "anon")
// }

func withCursor(t *testing.T) {
	gql := `query {
		Products(
			first: 20
			after: $cursor
			order_by: { price: desc }) {
			Name
		}
	}`

	vars := map[string]json.RawMessage{
		"cursor": json.RawMessage(`"0,1"`),
	}

	compileGQLToPSQL(t, gql, vars, "admin")
}

func jsonColumnAsTable(t *testing.T) {
	gql := `query {
		products {
			id
			name
			tag_count {
				count
				tags {
					name
				}
			}
		}
	}`

	compileGQLToPSQL(t, gql, nil, "admin")
}

func nullForAuthRequiredInAnon(t *testing.T) {
	gql := `query {
		products {
			id
			name
			user(where: { id: { eq: $user_id } }) {
				id
				email
			}
		}
	}`

	compileGQLToPSQL(t, gql, nil, "anon")
}

func blockedQuery(t *testing.T) {
	gql := `query {
		user(id: $id, where: { id: { gt: 3 } }) {
			id
			full_name
			email
		}
	}`

	compileGQLToPSQL(t, gql, nil, "bad_dude")
}

func blockedFunctions(t *testing.T) {
	gql := `query {
		users {
			count_id
			email
		}
	}`

	compileGQLToPSQL(t, gql, nil, "bad_dude")
}

func TestCompileQuery(t *testing.T) {
	t.Run("withComplexArgs", withComplexArgs)
	t.Run("withWhereIn", withWhereIn)
	t.Run("withWhereAndList", withWhereAndList)
	t.Run("withWhereIsNull", withWhereIsNull)
	t.Run("withWhereMultiOr", withWhereMultiOr)
	t.Run("fetchByID", fetchByID)
	t.Run("searchQuery", searchQuery)
	t.Run("oneToMany", oneToMany)
	t.Run("oneToManyReverse", oneToManyReverse)
	t.Run("oneToManyArray", oneToManyArray)
	t.Run("manyToMany", manyToMany)
	t.Run("manyToManyReverse", manyToManyReverse)
	t.Run("aggFunction", aggFunction)
	t.Run("aggFunctionBlockedByCol", aggFunctionBlockedByCol)
	t.Run("aggFunctionDisabled", aggFunctionDisabled)
	t.Run("aggFunctionWithFilter", aggFunctionWithFilter)
	t.Run("syntheticTables", syntheticTables)
	t.Run("queryWithVariables", queryWithVariables)
	t.Run("withWhereOnRelations", withWhereOnRelations)
	t.Run("multiRoot", multiRoot)
	t.Run("withFragment1", withFragment1)
	t.Run("withFragment2", withFragment2)
	t.Run("withFragment3", withFragment3)
	//t.Run("withInlineFragment", withInlineFragment)
	t.Run("jsonColumnAsTable", jsonColumnAsTable)
	t.Run("withCursor", withCursor)
	t.Run("nullForAuthRequiredInAnon", nullForAuthRequiredInAnon)
	t.Run("blockedQuery", blockedQuery)
	t.Run("blockedFunctions", blockedFunctions)
}

var benchGQL = []byte(`query {
	proDUcts(
		# returns only 30 items
		limit: 30,

		# starts from item 10, commented out for now
		# offset: 10,

		# orders the response items by highest price
		order_by: { price: desc },

		# only items with an id >= 30 and < 30 are returned
		where: { id: { and: { greater_or_equals: 20, lt: 28 } } }) {
		id
		NAME
		price
		user {
			full_name
			picture : avatar
		}
	}
}`)

func BenchmarkCompile(b *testing.B) {
	w := &bytes.Buffer{}

	b.ResetTimer()
	b.ReportAllocs()

	for n := 0; n < b.N; n++ {
		w.Reset()

		qc, err := qcompile.Compile(benchGQL, "user")
		if err != nil {
			b.Fatal(err)
		}

		_, err = pcompile.Compile(w, qc, nil)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkCompileParallel(b *testing.B) {
	b.ReportAllocs()

	b.RunParallel(func(pb *testing.PB) {
		w := &bytes.Buffer{}

		for pb.Next() {
			w.Reset()

			qc, err := qcompile.Compile(benchGQL, "user")
			if err != nil {
				b.Fatal(err)
			}

			_, err = pcompile.Compile(w, qc, nil)
			if err != nil {
				b.Fatal(err)
			}
		}
	})
}
