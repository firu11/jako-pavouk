import re
import random
import pytest
from playwright.sync_api import Page, expect
from mailtm import Email
import helpers.auth_email as e
from helpers.javascript import script

student_name: str = ""
ucitel_name: str = ""
password: str = ""
kod_tridy: str = ""


@pytest.fixture(name="ucitel")
def login_as_ucitel(page: Page):
    global ucitel_name, password
    if ucitel_name != "":
        page.goto("/prihlaseni")
        page.fill("input[type='text']", ucitel_name)
        page.fill("input[type='password']", password)
        page.click("button[type='submit']")

    # vytvorim pokud jeste neexistuje
    else:
        ucitel_name = "".join(
            random.choice("abcdefghijklmnopqrstuvwxyz1234567890") for _ in range(12)
        )
        password = "12345"

        tmp_email = Email()
        tmp_email.register(ucitel_name, password)

        page.goto("/registrace")
        page.fill("input[type='text']", ucitel_name)
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

    expect(page.locator("#ucet #jmeno")).to_have_text(ucitel_name)
    yield  # callne samotný test


@pytest.fixture(name="student")
def login_as_student(page: Page):
    global student_name, password
    if student_name != "":
        page.goto("/prihlaseni")
        page.fill("input[type='text']", student_name)
        page.fill("input[type='password']", password)
        page.click("button[type='submit']")

    # vytvorim pokud jeste neexistuje
    else:
        student_name = "".join(
            random.choice("abcdefghijklmnopqrstuvwxyz1234567890") for _ in range(12)
        )
        password = "12345"

        tmp_email = Email()
        tmp_email.register(student_name, password)

        page.goto("/registrace")
        page.fill("input[type='text']", student_name)
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

    expect(page.locator("#ucet #jmeno")).to_have_text(student_name)
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

    page.click("#pridat")

    page.locator("#rocnik").select_option("5.")
    page.locator("#trida").select_option("B")
    page.locator("#skupina").select_option("1")

    page.get_by_role("button", name="Vytvořit").click()

    trida_blok = page.locator("#tridy div.blok")
    expect(trida_blok).to_have_count(1)
    expect(trida_blok).to_contain_text("5.B ￨ 1")
    expect(page.locator("#alerty > *")).to_have_count(0, timeout=200)

    trida_blok.click()

    expect(page.locator("h1")).to_contain_text("Třída: 5.B ￨ 1")
    expect(page.locator("#alerty > *")).to_have_count(0, timeout=200)

    global kod_tridy
    kod_tridy = page.locator("#kod div #wrap #obsah span").inner_text()


@pytest.mark.order(3)
def test_zapis_zaka(page: Page, student):
    page.goto("/nastaveni")

    page.get_by_role("button", name="Zapsat se").click()
    expect(page).to_have_url("/zapis")

    page.fill("form > input", kod_tridy)
    page.get_by_role("button", name="Potvrdit").click()

    page.fill("form > input#jmeno", "Novák Jiří")
    page.get_by_role("button", name="Zapsat se").click()

    expect(page.locator("#alerty > *")).to_have_count(0, timeout=200)
    expect(page).to_have_url("/trida")
    expect(page.locator("h1")).to_contain_text("Třída: 5.B ￨ 1")


@pytest.mark.order(4)
def test_vytvoreni_prace(page: Page, ucitel):
    page.goto("/skola")
    page.click("#tridy div.blok")

    page.locator("#prepinac-tabu label:has-text('Práce')").click()
    page.click("div#pridat")

    page.click("#delka > div.cas-toggle > label:nth-child(1)")
    page.select_option("#typ-textu", "Zeměpis")
    page.click("#text>div>button")

    page.get_by_role("button", name="Zadat práci").click()

    expect(page.locator("#prace-kontejner .jedna-prace")).to_have_count(1, timeout=500)
    expect(page.locator("#prace-kontejner .jedna-prace")).to_contain_text("Práce 1")


@pytest.mark.order(5)
def test_psani_prace(page: Page, student):
    page.goto("/trida")
    expect(page.locator("#alerty > *")).to_have_count(0, timeout=200)

    prace = page.locator(".prace")
    expect(prace).to_have_count(1)

    prace.click()
    expect(page.locator("#alerty > *")).to_have_count(0, timeout=200)
    expect(page.get_by_role("heading", name="Práce ve třídě"))
    assert page.locator("#text>.slovo").count() > 0

    page.evaluate(script)

    page.wait_for_selector("#tlacitka-kontainer > button", timeout=61 * 1000)
