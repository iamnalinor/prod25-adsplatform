import uuid


def uuid_from_id(user_id: int) -> uuid.UUID:
    hex_id = "0" * (32 - len(str(user_id))) + str(user_id)
    return uuid.UUID(hex_id, version=4)
