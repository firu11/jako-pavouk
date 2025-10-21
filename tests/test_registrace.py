from ctypes import memmove
import json
import re, threading, random
from typing import Dict
from playwright.sync_api import Page, expect
from mailtm import Email

message_received = threading.Event()
received_message: str = ""


def listener(message: str):
    print("\ngot an email!")

    global received_message
    received_message = message
    message_received.set()


def get_code(message) -> str:
    plain_msg = message["text"]  # backup plain/text email body
    return re.findall("\\d{5}", plain_msg)[0]


def fill_in_details(page: Page, email: str):
    page.fill(
        "input[type='text']",
        "".join(
            random.choice("abcdefghijklmnopqrstuvwxyz1234567890") for _ in range(12)
        ),
    )
    page.fill("input[type='email']", email)
    page.fill("input[type='password']", "12345")
    page.click("button[type='submit']")


def fill_in_verification_code(page: Page, code: str):
    page.fill("input[type='text']", code)
    page.click("button[type='submit']")


def test_registrace(page: Page):
    email = Email()
    email.register()

    _ = page.goto("/registrace")

    fill_in_details(page, email.address)

    print("\nwaiting for email")
    email.start(listener)
    if not message_received.wait(timeout=25):
        raise TimeoutError("No message received within 20 seconds")
    email.stop()

    fill_in_verification_code(page, get_code(received_message))

    expect(page).to_have_url(re.compile("/klavesnice.*"))
