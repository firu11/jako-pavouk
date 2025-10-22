import re
from playwright.sync_api import Page, expect


def test_ma_h1(page: Page):
    page.goto("/")
    expect(page).to_have_title(re.compile("Jako Pavouk"))


def test_zacit_psat(page: Page):
    page.goto("/")

    page.get_by_role("button", name="Začít psát").click()
    expect(page.get_by_role("heading", name="První krůčky")).to_be_visible()
