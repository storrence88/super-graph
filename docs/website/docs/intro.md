---
id: intro
title: Introduction
sidebar_label: Introduction
---

import useBaseUrl from '@docusaurus/useBaseUrl'; // Add to the top of the file below the front matter.

Super Graph is a service that instantly and without code gives you a high performance and secure GraphQL API. Your GraphQL queries are auto translated into a single fast SQL query. No more spending weeks or months writing backend API code. Just make the query you need and Super Graph will do the rest.

Super Graph has a rich feature set like integrating with your existing Ruby on Rails apps, joining your DB with data from remote APIs, Role and Attribute based access control, Support for JWT tokens, DB migrations, seeding and a lot more.

## Features

- Role and Attribute based access control
- Works with existing Ruby-On-Rails apps
- Automatically learns database schemas and relationships
- Full text search and aggregations
- Rails authentication supported (Redis, Memcache, Cookie)
- JWT tokens supported (Auth0, etc)
- Join database with remote REST APIs
- Highly optimized and fast Postgres SQL queries
- GraphQL queries and mutations
- A simple config file
- High performance GO codebase
- Tiny docker image and low memory requirements
- Fuzz tested for security
- Database migrations tool
- Database seeding tool
- Works with Postgres and Yugabyte DB
- OpenCensus Support: Zipkin, Prometheus, X-Ray, Stackdriver

## Try the demo app

```bash
# clone the repository
git clone https://github.com/dosco/super-graph

# run db in background
docker-compose up -d db

# see logs and wait until DB is really UP
docker-compose logs db

# setup the demo rails app & database and run it
docker-compose run rails_app rake db:create db:migrate db:seed

# run the demo
docker-compose up

# signin to the demo app (user1@demo.com / 123456)
open http://localhost:3000

# try the super graph web ui
open http://localhost:8080
```

:::note Docker?
This demo requires `docker` you can either install it using `brew` or from the
docker website [https://docs.docker.com/docker-for-mac/install/](https://docs.docker.com/docker-for-mac/install/)
:::

## GraphQL

We fully support queries and mutations. For example the below GraphQL query would fetch two products that belong to the current user where the price is greater than 10.

### GraphQL Query

```graphql
query {
  users {
    id
    email
    picture: avatar
    password
    full_name
    products(limit: 2, where: { price: { gt: 10 } }) {
      id
      name
      description
      price
    }
  }
}
```

#### JSON Result

```json
{
  "data": {
    "users": [
      {
        "id": 1,
        "email": "odilia@west.info",
        "picture": "https://robohash.org/simur.png?size=300x300",
        "full_name": "Edwin Orn",
        "products": [
          {
            "id": 16,
            "name": "Sierra Nevada Style Ale",
            "description": "Belgian Abbey, 92 IBU, 4.7%, 17.4°Blg",
            "price": 16.47
          },
          ...
        ]
      }
    ]
  }
}
```

:::note Set User ID
In development mode you can use the `X-User-ID: 4` header to set a user id so you don't have to worries about cookies etc. This can be set using the _HTTP Headers_ tab at the bottom of the web UI.
:::

### GraphQL Mutation

In another example the below GraphQL mutation would insert a product into the database. The first part of the below example is the variable data and the second half is the GraphQL mutation. For mutations data has to always ben passed as a variable.

```json
{
  "data": {
    "name": "Art of Computer Programming",
    "description": "The Art of Computer Programming (TAOCP) is a comprehensive monograph written by computer scientist Donald Knuth",
    "price": 30.5
  }
}
```

```graphql
mutation {
  product(insert: $data) {
    id
    name
  }
}
```

### Built-in GraphQL Editor

Quickly craft and test your queries with a full-featured GraphQL editor. Auto-complete and schema documentation is automatically available.

<img alt="Zipkin Traces" src={useBaseUrl("img/webui.jpg")} />
