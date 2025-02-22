import contextlib
import itertools
import logging
import random
import threading
import time
import uuid

import requests

logging.basicConfig(level=logging.INFO)
logger = logging.getLogger(__name__)

BASE_URL = "http://localhost:8080"

CLIENTS_COUNT = 200
ADVERTISERS_COUNT = 200
ML_SCORES_COUNT = 800
CAMPAIGNS_COUNT = 1000
VIEW_ADS_PER_DAY_COUNT = 100
DAYS_COUNT = 20

CLIENT_GENDERS = ("MALE", "FEMALE")
MAX_AGE = 20
CLIENT_LOCATIONS = ("Moscow", "Kazan", "New York City", "London")

AGE_PAIRS = (
    (0, 7),
    (3, 10),
    (10, 18),
    (12, 18),
    (12, None),
    (15, None),
    (None, 15),
    (None, 18),
    (None, None),
)

log_block = threading.Lock()


@contextlib.contextmanager
def measure(name: str) -> None:
    start = time.perf_counter()
    yield
    seconds = time.perf_counter() - start
    logger.info(f"Request {name}: {seconds * 1000:.2f}ms")


def setup_clients() -> list[str]:
    client_ids = [str(uuid.uuid4()) for _ in range(CLIENTS_COUNT)]

    clients_bulk = [
        {
            "client_id": client_ids[i],
            "login": f"Client #{i}",
            "age": random.randint(1, MAX_AGE),
            "location": random.choice(CLIENT_LOCATIONS),
            "gender": random.choice(CLIENT_GENDERS),
        }
        for i in range(CLIENTS_COUNT)
    ]

    with measure("upsert clients"):
        resp = requests.post(f"{BASE_URL}/clients/bulk", json=clients_bulk)
    assert resp.status_code == 200

    return client_ids


def setup_advertisers() -> list[str]:
    advertiser_ids = [str(uuid.uuid4()) for _ in range(ADVERTISERS_COUNT)]

    advertisers_bulk = [
        {
            "advertiser_id": advertiser_ids[i],
            "name": f"Advertiser #{i}",
        }
        for i in range(ADVERTISERS_COUNT)
    ]

    with measure("upsert advertisers"):
        resp = requests.post(f"{BASE_URL}/advertisers/bulk", json=advertisers_bulk)
    assert resp.status_code == 200

    return advertiser_ids


def ml_scores_thread(name: str, pairs: list[tuple[str, str]]) -> None:
    times = []

    for client_id, advertiser_id in pairs:
        start = time.perf_counter()
        resp = requests.post(
            f"{BASE_URL}/ml-scores",
            json={
                "client_id": client_id,
                "advertiser_id": advertiser_id,
                "score": random.randint(0, 100),
            },
        )
        times.append((time.perf_counter() - start) * 1000)
        assert resp.status_code == 200

    with log_block:
        logger.info(
            f"Upserting ml_scores: thread {name}, executed {len(times)} requests, "
            f"{sum(times) / len(times):.2f}ms average",
        )


def setup_ml_scores(client_ids: list[str], advertiser_ids: list[str]) -> None:
    pairs = list(itertools.product(client_ids, advertiser_ids))
    random.shuffle(pairs)
    pairs = pairs[:ML_SCORES_COUNT]

    threads_count = 20
    chunk_size = len(pairs) // threads_count

    threads = [
        threading.Thread(
            target=ml_scores_thread,
            args=(f"#{i}", pairs[i * chunk_size : (i + 1) * chunk_size]),
        )
        for i in range(threads_count)
    ]

    for thread in threads:
        thread.start()

    for thread in threads:
        thread.join()


def set_date(date: int) -> None:
    with measure("set date"):
        resp = requests.post(f"{BASE_URL}/time/advance", json={"current_date": date})
    assert resp.status_code == 200


campaign_ids = []


def campaigns_thread(name: str, advertiser_ids: list[str], count: int) -> None:
    times = []

    for i in range(count):
        adv_id = random.choice(advertiser_ids)

        start_date = random.randint(1, DAYS_COUNT)
        start_age, end_age = random.choice(AGE_PAIRS)

        start = time.perf_counter()
        resp = requests.post(
            f"{BASE_URL}/advertisers/{adv_id}/campaigns",
            json={
                "ad_title": f"Campaign #{i} (thread {name})",
                "ad_text": "Hi!",
                "impressions_limit": random.randint(0, 10),
                "cost_per_impression": round(random.uniform(0, 10), 2),
                "clicks_limit": random.randint(0, 10),
                "cost_per_click": round(random.uniform(0, 10), 2),
                "start_date": start_date,
                "end_date": random.randint(start_date, DAYS_COUNT),
                "targeting": {
                    "age_from": start_age,
                    "age_to": end_age,
                    "gender": random.choice(CLIENT_GENDERS + ("ALL", None)),
                    "location": random.choice(CLIENT_LOCATIONS + (None,)),
                },
            },
        )
        times.append((time.perf_counter() - start) * 1000)

        assert resp.status_code == 200
        campaign_ids.append(resp.json()["campaign_id"])

    with log_block:
        logger.info(
            f"Creating campaigns: thread {name}, executed {len(times)} requests, "
            f"{sum(times) / len(times):.2f}ms average",
        )


def setup_campaigns(advertiser_ids: list[str]) -> None:
    threads_count = 5
    chunk_size = CAMPAIGNS_COUNT // threads_count

    threads = [
        threading.Thread(
            target=campaigns_thread, args=(f"#{i}", advertiser_ids, chunk_size)
        )
        for i in range(threads_count)
    ]

    for thread in threads:
        thread.start()

    for thread in threads:
        thread.join()


def test_stress_ads():
    requests.post(
        f"{BASE_URL}/ai/moderation/enabled", json={"enabled": False}
    ).raise_for_status()

    client_ids = setup_clients()
    advertiser_ids = setup_advertisers()
    setup_ml_scores(client_ids, advertiser_ids)
    set_date(0)
    setup_campaigns(advertiser_ids)

    for date in range(1, DAYS_COUNT):
        set_date(date)

        conversion = random.random()

        times = []

        for _ in range(VIEW_ADS_PER_DAY_COUNT):
            client_id = random.choice(client_ids)
            start = time.perf_counter()
            resp = requests.get(f"{BASE_URL}/ads?client_id={client_id}")
            times.append((time.perf_counter() - start) * 1000)
            assert resp.status_code in [200, 404]

            if resp.status_code == 200:
                ad = resp.json()
                logger.info(f"Viewed ad {ad} from {client_id}")
                if random.random() > conversion:
                    resp = requests.post(
                        f"{BASE_URL}/ads/{ad['ad_id']}/click",
                        json={"client_id": client_id},
                    )
                    assert resp.status_code == 204
                    logger.info(f"Clicked ad {ad} from {client_id}")

        with log_block:
            logger.info(
                f"Viewing on day {date}: executed {len(times)} requests, "
                f"{sum(times) / len(times):.2f}ms average",
            )


def test_stress_ads_cleanup():
    for campaign_id in campaign_ids:
        resp = requests.delete(
            f"{BASE_URL}/advertisers/{campaign_id}/campaigns/{campaign_id}",
            params={"testAdvertiserValidation": "skip"},
        )
        assert resp.status_code == 204
