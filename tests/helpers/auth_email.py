import re
import threading
from typing import Any

message_received = threading.Event()
received_message: dict[str, Any] = {}


def listener(message: dict):
    print("\ngot an email!")

    global received_message
    received_message = message
    message_received.set()


def get_code() -> str:
    if len(received_message) == 0:
        return ""

    plain_msg = received_message["text"]  # backup plain/text email body
    return re.findall("\\d{5}", plain_msg)[0]
