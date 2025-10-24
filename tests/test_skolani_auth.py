import re
import random
import pytest
from playwright.sync_api import Page, expect
from mailtm import Email
import helpers.auth_email as e

name: str = ""
password: str = ""


@pytest.fixture(name="ucitel")
def login_as_ucitel(page: Page):
    global name, password
    if name != "":
        page.goto("/prihlaseni")
        page.fill("input[type='text']", name)
        page.fill("input[type='password']", password)
        page.click("button[type='submit']")

    # vytvorim pokud jeste neexistuje
    else:
        name = "".join(
            random.choice("abcdefghijklmnopqrstuvwxyz1234567890") for _ in range(12)
        )
        password = "12345"

        tmp_email = Email()
        tmp_email.register(name, password)

        page.goto("/registrace")
        page.fill("input[type='text']", name)
        page.fill("input[type='email']", tmp_email.address)
        page.fill("input[type='password']", password)
        page.click("button[type='submit']")

        print("\nwaiting for email")
        tmp_email.start(e.listener)
        if not e.message_received.wait(timeout=120):
            tmp_email.stop()
            e.message_received.clear()
            raise TimeoutError("No message received within 120 seconds")
        tmp_email.stop()
        e.message_received.clear()

        page.fill("input[type='text']", e.get_code())
        page.click("button[type='submit']")

    expect(page.locator("#ucet #jmeno")).to_have_text(name)
    yield  # callne samotný test


@pytest.mark.order(1)
def test_registrace_skoly(page: Page, ucitel):
    testovaci_skola_name = "Testovací škola"

    page.goto("/skolni-system")
    expect(page).to_have_title(re.compile("Systém pro školy"))

    page.get_by_label("Jméno školy").fill(testovaci_skola_name)
    page.get_by_label("Telefonní číslo").fill("+420123456789")
    page.click("button[type='submit']")

    expect(page.locator("#formular")).to_contain_text("Díky za registraci školy!")

    page.goto("/skola")

    expect(page.get_by_role("heading", name=testovaci_skola_name))


@pytest.mark.order(2)
def test_vytvoreni_tridy(page: Page, ucitel):
    page.goto("/skola")

    page.locator("#pridat").click()

    page.locator("#rocnik").select_option("5.")
    page.locator("#trida").select_option("B")
    page.locator("#skupina").select_option("1")

    page.get_by_role("button", name="Vytvořit").click()

    trida_blok = page.locator("#tridy div.blok")
    expect(trida_blok).to_have_count(1)
    expect(trida_blok).to_contain_text("5.B ￨ 1")
    expect(page.locator("#alerty > *")).to_have_count(0, timeout=0.2)

    trida_blok.click()

    expect(page.locator("h1")).to_contain_text("Třída: 5.B ￨ 1")
    expect(page.locator("#alerty > *")).to_have_count(0, timeout=0.2)
