import os
from typing import Tuple
from dotenv import load_dotenv

load_dotenv()


def create_ucitel() -> Tuple[str, str]:
    name = os.getenv("UCITEL_JMENO")
    password = os.getenv("UCITEL_HESLO")
    if name is None or password is None:
        raise EnvironmentError("UCITEL_JMENO or UCITEL_HESLO variable doesnt exist")

    return name, password


def delete_ucitel():
    pass
