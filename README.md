# Warehouse Monorepo

Warehousing and ordering solution for producs and related articles.

## What I have done

### Microservice over Monolith architecture

In my opinion, a monolith would have been more suitable for a project of this limited scale. However, as it is supposed to represent real world knowledge of real systems, I thought it relevant to display my understanding and knowledge of distributed systsems, especially as I could see a non linear scalability need between articles, products, and order systems in a real world scenario.

### Memory storage over DB

I created storer interfaces to allow mocking and easier switching of DBs. I chose to implement a memory storage to simplify the design and to display my knowdlege of concurrent data access in Go.

### CLI

I built the warehousing solution as a set of APIs. However from my understanding, a functional requirement was to be able to upload articles and products directly from a file, and I therefore made CLI commands available for that. There is also a handy `ping` command to help validate system status.

### Containerization

I chose to deploy my microservices using docker. A simpler approach would have been to build and run locally using ex. a Makefile. However, as the spec said to "assume that you will have ownership of the delivery and operations of your product", I wanted to display my knowledge of deploying and running systems using Docker.

### Patching over inserting stock

I chose to implement the upload of articles as a stock "refill", rather than directly setting the new stock, which made more sense to me. This means that new stock will be the sum of current and incomming stock. I did this also as a challenge to allow concurrent reads and writes. I also made adding products an "upsert" endpoint to allow creating and modifying products within the same endpoint.

## What I have not done

-   Any type of API authentication or authorization.
-   A complete test suite - For me, in any real product development, tests are a required part of delivery and not optional. However, in the limited timeframe of this assignment I wanted to focus on functionality and architecture. I did however write some basic e2e tests.
-   SDK for cross service communication - Could have reduced type duplication and request clutter, but I did not prioritize that.
-   ENVs - request urls and ports are hard coded. A proper implementation would read from env.

## Requirements

-   Docker
-   Go
-   Python3 (for E2E)

## Instructions

_All instructions are based on standing in the root of the repo._

### Start all services

```
docker compose up
```

### E2E

```
python3 e2e.py
```

### CLI

#### Test status of services

```
go run cmd/cli.go ping
```

#### Post articles

_Absolute or relative file path. Examples are available in ./example-files/_

```
go run cmd/cli.go articles ./example-files/articles.json
```

#### Post products

_Absolute or relative file path. Examples are available in ./example-files/_

```
go run cmd/cli.go products ./example-files/products.json
```

### Basic API endpoints and curl examples

_CLI instructions further down_

#### Post Articles

```
curl --location 'localhost:8001/api/articles' \
--header 'Content-Type: application/json' \
--data '{
  "articles": [
    {
      "art_id": "1",
      "name": "leg",
      "stock": "12"
    },
    {
      "art_id": "2",
      "name": "screw",
      "stock": "17"
    },
    {
      "art_id": "3",
      "name": "seat",
      "stock": "2"
    },
    {
      "art_id": "4",
      "name": "table top",
      "stock": "1"
    }
  ]
}
'
```

#### Post Products

```
curl --location 'localhost:8002/api/products' \
--header 'Content-Type: application/json' \
--data '{
  "products": [
    {
      "name": "Dining Chair",
      "contain_articles": [
        {
          "art_id": "1",
          "amount_of": "4"
        },
        {
          "art_id": "2",
          "amount_of": "8"
        },
        {
          "art_id": "3",
          "amount_of": "1"
        }
      ],
      "price": 39.99
    },
    {
      "name": "Dining Table",
      "contain_articles": [
        {
          "art_id": "1",
          "amount_of": "4"
        },
        {
          "art_id": "2",
          "amount_of": "8"
        },
        {
          "art_id": "4",
          "amount_of": "1"
        }
      ],
      "price": 99.99
    }
  ]
}
'
```

#### Get products with availability

```
curl --location 'localhost:8002/api/products'
```

#### Make an order

```
curl --location 'localhost:8003/api/orders' \
--header 'Content-Type: application/json' \
--data '{
    "orders": [
        {
            "product_name": "Dining Chair",
            "amount": 1
        }
    ]
}'
```
