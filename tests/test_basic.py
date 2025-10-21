import re
from playwright.sync_api import Page, expect


def test_has_title(page: Page):
    _ = page.goto("/")

    # Expect a title "to contain" a substring.
    expect(page).to_have_title(re.compile("Jako Pavouk"))


def test_get_started_link(page: Page):
    _ = page.goto("/")

    # Click the get started link.
    page.get_by_role("button", name="Začít psát").click()

    # Expects page to have a heading with the name of Installation.
    expect(page.get_by_role("heading", name="První krůčky")).to_be_visible()
