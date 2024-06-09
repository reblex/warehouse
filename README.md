# Warehouse Monorepo

Blablabla

## What I have done

### Microservice over Monolith architecture

In my opinion, a monolith would have been more suitable for a project of this limited scale. However, as it is supposed to represent real world knowledge of real systems, I thought it relevant to display my understanding and knowledge of distributed systsems, especially as I could see a non linear scalability need between articles, products, and order systems in a real world scenario.

### Memory storage over DB

I created storer interfaces to allow mocking and easier switching of DBs. I chose to implement a memory storage to simplify the design and to display my knowdlege of concurrent data access in Go.

### CLI

I built the warehousing solution as a ... of APIs. However from my understanding, a functional requirement was to be able to upload articles and products directly from a file, and I therefore made CLI commands available for that. There is also a handy `ping` command to help validate system status.

### Containerization

I chose to deploy my microservices using docker. A simpler approach would have been to build and run locally using ex. a Makefile. However, as the spec said to "assume that you will have ownership of the delivery and operations of your product", I wanted to display my knowledge of deploying and running real systems using Docker.

## What I have not done

-   Any type of API authentication or authorization.
-   Any type of API documentation.
-   A complete test suite - For me, in any real product development, tests are a required part of delivery and not optional. However, in the limited timeframe of this assignment I wanted to focus on functionality and architecture.

## Requirements

-   Docker
-   Go

## Instructions

_All instructions are based on standing in the root of the repo._

### Start all services

```
docker compose up
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
            "product_name": "Dining Table",
            "amount": 1
        }
    ]
}'
```
