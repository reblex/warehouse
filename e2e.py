#!/usr/bin/python3

import http.client
import json


class TestFailedException(Exception):
    pass


def main():
    test(ping_article_service)
    test(ping_product_service)
    test(ping_order_service)
    test(post_articles)
    test(get_articles)
    test(reserve_articles)
    test(get_article_availability)
    test(post_products)
    test(get_products)
    test(reserve_products)
    test(post_order)


def test(func):
    try:
        func()
    except TestFailedException as e:
        print("x", func.__name__, "- Test failed:", e)
        return 1
    except Exception as e:
        print("x", func.__name__, "- Test failed, unexpected exception:", e)
        return 1

    print("\u2713", func.__name__)


def ping_article_service():
    conn = http.client.HTTPConnection("localhost", 8001)
    payload = ''
    headers = {}
    conn.request("GET", "/api/ping", payload, headers)
    res = conn.getresponse()

    if res.status != 200:
        raise TestFailedException("Non response from article-service")


def ping_product_service():
    conn = http.client.HTTPConnection("localhost", 8002)
    payload = ''
    headers = {}
    conn.request("GET", "/api/ping", payload, headers)
    res = conn.getresponse()

    if res.status != 200:
        raise TestFailedException("Non response from product-service")


def ping_order_service():
    conn = http.client.HTTPConnection("localhost", 8003)
    payload = ''
    headers = {}
    conn.request("GET", "/api/ping", payload, headers)
    res = conn.getresponse()

    if res.status != 200:
        raise TestFailedException("Non response from order-service")


def post_articles():
    conn = http.client.HTTPConnection("localhost", 8001)
    payload = json.dumps({
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
    })
    headers = {
        'Content-Type': 'application/json'
    }
    conn.request("POST", "/api/articles", payload, headers)
    res = conn.getresponse()

    if res.status != 200:
        raise TestFailedException("Non success response")


def get_articles():
    conn = http.client.HTTPConnection("localhost", 8001)
    conn.request("GET", "/api/articles")
    res = conn.getresponse()
    data = json.loads(res.read().decode())

    if res.status != 200:
        raise TestFailedException("Non success response")

    if len(data["articles"]) < 4:
        raise TestFailedException("Incorrect count of articles")


def reserve_articles():
    conn = http.client.HTTPConnection("localhost", 8001)
    payload = json.dumps({
        "reservations": [
            {
                "id": "1",
                "count": "2"
            }
        ]
    })
    headers = {
        'Content-Type': 'application/json'
    }
    conn.request("POST", "/api/articles/reserve", payload, headers)
    res = conn.getresponse()

    if res.status != 200:
        raise TestFailedException("Non success response")


def get_article_availability():
    conn = http.client.HTTPConnection("localhost", 8001)
    payload = json.dumps({
        "reservations": [
            {
                "id": "1",
                "Count": "4"
            },
            {
                "id": "2",
                "Count": "8"
            },
            {
                "id": "4",
                "Count": "1"
            }
        ]
    })
    headers = {
        'Content-Type': 'application/json'
    }
    conn.request("POST", "/api/articles/availability", payload, headers)
    res = conn.getresponse()
    data = json.loads(res.read().decode())

    if res.status != 200:
        raise TestFailedException("Non success response")

    if data["availability"] < 1:
        raise TestFailedException("Incorrect availability")


def post_products():
    conn = http.client.HTTPConnection("localhost", 8002)
    payload = json.dumps({
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
    })
    headers = {
        'Content-Type': 'application/json'
    }
    conn.request("POST", "/api/products", payload, headers)
    res = conn.getresponse()

    if res.status != 200:
        raise TestFailedException("Non success response")


def get_products():
    conn = http.client.HTTPConnection("localhost", 8002)
    payload = ''
    headers = {}
    conn.request("GET", "/api/products", payload, headers)
    res = conn.getresponse()
    data = json.loads(res.read().decode())

    if res.status != 200:
        raise TestFailedException("Non success response")

    if len(data["products"]) < 2:
        raise TestFailedException("Incorrect product count")


def reserve_products():
    conn = http.client.HTTPConnection("localhost", 8002)
    payload = json.dumps({
        "orders": [
            {
                "product_name": "Dining Table",
                "amount": 1
            }
        ]
    })
    headers = {
        'Content-Type': 'application/json'
    }
    conn.request("POST", "/api/products/reserve", payload, headers)
    res = conn.getresponse()

    if res.status != 200:
        raise TestFailedException("Non success response")


def post_order():
    conn = http.client.HTTPConnection("localhost", 8003)
    payload = json.dumps({
        "orders": [
            {
                "product_name": "Dining Chair",
                "amount": 1
            }
        ]
    })
    headers = {
        'Content-Type': 'application/json'
    }
    conn.request("POST", "/api/orders", payload, headers)
    res = conn.getresponse()

    if res.status != 200:
        raise TestFailedException("Non success response")


if __name__ == "__main__":
    main()
