import argparse
import random
import uuid
from http import HTTPStatus

import requests

BASE_URL = "http://localhost:8080"
CLIENTS_COUNT = 100
AD_VIEWS_COUNT = 100
CLICK_PROBABILITY = 0.5


def create_clients() -> list[str]:
    client_ids = [str(uuid.uuid4()) for _ in range(CLIENTS_COUNT)]
    clients_bulk = [
        {
            "client_id": client_ids[i],
            "login": f"Client #{i}",
            "age": random.randint(1, 100),
            "location": random.choice(["Moscow", "Kazan", "New York City", "London"]),
            "gender": random.choice(["MALE", "FEMALE"]),
        }
        for i in range(CLIENTS_COUNT)
    ]

    resp = requests.post(f"{BASE_URL}/clients/bulk", json=clients_bulk)
    assert resp.status_code == HTTPStatus.OK

    return client_ids


def add_ml_scores(client_ids: list[str], advertiser_id: str) -> None:
    for client_id in client_ids:
        resp = requests.post(
            f"{BASE_URL}/ml-scores",
            json={
                "client_id": client_id,
                "advertiser_id": advertiser_id,
                "score": random.randint(100_000, 200_000),
            },
        )
        assert resp.status_code == HTTPStatus.OK


def view_click_ads(client_ids: list[str]) -> None:
    for client_id in client_ids:
        resp = requests.get(f"{BASE_URL}/ads?client_id={client_id}")
        assert resp.status_code in [HTTPStatus.OK, HTTPStatus.NOT_FOUND]
        if resp.status_code == HTTPStatus.NOT_FOUND:
            print(f"Received 404 on {client_id}")
            continue

        ad = resp.json()
        print(f"Viewed ad {ad} from {client_id}")

        if random.random() > CLICK_PROBABILITY:
            click_resp = requests.post(
                f"{BASE_URL}/ads/{ad['ad_id']}/click", json={"client_id": client_id}
            )
            assert click_resp.status_code == HTTPStatus.NO_CONTENT
            print(f"Clicked ad {ad} from {client_id}")


def main() -> None:
    advertiser_id = input("Enter target advertiser_id (must be UUID): ")

    client_ids = create_clients()
    add_ml_scores(client_ids, advertiser_id)
    view_click_ads(client_ids)


if __name__ == "__main__":
    main()
