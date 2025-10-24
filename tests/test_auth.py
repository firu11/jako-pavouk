import re
import pytest
import random
from playwright.sync_api import Page, expect
from mailtm import Email
import helpers.auth_email as e

name = "".join(random.choice("abcdefghijklmnopqrstuvwxyz1234567890") for _ in range(12))
password = "12345"
new_pass = "54321"

tmp_email = Email()
tmp_email.register(name, password)


@pytest.mark.order(1)
def test_registrace(page: Page):
    page.goto("/registrace")

    fill_in_registration_details(page, tmp_email.address)

    print("\nwaiting for email")
    tmp_email.start(e.listener)
    if not e.message_received.wait(timeout=120):
        tmp_email.stop()
        e.message_received.clear()
        raise TimeoutError("No message received within 120 seconds")
    tmp_email.stop()
    e.message_received.clear()

    fill_in_verification_code(page, e.get_code())

    expect(page).to_have_url(re.compile("/klavesnice.*"))
    expect(page.locator("#alerty > *")).to_have_count(0, timeout=0.2)


@pytest.mark.order(2)
def test_login(page: Page):
    page.goto("/prihlaseni")
    fill_in_login_details(page)

    expect(page).to_have_url("/statistiky")
    expect(page.locator("#alerty > *")).to_have_count(0, timeout=0.2)
    expect(page.locator("#ucet #jmeno")).to_have_text(name)


@pytest.mark.order(3)
def test_zmenit_jmeno(page: Page):
    global password

    page.goto("/prihlaseni")
    fill_in_login_details(page)
    expect(page).to_have_url("/statistiky")

    page.locator("#ucet .kulate-tlacitko img[src*='nastaveni']").click()
    expect(page).to_have_url("/nastaveni")

    page.get_by_role("button", name="Změnit heslo").click()
    expect(page).to_have_url(re.compile("/zapomenute-heslo*"))

    page.get_by_role("button", name="Poslat e-mail").click()
    expect(page).to_have_url(re.compile("/zapomenute-heslo*"))

    print("\nwaiting for email")
    tmp_email.start(e.listener)
    if not e.message_received.wait(timeout=120):
        tmp_email.stop()
        e.message_received.clear()
        raise TimeoutError("No message received within 120 seconds")
    tmp_email.stop()
    e.message_received.clear()

    fill_in_code_and_new_password(page, e.get_code(), new_pass)

    expect(page.locator("form > h3")).to_have_count(2)

    page.goto("/prihlaseni")
    fill_in_login_details(page, new_pass)
    expect(page.locator("#alerty > *")).to_have_count(0, timeout=0.2)


@pytest.mark.order(4)
def test_smazat_ucet(page: Page):
    page.goto("/prihlaseni")
    fill_in_login_details(page, new_pass)
    expect(page).to_have_url("/statistiky")

    page.locator("#ucet .kulate-tlacitko img[src*='nastaveni']").click()
    expect(page).to_have_url("/nastaveni")

    page.get_by_role("button", name="Smazat účet").click()
    page.get_by_role("button", name="Opravdu smazat").click()
    expect(page).to_have_url("/prihlaseni")

    fill_in_login_details(page, new_pass)
    expect(page.locator("#alerty > *")).to_have_count(1, timeout=0.2)


def fill_in_verification_code(page: Page, code: str):
    page.fill("input[type='text']", code)
    page.click("button[type='submit']")


def fill_in_code_and_new_password(page: Page, code: str, new_password: str):
    page.fill("input[type='text']", code)
    page.fill("input[type='password']", new_password)
    page.click("button[type='submit']")


def fill_in_registration_details(page: Page, email: str):
    page.fill("input[type='text']", name)
    page.fill("input[type='email']", email)
    page.fill("input[type='password']", password)
    page.click("button[type='submit']")


def fill_in_login_details(page: Page, passw: str = password):
    page.fill("input[type='text']", name)
    page.fill("input[type='password']", passw)
    page.click("button[type='submit']")
